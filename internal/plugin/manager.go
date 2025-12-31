package plugin

import (
	"fmt"
	"sort"

	"github.com/ansoncodes/workshot/pkg/types"
)

// manager controls all capture plugins
// it registers plugins and runs them
type Manager struct {
	capturers []types.Capturer
}

// newmanager creates a plugin manager
func NewManager() *Manager {
	return &Manager{
		capturers: make([]types.Capturer, 0),
	}
}

// register adds a plugin to the manager
// plugins are added during startup
func (m *Manager) Register(c types.Capturer) {
	m.capturers = append(m.capturers, c)
}

// captureall runs all plugins and collects data
// some plugins can fail without stopping others
func (m *Manager) CaptureAll() (map[string]interface{}, error) {
	// run plugins based on priority
	sort.Slice(m.capturers, func(i, j int) bool {
		return m.capturers[i].Priority() < m.capturers[j].Priority()
	})

	pluginData := make(map[string]interface{})
	var errors []error

	for _, capturer := range m.capturers {
		data, err := capturer.Capture()
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", capturer.Name(), err))
			continue
		}

		// save data only if something was captured
		if data != nil && len(data) > 0 {
			pluginData[capturer.Name()] = data
		}
	}

	// fail only if everything failed
	if len(errors) > 0 && len(pluginData) == 0 {
		return nil, fmt.Errorf("all capture plugins failed: %v", errors)
	}

	return pluginData, nil
}

// restoreall restores data using plugins
// errors are collected but restore continues
func (m *Manager) RestoreAll(pluginData map[string]interface{}) []error {
	var errors []error

	for _, capturer := range m.capturers {
		data, exists := pluginData[capturer.Name()]
		if !exists {
			continue
		}

		// make sure data is in correct format
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			errors = append(errors, fmt.Errorf("%s: invalid data format", capturer.Name()))
			continue
		}

		// check if plugin can restore this data
		if !capturer.CanRestore(dataMap) {
			continue
		}

		// try restoring the data
		if err := capturer.Restore(dataMap); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", capturer.Name(), err))
		}
	}

	return errors
}

// listcapturers returns names of all plugins
func (m *Manager) ListCapturers() []string {
	names := make([]string, len(m.capturers))
	for i, c := range m.capturers {
		names[i] = c.Name()
	}
	return names
}
