package doctor

import (
	"io/fs"
	"os"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

func TestCheckResultString(t *testing.T) {
	result := CheckResult{
		Name:     "Homebrew",
		Passed:   true,
		Required: true,
		Message:  "Install from: https://brew.sh",
	}

	if result.Name != "Homebrew" {
		t.Errorf("CheckResult.Name = %q, want %q", result.Name, "Homebrew")
	}
	if !result.Passed {
		t.Error("CheckResult.Passed should be true")
	}
	if !result.Required {
		t.Error("CheckResult.Required should be true")
	}
	expectedMsg := "Install from: https://brew.sh"
	if result.Message != expectedMsg {
		t.Errorf("CheckResult.Message = %q, want %q", result.Message, expectedMsg)
	}
}

func TestNewHealthChecker(t *testing.T) {
	cfg := &elconfig.Config{
		Paths: elconfig.PathsConfig{
			LibrariesDir: "/test/libraries",
		},
	}

	hc := NewHealthChecker(nil, nil, cfg)

	if hc == nil {
		t.Fatal("NewHealthChecker returned nil")
	}
	if hc.cfg != cfg {
		t.Error("HealthChecker.cfg not set correctly")
	}
	if hc.brew == nil {
		t.Error("HealthChecker.brew should not be nil")
	}
}

func TestIsElfHMissing(t *testing.T) {
	// Create a mock filesystem
	cfg := &elconfig.Config{
		Paths: elconfig.PathsConfig{
			LibrariesDir: "/nonexistent/path",
		},
	}

	hc := &HealthChecker{
		cfg: cfg,
		fs:  &mockFS{existsMap: map[string]bool{}},
	}

	// elf.h doesn't exist, should return true
	if !hc.IsElfHMissing() {
		t.Error("IsElfHMissing() should return true when elf.h doesn't exist")
	}
}

func TestIsElfHExists(t *testing.T) {
	cfg := &elconfig.Config{
		Paths: elconfig.PathsConfig{
			LibrariesDir: "/test/libraries",
		},
	}

	hc := &HealthChecker{
		cfg: cfg,
		fs: &mockFS{existsMap: map[string]bool{
			"/test/libraries/elf.h": true,
		}},
	}

	// elf.h exists, should return false
	if hc.IsElfHMissing() {
		t.Error("IsElfHMissing() should return false when elf.h exists")
	}
}

// mockFS is a minimal mock for testing
type mockFS struct {
	existsMap map[string]bool
}

func (m *mockFS) Exists(path string) bool {
	return m.existsMap[path]
}

// Implement remaining FileSystem interface methods
func (m *mockFS) IsDir(path string) bool                                     { return false }
func (m *mockFS) Stat(name string) (os.FileInfo, error)                      { return nil, nil }
func (m *mockFS) ReadFile(name string) ([]byte, error)                       { return nil, nil }
func (m *mockFS) WriteFile(name string, data []byte, perm os.FileMode) error { return nil }
func (m *mockFS) MkdirAll(path string, perm os.FileMode) error               { return nil }
func (m *mockFS) ReadDir(name string) ([]fs.DirEntry, error)                 { return nil, nil }
func (m *mockFS) Remove(name string) error                                   { return nil }
func (m *mockFS) RemoveAll(path string) error                                { return nil }
func (m *mockFS) Getwd() (string, error)                                     { return "", nil }
func (m *mockFS) Create(name string) (*os.File, error)                       { return nil, nil }
func (m *mockFS) Open(name string) (*os.File, error)                         { return nil, nil }
