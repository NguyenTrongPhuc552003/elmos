// Package toolchain provides crosstool-ng integration for building cross-compilers.
package toolchain

import (
	"context"
	"fmt"
	"path/filepath"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// Manager handles crosstool-ng toolchain operations.
type Manager struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
}

// NewManager creates a new toolchain Manager.
func NewManager(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config) *Manager {
	return &Manager{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
	}
}

// ToolchainInfo contains information about a built toolchain.
type ToolchainInfo struct {
	Target    string // e.g., "riscv64-unknown-linux-gnu"
	Path      string // Full path to toolchain directory
	Installed bool   // Whether fully built
	Version   string // GCC version if installed
}

// Paths returns important toolchain-related paths.
func (m *Manager) Paths() ToolchainPaths {
	base := m.cfg.Paths.ToolchainsDir
	return ToolchainPaths{
		Base:        base,
		CrosstoolNG: filepath.Join(base, "crosstool-ng"),
		XTools:      filepath.Join(base, "x-tools"),
		Src:         filepath.Join(base, "src"),
		Configs:     filepath.Join(base, "configs"),
	}
}

// ToolchainPaths holds all toolchain directory paths.
type ToolchainPaths struct {
	Base        string // /Volumes/kernel-dev/toolchains
	CrosstoolNG string // ct-ng installation
	XTools      string // Built toolchains output
	Src         string // Downloaded tarballs cache
	Configs     string // Custom .config files
}

// IsInstalled checks if crosstool-ng is installed.
func (m *Manager) IsInstalled() bool {
	ctngBin := filepath.Join(m.Paths().CrosstoolNG, "bin", "ct-ng")
	return m.fs.Exists(ctngBin)
}

// GetCtNgPath returns the path to ct-ng binary.
func (m *Manager) GetCtNgPath() string {
	return filepath.Join(m.Paths().CrosstoolNG, "bin", "ct-ng")
}

// ListSamples lists available crosstool-ng sample configurations.
func (m *Manager) ListSamples(ctx context.Context) ([]string, error) {
	if !m.IsInstalled() {
		return nil, fmt.Errorf("crosstool-ng not installed, run 'elmos toolchains install'")
	}

	output, err := m.exec.Output(ctx, m.GetCtNgPath(), "list-samples")
	if err != nil {
		return nil, fmt.Errorf("failed to list samples: %w", err)
	}

	// Parse output to extract sample names
	// ct-ng list-samples outputs lines like "[L..]   riscv64-unknown-linux-gnu"
	var samples []string
	lines := splitLines(string(output))
	for _, line := range lines {
		if len(line) > 8 && line[0] == '[' {
			// Extract target name after the brackets
			parts := splitWhitespace(line)
			if len(parts) >= 2 {
				samples = append(samples, parts[len(parts)-1])
			}
		}
	}

	return samples, nil
}

// GetInstalledToolchains returns a list of installed toolchains.
func (m *Manager) GetInstalledToolchains() ([]ToolchainInfo, error) {
	xtoolsDir := m.Paths().XTools
	if !m.fs.IsDir(xtoolsDir) {
		return nil, nil
	}

	entries, err := m.fs.ReadDir(xtoolsDir)
	if err != nil {
		return nil, err
	}

	var toolchains []ToolchainInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		target := entry.Name()
		path := filepath.Join(xtoolsDir, target)

		// Check if bin directory exists (indicates successful build)
		binDir := filepath.Join(path, "bin")
		installed := m.fs.IsDir(binDir)

		toolchains = append(toolchains, ToolchainInfo{
			Target:    target,
			Path:      path,
			Installed: installed,
		})
	}

	return toolchains, nil
}

// splitLines splits a string by newlines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// splitWhitespace splits a string by whitespace.
func splitWhitespace(s string) []string {
	var parts []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				parts = append(parts, s[start:i])
				start = -1
			}
		} else {
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		parts = append(parts, s[start:])
	}
	return parts
}
