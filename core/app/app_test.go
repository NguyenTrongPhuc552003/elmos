// Package app provides the CLI application layer for elmos.
package app

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNew(t *testing.T) {
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	cfg := &config.Config{}

	app := New(exec, fs, cfg)

	if app == nil {
		t.Fatal("New() returned nil")
	}
	// App should have initialized fields - test BuildRootCommand works
	cmd := app.BuildRootCommand()
	if cmd == nil {
		t.Error("New() app.BuildRootCommand() returned nil")
	}
}

func TestApp_BuildRootCommand(t *testing.T) {
	exec := executor.NewMockExecutor()
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
		t.Errorf("BuildRootCommand().Use = %v, want elmos", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildRootCommand() should have Short description")
	}

	// Should have subcommands
	if len(cmd.Commands()) == 0 {
		t.Error("BuildRootCommand() should have subcommands")
	}
}
