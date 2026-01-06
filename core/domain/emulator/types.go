// Package emulator provides QEMU emulation orchestration for elmos.
// This file contains type definitions for the emulator package.
package emulator

// RunOptions contains options for running QEMU.
type RunOptions struct {
	Debug     bool // Enable GDB stub
	Graphical bool // Use graphical display instead of serial console
}
