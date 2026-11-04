package builder

import (
	"os"
	"path/filepath"
	"strings"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// getToolchainEnv returns the environment variables with modified PATH and CROSS_COMPILE
// if a toolchain is found for the given architecture.
// Returns env, crossCompilePrefix, error
func getToolchainEnv(ctx *elcontext.Context, cfg *elconfig.Config, tm *toolchain.Manager, fs filesystem.FileSystem, arch string) ([]string, string, error) {
	// Start with default make environment
	env := ctx.GetMakeEnv()
	defaultCross := cfg.Build.CrossCompile

	// Get arch config to find toolchain binary name
	archCfg := elconfig.GetArchConfig(arch)
	if archCfg == nil || archCfg.GCCBinary == "" {
		return env, defaultCross, nil
	}

	// Extract target tuple from GCCBinary (e.g., riscv64-unknown-linux-gnu-gcc -> riscv64-unknown-linux-gnu)
	binary := archCfg.GCCBinary
	if !strings.HasSuffix(binary, "-gcc") {
		return env, defaultCross, nil
	}
	target := strings.TrimSuffix(binary, "-gcc")

	// Check if toolchain is installed in x-tools
	xTools := tm.Paths().XTools
	binDir := filepath.Join(xTools, target, "bin")

	if !fs.IsDir(binDir) {
		// Not installed, fallback to default
		return env, defaultCross, nil
	}

	// Detected toolchain! Update environment.
	crossCompile := target + "-"
	newPath := binDir + string(os.PathListSeparator)

	var newEnv []string
	pathUpdated := false
	crossUpdated := false

	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			// Prepend new binDir to PATH
			newEnv = append(newEnv, "PATH="+newPath+e[5:])
			pathUpdated = true
		} else if strings.HasPrefix(e, "CROSS_COMPILE=") {
			// Override CROSS_COMPILE
			newEnv = append(newEnv, "CROSS_COMPILE="+crossCompile)
			crossUpdated = true
		} else {
			newEnv = append(newEnv, e)
		}
	}

	if !pathUpdated {
		// Should generally interpret existing PATH, but GetMakeEnv guarantees PATH is set.
		// Fallback just in case implementation changes.
		newEnv = append(newEnv, "PATH="+newPath+os.Getenv("PATH"))
	}
	if !crossUpdated {
		newEnv = append(newEnv, "CROSS_COMPILE="+crossCompile)
	}

	return newEnv, crossCompile, nil
}
