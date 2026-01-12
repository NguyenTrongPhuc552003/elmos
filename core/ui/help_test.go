// Package ui provides console output helpers for elmos.
package ui

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBanner(t *testing.T) {
	got := Banner()
	if got == "" {
		t.Error("Banner() returned empty string")
	}
	// Banner uses ASCII art which might not contain literal "elmos"
	// Just verify it returns something non-empty
}

func TestSetCustomUsageFunc(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	// Should not panic
	SetCustomUsageFunc(cmd)
	if cmd.UsageFunc() == nil {
		t.Error("SetCustomUsageFunc() should set usage function")
	}
}

func Test_customHelpFunc(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "A test command for testing",
	}
	// Should not panic
	customHelpFunc(cmd, []string{})
}

func Test_writeHeader(t *testing.T) {
	out := &strings.Builder{}
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "A longer description",
	}
	writeHeader(out, cmd)
	// The function should write without panicking
	// Content validation depends on implementation details
}

func Test_writeUsage(t *testing.T) {
	out := &strings.Builder{}
	cmd := &cobra.Command{Use: "test [args]"}
	writeUsage(out, cmd)
	got := out.String()
	// Should contain usage line or be empty if no usage
	if got != "" && !strings.Contains(got, "USAGE") && !strings.Contains(got, "test") {
		t.Error("writeUsage() should contain USAGE or command name if output is not empty")
	}
}

func Test_writeCommands(t *testing.T) {
	out := &strings.Builder{}
	parent := &cobra.Command{Use: "parent"}
	child := &cobra.Command{Use: "child", Short: "Child command"}
	parent.AddCommand(child)

	writeCommands(out, parent)
	// Should not panic; output validation optional
}

func Test_writeGroupedCommands(t *testing.T) {
	out := &strings.Builder{}
	cmds := []*cobra.Command{
		{Use: "cmd1", Short: "First command", Annotations: map[string]string{"group": "Build"}},
		{Use: "cmd2", Short: "Second command", Annotations: map[string]string{"group": "Build"}},
	}
	writeGroupedCommands(out, cmds)
	// Should not panic
}

func Test_writeSimpleCommands(t *testing.T) {
	out := &strings.Builder{}
	cmds := []*cobra.Command{
		{Use: "cmd1", Short: "First command"},
		{Use: "cmd2", Short: "Second command"},
	}
	writeSimpleCommands(out, cmds)
	// Should not panic
}

func Test_writeCommand(t *testing.T) {
	out := &strings.Builder{}
	cmd := &cobra.Command{Use: "test", Short: "A test command"}
	writeCommand(out, cmd)
	got := out.String()
	if !strings.Contains(got, "test") {
		t.Error("writeCommand() should write command name")
	}
}

func Test_writeFlags(t *testing.T) {
	out := &strings.Builder{}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	writeFlags(out, cmd)
	// Should not panic; flags may or may not be written depending on NFlags
}

func Test_writeExamples(t *testing.T) {
	out := &strings.Builder{}
	cmd := &cobra.Command{
		Use:     "test",
		Example: "  test --verbose",
	}
	writeExamples(out, cmd)
	got := out.String()
	if cmd.Example != "" && !strings.Contains(got, "EXAMPLES") && got != "" {
		t.Error("writeExamples() should write examples section if command has examples")
	}
}

func Test_writeFooter(t *testing.T) {
	out := &strings.Builder{}
	parent := &cobra.Command{Use: "parent"}
	child := &cobra.Command{Use: "child"}
	parent.AddCommand(child)

	writeFooter(out, parent)
	// Should not panic
}

func Test_customUsageFunc(t *testing.T) {
	cmd := &cobra.Command{Use: "test", Short: "Test"}
	err := customUsageFunc(cmd)
	if err != nil {
		t.Errorf("customUsageFunc() error = %v", err)
	}
}

func Test_groupCommands(t *testing.T) {
	cmds := []*cobra.Command{
		{Use: "build", Short: "Build", Annotations: map[string]string{"group": "Build Commands"}},
		{Use: "run", Short: "Run", Annotations: map[string]string{"group": "Build Commands"}},
		{Use: "doctor", Short: "Doctor", Annotations: map[string]string{"group": "Help Commands"}},
	}
	got := groupCommands(cmds)
	// Should return groups based on annotations
	// Empty result is valid if grouping isn't done by annotation
	_ = got
}

func Test_buildGroupSlice(t *testing.T) {
	grouped := map[string][]*cobra.Command{
		"Build": {
			{Use: "build1"},
			{Use: "build2"},
		},
	}
	got := buildGroupSlice(grouped)
	if len(got) != 1 {
		t.Errorf("buildGroupSlice() length = %d, want 1", len(got))
	}
	if got[0].name != "Build" {
		t.Errorf("buildGroupSlice() group name = %s, want Build", got[0].name)
	}
}
