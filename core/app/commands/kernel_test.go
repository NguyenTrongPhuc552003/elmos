package commands

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "codeberg.org/NguyenTrongPhuc552003/elmos/core/context"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/domain/patch"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/ui"
)

// gitAndPatchCalls drops the volume checks (hdiutil) that RunEWithContext makes.
func gitAndPatchCalls(rec *executor.Recorder) []string {
	var out []string
	for _, c := range rec.Calls {
		if strings.HasPrefix(c, "git ") || strings.HasPrefix(c, "patch ") {
			out = append(out, c)
		}
	}
	return out
}

func newKernelTestContext(t *testing.T) (*Context, *executor.Recorder) {
	t.Helper()
	cfg := &config.Config{}
	cfg.Paths.KernelDir = filepath.Join(t.TempDir(), "linux")
	cfg.Build.Arch = "arm64"
	rec := executor.NewRecorder()
	fs := filesystem.NewOSFileSystem()
	return &Context{Config: cfg, FS: fs, Exec: rec, Printer: ui.NewPrinter(),
		AppContext:   elcontext.New(cfg, rec, fs),
		PatchManager: patch.NewManager(rec, fs, cfg)}, rec
}

// Audit E10: the clone was a full, unpinned mainline history (12 minutes,
// several GB) and the shipped macOS patches were never applied.
func TestKernelCloneIsPinnedShallowAndPatched(t *testing.T) {
	ctx, rec := newKernelTestContext(t)
	cmd := buildKernelCloneCmd(ctx)
	cmd.SetContext(context.Background())
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	calls := gitAndPatchCalls(rec)
	if len(calls) == 0 {
		t.Fatal("no command ran")
	}
	want := "git clone --depth 1 --branch " + config.DefaultKernelRef + " " +
		config.DefaultKernelURL + " " + ctx.Config.Paths.KernelDir
	if calls[0] != want {
		t.Fatalf("clone = %q\nwant    %q", calls[0], want)
	}
	if len(patch.Builtin(config.DefaultKernelRef, "arm64")) == 0 {
		t.Fatalf("no built-in patches for %s", config.DefaultKernelRef)
	}
	patched := false
	for _, c := range calls[1:] {
		if strings.HasPrefix(c, "patch -p1 -i ") {
			patched = true
		}
	}
	if !patched {
		t.Fatalf("built-in patches not applied: %v", calls)
	}
}

func TestKernelCloneFullAndNoPatches(t *testing.T) {
	ctx, rec := newKernelTestContext(t)
	cmd := buildKernelCloneCmd(ctx)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"--full", "--no-patches", "--ref", "v6.12"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	calls := gitAndPatchCalls(rec)
	if len(calls) != 1 || strings.Contains(calls[0], "--depth") ||
		!strings.Contains(calls[0], "--branch v6.12") {
		t.Fatalf("calls = %v", calls)
	}
}
