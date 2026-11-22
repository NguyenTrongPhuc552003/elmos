package toolchain

import (
	"strings"
	"testing"
)

func TestReplaceConfigValue(t *testing.T) {
	in := "CT_PREFIX_DIR=\"${HOME}/x-tools/${CT_TARGET}\"\nCT_OTHER=\"keep\"\nCT_PREFIX_DIR=\"again\"\n"
	got := replaceConfigValue(in, "CT_PREFIX_DIR", "/w/x-tools")
	want := "CT_PREFIX_DIR=\"/w/x-tools\"\nCT_OTHER=\"keep\"\nCT_PREFIX_DIR=\"/w/x-tools\"\n"
	if got != want {
		t.Fatalf("got %q\nwant %q", got, want)
	}
	if got := replaceConfigValue("CT_X=\"unterminated\n", "CT_X", "v"); got != "CT_X=\"unterminated\n" {
		t.Fatalf("an unterminated value changed: %q", got)
	}
}

// A key must match a whole key: MY_CT_PREFIX_DIR is not CT_PREFIX_DIR.
func TestReplaceConfigValueWholeKey(t *testing.T) {
	in := "MY_CT_PREFIX_DIR=\"mine\"\nCT_PREFIX_DIR=\"old\"\n"
	got := replaceConfigValue(in, "CT_PREFIX_DIR", "new")
	if got != "MY_CT_PREFIX_DIR=\"mine\"\nCT_PREFIX_DIR=\"new\"\n" {
		t.Fatalf("got %q", got)
	}
}

func TestPatchConfigContent(t *testing.T) {
	t.Setenv("HOME", "/w/hh")
	in := strings.Join([]string{
		"CT_PREFIX_DIR=\"${CT_PREFIX:-${HOME}/x-tools}/${CT_TARGET}\"",
		"CT_LOCAL_TARBALLS_DIR=\"/w/hh/src\"",
		"CT_WORK_DIR=\"/w/hh/x-tools/work\"",
		"CT_COMP_TOOLS_M4=y",
		"CT_COMP_TOOLS_MAKE=y",
		"CT_COMP_TOOLS_AUTOCONF=y",
		"CT_COMP_TOOLS_AUTOMAKE=y",
		"CT_COMP_TOOLS_LIBTOOL=y",
		"CT_ARCH_ARM=y",
	}, "\n")
	got := patchConfigContent(in, ToolchainPaths{XTools: "/v/x-tools", Src: "/v/src"})
	for _, want := range []string{
		"CT_PREFIX_DIR=\"/v/x-tools/${CT_TARGET}\"",
		"CT_LOCAL_TARBALLS_DIR=\"/v/src\"",
		"CT_WORK_DIR=\"/v/x-tools/work\"",
		"# CT_COMP_TOOLS_M4 is not set",
		"# CT_COMP_TOOLS_LIBTOOL is not set",
		"CT_ARCH_ARM=y",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "/w/hh") || strings.Contains(got, "CT_COMP_TOOLS_MAKE=y") {
		t.Errorf("home path or companion tool left in\n%s", got)
	}
}

func TestAppendOrUpdateEnv(t *testing.T) {
	env := appendOrUpdateEnv([]string{"A=1"}, "LDFLAGS", " -L/x")
	if env[len(env)-1] != "LDFLAGS= -L/x" {
		t.Fatalf("new variable: %v", env)
	}
	env = appendOrUpdateEnv([]string{"LDFLAGS=-L/a"}, "LDFLAGS", " -L/x")
	if env[0] != "LDFLAGS=-L/a -L/x" {
		t.Fatalf("flags: %v", env)
	}
}

// A search path keeps its entries apart: PKG_CONFIG_PATH=/a plus /b is /a:/b.
func TestAppendOrUpdatePathKeepsTheSeparator(t *testing.T) {
	env := appendOrUpdatePath([]string{"PKG_CONFIG_PATH=/a"}, "PKG_CONFIG_PATH", "/b")
	if env[0] != "PKG_CONFIG_PATH=/a:/b" {
		t.Fatalf("got %v", env)
	}
	env = appendOrUpdatePath([]string{"X=1"}, "PKG_CONFIG_PATH", "/b")
	if env[1] != "PKG_CONFIG_PATH=/b" {
		t.Fatalf("got %v", env)
	}
}
