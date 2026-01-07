package tui

import (
	"testing"
)

func TestMenuItemFields(t *testing.T) {
	item := MenuItem{
		Label:            "Test",
		Desc:             "Test description",
		Action:           "test:action",
		Command:          "elmos test",
		Interactive:      true,
		NeedsInput:       true,
		InputPrompt:      "Enter value:",
		InputPlaceholder: "value",
	}

	if item.Label != "Test" {
		t.Errorf("MenuItem.Label = %q, want %q", item.Label, "Test")
	}
	if item.Action != "test:action" {
		t.Errorf("MenuItem.Action = %q, want %q", item.Action, "test:action")
	}
	if !item.Interactive {
		t.Error("MenuItem.Interactive should be true")
	}
	if !item.NeedsInput {
		t.Error("MenuItem.NeedsInput should be true")
	}
}

func TestMenuItemChildren(t *testing.T) {
	parent := MenuItem{
		Label: "Parent",
		Children: []MenuItem{
			{Label: "Child1"},
			{Label: "Child2"},
		},
	}

	if len(parent.Children) != 2 {
		t.Errorf("MenuItem.Children len = %d, want 2", len(parent.Children))
	}
	if parent.Children[0].Label != "Child1" {
		t.Errorf("Children[0].Label = %q, want %q", parent.Children[0].Label, "Child1")
	}
}

func TestCommandDoneMsg(t *testing.T) {
	msg := CommandDoneMsg{
		Action: "build:kernel",
		Err:    nil,
		Output: "Build successful",
	}

	if msg.Action != "build:kernel" {
		t.Errorf("CommandDoneMsg.Action = %q, want %q", msg.Action, "build:kernel")
	}
	if msg.Err != nil {
		t.Error("CommandDoneMsg.Err should be nil")
	}
	if msg.Output != "Build successful" {
		t.Errorf("CommandDoneMsg.Output = %q, want %q", msg.Output, "Build successful")
	}
}

func TestKeyMapBindings(t *testing.T) {
	// Verify key bindings are initialized
	if len(keys.Up.Keys()) == 0 {
		t.Error("keys.Up should have key bindings")
	}
	if len(keys.Down.Keys()) == 0 {
		t.Error("keys.Down should have key bindings")
	}
	if len(keys.Enter.Keys()) == 0 {
		t.Error("keys.Enter should have key bindings")
	}
	if len(keys.Quit.Keys()) == 0 {
		t.Error("keys.Quit should have key bindings")
	}
}
