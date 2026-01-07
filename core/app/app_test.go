package app

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNew(t *testing.T) {
	exec := executor.NewShellExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Paths: config.PathsConfig{
			ProjectRoot: "/test",
		},
	}

	app := New(exec, fs, cfg)

	if app == nil {
		t.Fatal("New() returned nil")
	}
	if app.Exec == nil {
		t.Error("App.Exec should not be nil")
	}
	if app.FS == nil {
		t.Error("App.FS should not be nil")
	}
	if app.Config != cfg {
		t.Error("App.Config not set correctly")
	}
	if app.Context == nil {
		t.Error("App.Context should not be nil")
	}
	if app.KernelBuilder == nil {
		t.Error("App.KernelBuilder should not be nil")
	}
	if app.Printer == nil {
		t.Error("App.Printer should not be nil")
	}
}

func TestBuildRootCommand(t *testing.T) {
	exec := executor.NewShellExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
	}

	app := New(exec, fs, cfg)
	cmd := app.BuildRootCommand()

	if cmd == nil {
		t.Fatal("BuildRootCommand() returned nil")
	}
	if cmd.Use != "elmos" {
		t.Errorf("RootCommand.Use = %q, want %q", cmd.Use, "elmos")
	}
	if cmd.Version == "" {
		t.Error("RootCommand.Version should not be empty")
	}
}
