package commands

import (
	"testing"
)

func TestGetSection(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Homebrew", "Package Manager"},
		{"Tap: messense/macos-cross-toolchains", "Homebrew Taps"},
		{"Package: llvm", "Homebrew Packages"},
		{"Package: qemu", "Homebrew Packages"},
		{"Header: elf.h", "Custom Headers"},
		{"Header: asm/", "Custom Headers"},
		{"GDB: arm64 (aarch64-unknown-linux-gnu-gdb)", "Cross Debuggers"},
		{"GCC: arm64 (aarch64-unknown-linux-gnu-gcc)", "Cross Compilers"},
		{"Unknown item", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSection(tt.name); got != tt.want {
				t.Errorf("getSection(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
