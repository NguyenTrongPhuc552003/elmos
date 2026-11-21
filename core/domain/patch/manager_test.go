package patch

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

const original = "int answer(void)\n{\n\treturn 41;\n}\n"

const fix = `--- a/usr/answer.c
+++ b/usr/answer.c
@@ -1,4 +1,4 @@
 int answer(void)
 {
-	return 41;
+	return 42;
 }
`

// tree makes a kernel directory with one source file and a patches
// directory holding one patch for v6.18/generic.
func tree(t *testing.T) *elconfig.Config {
	t.Helper()
	root := t.TempDir()
	cfg := &elconfig.Config{}
	cfg.Paths.KernelDir = filepath.Join(root, "linux")
	cfg.Paths.PatchesDir = filepath.Join(root, "patches")
	for _, d := range []string{filepath.Join(cfg.Paths.KernelDir, "usr"), filepath.Join(cfg.Paths.PatchesDir, "v6.18", "generic")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, filepath.Join(cfg.Paths.KernelDir, "usr", "answer.c"), original)
	write(t, filepath.Join(cfg.Paths.PatchesDir, "v6.18", "generic", "0001-answer.patch"), fix)
	write(t, filepath.Join(cfg.Paths.PatchesDir, "v6.18", "generic", "README"), "not a patch\n")
	return cfg
}

func write(t *testing.T, name, text string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func needPatch(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		// kernel sources are patched on macOS and Linux hosts only; on the
		// Windows runner, Git's patch(1) fails on this tree (exit 3; not examined)
		t.Skip("kernel patching runs on macOS and Linux hosts")
	}
	if _, err := exec.LookPath("patch"); err != nil {
		t.Skip("patch(1) is not on PATH")
	}
}

func TestBuiltinListsTheEmbeddedPatches(t *testing.T) {
	got := Builtin("v6.18", "arm64")
	if len(got) == 0 || got[0].Arch != "generic" || !strings.Contains(got[0].Name, "gen_init_cpio") {
		t.Fatalf("Builtin(v6.18, arm64) = %+v, want the generic gen_init_cpio patch first", got)
	}
	if n := len(Builtin("v0.0", "arm64")); n != 0 {
		t.Fatalf("a ref without patches gave %d", n)
	}
	var arm bool
	for _, p := range AllBuiltin() {
		if p.Version == "v6.0" && p.Arch == "arm" {
			arm = true
		}
	}
	if !arm {
		t.Fatal("AllBuiltin misses the v6.0 arm patch")
	}
}

func TestListOnlyPatchFiles(t *testing.T) {
	cfg := tree(t)
	m := NewManager(executor.NewShellExecutor(), filesystem.NewOSFileSystem(), cfg)
	all, err := m.List()
	if err != nil || len(all) != 1 || all[0].Name != "0001-answer.patch" || all[0].Version != "v6.18" || all[0].Arch != "generic" {
		t.Fatalf("List() = %+v, %v", all, err)
	}
	if v, _ := m.GetPatchesForVersion("v6.0"); len(v) != 0 {
		t.Fatalf("GetPatchesForVersion(v6.0) = %+v", v)
	}
}

func TestApplyAgainAndReverse(t *testing.T) {
	needPatch(t)
	cfg := tree(t)
	m := NewManager(executor.NewShellExecutor(), filesystem.NewOSFileSystem(), cfg)
	src := filepath.Join(cfg.Paths.KernelDir, "usr", "answer.c")
	ctx := context.Background()

	// the "patches/" prefix a user may type is accepted
	if err := m.Apply(ctx, "patches/v6.18/generic/0001-answer.patch"); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !strings.Contains(read(t, src), "return 42;") {
		t.Fatal("Apply did not change the source")
	}
	if err := m.Apply(ctx, "v6.18/generic/0001-answer.patch"); err == nil || !strings.Contains(err.Error(), "already") {
		t.Fatalf("second Apply: %v, want already applied", err)
	}
	if err := m.Reverse(ctx, filepath.Join(cfg.Paths.PatchesDir, "v6.18", "generic", "0001-answer.patch")); err != nil {
		t.Fatalf("Reverse: %v", err)
	}
	if read(t, src) != original {
		t.Fatal("Reverse did not restore the source")
	}
	// reversing a patch that is not applied fails and changes nothing
	if err := m.Reverse(ctx, "v6.18/generic/0001-answer.patch"); err == nil {
		t.Fatal("Reverse of a patch that is not applied succeeded")
	}
	if read(t, src) != original {
		t.Fatal("a refused Reverse changed the source")
	}
}

func TestApplyMissingPatch(t *testing.T) {
	cfg := tree(t)
	mock := executor.NewMockExecutor()
	m := NewManager(mock, filesystem.NewOSFileSystem(), cfg)
	for _, name := range []string{"v6.18/generic/none.patch", filepath.Join(cfg.Paths.PatchesDir, "none.patch")} {
		if err := m.Apply(context.Background(), name); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("Apply(%s): %v, want not found", name, err)
		}
		if err := m.Reverse(context.Background(), name); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("Reverse(%s): %v, want not found", name, err)
		}
	}
	if len(mock.Calls) != 0 {
		t.Fatalf("patch(1) ran for a missing patch: %+v", mock.Calls)
	}
}
