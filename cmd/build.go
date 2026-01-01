// Package cmd implements the Cobra CLI commands for elmos.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/NguyenTrongPhuc552003/elmos/internal/core"
	"github.com/spf13/cobra"
)

// kernelCmd - kernel configuration
var kernelCmd = &cobra.Command{
	Use:   "kernel",
	Short: "Kernel configuration commands",
	Long:  `Commands to configure the Linux kernel (defconfig, menuconfig, etc).`,
}

var kernelConfigCmd = &cobra.Command{
	Use:   "config [type]",
	Short: "Run kernel configuration",
	Long: `Configure the kernel. Types:
  defconfig       - Default configuration (default)
  menuconfig      - Interactive menu
  allnoconfig     - Minimal configuration
  kvm_guest.config - KVM guest support`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}

		configType := "defconfig"
		if len(args) > 0 {
			configType = args[0]
		}

		return runKernelConfig(configType)
	},
}

var kernelCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean kernel build artifacts (distclean)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runKernelClean()
	},
}

func init() {
	kernelCmd.AddCommand(kernelConfigCmd)
	kernelCmd.AddCommand(kernelCleanCmd)
}

func runKernelConfig(configType string) error {
	cfg := ctx.Config

	printStep("Running 'make %s' for ARCH=%s...", configType, cfg.Build.Arch)

	cmd := exec.Command("make",
		fmt.Sprintf("ARCH=%s", cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", cfg.Build.CrossCompile),
		configType,
	)
	cmd.Dir = cfg.Paths.KernelDir
	cmd.Env = ctx.GetMakeEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("configuration failed: %w", err)
	}

	// Additional config for KVM guest
	if configType == "kvm_guest.config" {
		if err := enableKVMConfig(); err != nil {
			printWarn("Failed to enable KVM configs: %v", err)
		}
	}

	printSuccess("Configuration complete")
	return nil
}

func enableKVMConfig() error {
	cfg := ctx.Config
	configScript := fmt.Sprintf("%s/scripts/config", cfg.Paths.KernelDir)

	cmd := exec.Command(configScript, "--file", ".config",
		"--enable", "CONFIG_DRM",
		"--enable", "CONFIG_DRM_VIRTIO_GPU",
		"--enable", "CONFIG_FB",
		"--enable", "CONFIG_FRAMEBUFFER_CONSOLE",
	)
	cmd.Dir = cfg.Paths.KernelDir

	return cmd.Run()
}

func runKernelClean() error {
	cfg := ctx.Config

	printStep("Running 'make distclean'...")

	cmd := exec.Command("make",
		fmt.Sprintf("ARCH=%s", cfg.Build.Arch),
		"distclean",
	)
	cmd.Dir = cfg.Paths.KernelDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clean failed: %w", err)
	}

	printSuccess("Kernel tree cleaned")
	return nil
}

// buildCmd - kernel build
var buildCmd = &cobra.Command{
	Use:   "build [targets...]",
	Short: "Build the Linux kernel",
	Long: `Build the Linux kernel and optional targets.

Default targets: Image dtbs modules

Examples:
  elmos build                    # Build Image, dtbs, modules
  elmos build -j8               # Build with 8 parallel jobs
  elmos build modules_prepare   # Only prepare for module building`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}

		// Get jobs flag
		jobs, _ := cmd.Flags().GetInt("jobs")
		if jobs == 0 {
			jobs = ctx.Config.Build.Jobs
		}

		// Get targets or use defaults (ARM32 needs zImage instead of Image)
		targets := args
		if len(targets) == 0 {
			if ctx.Config.Build.Arch == "arm" {
				targets = []string{"zImage", "dtbs", "modules"}
			} else {
				targets = []string{"Image", "dtbs", "modules"}
			}
		}

		return runBuild(jobs, targets)
	},
}

func init() {
	buildCmd.Flags().IntP("jobs", "j", 0, "Number of parallel jobs (default: auto)")
}

func runBuild(jobs int, targets []string) error {
	cfg := ctx.Config

	// Validate build targets
	for _, t := range targets {
		if !validBuildTargets[t] {
			return fmt.Errorf("invalid build target: %s", t)
		}
	}

	// Check for .config
	if !ctx.HasConfig() {
		return fmt.Errorf("kernel not configured - run 'elmos kernel config' first")
	}

	printStep("Building kernel for ARCH=%s with %d jobs...", cfg.Build.Arch, jobs)
	printInfo("Targets: %v", targets)

	// Build make arguments
	makeArgs := []string{
		fmt.Sprintf("-j%d", jobs),
		fmt.Sprintf("ARCH=%s", cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", cfg.Build.CrossCompile),
	}
	makeArgs = append(makeArgs, targets...)

	cmd := exec.Command("make", makeArgs...)
	cmd.Dir = cfg.Paths.KernelDir
	cmd.Env = ctx.GetMakeEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	printSuccess("Build complete!")

	// Show output paths
	if ctx.HasKernelImage() {
		printInfo("Kernel image: %s", ctx.GetKernelImage())
	}

	return nil
}

// patchCmd - patch management
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Manage kernel patches",
	Long:  `Apply and manage kernel patches.`,
}

