package toolchain

import (
	"fmt"
	"os"
	"path/filepath"
)

// getBuildEnv returns environment variables for building toolchains.
// Implements brianredbeard's solution from github.com/crosstool-ng/crosstool-ng/issues/2378.
func (m *Manager) getBuildEnv(paths ToolchainPaths) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("CT_PREFIX=%s", paths.XTools))

	// Export workspace name for crosstool-NG config variable substitution
	// This allows toolchain configs to use ${ELMOS_WORKSPACE} in paths
	if m.cfg != nil && m.cfg.Image.VolumeName != "" {
		env = append(env, fmt.Sprintf("ELMOS_WORKSPACE=%s", m.cfg.Image.VolumeName))
	}

	brewPrefix := getBrewPrefix()
	localBin := m.ensureLocalBin()
	m.ensureGCCSymlinks(localBin, brewPrefix)

	env = m.addLibraryFlags(env, brewPrefix)
	env = m.addGmake(env, brewPrefix)
	env = m.buildPathEnv(env, localBin, brewPrefix)

	return env
}

// getBrewPrefix returns the Homebrew prefix.
func getBrewPrefix() string {
	if prefix := os.Getenv("HOMEBREW_PREFIX"); prefix != "" {
		return prefix
	}
	return "/opt/homebrew"
}

// ensureLocalBin creates and returns the ~/.local/bin directory.
func (m *Manager) ensureLocalBin() string {
	home := os.Getenv("HOME")
	localBin := filepath.Join(home, ".local", "bin")
	_ = m.fs.MkdirAll(localBin, 0755)
	return localBin
}

// ensureGCCSymlinks creates gcc/g++ symlinks in localBin if needed.
func (m *Manager) ensureGCCSymlinks(localBin, brewPrefix string) {
	gccLink := filepath.Join(localBin, "gcc")
	if m.fs.Exists(gccLink) {
		return
	}
	for _, ver := range []string{"14", "15", "13", "12"} {
		gccPath := fmt.Sprintf("%s/bin/gcc-%s", brewPrefix, ver)
		if m.fs.Exists(gccPath) {
			_ = os.Symlink(gccPath, gccLink)
			gxxPath := fmt.Sprintf("%s/bin/g++-%s", brewPrefix, ver)
			gxxLink := filepath.Join(localBin, "g++")
			if !m.fs.Exists(gxxLink) {
				_ = os.Symlink(gxxPath, gxxLink)
			}
			break
		}
	}
}

// addLibraryFlags adds LDFLAGS and CPPFLAGS for Homebrew libraries.
func (m *Manager) addLibraryFlags(env []string, brewPrefix string) []string {
	ldPkgs := []string{"bison", "flex", "ncurses", "zlib"}
	cppPkgs := []string{"binutils", "flex", "ncurses", "zlib"}

	ldflags := m.buildFlagString(brewPrefix, ldPkgs, "lib", "-L")
	if ldflags != "" {
		env = appendOrUpdateEnv(env, "LDFLAGS", ldflags)
	}

	cppflags := m.buildFlagString(brewPrefix, cppPkgs, "include", "-I")
	if cppflags != "" {
		env = appendOrUpdateEnv(env, "CPPFLAGS", cppflags)
	}

	pkgConfigPath := fmt.Sprintf("%s/share/pkgconfig", brewPrefix)
	env = appendOrUpdateEnv(env, "PKG_CONFIG_PATH", pkgConfigPath)

	return env
}

// buildFlagString builds a flag string for the given packages.
func (m *Manager) buildFlagString(brewPrefix string, pkgs []string, subdir, flag string) string {
	var result string
	for _, pkg := range pkgs {
		pkgPath := fmt.Sprintf("%s/opt/%s/%s", brewPrefix, pkg, subdir)
		if m.fs.IsDir(pkgPath) {
			result += fmt.Sprintf(" %s%s", flag, pkgPath)
		}
	}
	return result
}

// addGmake sets CT_MAKE_FOR_BUILD to gmake if available.
func (m *Manager) addGmake(env []string, brewPrefix string) []string {
	gmakePath := fmt.Sprintf("%s/bin/gmake", brewPrefix)
	if m.fs.Exists(gmakePath) {
		env = append(env, "CT_MAKE_FOR_BUILD="+gmakePath)
	}
	return env
}

// buildPathEnv builds the PATH environment variable.
func (m *Manager) buildPathEnv(env []string, localBin, brewPrefix string) []string {
	currentPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s:%s/opt/bison/bin:%s/opt/flex/bin:%s/opt/coreutils/libexec/gnubin:%s/opt/gnu-tar/libexec/gnubin:%s/opt/gnu-sed/libexec/gnubin:%s/opt/libtool/libexec/gnubin:%s/opt/grep/libexec/gnubin:%s/opt/gawk/libexec/gnubin:%s",
		localBin, brewPrefix, brewPrefix, brewPrefix, brewPrefix, brewPrefix, brewPrefix, brewPrefix, brewPrefix, currentPath)

	for i, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			env[i] = "PATH=" + newPath
			return env
		}
	}
	return append(env, "PATH="+newPath)
}

// appendOrUpdateEnv appends to an existing env var or creates it.
func appendOrUpdateEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if len(e) >= len(prefix) && e[:len(prefix)] == prefix {
			env[i] = e + value
			return env
		}
	}
	return append(env, key+"="+value)
}
