package build

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"github.com/spf13/cobra"
)

// BuildApps creates the app command tree for userspace application management.
func BuildApps(ctx *types.Context) *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Manage userspace applications",
	}

	buildCmd := &cobra.Command{
		Use:   "build [name]",
		Short: "Build apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			ctx.Printer.Step("Building apps...")
			if err := ctx.AppBuilder.Build(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Apps built!")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			apps, err := ctx.AppBuilder.GetApps("")
			if err != nil {
				return err
			}
			if len(apps) == 0 {
				ctx.Printer.Info("No apps found")
				return nil
			}
			ctx.Printer.Print("Applications:")
			for i, app := range apps {
				status := ""
				if app.Built {
					status = " (built)"
				}
				ctx.Printer.Print("  %d. %s%s", i+1, app.Name, status)
			}
			return nil
		},
	}

	newCmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Create new app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppBuilder.CreateApp(args[0]); err != nil {
				return err
			}
			ctx.Printer.Success("Created app: %s", args[0])
			return nil
		},
	}

	cleanCmd := &cobra.Command{
		Use:   "clean [name]",
		Short: "Clean app build artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			ctx.Printer.Step("Cleaning apps...")
			if err := ctx.AppBuilder.Clean(cmd.Context(), name); err != nil {
				return err
			}
			ctx.Printer.Success("Apps cleaned!")
			return nil
		},
	}

	appCmd.AddCommand(buildCmd, listCmd, newCmd, cleanCmd)
	return appCmd
}
