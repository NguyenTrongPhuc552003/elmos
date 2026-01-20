// Package patch provides kernel patch management for elmos.
package patch

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// Manager handles kernel patch operations.
type Manager struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
}

// NewManager creates a new patch Manager.
func NewManager(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config) *Manager {
	return &Manager{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
	}
}

// Apply applies a patch file to the kernel source.
func (m *Manager) Apply(ctx context.Context, patchFile string) error {
	// Resolve patch path
	patchPath := patchFile

	// Handle absolute paths
	if filepath.IsAbs(patchPath) {
		if !m.fs.Exists(patchPath) {
			return fmt.Errorf("patch file not found: %s", patchFile)
		}
	} else {
		// Strip redundant patches/ prefix if user included it
		patchPath = strings.TrimPrefix(patchPath, "patches/")

		// Build full path from patches directory
		fullPath := filepath.Join(m.cfg.Paths.PatchesDir, patchPath)
		if !m.fs.Exists(fullPath) {
			return fmt.Errorf("patch file not found: %s", patchFile)
		}
		patchPath = fullPath
	}

	// Check if patch is already applied
	checkArgs := []string{"-p1", "--dry-run", "-i", patchPath}
	err := m.exec.RunInDir(ctx, m.cfg.Paths.KernelDir, "patch", checkArgs...)
	if err != nil {
		// Try reverse check to see if already applied
		reverseArgs := []string{"-p1", "--dry-run", "-R", "-i", patchPath}
		if m.exec.RunInDir(ctx, m.cfg.Paths.KernelDir, "patch", reverseArgs...) == nil {
			return fmt.Errorf("patch appears to already be applied")
		}
		return fmt.Errorf("patch cannot be applied cleanly: %w", err)
	}

	// Apply the patch
	applyArgs := []string{"-p1", "-i", patchPath}
	if err := m.exec.RunInDir(ctx, m.cfg.Paths.KernelDir, "patch", applyArgs...); err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	return nil
}

// Reverse reverses a previously applied patch.
func (m *Manager) Reverse(ctx context.Context, patchFile string) error {
	patchPath := patchFile

	if filepath.IsAbs(patchPath) {
		if !m.fs.Exists(patchPath) {
			return fmt.Errorf("patch file not found: %s", patchFile)
		}
	} else {
		patchPath = strings.TrimPrefix(patchPath, "patches/")
		fullPath := filepath.Join(m.cfg.Paths.PatchesDir, patchPath)
		if !m.fs.Exists(fullPath) {
			return fmt.Errorf("patch file not found: %s", patchFile)
		}
		patchPath = fullPath
	}

	reverseArgs := []string{"-p1", "-R", "-i", patchPath}
	if err := m.exec.RunInDir(ctx, m.cfg.Paths.KernelDir, "patch", reverseArgs...); err != nil {
		return fmt.Errorf("failed to reverse patch: %w", err)
	}

	return nil
}

// List returns all available patches.
func (m *Manager) List() ([]PatchInfo, error) {
	if !m.fs.Exists(m.cfg.Paths.PatchesDir) {
		return nil, nil
	}

	versionDirs, err := m.fs.ReadDir(m.cfg.Paths.PatchesDir)
	if err != nil {
		return nil, err
	}

	var patches []PatchInfo
	for _, versionEntry := range versionDirs {
		if !versionEntry.IsDir() {
			continue
		}
		versionPatches := m.listPatchesForVersion(versionEntry.Name())
		patches = append(patches, versionPatches...)
	}

	return patches, nil
}

// listPatchesForVersion lists all patches for a specific version directory.
func (m *Manager) listPatchesForVersion(version string) []PatchInfo {
	versionDir := filepath.Join(m.cfg.Paths.PatchesDir, version)
	archDirs, err := m.fs.ReadDir(versionDir)
	if err != nil {
		return nil
	}

	var patches []PatchInfo
	for _, archEntry := range archDirs {
		if !archEntry.IsDir() {
			continue
		}
		archPatches := m.listPatchesForArch(version, archEntry.Name())
		patches = append(patches, archPatches...)
	}
	return patches
}

// listPatchesForArch lists all patches for a specific arch directory.
func (m *Manager) listPatchesForArch(version, arch string) []PatchInfo {
	archDir := filepath.Join(m.cfg.Paths.PatchesDir, version, arch)
	patchFiles, err := m.fs.ReadDir(archDir)
	if err != nil {
		return nil
	}

	var patches []PatchInfo
	for _, pf := range patchFiles {
		if pf.IsDir() || !strings.HasSuffix(pf.Name(), ".patch") {
			continue
		}
		patches = append(patches, PatchInfo{
			Name:    pf.Name(),
			Path:    filepath.Join(archDir, pf.Name()),
			Version: version,
			Arch:    arch,
		})
	}
	return patches
}

// GetPatchesForVersion returns patches for a specific kernel version.
func (m *Manager) GetPatchesForVersion(version string) ([]PatchInfo, error) {
	allPatches, err := m.List()
	if err != nil {
		return nil, err
	}

	var filtered []PatchInfo
	for _, p := range allPatches {
		if p.Version == version {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}
