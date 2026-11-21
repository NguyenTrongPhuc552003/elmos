package emulator

import (
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

func runner(arch string, gdbPort int) (*QEMURunner, *elconfig.ArchConfig) {
	cfg := &elconfig.Config{}
	cfg.Build.Arch = arch
	cfg.QEMU.Memory = "2G"
	cfg.QEMU.SMP = 2
	cfg.QEMU.GDBPort = gdbPort
	cfg.QEMU.SSHPort = 2222
	cfg.Paths.DiskImage = "/w/rootfs.img"
	cfg.Paths.ModulesDir = "/w/modules"
	return &QEMURunner{cfg: cfg}, elconfig.GetArchConfig(arch)
}

// value returns the argument after the first occurrence of flag, or "".
func value(args []string, flag string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == flag {
			return args[i+1]
		}
	}
	return ""
}

func has(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func TestArgsPerArchitecture(t *testing.T) {
	for _, tc := range []struct{ arch, machine, cpu, console string }{
		{"arm64", "virt", "cortex-a72", "ttyAMA0"},
		{"arm", "virt,highmem=off", "cortex-a15", "ttyAMA0"},
		{"riscv", "virt", "rv64", "ttyS0"},
	} {
		q, a := runner(tc.arch, 1234)
		args := q.buildArgs(a, "/w/Image", RunOptions{})
		if value(args, "-machine") != tc.machine || value(args, "-cpu") != tc.cpu || value(args, "-kernel") != "/w/Image" {
			t.Errorf("%s: %v", tc.arch, args)
		}
		if got := value(args, "-append"); !strings.Contains(got, "console="+tc.console) || !strings.Contains(got, "root=/dev/vda") {
			t.Errorf("%s: -append %q", tc.arch, got)
		}
		if value(args, "-m") != "2G" || value(args, "-smp") != "2" || !strings.Contains(value(args, "-drive"), "file=/w/rootfs.img") {
			t.Errorf("%s: memory, CPUs or disk wrong: %v", tc.arch, args)
		}
		if has(args, "-S") {
			t.Errorf("%s: a normal run waits for a debugger: %v", tc.arch, args)
		}
	}
}

func TestOnlyRISCVGetsTheDefaultBIOS(t *testing.T) {
	for arch, want := range map[string]bool{"arm64": false, "arm": false, "riscv": true} {
		q, a := runner(arch, 1234)
		if got := value(q.buildArgs(a, "/w/Image", RunOptions{}), "-bios") == "default"; got != want {
			t.Errorf("%s: -bios default = %v, want %v", arch, got, want)
		}
	}
}

// The debugger elmos starts (ConnectGDB) connects to QEMU.GDBPort, so the
// stub must listen there, not on QEMU's -s default of 1234.
func TestDebugListensOnTheConfiguredGDBPort(t *testing.T) {
	q, a := runner("arm64", 1235)
	args := q.buildArgs(a, "/w/Image", RunOptions{Debug: true})
	if value(args, "-gdb") != "tcp::1235" || !has(args, "-S") {
		t.Fatalf("debug args %v, want -gdb tcp::1235 and -S", args)
	}
	if has(args, "-s") {
		t.Fatalf("-s opens port 1234 as well: %v", args)
	}
}

// cocoa exists only in macOS builds of QEMU; elsewhere QEMU picks its own
// default display.
func TestGraphicalDisplayPerHost(t *testing.T) {
	for goos, want := range map[string]string{"darwin": "cocoa", "linux": "", "windows": ""} {
		if got := value(displayArgs(goos), "-display"); got != want {
			t.Errorf("%s: -display %q, want %q", goos, got, want)
		}
	}
}

func TestInstallHintPerHost(t *testing.T) {
	for goos, want := range map[string]string{"darwin": "brew install qemu", "linux": "package manager", "windows": "qemu.org"} {
		if !strings.Contains(qemuInstallHint(goos), want) {
			t.Errorf("%s: hint %q lacks %q", goos, qemuInstallHint(goos), want)
		}
	}
}
