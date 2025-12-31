// Package cmd implements the Cobra CLI commands for elmos.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd generates shell completions
var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for elmos.

Supported shells: bash, zsh, fish, powershell

Examples:

  # Bash (add to ~/.bashrc):
  source <(elmos completion bash)

  # Zsh (add to ~/.zshrc):
  source <(elmos completion zsh)

  # Or install to completions directory:
  elmos completion zsh > "${fpath[1]}/_elmos"

  # Fish:
  elmos completion fish > ~/.config/fish/completions/elmos.fish

  # PowerShell:
  elmos completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
