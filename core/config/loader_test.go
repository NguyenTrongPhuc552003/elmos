// Package config provides configuration management for elmos.
// This file contains configuration loading and default handling.
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file for testing
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.yaml")
	configContent := []byte(`
build:
  arch: "arm64"
`)
	if err := os.WriteFile(configFile, configContent, 0644); err != nil {
		t.Fatal(err)
	}

	type args struct {
		configPath string
	}
	tests := []struct {
		name     string
		args     args
		wantArch string
		wantErr  bool
	}{
		{
			name:     "Load default config (empty path)",
			args:     args{configPath: ""},
			wantArch: DefaultArch, // Should fallback to default
			wantErr:  false,
		},
		{
			name:     "Load specific config file",
			args:     args{configPath: configFile},
			wantArch: "arm64",
			wantErr:  false,
		},
		{
			name:     "Load non-existent config file",
			args:     args{configPath: "/non/existent/path.yaml"},
			wantArch: "", // Won't match valid config
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset() // Reset singleton state
			got, err := Load(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.Build.Arch != tt.wantArch {
				t.Errorf("Load() Arch = %v, want %v", got.Build.Arch, tt.wantArch)
			}
			if got.Paths.ProjectRoot == "" {
				t.Error("Load() ProjectRoot should not be empty")
			}
		})
	}
}

func TestGet(t *testing.T) {
	Reset()
	cfg1 := Get()
	if cfg1 == nil {
		t.Error("Get() returned nil")
	}
	cfg2 := Get()
	if cfg1 != cfg2 {
		t.Error("Get() did not return singleton instance")
	}
}

func TestReset(t *testing.T) {
	Get() // Ensure initialized
	if configInstance == nil {
		t.Fatal("Setup failed: configInstance is nil")
	}
	Reset()
	if configInstance != nil {
		t.Error("Reset() did not clear configInstance")
	}
}

func Test_setDefaults(t *testing.T) {
	v := viper.New()
	setDefaults(v)

	tests := []struct {
		key      string
		expected any
	}{
		{"image.volume_name", DefaultVolumeName},
		{"image.size", DefaultImageSize},
		{"build.arch", DefaultArch},
		{"build.jobs", runtime.NumCPU()},
		{"build.llvm", true},
		{"qemu.memory", DefaultMemory},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := v.Get(tt.key)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("setDefaults() %s = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func Test_applyComputedDefaults(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{},
		Image: ImageConfig{},
	}
	applyComputedDefaults(cfg)

	if cfg.Paths.ProjectRoot == "" {
		t.Error("ProjectRoot not set")
	}
	if cfg.Image.Path == "" {
		t.Error("Image Path not set")
	}
	if cfg.Paths.KernelDir == "" {
		t.Error("KernelDir not set")
	}
}

func Test_applyProjectRoot(t *testing.T) {
	cfg := &Config{Paths: PathsConfig{}}
	applyProjectRoot(cfg)
	cwd, _ := os.Getwd()
	if cfg.Paths.ProjectRoot != cwd {
		t.Errorf("applyProjectRoot() = %v, want %v", cfg.Paths.ProjectRoot, cwd)
	}

	// Should not overwrite if set
	cfg.Paths.ProjectRoot = "/custom/root"
	applyProjectRoot(cfg)
	if cfg.Paths.ProjectRoot != "/custom/root" {
		t.Error("applyProjectRoot() overwrite existing value")
	}
}

func Test_applyImageDefaults(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{ProjectRoot: "/root"},
		Image: ImageConfig{VolumeName: "TESTVOL"},
	}
	applyImageDefaults(cfg)

	wantPath := filepath.Join("/root", "data", "img.sparseimage")
	if cfg.Image.Path != wantPath {
		t.Errorf("applyImageDefaults() Path = %v, want %v", cfg.Image.Path, wantPath)
	}

	wantMount := filepath.Join("/Volumes", "TESTVOL")
	if cfg.Image.MountPoint != wantMount {
		t.Errorf("applyImageDefaults() MountPoint = %v, want %v", cfg.Image.MountPoint, wantMount)
	}
}

func Test_applyPathDefaults(t *testing.T) {
	cfg := &Config{
		Paths: PathsConfig{ProjectRoot: "/root"},
		Image: ImageConfig{MountPoint: "/mnt"},
	}
	applyPathDefaults(cfg)

	tests := []struct {
		got  string
		want string
	}{
		{cfg.Paths.KernelDir, filepath.Join("/mnt", "linux")},
		{cfg.Paths.ModulesDir, filepath.Join("/root", "modules")},
		{cfg.Paths.AppsDir, filepath.Join("/root", "apps")},
		{cfg.Paths.LibrariesDir, filepath.Join("/root", "libraries")},
		{cfg.Paths.RootfsDir, filepath.Join("/mnt", "rootfs")},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("path default mismatch: got %v, want %v", tt.got, tt.want)
		}
	}
}

func Test_setIfEmpty(t *testing.T) {
	val := ""
	setIfEmpty(&val, "default")
	if val != "default" {
		t.Error("setIfEmpty did not set empty value")
	}

	val = "existing"
	setIfEmpty(&val, "default")
	if val != "existing" {
		t.Error("setIfEmpty overwrote existing value")
	}
}
