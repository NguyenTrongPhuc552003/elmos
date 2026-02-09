// Package volume provides abstractions for workspace volume management across platforms.
package volume

import (
	"context"
)

// Manager abstracts platform-specific volume/disk image management.
// Implementations handle creating, mounting, and managing workspace volumes.
type Manager interface {
	// Create creates a new volume/disk image with the specified name and size.
	Create(ctx context.Context, name, size, path string) error

	// Mount mounts the volume at the specified mount point.
	Mount(ctx context.Context, volumePath, mountPoint string) error

	// Unmount unmounts the volume.
	Unmount(ctx context.Context, mountPoint string, force bool) error

	// IsMounted checks if a volume is currently mounted.
	IsMounted(ctx context.Context, volumePath string) (bool, error)

	// Exists checks if a volume file exists.
	Exists(volumePath string) bool
}
