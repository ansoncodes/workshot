package plugin

import (
	"fmt"
	"testing"
)

// Mock capturer for testing - implements the Capturer interface
type mockCapturer struct {
	name         string
	priority     int
	captureData  map[string]interface{}
	captureError error
	canRestore   bool
	restoreError error
}

func (m *mockCapturer) Name() string {
	return m.name
}

func (m *mockCapturer) Priority() int {
	return m.priority
}

func (m *mockCapturer) Capture() (map[string]interface{}, error) {
	return m.captureData, m.captureError
}

func (m *mockCapturer) Restore(data map[string]interface{}) error {
	return m.restoreError
}

func (m *mockCapturer) CanRestore(data map[string]interface{}) bool {
	return m.canRestore
}

func TestManagerRegister(t *testing.T) {
	manager := NewManager()
	
	mock1 := &mockCapturer{name: "mock1", priority: 10}
	mock2 := &mockCapturer{name: "mock2", priority: 20}
	
	manager.Register(mock1)
	manager.Register(mock2)
	
	names := manager.ListCapturers()
	if len(names) != 2 {
		t.Errorf("Expected 2 capturers, got %d", len(names))
	}
}

func TestManagerCaptureAll(t *testing.T) {
	manager := NewManager()
	
	mock1 := &mockCapturer{
		name:     "mock1",
		priority: 10,
		captureData: map[string]interface{}{
			"key1": "value1",
		},
	}
	
	mock2 := &mockCapturer{
		name:     "mock2",
		priority: 20,
		captureData: map[string]interface{}{
			"key2": "value2",
		},
	}
	
	manager.Register(mock1)
	manager.Register(mock2)
	
	data, err := manager.CaptureAll()
	if err != nil {
		t.Fatalf("CaptureAll failed: %v", err)
	}
	
	if len(data) != 2 {
		t.Errorf("Expected 2 plugin data entries, got %d", len(data))
	}
	
	if _, ok := data["mock1"]; !ok {
		t.Error("Missing mock1 data")
	}
	
	if _, ok := data["mock2"]; !ok {
		t.Error("Missing mock2 data")
	}
}

func TestManagerCaptureAllWithError(t *testing.T) {
	manager := NewManager()
	
	mock1 := &mockCapturer{
		name:         "mock1",
		priority:     10,
		captureError: fmt.Errorf("capture failed"),
	}
	
	mock2 := &mockCapturer{
		name:     "mock2",
		priority: 20,
		captureData: map[string]interface{}{
			"key2": "value2",
		},
	}
	
	manager.Register(mock1)
	manager.Register(mock2)
	
	data, err := manager.CaptureAll()
	
	if len(data) != 1 {
		t.Errorf("Expected 1 successful capture, got %d", len(data))
	}
	
	if err != nil {
		t.Errorf("Expected nil error for partial success, got: %v", err)
	}
}

func TestManagerRestoreAll(t *testing.T) {
	manager := NewManager()
	
	mock1 := &mockCapturer{
		name:       "mock1",
		priority:   10,
		canRestore: true,
	}
	
	mock2 := &mockCapturer{
		name:       "mock2",
		priority:   20,
		canRestore: false,
	}
	
	manager.Register(mock1)
	manager.Register(mock2)
	
	pluginData := map[string]interface{}{
		"mock1": map[string]interface{}{"key": "value"},
		"mock2": map[string]interface{}{"key": "value"},
	}
	
	errors := manager.RestoreAll(pluginData)
	
	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errors), errors)
	}
}