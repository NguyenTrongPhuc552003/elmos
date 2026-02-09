package doctor

import (
	"context"
	"fmt"
	"path/filepath"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// CheckHomebrew checks if Homebrew is installed.
func (h *HealthChecker) CheckHomebrew(ctx context.Context) CheckResult {
	_, err := h.exec.LookPath("brew")
	return CheckResult{
		Name:     "Homebrew",
		Passed:   err == nil,
		Required: true,
		Message:  "Install from: https://brew.sh",
	}
}

// CheckHeaders checks if required header files exist.
func (h *HealthChecker) CheckHeaders(ctx context.Context) []CheckResult {
	var results []CheckResult

	headersDir := h.cfg.Paths.LibrariesDir

	if !h.fs.IsDir(headersDir) {
		return []CheckResult{{
			Name:     "Headers directory",
			Passed:   false,
			Required: true,
			Message:  fmt.Sprintf("Directory not found: %s", headersDir),
		}}
	}

	for _, header := range elconfig.RequiredHeaders {
		headerPath := filepath.Join(headersDir, header)
		passed := h.fs.Exists(headerPath)
		results = append(results, CheckResult{
			Name:     fmt.Sprintf("Custom Headers: %s", header),
			Passed:   passed,
			Required: true,
			Message:  "",
		})
	}

	// Check asm directory
	asmDir := filepath.Join(headersDir, "asm")
	results = append(results, CheckResult{
		Name:     "Custom Headers: asm/",
		Passed:   h.fs.IsDir(asmDir),
		Required: true,
		Message:  "",
	})

	return results
}
