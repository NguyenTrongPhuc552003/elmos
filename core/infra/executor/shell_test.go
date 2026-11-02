package executor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A tool placed on the child PATH must win over the one on the PATH of the
// elmos process (audit E11: the kernel build ran macOS make 3.81 although
// Homebrew GNU make was prepended to the child PATH).
func TestRunWithEnvUsesChildPath(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	stub := filepath.Join(dir, "elmos-test-tool")
	script := "#!/bin/sh\necho child-path > " + out + "\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	e := &ShellExecutor{}
	env := []string{"PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")}
	if err := e.RunWithEnv(context.Background(), env, "elmos-test-tool"); err != nil {
		t.Fatalf("tool on the child PATH was not found: %v", err)
	}
	got, err := os.ReadFile(out)
	if err != nil || strings.TrimSpace(string(got)) != "child-path" {
		t.Fatalf("stub did not run: %q %v", got, err)
	}
}

func TestResolveInEnvKeepsExplicitPaths(t *testing.T) {
	if got := resolveInEnv("/bin/sh", []string{"PATH=/nowhere"}); got != "/bin/sh" {
		t.Fatalf("explicit path rewritten to %q", got)
	}
	if got := resolveInEnv("sh", nil); got != "sh" {
		t.Fatalf("no env PATH should leave the name alone, got %q", got)
	}
}