var patchApplyCmd = &cobra.Command{
	Use:   "apply [patch-file]",
	Short: "Apply a patch file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		return runPatchApply(args[0])
	},
}

var patchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available patches",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPatchList()
	},
}

func init() {
	patchCmd.AddCommand(patchApplyCmd)
	patchCmd.AddCommand(patchListCmd)
}

func runPatchApply(patchFile string) error {
	cfg := ctx.Config

	// Resolve patch path
	fullPath := patchFile
	if patchFile[0] != '/' {
		fullPath = fmt.Sprintf("%s/%s", cfg.Paths.ProjectRoot, patchFile)
	}

	// Check file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("patch file not found: %s", fullPath)
	}

	printStep("Checking patch applicability...")

	// Clean any aborted am session
	exec.Command("git", "am", "--abort").Run()

	// Test with dry-run
	testCmd := exec.Command("git", "am", "--3way", "--dry-run", fullPath)
	testCmd.Dir = cfg.Paths.KernelDir
	if err := testCmd.Run(); err != nil {
		printWarn("git am dry-run failed, trying git apply...")
		testCmd2 := exec.Command("git", "apply", "--3way", "--check", fullPath)
		testCmd2.Dir = cfg.Paths.KernelDir
		if err := testCmd2.Run(); err != nil {
			return fmt.Errorf("patch cannot be applied cleanly")
		}
	}

	printStep("Applying patch...")

	applyCmd := exec.Command("git", "am", "--3way", "--signoff", fullPath)
	applyCmd.Dir = cfg.Paths.KernelDir
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Run(); err != nil {
		printError("Patch application failed")
		printInfo("Run 'git am --abort' to cancel, or 'git am --continue' after resolving")
		return err
	}

	printSuccess("Patch applied successfully")
	return nil
}

func runPatchList() error {
	cfg := ctx.Config

	entries, err := os.ReadDir(cfg.Paths.PatchesDir)
	if err != nil {
		return fmt.Errorf("failed to read patches directory: %w", err)
	}

	if len(entries) == 0 {
		printInfo("No patches found in %s", cfg.Paths.PatchesDir)
		return nil
	}

	fmt.Println("Available patches:")
	for i, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  %d. %s/\n", i+1, entry.Name())
			// List patches in subdirectory
			subPath := fmt.Sprintf("%s/%s", cfg.Paths.PatchesDir, entry.Name())
			subEntries, _ := os.ReadDir(subPath)
			for _, sub := range subEntries {
				if !sub.IsDir() {
					fmt.Printf("       - %s\n", sub.Name())
				}
			}
		} else {
			fmt.Printf("  %d. %s\n", i+1, entry.Name())
		}
	}

	return nil
}

// rootfsCmd - rootfs management
var rootfsCmd = &cobra.Command{
	Use:   "rootfs",
	Short: "Manage root filesystem",
	Long:  `Create and manage the Debian root filesystem for QEMU.`,
}

var rootfsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ext4 disk image with Debian rootfs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ctx.EnsureMounted(); err != nil {
			return err
		}
		size, _ := cmd.Flags().GetString("size")
		return runRootfsCreate(size)
	},
}

func init() {
	rootfsCmd.AddCommand(rootfsCreateCmd)
	rootfsCreateCmd.Flags().StringP("size", "s", "5G", "Disk image size")
}

