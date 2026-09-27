// Package config provides configuration management for elmos.
// This file contains all configuration struct definitions.
package config

// Config holds the application configuration.
type Config struct {
	// ConfigFile is the path to the loaded configuration file
	ConfigFile string `yaml:"-"`
	// ExplicitPaths records workspace paths supplied before defaults were applied.
	ExplicitPaths ExplicitWorkspacePaths `mapstructure:"-" yaml:"-"`

	// Image settings
	Image ImageConfig `mapstructure:"image"`

	// Build settings
	Build BuildConfig `mapstructure:"build"`

	// QEMU settings
	QEMU QEMUConfig `mapstructure:"qemu"`

	// Paths
	Paths PathsConfig `mapstructure:"paths"`

	// Profiles for different configurations
	Profiles map[string]ProfileConfig `mapstructure:"profiles"`
}

// ExplicitWorkspacePaths distinguishes configured paths from computed defaults.
type ExplicitWorkspacePaths struct {
	KernelDir bool
	RootfsDir bool
	DiskImage bool
}

// ImageConfig holds disk image configuration.
type ImageConfig struct {
	Path       string `mapstructure:"path" yaml:"path"`
	VolumeName string `mapstructure:"volume_name" yaml:"volume_name"`
	Size       string `mapstructure:"size" yaml:"size"`
	MountPoint string `mapstructure:"mount_point" yaml:"mount_point"`
}

// BuildConfig holds kernel build configuration.
type BuildConfig struct {
	Arch         string `mapstructure:"arch" yaml:"arch"`
	Jobs         int    `mapstructure:"jobs" yaml:"jobs"`
	LLVM         bool   `mapstructure:"llvm" yaml:"llvm"`
	CrossCompile string `mapstructure:"cross_compile" yaml:"cross_compile"`
}

// QEMUConfig holds QEMU configuration.
type QEMUConfig struct {
	Memory  string `mapstructure:"memory" yaml:"memory"`
	GDBPort int    `mapstructure:"gdb_port" yaml:"gdb_port"`
	SSHPort int    `mapstructure:"ssh_port" yaml:"ssh_port"`
	SMP     int    `mapstructure:"smp" yaml:"smp"`
}

// PathsConfig holds important paths.
type PathsConfig struct {
	ProjectRoot   string `mapstructure:"project_root" yaml:"project_root"`
	KernelDir     string `mapstructure:"kernel_dir" yaml:"kernel_dir"`
	ModulesDir    string `mapstructure:"modules_dir" yaml:"modules_dir"`
	AppsDir       string `mapstructure:"apps_dir" yaml:"apps_dir"`
	LibrariesDir  string `mapstructure:"libraries_dir" yaml:"libraries_dir"`
	PatchesDir    string `mapstructure:"patches_dir" yaml:"patches_dir"`
	RootfsDir     string `mapstructure:"rootfs_dir" yaml:"rootfs_dir"`
	DiskImage     string `mapstructure:"disk_image" yaml:"disk_image"`
	DebianMirror  string `mapstructure:"debian_mirror" yaml:"debian_mirror"`
	ToolchainsDir string `mapstructure:"toolchains_dir" yaml:"toolchains_dir"`
}

// ProfileConfig holds a named configuration profile.
type ProfileConfig struct {
	Arch         string `mapstructure:"arch" yaml:"arch"`
	Jobs         int    `mapstructure:"jobs" yaml:"jobs"`
	Memory       string `mapstructure:"memory" yaml:"memory"`
	CrossCompile string `mapstructure:"cross_compile" yaml:"cross_compile"`
}
