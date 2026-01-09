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
		{"Tap: messense/macos-cross-toolchains", "Tap"},
		{"Package: llvm", "Package"},
		{"Package: qemu", "Package"},
		{"Header: elf.h", "Header"},
		{"Header: asm/", "Header"},
		{"GDB: arm64 (aarch64-unknown-linux-gnu-gdb)", "GDB"},
		{"GCC: arm64 (aarch64-unknown-linux-gnu-gcc)", "GCC"},
		{"Toolchain: riscv64", "Toolchains"},
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
