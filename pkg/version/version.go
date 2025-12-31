// Package version provides version information for the elmos CLI.
package version

import (
	"fmt"
	"runtime"
)

// Build information (set via ldflags at build time)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the current version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("elmos %s (%s) built on %s with %s",
		i.Version, i.Commit[:min(7, len(i.Commit))], i.BuildDate, i.GoVersion)
}

// Short returns just the version number
func (i Info) Short() string {
	return i.Version
}
