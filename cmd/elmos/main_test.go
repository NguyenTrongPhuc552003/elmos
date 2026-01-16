// Package main is the entry point for the elmos CLI application.
// ELMOS (Embedded Linux on MacOS) provides native Linux kernel build tools for macOS.
package main

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
