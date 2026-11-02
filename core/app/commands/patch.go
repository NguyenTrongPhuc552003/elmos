package commands

import (
	"github.com/spf13/cobra"

	"codeberg.org/NguyenTrongPhuc552003/elmos/core/domain/patch"
)

// BuildPatch creates the patch command tree for kernel patch management.
func BuildPatch(ctx *Context) *cobra.Command {
	patchCmd := &cobra.Command{
		Use:   "patch",
		Short: "Manage kernel patches",
	}

	applyCmd := &cobra.Command{
		Use:   "apply [file]",
		Short: "Apply patch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.AppContext.EnsureMounted(); err != nil {
				return err
			}
			ctx.Printer.Step("Applying patch: %s", args[0])
			if err := ctx.PatchManager.Apply(cmd.Context(), args[0]); err != nil {
				return err
			}
			ctx.Printer.Success("Patch applied!")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List patches",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := ctx.PatchManager.List()
			if err != nil {
				return err
			}
			builtin := patch.AllBuiltin()
			if len(builtin) > 0 {
				ctx.Printer.Print("Built-in patches (applied by 'elmos kernel clone'):")
				seen := map[string]bool{}
				for _, p := range builtin {
					if !seen[p.Path] {
						seen[p.Path] = true
						ctx.Printer.Print("  %s/%s/%s", p.Version, p.Arch, p.Name)
					}
				}
			}
			if len(patches) == 0 && len(builtin) == 0 {
				ctx.Printer.Info("No patches")
				return nil
			}
			if len(patches) == 0 {
				return nil
			}
			ctx.Printer.Print("Project patches (%s):", ctx.Config.Paths.PatchesDir)
			for _, p := range patches {
				ctx.Printer.Print("  %s/%s/%s", p.Version, p.Arch, p.Name)
			}
			return nil
		},
	}

	patchCmd.AddCommand(applyCmd, listCmd)
	return patchCmd
}
