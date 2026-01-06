package commands

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/emulator"
)

// BuildQEMU creates the qemu command tree for QEMU emulation.
func BuildQEMU(ctx *Context) *cobra.Command {
	qemuCmd := &cobra.Command{
		Use:   "qemu",
		Short: "Run and debug kernel in QEMU",
	}
	var graphical bool

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run kernel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Starting QEMU...")
			return ctx.QEMURunner.Run(cmd.Context(), emulator.RunOptions{Graphical: graphical})
		},
	}
	runCmd.Flags().BoolVarP(&graphical, "graphical", "g", false, "Graphical mode")

	debugCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug kernel with GDB server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Starting QEMU in debug mode...")
			return ctx.QEMURunner.Debug(cmd.Context(), graphical)
		},
	}

	qemuCmd.AddCommand(runCmd, debugCmd)
	return qemuCmd
}

// BuildGDB creates the gdb command for connecting to QEMU debug session.
func BuildGDB(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "gdb",
		Short: "Connect GDB to running QEMU debug session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ctx.QEMURunner.ConnectGDB()
		},
	}
}
