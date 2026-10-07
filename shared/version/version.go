package version

import (
	"runtime"
	"runtime/debug"
)

// Info holds version information.
type Info struct {
	Version   string
	Commit    string
	BuildTime string
	GoVersion string
}

// Get returns version information.
func Get() Info {
	info := Info{
		GoVersion: runtime.Version(),
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.GoVersion = buildInfo.GoVersion
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				info.Commit = setting.Value
			case "vcs.time":
				info.BuildTime = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					info.Commit += "-dirty"
				}
			}
		}
	}

	// If version is not set from build info, use a default
	if info.Version == "" {
		info.Version = "dev"
	}

	return info
}

// String returns a string representation of the version info.
func (i Info) String() string {
	if i.Commit != "" {
		return i.Version + " (" + i.Commit + ")"
	}
	return i.Version
}
