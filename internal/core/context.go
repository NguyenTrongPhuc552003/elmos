// Package core provides core types, configuration, and context for elmos.
package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Context holds the current build context and state
type Context struct {
	Config    *Config
	Mounted   bool
	KernelDir string
	Verbose   bool
}

// NewContext creates a new build context with the given configuration
func NewContext(cfg *Config) *Context {
	return &Context{
		Config:    cfg,
		KernelDir: cfg.Paths.KernelDir,
	}
}

// IsMounted checks if the kernel volume is currently mounted
func (ctx *Context) IsMounted() bool {
	// Check if mount point exists and is a mount
	info, err := os.Stat(ctx.Config.Image.MountPoint)
	if err != nil || !info.IsDir() {
		return false
	}

	// Use mount command to verify
	out, err := exec.Command("mount").Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), ctx.Config.Image.MountPoint)
}

// EnsureMounted ensures the kernel volume is mounted
func (ctx *Context) EnsureMounted() error {
	if !ctx.IsMounted() {
		return ImageError("kernel volume not mounted", ErrNotMounted)
	}
	return nil
}

// KernelExists checks if the kernel source directory exists
func (ctx *Context) KernelExists() bool {
	gitDir := filepath.Join(ctx.KernelDir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// HasConfig checks if the kernel has been configured (.config exists)
func (ctx *Context) HasConfig() bool {
	configFile := filepath.Join(ctx.KernelDir, ".config")
	_, err := os.Stat(configFile)
	return err == nil
}

// GetKernelImage returns the path to the built kernel image for the current arch
func (ctx *Context) GetKernelImage() string {
	arch := ctx.Config.Build.Arch

	// ARM32 uses zImage, others use Image
	imageName := "Image"
	if arch == "arm" {
		imageName = "zImage"
	}

	return filepath.Join(ctx.KernelDir, "arch", arch, "boot", imageName)
}

// GetVmlinux returns the path to vmlinux (for debugging)
func (ctx *Context) GetVmlinux() string {
	return filepath.Join(ctx.KernelDir, "vmlinux")
}

// HasKernelImage checks if the kernel image has been built
func (ctx *Context) HasKernelImage() bool {
	_, err := os.Stat(ctx.GetKernelImage())
	return err == nil
}

// GetMakeEnv returns environment variables for kernel make commands
func (ctx *Context) GetMakeEnv() []string {
	cfg := ctx.Config
	// filter existing PATH from os.Environ to avoid duplicates if we want to be clean,
	// checking if we should build a fresh map or just append.
	// For simplicity and safety, we'll build a new list where valid keys are mostly preserved.

	var env []string
	originalPath := os.Getenv("PATH")

	// Prepend GNU tools and LLVM to PATH
	newPath := originalPath

	// GNU tools
	gnuSed := GetBrewLibexecBin("gnu-sed")
	if gnuSed != "" {
		newPath = gnuSed + string(os.PathListSeparator) + newPath
	}
	coreutils := GetBrewLibexecBin("coreutils")
	if coreutils != "" {
		newPath = coreutils + string(os.PathListSeparator) + newPath
	}

	// LLVM toolchain (KEY FIX: required for llvm-objdump, llvm-ar, etc)
	llvmBin := GetBrewBin("llvm")
	if llvmBin != "" {
		newPath = llvmBin + string(os.PathListSeparator) + newPath
	}
	lldBin := GetBrewBin("lld")
	if lldBin != "" {
		newPath = lldBin + string(os.PathListSeparator) + newPath
	}

	// e2fsprogs (sbin)
	e2fsBin := GetBrewSbin("e2fsprogs")
	if e2fsBin != "" {
		newPath = e2fsBin + string(os.PathListSeparator) + newPath
	}

	// Reconstruct env, skipping original PATH
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "PATH=") {
			env = append(env, e)
		}
	}
	// Add modified PATH
	env = append(env, "PATH="+newPath)

	// Add our build-specific environment
	env = append(env,
		"ARCH="+cfg.Build.Arch,
		"LLVM=1",
		"CROSS_COMPILE="+cfg.Build.CrossCompile,
	)

	// Add HOSTCFLAGS for macOS compatibility
	hostcflags := ctx.buildHostCFlags()
	if hostcflags != "" {
		env = append(env, "HOSTCFLAGS="+hostcflags)
	}

	return env
}

// buildHostCFlags constructs the HOSTCFLAGS for macOS kernel builds
func (ctx *Context) buildHostCFlags() string {
	var flags []string

	// Custom macOS headers
	if ctx.Config.Paths.LibrariesDir != "" {
		flags = append(flags, "-I"+ctx.Config.Paths.LibrariesDir)
	}

	// libelf include path (from Homebrew)
	libelfInclude := getBrewInclude("libelf")
	if libelfInclude != "" {
		flags = append(flags, "-I"+libelfInclude)
	}

	// macOS compatibility flags
	flags = append(flags,
		"-D_UUID_T",
		"-D__GETHOSTUUID_H",
		"-D_DARWIN_C_SOURCE",
		"-D_FILE_OFFSET_BITS=64",
	)

	return strings.Join(flags, " ")
}

// getBrewInclude returns the include path for a Homebrew package
func getBrewInclude(pkg string) string {
	out, err := exec.Command("brew", "--prefix", pkg).Output()
	if err != nil {
		return ""
	}
	prefix := strings.TrimSpace(string(out))
	return filepath.Join(prefix, "include")
}

// GetBrewBin returns the bin path for a Homebrew package
func GetBrewBin(pkg string) string {
	out, err := exec.Command("brew", "--prefix", pkg).Output()
	if err != nil {
		return ""
	}
	prefix := strings.TrimSpace(string(out))
	return filepath.Join(prefix, "bin")
}

// GetBrewLibexecBin returns the libexec/gnubin path for GNU tools
func GetBrewLibexecBin(pkg string) string {
	out, err := exec.Command("brew", "--prefix", pkg).Output()
	if err != nil {
		return ""
	}
	prefix := strings.TrimSpace(string(out))
	return filepath.Join(prefix, "libexec", "gnubin")
}

// GetBrewSbin returns the sbin path for a Homebrew package
func GetBrewSbin(pkg string) string {
	out, err := exec.Command("brew", "--prefix", pkg).Output()
	if err != nil {
		return ""
	}
	prefix := strings.TrimSpace(string(out))
	return filepath.Join(prefix, "sbin")
}
