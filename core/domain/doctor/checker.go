// Package doctor provides dependency checking and environment validation for elmos.
package doctor

import (
	"context"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/homebrew"
)

// CheckResult represents the result of a single check.
type CheckResult struct {
	Name     string
	Passed   bool
	Required bool
	Message  string
}

// HealthChecker validates the development environment.
type HealthChecker struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
	brew *homebrew.Resolver
	tm   *toolchain.Manager
}

// NewHealthChecker creates a new HealthChecker with the given dependencies.
func NewHealthChecker(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config, tm *toolchain.Manager) *HealthChecker {
	return &HealthChecker{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
		brew: homebrew.NewResolver(exec),
		tm:   tm,
	}
}

// CheckAll runs all environment checks and returns the results.
func (h *HealthChecker) CheckAll(ctx context.Context) ([]CheckResult, int) {
	var results []CheckResult
	issues := 0

	// Check Homebrew
	result := h.CheckHomebrew(ctx)
	results = append(results, result)
	if !result.Passed && result.Required {
		issues++
	}

	// Check packages
	pkgResults := h.CheckPackages(ctx)
	for _, r := range pkgResults {
		results = append(results, r)
		if !r.Passed && r.Required {
			issues++
		}
	}

	// Check headers
	headerResults := h.CheckHeaders(ctx)
	for _, r := range headerResults {
		results = append(results, r)
		if !r.Passed && r.Required {
			issues++
		}
	}

	// Check cross debuggers
	gdbResults := h.CheckCrossGDB(ctx)
	results = append(results, gdbResults...)

	// Check cross compilers
	gccResults := h.CheckCrossGCC(ctx)
	results = append(results, gccResults...)

	// Check toolchains (ct-ng and installed targets)
	tcResults := h.CheckToolchains(ctx)
	results = append(results, tcResults...)

	return results, issues
}
