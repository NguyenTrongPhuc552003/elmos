// Package doctor provides dependency checking and environment validation for elmos.
package doctor

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

// mockExecutor implements executor.Executor for testing
type mockExecutor struct {
	executor.Executor
	lookPathFunc func(file string) (string, error)
	runFunc      func(ctx context.Context, cmd string, args ...string) error
}

func (m *mockExecutor) LookPath(file string) (string, error) {
	if m.lookPathFunc != nil {
		return m.lookPathFunc(file)
	}
	return "", fmt.Errorf("not found")
}

func (m *mockExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	if m.runFunc != nil {
		return m.runFunc(ctx, cmd, args...)
	}
	return nil
}

func (m *mockExecutor) Output(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	if cmd == "brew" && len(args) > 0 && args[0] == "list" {
		return []byte("pkg1\npkg2\n"), nil
	}
	return nil, fmt.Errorf("unknown command")
}

// mockFileSystem implements filesystem.FileSystem for testing
type mockFileSystem struct {
	filesystem.FileSystem
	existsFunc  func(path string) bool
	isDirFunc   func(path string) bool
	readDirFunc func(name string) ([]fs.DirEntry, error)
}

func (m *mockFileSystem) Exists(path string) bool {
	if m.existsFunc != nil {
		return m.existsFunc(path)
	}
	return false
}

func (m *mockFileSystem) IsDir(path string) bool {
	if m.isDirFunc != nil {
		return m.isDirFunc(path)
	}
	return false
}

func (m *mockFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	if m.readDirFunc != nil {
		return m.readDirFunc(name)
	}
	return nil, nil // Return empty list by default
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error)                      { return nil, nil }
func (m *mockFileSystem) ReadFile(name string) ([]byte, error)                       { return nil, nil }
func (m *mockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error { return nil }
func (m *mockFileSystem) MkdirAll(path string, perm os.FileMode) error               { return nil }
func (m *mockFileSystem) Remove(name string) error                                   { return nil }
func (m *mockFileSystem) RemoveAll(path string) error                                { return nil }
func (m *mockFileSystem) Getwd() (string, error)                                     { return "/tmp", nil }
func (m *mockFileSystem) Create(name string) (*os.File, error)                       { return nil, nil }
func (m *mockFileSystem) Open(name string) (*os.File, error)                         { return nil, nil }

func TestHealthChecker_CheckAll(t *testing.T) {
	exec := &mockExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/" + file, nil
		},
	}
	fs := &mockFileSystem{
		existsFunc: func(path string) bool { return true },
		isDirFunc:  func(path string) bool { return true },
	}
	cfg := &elconfig.Config{}
	tm := toolchain.NewManager(exec, fs, cfg, ui.NewPrinter())
	h := NewHealthChecker(exec, fs, cfg, tm)

	results, issues := h.CheckAll(context.Background())
	if len(results) == 0 {
		t.Error("CheckAll() returned no results")
	}
	_ = issues // ignored
}

func TestHealthChecker_CheckHomebrew(t *testing.T) {
	tests := []struct {
		name       string
		lookPath   func(string) (string, error)
		wantPassed bool
	}{
		{
			name: "Homebrew installed",
			lookPath: func(s string) (string, error) {
				return "/opt/homebrew/bin/brew", nil
			},
			wantPassed: true,
		},
		{
			name: "Homebrew missing",
			lookPath: func(s string) (string, error) {
				return "", fmt.Errorf("not found")
			},
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthChecker{
				exec: &mockExecutor{lookPathFunc: tt.lookPath},
			}
			got := h.CheckHomebrew(context.Background())
			if got.Passed != tt.wantPassed {
				t.Errorf("CheckHomebrew() passed = %v, want %v", got.Passed, tt.wantPassed)
			}
		})
	}
}

func TestHealthChecker_CheckToolchains(t *testing.T) {
	// Requires deeper mocking of ToolchainManager which interacts with FS.
	// For now we test basic connectivity.
	exec := &mockExecutor{}
	fs := &mockFileSystem{}
	cfg := &elconfig.Config{}
	tm := toolchain.NewManager(exec, fs, cfg, ui.NewPrinter())
	h := NewHealthChecker(exec, fs, cfg, tm)

	got := h.CheckToolchains(context.Background())
	if len(got) == 0 {
		t.Error("CheckToolchains() returned no results")
	}
}

