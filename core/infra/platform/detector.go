// Package platform provides platform detection and factory services.
package platform

import (
	"os"
	"runtime"
	"strings"
)

// Type represents the detected platform type.
type Type string

const (
	MacOS Type = "macos"
	Linux Type = "linux"
	WSL2  Type = "wsl2"
)

// Detector detects the current platform.
type Detector struct{}

// NewDetector creates a new platform detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect returns the current platform type.
func (d *Detector) Detect() Type {
	switch runtime.GOOS {
	case "darwin":
		return MacOS
	case "linux":
		if d.isWSL2() {
			return WSL2
		}
		return Linux
	default:
		// Default to Linux for other Unix-like systems
		return Linux
	}
}

// isWSL2 checks if running under WSL2.
func (d *Detector) isWSL2() bool {
	// Check for WSL-specific files
	if d.fileExists("/proc/sys/fs/binfmt_misc/WSLInterop") {
		return true
	}

	// Check /proc/version for WSL indicators
	if data, err := os.ReadFile("/proc/version"); err == nil {
		version := strings.ToLower(string(data))
		if strings.Contains(version, "microsoft") || strings.Contains(version, "wsl") {
			return true
		}
	}

	// Check for /run/WSL directory
	if d.dirExists("/run/WSL") {
		return true
	}

	return false
}

// fileExists checks if a file exists.
func (d *Detector) fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// dirExists checks if a directory exists.
func (d *Detector) dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// String returns the string representation of the platform type.
func (t Type) String() string {
	return string(t)
}
