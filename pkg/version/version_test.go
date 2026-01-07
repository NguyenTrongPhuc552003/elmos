package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GoVersion != runtime.Version() {
		t.Errorf("GoVersion = %q, want %q", info.GoVersion, runtime.Version())
	}

	if info.OS != runtime.GOOS {
		t.Errorf("OS = %q, want %q", info.OS, runtime.GOOS)
	}

	if info.Arch != runtime.GOARCH {
		t.Errorf("Arch = %q, want %q", info.Arch, runtime.GOARCH)
	}
}

func TestInfoString(t *testing.T) {
	info := Info{
		Version:   "v1.0.0",
		Commit:    "abc1234567890",
		BuildDate: "2026-01-01",
		GoVersion: "go1.21.0",
	}

	got := info.String()

	if !strings.Contains(got, "v1.0.0") {
		t.Errorf("String() missing version: %q", got)
	}
	if !strings.Contains(got, "abc1234") {
		t.Errorf("String() missing short commit: %q", got)
	}
	if !strings.Contains(got, "2026-01-01") {
		t.Errorf("String() missing build date: %q", got)
	}
}

func TestInfoShort(t *testing.T) {
	info := Info{Version: "v2.5.0"}

	if got := info.Short(); got != "v2.5.0" {
		t.Errorf("Short() = %q, want %q", got, "v2.5.0")
	}
}

func TestInfoStringShortCommit(t *testing.T) {
	// Test with commit shorter than 7 chars
	info := Info{
		Version:   "dev",
		Commit:    "abc",
		BuildDate: "unknown",
		GoVersion: "go1.21.0",
	}

	got := info.String()
	if !strings.Contains(got, "abc") {
		t.Errorf("String() should handle short commit: %q", got)
	}
}
