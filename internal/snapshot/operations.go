package snapshot

import (
	"fmt"
	"os"

	"github.com/ansoncodes/workshot/internal/plugin"
	"github.com/ansoncodes/workshot/internal/storage"
	"github.com/ansoncodes/workshot/pkg/types"
)

// freeze saves the current work context
func Freeze(name string, manager *plugin.Manager) error {
	// create a new snapshot
	snap := types.NewSnapshot(name)

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	snap.WorkingDir = cwd

	// run all capture plugins
	pluginData, err := manager.CaptureAll()
	if err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}
	snap.PluginData = pluginData

	// copy git data to top level fields
	if gitData, ok := pluginData["git"].(map[string]interface{}); ok {
		if branch, ok := gitData["branch"].(string); ok {
			snap.GitBranch = branch
		}
		if remote, ok := gitData["remote"].(string); ok {
			snap.GitRemote = remote
		}
		if dirty, ok := gitData["dirty"].(bool); ok {
			snap.GitDirty = dirty
		}
	}

	// create storage handler
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// check if snapshot name already exists
	if store.Exists(name) {
		return fmt.Errorf("workshot '%s' already exists (use 'workshot delete %s' first)", name, name)
	}

	// save snapshot to disk
	if err := store.Save(snap); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	return nil
}

// restore loads a saved snapshot and applies it
func Restore(name string, manager *plugin.Manager) (*types.Snapshot, []error) {
	// create storage handler
	store, err := storage.New()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to initialize storage: %w", err)}
	}

	// load snapshot from disk
	snap, err := store.Load(name)
	if err != nil {
		return nil, []error{err}
	}

	var errors []error

	// move to saved working directory
	if snap.WorkingDir != "" {
		if err := os.Chdir(snap.WorkingDir); err != nil {
			errors = append(errors, fmt.Errorf("failed to change directory: %w", err))
		}
	}

	// run restore on all plugins
	restoreErrors := manager.RestoreAll(snap.PluginData)
	errors = append(errors, restoreErrors...)

	return snap, errors
}
