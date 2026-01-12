package commands

import (
	"context"
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

func Test_validateVolumeSize(t *testing.T) {
	type args struct {
		size    string
		printer *ui.Printer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVolumeSize(tt.args.size, tt.args.printer); (err != nil) != tt.wantErr {
				t.Errorf("validateVolumeSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findDiskDevice(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findDiskDevice(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findDiskDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("findDiskDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDiskDeviceFromHdiutil(t *testing.T) {
	type args struct {
		output    string
		imagePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDiskDeviceFromHdiutil(tt.args.output, tt.args.imagePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDiskDeviceFromHdiutil() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("parseDiskDeviceFromHdiutil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractDiskDevice(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractDiskDevice(tt.args.line); got != tt.want {
				t.Errorf("extractDiskDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}
