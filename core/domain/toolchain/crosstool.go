package toolchain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Install installs crosstool-ng from the latest git source.
func (m *Manager) Install(ctx context.Context) error {
	paths := m.Paths()

	// Ensure base directory exists
	if err := m.fs.MkdirAll(paths.Base, 0755); err != nil {
		return fmt.Errorf("failed to create toolchains directory: %w", err)
	}

	// Create src directory for downloads
	if err := m.fs.MkdirAll(paths.Src, 0755); err != nil {
		return fmt.Errorf("failed to create src directory: %w", err)
	}

	srcDir := filepath.Join(paths.Base, "crosstool-ng-src")

	// Clone crosstool-ng if not exists
	if !m.fs.IsDir(srcDir) {
		err := m.exec.Run(ctx, "git", "clone",
			"https://github.com/crosstool-ng/crosstool-ng.git",
			srcDir)
		if err != nil {
			return fmt.Errorf("failed to clone crosstool-ng: %w", err)
		}
	}

	// Get install environment with brew paths
	env := m.getInstallEnv()

	// Bootstrap
	if err := m.exec.RunWithEnvInDir(ctx, env, srcDir, "./bootstrap"); err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}

	// Configure with prefix
	configArgs := []string{
		fmt.Sprintf("--prefix=%s", paths.CrosstoolNG),
	}
	if err := m.exec.RunWithEnvInDir(ctx, env, srcDir, "./configure", configArgs...); err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	// Build
	if err := m.exec.RunWithEnvInDir(ctx, env, srcDir, "make"); err != nil {
		return fmt.Errorf("make failed: %w", err)
	}

	// Install
	if err := m.exec.RunWithEnvInDir(ctx, env, srcDir, "make", "install"); err != nil {
		return fmt.Errorf("make install failed: %w", err)
	}

	return nil
}

// getInstallEnv returns environment variables for installing ct-ng.
// Includes brew binutils and bison paths for macOS.
func (m *Manager) getInstallEnv() []string {
	env := os.Environ()

	// Get brew prefix (typically /opt/homebrew on Apple Silicon)
	brewPrefix := os.Getenv("HOMEBREW_PREFIX")
	if brewPrefix == "" {
		brewPrefix = "/opt/homebrew" // Default for Apple Silicon
	}

	// Add brew binutils and bison to PATH (required for objcopy, etc.)
	currentPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s/opt/binutils/bin:%s/opt/bison/bin:%s",
		brewPrefix, brewPrefix, currentPath)

	// Update PATH in environment
	for i, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			env[i] = "PATH=" + newPath
			return env
		}
	}

	// PATH not found, add it
	env = append(env, "PATH="+newPath)
	return env
}

// SelectTarget configures crosstool-ng for the specified target.
// It checks for a custom config first, then falls back to ct-ng samples.
func (m *Manager) SelectTarget(ctx context.Context, target string) error {
	m.printer.Step("Configuring toolchain for %s...", target)

	if !m.IsInstalled() {
		return fmt.Errorf("crosstool-ng not installed, run 'elmos toolchains install'")
	}

	paths := m.Paths()

	// Create x-tools directory with proper permissions
	if err := m.fs.MkdirAll(paths.XTools, 0755); err != nil {
		return fmt.Errorf("failed to create x-tools directory: %w", err)
	}

	// Check if we have a custom config for this target
	customConfig := m.GetCustomConfigPath(target)
	if customConfig != "" {
		m.printer.Info("Using custom configuration: %s", customConfig)

		// Read custom config
		content, err := os.ReadFile(customConfig)
		if err != nil {
			return fmt.Errorf("failed to read custom config: %w", err)
		}

		// Write to .config in toolchains dir
		configPath := filepath.Join(paths.Base, ".config")
		if err := os.WriteFile(configPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write .config: %w", err)
		}
	} else {
		// Use ct-ng sample
		// Run ct-ng <target> in toolchains directory
		if err := m.exec.RunWithEnvInDir(ctx, m.getBuildEnv(paths), paths.Base, m.GetCtNgPath(), target); err != nil {
			return fmt.Errorf("failed to select target %s: %w", target, err)
		}
	}

	// Update .config to use our paths
	configFile := filepath.Join(paths.Base, ".config")
	if m.fs.Exists(configFile) {
		if err := m.patchConfig(configFile, paths); err != nil {
			return fmt.Errorf("failed to patch config: %w", err)
		}
	}

	return nil
}

// patchConfig updates the .config file to use our paths.
func (m *Manager) patchConfig(configFile string, paths ToolchainPaths) error {
	content, err := m.fs.ReadFile(configFile)
	if err != nil {
		return err
	}

	patched := patchConfigContent(string(content), paths)
	return m.fs.WriteFile(configFile, []byte(patched), 0644)
}

// patchConfigContent replaces home directory paths with our paths.
func patchConfigContent(content string, paths ToolchainPaths) string {
	home := os.Getenv("HOME")

	// Replace prefix directory
	content = replaceConfigValue(content, "CT_PREFIX_DIR",
		fmt.Sprintf("%s/${CT_TARGET}", paths.XTools))

	// Replace source directory
	content = replaceConfigValue(content, "CT_LOCAL_TARBALLS_DIR", paths.Src)

	// Replace home paths with our base
	if home != "" {
		content = replaceAll(content, home+"/x-tools", paths.XTools)
		content = replaceAll(content, home+"/src", paths.Src)
	}

	// On macOS: Disable building companion tools that fail with Clang/GCC mixing
	// Use system versions from Homebrew instead
	content = replaceAll(content, "CT_COMP_TOOLS_M4=y", "# CT_COMP_TOOLS_M4 is not set")
	content = replaceAll(content, "CT_COMP_TOOLS_MAKE=y", "# CT_COMP_TOOLS_MAKE is not set")
	content = replaceAll(content, "CT_COMP_TOOLS_AUTOCONF=y", "# CT_COMP_TOOLS_AUTOCONF is not set")
	content = replaceAll(content, "CT_COMP_TOOLS_AUTOMAKE=y", "# CT_COMP_TOOLS_AUTOMAKE is not set")
	content = replaceAll(content, "CT_COMP_TOOLS_LIBTOOL=y", "# CT_COMP_TOOLS_LIBTOOL is not set")

	return content
}

// replaceConfigValue replaces a config value in the content.
func replaceConfigValue(content, key, value string) string {
	// Find and replace CT_KEY="..." pattern
	prefix := key + "=\""
	start := 0
	for {
		idx := indexOf(content[start:], prefix)
		if idx < 0 {
			break
		}
		idx += start
		endIdx := indexOf(content[idx+len(prefix):], "\"")
		if endIdx < 0 {
			break
		}
		endIdx += idx + len(prefix)
		content = content[:idx] + key + "=\"" + value + "\"" + content[endIdx+1:]
		start = idx + len(key) + len(value) + 3
	}
	return content
}

// replaceAll replaces all occurrences of old with new.
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx < 0 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
