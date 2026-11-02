package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/NguyenTrongPhuc552003/elmos/core/config"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
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

func TestInitReopensConfiguredImage(t *testing.T) {
	root := t.TempDir()
	changeDir(t, root)
	t.Setenv("HOME", filepath.Join(root, "home"))
	image := filepath.Join(root, "elsewhere", "demo.sparseimage")
	if err := os.MkdirAll(filepath.Dir(image), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(image, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "elmos.yaml")
	yaml := "image:\n  path: " + image + "\n  volume_name: demo\n  size: 40G\n  mount_point: /Volumes/demo\n"
	if err := os.WriteFile(configPath, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	application, mock := newTestApp(t)
	cmd := application.BuildRootCommand()
	cmd.SetArgs([]string{"--config", configPath, "init", "demo"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if application.Config.Image.Path != image {
		t.Fatalf("image path = %q, want %q", application.Config.Image.Path, image)
	}
	for _, c := range mock.Calls {
		if len(c.Args) > 0 && c.Args[0] == "create" {
			t.Fatalf("re-init created a new image: %v", c.Args)
		}
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

// Every command shown in the root help's "Common workflow" must exist
// (audit E1: the help once showed "elmos build", which never existed).
func TestHelpWorkflowCommandsExist(t *testing.T) {
	root := t.TempDir()
	changeDir(t, root)
	t.Setenv("HOME", filepath.Join(root, "home"))
	application, _ := newTestApp(t)
	cmd := application.BuildRootCommand()
	seen := 0
	for _, line := range strings.Split(cmd.Long, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "elmos ") {
			continue
		}
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		args := strings.Fields(line)[1:]
		found, rest, err := cmd.Find(args)
		if err != nil || found == cmd || len(rest) != 0 {
			t.Errorf("help workflow shows %q, which is not a command", line)
		}
		seen++
	}
	if seen == 0 {
		t.Fatal("no workflow lines found in the root help")
	}
}
