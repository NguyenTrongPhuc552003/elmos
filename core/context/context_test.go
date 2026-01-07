package context

import (
	"io/fs"
	"os"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// mockFileSystem implements filesystem.FileSystem for testing.
type mockFileSystem struct {
	existsMap map[string]bool
	isDirMap  map[string]bool
}

func (m *mockFileSystem) Exists(path string) bool {
	if m.existsMap == nil {
		return false
	}
	return m.existsMap[path]
}

func (m *mockFileSystem) IsDir(path string) bool {
	if m.isDirMap == nil {
		return false
	}
	return m.isDirMap[path]
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error)                      { return nil, nil }
func (m *mockFileSystem) ReadFile(name string) ([]byte, error)                       { return nil, nil }
func (m *mockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error { return nil }
func (m *mockFileSystem) MkdirAll(path string, perm os.FileMode) error               { return nil }
func (m *mockFileSystem) ReadDir(name string) ([]fs.DirEntry, error)                 { return nil, nil }
func (m *mockFileSystem) Remove(name string) error                                   { return nil }
func (m *mockFileSystem) RemoveAll(path string) error                                { return nil }
func (m *mockFileSystem) Getwd() (string, error)                                     { return "", nil }
func (m *mockFileSystem) Create(name string) (*os.File, error)                       { return nil, nil }
func (m *mockFileSystem) Open(name string) (*os.File, error)                         { return nil, nil }

func TestKernelExists(t *testing.T) {
	tests := []struct {
		name      string
		gitExists bool
		want      bool
	}{
		{"git dir exists", true, true},
		{"git dir missing", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Paths: config.PathsConfig{
					KernelDir: "/test/kernel",
				},
			}
			fs := &mockFileSystem{
				existsMap: map[string]bool{
					"/test/kernel/.git": tt.gitExists,
				},
			}
			ctx := &Context{Config: cfg, FS: fs}

			if got := ctx.KernelExists(); got != tt.want {
				t.Errorf("KernelExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasConfig(t *testing.T) {
	tests := []struct {
		name         string
		configExists bool
		want         bool
	}{
		{"config exists", true, true},
		{"config missing", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Paths: config.PathsConfig{
					KernelDir: "/test/kernel",
				},
			}
			fs := &mockFileSystem{
				existsMap: map[string]bool{
					"/test/kernel/.config": tt.configExists,
				},
			}
			ctx := &Context{Config: cfg, FS: fs}

			if got := ctx.HasConfig(); got != tt.want {
				t.Errorf("HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKernelImage(t *testing.T) {
	cfg := &config.Config{
		Build: config.BuildConfig{Arch: "arm64"},
		Paths: config.PathsConfig{KernelDir: "/test/kernel"},
	}
	ctx := &Context{Config: cfg}

	got := ctx.GetKernelImage()
	want := "/test/kernel/arch/arm64/boot/Image"

	if got != want {
		t.Errorf("GetKernelImage() = %q, want %q", got, want)
	}
}

func TestGetVmlinux(t *testing.T) {
	cfg := &config.Config{
		Paths: config.PathsConfig{KernelDir: "/test/kernel"},
	}
	ctx := &Context{Config: cfg}

	got := ctx.GetVmlinux()
	want := "/test/kernel/vmlinux"

	if got != want {
		t.Errorf("GetVmlinux() = %q, want %q", got, want)
	}
}

func TestGetDefaultTargets(t *testing.T) {
	tests := []struct {
		name    string
		arch    string
		wantLen int
	}{
		{"arm64", "arm64", 3},
		{"arm", "arm", 3},
		{"riscv", "riscv", 3},
		{"invalid", "invalid", 3}, // Falls back to default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Build: config.BuildConfig{Arch: tt.arch},
			}
			ctx := &Context{Config: cfg}

			got := ctx.GetDefaultTargets()
			if len(got) != tt.wantLen {
				t.Errorf("GetDefaultTargets() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}
