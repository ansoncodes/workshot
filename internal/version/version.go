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
	// First, try to get version from Go module info (go install @vX.X.X)
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version := info.Main.Version
			
			// Go adds "v" prefix, keep it
			if strings.HasPrefix(version, "v") {
				// If it's a clean tag like v0.1.0, return it
				if !strings.Contains(version, "-") {
					return version
				}
				
				// If it's a pseudo-version, extract base
				parts := strings.Split(version, "-")
				if len(parts) > 1 && parts[0] != "v0.0.0" {
					return parts[0] + "-dev"
				}
			}
			
			return version
		}
	}

	// If no module info, check if Version was set by ldflags (GoReleaser)
	if Version != "" && Version != "dev" {
		return Version
	}

	// Fallback to "dev"
	return "dev"
}