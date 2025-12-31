package capture

import (
	"testing"
)

func TestTerminalSensitiveFiltering(t *testing.T) {
	tc := NewTerminalCapturer().(*TerminalCapturer)

	tests := []struct {
		command   string
		sensitive bool
		reason    string
	}{
		{"ls -la", false, "basic command"},
		{"git status", false, "git command"},
		{"npm install", false, "package manager"},
		{"go test ./...", false, "test command"},
		{"python script.py", false, "python command"},
		
		{"export PASSWORD=secret123", true, "password env var"},
		{"echo password=test", true, "password in echo"},
		{"set PASSWORD=secret", true, "Windows set password"},
		
		{"export AWS_SECRET_ACCESS_KEY=xyz", true, "AWS secret"},
		{"export GITHUB_TOKEN=ghp_123", true, "GitHub token"},
		{"export API_KEY=abc123", true, "API key"},
		{"set API_KEY=secret", true, "Windows API key"},
		
		{"curl -H 'Authorization: Bearer token123'", true, "bearer token"},
		{"curl --header 'authorization: Basic xyz'", true, "basic auth"},
		
		{"ssh user@server", false, "basic ssh"},
		{"export PATH=/usr/bin:$PATH", false, "PATH export"},
		{"export NODE_ENV=production", false, "NODE_ENV export"},
	}

	for _, tt := range tests {
		result := tc.isSensitive(tt.command)
		if result != tt.sensitive {
			t.Errorf("isSensitive(%q) = %v, want %v (reason: %s)", 
				tt.command, result, tt.sensitive, tt.reason)
		}
	}
}

func TestTerminalCapturerName(t *testing.T) {
	tc := NewTerminalCapturer()
	
	if tc.Name() != "terminal" {
		t.Errorf("Expected name 'terminal', got %s", tc.Name())
	}
}

func TestTerminalCapturerPriority(t *testing.T) {
	tc := NewTerminalCapturer()
	
	if tc.Priority() != 30 {
		t.Errorf("Expected priority 30, got %d", tc.Priority())
	}
}

func TestTerminalCanRestore(t *testing.T) {
	tc := NewTerminalCapturer()
	
	data := map[string]interface{}{
		"recent_commands": []string{"ls", "pwd"},
	}
	
	if tc.CanRestore(data) {
		t.Error("Terminal capturer should not be able to restore")
	}
}

func TestTerminalCleanHistoryLine(t *testing.T) {
	tc := NewTerminalCapturer().(*TerminalCapturer)
	
	tests := []struct {
		input    string
		expected string
	}{
		{": 1640000000:0;git status", "git status"},
		{": 1640000000:0;ls -la", "ls -la"},
		{"git status", "git status"},
		{"  git status  ", "git status"},
		{":;command", "command"},
	}
	
	for _, tt := range tests {
		result := tc.cleanHistoryLine(tt.input)
		if result != tt.expected {
			t.Errorf("cleanHistoryLine(%q) = %q, want %q", 
				tt.input, result, tt.expected)
		}
	}
}