package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// BuildStatus creates the status command for workspace status display.
func BuildStatus(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show workspace status (volume mount info)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if mounted
			if !ctx.AppContext.IsMounted() {
				ctx.Printer.Info("Workspace not mounted")
				return nil
			}

			// Get actual mount point
			mountPoint, err := ctx.AppContext.GetActualMountPoint()
			if err != nil {
				mountPoint = ctx.Config.Image.MountPoint
			}

			ctx.Printer.Success("Workspace mounted at %s", mountPoint)
			ctx.Printer.Print("")
			ctx.Printer.Step("Volume info:")

			// Run hdiutil info and display relevant parts
			out, err := ctx.Exec.Output(cmd.Context(), "hdiutil", "info")
			if err != nil {
				return fmt.Errorf("failed to get hdiutil info: %w", err)
			}

			// Print the output (filtered to our image)
			lines := strings.Split(string(out), "\n")
			inOurImage := false
			for _, line := range lines {
				if strings.Contains(line, ctx.Config.Image.Path) {
					inOurImage = true
				}
				if inOurImage {
					ctx.Printer.Print("  %s", line)
					if strings.HasPrefix(line, "/dev/disk") && strings.Contains(line, "/Volumes/") {
						break
					}
				}
			}

			return nil
		},
	}
}
