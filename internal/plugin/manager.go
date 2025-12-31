package plugin

import (
    "fmt"
    "sort"

    "github.com/ansoncodes/workshot/pkg/types"
)

type Manager struct {
    capturers []types.Capturer
}

func NewManager() *Manager {
    return &Manager{
        capturers: make([]types.Capturer, 0),
    }
}

func (m *Manager) Register(c types.Capturer) {
    m.capturers = append(m.capturers, c)
}

func (m *Manager) CaptureAll() (map[string]interface{}, error) {
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

        if data != nil && len(data) > 0 {
            pluginData[capturer.Name()] = data
        }
    }

    if len(errors) > 0 && len(pluginData) == 0 {
        return nil, fmt.Errorf("all capture plugins failed: %v", errors)
    }

    return pluginData, nil
}

func (m *Manager) RestoreAll(pluginData map[string]interface{}) []error {
    var errors []error

    for _, capturer := range m.capturers {
        data, exists := pluginData[capturer.Name()]
        if !exists {
            continue
        }

        dataMap, ok := data.(map[string]interface{})
        if !ok {
            errors = append(errors, fmt.Errorf("%s: invalid data format", capturer.Name()))
            continue
        }

        if !capturer.CanRestore(dataMap) {
            continue
        }

        if err := capturer.Restore(dataMap); err != nil {
            errors = append(errors, fmt.Errorf("%s: %w", capturer.Name(), err))
        }
    }

    return errors
}

func (m *Manager) ListCapturers() []string {
    names := make([]string, len(m.capturers))
    for i, c := range m.capturers {
        names[i] = c.Name()
    }
    return names
}
