package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestInformationalCommandsDoNotCreateConfig(t *testing.T) {
	for _, args := range [][]string{{"version"}, {"help"}, {"arch"}} {
		t.Run(args[0], func(t *testing.T) {
			root := t.TempDir()
			changeDir(t, root)
			t.Setenv("HOME", filepath.Join(root, "home"))
			application, _ := newTestApp(t)
			cmd := application.BuildRootCommand()
			cmd.SetArgs(args)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			assertMissing(t, filepath.Join(root, "build"))
			assertMissing(t, filepath.Join(root, "elmos.yaml"))
		})
	}
}

func TestExplicitConfigFlagLoadsSelectedFile(t *testing.T) {
	root := t.TempDir()
	changeDir(t, root)
	t.Setenv("HOME", filepath.Join(root, "home"))
	if err := os.WriteFile(filepath.Join(root, "elmos.yaml"), []byte("build:\n  arch: arm\n"), 0644); err != nil {
		t.Fatal(err)
	}
	custom := filepath.Join(root, "custom.yaml")
	if err := os.WriteFile(custom, []byte("build:\n  arch: riscv\n"), 0644); err != nil {
		t.Fatal(err)
	}
	application, _ := newTestApp(t)
	cmd := application.BuildRootCommand()
	cmd.SetArgs([]string{"--config", custom, "arch", "show"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if application.Config.Build.Arch != "riscv" || application.Config.ConfigFile != custom {
		t.Fatalf("selected config = (%q, %q)", application.Config.ConfigFile, application.Config.Build.Arch)
	}
	assertMissing(t, filepath.Join(root, "build"))
}

func TestInitCreatesConfigOnlyOnExplicitCommand(t *testing.T) {
	for _, custom := range []bool{false, true} {
		name := "default path"
		if custom {
			name = "custom path"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			changeDir(t, root)
			t.Setenv("HOME", filepath.Join(root, "home"))
			application, mock := newTestApp(t)
			args := []string{"init", "renamed"}
			configPath := filepath.Join(root, "elmos.yaml")
			if custom {
				configPath = filepath.Join(root, "custom", "elmos.yaml")
				args = append([]string{"--config", configPath}, args...)
			}
			cmd := application.BuildRootCommand()
			cmd.SetArgs(args)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			cfg, err := config.Load(configPath)
			if err != nil || cfg.Image.VolumeName != "renamed" {
				t.Fatalf("initialized config: volume=%q, err=%v", cfg.Image.VolumeName, err)
			}
			assertMissing(t, filepath.Join(root, "build"))
			assertMissing(t, cfg.Image.Path)
			if len(mock.Calls) == 0 {
				t.Fatal("expected mocked disk operations")
			}
		})
	}
}

func newTestApp(t *testing.T) (*App, *executor.MockExecutor) {
	t.Helper()
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	mock := executor.NewMockExecutor()
	return New(mock, filesystem.NewOSFileSystem(), cfg), mock
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

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be absent, stat error: %v", path, err)
	}
}
