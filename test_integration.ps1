package capture

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/ansoncodes/workshot/pkg/types"
)

type TerminalCapturer struct {
	maxCommands int
}

func NewTerminalCapturer() types.Capturer {
	return &TerminalCapturer{
		maxCommands: 20,
	}
}

func (t *TerminalCapturer) Name() string {
	return "terminal"
}

func (t *TerminalCapturer) Priority() int {
	return 30
}

func (t *TerminalCapturer) Capture() (map[string]interface{}, error) {
	commands := t.getRecentCommands()
	if len(commands) == 0 {
		return nil, nil
	}

	data := make(map[string]interface{})
	data["recent_commands"] = commands

	return data, nil
}

func (t *TerminalCapturer) Restore(data map[string]interface{}) error {
	return nil
}

func (t *TerminalCapturer) CanRestore(data map[string]interface{}) bool {
	return false
}

func (t *TerminalCapturer) getRecentCommands() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	var historyFiles []string

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		
		// PowerShell history FIRST
		historyFiles = []string{
			filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt"),
			filepath.Join(home, ".bash_history"),
			filepath.Join(home, ".history"),
		}
	} else {
		historyFiles = []string{
			filepath.Join(home, ".zsh_history"),
			filepath.Join(home, ".bash_history"),
			filepath.Join(home, ".history"),
		}
	}

	for _, histFile := range historyFiles {
		if commands := t.readHistoryFile(histFile); len(commands) > 0 {
			return commands
		}
	}

	return []string{}
}

func (t *TerminalCapturer) readHistoryFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := t.cleanHistoryLine(scanner.Text())
		if line != "" && !t.isSensitive(line) {
			lines = append(lines, line)
		}
	}

	// Return last N lines
	if len(lines) > t.maxCommands {
		lines = lines[len(lines)-t.maxCommands:]
	}

	return lines
}

func (t *TerminalCapturer) cleanHistoryLine(line string) string {
	// Zsh history format: ": timestamp:0;command"
	if strings.HasPrefix(line, ":") {
		parts := strings.SplitN(line, ";", 2)
		if len(parts) == 2 {
			line = parts[1]
		}
	}

	// Remove PowerShell line continuation backticks
	line = strings.TrimSuffix(line, "`")

	return strings.TrimSpace(line)
}

func (t *TerminalCapturer) isSensitive(line string) bool {
	sensitivePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(password|passwd|pwd)\s*=`),
		regexp.MustCompile(`(?i)(secret|token|key|api[-_]?key)\s*=`),
		regexp.MustCompile(`(?i)export\s+(AWS|GITHUB|GITLAB|API|AUTH)_`),
		regexp.MustCompile(`(?i)\$env:[A-Z_]*?(PASSWORD|SECRET|TOKEN|KEY)`),
		regexp.MustCompile(`(?i)(bearer|auth(orization)?)\s+[a-zA-Z0-9+/=]{20,}`),
		regexp.MustCompile(`(?i)curl.*(-H|--header).*authorization`),
		regexp.MustCompile(`(?i)(ssh|scp|rsync).*password`),
		regexp.MustCompile(`(?i)(postgres|mysql|mongodb|redis):\/\/.*:.*@`),
		regexp.MustCompile(`(?i)SECRET_KEY\s*=`),
		regexp.MustCompile(`(?i)DJANGO_(SECRET|DB_PASSWORD)`),
	}

	for _, pattern := range sensitivePatterns {
		if pattern.MatchString(line) {
			return true
		}
	}

	return false
}