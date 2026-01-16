package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
	"github.com/spf13/cobra"
)

// BuildInit creates the init command for workspace initialization.
func BuildInit(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "init [workspace_name] [size]",
		Short: "Initialize workspace (mount volume)",
		Long: `Initialize workspace and mount volume.

Arguments:
  workspace_name  Optional name for the workspace volume (default: "elmos")
  size           Optional volume size (default: "40G", minimum: 40G)

Examples:
  elmos init                    # Create /Volumes/elmos/ with 40GB
  elmos init my_workspace       # Create /Volumes/my_workspace/ with 40GB
  elmos init my_workspace 50G   # Create /Volumes/my_workspace/ with 50GB`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(ctx, cmd, args)
		},
	}
}

// runInit executes the init command logic.
func runInit(ctx *Context, cmd *cobra.Command, args []string) error {
	// Parse arguments
	workspaceName := ctx.Config.Image.VolumeName
	if workspaceName == "" {
		workspaceName = config.DefaultVolumeName
	}
	volumeSize := ctx.Config.Image.Size
	if volumeSize == "" {
		volumeSize = config.DefaultImageSize
	}

	if len(args) > 0 {
		workspaceName = args[0]
	}
	if len(args) > 1 {
		volumeSize = args[1]
		// Validate size
		if err := validateVolumeSize(volumeSize, ctx.Printer); err != nil {
			return err
		}
	}

	// Check/Update configuration
	if err := updateInitConfig(ctx, workspaceName, volumeSize); err != nil {
		return err
	}

	// Create and mount volume
	if err := ensureWorkspaceVolume(ctx, cmd); err != nil {
		return err
	}

	ctx.Printer.Success("Workspace initialized! Volume mounted at %s", ctx.Config.Image.MountPoint)
	return nil
}

// ensureWorkspaceVolume creates and mounts the disk image if needed.
func ensureWorkspaceVolume(ctx *Context, cmd *cobra.Command) error {
	// Create disk image if it doesn't exist
	if !ctx.FS.Exists(ctx.Config.Image.Path) {
		ctx.Printer.Step("Creating sparse disk image...")
		if err := ctx.Exec.Run(cmd.Context(), "hdiutil", "create",
			"-size", ctx.Config.Image.Size,
			"-fs", "Case-sensitive APFS",
			"-volname", ctx.Config.Image.VolumeName,
			"-type", "SPARSE",
			ctx.Config.Image.Path,
		); err != nil {
			return fmt.Errorf("failed to create disk image: %w", err)
		}
		ctx.Printer.Success("Disk image created!")
	}

	// Mount volume if not already mounted
	if !ctx.AppContext.IsMounted() {
		ctx.Printer.Step("Mounting volume...")
		if err := ctx.Exec.Run(cmd.Context(), "hdiutil", "attach",
			"-mountpoint", ctx.Config.Image.MountPoint,
			ctx.Config.Image.Path,
		); err != nil {
			return fmt.Errorf("failed to mount: %w", err)
		}
	}
	return nil
}

// updateInitConfig updates the configuration and saves it if necessary.
func updateInitConfig(ctx *Context, workspaceName, volumeSize string) error {
	// Check if we need to update config
	configChanged := false
	if ctx.Config.Image.VolumeName != workspaceName {
		ctx.Config.Image.VolumeName = workspaceName
		configChanged = true
	}
	if ctx.Config.Image.Size != volumeSize {
		ctx.Config.Image.Size = volumeSize
		configChanged = true
	}

	// Update derived paths
	ctx.Config.Image.MountPoint = fmt.Sprintf("/Volumes/%s", workspaceName)
	ctx.Config.Image.Path = fmt.Sprintf("%s/data/%s.sparseimage",
		ctx.Config.Paths.ProjectRoot, workspaceName)
	ctx.Config.Paths.ToolchainsDir = fmt.Sprintf("/Volumes/%s/toolchains", workspaceName)

	// Determine config file path
	configPath := ctx.Config.ConfigFile
	if configPath == "" {
		configPath = fmt.Sprintf("%s/elmos.yaml", ctx.Config.Paths.ProjectRoot)
	}

	// Save only if changed OR file doesn't exist
	shouldSave := configChanged
	if !ctx.FS.Exists(configPath) {
		shouldSave = true
	}

	if shouldSave {
		if err := ctx.Config.Save(configPath); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}
	return nil
}

