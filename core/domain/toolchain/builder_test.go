// Package toolchain provides crosstool-ng integration for building cross-compilers.
package toolchain

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

// --- Local Mocks ---

type mockExecutor struct {
	executor.Executor
	runFunc             func(ctx context.Context, cmd string, args ...string) error
	runInDirFunc        func(ctx context.Context, dir string, cmd string, args ...string) error
	runWithEnvInDirFunc func(ctx context.Context, env []string, dir string, cmd string, args ...string) error
	outputFunc          func(ctx context.Context, cmd string, args ...string) ([]byte, error)
}

func (m *mockExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	if m.runFunc != nil {
		return m.runFunc(ctx, cmd, args...)
	}
	return nil
}

func (m *mockExecutor) RunInDir(ctx context.Context, dir string, cmd string, args ...string) error {
	if m.runInDirFunc != nil {
		return m.runInDirFunc(ctx, dir, cmd, args...)
	}
	return nil
}

func (m *mockExecutor) RunWithEnvInDir(ctx context.Context, env []string, dir string, cmd string, args ...string) error {
	if m.runWithEnvInDirFunc != nil {
		return m.runWithEnvInDirFunc(ctx, env, dir, cmd, args...)
	}
	return nil
}

func (m *mockExecutor) Output(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	if m.outputFunc != nil {
		return m.outputFunc(ctx, cmd, args...)
	}
	return []byte{}, nil
}

type mockFileSystem struct {
	filesystem.FileSystem
	existsFunc    func(path string) bool
	isDirFunc     func(path string) bool
	mkdirAllFunc  func(path string, perm os.FileMode) error
	readFileFunc  func(name string) ([]byte, error)
	writeFileFunc func(name string, data []byte, perm os.FileMode) error
	readDirFunc   func(name string) ([]fs.DirEntry, error)
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

func (m *mockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.mkdirAllFunc != nil {
		return m.mkdirAllFunc(path, perm)
	}
	return nil
}

func (m *mockFileSystem) ReadFile(name string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(name)
	}
	return nil, fmt.Errorf("file not found: %s", name)
}

func (m *mockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(name, data, perm)
	}
	return nil
}

func (m *mockFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	if m.readDirFunc != nil {
		return m.readDirFunc(name)
	}
	return nil, nil
}

// Stubs to satisfy interface
func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) { return nil, nil }
func (m *mockFileSystem) Remove(name string) error              { return nil }
func (m *mockFileSystem) RemoveAll(path string) error           { return nil }
func (m *mockFileSystem) Getwd() (string, error)                { return "/tmp", nil }
func (m *mockFileSystem) Create(name string) (*os.File, error)  { return nil, nil }
func (m *mockFileSystem) Open(name string) (*os.File, error)    { return nil, nil }

// --- Tests ---

