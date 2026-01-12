package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildPatch(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Paths: config.PathsConfig{KernelDir: "/mnt/linux"},
	}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	cmd := BuildPatch(ctx)

	if cmd == nil {
		t.Fatal("BuildPatch() returned nil")
	}
	if cmd.Use != "patch" {
		t.Errorf("BuildPatch().Use = %v, want patch", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildPatch() should have Short description")
	}

	// Should have subcommands (apply, create, list, etc.)
	if len(cmd.Commands()) == 0 {
		t.Error("BuildPatch() should have subcommands")
	}
}
