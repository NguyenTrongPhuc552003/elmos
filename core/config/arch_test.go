// Package config provides configuration management for elmos.
package config

import (
	"testing"
)

func TestGetArchConfig(t *testing.T) {
	tests := []struct {
		name     string
		arch     string
		wantNil  bool
		wantName string
	}{
		{
			name:     "arm64 architecture",
			arch:     "arm64",
			wantNil:  false,
			wantName: "arm64",
		},
		{
			name:     "arm architecture",
			arch:     "arm",
			wantNil:  false,
			wantName: "arm",
		},
		{
			name:     "riscv architecture",
			arch:     "riscv",
			wantNil:  false,
			wantName: "riscv",
		},
		{
			name:    "invalid architecture",
			arch:    "x86",
			wantNil: true,
		},
		{
			name:    "empty architecture",
			arch:    "",
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetArchConfig(tt.arch)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetArchConfig() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("GetArchConfig() = nil, want non-nil")
			}
			if got.Name != tt.wantName {
				t.Errorf("GetArchConfig().Name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

func TestSupportedArchitectures(t *testing.T) {
	got := SupportedArchitectures()

	// Should return at least 3 architectures
	if len(got) < 3 {
		t.Errorf("SupportedArchitectures() returned %d archs, want at least 3", len(got))
	}

	// Should contain expected architectures
	expected := map[string]bool{"arm64": false, "arm": false, "riscv": false}
	for _, arch := range got {
		if _, ok := expected[arch]; ok {
			expected[arch] = true
		}
	}

	for arch, found := range expected {
		if !found {
			t.Errorf("SupportedArchitectures() missing %s", arch)
		}
	}
}

func TestIsValidArch(t *testing.T) {
	tests := []struct {
		name string
		arch string
		want bool
	}{
		{name: "arm64 valid", arch: "arm64", want: true},
		{name: "arm valid", arch: "arm", want: true},
		{name: "riscv valid", arch: "riscv", want: true},
		{name: "x86 invalid", arch: "x86", want: false},
		{name: "empty invalid", arch: "", want: false},
		{name: "random invalid", arch: "foobar", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidArch(tt.arch); got != tt.want {
				t.Errorf("IsValidArch() = %v, want %v", got, tt.want)
			}
		})
	}
}
