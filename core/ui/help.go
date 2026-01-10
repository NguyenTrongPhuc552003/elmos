// Package ui provides console output helpers for elmos.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Tokyo Night color palette
var (
	headerColor  = lipgloss.Color("141") // Purple
	commandColor = lipgloss.Color("120") // Green
	flagColor    = lipgloss.Color("214") // Orange
	descColor    = lipgloss.Color("245") // Grey
	exampleColor = lipgloss.Color("51")  // Cyan
	sectionColor = lipgloss.Color("255") // White
	subtleColor  = lipgloss.Color("238") // Dark grey
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(headerColor)
	commandStyle = lipgloss.NewStyle().Foreground(commandColor)
	flagStyle    = lipgloss.NewStyle().Foreground(flagColor)
	descStyle    = lipgloss.NewStyle().Foreground(descColor)
	sectionStyle = lipgloss.NewStyle().Bold(true).Foreground(sectionColor).MarginTop(1)
	exampleStyle = lipgloss.NewStyle().Foreground(exampleColor)
	subtleStyle  = lipgloss.NewStyle().Foreground(subtleColor)
)

// Banner returns a styled ASCII art banner for ELMOS.
func Banner() string {
	banner := `
 ███████╗██╗     ███╗   ███╗ ██████╗ ███████╗
 ██╔════╝██║     ████╗ ████║██╔═══██╗██╔════╝
 █████╗  ██║     ██╔████╔██║██║   ██║███████╗
 ██╔══╝  ██║     ██║╚██╔╝██║██║   ██║╚════██║
 ███████╗███████╗██║ ╚═╝ ██║╚██████╔╝███████║
 ╚══════╝╚══════╝╚═╝     ╚═╝ ╚═════╝ ╚══════╝`
	return headerStyle.Render(banner)
}

// SetCustomUsageFunc sets a custom usage function for a Cobra command.
func SetCustomUsageFunc(cmd *cobra.Command) {
	cmd.SetUsageFunc(customUsageFunc)
	cmd.SetHelpFunc(customHelpFunc)
}

func customHelpFunc(cmd *cobra.Command, args []string) {
	var out strings.Builder

	writeHeader(&out, cmd)
	writeUsage(&out, cmd)
	writeCommands(&out, cmd)
	writeFlags(&out, cmd)
	writeExamples(&out, cmd)
	writeFooter(&out, cmd)

	_, _ = fmt.Fprint(cmd.OutOrStdout(), out.String())
}

// writeHeader writes the banner and description section.
func writeHeader(out *strings.Builder, cmd *cobra.Command) {
	if !cmd.HasParent() {
		out.WriteString(Banner())
		out.WriteString("\n\n")
	}
	if cmd.Short != "" {
		out.WriteString(headerStyle.Render(cmd.Short))
		out.WriteString("\n")
	}
	if cmd.Long != "" {
		out.WriteString(subtleStyle.Render(cmd.Long))
		out.WriteString("\n")
	}
}

// writeUsage writes the usage section.
func writeUsage(out *strings.Builder, cmd *cobra.Command) {
	out.WriteString("\n")
	out.WriteString(sectionStyle.Render("USAGE"))
	out.WriteString("\n")
	out.WriteString("  " + commandStyle.Render(cmd.UseLine()))
	out.WriteString("\n")
}

// writeCommands writes the commands section.
func writeCommands(out *strings.Builder, cmd *cobra.Command) {
	if !cmd.HasAvailableSubCommands() {
		return
	}

	out.WriteString("\n")
	out.WriteString(sectionStyle.Render("COMMANDS"))
	out.WriteString("\n")

	if !cmd.HasParent() {
		writeGroupedCommands(out, cmd.Commands())
	} else {
		writeSimpleCommands(out, cmd.Commands())
	}
}

