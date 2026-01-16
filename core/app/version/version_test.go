// Package version provides version information for the elmos CLI.
package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	got := Get()

	// Version should not be empty
	if got.Version == "" {
		t.Error("Get().Version is empty")
	}

	// GoVersion should match runtime
	if got.GoVersion != runtime.Version() {
		t.Errorf("Get().GoVersion = %v, want %v", got.GoVersion, runtime.Version())
	}

	// OS should match runtime
	if got.OS != runtime.GOOS {
		t.Errorf("Get().OS = %v, want %v", got.OS, runtime.GOOS)
	}

	// Arch should match runtime
	if got.Arch != runtime.GOARCH {
		t.Errorf("Get().Arch = %v, want %v", got.Arch, runtime.GOARCH)
	}
}

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name string
		i    Info
		want string
	}{
		{
			name: "Full info",
			i: Info{
				Version:   "v1.0.0",
				Commit:    "abc1234567",
				BuildDate: "2024-01-01",
				GoVersion: "go1.21.0",
			},
			want: "elmos v1.0.0 (abc1234) built on 2024-01-01 with go1.21.0",
		},
		{
			name: "Short commit",
			i: Info{
				Version:   "v2.0.0",
				Commit:    "def",
				BuildDate: "2024-06-01",
				GoVersion: "go1.22.0",
			},
			want: "elmos v2.0.0 (def) built on 2024-06-01 with go1.22.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.i.String()
			if got != tt.want {
				t.Errorf("Info.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_Short(t *testing.T) {
	tests := []struct {
		name string
		i    Info
		want string
	}{
		{
			name: "Returns version",
			i:    Info{Version: "v1.2.3"},
			want: "v1.2.3",
		},
		{
			name: "Dev version",
			i:    Info{Version: "dev"},
			want: "dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Short(); got != tt.want {
				t.Errorf("Info.Short() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_String_ContainsElmos(t *testing.T) {
	info := Get()
	s := info.String()
	if !strings.Contains(s, "elmos") {
		t.Error("Info.String() should contain 'elmos'")
	}
}
