package config

import (
	"testing"
)

func TestGetArchConfig(t *testing.T) {
	tests := []struct {
		name     string
		arch     string
		wantNil  bool
		wantName string
	}{
		{"arm64 exists", "arm64", false, "arm64"},
		{"arm exists", "arm", false, "arm"},
		{"riscv exists", "riscv", false, "riscv"},
		{"invalid arch", "x86", true, ""},
		{"empty arch", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetArchConfig(tt.arch)
			if tt.wantNil && got != nil {
				t.Errorf("GetArchConfig(%q) = %v, want nil", tt.arch, got)
			}
			if !tt.wantNil && got == nil {
				t.Errorf("GetArchConfig(%q) = nil, want non-nil", tt.arch)
			}
			if !tt.wantNil && got != nil && got.Name != tt.wantName {
				t.Errorf("GetArchConfig(%q).Name = %q, want %q", tt.arch, got.Name, tt.wantName)
			}
		})
	}
}

func TestIsValidArch(t *testing.T) {
	tests := []struct {
		arch string
		want bool
	}{
		{"arm64", true},
		{"arm", true},
		{"riscv", true},
		{"x86", false},
		{"x86_64", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.arch, func(t *testing.T) {
			if got := IsValidArch(tt.arch); got != tt.want {
				t.Errorf("IsValidArch(%q) = %v, want %v", tt.arch, got, tt.want)
			}
		})
	}
}

func TestSupportedArchitectures(t *testing.T) {
	archs := SupportedArchitectures()

	if len(archs) != 3 {
		t.Errorf("SupportedArchitectures() returned %d archs, want 3", len(archs))
	}

	// Check that all expected archs are present
	expected := map[string]bool{"arm64": false, "arm": false, "riscv": false}
	for _, a := range archs {
		if _, ok := expected[a]; ok {
			expected[a] = true
		}
	}
	for arch, found := range expected {
		if !found {
			t.Errorf("SupportedArchitectures() missing %q", arch)
		}
	}
}

func TestApplyComputedDefaults(t *testing.T) {
	cfg := &Config{
		Image: ImageConfig{
			VolumeName: "test-volume",
		},
		Paths: PathsConfig{
			ProjectRoot: "/tmp/test-project",
		},
	}

	applyComputedDefaults(cfg)

	// Check image path (now in build/ subdirectory)
	if cfg.Image.Path != "/tmp/test-project/build/img.sparseimage" {
		t.Errorf("Image.Path = %q, want %q", cfg.Image.Path, "/tmp/test-project/build/img.sparseimage")
	}

	// Check mount point
	if cfg.Image.MountPoint != "/Volumes/test-volume" {
		t.Errorf("Image.MountPoint = %q, want %q", cfg.Image.MountPoint, "/Volumes/test-volume")
	}

	// Check kernel dir
	if cfg.Paths.KernelDir != "/Volumes/test-volume/linux" {
		t.Errorf("Paths.KernelDir = %q, want %q", cfg.Paths.KernelDir, "/Volumes/test-volume/linux")
	}

	// Check modules dir
	if cfg.Paths.ModulesDir != "/tmp/test-project/modules" {
		t.Errorf("Paths.ModulesDir = %q, want %q", cfg.Paths.ModulesDir, "/tmp/test-project/modules")
	}
}

func TestApplyComputedDefaultsPreservesExisting(t *testing.T) {
	cfg := &Config{
		Image: ImageConfig{
			VolumeName: "test-volume",
			Path:       "/custom/path.sparseimage",
			MountPoint: "/custom/mount",
		},
		Paths: PathsConfig{
			ProjectRoot: "/tmp/test-project",
			KernelDir:   "/custom/kernel",
		},
	}

	applyComputedDefaults(cfg)

	// Custom values should be preserved
	if cfg.Image.Path != "/custom/path.sparseimage" {
		t.Errorf("Image.Path was overwritten, got %q", cfg.Image.Path)
	}
	if cfg.Image.MountPoint != "/custom/mount" {
		t.Errorf("Image.MountPoint was overwritten, got %q", cfg.Image.MountPoint)
	}
	if cfg.Paths.KernelDir != "/custom/kernel" {
		t.Errorf("Paths.KernelDir was overwritten, got %q", cfg.Paths.KernelDir)
	}
}

func TestReset(t *testing.T) {
	// Set a config instance
	configInstance = &Config{Build: BuildConfig{Arch: "test"}}

	// Reset should clear it
	Reset()

	if configInstance != nil {
		t.Error("Reset() did not clear configInstance")
	}
}
