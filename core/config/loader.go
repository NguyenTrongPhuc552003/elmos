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
		v.AddConfigPath(filepath.Join(".", "build"))                          // Build directory
		v.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "elmos")) // User config
		v.AddConfigPath("/etc/elmos")                                         // System config

		// Auto-create elmos.yaml from embedded template in build/ if it doesn't exist
		cwd, _ := os.Getwd()
		buildDir := filepath.Join(cwd, "build")
		configFile := filepath.Join(buildDir, "elmos.yaml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			// Ensure build directory exists
			_ = os.MkdirAll(buildDir, 0755)
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
	applyProjectRoot(cfg)
	applyImageDefaults(cfg)
	applyPathDefaults(cfg)
}

// applyProjectRoot sets the project root if not already set.
func applyProjectRoot(cfg *Config) {
	if cfg.Paths.ProjectRoot == "" {
		if cwd, err := os.Getwd(); err == nil {
			cfg.Paths.ProjectRoot = cwd
		}
	}
}

// applyImageDefaults sets image-related defaults.
func applyImageDefaults(cfg *Config) {
	root := cfg.Paths.ProjectRoot
	if cfg.Image.Path == "" {
		cfg.Image.Path = filepath.Join(root, "data", fmt.Sprintf("%s.sparseimage", cfg.Image.VolumeName))
	}
	if cfg.Image.MountPoint == "" {
		cfg.Image.MountPoint = filepath.Join("/Volumes", cfg.Image.VolumeName)
	}
}

// applyPathDefaults sets path-related defaults.
func applyPathDefaults(cfg *Config) {
	root := cfg.Paths.ProjectRoot
	mount := cfg.Image.MountPoint

	setIfEmpty(&cfg.Paths.KernelDir, filepath.Join(mount, "linux"))
	setIfEmpty(&cfg.Paths.ModulesDir, filepath.Join(root, "modules"))
	setIfEmpty(&cfg.Paths.AppsDir, filepath.Join(root, "apps"))
	setIfEmpty(&cfg.Paths.LibrariesDir, filepath.Join(root, "libraries"))
	setIfEmpty(&cfg.Paths.PatchesDir, filepath.Join(root, "patches"))
	setIfEmpty(&cfg.Paths.RootfsDir, filepath.Join(mount, "rootfs"))
	setIfEmpty(&cfg.Paths.DiskImage, filepath.Join(mount, "disk.img"))
	setIfEmpty(&cfg.Paths.ToolchainsDir, filepath.Join(mount, "toolchains"))
}

// setIfEmpty sets the target to value if target is empty.
func setIfEmpty(target *string, value string) {
	if *target == "" {
		*target = value
	}
}
