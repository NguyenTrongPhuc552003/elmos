package packages

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
)

// PacmanResolver implements package resolution for Pacman-based systems (Arch Linux, Manjaro).
type PacmanResolver struct {
	exec  executor.Executor
	cache map[string]string
}

// NewPacmanResolver creates a new Pacman package resolver.
func NewPacmanResolver(exec executor.Executor) Resolver {
	return &PacmanResolver{
		exec:  exec,
		cache: make(map[string]string),
	}
}

// GetPrefix returns the installation prefix for a package.
// On Arch, packages typically install to /usr.
func (r *PacmanResolver) GetPrefix(pkg string) string {
	if cached, ok := r.cache[pkg]; ok {
		return cached
	}

	// Check if package is installed
	out, err := r.exec.Output(context.Background(), "pacman", "-Ql", pkg)
	if err != nil {
		return ""
	}

	// Parse output to find common prefix (usually /usr)
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		prefix := "/usr"
		r.cache[pkg] = prefix
		return prefix
	}

	return ""
}

// GetBin returns the bin directory path.
// Standard location is /usr/bin on Arch Linux.
func (r *PacmanResolver) GetBin(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/bin" // Default fallback
	}
	return filepath.Join(prefix, "bin")
}

// GetSbin returns the sbin directory path.
// Note: Arch Linux merged /usr/sbin into /usr/bin, but we keep this for compatibility.
func (r *PacmanResolver) GetSbin(pkg string) string {
	return "/usr/bin" // Arch merged sbin into bin
}

// GetInclude returns the include directory path.
// Standard location is /usr/include on Arch Linux.
func (r *PacmanResolver) GetInclude(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/include" // Default fallback
	}
	return filepath.Join(prefix, "include")
}

// GetLib returns the lib directory path.
// Standard location is /usr/lib on Arch Linux.
func (r *PacmanResolver) GetLib(pkg string) string {
	prefix := r.GetPrefix(pkg)
	if prefix == "" {
		return "/usr/lib" // Default fallback
	}
	return filepath.Join(prefix, "lib")
}

// GetLibexecBin returns empty string as Linux doesn't use libexec/gnubin pattern.
// GNU tools are directly in /usr/bin on Arch Linux.
func (r *PacmanResolver) GetLibexecBin(pkg string) string {
	return "" // Not applicable on Arch Linux
}

// ListInstalled returns a list of installed packages.
func (r *PacmanResolver) ListInstalled() ([]string, error) {
	out, err := r.exec.Output(context.Background(), "pacman", "-Qq")
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

// IsInstalled checks if a package is installed.
func (r *PacmanResolver) IsInstalled(pkg string) bool {
	// Use pacman -Q for fast check
	err := r.exec.Run(context.Background(), "pacman", "-Q", pkg)
	return err == nil
}

// ClearCache clears the package prefix cache.
func (r *PacmanResolver) ClearCache() {
	r.cache = make(map[string]string)
}
