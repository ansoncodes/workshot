package capture

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ansoncodes/workshot/pkg/types"
)

// terminalcapturer captures recent terminal commands
type TerminalCapturer struct {
	maxCommands int
}

// newterminalcapturer creates a terminal capturer
func NewTerminalCapturer() types.Capturer {
	return &TerminalCapturer{
		maxCommands: 20, // capture last 20 commands
	}
}

func (t *TerminalCapturer) Name() string {
	return "terminal"
}

func (t *TerminalCapturer) Priority() int {
	return 30 // runs later because it is only extra info
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
	// terminal history cannot be restored
	return nil
}

func (t *TerminalCapturer) CanRestore(data map[string]interface{}) bool {
	return false // terminal data is read only
}

// read recent commands from shell history
func (t *TerminalCapturer) getRecentCommands() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	// common history files to check
	historyFiles := []string{
		filepath.Join(home, ".zsh_history"),
		filepath.Join(home, ".bash_history"),
		filepath.Join(home, ".history"),
	}

	for _, histFile := range historyFiles {
		if commands := t.readHistoryFile(histFile); len(commands) > 0 {
			return commands
		}
	}

	return []string{}
}

// read and parse a history file
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

	// keep only last n commands
	if len(lines) > t.maxCommands {
		lines = lines[len(lines)-t.maxCommands:]
	}

	return lines
}

// clean shell specific formatting
func (t *TerminalCapturer) cleanHistoryLine(line string) string {
	// zsh history format example
	if strings.HasPrefix(line, ":") {
		parts := strings.SplitN(line, ";", 2)
		if len(parts) == 2 {
			line = parts[1]
		}
	}

	return strings.TrimSpace(line)
}

// check if command may contain sensitive data
func (t *TerminalCapturer) isSensitive(line string) bool {
	// patterns that usually mean secrets
	sensitivePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(password|passwd|pwd)\s*=`),
		regexp.MustCompile(`(?i)(secret|token|key|api[-_]?key)\s*=`),
		regexp.MustCompile(`(?i)export\s+(AWS|GITHUB|GITLAB|API|AUTH)_`),
		regexp.MustCompile(`(?i)(bearer|auth(orization)?)\s+[a-zA-Z0-9+/=]{20,}`),
		regexp.MustCompile(`(?i)curl.*(-H|--header).*authorization`),
		regexp.MustCompile(`(?i)(ssh|scp|rsync).*password`),
	}

	for _, pattern := range sensitivePatterns {
		if pattern.MatchString(line) {
			return true
		}
	}

	return false
}
