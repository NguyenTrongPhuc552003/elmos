// Package main is the entry point for the elmos CLI application.
// ELMOS (Embedded Linux on MacOS) provides native Linux kernel build tools for macOS.
package main

import (
	"os"

	"github.com/NguyenTrongPhuc552003/elmos/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
