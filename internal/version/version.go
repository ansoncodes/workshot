package version

import (
	"runtime/debug"
	"strings"
)

var (
	// Version is set by ldflags during release builds
	Version = "dev"
)

// Get returns the actual version, checking multiple sources
func Get() string {
	// If Version was set by ldflags (GoReleaser), use it
	if Version != "" && Version != "dev" {
		return Version
	}

	// Try to get version from Go module info (go install @vX.X.X)
	if info, ok := debug.ReadBuildInfo(); ok {
		// Check if this was built with a specific version tag
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			// Clean up the version string
			version := info.Main.Version
			
			// Go adds "v" prefix, keep it
			// Example: "v0.1.0" or "v0.0.0-20231231120000-abc123def456"
			if strings.HasPrefix(version, "v") {
				// If it's a clean tag like v0.1.0, return it
				if !strings.Contains(version, "-") {
					return version
				}
				
				// If it's a pseudo-version, extract the base version
				// v0.0.0-20231231-abc123 -> v0.0.0 (not useful)
				// v0.1.1-0.20231231-abc -> v0.1.1-dev
				parts := strings.Split(version, "-")
				if len(parts) > 1 && parts[0] != "v0.0.0" {
					return parts[0] + "-dev"
				}
			}
			
			return version
		}
	}

	// Fallback to "dev" if no version info available
	return "dev"
}