func runRootfsCreate(size string) error {
	cfg := ctx.Config

	// Map architecture for debootstrap
	archMap := map[string]string{
		"arm64": "arm64",
		"riscv": "riscv64",
		"arm":   "armhf",
	}

	debArch, ok := archMap[cfg.Build.Arch]
	if !ok {
		return fmt.Errorf("unsupported architecture for debootstrap: %s", cfg.Build.Arch)
	}

	printStep("Creating Debian rootfs for %s...", debArch)

	// Create rootfs directory
	os.RemoveAll(cfg.Paths.RootfsDir)
	os.MkdirAll(cfg.Paths.RootfsDir, 0755)

	// Run debootstrap
	debootstrapDir := fmt.Sprintf("%s/tools/debootstrap", cfg.Paths.ProjectRoot)
	debootstrapPath := filepath.Join(debootstrapDir, "debootstrap")

	// Check if debootstrap exists, clone if not
	if _, err := os.Stat(debootstrapPath); os.IsNotExist(err) {
		printStep("Debootstrap not found. Cloning from upstream...")
		if err := os.MkdirAll(filepath.Dir(debootstrapDir), 0755); err != nil {
			return fmt.Errorf("failed to create tools directory: %w", err)
		}

		cloneCmd := exec.Command("git", "clone", "--depth=1", "https://salsa.debian.org/installer-team/debootstrap.git", debootstrapDir)
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr
		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("failed to clone debootstrap: %w", err)
		}
		printSuccess("Debootstrap cloned")
	}

	printStep("Running debootstrap stage 1...")

	// Build environment with proper PATH (matching common.env)
	env := ctx.GetMakeEnv() // This includes gnu-sed, llvm, e2fsprogs, coreutils in PATH
	env = append(env, fmt.Sprintf("DEBOOTSTRAP_DIR=%s", debootstrapDir))

	// Run debootstrap exactly like original: sudo env ... fakeroot debootstrap ...
	cmd := exec.Command("sudo", "-E",
		"fakeroot", debootstrapPath,
		"--foreign",
		"--arch="+debArch,
		"--no-check-gpg",
		"stable",
		cfg.Paths.RootfsDir,
		cfg.Paths.DebianMirror,
	)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("debootstrap failed: %w", err)
	}

	// Create init script
	if err := createInitScript(cfg.Paths.RootfsDir); err != nil {
		return err
	}

	// Create disk image
	printStep("Creating ext4 disk image (%s)...", size)

	// Get mke2fs path from Homebrew e2fsprogs
	e2fsSbin := core.GetBrewSbin("e2fsprogs")
	mke2fsPath := "mke2fs"
	if e2fsSbin != "" {
		mke2fsPath = filepath.Join(e2fsSbin, "mke2fs")
	}

	mke2fsCmd := exec.Command(mke2fsPath,
		"-t", "ext4",
		"-E", "lazy_itable_init=0,lazy_journal_init=0",
		"-d", cfg.Paths.RootfsDir,
		cfg.Paths.DiskImage,
		size,
	)
	mke2fsCmd.Stdout = os.Stdout
	mke2fsCmd.Stderr = os.Stderr

	if err := mke2fsCmd.Run(); err != nil {
		return fmt.Errorf("mke2fs failed: %w", err)
	}

	printSuccess("Disk image created: %s", cfg.Paths.DiskImage)
	return nil
}

func createInitScript(rootfsDir string) error {
	initContent := `#!/bin/sh

MARKER="/.rootfs-setup-complete"

echo "Booting Debian root filesystem..."

if [ ! -f "$MARKER" ]; then
    echo "First boot detected – running debootstrap second stage..."
    /debootstrap/debootstrap --second-stage
    if [ $? -eq 0 ]; then
        touch "$MARKER"
        echo "Second stage completed successfully."
    else
        echo "Second stage failed – dropping to emergency shell."
        exec /bin/sh
    fi
else
    echo "Root filesystem already set up."
fi

# Mount essential virtual filesystems
mount -t proc  proc  /proc
mount -t sysfs sys   /sys
mount -t devtmpfs dev /dev 2>/dev/null || mount -t tmpfs dev /dev
[ -d /dev/pts ] || mkdir /dev/pts
mount -t devpts devpts /dev/pts

# Network configuration
ip link set lo up
ip link set eth0 up
ip addr add 10.0.2.15/24 dev eth0
ip route add default via 10.0.2.2
echo "nameserver 8.8.8.8" > /etc/resolv.conf

# Mount 9p share for modules
mkdir -p /mnt/modules
mount -t 9p -o trans=virtio,version=9p2000.L modules_mount /mnt/modules 2>/dev/null

# Execute module sync script if present
if [ -f /mnt/modules/guesync.sh ]; then
    /mnt/modules/guesync.sh
fi

echo "System ready."
exec /bin/sh
`

	initPath := fmt.Sprintf("%s/init", rootfsDir)
	if err := os.WriteFile(initPath, []byte(initContent), 0755); err != nil {
		return fmt.Errorf("failed to create init script: %w", err)
	}

	return nil
}

// Validate targets
var validBuildTargets = map[string]bool{
	"Image":           true,
	"zImage":          true, // ARM32
	"dtbs":            true,
	"modules":         true,
	"modules_prepare": true,
	"all":             true,
	"vmlinux":         true,
}
