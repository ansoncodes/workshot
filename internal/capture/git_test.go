package capture

import (
	"testing"
)

func TestGitCapturerName(t *testing.T) {
	gc := NewGitCapturer()
	
	if gc.Name() != "git" {
		t.Errorf("Expected name 'git', got %s", gc.Name())
	}
}

func TestGitCapturerPriority(t *testing.T) {
	gc := NewGitCapturer()
	
	if gc.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", gc.Priority())
	}
}

func TestGitCanRestore(t *testing.T) {
	gc := NewGitCapturer()
	
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected bool
	}{
		{
			"with branch",
			map[string]interface{}{"branch": "main"},
			true,
		},
		{
			"without branch",
			map[string]interface{}{"commit": "abc123"},
			false,
		},
		{
			"empty data",
			map[string]interface{}{},
			false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gc.CanRestore(tt.data)
			if result != tt.expected {
				t.Errorf("CanRestore(%v) = %v, want %v", tt.data, result, tt.expected)
			}
		})
	}
}