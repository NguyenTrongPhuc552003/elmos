package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildVersion(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{}
	appCtx := context.New(cfg, exec, fs)
	printer := ui.NewPrinter()

	ctx := &Context{
		AppContext: appCtx,
		Config:     cfg,
		Exec:       exec,
		Printer:    printer,
	}

	cmd := BuildVersion(ctx)

	if cmd == nil {
		t.Fatal("BuildVersion() returned nil")
	}
	if cmd.Use != "version" {
		t.Errorf("BuildVersion().Use = %v, want version", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildVersion() should have Short description")
	}

	// Should not error when executed
	if err := cmd.Execute(); err != nil {
		t.Errorf("BuildVersion().Execute() error = %v", err)
	}
}
