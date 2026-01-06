// Package config provides configuration management for elmos.
// This file contains configuration saving and profile management.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Save saves the configuration to a YAML file.
func (cfg *Config) Save(path string) error {
	v := viper.New()
	v.SetConfigType("yaml")

	// Calculate defaults to compare against
	defaults := &Config{}
	// If project root is set in current config, use it for defaults calculation base
	if cfg.Paths.ProjectRoot != "" {
		defaults.Paths.ProjectRoot = cfg.Paths.ProjectRoot
	}
	applyComputedDefaults(defaults)

	// Create a copy to modify for saving
	saveCfg := *cfg
	savePaths := saveCfg.Paths

	// Only save paths that differ from defaults
	if savePaths.ProjectRoot == defaults.Paths.ProjectRoot {
		savePaths.ProjectRoot = ""
	}
	if savePaths.KernelDir == defaults.Paths.KernelDir {
		savePaths.KernelDir = ""
	}
	if savePaths.ModulesDir == defaults.Paths.ModulesDir {
		savePaths.ModulesDir = ""
	}
	if savePaths.AppsDir == defaults.Paths.AppsDir {
		savePaths.AppsDir = ""
	}
	if savePaths.LibrariesDir == defaults.Paths.LibrariesDir {
		savePaths.LibrariesDir = ""
	}
	if savePaths.PatchesDir == defaults.Paths.PatchesDir {
		savePaths.PatchesDir = ""
	}
	if savePaths.RootfsDir == defaults.Paths.RootfsDir {
		savePaths.RootfsDir = ""
	}
	if savePaths.DiskImage == defaults.Paths.DiskImage {
		savePaths.DiskImage = ""
	}

	// Verify image paths too
	if saveCfg.Image.Path == defaults.Image.Path {
		saveCfg.Image.Path = ""
	}
	if saveCfg.Image.MountPoint == defaults.Image.MountPoint {
		saveCfg.Image.MountPoint = ""
	}

	saveCfg.Paths = savePaths

	// Set all values
	v.Set("image", saveCfg.Image)
	v.Set("build", saveCfg.Build)
	v.Set("qemu", saveCfg.QEMU)
	v.Set("paths", saveCfg.Paths)
	v.Set("profiles", saveCfg.Profiles)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	// Write config
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ApplyProfile applies a named profile to the current configuration.
func (cfg *Config) ApplyProfile(name string) error {
	profile, ok := cfg.Profiles[name]
	if !ok {
		return fmt.Errorf("profile not found: %s", name)
	}

	if profile.Arch != "" {
		cfg.Build.Arch = profile.Arch
	}
	if profile.Jobs > 0 {
		cfg.Build.Jobs = profile.Jobs
	}
	if profile.Memory != "" {
		cfg.QEMU.Memory = profile.Memory
	}
	if profile.CrossCompile != "" {
		cfg.Build.CrossCompile = profile.CrossCompile
	}

	return nil
}

// GetArchConfig returns the architecture configuration for the current build arch.
func (cfg *Config) GetArchConfig() *ArchConfig {
	return GetArchConfig(cfg.Build.Arch)
}
