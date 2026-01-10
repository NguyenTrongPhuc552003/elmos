package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
)

// BuildKernel creates the kernel command tree for kernel management.
func BuildKernel(ctx *Context) *cobra.Command {
	kernelCmd := &cobra.Command{
		Use:   "kernel",
		Short: "Kernel configuration commands",
	}

	configCmd := &cobra.Command{
		Use:   "config [type]",
		Short: "Configure the kernel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
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
		},
	}

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean kernel build artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Cleaning...")
			if err := ctx.KernelBuilder.Clean(cmd.Context()); err != nil {
				return err
			}
			ctx.Printer.Success("Cleaned!")
			return nil
		},
	}

	cloneCmd := &cobra.Command{
		Use:   "clone [git-url]",
		Short: "Clone the Linux kernel source",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			if ctx.AppContext.KernelExists() {
				ctx.Printer.Info("Kernel source already exists at %s", ctx.Config.Paths.KernelDir)
				return nil
			}
			url := "https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git"
			if len(args) > 0 {
				url = args[0]
			}
			ctx.Printer.Step("Cloning kernel from %s...", url)
			if err := ctx.Exec.Run(cmd.Context(), "git", "clone", url, ctx.Config.Paths.KernelDir); err != nil {
				return fmt.Errorf("failed to clone: %w", err)
			}
			ctx.Printer.Success("Kernel cloned to %s", ctx.Config.Paths.KernelDir)
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show kernel source status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}

			if !ctx.AppContext.KernelExists() {
				ctx.Printer.Info("Kernel source not found at %s", ctx.Config.Paths.KernelDir)
				ctx.Printer.Print("  Run 'elmos kernel clone' to download kernel source")
				return nil
			}

			ctx.Printer.Success("Kernel source found at %s", ctx.Config.Paths.KernelDir)
			ctx.Printer.Print("")
			printKernelGitInfo(ctx, cmd)
			printKernelBuildStatus(ctx)
			return nil
		},
	}

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset kernel source (reclone completely)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			if ctx.AppContext.KernelExists() {
				ctx.Printer.Step("Removing existing kernel source...")
				if err := os.RemoveAll(ctx.Config.Paths.KernelDir); err != nil {
					return fmt.Errorf("failed to remove kernel: %w", err)
				}
			}
			url := "https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git"
			ctx.Printer.Step("Cloning kernel from %s...", url)
			if err := ctx.Exec.Run(cmd.Context(), "git", "clone", url, ctx.Config.Paths.KernelDir); err != nil {
				return fmt.Errorf("failed to clone: %w", err)
			}
			ctx.Printer.Success("Kernel reset complete!")
			return nil
		},
	}

	switchCmd := &cobra.Command{
		Use:   "switch [ref]",
		Short: "List or switch branch/tag (auto-detects)",
		Long: `List all branches and tags, or switch to a specific ref.
Automatically detects whether the ref is a branch or tag.

Examples:
  elmos kernel switch           # List all refs
  elmos kernel switch master    # Switch to branch
  elmos kernel switch v6.7      # Switch to tag`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			if !ctx.AppContext.KernelExists() {
				ctx.Printer.Info("Kernel source not found. Run 'elmos kernel clone' first.")
				return nil
			}

			if len(args) == 0 {
				return listKernelRefs(ctx, cmd)
			}
			return switchKernelRef(ctx, cmd, args[0])
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Update kernel source",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			if !ctx.AppContext.KernelExists() {
				ctx.Printer.Info("Kernel source not found. Run 'elmos kernel clone' first.")
				return nil
			}
			ctx.Printer.Step("Updating kernel source...")
			if err := ctx.Exec.Run(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "pull"); err != nil {
				return fmt.Errorf("failed to update: %w", err)
			}
			ctx.Printer.Success("Kernel updated!")
			return nil
		},
	}

	var jobs int
	buildCmd := &cobra.Command{
		Use:   "build [targets...]",
		Short: "Build the Linux kernel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
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
		},
	}
	buildCmd.Flags().IntVarP(&jobs, "jobs", "j", 0, "Number of parallel build jobs")

	kernelCmd.AddCommand(configCmd, cleanCmd, cloneCmd, statusCmd, resetCmd, switchCmd, pullCmd, buildCmd)
	return kernelCmd
}

// --- Helper functions to reduce RunE complexity ---

// printKernelGitInfo prints git branch/tag and commit info for status command.
func printKernelGitInfo(ctx *Context, cmd *cobra.Command) {
	ctx.Printer.Step("Git info:")
	branch, err := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "symbolic-ref", "-q", "--short", "HEAD")
	if err == nil {
		ctx.Printer.Print("  Branch: %s", strings.TrimSpace(string(branch)))
	} else {
		tag, err := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "describe", "--tags", "--exact-match")
		if err == nil {
			ctx.Printer.Print("  Tag: %s", strings.TrimSpace(string(tag)))
		} else {
			ctx.Printer.Print("  Branch: <detached>")
		}
	}
	commit, err := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "log", "-1", "--format=%h %s")
	if err == nil {
		ctx.Printer.Print("  Commit: %s", strings.TrimSpace(string(commit)))
	}
}

// printKernelBuildStatus prints kernel config and image status.
func printKernelBuildStatus(ctx *Context) {
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

// listKernelRefs lists available branches and tags for switch command.
func listKernelRefs(ctx *Context, cmd *cobra.Command) error {
	ctx.Printer.Step("Branches:")
	branches, _ := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "branch", "-a", "--format=%(refname:short)")
	for _, b := range strings.Split(string(branches), "\n") {
		if b != "" {
			ctx.Printer.Print("  %s", b)
		}
	}
	ctx.Printer.Print("")
	ctx.Printer.Step("Tags (latest 10):")
	tags, _ := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "tag", "-l", "--sort=-v:refname", "v*")
	for i, t := range strings.Split(string(tags), "\n") {
		if i >= 10 || t == "" {
			break
		}
		ctx.Printer.Print("  %s", t)
	}
	return nil
}

// switchKernelRef switches to a specific branch or tag.
func switchKernelRef(ctx *Context, cmd *cobra.Command, ref string) error {
	ctx.Printer.Step("Switching to: %s", ref)
	if err := ctx.Exec.Run(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "checkout", ref); err != nil {
		ctx.Printer.Info("Not found locally, fetching...")
		_ = ctx.Exec.Run(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "fetch", "--all", "--tags")
		if err := ctx.Exec.Run(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "checkout", ref); err != nil {
			return fmt.Errorf("failed to switch: %w", err)
		}
	}
	ctx.Printer.Success("Now on: %s", ref)
	return nil
}
