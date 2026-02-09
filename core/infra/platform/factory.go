package platform

import (
	"context"
	"fmt"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/homebrew"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/packages"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/volume"
)

// Factory creates platform-specific implementations.
type Factory struct {
	detector *Detector
	exec     executor.Executor
	fs       filesystem.FileSystem
}

// NewFactory creates a new platform factory.
func NewFactory(exec executor.Executor, fs filesystem.FileSystem) *Factory {
	return &Factory{
		detector: NewDetector(),
		exec:     exec,
		fs:       fs,
	}
}

// GetVolumeManager returns the appropriate volume manager for the current platform.
func (f *Factory) GetVolumeManager() (volume.Manager, error) {
	platformType := f.detector.Detect()

	switch platformType {
	case MacOS:
		return volume.NewMacOSManager(f.exec, f.fs), nil
	case Linux:
		return volume.NewLinuxManager(f.exec, f.fs), nil
	case WSL2:
		return volume.NewWSL2Manager(f.exec, f.fs), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platformType)
	}
}

// GetPackageResolver returns the appropriate package resolver for the current platform.
func (f *Factory) GetPackageResolver() (packages.Resolver, error) {
	platformType := f.detector.Detect()

	switch platformType {
	case MacOS:
		// Use existing Homebrew resolver wrapped in the interface
		brew := homebrew.NewResolver(f.exec)
		return packages.NewHomebrewResolver(brew), nil

	case Linux, WSL2:
		// Detect available package manager
		pkgMgr := f.detectPackageManager()
		switch pkgMgr {
		case "apt":
			return packages.NewAptResolver(f.exec), nil
		case "pacman":
			return packages.NewPacmanResolver(f.exec), nil
		default:
			// Default to APT as it's most common
			return packages.NewAptResolver(f.exec), nil
		}

	default:
		return nil, fmt.Errorf("unsupported platform: %s", platformType)
	}
}

// GetPlatformType returns the detected platform type.
func (f *Factory) GetPlatformType() Type {
	return f.detector.Detect()
}

// detectPackageManager detects which package manager is available on Linux systems.
func (f *Factory) detectPackageManager() string {
	ctx := context.Background()

	// Check for apt-get (Debian/Ubuntu)
	if err := f.exec.Run(ctx, "which", "apt-get"); err == nil {
		return "apt"
	}

	// Check for pacman (Arch Linux)
	if err := f.exec.Run(ctx, "which", "pacman"); err == nil {
		return "pacman"
	}

	// Check for dnf (Fedora/RHEL)
	if err := f.exec.Run(ctx, "which", "dnf"); err == nil {
		return "dnf"
	}

	// Check for yum (older RHEL/CentOS)
	if err := f.exec.Run(ctx, "which", "yum"); err == nil {
		return "yum"
	}

	// Default to apt
	return "apt"
}
