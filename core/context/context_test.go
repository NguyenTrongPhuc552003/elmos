// Package context provides build context management for elmos.
package context

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/homebrew"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()

	ctx := New(cfg, exec, fs)

	if ctx == nil {
		t.Fatal("New() returned nil")
	}
	if ctx.Config != cfg {
		t.Error("New() Config not set correctly")
	}
	if ctx.Exec != exec {
		t.Error("New() Exec not set correctly")
	}
	if ctx.FS != fs {
		t.Error("New() FS not set correctly")
	}
	if ctx.Brew == nil {
		t.Error("New() Brew should be initialized")
	}
}

func TestContext_IsMounted(t *testing.T) {
	tmpDir := t.TempDir()
	mountPoint := filepath.Join(tmpDir, "mount")

	// Create mock executor that returns empty output
	exec := executor.NewMockExecutor()
	exec.OutputResponses["mount"] = []byte("")
	exec.OutputResponses["hdiutil"] = []byte("")

	fs := filesystem.NewOSFileSystem()

	tests := []struct {
		name       string
		mountPoint string
		setupDir   bool
		want       bool
	}{
		{
			name:       "Mount point does not exist",
			mountPoint: "/non/existent/path",
			setupDir:   false,
			want:       false,
		},
		{
			name:       "Mount point exists but not in mount output",
			mountPoint: mountPoint,
			setupDir:   true,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupDir {
				_ = fs.MkdirAll(tt.mountPoint, 0755)
			}

			ctx := &Context{
				Config: &config.Config{
					Image: config.ImageConfig{
						MountPoint: tt.mountPoint,
						Path:       "test.sparseimage",
					},
				},
				Exec: exec,
				FS:   fs,
			}
			if got := ctx.IsMounted(); got != tt.want {
				t.Errorf("Context.IsMounted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_EnsureMounted(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["mount"] = []byte("")
	exec.OutputResponses["hdiutil"] = []byte("")
	fs := filesystem.NewOSFileSystem()

	tests := []struct {
		name    string
		ctx     *Context
		wantErr bool
	}{
		{
			name: "Not mounted returns error",
			ctx: &Context{
				Config: &config.Config{
					Image: config.ImageConfig{
						MountPoint: "/non/existent",
						Path:       "test.sparseimage",
					},
				},
				Exec: exec,
				FS:   fs,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ctx.EnsureMounted(); (err != nil) != tt.wantErr {
				t.Errorf("Context.EnsureMounted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContext_GetActualMountPoint(t *testing.T) {
	tmpDir := t.TempDir()
	mountPoint := filepath.Join(tmpDir, "mount")

	exec := executor.NewMockExecutor()
	fs := filesystem.NewOSFileSystem()
	_ = fs.MkdirAll(mountPoint, 0755)

	tests := []struct {
		name    string
		ctx     *Context
		want    string
		wantErr bool
	}{
		{
			name: "Mount point exists returns fast path",
			ctx: &Context{
				Config: &config.Config{
					Image: config.ImageConfig{MountPoint: mountPoint},
				},
				FS:   fs,
				Exec: exec,
			},
			want:    mountPoint,
			wantErr: false,
		},
		{
			name: "Mount point does not exist returns error",
			ctx: &Context{
				Config: &config.Config{
					Image: config.ImageConfig{
						MountPoint: "/non/existent",
						Path:       "test.sparseimage",
					},
				},
				FS:   fs,
				Exec: exec,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ctx.GetActualMountPoint()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Context.GetActualMountPoint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Context.GetActualMountPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_KernelExists(t *testing.T) {
	tmpDir := t.TempDir()
	kernelDir := filepath.Join(tmpDir, "linux")
	gitDir := filepath.Join(kernelDir, ".git")

	fs := filesystem.NewOSFileSystem()
	_ = fs.MkdirAll(gitDir, 0755)

	tests := []struct {
		name      string
		kernelDir string
		want      bool
	}{
		{
			name:      "Kernel exists with .git",
			kernelDir: kernelDir,
			want:      true,
		},
		{
			name:      "Kernel does not exist",
			kernelDir: "/non/existent",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				Config: &config.Config{
					Paths: config.PathsConfig{KernelDir: tt.kernelDir},
				},
				FS: fs,
			}
			if got := ctx.KernelExists(); got != tt.want {
				t.Errorf("Context.KernelExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_HasConfig(t *testing.T) {
	tmpDir := t.TempDir()
	kernelDir := filepath.Join(tmpDir, "linux")
	configFile := filepath.Join(kernelDir, ".config")

	fs := filesystem.NewOSFileSystem()
	_ = fs.MkdirAll(kernelDir, 0755)
	_ = fs.WriteFile(configFile, []byte("# Kernel config"), 0644)

	tests := []struct {
		name      string
		kernelDir string
		want      bool
	}{
		{
			name:      "Config exists",
			kernelDir: kernelDir,
			want:      true,
		},
		{
			name:      "Config does not exist",
			kernelDir: "/non/existent",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				Config: &config.Config{
					Paths: config.PathsConfig{KernelDir: tt.kernelDir},
				},
				FS: fs,
			}
			if got := ctx.HasConfig(); got != tt.want {
				t.Errorf("Context.HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetKernelImage(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want string
	}{
		{
			name: "arm64 kernel image",
			ctx: &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: "arm64"},
					Paths: config.PathsConfig{KernelDir: "/mnt/linux"},
				},
			},
			want: "/mnt/linux/arch/arm64/boot/Image",
		},
		{
			name: "riscv kernel image",
			ctx: &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: "riscv"},
					Paths: config.PathsConfig{KernelDir: "/mnt/linux"},
				},
			},
			want: "/mnt/linux/arch/riscv/boot/Image",
		},
		{
			name: "Unknown arch returns empty",
			ctx: &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: "unknown"},
					Paths: config.PathsConfig{KernelDir: "/mnt/linux"},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetKernelImage(); got != tt.want {
				t.Errorf("Context.GetKernelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetVmlinux(t *testing.T) {
	ctx := &Context{
		Config: &config.Config{
			Paths: config.PathsConfig{KernelDir: "/mnt/linux"},
		},
	}
	want := "/mnt/linux/vmlinux"
	if got := ctx.GetVmlinux(); got != want {
		t.Errorf("Context.GetVmlinux() = %v, want %v", got, want)
	}
}

func TestContext_HasKernelImage(t *testing.T) {
	tmpDir := t.TempDir()
	kernelDir := filepath.Join(tmpDir, "linux")
	imageDir := filepath.Join(kernelDir, "arch", "arm64", "boot")

	fs := filesystem.NewOSFileSystem()
	_ = fs.MkdirAll(imageDir, 0755)
	_ = fs.WriteFile(filepath.Join(imageDir, "Image"), []byte("kernel"), 0644)

	tests := []struct {
		name      string
		kernelDir string
		arch      string
		want      bool
	}{
		{
			name:      "Image exists",
			kernelDir: kernelDir,
			arch:      "arm64",
			want:      true,
		},
		{
			name:      "Image does not exist",
			kernelDir: "/non/existent",
			arch:      "arm64",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: tt.arch},
					Paths: config.PathsConfig{KernelDir: tt.kernelDir},
				},
				FS: fs,
			}
			if got := ctx.HasKernelImage(); got != tt.want {
				t.Errorf("Context.HasKernelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetMakeEnv(t *testing.T) {
	exec := executor.NewMockExecutor()
	ctx := &Context{
		Config: &config.Config{
			Build: config.BuildConfig{
				Arch:         "arm64",
				CrossCompile: "aarch64-linux-gnu-",
			},
			Paths: config.PathsConfig{LibrariesDir: "/libs"},
		},
		Exec: exec,
		Brew: homebrew.NewResolver(exec),
	}

	env := ctx.GetMakeEnv()

	// Check required env vars
	hasArch := false
	hasLLVM := false
	hasCrossCompile := false
	for _, e := range env {
		if e == "ARCH=arm64" {
			hasArch = true
		}
		if e == "LLVM=1" {
			hasLLVM = true
		}
		if strings.HasPrefix(e, "CROSS_COMPILE=") {
			hasCrossCompile = true
		}
	}

	if !hasArch {
		t.Error("GetMakeEnv() missing ARCH")
	}
	if !hasLLVM {
		t.Error("GetMakeEnv() missing LLVM=1")
	}
	if !hasCrossCompile {
		t.Error("GetMakeEnv() missing CROSS_COMPILE")
	}
}

func TestContext_buildHostCFlags(t *testing.T) {
	exec := executor.NewMockExecutor()
	tests := []struct {
		name         string
		librariesDir string
		wantContains []string
	}{
		{
			name:         "With libraries dir",
			librariesDir: "/mylibs",
			wantContains: []string{"-I/mylibs", "-D_UUID_T", "-D__GETHOSTUUID_H"},
		},
		{
			name:         "Without libraries dir",
			librariesDir: "",
			wantContains: []string{"-D_UUID_T", "-D__GETHOSTUUID_H"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				Config: &config.Config{
					Paths: config.PathsConfig{LibrariesDir: tt.librariesDir},
				},
				Exec: exec,
				Brew: homebrew.NewResolver(exec),
			}
			got := ctx.buildHostCFlags()
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("buildHostCFlags() = %v, want to contain %v", got, want)
				}
			}
		})
	}
}

func TestContext_GetDefaultTargets(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want []string
	}{
		{
			name: "arm64 targets",
			ctx: &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: "arm64"},
				},
			},
			want: []string{"Image", "dtbs", "modules"},
		},
		{
			name: "Unknown arch returns default",
			ctx: &Context{
				Config: &config.Config{
					Build: config.BuildConfig{Arch: "unknown"},
				},
			},
			want: []string{"Image", "dtbs", "modules"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetDefaultTargets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Context.GetDefaultTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test helper functions
func Test_parseMountPointFromHdiutil(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		imagePath string
		want      string
		wantErr   bool
	}{
		{
			name: "Found mount point",
			output: `image-path      : /path/to/test.sparseimage
/dev/disk4      GUID_partition_scheme
/dev/disk4s1    C12A7328-F81F-11D2-BA4B-00A0C93EC93B
/dev/disk5s1    41504653-0000-11AA-AA11-00306543ECAC    /Volumes/test`,
			imagePath: "test.sparseimage",
			want:      "/Volumes/test",
			wantErr:   false,
		},
		{
			name:      "Image not found",
			output:    "some other output",
			imagePath: "test.sparseimage",
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMountPointFromHdiutil(tt.output, tt.imagePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseMountPointFromHdiutil() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("parseMountPointFromHdiutil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findMountPointInLines(t *testing.T) {
	lines := []string{
		"/dev/disk4      GUID",
		"/dev/disk5s1    UUID    /Volumes/myvolume",
		"some other line",
	}

	got := findMountPointInLines(lines, 0, 3)
	want := "/Volumes/myvolume"
	if got != want {
		t.Errorf("findMountPointInLines() = %v, want %v", got, want)
	}

	got = findMountPointInLines(lines, 2, 3)
	if got != "" {
		t.Errorf("findMountPointInLines() = %v, want empty", got)
	}
}
