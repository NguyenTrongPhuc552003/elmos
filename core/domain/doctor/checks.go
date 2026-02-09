package doctor

import (
	"path/filepath"
)

// IsElfHMissing checks if elf.h is missing.
func (h *HealthChecker) IsElfHMissing() bool {
	elfPath := filepath.Join(h.cfg.Paths.LibrariesDir, "elf.h")
	return !h.fs.Exists(elfPath)
}
