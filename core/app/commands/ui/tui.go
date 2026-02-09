package ui

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/ui/terminal"
)

// BuildTUI creates the tui command for launching interactive mode.
func BuildTUI(ctx *types.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive Text User Interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminal.Run()
		},
	}
}
