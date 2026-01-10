// Package toolchain provides crosstool-ng integration for building cross-compilers.
package toolchain

import (
	"context"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestNewManager(t *testing.T) {
	mockEx := &mockExecutor{}
	mockFS := &mockFileSystem{}
	cfg := &elconfig.Config{}
	printer := ui.NewPrinter()

	type args struct {
		exec    executor.Executor
		fs      filesystem.FileSystem
		cfg     *elconfig.Config
		printer *ui.Printer
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "Initialization",
			args: args{mockEx, mockFS, cfg, printer},
			want: &Manager{exec: mockEx, fs: mockFS, cfg: cfg, printer: printer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.printer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Paths(t *testing.T) {
	cfg := &elconfig.Config{
		Paths: elconfig.PathsConfig{
			ToolchainsDir: "/toolchains",
		},
	}
	m := &Manager{cfg: cfg}

	want := ToolchainPaths{
		Base:        "/toolchains",
		CrosstoolNG: "/toolchains/crosstool-ng",
		XTools:      "/toolchains/x-tools",
		Src:         "/toolchains/src",
		Configs:     "/toolchains/configs",
	}

	if got := m.Paths(); !reflect.DeepEqual(got, want) {
		t.Errorf("Manager.Paths() = %v, want %v", got, want)
	}
}

func TestManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		want   bool
	}{
		{"Installed", true, true},
		{"Not Installed", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockFileSystem{
				existsFunc: func(path string) bool {
					return tt.exists
				},
			}
			m := &Manager{fs: mockFS, cfg: &elconfig.Config{}}
			if got := m.IsInstalled(); got != tt.want {
				t.Errorf("Manager.IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCtNgPath(t *testing.T) {
	m := &Manager{cfg: &elconfig.Config{
		Paths: elconfig.PathsConfig{ToolchainsDir: "/tc"},
	}}
	want := "/tc/crosstool-ng/bin/ct-ng"
	if got := m.GetCtNgPath(); got != want {
		t.Errorf("Manager.GetCtNgPath() = %v, want %v", got, want)
	}
}

func TestManager_ListSamples(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		installed bool
		output    string
		want      []string
		wantErr   bool
	}{
		{
			name:      "Success",
			installed: true,
			output:    "[L..]   riscv64-unknown-linux-gnu\n[G..]   arm-unknown-linux-gnueabi",
			want:      []string{"riscv64-unknown-linux-gnu", "arm-unknown-linux-gnueabi"},
			wantErr:   false,
		},
		{
			name:      "Not Installed",
			installed: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockFileSystem{existsFunc: func(s string) bool { return tt.installed }}
			mockEx := &mockExecutor{
				outputFunc: func(ctx context.Context, cmd string, args ...string) ([]byte, error) {
					return []byte(tt.output), nil
				},
			}
			m := &Manager{exec: mockEx, fs: mockFS, cfg: &elconfig.Config{}}

			got, err := m.ListSamples(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.ListSamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.ListSamples() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCustomConfigPath(t *testing.T) {
	cfg := &elconfig.Config{
		Paths: elconfig.PathsConfig{
			ProjectRoot:   "/root",
			ToolchainsDir: "/tc",
		},
	}

	tests := []struct {
		name         string
		target       string
		existsInRoot bool
		existsInTc   bool
		want         string
	}{
		{
			name:         "Found in Project Root",
			target:       "t1",
			existsInRoot: true,
			want:         "/root/tools/toolchains/configs/t1.config",
		},
		{
			name:         "Found in Toolchains Dir",
			target:       "t2",
			existsInRoot: false,
			existsInTc:   true,
			want:         "/tc/configs/t2.config",
		},
		{
			name:         "Not Found",
			target:       "t3",
			existsInRoot: false,
			existsInTc:   false,
			want:         "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockFileSystem{
				existsFunc: func(path string) bool {
					if path == filepath.Join("/root/tools/toolchains/configs", tt.target+".config") {
						return tt.existsInRoot
					}
					if path == filepath.Join("/tc/configs", tt.target+".config") {
						return tt.existsInTc
					}
					return false
				},
			}
			m := &Manager{fs: mockFS, cfg: cfg}
			if got := m.GetCustomConfigPath(tt.target); got != tt.want {
				t.Errorf("Manager.GetCustomConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// simpleDirEntry implements fs.DirEntry for testing
type simpleDirEntry struct {
	name  string
	isDir bool
}

func (d simpleDirEntry) Name() string               { return d.name }
func (d simpleDirEntry) IsDir() bool                { return d.isDir }
func (d simpleDirEntry) Type() fs.FileMode          { return 0 }
func (d simpleDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func TestManager_GetInstalledToolchains(t *testing.T) {
	tests := []struct {
		name      string
		dirExists bool
		entries   []fs.DirEntry
		binExists bool
		want      []ToolchainInfo
		wantErr   bool
	}{
		{
			name:      "Success",
			dirExists: true,
			entries: []fs.DirEntry{
				simpleDirEntry{name: "arch1", isDir: true},
				simpleDirEntry{name: "file1", isDir: false},
			},
			binExists: true,
			want: []ToolchainInfo{
				{Target: "arch1", Path: "/tc/x-tools/arch1", Installed: true},
			},
			wantErr: false,
		},
		{
			name:      "No XTools Dir",
			dirExists: false,
			want:      nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockFileSystem{
				isDirFunc: func(path string) bool {
					if path == "/tc/x-tools" {
						return tt.dirExists
					}
					if strings.HasSuffix(path, "bin") {
						return tt.binExists
					}
					return false
				},
				readDirFunc: func(name string) ([]fs.DirEntry, error) {
					return tt.entries, nil
				},
			}
			m := &Manager{fs: mockFS, cfg: &elconfig.Config{
				Paths: elconfig.PathsConfig{ToolchainsDir: "/tc"},
			}}

			got, err := m.GetInstalledToolchains()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetInstalledToolchains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetInstalledToolchains() = %v, want %v", got, tt.want)
			}
		})
	}
}
