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

// LinuxManager implements volume management using loopback devices for native Linux.
type LinuxManager struct {
	exec executor.Executor
	fs   filesystem.FileSystem
}

// NewLinuxManager creates a new Linux volume manager.
func NewLinuxManager(exec executor.Executor, fs filesystem.FileSystem) Manager {
	return &LinuxManager{
		exec: exec,
		fs:   fs,
	}
}

// Create creates a sparse ext4 filesystem image.
func (m *LinuxManager) Create(ctx context.Context, name, size, path string) error {
	// Parse size (e.g., "20G" -> 20 GB)
	sizeBytes, err := parseSize(size)
	if err != nil {
		return fmt.Errorf("invalid size format: %w", err)
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

	// Format as ext4 with case-sensitive filesystem
	// -L: volume label, -F: force even if it looks mounted
	if err := m.exec.Run(ctx, "mkfs.ext4", "-L", name, "-F", path); err != nil {
		return fmt.Errorf("failed to format volume: %w", err)
	}

	return nil
}

// Mount mounts the loop device at the specified mount point.
func (m *LinuxManager) Mount(ctx context.Context, volumePath, mountPoint string) error {
	// Ensure mount point exists
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}

	// Find free loop device
	loopDev, err := m.setupLoopDevice(ctx, volumePath)
	if err != nil {
		return err
	}

	// Mount the loop device
	if err := m.exec.Run(ctx, "mount", loopDev, mountPoint); err != nil {
		// Cleanup loop device if mount fails
		_ = m.exec.Run(ctx, "losetup", "-d", loopDev)
		return fmt.Errorf("failed to mount loop device: %w", err)
	}

	return nil
}

// Unmount unmounts the loop device.
func (m *LinuxManager) Unmount(ctx context.Context, mountPoint string, force bool) error {
	// Find loop device from mount point
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
func (m *LinuxManager) IsMounted(ctx context.Context, volumePath string) (bool, error) {
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

	// Also check losetup output for the backing file
	loopOut, _ := m.exec.Output(ctx, "losetup", "-a")

	return strings.Contains(string(out), absPath) || strings.Contains(string(loopOut), absPath), nil
}

// Exists checks if the volume file exists.
func (m *LinuxManager) Exists(volumePath string) bool {
	return m.fs.Exists(volumePath)
}

// setupLoopDevice sets up a loop device for the volume file.
func (m *LinuxManager) setupLoopDevice(ctx context.Context, volumePath string) (string, error) {
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
func (m *LinuxManager) findLoopDevice(ctx context.Context, mountPoint string) (string, error) {
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

// parseSize parses size strings like "20G", "500M", "1T" to bytes.
func parseSize(size string) (int64, error) {
	size = strings.ToUpper(strings.TrimSpace(size))
	if len(size) < 2 {
		return 0, fmt.Errorf("invalid size format: %s", size)
	}

	// Extract number and unit
	numStr := size[:len(size)-1]
	unit := size[len(size)-1:]

	var multiplier int64
	switch unit {
	case "K":
		multiplier = 1024
	case "M":
		multiplier = 1024 * 1024
	case "G":
		multiplier = 1024 * 1024 * 1024
	case "T":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("invalid size unit: %s (expected K, M, G, or T)", unit)
	}

	// Parse number
	var num int64
	if _, err := fmt.Sscanf(numStr, "%d", &num); err != nil {
		return 0, fmt.Errorf("invalid size number: %s", numStr)
	}

	return num * multiplier, nil
}
