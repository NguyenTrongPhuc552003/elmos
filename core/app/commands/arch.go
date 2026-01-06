package commands

import (
	"fmt"
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

			// Set architecture
			if !config.IsValidArch(target) {
				return fmt.Errorf("invalid architecture: %s (use: arm64, arm, riscv)", target)
			}
			ctx.Config.Build.Arch = target
			configPath := ctx.Config.ConfigFile
			if configPath == "" {
				configPath = filepath.Join(ctx.Config.Paths.ProjectRoot, "elmos.yaml")
			}
			if err := ctx.Config.Save(configPath); err != nil {
				return err
			}
			ctx.Printer.Success("Architecture set to: %s", target)
			return nil
		},
	}
}
