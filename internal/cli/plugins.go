package cli

import (
	"github.com/ansoncodes/workshot/internal/capture"
	"github.com/ansoncodes/workshot/internal/plugin"
)

// create plugin manager and register plugins
func initPluginManager() *plugin.Manager {
	manager := plugin.NewManager()

	// register all plugins
	// lower priority runs first
	manager.Register(capture.NewGitCapturer())       // priority 10
	manager.Register(capture.NewTerminalCapturer())  // priority 30

	return manager
}
