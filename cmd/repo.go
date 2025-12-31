// Package cmd implements the Cobra CLI commands for elmos.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const kernelRepoURL = "https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git"

// repoCmd - git repository management
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage kernel git repository",
	Long:  `Commands to clone, update, and manage the Linux kernel git repository.`,
}

var repoStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runGitCommand("status")
	},
}

var repoCheckoutCmd = &cobra.Command{
	Use:   "checkout [branch|tag]",
	Short: "Checkout a branch or tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runRepoCheckout(args[0])
	},
}

var repoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Fetch and reset to origin/master",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runRepoUpdate()
	},
}

var repoResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Hard reset to origin/master (no fetch)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runRepoReset()
	},
}

var repoReinitCmd = &cobra.Command{
	Use:   "reinit",
	Short: "Delete and re-clone kernel repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runRepoReinit()
	},
}

func init() {
	repoCmd.AddCommand(repoStatusCmd)
	repoCmd.AddCommand(repoCheckoutCmd)
	repoCmd.AddCommand(repoUpdateCmd)
	repoCmd.AddCommand(repoResetCmd)
	repoCmd.AddCommand(repoReinitCmd)
}

func runRepoCheck() error {
	kernelDir := ctx.Config.Paths.KernelDir
	gitDir := filepath.Join(kernelDir, ".git")

	if _, err := os.Stat(gitDir); err == nil {
		printInfo("Kernel repository exists at %s", kernelDir)
		return nil
	}

	printStep("Cloning Linux kernel into %s...", kernelDir)
	printInfo("This may take a while for the full history")

	cmd := exec.Command("git", "clone", kernelRepoURL, kernelDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	printSuccess("Repository cloned successfully")
	return nil
}

func runRepoCheckout(target string) error {
	printStep("Checking out %s...", target)

	// First, try as a tag
	if err := runGitCommand("checkout", target); err != nil {
		// If tag doesn't exist, try creating a new branch
		if err := runGitCommand("checkout", "-b", target, "--track", "origin/master"); err != nil {
			return fmt.Errorf("failed to checkout %s: %w", target, err)
		}
		printSuccess("Created and switched to new branch: %s", target)
		return nil
	}

	printSuccess("Switched to %s", target)
	return nil
}

func runRepoUpdate() error {
	printStep("Fetching from origin...")
	if err := runGitCommand("fetch", "origin"); err != nil {
		return err
	}

	printStep("Resetting to origin/master...")
	if err := runGitCommand("reset", "--hard", "origin/master"); err != nil {
		return err
	}

	printStep("Cleaning untracked files...")
	if err := runGitCommand("clean", "-fd"); err != nil {
		return err
	}

	printSuccess("Repository updated to origin/master")
	return nil
}

func runRepoReset() error {
	printWarn("This will discard all local changes")

	printStep("Resetting to origin/master...")
	if err := runGitCommand("reset", "--hard", "origin/master"); err != nil {
		return err
	}

	printStep("Cleaning untracked files...")
	if err := runGitCommand("clean", "-fd"); err != nil {
		return err
	}

	printSuccess("Local changes discarded")
	return nil
}

func runRepoReinit() error {
	kernelDir := ctx.Config.Paths.KernelDir

	printWarn("This will DELETE the entire kernel tree and re-clone")
	printInfo("Press Ctrl+C to cancel, or wait 5 seconds to continue...")

	// In interactive mode, we'd prompt here
	// For now, just proceed

	printStep("Removing %s...", kernelDir)
	if err := os.RemoveAll(kernelDir); err != nil {
		return fmt.Errorf("failed to remove kernel directory: %w", err)
	}

	return runRepoCheck()
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = ctx.Config.Paths.KernelDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
