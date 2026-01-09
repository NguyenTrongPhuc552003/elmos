// Package rootfs provides rootfs creation functionality for elmos.
package rootfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// Creator handles rootfs creation.
type Creator struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
}

// NewCreator creates a new rootfs Creator.
func NewCreator(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config) *Creator {
	return &Creator{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
	}
}

// Create creates an ext4 disk image with Debian rootfs.
func (c *Creator) Create(ctx context.Context, opts CreateOptions) error {
	size := opts.Size
	if size == "" {
		size = "5G"
	}

	diskImage := c.cfg.Paths.DiskImage
	rootfsDir := c.cfg.Paths.RootfsDir

	// Step 1: Clean old rootfs directory (must use sudo since debootstrap creates root-owned files)
	if c.fs.Exists(rootfsDir) {
		if err := c.exec.Run(ctx, "sudo", "rm", "-rf", rootfsDir); err != nil {
			return fmt.Errorf("failed to clean old rootfs: %w", err)
		}
	}

	// Create fresh rootfs directory
	if err := c.exec.Run(ctx, "mkdir", "-p", rootfsDir); err != nil {
		return fmt.Errorf("failed to create rootfs directory: %w", err)
	}

	// Step 2: Get Debian architecture and run debootstrap
	arch := c.getDebianArch()
	debootstrapDir := filepath.Join(c.cfg.Paths.ProjectRoot, "tools", "debootstrap")
	debootstrapPath := filepath.Join(debootstrapDir, "debootstrap")
	if !c.fs.Exists(debootstrapPath) {
		return fmt.Errorf("debootstrap not found at %s", debootstrapPath)
	}

	// Execute debootstrap with DEBOOTSTRAP_DIR set via env command
	if err := c.exec.Run(ctx,
		"sudo", "-E", "DEBOOTSTRAP_DIR="+debootstrapDir,
		"fakeroot", debootstrapPath,
		"--foreign",
		"--arch="+arch,
		"--no-check-gpg",
		"stable",
		rootfsDir,
		c.cfg.Paths.DebianMirror,
	); err != nil {
		return fmt.Errorf("debootstrap failed: %w", err)
	}

	// Step 3: Create init script in rootfs directory
	if err := c.createInitScript(rootfsDir); err != nil {
		return fmt.Errorf("failed to create init script: %w", err)
	}

	// Step 4: Remove old disk image if exists
	if c.fs.Exists(diskImage) {
		if err := c.exec.Run(ctx, "rm", "-f", diskImage); err != nil {
			return fmt.Errorf("failed to remove old disk image: %w", err)
		}
	}

	// Step 4.5: Fix apt list filenames (host debootstrap may include protocol prefix, guest might not expect it)
	if err := c.fixAptLists(rootfsDir); err != nil {
		fmt.Printf("Warning: failed to fix apt lists: %v\n", err)
	}

	// Step 5: Create ext4 disk image and populate it from rootfs directory using mke2fs -d
	// Using sudo since debootstrap created root-owned files that mke2fs needs to copy
	if err := c.exec.Run(ctx,
		"sudo", "mke2fs", "-t", "ext4",
		"-E", "lazy_itable_init=0,lazy_journal_init=0",
		"-d", rootfsDir,
		diskImage, size,
	); err != nil {
		return fmt.Errorf("failed to create disk image: %w", err)
	}

	return nil
}

// fixAptLists smoothes over differences between debootstrap versions by symlinking
// http:__ prefix files to their non-prefixed counterparts.
func (c *Creator) fixAptLists(rootfsDir string) error {
	listsDir := filepath.Join(rootfsDir, "var", "lib", "apt", "lists")
	entries, err := os.ReadDir(listsDir)
	if err != nil {
		return nil // Directory might not exist yet, ignore
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			continue
		}

		newName := name
		if len(name) > 7 && name[:7] == "http:__" {
			newName = name[7:]
		} else if len(name) > 8 && name[:8] == "https:__" {
			newName = name[8:]
		}

		if newName != name {
			newPath := filepath.Join(listsDir, newName)
			// Create symlink if it doesn't match
			// We need to use sudo because the directory is owned by root (created by sudo debootstrap)
			// Use absolute path for target to be safe, or just filename if relative
			// Here we use just name because they are in the same directory
			_ = c.exec.Run(context.Background(), "sudo", "ln", "-sf", name, newPath)
		}
	}
	return nil
}

// getDebianArch returns the Debian architecture name.
func (c *Creator) getDebianArch() string {
	switch c.cfg.Build.Arch {
	case "arm64":
		return "arm64"
	case "arm":
		return "armhf"
	case "riscv":
		return "riscv64"
	default:
		return "arm64"
	}
}

// createInitScript creates the /init script for the rootfs.
func (c *Creator) createInitScript(rootfsDir string) error {
	initPath := filepath.Join(rootfsDir, "init")

	// Read init script from scripts/init file
	scriptPath := filepath.Join(c.cfg.Paths.ProjectRoot, "scripts", "init")
	content, err := c.fs.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read init script from %s: %w", scriptPath, err)
	}

	if err := c.fs.WriteFile(initPath, content, 0755); err != nil {
		return err
	}

	return nil
}

// RootfsInfo contains information about the rootfs.
type RootfsInfo struct {
	DiskImageExists bool
	DiskImagePath   string
	DiskImageSize   int64
	RootfsDirExists bool
	RootfsDirPath   string
	Architecture    string
}

// Status returns information about the current rootfs.
func (c *Creator) Status() (*RootfsInfo, error) {
	info := &RootfsInfo{
		DiskImagePath: c.cfg.Paths.DiskImage,
		RootfsDirPath: c.cfg.Paths.RootfsDir,
		Architecture:  c.getDebianArch(),
	}

	info.DiskImageExists = c.fs.Exists(c.cfg.Paths.DiskImage)
	info.RootfsDirExists = c.fs.Exists(c.cfg.Paths.RootfsDir)

	if info.DiskImageExists {
		if fi, err := os.Stat(c.cfg.Paths.DiskImage); err == nil {
			info.DiskImageSize = fi.Size()
		}
	}

	return info, nil
}

// Clean removes the rootfs directory and disk image.
func (c *Creator) Clean(ctx context.Context) error {
	// Remove disk image
	if c.fs.Exists(c.cfg.Paths.DiskImage) {
		if err := c.exec.Run(ctx, "rm", "-f", c.cfg.Paths.DiskImage); err != nil {
			return fmt.Errorf("failed to remove disk image: %w", err)
		}
	}

	// Remove rootfs directory (needs sudo due to root-owned files)
	if c.fs.Exists(c.cfg.Paths.RootfsDir) {
		if err := c.exec.Run(ctx, "sudo", "rm", "-rf", c.cfg.Paths.RootfsDir); err != nil {
			return fmt.Errorf("failed to remove rootfs directory: %w", err)
		}
	}

	return nil
}

// Exists returns true if the rootfs disk image exists.
func (c *Creator) Exists() bool {
	return c.fs.Exists(c.cfg.Paths.DiskImage)
}
