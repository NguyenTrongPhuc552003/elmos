// Package doctor provides dependency checking and environment validation for elmos.
package doctor

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"codeberg.org/NguyenTrongPhuc552003/elmos/assets"

	elconfig "codeberg.org/NguyenTrongPhuc552003/elmos/core/config"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// AutoFixer provides methods to automatically fix common issues.
type AutoFixer struct {
	fs  filesystem.FileSystem
	cfg *elconfig.Config
}

// NewAutoFixer creates a new AutoFixer with the given dependencies.
func NewAutoFixer(fs filesystem.FileSystem, cfg *elconfig.Config) *AutoFixer {
	return &AutoFixer{
		fs:  fs,
		cfg: cfg,
	}
}

// MissingHeaders lists the embedded host headers that are not present in
// the project's libraries directory.
func (f *AutoFixer) MissingHeaders() []string {
	var missing []string
	_ = fs.WalkDir(assets.Libraries, "libraries", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel := strings.TrimPrefix(p, "libraries/")
		if !f.fs.Exists(filepath.Join(f.cfg.Paths.LibrariesDir, filepath.FromSlash(rel))) {
			missing = append(missing, rel)
		}
		return nil
	})
	return missing
}

// FixHeaders writes the missing host headers from the copies embedded in
// elmos (audit E14: a project made by 'elmos init' had none of them, so the
// kernel host tools failed on byteswap.h). Existing files are left alone.
func (f *AutoFixer) FixHeaders() ([]string, error) {
	missing := f.MissingHeaders()
	for _, rel := range missing {
		data, err := fs.ReadFile(assets.Libraries, "libraries/"+rel)
		if err != nil {
			return nil, err
		}
		dst := filepath.Join(f.cfg.Paths.LibrariesDir, filepath.FromSlash(rel))
		if err := f.fs.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		if err := f.fs.WriteFile(dst, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", rel, err)
		}
	}
	return missing, nil
}
