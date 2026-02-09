package volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// MacOSManager implements volume management using hdiutil for macOS sparse images.
type MacOSManager struct {
	exec executor.Executor
	fs   filesystem.FileSystem
}

// NewMacOSManager creates a new macOS volume manager.
func NewMacOSManager(exec executor.Executor, fs filesystem.FileSystem) Manager {
	return &MacOSManager{
		exec: exec,
		fs:   fs,
	}
}

// Create creates a sparse APFS disk image.
func (m *MacOSManager) Create(ctx context.Context, name, size, path string) error {
	return m.exec.Run(ctx, "hdiutil", "create",
		"-size", size,
		"-fs", "Case-sensitive APFS",
		"-volname", name,
		"-type", "SPARSE",
		path,
	)
}

// Mount mounts the disk image at the specified mount point.
func (m *MacOSManager) Mount(ctx context.Context, volumePath, mountPoint string) error {
	return m.exec.Run(ctx, "hdiutil", "attach",
		"-mountpoint", mountPoint,
		volumePath,
	)
}

// Unmount unmounts the disk image.
func (m *MacOSManager) Unmount(ctx context.Context, mountPoint string, force bool) error {
	// Find device from mount point
	device, err := m.findDevice(ctx, mountPoint)
	if err != nil {
		return err
	}

	args := []string{"detach", device}
	if force {
		args = append(args, "-force")
	}

	return m.exec.Run(ctx, "hdiutil", args...)
}

// IsMounted checks if the volume is currently mounted.
func (m *MacOSManager) IsMounted(ctx context.Context, volumePath string) (bool, error) {
	out, err := m.exec.Output(ctx, "hdiutil", "info")
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), volumePath), nil
}

// Exists checks if the volume file exists.
func (m *MacOSManager) Exists(volumePath string) bool {
	return m.fs.Exists(volumePath)
}

// findDevice finds the device path for a mounted volume.
func (m *MacOSManager) findDevice(ctx context.Context, mountPoint string) (string, error) {
	out, err := m.exec.Output(ctx, "mount")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, mountPoint) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf("device not found for mount point: %s", mountPoint)
}
