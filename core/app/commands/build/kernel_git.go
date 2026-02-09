package build

import (
	"fmt"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// buildKernelCloneCmd creates the kernel clone subcommand.
func buildKernelCloneCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "clone [git-url]",
		Short: "Clone the Linux kernel source",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
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
		}),
	}
}

// buildKernelStatusCmd creates the kernel status subcommand.
func buildKernelStatusCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show kernel source status",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
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
		}),
	}
}

// buildKernelResetCmd creates the kernel reset subcommand.
func buildKernelResetCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset kernel source (reclone completely)",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
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
		}),
	}
}

// buildKernelSwitchCmd creates the kernel switch subcommand.
func buildKernelSwitchCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "switch [ref]",
		Short: "List or switch branch/tag (auto-detects)",
		Long: `List all branches and tags, or switch to a specific ref.
Automatically detects whether the ref is a branch or tag.

Examples:
  elmos kernel switch           # List all refs
  elmos kernel switch master    # Switch to branch
  elmos kernel switch v6.7      # Switch to tag`,
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
			if !ctx.AppContext.KernelExists() {
				ctx.Printer.Info("Kernel source not found. Run 'elmos kernel clone' first.")
				return nil
			}
			if len(args) == 0 {
				return listKernelRefs(ctx, cmd)
			}
			return switchKernelRef(ctx, cmd, args[0])
		}),
	}
}

// buildKernelPullCmd creates the kernel pull subcommand.
func buildKernelPullCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Update kernel source",
		RunE: RunEWithContext(ctx, func(cmd *cobra.Command, args []string) error {
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
		}),
	}
}

// printKernelGitInfo prints git branch/tag and commit info for status command.
func printKernelGitInfo(ctx *types.Context, cmd *cobra.Command) {
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

// listKernelRefs lists available branches and tags for switch command.
func listKernelRefs(ctx *types.Context, cmd *cobra.Command) error {
	ctx.Printer.Step("Branches:")
	// Use for-each-ref for structured output
	branches, _ := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	for _, b := range strings.Split(string(branches), "\n") {
		if b != "" {
			ctx.Printer.Print("  %s", b)
		}
	}
	ctx.Printer.Print("")
	ctx.Printer.Step("Tags (latest 10):")
	// Optimized tag listing
	tags, _ := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "describe", "--tags", "--abbrev=0")
	// Fallback to list if describe only shows one
	if tags == nil || len(string(tags)) < 2 {
		tagsRaw, _ := ctx.Exec.Output(cmd.Context(), "git", "-C", ctx.Config.Paths.KernelDir, "tag", "-l", "--sort=-v:refname", "v*")
		lines := strings.Split(string(tagsRaw), "\n")
		for i, t := range lines {
			if i >= 10 || t == "" {
				break
			}
			ctx.Printer.Print("  %s", t)
		}
	} else {
		ctx.Printer.Print("  Latest: %s", strings.TrimSpace(string(tags)))
	}
	return nil
}

// switchKernelRef switches to a specific branch or tag.
func switchKernelRef(ctx *types.Context, cmd *cobra.Command, ref string) error {
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
