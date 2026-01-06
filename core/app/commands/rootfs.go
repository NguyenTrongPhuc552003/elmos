package commands

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/rootfs"
)

// BuildRootfs creates the rootfs command tree for root filesystem management.
func BuildRootfs(ctx *Context) *cobra.Command {
	rootfsCmd := &cobra.Command{
		Use:   "rootfs",
		Short: "Manage root filesystem",
	}
	var size string

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create rootfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Creating rootfs...")
			if err := ctx.RootfsCreator.Create(cmd.Context(), rootfs.CreateOptions{Size: size}); err != nil {
				return err
			}
			ctx.Printer.Success("Rootfs created!")
			return nil
		},
	}
	createCmd.Flags().StringVarP(&size, "size", "s", "5G", "Disk size")

	rootfsCmd.AddCommand(createCmd)
	return rootfsCmd
}
