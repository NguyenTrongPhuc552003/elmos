// Package types provides shared types for command packages.
package types

import (
	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/doctor"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/emulator"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/patch"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/rootfs"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/packages"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/platform"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/volume"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

// Context provides dependencies for command builders.
// This avoids circular imports by passing only what commands need.
type Context struct {
	Exec             executor.Executor
	FS               filesystem.FileSystem
	Config           *config.Config
	AppContext       *elcontext.Context
	PlatformFactory  *platform.Factory
	VolumeManager    volume.Manager
	PackageResolver  packages.Resolver
	KernelBuilder    *builder.KernelBuilder
	ModuleBuilder    *builder.ModuleBuilder
	AppBuilder       *builder.AppBuilder
	QEMURunner       *emulator.QEMURunner
	HealthChecker    *doctor.HealthChecker
	AutoFixer        *doctor.AutoFixer
	RootfsCreator    *rootfs.Creator
	PatchManager     *patch.Manager
	ToolchainManager *toolchain.Manager
	Printer          *ui.Printer

	// Flags that can be modified
	Verbose    *bool
	ConfigFile *string
}