// writeGroupedCommands writes commands organized into groups.
func writeGroupedCommands(out *strings.Builder, cmds []*cobra.Command) {
	groups := groupCommands(cmds)
	for _, group := range groups {
		if group.name != "" {
			out.WriteString("  " + subtleStyle.Render("─── "+group.name+" ───"))
			out.WriteString("\n")
		}
		for _, sub := range group.commands {
			if sub.IsAvailableCommand() {
				writeCommand(out, sub)
			}
		}
	}
}

// writeSimpleCommands writes a simple list of commands.
func writeSimpleCommands(out *strings.Builder, cmds []*cobra.Command) {
	for _, sub := range cmds {
		if sub.IsAvailableCommand() {
			writeCommand(out, sub)
		}
	}
}

// writeCommand writes a single command entry.
func writeCommand(out *strings.Builder, sub *cobra.Command) {
	name := commandStyle.Render(fmt.Sprintf("%-12s", sub.Name()))
	desc := descStyle.Render(sub.Short)
	out.WriteString(fmt.Sprintf("  %s  %s\n", name, desc))
}

// writeFlags writes the flags section.
func writeFlags(out *strings.Builder, cmd *cobra.Command) {
	if !cmd.HasAvailableLocalFlags() && !cmd.HasAvailablePersistentFlags() {
		return
	}

	out.WriteString("\n")
	out.WriteString(sectionStyle.Render("FLAGS"))
	out.WriteString("\n")

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		var name string
		if f.Shorthand != "" {
			name = flagStyle.Render(fmt.Sprintf("  -%s, --%s", f.Shorthand, f.Name))
		} else {
			name = flagStyle.Render(fmt.Sprintf("      --%s", f.Name))
		}
		desc := descStyle.Render(f.Usage)
		out.WriteString(fmt.Sprintf("%-30s  %s\n", name, desc))
	})
}

// writeExamples writes the examples section.
func writeExamples(out *strings.Builder, cmd *cobra.Command) {
	if cmd.Example == "" {
		return
	}

	out.WriteString("\n")
	out.WriteString(sectionStyle.Render("EXAMPLES"))
	out.WriteString("\n")
	for _, line := range strings.Split(cmd.Example, "\n") {
		out.WriteString("  " + exampleStyle.Render(line) + "\n")
	}
}

// writeFooter writes the help footer.
func writeFooter(out *strings.Builder, cmd *cobra.Command) {
	if !cmd.HasAvailableSubCommands() {
		return
	}

	out.WriteString("\n")
	out.WriteString(subtleStyle.Render(fmt.Sprintf("Use \"%s [command] --help\" for more information about a command.", cmd.CommandPath())))
	out.WriteString("\n")
}

func customUsageFunc(cmd *cobra.Command) error {
	customHelpFunc(cmd, nil)
	return nil
}

// commandGroup represents a group of related commands.
type commandGroup struct {
	name     string
	commands []*cobra.Command
}

// groupCommands organizes commands into logical groups.
func groupCommands(cmds []*cobra.Command) []commandGroup {
	core := []*cobra.Command{}
	build := []*cobra.Command{}
	runtime := []*cobra.Command{}
	config := []*cobra.Command{}
	other := []*cobra.Command{}

	for _, cmd := range cmds {
		if !cmd.IsAvailableCommand() {
			continue
		}
		switch cmd.Name() {
		case "init", "exit", "doctor", "version", "tui", "status", "arch":
			core = append(core, cmd)
		case "kernel", "module", "app", "rootfs", "patch":
			build = append(build, cmd)
		case "qemu", "gdb":
			runtime = append(runtime, cmd)
		default:
			other = append(other, cmd)
		}
	}

	groups := []commandGroup{}
	if len(core) > 0 {
		groups = append(groups, commandGroup{name: "Core", commands: core})
	}
	if len(build) > 0 {
		groups = append(groups, commandGroup{name: "Build", commands: build})
	}
	if len(runtime) > 0 {
		groups = append(groups, commandGroup{name: "Runtime", commands: runtime})
	}
	if len(config) > 0 {
		groups = append(groups, commandGroup{name: "Config", commands: config})
	}
	if len(other) > 0 {
		groups = append(groups, commandGroup{name: "", commands: other})
	}

	return groups
}
