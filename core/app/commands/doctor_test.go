package commands

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestBuildDoctor(t *testing.T) {
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

	cmd := BuildDoctor(ctx)

	if cmd == nil {
		t.Fatal("BuildDoctor() returned nil")
	}
	if cmd.Use != "doctor" {
		t.Errorf("BuildDoctor().Use = %v, want doctor", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("BuildDoctor() should have Short description")
	}
}

func Test_getSection(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "Environment", arg: "Environment"},
		{name: "Toolchain", arg: "Toolchain"},
		{name: "Custom", arg: "MySection"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSection(tt.arg)
			// getSection should return a non-empty formatted string
			if got == "" {
				t.Error("getSection() returned empty string")
			}
		})
	}
}
