// Package packages provides abstractions for system package management across platforms.
package packages

// Resolver abstracts platform-specific package management operations.
// Implementations handle resolving package paths, checking installations, and listing packages.
type Resolver interface {
	// GetPrefix returns the installation prefix for a package.
	GetPrefix(pkg string) string

	// GetBin returns the bin directory path for a package.
	GetBin(pkg string) string

	// GetSbin returns the sbin directory path for a package (if applicable).
	GetSbin(pkg string) string

	// GetInclude returns the include directory path for a package.
	GetInclude(pkg string) string

	// GetLib returns the lib directory path for a package.
	GetLib(pkg string) string

	// GetLibexecBin returns the libexec/gnubin path for GNU tools (macOS-specific, optional).
	GetLibexecBin(pkg string) string

	// ListInstalled returns a list of installed packages.
	ListInstalled() ([]string, error)

	// IsInstalled checks if a package is installed.
	IsInstalled(pkg string) bool

	// ClearCache clears any cached package information.
	ClearCache()
}
