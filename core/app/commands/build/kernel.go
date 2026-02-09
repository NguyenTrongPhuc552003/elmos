package build

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"github.com/spf13/cobra"
)

// BuildKernel creates the kernel command tree for kernel management.
func BuildKernel(ctx *types.Context) *cobra.Command {
	kernelCmd := &cobra.Command{
		Use:   "kernel",
		Short: "Kernel configuration commands",
	}

	kernelCmd.AddCommand(
		buildKernelConfigCmd(ctx),
		buildKernelCleanCmd(ctx),
		buildKernelCloneCmd(ctx),
		buildKernelStatusCmd(ctx),
		buildKernelResetCmd(ctx),
		buildKernelSwitchCmd(ctx),
		buildKernelPullCmd(ctx),
		buildKernelBuildCmd(ctx),
	)

	return kernelCmd
}

// RunEWithContext is a helper to run commands that require the AppContext to be mounted.
func RunEWithContext(ctx *types.Context, run func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := ctx.AppContext.EnsureMounted(); err != nil {
			return err
		}
		return run(cmd, args)
	}
}
