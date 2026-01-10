// Package config provides configuration management for elmos.
// This file contains configuration saving and profile management.
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "saved_config.yaml")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		cfg     *Config
		args    args
		wantErr bool
	}{
		{
			name: "Save valid config",
			cfg: &Config{
				Build: BuildConfig{Arch: "x86_64"},
				Paths: PathsConfig{ProjectRoot: tmpDir},
			},
			args:    args{path: configPath},
			wantErr: false,
		},
		{
			name:    "Save to invalid path (permission error)",
			cfg:     &Config{},
			args:    args{path: "/root/invalid/path.yaml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Save(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Config.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := os.Stat(tt.args.path); os.IsNotExist(err) {
					t.Error("Config.Save() did not create file")
				}
			}
		})
	}
}

func TestConfig_computeDefaults(t *testing.T) {
	cfg := &Config{Paths: PathsConfig{ProjectRoot: "/test"}}
	got := cfg.computeDefaults()

	if got.Paths.ProjectRoot != "/test" {
		t.Errorf("computeDefaults() ProjectRoot = %v, want /test", got.Paths.ProjectRoot)
	}
	if got.Image.Path == "" {
		t.Error("computeDefaults() did not set defaults")
	}
}

func TestConfig_prepareForSave(t *testing.T) {
	defaults := &Config{
		Paths: PathsConfig{KernelDir: "/default/kernel"},
		Image: ImageConfig{Path: "/default/image"},
	}
	cfg := &Config{
		Paths: PathsConfig{KernelDir: "/default/kernel", ModulesDir: "/custom/modules"},
		Image: ImageConfig{Path: "/custom/image"},
	}
	got := cfg.prepareForSave(defaults)

	if got.Paths.KernelDir != "" {
		t.Errorf("prepareForSave() KernelDir = %v, want empty (matched default)", got.Paths.KernelDir)
	}
	if got.Paths.ModulesDir != "/custom/modules" {
		t.Errorf("prepareForSave() ModulesDir = %v, want /custom/modules", got.Paths.ModulesDir)
	}
	if got.Image.Path != "/custom/image" {
		t.Errorf("prepareForSave() ImagePath = %v, want /custom/image", got.Image.Path)
	}
}

func Test_clearDefaultPaths(t *testing.T) {
	defaults := PathsConfig{KernelDir: "/def", AppsDir: "/def"}
	paths := PathsConfig{KernelDir: "/def", AppsDir: "/custom"}

	got := clearDefaultPaths(paths, defaults)
	if got.KernelDir != "" {
		t.Error("clearDefaultPaths failed to clear matching kernel dir")
	}
	if got.AppsDir != "/custom" {
		t.Error("clearDefaultPaths cleared non-matching apps dir")
	}
}

func Test_clearDefaultImage(t *testing.T) {
	defaults := ImageConfig{Path: "/def/img"}
	image := ImageConfig{Path: "/def/img"}

	got := clearDefaultImage(image, defaults)
	if got.Path != "" {
		t.Error("clearDefaultImage failed to clear matching path")
	}
}

func Test_ensureDir(t *testing.T) {
	tmpDir := t.TempDir()
	validDir := filepath.Join(tmpDir, "newdir")

	if err := ensureDir(validDir); err != nil {
		t.Errorf("ensureDir() error = %v", err)
	}

	if _, err := os.Stat(validDir); os.IsNotExist(err) {
		t.Error("ensureDir() failed to create directory")
	}

	// Test invalid path (if running as non-root, this usually fails on system dirs, but harder to reliable test cross-platform without mocking FS.
	// We'll skip forcing an error here to match common practice unless mocking syscalls).
}

func TestConfig_ApplyProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]ProfileConfig{
			"debug": {Arch: "x86_64", Jobs: 2},
		},
		Build: BuildConfig{Arch: "arm64", Jobs: 4},
	}

	// Success case
	if err := cfg.ApplyProfile("debug"); err != nil {
		t.Errorf("ApplyProfile() error = %v", err)
	}
	if cfg.Build.Arch != "x86_64" {
		t.Errorf("ApplyProfile() Arch = %v, want x86_64", cfg.Build.Arch)
	}
	if cfg.Build.Jobs != 2 {
		t.Errorf("ApplyProfile() Jobs = %v, want 2", cfg.Build.Jobs)
	}

	// Failure case
	if err := cfg.ApplyProfile("unknown"); err == nil {
		t.Error("ApplyProfile() expected error for unknown profile")
	}
}

func TestConfig_GetArchConfig(t *testing.T) {
	cfg := &Config{Build: BuildConfig{Arch: "arm64"}}
	got := cfg.GetArchConfig()
	if got.Name != "arm64" {
		t.Errorf("GetArchConfig() Name = %v, want arm64", got.Name)
	}
}
