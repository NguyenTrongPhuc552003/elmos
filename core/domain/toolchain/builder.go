// Package toolchain provides crosstool-ng integration for building cross-compilers.
package toolchain

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
)

// Build builds the currently configured toolchain.
func (m *Manager) Build(ctx context.Context, jobs int) error {
	if !m.IsInstalled() {
		return fmt.Errorf("crosstool-ng not installed, run 'elmos toolchains install'")
	}

	paths := m.Paths()

	// Check for .config
	configFile := filepath.Join(paths.Base, ".config")
	if !m.fs.Exists(configFile) {
		return fmt.Errorf("no target selected, run 'elmos toolchains <target>' first")
	}

	// Setup environment
	env := m.getBuildEnv(paths)

	// Build with specified jobs
	if jobs <= 0 {
		jobs = runtime.NumCPU()
	}

	buildTarget := fmt.Sprintf("build.%d", jobs)
	if err := m.exec.RunWithEnvInDir(ctx, env, paths.Base, m.GetCtNgPath(), buildTarget); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}

// Clean cleans the build artifacts.
func (m *Manager) Clean(ctx context.Context) error {
	if !m.IsInstalled() {
		return nil
	}

	paths := m.Paths()
	return m.exec.RunInDir(ctx, paths.Base, m.GetCtNgPath(), "clean")
}

// Menuconfig opens the interactive configuration menu.
func (m *Manager) Menuconfig(ctx context.Context) error {
	if !m.IsInstalled() {
		return fmt.Errorf("crosstool-ng not installed, run 'elmos toolchains install'")
	}

	paths := m.Paths()
	return m.exec.RunInDir(ctx, paths.Base, m.GetCtNgPath(), "menuconfig")
}
