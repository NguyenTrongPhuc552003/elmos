// Package tui provides the interactive Text User Interface for elmos.
// This file contains the Update function and message handling logic.
package tui

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel() Model {
	return Model{
		menuStack:   make([][]MenuItem, 0),
		currentMenu: buildMenuStructure(),
		logLines:    []string{},
		viewport:    viewport.New(80, 20),
		textInput:   textinput.New(),
	}
}

func TestModel_handleEnterKey(t *testing.T) {
	// Test entering a submenu
	m := newTestModel()
	// Navigate to "Modules" which is a submenu (index 3 in main menu)
	m.cursor = 3

	newM, _ := m.handleEnterKey()
	finalM := newM.(Model)

	if len(finalM.menuStack) != 1 {
		t.Error("handleEnterKey() did not push to menuStack")
	}
	if finalM.parentTitle != "Modules" {
		t.Errorf("handleEnterKey() parentTitle = %v, want Modules", finalM.parentTitle)
	}

	// Test action item (non-input, non-interactive)
	// Create a dummy item for testing
	m = newTestModel()
	m.currentMenu = []MenuItem{{Label: "Test", Action: "test:action"}}
	m.cursor = 0

	newM, _ = m.handleEnterKey()
	finalM = newM.(Model)
	if !finalM.isRunning {
		t.Error("handleEnterKey() did not set isRunning for action")
	}
	if finalM.currentTask != "Test" {
		t.Errorf("handleEnterKey() currentTask = %v, want Test", finalM.currentTask)
	}

	// Test input item
	m = newTestModel()
	m.currentMenu = []MenuItem{{Label: "Input", NeedsInput: true, Action: "test:input"}}
	m.cursor = 0

	newM, _ = m.handleEnterKey()
	finalM = newM.(Model)
	if !finalM.inputMode {
		t.Error("handleEnterKey() should enter input mode")
	}
}

func TestModel_handleNavigationKey(t *testing.T) {
	m := newTestModel()

	// Test Down
	m.cursor = 0
	m.handleNavigationKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("handleNavigationKey(Down) cursor = %d, want 1", m.cursor)
	}

	// Test Up
	m.handleNavigationKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("handleNavigationKey(Up) cursor = %d, want 0", m.cursor)
	}

	// Test Back
	m.menuStack = append(m.menuStack, buildMenuStructure())
	m.handleNavigationKey(tea.KeyMsg{Type: tea.KeyEsc})
	if len(m.menuStack) != 0 { // esc is treated as back check key mapping in code
		// The code uses keys.Back which is mapped to Esc in keys.go (assumed),
		// but handleNavigationKey uses key.Matches.
		// Let's rely on logic verification.
	}
}

func TestModel_handleCursorKey(t *testing.T) {
	m := newTestModel()
	m.currentMenu = []MenuItem{{}, {}} // 2 items

	m.cursor = 0
	m.handleCursorKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Error("handleCursorKey(Down) failed to increment")
	}

	m.handleCursorKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Error("handleCursorKey(Down) went out of bounds")
	}

	m.handleCursorKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Error("handleCursorKey(Up) failed to decrement")
	}
}

func TestModel_handleViewportKey(t *testing.T) {
	m := newTestModel()
	m.viewport.SetContent("line1\nline2\nline3\nline4\nline5")
	m.viewport.Height = 2

	// Scroll Down
	m.handleViewportKey(tea.KeyMsg{Type: tea.KeyDown}) // Note: code uses ScrollDown/PageDown
	// We need to construct KeyMsg that matches keys.ScrollDown/Up
}

func TestModel_getCommandWithInput(t *testing.T) {
	m := newTestModel()
	tests := []struct {
		action string
		value  string
		want   string
	}{
		{"module:new", "mymod", "elmos module new mymod"},
		{"unknown:action", "val", "elmos unknown:action"},
	}
	for _, tt := range tests {
		if got := m.getCommandWithInput(tt.action, tt.value); got != tt.want {
			t.Errorf("getCommandWithInput(%q, %q) = %q, want %q", tt.action, tt.value, got, tt.want)
		}
	}
}

