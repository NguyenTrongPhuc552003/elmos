package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildStatus(t *testing.T) {
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

	cmd := BuildStatus(ctx)

	if cmd == nil {
		t.Fatal("BuildStatus() returned nil")
	}
	if cmd.Use != "status" {
		t.Errorf("BuildStatus().Use = %v, want status", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildStatus() should have Short description")
	}
}
