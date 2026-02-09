// Package environment provides build environment construction services.
package environment

import (
	"os"
	"runtime"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/packages"
)

// Builder constructs environment variables for kernel builds and toolchains.
type Builder struct {
	config   *config.Config
	pkgRes   packages.Resolver
	platform string
}

// New creates a new environment builder.
func New(cfg *config.Config, pkgRes packages.Resolver) *Builder {
	return &Builder{
		config:   cfg,
		pkgRes:   pkgRes,
		platform: runtime.GOOS,
	}
}

// BuildMakeEnv returns environment variables for kernel make commands.
func (b *Builder) BuildMakeEnv() []string {
	var env []string
	originalPath := os.Getenv("PATH")
	newPath := originalPath

	// Platform-specific PATH modifications
	if b.platform == "darwin" {
		// macOS: Prepend GNU tools from Homebrew
		if gnuSed := b.pkgRes.GetLibexecBin("gnu-sed"); gnuSed != "" {
			newPath = gnuSed + string(os.PathListSeparator) + newPath
		}
		if coreutils := b.pkgRes.GetLibexecBin("coreutils"); coreutils != "" {
			newPath = coreutils + string(os.PathListSeparator) + newPath
		}

		// LLVM toolchain from Homebrew
		if llvmBin := b.pkgRes.GetBin("llvm"); llvmBin != "" {
			newPath = llvmBin + string(os.PathListSeparator) + newPath
		}
		if lldBin := b.pkgRes.GetBin("lld"); lldBin != "" {
			newPath = lldBin + string(os.PathListSeparator) + newPath
		}

		// e2fsprogs (sbin)
		if e2fsBin := b.pkgRes.GetSbin("e2fsprogs"); e2fsBin != "" {
			newPath = e2fsBin + string(os.PathListSeparator) + newPath
		}
	} else {
		// Linux/WSL2: System packages are already in standard locations
		// Ensure /usr/bin and /usr/sbin are in PATH
		standardPaths := []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin"}
		for _, p := range standardPaths {
			if !strings.Contains(newPath, p) {
				newPath = p + string(os.PathListSeparator) + newPath
			}
		}
	}

	// Reconstruct env, skipping original PATH
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "PATH=") {
			env = append(env, e)
		}
	}
	env = append(env, "PATH="+newPath)

	// Add build-specific environment
	env = append(env,
		"ARCH="+b.config.Build.Arch,
		"LLVM=1",
		"CROSS_COMPILE="+b.config.Build.CrossCompile,
	)

	// Add platform-specific HOSTCFLAGS
	hostcflags := b.buildHostCFlags()
	if hostcflags != "" {
		env = append(env, "HOSTCFLAGS="+hostcflags)
	}

	return env
}

// buildHostCFlags constructs the HOSTCFLAGS for kernel builds.
func (b *Builder) buildHostCFlags() string {
	var flags []string

	if b.platform == "darwin" {
		// macOS-specific flags

		// Custom macOS headers
		if b.config.Paths.LibrariesDir != "" {
			flags = append(flags, "-I"+b.config.Paths.LibrariesDir)
		}

		// libelf include path (from Homebrew)
		if libelfInclude := b.pkgRes.GetInclude("libelf"); libelfInclude != "" {
			flags = append(flags, "-I"+libelfInclude)
		}

		// macOS compatibility flags
		flags = append(flags,
			"-D_UUID_T",
			"-D__GETHOSTUUID_H",
			"-D_DARWIN_C_SOURCE",
			"-D_FILE_OFFSET_BITS=64",
		)
	} else {
		// Linux/WSL2: Minimal flags, system headers are already accessible
		// Only add large file support
		flags = append(flags, "-D_FILE_OFFSET_BITS=64")
	}

	return strings.Join(flags, " ")
}
