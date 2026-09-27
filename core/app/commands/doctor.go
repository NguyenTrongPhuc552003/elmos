package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// BuildDoctor creates the doctor command for environment checking.
func BuildDoctor(ctx *Context) *cobra.Command {
	var fix bool
	doctorCmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Printer.Info("ELMOS Doctor - Environment Check")
			ctx.Printer.Print("")
			var fixErr error
			if fix && ctx.AutoFixer.CanFixElfH() {
				ctx.Printer.Step("Downloading missing elf.h...")
				fixErr = ctx.AutoFixer.FixElfH()
				if fixErr == nil {
					ctx.Printer.Success("elf.h downloaded")
				}
				ctx.Printer.Print("")
			}
			results, issues := ctx.HealthChecker.CheckAll(cmd.Context())
			currentSection := ""
			for _, r := range results {
				section := getSection(r.Name)
				if section != currentSection {
					ctx.Printer.Step("Checking %s...", section)
					currentSection = section
				}

				// Strip section prefix if present
				displayName := r.Name
				prefix := section + ": "
				if len(displayName) > len(prefix) && displayName[:len(prefix)] == prefix {
					displayName = displayName[len(prefix):]
				}

				if r.Passed {
					ctx.Printer.Print("  ✓ %s", displayName)
				} else if r.Required {
					ctx.Printer.Print("  ✗ %s (missing)", displayName)
				} else {
					ctx.Printer.Print("  ○ %s - optional", displayName)
				}
			}
			if !fix && ctx.AutoFixer.CanFixElfH() {
				ctx.Printer.Print("")
				ctx.Printer.Info("Run 'elmos doctor --fix' to download missing elf.h")
			}
			ctx.Printer.Print("")
			if issues == 0 {
				ctx.Printer.Success("All checks passed!")
			} else {
				ctx.Printer.Warn("Found %d issue(s)", issues)
			}
			if fixErr != nil {
				return fmt.Errorf("failed to download elf.h: %w", fixErr)
			}
			return nil
		},
	}
	doctorCmd.Flags().BoolVar(&fix, "fix", false, "download missing elf.h")
	return doctorCmd
}

// getSection extracts section name from a check name.
func getSection(name string) string {
	// Match exact names first
	if name == "Homebrew" {
		return "Package Manager"
	}
	if name == "crosstool-ng" || strings.HasPrefix(name, "Toolchain:") {
		return "Toolchains"
	}

	// Dynamic sections based on prefix "Scetion: Item"
	if idx := strings.Index(name, ": "); idx != -1 {
		return name[:idx]
	}

	// fallback
	return "Other"
}