func TestManager_Install(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		mockEx  *mockExecutor
		mockFS  *mockFileSystem
		wantErr bool
	}{
		{
			name: "Success - Clone and Install",
			mockEx: &mockExecutor{
				runFunc: func(ctx context.Context, cmd string, args ...string) error {
					return nil // git clone success
				},
				runWithEnvInDirFunc: func(ctx context.Context, env []string, dir, cmd string, args ...string) error {
					return nil // make/bootstrap/configure success
				},
			},
			mockFS: &mockFileSystem{
				mkdirAllFunc: func(path string, perm os.FileMode) error { return nil },
				isDirFunc:    func(path string) bool { return false }, // Not cloned yet
			},
			wantErr: false,
		},
		{
			name: "Failure - Mkdir",
			mockFS: &mockFileSystem{
				mkdirAllFunc: func(path string, perm os.FileMode) error { return fmt.Errorf("perm error") },
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				exec: tt.mockEx,
				fs:   tt.mockFS,
				cfg:  &elconfig.Config{},
			}
			if tt.mockEx == nil {
				m.exec = &mockExecutor{}
			}
			if err := m.Install(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SelectTarget(t *testing.T) {
	ctx := context.Background()

	// Create a temp file for custom config test
	tmpFile, err := os.CreateTemp("", "elmos-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("CT_TARGET=custom")
	tmpFile.Close()

	tests := []struct {
		name         string
		target       string
		mockEx       *mockExecutor
		mockFS       *mockFileSystem
		customConfig string // Path to real file for custom config test
		wantErr      bool
	}{
		{
			name:   "Success - Standard Target",
			target: "x86_64-unknown-linux-gnu",
			mockEx: &mockExecutor{
				outputFunc: func(ctx context.Context, cmd string, args ...string) ([]byte, error) {
					return []byte("bin/ct-ng"), nil // LookPath
				},
				runWithEnvInDirFunc: func(ctx context.Context, env []string, dir, cmd string, args ...string) error {
					return nil // ct-ng run success
				},
			},
			mockFS: &mockFileSystem{
				mkdirAllFunc: func(path string, perm os.FileMode) error { return nil },
				existsFunc: func(path string) bool {
					// Only return true for installed check, NOT for config file check
					return strings.Contains(path, "bin/ct-ng")
				},
				readFileFunc:  func(name string) ([]byte, error) { return []byte(""), nil }, // patchConfig
				writeFileFunc: func(name string, data []byte, perm os.FileMode) error { return nil },
			},
			wantErr: false,
		},
		{
			name:         "Success - Custom Config",
			target:       "custom-arch",
			customConfig: tmpFile.Name(),
			mockEx:       &mockExecutor{},
			mockFS: &mockFileSystem{
				mkdirAllFunc: func(path string, perm os.FileMode) error { return nil },
				existsFunc: func(path string) bool {
					// Return true for the temp file path so GetCustomConfigPath finds it
					if path == tmpFile.Name() {
						return true
					}
					// Return true for .config existence check (installed check)
					if strings.Contains(path, "bin/ct-ng") {
						return true
					}
					return false
				},
				// patchConfig uses ReadFile
				readFileFunc: func(name string) ([]byte, error) {
					return []byte("CT_TARGET=custom"), nil
				},
				writeFileFunc: func(name string, data []byte, perm os.FileMode) error { return nil },
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				exec:    tt.mockEx,
				fs:      tt.mockFS,
				cfg:     &elconfig.Config{},
				printer: ui.NewPrinter(),
			}

			// Setup for custom config test
			if tt.customConfig != "" {
				// We need GetCustomConfigPath to return tt.customConfig.
				// GetCustomConfigPath checks cfg.Paths.ProjectRoot + ... and m.Paths().Configs + ...
				// We can hijack this by setting cfg.Paths.ProjectRoot to the dir of tmpFile?
				// Complex.
				// Alternative: The test logic above mocks Exists(tmpFile.Name()) -> true.
				// But GetCustomConfigPath constructs path: `filepath.Join(projectConfigs, target+".config")`.
				// `projectConfigs` is `Root/tools/toolchains/configs`.
				// So we need `Root/.../custom-arch.config` to equal `tmpFile.Name()`.
				// This is only possible if we create the file AT the expected location.

				// Better approach: skip deep custom config path logic test and trust unit tests for GetCustomConfigPath elsewhere,
				// OR create the directory structure expected.
				// For this unit test, let's trick it. We can set Paths.Configs to the dir containing our temp file.
				// And ensure the temp file uses the correct name format.

				dir := filepath.Dir(tt.customConfig)
				// base := filepath.Base(tt.customConfig)
				// targetName := strings.TrimSuffix(base, ".config")

				m.cfg.Paths.ToolchainsDir = dir // ToolchainsDir/configs is checked
				// We need to match the filename which is randomized by CreateTemp.
				// Let's just create a file with the correct name in a temp dir.

				tmpDir, _ := os.MkdirTemp("", "elmos-toolchain")
				defer os.RemoveAll(tmpDir)
				m.cfg.Paths.ToolchainsDir = tmpDir

				// Create the custom config file where checking logic expects it: tmpDir/configs/target.config
				configDir := filepath.Join(tmpDir, "configs")
				os.MkdirAll(configDir, 0755)
				customConfPath := filepath.Join(configDir, tt.target+".config")
				os.WriteFile(customConfPath, []byte("CT_CUSTOM=y"), 0644)

				// Update mockFS to allow reading this real file? No, we need m.fs.Exists to return true for it.
				// But we also need os.ReadFile to read it (which it will).
				// So we need mockFS.Exists to delegate to real FS or return true for this path.

				// Let's define specific mock behavior inside the run loop
				tt.mockFS.existsFunc = func(path string) bool {
					if path == customConfPath {
						return true
					}
					if strings.Contains(path, "bin/ct-ng") {
						return true
					} // Installed check
					return false
				}
			}

			if err := m.SelectTarget(ctx, tt.target); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SelectTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Build(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		mockEx  *mockExecutor
		mockFS  *mockFileSystem
		jobs    int
		wantErr bool
	}{
		{
			name: "Success",
			mockEx: &mockExecutor{
				runWithEnvInDirFunc: func(ctx context.Context, env []string, dir, cmd string, args ...string) error {
					return nil
				},
			},
			mockFS: &mockFileSystem{
				existsFunc: func(path string) bool { return true }, // Installed + .config exists
			},
			jobs:    4,
			wantErr: false,
		},
		{
			name: "Not Configured",
			mockFS: &mockFileSystem{
				existsFunc: func(path string) bool {
					if strings.Contains(path, "bin/ct-ng") {
						return true
					} // Installed
					return false // .config missing
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{exec: tt.mockEx, fs: tt.mockFS, cfg: &elconfig.Config{}}
			// Fix for jobs=0 defaults to runtime.NumCPU()
			jobs := tt.jobs
			if jobs == 0 {
				jobs = runtime.NumCPU()
			}

			if err := m.Build(ctx, tt.jobs); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Clean(t *testing.T) {
	ctx := context.Background()
	mockEx := &mockExecutor{
		runInDirFunc: func(ctx context.Context, dir, cmd string, args ...string) error {
			if !strings.HasSuffix(cmd, "ct-ng") || args[0] != "clean" {
				return fmt.Errorf("unexpected command: %s %v", cmd, args)
			}
			return nil
		},
	}
	mockFS := &mockFileSystem{existsFunc: func(p string) bool { return true }}

	m := &Manager{exec: mockEx, fs: mockFS, cfg: &elconfig.Config{}}
	if err := m.Clean(ctx); err != nil {
		t.Errorf("Manager.Clean() error = %v", err)
	}
}

func TestManager_patchConfig(t *testing.T) {
	mockFS := &mockFileSystem{
		readFileFunc: func(name string) ([]byte, error) {
			return []byte("CT_PREFIX_DIR=\"old\""), nil
		},
		writeFileFunc: func(name string, data []byte, perm os.FileMode) error {
			if !strings.Contains(string(data), "CT_PREFIX_DIR=\"xtools/${CT_TARGET}\"") {
				return fmt.Errorf("unexpected data: %s", string(data))
			}
			return nil
		},
	}
	m := &Manager{fs: mockFS}
	paths := ToolchainPaths{XTools: "xtools"}

	if err := m.patchConfig("conf", paths); err != nil {
		t.Errorf("patchConfig() error = %v", err)
	}
}

func Test_getInstallEnv(t *testing.T) {
	m := &Manager{}
	env := m.getInstallEnv()
	if len(env) == 0 {
		t.Error("getInstallEnv() returned empty env")
	}
	foundPath := false
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			foundPath = true
			if !strings.Contains(e, "opt/binutils/bin") {
				t.Error("getInstallEnv() PATH missing binutils")
			}
		}
	}
	if !foundPath {
		t.Error("getInstallEnv() missing PATH")
	}
}

func TestManager_getInstallEnv(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getInstallEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getInstallEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Menuconfig(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Menuconfig(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Menuconfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_patchConfigContent(t *testing.T) {
	type args struct {
		content string
		paths   ToolchainPaths
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := patchConfigContent(tt.args.content, tt.args.paths); got != tt.want {
				t.Errorf("patchConfigContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceConfigValue(t *testing.T) {
	type args struct {
		content string
		key     string
		value   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceConfigValue(tt.args.content, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("replaceConfigValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceAll(t *testing.T) {
	type args struct {
		s   string
		old string
		new string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceAll(tt.args.s, tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("replaceAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexOf(t *testing.T) {
	type args struct {
		s      string
		substr string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexOf(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("indexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getBuildEnv(t *testing.T) {
	type args struct {
		paths ToolchainPaths
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getBuildEnv(tt.args.paths); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getBuildEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBrewPrefix(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBrewPrefix(); got != tt.want {
				t.Errorf("getBrewPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_ensureLocalBin(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ensureLocalBin(); got != tt.want {
				t.Errorf("Manager.ensureLocalBin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_ensureGCCSymlinks(t *testing.T) {
	type args struct {
		localBin   string
		brewPrefix string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.ensureGCCSymlinks(tt.args.localBin, tt.args.brewPrefix)
		})
	}
}

func TestManager_addLibraryFlags(t *testing.T) {
	type args struct {
		env        []string
		brewPrefix string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.addLibraryFlags(tt.args.env, tt.args.brewPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.addLibraryFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_buildFlagString(t *testing.T) {
	type args struct {
		brewPrefix string
		pkgs       []string
		subdir     string
		flag       string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.buildFlagString(tt.args.brewPrefix, tt.args.pkgs, tt.args.subdir, tt.args.flag); got != tt.want {
				t.Errorf("Manager.buildFlagString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_addGmake(t *testing.T) {
	type args struct {
		env        []string
		brewPrefix string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.addGmake(tt.args.env, tt.args.brewPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.addGmake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_buildPathEnv(t *testing.T) {
	type args struct {
		env        []string
		localBin   string
		brewPrefix string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.buildPathEnv(tt.args.env, tt.args.localBin, tt.args.brewPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.buildPathEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendOrUpdateEnv(t *testing.T) {
	type args struct {
		env   []string
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendOrUpdateEnv(tt.args.env, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendOrUpdateEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
