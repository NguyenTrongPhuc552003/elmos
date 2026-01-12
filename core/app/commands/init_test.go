package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildInit(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Image: config.ImageConfig{VolumeName: "test", MountPoint: "/Volumes/test"},
		Paths: config.PathsConfig{ProjectRoot: t.TempDir()},
	}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	cmd := BuildInit(ctx)

	if cmd == nil {
		t.Fatal("BuildInit() returned nil")
	}
	if cmd.Use != "init [workspace_name] [size]" {
		t.Errorf("BuildInit().Use = %v", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildInit() should have Short description")
	}
}

func TestBuildExit(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Image: config.ImageConfig{MountPoint: "/Volumes/test"},
	}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	cmd := BuildExit(ctx)

	if cmd == nil {
		t.Fatal("BuildExit() returned nil")
	}
	if cmd.Use != "exit" {
		t.Errorf("BuildExit().Use = %v, want exit", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildExit() should have Short description")
	}

	// Should have force flag
	forceFlag := cmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("BuildExit() should have --force flag")
	}
}
