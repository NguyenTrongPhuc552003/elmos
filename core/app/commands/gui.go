package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// BuildGUI creates the gui command for launching native macOS GUI.
func BuildGUI(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Launch native macOS GUI application",
		Long:  "Launch the native SwiftUI-based graphical user interface for ELMOS.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find GUI binary (Swift executable)
			guiBinary := filepath.Join("build", "gui", "elmos")

			// Check if GUI binary exists
			if _, err := os.Stat(guiBinary); os.IsNotExist(err) {
				return fmt.Errorf("GUI binary not found at %s. Run 'task gui:build' first", guiBinary)
			}

			// Execute GUI app
			guiCmd := exec.Command(guiBinary)
			guiCmd.Stdout = os.Stdout
			guiCmd.Stderr = os.Stderr
			guiCmd.Env = os.Environ()

			return guiCmd.Run()
		},
	}
}
