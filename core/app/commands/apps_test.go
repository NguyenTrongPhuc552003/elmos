package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildApps(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Paths: config.PathsConfig{AppsDir: "/apps"},
	}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	cmd := BuildApps(ctx)

	if cmd == nil {
		t.Fatal("BuildApps() returned nil")
	}
	if cmd.Use != "app" {
		t.Errorf("BuildApps().Use = %v, want app", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildApps() should have Short description")
	}

	// Should have subcommands (build, list, etc.)
	if len(cmd.Commands()) == 0 {
		t.Error("BuildApps() should have subcommands")
	}
}
