package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// BuildArch creates the arch command for architecture management.
func BuildArch(ctx *Context) *cobra.Command {
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
				// No args = show or init
				if ctx.Config.ConfigFile == "" {
					// Init default config
					cwd, _ := os.Getwd()
					configPath := filepath.Join(cwd, "elmos.yaml")
					if !ctx.FS.Exists(configPath) {
						cfg := &config.Config{
							Build: config.BuildConfig{Arch: "arm64", LLVM: true, CrossCompile: "llvm-"},
							Image: config.ImageConfig{Size: "20G", VolumeName: "kernel-dev"},
							QEMU:  config.QEMUConfig{Memory: "2G", GDBPort: 1234, SSHPort: 2222},
						}
						if err := cfg.Save(configPath); err != nil {
							return err
						}
						ctx.Printer.Success("Initialized config: %s", configPath)
						return nil
					}
				}
				// Show current arch
				ctx.Printer.Print("Architecture: %s", ctx.Config.Build.Arch)
				return nil
			}

			target := args[0]
			var archName string
			var toolchainTarget string

			// Handle "show" subcommand
			if target == "show" {
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

			// Check if it's a known short architecture name
			if config.IsValidArch(target) {
				archName = target
				// Lookup default toolchain for this arch
				if cfg := config.GetArchConfig(target); cfg != nil {
					// Map short arch to ct-ng target if known
					// For riscv, defaults to riscv64-unknown-linux-gnu
					if cfg.GCCBinary != "" {
						// Extract toolchain prefix from GCC binary (remove -gcc suffix)
						// e.g. riscv64-unknown-linux-gnu-gcc -> riscv64-unknown-linux-gnu
						if len(cfg.GCCBinary) > 4 && cfg.GCCBinary[len(cfg.GCCBinary)-4:] == "-gcc" {
							toolchainTarget = cfg.GCCBinary[:len(cfg.GCCBinary)-4]
						}
					}
				}
			} else {
				// Assume it's a toolchain target (e.g. riscv64-unknown-linux-gnu)
				// Try to select it using ToolchainManager
				if err := ctx.AppContext.EnsureMounted(); err != nil {
					return err
				}

				// Try to infer architecture from toolchain name
				if containsIgnoreCase(target, "riscv") {
					archName = "riscv"
				} else if containsIgnoreCase(target, "arm64") || containsIgnoreCase(target, "aarch64") {
					archName = "arm64"
				} else if containsIgnoreCase(target, "arm") {
					archName = "arm"
				}

				toolchainTarget = target
			}

			// 1. Set Architecture in Config
			if archName != "" {
				ctx.Config.Build.Arch = archName
				configPath := ctx.Config.ConfigFile
				if configPath == "" {
					configPath = filepath.Join(ctx.Config.Paths.ProjectRoot, "elmos.yaml")
				}
				if err := ctx.Config.Save(configPath); err != nil {
					return err
				}
				ctx.Printer.Success("Architecture set to: %s", archName)
			} else {
				ctx.Printer.Warn("Could not infer architecture from '%s', only toolchain will be set", target)
			}

			// 2. Select Toolchain (if applicable)
			if toolchainTarget != "" {
				if err := ctx.AppContext.EnsureMounted(); err != nil {
					return err
				}

				ctx.Printer.Step("Selecting toolchain: %s", toolchainTarget)
				if err := ctx.ToolchainManager.SelectTarget(cmd.Context(), toolchainTarget); err != nil {
					// If exact match fails, maybe it was just an arch name without mapped toolchain?
					// Or the toolchain sample doesn't exist?
					return err
				}
				ctx.Printer.Success("Toolchain selected: %s", toolchainTarget)
				ctx.Printer.Print("  Run 'elmos toolchains build' to build it")
			}

			return nil
		},
	}
}
