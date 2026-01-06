// Package homebrew provides utilities for resolving Homebrew package paths.
// This file contains type definitions and interfaces for the homebrew package.
package homebrew

// PathResolver is the interface for resolving Homebrew package paths.
type PathResolver interface {
	// GetPrefix returns the installation prefix for a Homebrew package.
	GetPrefix(pkg string) string
	// GetBin returns the bin directory path for a Homebrew package.
	GetBin(pkg string) string
	// GetSbin returns the sbin directory path for a Homebrew package.
	GetSbin(pkg string) string
	// GetInclude returns the include directory path for a Homebrew package.
	GetInclude(pkg string) string
	// GetLib returns the lib directory path for a Homebrew package.
	GetLib(pkg string) string
	// GetLibexecBin returns the libexec/gnubin path for GNU tools.
	GetLibexecBin(pkg string) string
	// ListInstalled returns a list of installed Homebrew formulae.
	ListInstalled() ([]string, error)
	// ListTaps returns a list of tapped Homebrew repositories.
	ListTaps() ([]string, error)
	// IsInstalled checks if a Homebrew package is installed.
	IsInstalled(pkg string) bool
	// IsTapped checks if a Homebrew tap is tapped.
	IsTapped(tap string) bool
	// ClearCache clears the prefix cache.
	ClearCache()
}
