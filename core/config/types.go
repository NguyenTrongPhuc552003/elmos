// Package config provides configuration management for elmos.
// This file contains all configuration struct definitions.
package config

// Config holds the application configuration.
type Config struct {
	// ConfigFile is the path to the loaded configuration file
	ConfigFile string `yaml:"-"`

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

// ImageConfig holds disk image configuration.
type ImageConfig struct {
	Path       string `mapstructure:"path"`
	VolumeName string `mapstructure:"volume_name"`
	Size       string `mapstructure:"size"`
	MountPoint string `mapstructure:"mount_point"`
}

// BuildConfig holds kernel build configuration.
type BuildConfig struct {
	Arch         string `mapstructure:"arch"`
	Jobs         int    `mapstructure:"jobs"`
	LLVM         bool   `mapstructure:"llvm"`
	CrossCompile string `mapstructure:"cross_compile"`
}

// QEMUConfig holds QEMU configuration.
type QEMUConfig struct {
	Memory  string `mapstructure:"memory"`
	GDBPort int    `mapstructure:"gdb_port"`
	SSHPort int    `mapstructure:"ssh_port"`
	SMP     int    `mapstructure:"smp"`
}

// PathsConfig holds important paths.
type PathsConfig struct {
	ProjectRoot   string `mapstructure:"project_root"`
	KernelDir     string `mapstructure:"kernel_dir"`
	ModulesDir    string `mapstructure:"modules_dir"`
	AppsDir       string `mapstructure:"apps_dir"`
	LibrariesDir  string `mapstructure:"libraries_dir"`
	PatchesDir    string `mapstructure:"patches_dir"`
	RootfsDir     string `mapstructure:"rootfs_dir"`
	DiskImage     string `mapstructure:"disk_image"`
	DebianMirror  string `mapstructure:"debian_mirror"`
	ToolchainsDir string `mapstructure:"toolchains_dir"`
}

// ProfileConfig holds a named configuration profile.
type ProfileConfig struct {
	Arch         string `mapstructure:"arch"`
	Jobs         int    `mapstructure:"jobs"`
	Memory       string `mapstructure:"memory"`
	CrossCompile string `mapstructure:"cross_compile"`
}