// validateVolumeSize checks if the provided size meets minimum requirements.
func validateVolumeSize(size string, printer *ui.Printer) error {
	// Parse size string (e.g., "40G", "50G", "1T")
	var numericValue int
	var unit string

	if _, err := fmt.Sscanf(size, "%d%s", &numericValue, &unit); err != nil {
		return fmt.Errorf("invalid size format: %s (expected format: 40G, 50G, etc.)", size)
	}

	// Convert to GB for comparison
	sizeInGB := numericValue
	switch unit {
	case "G", "g":
		// Already in GB
	case "T", "t":
		sizeInGB = numericValue * 1024
	case "M", "m":
		sizeInGB = numericValue / 1024
	default:
		return fmt.Errorf("invalid size unit: %s (use G for gigabytes or T for terabytes)", unit)
	}

	// Validate minimum size
	if sizeInGB < config.MinimumImageSize {
		printer.Warn("⚠️  Volume size %s is less than the recommended minimum of %dG", size, config.MinimumImageSize)
		printer.Warn("   This may cause issues with toolchain builds and kernel compilation")
		printer.Warn("   Consider using at least %dG for optimal performance", config.MinimumImageSize)
	}

	return nil
}

// BuildExit creates the exit command for unmounting the workspace.
func BuildExit(ctx *Context) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "exit",
		Short: "Exit workspace (unmount volume)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.AppContext.IsMounted() {
				ctx.Printer.Info("Volume not mounted")
				return nil
			}
			ctx.Printer.Step("Unmounting volume...")

			// Find the disk device for our image from hdiutil info
			diskDevice, err := findDiskDevice(ctx)
			if err != nil {
				return fmt.Errorf("failed to find disk device: %w", err)
			}

			// Prepare args - use disk device which is reliable
			runArgs := []string{"detach", diskDevice}
			if force {
				runArgs = append(runArgs, "-force")
			}

			if err := ctx.Exec.Run(cmd.Context(), "hdiutil", runArgs...); err != nil {
				return fmt.Errorf("failed to unmount: %w", err)
			}
			ctx.Printer.Success("Volume unmounted")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force unmount (needed if resource is busy)")
	return cmd
}

// findDiskDevice finds the disk device for our mounted image from hdiutil info.
// Returns the disk device path like "/dev/disk4" that can be used with hdiutil detach.
func findDiskDevice(ctx *Context) (string, error) {
	out, err := ctx.Exec.Output(context.Background(), "hdiutil", "info")
	if err != nil {
		return "", err
	}

	return parseDiskDeviceFromHdiutil(string(out), ctx.Config.Image.Path)
}

// parseDiskDeviceFromHdiutil extracts the disk device for an image from hdiutil info output.
func parseDiskDeviceFromHdiutil(output, imagePath string) (string, error) {
	// Simple O(N) regex search instead of multi-pass line parsing
	// Look for: image-path ... /dev/diskXsY or /dev/diskX
	// But hdiutil output structure is:
	// image-path: ...
	// /dev/disk...
	//
	// We need to find the block for our image.

	// Split by block for safety
	blocks := strings.Split(output, "===")
	for _, block := range blocks {
		if strings.Contains(block, imagePath) {
			re := regexp.MustCompile(`/dev/disk\d+`)
			if match := re.FindString(block); match != "" {
				return match, nil
			}
		}
	}

	return "", fmt.Errorf("disk device not found for image: %s", imagePath)
}
