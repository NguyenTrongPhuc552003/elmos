package packages

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/homebrew"
)

// HomebrewResolver adapts the existing Homebrew resolver to the Resolver interface.
type HomebrewResolver struct {
	brew *homebrew.Resolver
}

// NewHomebrewResolver creates a new Homebrew package resolver.
func NewHomebrewResolver(brew *homebrew.Resolver) Resolver {
	return &HomebrewResolver{
		brew: brew,
	}
}

// GetPrefix returns the installation prefix for a Homebrew package.
func (r *HomebrewResolver) GetPrefix(pkg string) string {
	return r.brew.GetPrefix(pkg)
}

// GetBin returns the bin directory path for a Homebrew package.
func (r *HomebrewResolver) GetBin(pkg string) string {
	return r.brew.GetBin(pkg)
}

// GetSbin returns the sbin directory path for a Homebrew package.
func (r *HomebrewResolver) GetSbin(pkg string) string {
	return r.brew.GetSbin(pkg)
}

// GetInclude returns the include directory path for a Homebrew package.
func (r *HomebrewResolver) GetInclude(pkg string) string {
	return r.brew.GetInclude(pkg)
}

// GetLib returns the lib directory path for a Homebrew package.
func (r *HomebrewResolver) GetLib(pkg string) string {
	return r.brew.GetLib(pkg)
}

// GetLibexecBin returns the libexec/gnubin path for GNU tools.
func (r *HomebrewResolver) GetLibexecBin(pkg string) string {
	return r.brew.GetLibexecBin(pkg)
}

// ListInstalled returns a list of installed Homebrew formulae.
func (r *HomebrewResolver) ListInstalled() ([]string, error) {
	return r.brew.ListInstalled()
}

// IsInstalled checks if a Homebrew package is installed.
func (r *HomebrewResolver) IsInstalled(pkg string) bool {
	return r.brew.IsInstalled(pkg)
}

// ClearCache clears the Homebrew prefix cache.
func (r *HomebrewResolver) ClearCache() {
	r.brew.ClearCache()
}
