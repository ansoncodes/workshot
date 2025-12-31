package types

import "time"

const (
	// SchemaVersion represents the current snapshot format version.
	// This should be incremented ONLY when a breaking change is made
	// to the snapshot structure.
	SchemaVersion = 1
)

// Snapshot represents a saved development context at a point in time.
// It contains both core metadata and extensible plugin-captured data.
type Snapshot struct {
	SchemaVersion int       `json:"schema_version"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	WorkingDir    string    `json:"working_dir"`

	// Core Git-related fields extracted for quick access and indexing.
	// These are commonly used values and should not require plugin parsing.
	GitBranch string `json:"git_branch,omitempty"`
	GitRemote string `json:"git_remote,omitempty"`
	GitDirty  bool   `json:"git_dirty,omitempty"`

	// PluginData stores data captured by individual plugins.
	// The key is the plugin name, and the value is plugin-specific data.
	// This allows the snapshot format to remain flexible and extensible.
	PluginData map[string]interface{} `json:"plugin_data,omitempty"`
}

// Capturer defines the contract that all capture plugins must follow.
// This interface enables a clean, extensible plugin architecture.
type Capturer interface {
	// Name returns a unique identifier for the capturer.
	// This is used as the key inside Snapshot.PluginData.
	Name() string

	// Priority determines execution order.
	// Lower values run earlier, allowing dependency-based ordering.
	Priority() int

	// Capture collects relevant data from the current environment.
	// It should return (nil, nil) if there is nothing to capture.
	Capture() (map[string]interface{}, error)

	// Restore applies previously captured data back to the environment.
	Restore(data map[string]interface{}) error

	// CanRestore determines whether the capturer can safely restore
	// the provided data. This enables graceful handling of
	// missing, incompatible, or partial snapshot data.
	CanRestore(data map[string]interface{}) bool
}

// NewSnapshot creates a new Snapshot instance with sensible defaults.
// PluginData is always initialized to avoid nil map checks later.
func NewSnapshot(name string) *Snapshot {
	return &Snapshot{
		SchemaVersion: SchemaVersion,
		Name:          name,
		CreatedAt:     time.Now(),
		PluginData:    make(map[string]interface{}),
	}
}
