// Package rootfs provides rootfs creation functionality for elmos.
// This file contains type definitions for the rootfs package.
package rootfs

// CreateOptions contains options for creating a rootfs.
type CreateOptions struct {
	Size string // Disk image size, e.g., "5G"
}
