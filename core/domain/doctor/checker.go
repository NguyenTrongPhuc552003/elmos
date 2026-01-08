// Package doctor provides dependency checking and environment validation for elmos.
package doctor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

// CheckToolchains checks if crosstool-ng and toolchains are installed.
func (h *HealthChecker) CheckToolchains(ctx context.Context) []CheckResult {
	var results []CheckResult

	// Check ct-ng installation
	ctngInstalled := h.tm.IsInstalled()
	results = append(results, CheckResult{
		Name:     "crosstool-ng",
		Passed:   ctngInstalled,
		Required: false, // Optional for building, but good to have
		Message:  "Run: elmos toolchains install",
	})

	if ctngInstalled {
		// Check installed toolchains
		toolchains, err := h.tm.GetInstalledToolchains()
		if err == nil {
			for _, tc := range toolchains {
				msg := ""
				if !tc.Installed {
					msg = "Toolchain built but not fully installed"
				}
				results = append(results, CheckResult{
					Name:     fmt.Sprintf("Toolchain: %s", tc.Target),
					Passed:   tc.Installed,
					Required: false,
					Message:  msg,
				})
			}
		}
	}

	return results
}

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

// CheckPackages checks if required Homebrew packages are installed.
func (h *HealthChecker) CheckPackages(ctx context.Context) []CheckResult {
	var results []CheckResult

	installedPkgs, err := h.brew.ListInstalled()
	if err != nil {
		return []CheckResult{{
			Name:     "Homebrew Packages",
			Passed:   false,
			Required: true,
			Message:  "Failed to list packages",
		}}
	}

	pkgSet := make(map[string]bool)
	for _, p := range installedPkgs {
		pkgSet[p] = true
	}

	cats := []string{"Build Tools", "Virtualization", "Toolchain Dependencies"}
	grouped := make(map[string][]elconfig.RequiredPackage)

	for _, pkg := range elconfig.RequiredPackages {
		cat := pkg.Category
		if cat == "" {
			cat = "Other"
		}
		grouped[cat] = append(grouped[cat], pkg)
	}

	for _, catName := range cats {
		pkgs := grouped[catName]
		if len(pkgs) == 0 {
			continue
		}

		var missing []string
		for _, pkg := range pkgs {
			passed := pkgSet[pkg.Name]
			results = append(results, CheckResult{
				Name:     fmt.Sprintf("Homebrew Packages: [%s] %s", catName, pkg.Name),
				Passed:   passed,
				Required: pkg.Required,
				Message:  pkg.Description,
			})
			if !passed && pkg.Required {
				missing = append(missing, pkg.Name)
			}
		}

		if len(missing) > 0 {
			results = append(results, CheckResult{
				Name:     "  Fix missing packages",
				Passed:   false,
				Required: false,
				Message:  fmt.Sprintf("brew install %s", strings.Join(missing, " ")),
			})
		}
	}

	return results
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

// CheckCrossGDB checks for cross-architecture GDB binaries.
func (h *HealthChecker) CheckCrossGDB(ctx context.Context) []CheckResult {
	var results []CheckResult

	for _, arch := range elconfig.SupportedArchitectures() {
		archCfg := elconfig.GetArchConfig(arch)
		if archCfg == nil || archCfg.GDBBinary == "" {
			continue
		}

		binary := archCfg.GDBBinary
		passed := false
		message := ""

		// Check elmos toolchains only (strict check)
		// Extract target tuple: binary minus "-gdb"
		if strings.HasSuffix(binary, "-gdb") {
			target := strings.TrimSuffix(binary, "-gdb")
			toolchainBin := filepath.Join(h.tm.Paths().XTools, target, "bin", binary)
			if h.fs.Exists(toolchainBin) {
				passed = true
				message = "Found in elmos toolchains"
			}
		}

		results = append(results, CheckResult{
			Name:     fmt.Sprintf("Cross Debuggers: %s (%s)", arch, binary),
			Passed:   passed,
			Required: false, // GDB is optional
			Message:  message,
		})
	}

	return results
}

// CheckCrossGCC checks for cross-architecture GCC binaries.
func (h *HealthChecker) CheckCrossGCC(ctx context.Context) []CheckResult {
	var results []CheckResult

	for _, arch := range elconfig.SupportedArchitectures() {
		archCfg := elconfig.GetArchConfig(arch)
		if archCfg == nil || archCfg.GCCBinary == "" {
			continue
		}

		binary := archCfg.GCCBinary
		passed := false
		message := ""

		// Check elmos toolchains only (strict check)
		// Extract target tuple: binary minus "-gcc"
		if strings.HasSuffix(binary, "-gcc") {
			target := strings.TrimSuffix(binary, "-gcc")
			toolchainBin := filepath.Join(h.tm.Paths().XTools, target, "bin", binary)
			if h.fs.Exists(toolchainBin) {
				passed = true
				message = "Found in elmos toolchains"
			}
		}

		results = append(results, CheckResult{
			Name:     fmt.Sprintf("Cross Compilers: %s (%s)", arch, binary),
			Passed:   passed,
			Required: false, // GCC is optional when using LLVM
			Message:  message,
		})
	}

	return results
}

// IsElfHMissing checks if elf.h is missing.
func (h *HealthChecker) IsElfHMissing() bool {
	elfPath := filepath.Join(h.cfg.Paths.LibrariesDir, "elf.h")
	return !h.fs.Exists(elfPath)
}
