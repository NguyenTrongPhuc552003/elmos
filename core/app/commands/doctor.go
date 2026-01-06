package commands

import (
	"github.com/spf13/cobra"
)

// BuildDoctor creates the doctor command for environment checking.
func BuildDoctor(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Printer.Info("ELMOS Doctor - Environment Check")
			ctx.Printer.Print("")
			results, issues := ctx.HealthChecker.CheckAll(cmd.Context())
			currentSection := ""
			for _, r := range results {
				section := getSection(r.Name)
				if section != currentSection {
					ctx.Printer.Step("Checking %s...", section)
					currentSection = section
				}
				if r.Passed {
					ctx.Printer.Print("  ✓ %s", r.Name)
				} else if r.Required {
					ctx.Printer.Print("  ✗ %s (missing)", r.Name)
				} else {
					ctx.Printer.Print("  ○ %s - optional", r.Name)
				}
			}
			if ctx.AutoFixer.CanFixElfH() {
				ctx.Printer.Print("")
				ctx.Printer.Step("Downloading missing elf.h...")
				if err := ctx.AutoFixer.FixElfH(); err != nil {
					ctx.Printer.Error("Failed to download elf.h: %v", err)
				} else {
					ctx.Printer.Success("elf.h downloaded")
					issues--
				}
			}
			ctx.Printer.Print("")
			if issues == 0 {
				ctx.Printer.Success("All checks passed!")
			} else {
				ctx.Printer.Warn("Found %d issue(s)", issues)
			}
			return nil
		},
	}
}

// getSection extracts section name from a check name.
func getSection(name string) string {
	// Match exact names first
	if name == "Homebrew" {
		return "Package Manager"
	}

	// Match prefixes
	prefixSections := []struct {
		prefix  string
		section string
	}{
		{"Tap:", "Homebrew Taps"},
		{"Package:", "Homebrew Packages"},
		{"Header:", "Custom Headers"},
		{"GDB:", "Cross Debuggers"},
		{"GCC:", "Cross Compilers"},
	}

	for _, ps := range prefixSections {
		if len(name) >= len(ps.prefix) && name[:len(ps.prefix)] == ps.prefix {
			return ps.section
		}
	}

	return "Other"
}
