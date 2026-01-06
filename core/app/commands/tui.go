package commands

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/ui/tui"
)

// BuildTUI creates the tui command for launching interactive mode.
func BuildTUI(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive Text User Interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run()
		},
	}
}
