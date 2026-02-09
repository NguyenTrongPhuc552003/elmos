package packages

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
)

// AptResolver implements package resolution for APT-based systems (Debian/Ubuntu).
type AptResolver struct {
	exec  executor.Executor
	cache map[string]string
}

// NewAptResolver creates a new APT package resolver.
func NewAptResolver(exec executor.Executor) Resolver {
	return &AptResolver{
		exec:  exec,
		cache: make(map[string]string),
	}
}

// GetPrefix returns the installation prefix for a package.
// On Linux, most packages install to /usr, but we return a cached lookup.
func (r *AptResolver) GetPrefix(pkg string) string {
	if cached, ok := r.cache[pkg]; ok {
		return cached
	}

	// Try to find package installation path using dpkg
	out, err := r.exec.Output(context.Background(), "dpkg", "-L", pkg)
	if err != nil {
		return ""
	}

	// Parse output to find common prefix (usually /usr)
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		// Most system packages install to /usr
		prefix := "/usr"
		r.cache[pkg] = prefix
		return prefix
	}

	return ""
}

// GetBin returns the bin directory path.
// Standard location is /usr/bin on Linux.
func (r *AptResolver) GetBin(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/bin" // Default fallback
	}
	return filepath.Join(prefix, "bin")
}

// GetSbin returns the sbin directory path.
// Standard location is /usr/sbin on Linux.
func (r *AptResolver) GetSbin(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/sbin" // Default fallback
	}
	return filepath.Join(prefix, "sbin")
}

// GetInclude returns the include directory path.
// Standard location is /usr/include on Linux.
func (r *AptResolver) GetInclude(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/include" // Default fallback
	}
	return filepath.Join(prefix, "include")
}

// GetLib returns the lib directory path.
// Standard location is /usr/lib or /usr/lib/x86_64-linux-gnu on Linux.
func (r *AptResolver) GetLib(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/lib" // Default fallback
	}
	return filepath.Join(prefix, "lib")
}

// GetLibexecBin returns empty string as Linux doesn't use libexec/gnubin pattern.
// GNU tools are directly in /usr/bin on Linux.
func (r *AptResolver) GetLibexecBin(pkg string) string {
	return "" // Not applicable on Linux
}

// ListInstalled returns a list of installed packages.
func (r *AptResolver) ListInstalled() ([]string, error) {
	out, err := r.exec.Output(context.Background(), "dpkg-query", "-f", "${Package}\n", "-W")
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

// IsInstalled checks if a package is installed.
func (r *AptResolver) IsInstalled(pkg string) bool {
	// Use dpkg-query for faster check
	err := r.exec.Run(context.Background(), "dpkg-query", "-W", "-f=${Status}", pkg)
	return err == nil
}

// ClearCache clears the package prefix cache.
func (r *AptResolver) ClearCache() {
	r.cache = make(map[string]string)
}
