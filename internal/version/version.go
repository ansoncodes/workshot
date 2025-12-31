package version

import "runtime/debug"

// default for local builds
var Version = "dev"

func Get() string {
	// If injected via -ldflags, prefer it
	if Version != "dev" {
		return Version
	}

	// Read version from go install @version / @latest
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}
