package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSearchOrderDoesNotCreateConfig(t *testing.T) {
	root := t.TempDir()
	changeDir(t, root)
	home := filepath.Join(root, "home")
	t.Setenv("HOME", home)

	// With no config, loading uses defaults and leaves the workspace untouched.
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ConfigFile != "" || cfg.Build.Arch != DefaultArch {
		t.Fatalf("no-config load = (%q, %q)", cfg.ConfigFile, cfg.Build.Arch)
	}
	assertMissing(t, filepath.Join(root, "build"))
	assertMissing(t, filepath.Join(root, "elmos.yaml"))

	homeConfig := filepath.Join(home, ".config", "elmos", "elmos.yaml")
	writeConfig(t, homeConfig, "riscv")
	cfg, err = Load("")
	if err != nil || cfg.Build.Arch != "riscv" {
		t.Fatalf("home config: arch=%q, err=%v", cfg.Build.Arch, err)
	}
	assertMissing(t, filepath.Join(root, "build"))

	writeConfig(t, filepath.Join(root, "build", "elmos.yaml"), "arm")
	cfg, err = Load("")
	if err != nil || cfg.Build.Arch != "arm" {
		t.Fatalf("build config: arch=%q, err=%v", cfg.Build.Arch, err)
	}

	writeConfig(t, filepath.Join(root, "elmos.yaml"), "arm64")
	cfg, err = Load("")
	if err != nil || cfg.Build.Arch != "arm64" {
		t.Fatalf("cwd config: arch=%q, err=%v", cfg.Build.Arch, err)
	}
}

func changeDir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func writeConfig(t *testing.T, path, arch string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("build:\n  arch: "+arch+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be absent, stat error: %v", path, err)
	}
}
