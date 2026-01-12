package commands

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/version"
)

// BuildVersion creates the version command.
func BuildVersion(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			info := version.Get()
			ctx.Printer.Print("ELMOS - Embedded Linux on MacOS")
			ctx.Printer.Print("Version:    %s", ui.AccentStyle.Render(info.Version))
			ctx.Printer.Print("Commit:     %s", info.Commit)
			ctx.Printer.Print("Built:      %s", info.BuildDate)
			ctx.Printer.Print("Go version: %s", info.GoVersion)
			ctx.Printer.Print("OS/Arch:    %s/%s", info.OS, info.Arch)
		},
	}
}
