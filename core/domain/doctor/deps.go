package doctor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// CheckPackages checks if required Homebrew packages are installed.
func (h *HealthChecker) CheckPackages(ctx context.Context) []CheckResult {
	pkgSet, err := h.getInstalledPackageSet()
	if err != nil {
		return []CheckResult{{
			Name:     "Homebrew Packages",
			Passed:   false,
			Required: true,
			Message:  "Failed to list packages",
		}}
	}

	grouped := h.groupPackagesByCategory()
	cats := []string{"Build Tools", "Virtualization", "Toolchain Dependencies"}

	var results []CheckResult
	for _, catName := range cats {
		results = append(results, h.checkCategoryPackages(catName, grouped[catName], pkgSet)...)
	}

	return results
}

// getInstalledPackageSet returns a set of installed package names.
func (h *HealthChecker) getInstalledPackageSet() (map[string]bool, error) {
	installed, err := h.brew.ListInstalled()
	if err != nil {
		return nil, err
	}
	pkgSet := make(map[string]bool)
	for _, p := range installed {
		pkgSet[p] = true
	}
	return pkgSet, nil
}

// groupPackagesByCategory groups required packages by their category.
func (h *HealthChecker) groupPackagesByCategory() map[string][]elconfig.RequiredPackage {
	grouped := make(map[string][]elconfig.RequiredPackage)
	for _, pkg := range elconfig.RequiredPackages {
		cat := pkg.Category
		if cat == "" {
			cat = "Other"
		}
		grouped[cat] = append(grouped[cat], pkg)
	}
	return grouped
}

// checkCategoryPackages checks packages in a category and returns results.
func (h *HealthChecker) checkCategoryPackages(catName string, pkgs []elconfig.RequiredPackage, pkgSet map[string]bool) []CheckResult {
	if len(pkgs) == 0 {
		return nil
	}

	var results []CheckResult
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

	return results
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
