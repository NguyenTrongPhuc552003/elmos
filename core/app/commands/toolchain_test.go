package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
)

const testToolchainTarget = "riscv64-unknown-linux-gnu"

func TestToolchainsDirectSelectValidTarget(t *testing.T) {
	ctx, exec, root, _ := newToolchainSelectionContext(t)
	writeTestTargetConfig(t, root)

	if err := runToolchainSelection(ctx, testToolchainTarget); err != nil {
		t.Fatal(err)
	}
	selected := filepath.Join(ctx.Config.Paths.ToolchainsDir, ".config")
	content, err := os.ReadFile(selected)
	if err != nil || !strings.Contains(string(content), "CT_TARGET=\""+testToolchainTarget+"\"") {
		t.Fatalf("selected config = %q, err=%v", content, err)
	}
	if ctx.Config.Build.Arch != "riscv" {
		t.Fatalf("selected architecture = %q, want riscv", ctx.Config.Build.Arch)
	}
	for _, call := range exec.Calls {
		if call.Cmd != "mount" {
			t.Fatalf("unexpected external command: %s %v", call.Cmd, call.Args)
		}
	}
}

func TestToolchainsDirectSelectInvalidTarget(t *testing.T) {
	ctx, exec, _, configPath := newToolchainSelectionContext(t)
	exec.RunError = errors.New("unknown ct-ng sample")

	err := runToolchainSelection(ctx, "riscv64-invalid")
	if err == nil || !strings.Contains(err.Error(), "failed to select target riscv64-invalid") {
		t.Fatalf("selection error = %v", err)
	}
	if ctx.Config.Build.Arch != "arm64" {
		t.Fatalf("architecture changed after failed selection: %q", ctx.Config.Build.Arch)
	}
	selectedByCtNg := false
	for _, call := range exec.Calls {
		if strings.HasSuffix(call.Cmd, "ct-ng") && len(call.Args) == 1 && call.Args[0] == "riscv64-invalid" {
			selectedByCtNg = true
		}
	}
	if !selectedByCtNg {
		t.Fatal("invalid target was not passed to ct-ng for validation")
	}
	if _, err := os.Stat(filepath.Join(ctx.Config.Paths.ToolchainsDir, ".config")); !os.IsNotExist(err) {
		t.Fatalf("unexpected selected config after failure: %v", err)
	}
	reloaded, err := config.Load(configPath)
	if err != nil || reloaded.Build.Arch != "arm64" {
		t.Fatalf("persisted architecture after failure = %q, err=%v", reloaded.Build.Arch, err)
	}
}

func TestToolchainsDirectSelectionPersists(t *testing.T) {
	ctx, _, root, configPath := newToolchainSelectionContext(t)
	writeTestTargetConfig(t, root)
	if err := runToolchainSelection(ctx, testToolchainTarget); err != nil {
		t.Fatal(err)
	}

	reloaded, err := config.Load(configPath)
	if err != nil || reloaded.Build.Arch != "riscv" {
		t.Fatalf("persisted architecture = %q, err=%v", reloaded.Build.Arch, err)
	}
	if _, err := os.Stat(filepath.Join(ctx.Config.Paths.ToolchainsDir, ".config")); err != nil {
		t.Fatalf("selected ct-ng config was not persisted: %v", err)
	}
}

func newToolchainSelectionContext(t *testing.T) (*Context, *executor.MockExecutor, string, string) {
	t.Helper()
	root := t.TempDir()
	mountPoint := filepath.Join(root, "mounted")
	toolchainsDir := filepath.Join(root, "toolchains")
	configPath := filepath.Join(root, "elmos.yaml")
	if err := os.Mkdir(mountPoint, 0755); err != nil {
		t.Fatal(err)
	}
	content := "build:\n  arch: arm64\nimage:\n  mount_point: " + mountPoint +
		"\npaths:\n  project_root: " + root + "\n  toolchains_dir: " + toolchainsDir + "\n"
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	ctngPath := filepath.Join(toolchainsDir, "crosstool-ng", "bin", "ct-ng")
	if err := os.MkdirAll(filepath.Dir(ctngPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ctngPath, nil, 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	exec := executor.NewMockExecutor()
	exec.OutputResponses["mount"] = []byte("test volume on " + mountPoint)
	fs := filesystem.NewOSFileSystem()
	printer := ui.NewPrinter()
	ctx := &Context{
		Config:           cfg,
		AppContext:       elcontext.New(cfg, exec, fs),
		ToolchainManager: toolchain.NewManager(exec, fs, cfg, printer),
		Printer:          printer,
	}
	return ctx, exec, root, configPath
}

func writeTestTargetConfig(t *testing.T, root string) {
	t.Helper()
	configDir := filepath.Join(root, "assets", "toolchains", "configs")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "CT_TARGET=\"" + testToolchainTarget + "\"\nCT_PREFIX_DIR=\"/old/x-tools\"\n"
	if err := os.WriteFile(filepath.Join(configDir, testToolchainTarget+".config"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func runToolchainSelection(ctx *Context, target string) error {
	cmd := BuildToolchains(ctx)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{target})
	return cmd.Execute()
}
