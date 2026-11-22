package doctor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func checker(t *testing.T, goos string, onPath ...string) (*HealthChecker, *elconfig.Config) {
	t.Helper()
	root := t.TempDir()
	cfg := &elconfig.Config{}
	cfg.Paths.LibrariesDir = filepath.Join(root, "libraries")
	cfg.Paths.ToolchainsDir = filepath.Join(root, "toolchains")
	mock := executor.NewMockExecutor()
	for _, c := range onPath {
		mock.LookPathResponses[c] = "/usr/bin/" + c
	}
	fs := filesystem.NewOSFileSystem()
	h := NewHealthChecker(mock, fs, cfg, toolchain.NewManager(mock, fs, cfg, nil))
	h.goos = goos
	return h, cfg
}

func find(results []CheckResult, name string) *CheckResult {
	for i := range results {
		if results[i].Name == name {
			return &results[i]
		}
	}
	return nil
}

func TestHomebrewIsRequiredOnMacOSOnly(t *testing.T) {
	mac, _ := checker(t, "darwin")
	if r := mac.CheckHomebrew(context.Background()); r.Passed || !r.Required {
		t.Fatalf("macOS without brew: %+v", r)
	}
	linux, _ := checker(t, "linux")
	if r := linux.CheckHomebrew(context.Background()); r.Required || !strings.Contains(r.Message, "not needed on linux") {
		t.Fatalf("Linux without brew: %+v", r)
	}
}

func TestLinuxChecksHostToolsOnPath(t *testing.T) {
	h, _ := checker(t, "linux", "make", "git")
	res := h.CheckPackages(context.Background())
	if r := find(res, "Host tools: make"); r == nil || !r.Passed {
		t.Fatalf("make: %+v", res)
	}
	q := find(res, "Host tools: qemu-system-aarch64")
	if q == nil || q.Passed || !q.Required || !strings.Contains(q.Message, "package manager") {
		t.Fatalf("qemu: %+v", q)
	}
	if c := find(res, "Host tools: clang"); c == nil || c.Required {
		t.Fatalf("clang is optional: %+v", c)
	}
	for _, r := range res {
		if strings.HasPrefix(r.Name, "Homebrew") {
			t.Fatalf("Homebrew package check on Linux: %+v", r)
		}
	}
}

func TestIssuesCountedPerHost(t *testing.T) {
	// a Linux host with its tools and no Homebrew or header shims: no issue
	linux, _ := checker(t, "linux", "qemu-system-aarch64", "make", "git")
	if _, issues := linux.CheckAll(context.Background()); issues != 0 {
		t.Fatalf("Linux with its tools: %d issues", issues)
	}
	// the same on macOS: Homebrew, its package list and the shims are required
	mac, _ := checker(t, "darwin", "qemu-system-aarch64", "make", "git")
	if _, issues := mac.CheckAll(context.Background()); issues == 0 {
		t.Fatal("macOS without Homebrew or header shims: no issue")
	}
}

func TestFixHeadersThenCheckHeadersPasses(t *testing.T) {
	h, cfg := checker(t, "darwin")
	if r := h.CheckHeaders(context.Background()); len(r) != 1 || r[0].Passed || !r[0].Required {
		t.Fatalf("no libraries directory: %+v", r)
	}
	if !h.IsElfHMissing() {
		t.Fatal("elf.h reported present before the fix")
	}
	written, err := NewAutoFixer(filesystem.NewOSFileSystem(), cfg).FixHeaders()
	if err != nil || len(written) == 0 {
		t.Fatalf("FixHeaders: %v, %v", written, err)
	}
	for _, r := range h.CheckHeaders(context.Background()) {
		if !r.Passed {
			t.Fatalf("after FixHeaders: %+v", r)
		}
	}
	// a second fix writes nothing and keeps an edited header
	edited := filepath.Join(cfg.Paths.LibrariesDir, "elf.h")
	if err := os.WriteFile(edited, []byte("/* mine */\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if again, _ := NewAutoFixer(filesystem.NewOSFileSystem(), cfg).FixHeaders(); len(again) != 0 {
		t.Fatalf("second FixHeaders wrote %v", again)
	}
	if b, _ := os.ReadFile(edited); string(b) != "/* mine */\n" {
		t.Fatal("FixHeaders overwrote an existing header")
	}
}

func TestCrossGCCFoundOnPath(t *testing.T) {
	h, _ := checker(t, "linux", "aarch64-unknown-linux-gnu-gcc")
	r := find(h.CheckCrossGCC(context.Background()), "Cross Compilers: arm64 (aarch64-unknown-linux-gnu-gcc)")
	if r == nil || !r.Passed || !strings.Contains(r.Message, "Found on PATH") {
		t.Fatalf("%+v", r)
	}
}
