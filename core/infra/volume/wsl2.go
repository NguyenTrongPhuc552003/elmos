package volume

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// WSL2Manager implements volume management for WSL2 (Windows Subsystem for Linux).
// Uses a hybrid approach: native Linux ext4 images on WSL filesystem.
type WSL2Manager struct {
	exec executor.Executor
	fs   filesystem.FileSystem
}

// NewWSL2Manager creates a new WSL2 volume manager.
func NewWSL2Manager(exec executor.Executor, fs filesystem.FileSystem) Manager {
	return &WSL2Manager{
		exec: exec,
		fs:   fs,
	}
}

// Create creates a sparse ext4 filesystem image.
// Similar to Linux but optimized for WSL2's filesystem characteristics.
func (m *WSL2Manager) Create(ctx context.Context, name, size, path string) error {
	// Parse size
	sizeBytes, err := parseSize(size)
	if err != nil {
		return fmt.Errorf("invalid size format: %w", err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Create sparse file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create volume file: %w", err)
	}
	defer file.Close()

	// Truncate to desired size (creates sparse file)
	if err := file.Truncate(sizeBytes); err != nil {
		return fmt.Errorf("failed to truncate volume file: %w", err)
	}

	// Format as ext4
	// WSL2 supports full ext4 features
	if err := m.exec.Run(ctx, "mkfs.ext4", "-L", name, "-F", path); err != nil {
		return fmt.Errorf("failed to format volume: %w", err)
	}

	return nil
}

// Mount mounts the loop device at the specified mount point.
// WSL2 has full loop device support.
func (m *WSL2Manager) Mount(ctx context.Context, volumePath, mountPoint string) error {
	// Ensure mount point exists
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}

	// Setup loop device
	loopDev, err := m.setupLoopDevice(ctx, volumePath)
	if err != nil {
		return err
	}

	// Mount with WSL2-friendly options
	// -o rw: read-write mode
	// No special WSL2 mount options needed as it supports full Linux syscalls
	if err := m.exec.Run(ctx, "mount", "-o", "rw", loopDev, mountPoint); err != nil {
		// Cleanup loop device if mount fails
		_ = m.exec.Run(ctx, "losetup", "-d", loopDev)
		return fmt.Errorf("failed to mount loop device: %w", err)
	}

	return nil
}

// Unmount unmounts the loop device.
func (m *WSL2Manager) Unmount(ctx context.Context, mountPoint string, force bool) error {
	// Find loop device
	loopDev, err := m.findLoopDevice(ctx, mountPoint)
	if err != nil {
		return err
	}

	// Unmount
	args := []string{mountPoint}
	if force {
		args = append([]string{"-f"}, args...)
	}
	if err := m.exec.Run(ctx, "umount", args...); err != nil {
		return fmt.Errorf("failed to unmount: %w", err)
	}

	// Detach loop device
	if err := m.exec.Run(ctx, "losetup", "-d", loopDev); err != nil {
		return fmt.Errorf("failed to detach loop device: %w", err)
	}

	return nil
}

// IsMounted checks if the volume is currently mounted.
func (m *WSL2Manager) IsMounted(ctx context.Context, volumePath string) (bool, error) {
	// Get absolute path
	absPath, err := filepath.Abs(volumePath)
	if err != nil {
		return false, err
	}

	// Check /proc/mounts
	out, err := m.exec.Output(ctx, "cat", "/proc/mounts")
	if err != nil {
		return false, err
	}

	// Also check losetup
	loopOut, _ := m.exec.Output(ctx, "losetup", "-a")

	return strings.Contains(string(out), absPath) || strings.Contains(string(loopOut), absPath), nil
}

// Exists checks if the volume file exists.
func (m *WSL2Manager) Exists(volumePath string) bool {
	return m.fs.Exists(volumePath)
}

// setupLoopDevice sets up a loop device for the volume file.
func (m *WSL2Manager) setupLoopDevice(ctx context.Context, volumePath string) (string, error) {
	// Use losetup -f to find and setup free loop device
	out, err := m.exec.Output(ctx, "losetup", "-f", "--show", volumePath)
	if err != nil {
		return "", fmt.Errorf("failed to setup loop device: %w", err)
	}

	loopDev := strings.TrimSpace(string(out))
	if loopDev == "" {
		return "", fmt.Errorf("no loop device returned")
	}

	return loopDev, nil
}

// findLoopDevice finds the loop device for a mounted volume.
func (m *WSL2Manager) findLoopDevice(ctx context.Context, mountPoint string) (string, error) {
	out, err := m.exec.Output(ctx, "mount")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, mountPoint) {
			parts := strings.Fields(line)
			if len(parts) > 0 && strings.HasPrefix(parts[0], "/dev/loop") {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf("loop device not found for mount point: %s", mountPoint)
}
