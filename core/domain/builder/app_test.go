package builder

import (
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
)

// Audit E12: with only LLVM installed, "llvm-" + "gcc" ran the host
// llvm-gcc, which links for macOS and failed on crt0.o. The build must stop
// with an actionable message instead.
func TestAppCompilerNeedsALinuxToolchain(t *testing.T) {
	cfg := &elconfig.Config{}
	cfg.Build.Arch = "arm64"
	gcc := elconfig.GetArchConfig("arm64").GCCBinary
	rec := executor.NewRecorder()
	rec.NotOnPath = map[string]bool{gcc: true}
	a := &AppBuilder{exec: rec, cfg: cfg}

	if c, err := a.getCrossCompiler("llvm-"); err == nil || !strings.Contains(err.Error(), "toolchains build") {
		t.Fatalf("compiler = %q, err = %v; want a toolchain error", c, err)
	}
	if c, err := a.getCrossCompiler("aarch64-unknown-linux-gnu-"); err != nil || c != "aarch64-unknown-linux-gnu-gcc" {
		t.Fatalf("triple prefix: %q %v", c, err)
	}
	rec.NotOnPath = nil
	if c, err := a.getCrossCompiler("llvm-"); err != nil || c != "/usr/bin/"+gcc {
		t.Fatalf("gcc on PATH: %q %v", c, err)
	}
}