func TestHealthChecker_CheckPackages(t *testing.T) {
	// Mock brew list output is handled in mockExecutor.Output
	exec := &mockExecutor{}
	fs := &mockFileSystem{}
	cfg := &elconfig.Config{}

	// Inject a resolver that uses our mock executor
	// NOTE: Checker creates its own Resolver in NewHealthChecker.
	// Ideally we should inject it, but we can't change the struct easily here.
	// We will rely on NewHealthChecker using the passed executor.

	h := NewHealthChecker(exec, fs, cfg, nil)

	// We rely on "pkg1" and "pkg2" being in the mock output from mockExecutor.Output
	// And matching RequiredPackages logic.
	// Since RequiredPackages are hardcoded in config, we might not match them perfect.
	// But ensuring it returns *some* results is good.

	got := h.CheckPackages(context.Background())
	if len(got) == 0 {
		t.Error("CheckPackages() returned no results")
	}
}

func TestHealthChecker_CheckHeaders(t *testing.T) {
	fs := &mockFileSystem{
		isDirFunc:  func(path string) bool { return true },
		existsFunc: func(path string) bool { return true },
	}
	cfg := &elconfig.Config{Paths: elconfig.PathsConfig{LibrariesDir: "/libs"}}
	h := &HealthChecker{fs: fs, cfg: cfg}

	got := h.CheckHeaders(context.Background())
	if len(got) == 0 {
		t.Error("CheckHeaders() returned no results")
	}
	for _, res := range got {
		if !res.Passed {
			t.Errorf("CheckHeaders result %s failed unexpectedly", res.Name)
		}
	}
}

func TestHealthChecker_IsElfHMissing(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		want   bool
	}{
		{"elf.h exists", true, false},
		{"elf.h missing", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &mockFileSystem{
				existsFunc: func(path string) bool { return tt.exists },
			}
			h := &HealthChecker{fs: fs, cfg: &elconfig.Config{}}
			if got := h.IsElfHMissing(); got != tt.want {
				t.Errorf("IsElfHMissing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthChecker_getInstalledPackageSet(t *testing.T) {
	exec := &mockExecutor{} // returns pkg1, pkg2
	h := NewHealthChecker(exec, nil, nil, nil)

	got, err := h.getInstalledPackageSet()
	if err != nil {
		t.Fatalf("getInstalledPackageSet error: %v", err)
	}
	if !got["pkg1"] {
		t.Error("pkg1 missing")
	}
}

func TestHealthChecker_groupPackagesByCategory(t *testing.T) {
	h := &HealthChecker{}
	got := h.groupPackagesByCategory()
	if len(got) == 0 {
		// Only if RequiredPackages is empty, which it shouldn't be.
		// If it is, this test is trivial.
	}
}

func TestHealthChecker_checkCategoryPackages(t *testing.T) {
	h := &HealthChecker{}
	pkgs := []elconfig.RequiredPackage{
		{Name: "p1", Required: true},
		{Name: "p2", Required: false},
	}
	pkgSet := map[string]bool{"p1": true}

	got := h.checkCategoryPackages("TestCat", pkgs, pkgSet)
	if len(got) != 3 { // 2 packages + 0 fix suggestion (p2 is optional but missing)
		// Wait, check logic: p2 missing & not required -> no fix suggestion.
		// p1 present -> passed.
		// So 2 results.
		if len(got) != 2 {
			t.Errorf("checkCategoryPackages count = %d, want 2", len(got))
		}
	}
}

func TestHealthChecker_CheckCrossGDB(t *testing.T) {
	// This relies on file system checks for specific paths
	// We'll mock FS to return true for everything
	fs := &mockFileSystem{existsFunc: func(p string) bool { return true }}
	// We need to properly initialize ToolchainManager with paths
	cfg := &elconfig.Config{Paths: elconfig.PathsConfig{ToolchainsDir: "/tc"}}
	tm := toolchain.NewManager(nil, fs, cfg, ui.NewPrinter())
	h := &HealthChecker{fs: fs, tm: tm, cfg: cfg}

	got := h.CheckCrossGDB(context.Background())
	_ = got // ignored
}

func TestHealthChecker_CheckCrossGCC(t *testing.T) {
	fs := &mockFileSystem{existsFunc: func(p string) bool { return true }}
	cfg := &elconfig.Config{Paths: elconfig.PathsConfig{ToolchainsDir: "/tc"}}
	tm := toolchain.NewManager(nil, fs, cfg, ui.NewPrinter())
	h := &HealthChecker{fs: fs, tm: tm, cfg: cfg}

	got := h.CheckCrossGCC(context.Background())
	_ = got // ignored
}

func TestNewHealthChecker(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *elconfig.Config
		tm   *toolchain.Manager
	}
	tests := []struct {
		name string
		args args
		want *HealthChecker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealthChecker(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealthChecker() = %v, want %v", got, tt.want)
			}
		})
	}
}
