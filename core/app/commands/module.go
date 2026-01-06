package commands

import (
	"github.com/spf13/cobra"
)

// BuildModule creates the module command tree for kernel module management.
func BuildModule(ctx *Context) *cobra.Command {
	modCmd := &cobra.Command{
		Use:   "module",
		Short: "Manage kernel modules",
	}

	buildCmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Build kernel modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			ctx.Printer.Step("Building modules...")
			if err := ctx.ModuleBuilder.Build(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Modules built!")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			mods, err := ctx.ModuleBuilder.GetModules("")
			if err != nil {
				return err
			}
			if len(mods) == 0 {
				ctx.Printer.Info("No modules found")
				return nil
			}
			ctx.Printer.Print("Modules:")
			for i, m := range mods {
				status := ""
				if m.Built {
					status = " (built)"
				}
				ctx.Printer.Print("  %d. %s%s", i+1, m.Name, status)
			}
			return nil
		},
	}

	newCmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Create new module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.ModuleBuilder.CreateModule(args[0]); err != nil {
				return err
			}
			ctx.Printer.Success("Created module: %s", args[0])
			return nil
		},
	}

	cleanCmd := &cobra.Command{
		Use:   "clean [name]",
		Short: "Clean module build artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			ctx.Printer.Step("Cleaning modules...")
			if err := ctx.ModuleBuilder.Clean(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Modules cleaned!")
			return nil
		},
	}

	modCmd.AddCommand(buildCmd, listCmd, newCmd, cleanCmd)
	return modCmd
}
