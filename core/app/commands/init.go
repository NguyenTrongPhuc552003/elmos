package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildInit creates the init command for workspace initialization.
func BuildInit(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize workspace (mount volume)",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if !ctx.AppContext.IsMounted() {
				ctx.Printer.Step("Mounting volume...")
				if err := ctx.Exec.Run(cmd.Context(), "hdiutil", "attach", ctx.Config.Image.Path); err != nil {
					return fmt.Errorf("failed to mount: %w", err)
				}
			}
			ctx.Printer.Success("Workspace initialized! Volume mounted at %s", ctx.Config.Image.MountPoint)
			return nil
		},
	}
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
			mountPoint, err := ctx.AppContext.GetActualMountPoint()
			if err != nil {
				// Fallback to config path if detection fails but IsMounted passed
				mountPoint = ctx.Config.Image.MountPoint
			}

			// Prepare args
			runArgs := []string{"detach", mountPoint}
			if force {
				runArgs = append(runArgs, "-force")
			}

			if err := ctx.Exec.Run(cmd.Context(), "hdiutil", runArgs...); err != nil {
				return fmt.Errorf("failed to unmount: %w", err)
			}
			ctx.Printer.Success("Volume unmounted from %s", mountPoint)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force unmount (needed if resource is busy)")
	return cmd
}