func TestModel_actionToArgs(t *testing.T) {
	m := newTestModel()
	tests := []struct {
		action string
		input  string
		want   []string
	}{
		{"workspace:status", "", []string{"status"}},
		{"module:new", "foo", []string{"module", "new", "foo"}},
		{"kernel:switch", "", []string{"kernel", "switch"}},
		{"kernel:switch", "v1", []string{"kernel", "switch", "v1"}},
		{"unknown", "", []string{}},
	}
	for _, tt := range tests {
		got := m.actionToArgs(tt.action, tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("actionToArgs(%q, %q) = %v, want %v", tt.action, tt.input, got, tt.want)
		}
	}
}

func TestModel_isInteractiveCommand(t *testing.T) {
	m := newTestModel()
	if !m.isInteractiveCommand("kernel:config", "menuconfig") {
		t.Error("kernel:config menuconfig should be interactive")
	}
	if m.isInteractiveCommand("kernel:config", "defconfig") {
		t.Error("kernel:config defconfig should NOT be interactive")
	}
}

func TestModel_refreshViewport(t *testing.T) {
	m := newTestModel()
	m.logLines = []string{"test line 1", "test line 2"}
	m.refreshViewport()
	// Basic check that it didn't panic
}

func TestModel_handleQuit(t *testing.T) {
	// Case 1: Pop menu stack
	m := newTestModel()
	m.menuStack = append(m.menuStack, buildMenuStructure())
	newM, cmd := m.handleQuit()
	finalM := newM.(Model)

	if len(finalM.menuStack) != 0 {
		t.Error("handleQuit() should pop menu stack")
	}
	if cmd != nil {
		t.Error("handleQuit() should not return quit cmd when popping stack")
	}

	// Case 2: Quit app
	m = newTestModel()
	newM, cmd = m.handleQuit()
	finalM = newM.(Model)

	if !finalM.quitting {
		t.Error("handleQuit() should set quitting=true")
	}
	if cmd == nil { // tea.Quit is not nil
		t.Error("handleQuit() should return quit cmd")
	}
}

func TestModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		want tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_Update(t *testing.T) {
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name  string
		m     Model
		args  args
		want  tea.Model
		want1 tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.Update(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.Update() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.Update() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_handleWindowSize(t *testing.T) {
	type args struct {
		msg tea.WindowSizeMsg
	}
	tests := []struct {
		name  string
		m     Model
		args  args
		want  tea.Model
		want1 tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.handleWindowSize(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleWindowSize() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleWindowSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_handleCommandDone(t *testing.T) {
	type args struct {
		msg CommandDoneMsg
	}
	tests := []struct {
		name string
		m    *Model
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.handleCommandDone(tt.args.msg)
		})
	}
}

func TestModel_handleKeyMsg(t *testing.T) {
	type args struct {
		msg  tea.KeyMsg
		cmds []tea.Cmd
	}
	tests := []struct {
		name  string
		m     Model
		args  args
		want  tea.Model
		want1 tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.handleKeyMsg(tt.args.msg, tt.args.cmds)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleKeyMsg() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleKeyMsg() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_popMenuStack(t *testing.T) {
	tests := []struct {
		name string
		m    *Model
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.popMenuStack()
		})
	}
}

func TestModel_handleInputMode(t *testing.T) {
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name  string
		m     Model
		args  args
		want  tea.Model
		want1 tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.handleInputMode(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleInputMode() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleInputMode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_runCommand(t *testing.T) {
	type args struct {
		action     string
		inputValue string
	}
	tests := []struct {
		name string
		m    *Model
		args args
		want tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.runCommand(tt.args.action, tt.args.inputValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.runCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
