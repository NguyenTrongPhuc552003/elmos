// Package commands provides individual CLI command builders for elmos.
// Each command group is in its own file for maintainability.
package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
	"github.com/spf13/cobra"
)

func TestRegister(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
	}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	rootCmd := &cobra.Command{Use: "elmos"}

	// Should not panic
	Register(ctx, rootCmd)

	// Should have subcommands after registration
	if len(rootCmd.Commands()) == 0 {
		t.Error("Register() should add subcommands to rootCmd")
	}
}
