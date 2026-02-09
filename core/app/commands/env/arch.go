package env

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// BuildArch creates the arch command for architecture management.
func BuildArch(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "arch [target]",
		Short: "Set or show target architecture",
		Long: `Manage target architecture for cross-compilation.

Examples:
  elmos arch           # Show current config (or init if none)
  elmos arch arm64     # Set architecture to arm64
  elmos arch show      # Show detailed configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return handleArchNoArgs(ctx)
			}
			target := args[0]
			if target == "show" {
				return showArchConfig(ctx)
			}
			return setArchTarget(ctx, cmd, target)
		},
	}
}

// --- Helper functions to reduce BuildArch complexity ---

// handleArchNoArgs handles the case when no arguments are provided.
func handleArchNoArgs(ctx *types.Context) error {
	if ctx.Config.ConfigFile == "" {
		cwd, _ := os.Getwd()
		configPath := filepath.Join(cwd, "elmos.yaml")
		if !ctx.FS.Exists(configPath) {
			cfg := &config.Config{
				Build: config.BuildConfig{Arch: "arm64", LLVM: true, CrossCompile: "llvm-"},
				Image: config.ImageConfig{Size: config.DefaultImageSize, VolumeName: config.DefaultVolumeName},
				QEMU:  config.QEMUConfig{Memory: "2G", GDBPort: 1234, SSHPort: 2222},
			}
			if err := cfg.Save(configPath); err != nil {
				return err
			}
			ctx.Printer.Success("Initialized config: %s", configPath)
			return nil
		}
	}
	ctx.Printer.Print("Architecture: %s", ctx.Config.Build.Arch)
	return nil
}

// showArchConfig displays detailed configuration information.
func showArchConfig(ctx *types.Context) error {
	ctx.Printer.Print("Current Configuration:")
	ctx.Printer.Print("  Architecture:  %s", ctx.Config.Build.Arch)
	ctx.Printer.Print("  Jobs:          %d", ctx.Config.Build.Jobs)
	ctx.Printer.Print("  LLVM:          %v", ctx.Config.Build.LLVM)
	ctx.Printer.Print("  Memory:        %s", ctx.Config.QEMU.Memory)
	ctx.Printer.Print("  Project Root:  %s", ctx.Config.Paths.ProjectRoot)
	ctx.Printer.Print("  Volume:        %s", ctx.Config.Image.MountPoint)
	ctx.Printer.Print("  Config File:   %s", ctx.Config.ConfigFile)
	return nil
}

// setArchTarget sets the architecture and optionally configures the toolchain.
func setArchTarget(ctx *types.Context, cmd *cobra.Command, target string) error {
	archName, toolchainTarget := resolveArchAndToolchain(target)

	// Set Architecture in Config
	if archName != "" {
		if err := saveArchConfig(ctx, archName); err != nil {
			return err
		}
	} else {
		ctx.Printer.Warn("Could not infer architecture from '%s', only toolchain will be set", target)
	}

	// Select Toolchain (if applicable)
	if toolchainTarget != "" {
		return selectToolchain(ctx, cmd, toolchainTarget)
	}
	return nil
}

// resolveArchAndToolchain determines arch name and toolchain from target.
func resolveArchAndToolchain(target string) (archName, toolchainTarget string) {
	if config.IsValidArch(target) {
		archName = target
		if cfg := config.GetArchConfig(target); cfg != nil && cfg.GCCBinary != "" {
			if len(cfg.GCCBinary) > 4 && cfg.GCCBinary[len(cfg.GCCBinary)-4:] == "-gcc" {
				toolchainTarget = cfg.GCCBinary[:len(cfg.GCCBinary)-4]
			}
		}
	} else {
		toolchainTarget = target
		archName = inferArchFromToolchain(target)
	}
	return
}

// inferArchFromToolchain guesses architecture from toolchain name.
func inferArchFromToolchain(target string) string {
	if containsIgnoreCase(target, "riscv") {
		return "riscv"
	} else if containsIgnoreCase(target, "arm64") || containsIgnoreCase(target, "aarch64") {
		return "arm64"
	} else if containsIgnoreCase(target, "arm") {
		return "arm"
	}
	return ""
}

// saveArchConfig saves the architecture configuration to file.
func saveArchConfig(ctx *types.Context, archName string) error {
	ctx.Config.Build.Arch = archName
	configPath := ctx.Config.ConfigFile
	if configPath == "" {
		configPath = filepath.Join(ctx.Config.Paths.ProjectRoot, "elmos.yaml")
	}
	if err := ctx.Config.Save(configPath); err != nil {
		return err
	}
	ctx.Printer.Success("Architecture set to: %s", archName)
	return nil
}

// selectToolchain selects and configures the toolchain.
func selectToolchain(ctx *types.Context, cmd *cobra.Command, toolchainTarget string) error {
	if err := ctx.AppContext.EnsureMounted(); err != nil {
		return err
	}
	ctx.Printer.Step("Selecting toolchain: %s", toolchainTarget)
	if err := ctx.ToolchainManager.SelectTarget(cmd.Context(), toolchainTarget); err != nil {
		return err
	}
	ctx.Printer.Success("Toolchain selected: %s", toolchainTarget)
	ctx.Printer.Print("  Run 'elmos toolchains build' to build it")
	return nil
}
