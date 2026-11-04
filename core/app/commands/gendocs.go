package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// BuildGenDocs creates a hidden command that writes the command reference
// as Markdown pages with Starlight front matter, one page per command. The
// documentation site runs it at build time, so the reference never drifts
// from the code (audit E5).
func BuildGenDocs(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "gendocs DIR",
		Short:  "Write the command reference as Markdown (for the docs site)",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(args[0], 0o755); err != nil {
				return err
			}
			return writeCommandDocs(root, args[0])
		},
	}
}

func commandSlug(c *cobra.Command) string {
	return strings.ReplaceAll(c.CommandPath(), " ", "-")
}

func writeCommandDocs(c *cobra.Command, dir string) error {
	if c.Hidden || c.Name() == "help" {
		return nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\ntitle: %q\ndescription: %q\n---\n\n", c.CommandPath(), c.Short)
	if c.Long != "" {
		if strings.Contains(c.Long, "\n  ") {
			// indented examples and tables: keep the layout
			fmt.Fprintf(&b, "```text\n%s\n```\n\n", strings.TrimSpace(c.Long))
		} else {
			fmt.Fprintf(&b, "%s\n\n", c.Long)
		}
	} else if c.Short != "" {
		fmt.Fprintf(&b, "%s\n\n", c.Short)
	}
	if c.Runnable() {
		fmt.Fprintf(&b, "## Usage\n\n```\n%s\n```\n\n", c.UseLine())
	}
	if f := c.LocalFlags().FlagUsages(); f != "" {
		fmt.Fprintf(&b, "## Flags\n\n```\n%s```\n\n", f)
	}
	if f := c.InheritedFlags().FlagUsages(); f != "" {
		fmt.Fprintf(&b, "## Inherited flags\n\n```\n%s```\n\n", f)
	}
	var subs []*cobra.Command
	for _, s := range c.Commands() {
		if !s.Hidden && s.Name() != "help" {
			subs = append(subs, s)
		}
	}
	if len(subs) > 0 {
		b.WriteString("## Commands\n\n| Command | What it does |\n| - | - |\n")
		for _, s := range subs {
			fmt.Fprintf(&b, "| [`%s`](../%s/) | %s |\n", s.CommandPath(), commandSlug(s), s.Short)
		}
		b.WriteString("\n")
	}
	b.WriteString("_Generated from the elmos command definitions by `elmos gendocs`._\n")
	if err := os.WriteFile(filepath.Join(dir, commandSlug(c)+".md"), []byte(b.String()), 0o644); err != nil {
		return err
	}
	for _, s := range subs {
		if err := writeCommandDocs(s, dir); err != nil {
			return err
		}
	}
	return nil
}
