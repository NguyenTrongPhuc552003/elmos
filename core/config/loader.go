// Package config provides configuration management for elmos.
// This file contains configuration loading and default handling.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/NguyenTrongPhuc552003/elmos/assets"
	"github.com/spf13/viper"
)

var (
	configInstance *Config
)

// Load loads configuration from files and environment.
// This is the preferred method for new code using dependency injection.
func Load(configPath string) (*Config, error) {
	var loadErr error
	var cfg *Config

	v := viper.New()

	if configPath != "" {
		// Use specific config file
		v.SetConfigFile(configPath)
	} else {
		// Set config name and search paths
		v.SetConfigName("elmos")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")                                                  // Current directory
		v.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "elmos")) // User config
		v.AddConfigPath("/etc/elmos")                                         // System config

		// Auto-create elmos.yaml from embedded template if it doesn't exist
		cwd, _ := os.Getwd()
		configFile := filepath.Join(cwd, "elmos.yaml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			// Use embedded template
			if tmplData, err := assets.GetConfigTemplate(); err == nil {
				_ = os.WriteFile(configFile, tmplData, 0644)
			}
		}
	}

	// Environment variables
	v.SetEnvPrefix("ELMOS")
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		// If specific file provided, error out. Otherwise ignore if not found.
		if configPath != "" {
			loadErr = fmt.Errorf("failed to read config file %s: %w", configPath, err)
		} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// If it's not a "not found" error, report it
			loadErr = fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal into struct
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply computed defaults
	applyComputedDefaults(cfg)

	// Save which file was used
	cfg.ConfigFile = v.ConfigFileUsed()

	// Update singleton
	configInstance = cfg
	return cfg, loadErr
}

// Get returns the current configuration.
// If Load() hasn't been called, it will be called automatically.
func Get() *Config {
	if configInstance == nil {
		cfg, _ := Load("")
		return cfg
	}
	return configInstance
}

// Reset clears the cached configuration (for testing).
func Reset() {
	configInstance = nil
}

// setDefaults sets default values for configuration.
func setDefaults(v *viper.Viper) {
	// Image defaults
	v.SetDefault("image.volume_name", DefaultVolumeName)
	v.SetDefault("image.size", DefaultImageSize)

	// Build defaults
	v.SetDefault("build.arch", DefaultArch)
	v.SetDefault("build.jobs", runtime.NumCPU())
	v.SetDefault("build.llvm", true)
	v.SetDefault("build.cross_compile", DefaultCrossPrefix)

	// QEMU defaults
	v.SetDefault("qemu.memory", DefaultMemory)
	v.SetDefault("qemu.gdb_port", DefaultGDBPort)
	v.SetDefault("qemu.ssh_port", DefaultSSHPort)
	v.SetDefault("qemu.smp", runtime.NumCPU())

	// Paths defaults
	v.SetDefault("paths.debian_mirror", DefaultDebianMirror)
}

// applyComputedDefaults fills in paths based on project root.
func applyComputedDefaults(cfg *Config) {
	// Find project root (where go.mod or elmos.yaml exists)
	if cfg.Paths.ProjectRoot == "" {
		if cwd, err := os.Getwd(); err == nil {
			cfg.Paths.ProjectRoot = cwd
		}
	}

	root := cfg.Paths.ProjectRoot

	// Image path
	if cfg.Image.Path == "" {
		cfg.Image.Path = filepath.Join(root, "img.sparseimage")
	}

	// Mount point
	if cfg.Image.MountPoint == "" {
		cfg.Image.MountPoint = filepath.Join("/Volumes", cfg.Image.VolumeName)
	}

	// Kernel directory (inside mount)
	if cfg.Paths.KernelDir == "" {
		cfg.Paths.KernelDir = filepath.Join(cfg.Image.MountPoint, "linux")
	}

	// Modules directory (project root)
	if cfg.Paths.ModulesDir == "" {
		cfg.Paths.ModulesDir = filepath.Join(root, "modules")
	}

	// Apps directory (project root)
	if cfg.Paths.AppsDir == "" {
		cfg.Paths.AppsDir = filepath.Join(root, "apps")
	}

	// Libraries directory (project root)
	if cfg.Paths.LibrariesDir == "" {
		cfg.Paths.LibrariesDir = filepath.Join(root, "libraries")
	}

	// Patches directory (project root)
	if cfg.Paths.PatchesDir == "" {
		cfg.Paths.PatchesDir = filepath.Join(root, "patches")
	}

	// Rootfs directory (inside mount)
	if cfg.Paths.RootfsDir == "" {
		cfg.Paths.RootfsDir = filepath.Join(cfg.Image.MountPoint, "rootfs")
	}

	// Disk image (inside mount)
	if cfg.Paths.DiskImage == "" {
		cfg.Paths.DiskImage = filepath.Join(cfg.Image.MountPoint, "disk.img")
	}

	// Toolchains directory (inside mount for case-sensitivity)
	if cfg.Paths.ToolchainsDir == "" {
		cfg.Paths.ToolchainsDir = filepath.Join(cfg.Image.MountPoint, "toolchains")
	}
}
