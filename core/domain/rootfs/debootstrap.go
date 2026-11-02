// Package rootfs provides rootfs creation functionality for elmos.
// This file contains type definitions for the rootfs package.
package rootfs

// Root file system engines.
const (
	// EngineOCI builds the rootfs from a container image (default).
	EngineOCI = "oci"
	// EngineDebootstrap builds the rootfs with debootstrap (needs sudo).
	EngineDebootstrap = "debootstrap"
	// DefaultOCIImage is the image the OCI engine uses by default. Alpine's
	// busybox provides insmod and ip, which the guest /init and the module
	// sync script need; debian:stable-slim ships neither (audit E15).
	DefaultOCIImage = "alpine:3"
)

// CreateOptions contains options for creating a rootfs.
type CreateOptions struct {
	Size   string // Disk image size, e.g., "5G"
	Engine string // EngineOCI (default) or EngineDebootstrap
	Image  string // OCI image reference for EngineOCI
}
