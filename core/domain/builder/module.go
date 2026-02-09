// Package builder provides kernel and module build orchestration for elmos.
package builder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// ModuleInfo contains information about a kernel module.
type ModuleInfo struct {
	Name        string
	Path        string
	Description string
	Built       bool
}

// ModuleBuilder orchestrates kernel module build operations.
type ModuleBuilder struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
	ctx  *elcontext.Context
	tm   *toolchain.Manager
}

// NewModuleBuilder creates a new ModuleBuilder with the given dependencies.
func NewModuleBuilder(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config, ctx *elcontext.Context, tm *toolchain.Manager) *ModuleBuilder {
	return &ModuleBuilder{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
		ctx:  ctx,
		tm:   tm,
	}
}

// Build builds one or all kernel modules.
func (m *ModuleBuilder) Build(ctx context.Context, name string) error {
	modules, err := m.GetModules(name)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return nil // No modules to build
	}

	for _, mod := range modules {
		if err := m.buildModule(ctx, mod); err != nil {
			return err
		}
	}

	return nil
}

// buildModule builds a single module.
func (m *ModuleBuilder) buildModule(ctx context.Context, mod ModuleInfo) error {
	// Get environment with correct toolchain
	env, crossCompile, err := getToolchainEnv(m.ctx, m.cfg, m.tm, m.fs, m.cfg.Build.Arch)
	if err != nil {
		return fmt.Errorf("failed to configure toolchain environment: %w", err)
	}

	args := []string{
		"-C", m.cfg.Paths.KernelDir,
		fmt.Sprintf("M=%s", mod.Path),
		fmt.Sprintf("ARCH=%s", m.cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", crossCompile),
		"modules",
	}

	return m.exec.RunWithEnv(ctx, env, "make", args...)
}

// Clean cleans one or all kernel modules.
func (m *ModuleBuilder) Clean(ctx context.Context, name string) error {
	modules, err := m.GetModules(name)
	if err != nil {
		return err
	}

	for _, mod := range modules {
		args := []string{
			"-C", m.cfg.Paths.KernelDir,
			fmt.Sprintf("M=%s", mod.Path),
			fmt.Sprintf("ARCH=%s", m.cfg.Build.Arch),
			"clean",
		}

		// Ignore errors during clean
		_ = m.exec.Run(ctx, "make", args...)
	}

	return nil
}

// GetModules returns a list of modules, optionally filtered by name.
func (m *ModuleBuilder) GetModules(name string) ([]ModuleInfo, error) {
	if name != "" {
		return m.getSpecificModule(name)
	}
	return m.getAllModules()
}

// getSpecificModule returns a single module by name.
func (m *ModuleBuilder) getSpecificModule(name string) ([]ModuleInfo, error) {
	modPath := filepath.Join(m.cfg.Paths.ModulesDir, name)
	if !m.fs.Exists(modPath) {
		return nil, fmt.Errorf("module not found: %s", name)
	}

	info := m.getModuleInfo(name, modPath)
	return []ModuleInfo{info}, nil
}

// getAllModules returns all modules in the modules directory.
func (m *ModuleBuilder) getAllModules() ([]ModuleInfo, error) {
	entries, err := m.fs.ReadDir(m.cfg.Paths.ModulesDir)
	if err != nil {
		return nil, err
	}

	var modules []ModuleInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		modPath := filepath.Join(m.cfg.Paths.ModulesDir, name)

		// Check for Makefile
		makePath := filepath.Join(modPath, "Makefile")
		if !m.fs.Exists(makePath) {
			continue
		}

		info := m.getModuleInfo(name, modPath)
		modules = append(modules, info)
	}

	return modules, nil
}

// getModuleInfo builds ModuleInfo for a module.
func (m *ModuleBuilder) getModuleInfo(name, path string) ModuleInfo {
	info := ModuleInfo{
		Name: name,
		Path: path,
	}

	// Check if built
	koFile := filepath.Join(path, name+".ko")
	info.Built = m.fs.Exists(koFile)

	// Extract description from source file
	srcFile := filepath.Join(path, name+".c")
	if content, err := m.fs.ReadFile(srcFile); err == nil {
		info.Description = extractModuleDescription(string(content))
	}

	return info
}

// extractModuleDescription extracts MODULE_DESCRIPTION from source code.
func extractModuleDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "MODULE_DESCRIPTION") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start >= 0 && end > start {
				return line[start+1 : end]
			}
		}
	}
	return ""
}

// PrepareHeaders runs modules_prepare to set up headers for module building.
func (m *ModuleBuilder) PrepareHeaders(ctx context.Context) error {
	// Get environment with correct toolchain
	env, crossCompile, err := getToolchainEnv(m.ctx, m.cfg, m.tm, m.fs, m.cfg.Build.Arch)
	if err != nil {
		return fmt.Errorf("failed to configure toolchain environment: %w", err)
	}

	args := []string{
		"-C", m.cfg.Paths.KernelDir,
		fmt.Sprintf("-j%d", m.cfg.Build.Jobs),
		fmt.Sprintf("ARCH=%s", m.cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", crossCompile),
		"modules_prepare",
	}

	return m.exec.RunWithEnv(ctx, env, "make", args...)
}
