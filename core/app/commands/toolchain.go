package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// BuildToolchains creates the toolchains command tree for crosstool-ng management.
func BuildToolchains(ctx *Context) *cobra.Command {
	toolchainsCmd := &cobra.Command{
		Use:   "toolchains",
		Short: "Manage cross-compiler toolchains (crosstool-ng)",
		Long: `Manage cross-compiler toolchains using crosstool-ng.

Subcommands allow you to install crosstool-ng, list available targets,
select a target configuration, build toolchains, and more.

Examples:
  elmos toolchains install              # Install crosstool-ng
  elmos toolchains list                 # List available target samples
  elmos toolchains riscv64-unknown-linux-gnu  # Select target
  elmos toolchains build                # Build the selected toolchain
  elmos toolchains build -j8            # Build with 8 parallel jobs`,
	}

	// Install command
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install crosstool-ng from latest git",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Installing crosstool-ng...")
			if err := ctx.ToolchainManager.Install(cmd.Context()); err != nil {
				return err
			}
			ctx.Printer.Success("crosstool-ng installed!")
			return nil
		},
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available toolchain samples",
		RunE: func(cmd *cobra.Command, args []string) error {
			samples, err := ctx.ToolchainManager.ListSamples(cmd.Context())
			if err != nil {
				return err
			}
			if len(samples) == 0 {
				ctx.Printer.Info("No samples found")
				return nil
			}
			ctx.Printer.Print("Available targets:")
			for _, s := range samples {
				ctx.Printer.Print("  %s", s)
			}
			return nil
		},
	}

	// Status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show installed toolchains status",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if ct-ng is installed
			if !ctx.ToolchainManager.IsInstalled() {
				ctx.Printer.Warn("crosstool-ng not installed")
				ctx.Printer.Print("  Run: elmos toolchains install")
				return nil
			}
			ctx.Printer.Success("crosstool-ng installed at %s", ctx.ToolchainManager.Paths().CrosstoolNG)

			// List installed toolchains
			toolchains, err := ctx.ToolchainManager.GetInstalledToolchains()
			if err != nil {
				return err
			}
			if len(toolchains) == 0 {
				ctx.Printer.Info("No toolchains built yet")
				return nil
			}
			ctx.Printer.Print("")
			ctx.Printer.Print("Installed toolchains:")
			for _, tc := range toolchains {
				status := "✓"
				if !tc.Installed {
					status = "○"
				}
				ctx.Printer.Print("  %s %s", status, tc.Target)
			}
			return nil
		},
	}

	// Build command
	var buildJobs int
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the currently configured toolchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Building toolchain...")
			if err := ctx.ToolchainManager.Build(cmd.Context(), buildJobs); err != nil {
				return err
			}
			ctx.Printer.Success("Toolchain built!")
			return nil
		},
	}
	buildCmd.Flags().IntVarP(&buildJobs, "jobs", "j", 0, "Number of parallel jobs (default: CPU count)")

	// Menuconfig command
	menuconfigCmd := &cobra.Command{
		Use:   "menuconfig",
		Short: "Open interactive configuration menu",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ctx.ToolchainManager.Menuconfig(cmd.Context())
		},
	}

	// Clean command
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean build artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Printer.Step("Cleaning build artifacts...")
			if err := ctx.ToolchainManager.Clean(cmd.Context()); err != nil {
				return err
			}
			ctx.Printer.Success("Cleaned!")
			return nil
		},
	}

	// Env command
	envCmd := &cobra.Command{
		Use:   "env",
		Short: "Print environment variables for shell integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := ctx.ToolchainManager.Paths()
			// Print shell-compatible export statements
			fmt.Printf("export PATH=\"%s/bin:$PATH\"\n", paths.CrosstoolNG)

			toolchains, _ := ctx.ToolchainManager.GetInstalledToolchains()
			for _, tc := range toolchains {
				if tc.Installed {
					fmt.Printf("export PATH=\"%s/bin:$PATH\"\n", tc.Path)
				}
			}
			return nil
		},
	}

	// Add subcommands
	toolchainsCmd.AddCommand(installCmd, listCmd, statusCmd, buildCmd, menuconfigCmd, cleanCmd, envCmd)

	return toolchainsCmd
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
