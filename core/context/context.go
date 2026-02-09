// Package context provides build context management for elmos.
package context

import (
	gocontext "context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/environment"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/packages"
)

// Context holds the current build context and state.
type Context struct {
	Config     *config.Config
	Exec       executor.Executor
	FS         filesystem.FileSystem
	PkgRes     packages.Resolver
	EnvBuilder *environment.Builder
	Verbose    bool
}

// New creates a new build context with the given dependencies.
func New(cfg *config.Config, exec executor.Executor, fs filesystem.FileSystem, pkgRes packages.Resolver) *Context {
	return &Context{
		Config:     cfg,
		Exec:       exec,
		FS:         fs,
		PkgRes:     pkgRes,
		EnvBuilder: environment.New(cfg, pkgRes),
	}
}

// IsMounted checks if the kernel volume is currently mounted.
func (ctx *Context) IsMounted() bool {
	mountPoint := ctx.Config.Image.MountPoint

	// Check if the directory exists first
	if !ctx.FS.IsDir(mountPoint) {
		// It might be mounted at a different location
		// Check hdiutil info for our image file path
		out, err := ctx.Exec.Output(gocontext.Background(), "hdiutil", "info")
		if err != nil {
			return false
		}
		// Check if our image file is mounted (regardless of mount point)
		return strings.Contains(string(out), ctx.Config.Image.Path)
	}

	// Verify it's actually a mount point using 'mount'
	out, err := ctx.Exec.Output(gocontext.Background(), "mount")
	if err != nil {
		return false
	}

	return strings.Contains(string(out), mountPoint)
}

// EnsureMounted ensures the kernel volume is mounted.
func (ctx *Context) EnsureMounted() error {
	if !ctx.IsMounted() {
		return ImageError("kernel volume not mounted", ErrNotMounted)
	}
	return nil
}

// GetMountPoint returns the actual mount point path of the kernel volume.
// This handles cases where the volume is mounted at a different location (e.g. " 1" suffix).
func (ctx *Context) GetActualMountPoint() (string, error) {
	mountPoint := ctx.Config.Image.MountPoint

	// Fast path: if configured path exists
	if ctx.FS.IsDir(mountPoint) {
		return mountPoint, nil
	}

	// Slow path: check hdiutil info for our specific image file
	out, err := ctx.Exec.Output(gocontext.Background(), "hdiutil", "info")
	if err != nil {
		return "", err
	}

	return parseMountPointFromHdiutil(string(out), ctx.Config.Image.Path)
}

// parseMountPointFromHdiutil extracts the mount point for an image from hdiutil info output.
func parseMountPointFromHdiutil(output, imagePath string) (string, error) {
	lines := strings.Split(output, "\n")

	// Find the image block and look for mount point
	foundImage := false
	for i, line := range lines {
		if strings.Contains(line, imagePath) {
			foundImage = true
			// Look for /Volumes/ in subsequent lines (up to 20 lines)
			if mp := findMountPointInLines(lines, i+1, i+20); mp != "" {
				return mp, nil
			}
			break
		}
	}

	if !foundImage {
		return "", fmt.Errorf("image not mounted: %s", imagePath)
	}
	return "", fmt.Errorf("volume not found")
}

// findMountPointInLines searches for a /Volumes/ path in a range of lines.
func findMountPointInLines(lines []string, start, end int) string {
	for j := start; j < len(lines) && j < end; j++ {
		if idx := strings.Index(lines[j], "/Volumes/"); idx != -1 {
			mountStr := strings.TrimSpace(lines[j][idx:])
			if parts := strings.Fields(mountStr); len(parts) > 0 {
				return parts[0]
			}
		}
	}
	return ""
}

// KernelExists checks if the kernel source directory exists.
func (ctx *Context) KernelExists() bool {
	gitDir := filepath.Join(ctx.Config.Paths.KernelDir, ".git")
	return ctx.FS.Exists(gitDir)
}

// HasConfig checks if the kernel has been configured (.config exists).
func (ctx *Context) HasConfig() bool {
	configFile := filepath.Join(ctx.Config.Paths.KernelDir, ".config")
	return ctx.FS.Exists(configFile)
}

// GetKernelImage returns the path to the built kernel image for the current arch.
func (ctx *Context) GetKernelImage() string {
	archCfg := ctx.Config.GetArchConfig()
	if archCfg == nil {
		return ""
	}
	return filepath.Join(ctx.Config.Paths.KernelDir, "arch", archCfg.KernelArch, "boot", archCfg.KernelImage)
}

// GetVmlinux returns the path to vmlinux (for debugging).
func (ctx *Context) GetVmlinux() string {
	return filepath.Join(ctx.Config.Paths.KernelDir, "vmlinux")
}

// HasKernelImage checks if the kernel image has been built.
func (ctx *Context) HasKernelImage() bool {
	return ctx.FS.Exists(ctx.GetKernelImage())
}

// GetMakeEnv returns environment variables for kernel make commands.
// GetMakeEnv returns environment variables for kernel make commands.
// Delegates to EnvironmentBuilder for actual construction.
func (ctx *Context) GetMakeEnv() []string {
	return ctx.EnvBuilder.BuildMakeEnv()
}

// GetDefaultTargets returns the default build targets for the current architecture.
func (ctx *Context) GetDefaultTargets() []string {
	archCfg := ctx.Config.GetArchConfig()
	if archCfg == nil {
		return []string{"Image", "dtbs", "modules"}
	}
	return archCfg.DefaultTargets
}
