// Package patch provides kernel patch management for elmos.
// This file contains type definitions for the patch package.
package patch

// PatchInfo contains information about a patch file.
type PatchInfo struct {
	Name    string // Name of the patch file
	Path    string // Full path to the patch file
	Version string // Kernel version this patch applies to
}
