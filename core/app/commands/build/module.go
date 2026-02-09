package build

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"github.com/spf13/cobra"
)

// BuildModule creates the module command tree for kernel module management.
func BuildModule(ctx *types.Context) *cobra.Command {
	modCmd := &cobra.Command{
		Use:   "module",
		Short: "Manage kernel modules",
	}

	modCmd.AddCommand(
		buildModuleBuildCmd(ctx),
		buildModuleListCmd(ctx),
		buildModuleNewCmd(ctx),
		buildModuleCleanCmd(ctx),
		buildModuleHeaderCmd(ctx),
	)
	return modCmd
}

// buildModuleBuildCmd creates the module build subcommand.
func buildModuleBuildCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "build [name]",
		Short: "Build kernel modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			name := getOptionalArg(args)
			ctx.Printer.Step("Building modules...")
			if err := ctx.ModuleBuilder.Build(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Modules built!")
			return nil
		},
	}
}

// buildModuleListCmd creates the module list subcommand.
func buildModuleListCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
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
}

// buildModuleNewCmd creates the module new subcommand.
func buildModuleNewCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
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
}

// buildModuleCleanCmd creates the module clean subcommand.
func buildModuleCleanCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "clean [name]",
		Short: "Clean module build artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			name := getOptionalArg(args)
			ctx.Printer.Step("Cleaning modules...")
			if err := ctx.ModuleBuilder.Clean(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Modules cleaned!")
			return nil
		},
	}
}

// buildModuleHeaderCmd creates the module header subcommand.
func buildModuleHeaderCmd(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "header",
		Short: "Prepare kernel headers for module building",
		Long:  "Runs 'make modules_prepare' to set up kernel headers required for external module compilation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Preparing kernel headers...")
			if err := ctx.ModuleBuilder.PrepareHeaders(cmd.Context()); err != nil {
				return err
			}
			ctx.Printer.Success("Kernel headers prepared!")
			return nil
		},
	}
}

// getOptionalArg returns the first argument or empty string if none provided.
func getOptionalArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
