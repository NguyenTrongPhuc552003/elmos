package build

import (
	"fmt"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
)

// buildKernelConfigCmd creates the kernel config subcommand.
func buildKernelConfigCmd(ctx *types.Context) *cobra.Command {
	var enableOpts []string
	cmd := &cobra.Command{
		Use:   "config [type]",
		Short: "Configure the kernel",
		Long: `Configure the kernel with a config target or enable specific options.

Examples:
  elmos kernel config              # Run defconfig
  elmos kernel config menuconfig   # Interactive config
  elmos kernel config -E NETFILTER # Enable CONFIG_NETFILTER`,
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
			// Handle --enable options
			if len(enableOpts) > 0 {
				return enableKernelConfigs(ctx, cmd, enableOpts)
			}

			configType := "defconfig"
			if len(args) > 0 {
				configType = args[0]
			}
			ctx.Printer.Step("Running kernel %s...", configType)
			if err := ctx.KernelBuilder.Configure(cmd.Context(), configType); err != nil {
				return err
			}
			ctx.Printer.Success("Kernel configured!")
			return nil
		}),
	}
	cmd.Flags().StringArrayVarP(&enableOpts, "enable", "E", nil, "Enable kernel config option (e.g., NETFILTER)")
	return cmd
}

// enableKernelConfigs enables specific kernel config options and runs oldconfig.
func enableKernelConfigs(ctx *types.Context, cmd *cobra.Command, opts []string) error {
	configPath := ctx.Config.Paths.KernelDir + "/.config"
	scriptsConfig := ctx.Config.Paths.KernelDir + "/scripts/config"

	// Enable each option
	for _, opt := range opts {
		// Add CONFIG_ prefix if not present
		configName := opt
		if !strings.HasPrefix(opt, "CONFIG_") {
			configName = "CONFIG_" + opt
		}
		ctx.Printer.Step("Enabling %s...", configName)
		if err := ctx.Exec.Run(cmd.Context(), scriptsConfig, "--file", configPath, "--enable", configName); err != nil {
			return fmt.Errorf("failed to enable %s: %w", configName, err)
		}
	}

	// Run oldconfig to resolve dependencies
	ctx.Printer.Step("Updating config...")
	if err := ctx.KernelBuilder.Configure(cmd.Context(), "olddefconfig"); err != nil {
		return err
	}

	ctx.Printer.Success("Config updated! Run 'elmos kernel build' to apply changes.")
	return nil
}

// buildKernelCleanCmd creates the kernel clean subcommand.
func buildKernelCleanCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean kernel build artifacts",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
			ctx.Printer.Step("Cleaning...")
			if err := ctx.KernelBuilder.Clean(cmd.Context()); err != nil {
				return err
			}
			ctx.Printer.Success("Cleaned!")
			return nil
		}),
	}
}

// buildKernelBuildCmd creates the kernel build subcommand.
func buildKernelBuildCmd(ctx *types.Context) *cobra.Command {
	var jobs int
	cmd := &cobra.Command{
		Use:   "build [targets...]",
		Short: "Build the Linux kernel",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
			targets := args
			if len(targets) == 0 {
				targets = ctx.KernelBuilder.GetDefaultTargets()
			}
			ctx.Printer.Step("Building kernel for %s...", ctx.Config.Build.Arch)
			if err := ctx.KernelBuilder.Build(cmd.Context(), builder.BuildOptions{Jobs: jobs, Targets: targets}); err != nil {
				return err
			}
			ctx.Printer.Success("Build complete!")
			return nil
		}),
	}
	cmd.Flags().IntVarP(&jobs, "jobs", "j", 0, "Number of parallel build jobs")
	return cmd
}

// printKernelBuildStatus prints kernel config and image status.
func printKernelBuildStatus(ctx *types.Context) {
	ctx.Printer.Print("")
	ctx.Printer.Step("Build status:")
	if ctx.AppContext.HasConfig() {
		ctx.Printer.Print("  ✓ Kernel configured (.config exists)")
	} else {
		ctx.Printer.Print("  ○ Not configured (run 'elmos kernel config')")
	}
	if ctx.AppContext.HasKernelImage() {
		ctx.Printer.Print("  ✓ Kernel image built")
	} else {
		ctx.Printer.Print("  ○ Kernel not built (run 'elmos build')")
	}
}
