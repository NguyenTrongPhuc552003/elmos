// Package config provides configuration management for elmos.
package config

// Default values for configuration.
const (
	// DefaultImageSize is the default sparse image size.
	DefaultImageSize = "20G"
	// DefaultVolumeName is the default volume name for the disk image.
	DefaultVolumeName = "kernel-dev"
	// DefaultArch is the default target architecture.
	DefaultArch = "arm64"
	// DefaultCrossPrefix is the default cross-compile prefix.
	DefaultCrossPrefix = "llvm-"
	// DefaultMemory is the default QEMU memory allocation.
	DefaultMemory = "2G"
	// DefaultGDBPort is the default GDB debugging port.
	DefaultGDBPort = 1234
	// DefaultSSHPort is the default SSH forwarding port.
	DefaultSSHPort = 2222
	// DefaultDebianMirror is the default Debian package mirror.
	DefaultDebianMirror = "http://deb.debian.org/debian"
	// DefaultGlibcVersion is the glibc version used for downloading elf.h.
	DefaultGlibcVersion = "2.42"
)

// RequiredPackage represents a Homebrew package dependency.
type RequiredPackage struct {
	Name        string
	Description string
	Required    bool
}

// RequiredPackages lists all Homebrew dependencies for elmos.
var RequiredPackages = []RequiredPackage{
	// Core build tools
	{"llvm", "LLVM/Clang toolchain", true},
	{"lld", "LLVM linker", true},
	{"gnu-sed", "GNU sed (kernel requires it)", true},
	{"make", "GNU make 4.0+", true},
	{"libelf", "ELF library", true},
	{"git", "Git version control", true},
	{"qemu", "QEMU emulator", true},
	{"fakeroot", "Fake root for packaging", true},
	{"e2fsprogs", "ext4 filesystem tools", true},
	{"wget", "File downloader", false},
	{"coreutils", "GNU core utilities", true},
	{"go", "Go programming language", true},
	{"go-task", "Go task runner", true},
	// Crosstool-ng dependencies (optional, for building custom toolchains)
	{"binutils", "GNU binary utilities (objcopy)", false},
	{"gcc", "GNU Compiler Collection (for ct-ng builds)", false},
	{"gmp", "GNU Multiple Precision library", false},
	{"mpfr", "GNU MPFR library", false},
	{"libmpc", "GNU MPC library", false},
	{"isl", "Integer Set Library", false},
	{"texinfo", "GNU documentation system", false},
	{"bison", "Parser generator", false},
	{"gawk", "GNU AWK", false},
	{"autoconf", "Autoconf for ct-ng bootstrap", false},
	{"automake", "Automake for ct-ng bootstrap", false},
	{"libtool", "GNU Libtool", false},
	{"ncurses", "Terminal UI library (menuconfig)", false},
	{"xz", "XZ compression", false},
}

// RequiredTaps lists required Homebrew taps.
var RequiredTaps = []string{
	"messense/macos-cross-toolchains",
}

// RequiredHeaders lists header files that should exist in libraries/.
var RequiredHeaders = []string{
	"elf.h",
	"byteswap.h",
}

// ValidBuildTargets lists valid kernel build targets.
var ValidBuildTargets = map[string]bool{
	"Image":           true,
	"zImage":          true, // ARM32
	"dtbs":            true,
	"modules":         true,
	"modules_prepare": true,
	"all":             true,
	"vmlinux":         true,
}

// KernelConfigTypes lists valid kernel configuration types.
var KernelConfigTypes = []string{
	"defconfig",
	"tinyconfig",
	"kvm_guest.config",
	"menuconfig",
	"xconfig",
	"nconfig",
	"oldconfig",
	"olddefconfig",
	"allnoconfig",
	"allyesconfig",
	"allmodconfig",
	"localmodconfig",
	"localyesconfig",
}
