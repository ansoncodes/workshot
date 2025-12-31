package capture

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ansoncodes/workshot/pkg/types"
)

// gitcapturer captures and restores git state
type GitCapturer struct{}

// newgitcapturer creates a git capturer
func NewGitCapturer() types.Capturer {
	return &GitCapturer{}
}

func (g *GitCapturer) Name() string {
	return "git"
}

func (g *GitCapturer) Priority() int {
	return 10 // git runs early because it is important
}

func (g *GitCapturer) Capture() (map[string]interface{}, error) {
	// check if current folder is a git repo
	if !isGitRepo() {
		return nil, nil // nothing to capture if not a git repo
	}

	data := make(map[string]interface{})

	// save current branch
	if branch := getGitBranch(); branch != "" {
		data["branch"] = branch
	}

	// save remote url
	if remote := getGitRemote(); remote != "" {
		data["remote"] = remote
	}

	// check if repo has uncommitted changes
	data["dirty"] = isGitDirty()

	// save current commit hash
	if commit := getGitCommit(); commit != "" {
		data["commit"] = commit
	}

	// save stash count if any
	if stashCount := getGitStashCount(); stashCount > 0 {
		data["stash_count"] = stashCount
	}

	return data, nil
}

func (g *GitCapturer) Restore(data map[string]interface{}) error {
	branch, ok := data["branch"].(string)
	if !ok || branch == "" {
		return nil
	}

	// check if already on the same branch
	currentBranch := getGitBranch()
	if currentBranch == branch {
		return nil // already correct
	}

	// switch to saved branch
	cmd := exec.Command("git", "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch '%s': %s", branch, string(output))
	}

	return nil
}

func (g *GitCapturer) CanRestore(data map[string]interface{}) bool {
	_, hasBranch := data["branch"]
	return hasBranch && isGitRepo()
}

// helper functions

// check if current folder is a git repo
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// get current git branch name
func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// get origin remote url
func getGitRemote() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// check if repo has uncommitted changes
func isGitDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// get current commit hash
func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	commit := strings.TrimSpace(string(output))
	if len(commit) > 7 {
		return commit[:7] // use short hash
	}
	return commit
}

// get number of stashed changes
func getGitStashCount() int {
	cmd := exec.Command("git", "stash", "list")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0
	}
	return len(lines)
}
