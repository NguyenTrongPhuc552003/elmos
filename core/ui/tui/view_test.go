// Package tui provides the interactive Text User Interface for elmos.
// This file contains tests for the View function and UI rendering.
package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

func TestModel_View(t *testing.T) {
	m := newTestModel()
	m.width = 100
	m.height = 30
	m.leftWidth = 30
	m.rightWidth = 70

	// Test normal view
	view := m.View()
	if !strings.Contains(view, "ELMOS") {
		t.Error("View() missing title ELMOS")
	}
	if !strings.Contains(view, "Quit") {
		t.Error("View() missing footer")
	}

	// Test quitting view
	m.quitting = true
	if got := m.View(); got != "" {
		t.Error("View() should be empty when quitting")
	}
}

func TestModel_renderLeftPanel(t *testing.T) {
	m := newTestModel()
	m.parentTitle = "Submenu"
	m.leftWidth = 30

	output := m.renderLeftPanel(20)

	if !strings.Contains(output, "Submenu") {
		t.Error("renderLeftPanel() missing parent title")
	}
	// Menu items should be present (default newTestModel has main menu)
	if !strings.Contains(output, "Kernel") {
		t.Error("renderLeftPanel() missing menu items")
	}
}

func TestModel_getMenuItemPrefix(t *testing.T) {
	tests := []struct {
		item MenuItem
		want string
	}{
		{MenuItem{Children: []MenuItem{{}}}, "▸ "},
		{MenuItem{Interactive: true}, "⚡"},
		{MenuItem{NeedsInput: true}, "✎ "},
		{MenuItem{}, "• "},
	}

	for _, tt := range tests {
		if got := getMenuItemPrefix(tt.item); got != tt.want {
			t.Errorf("getMenuItemPrefix() = %q, want %q", got, tt.want)
		}
	}
}

func TestModel_renderRightPanel(t *testing.T) {
	// Case 1: Normal with Hint
	m := newTestModel()
	m.rightWidth = 50
	m.currentMenu = []MenuItem{{Label: "Test", Command: "test cmd", Desc: "Description"}}
	m.cursor = 0

	output := m.renderRightPanel(20)
	if !strings.Contains(output, "test cmd") {
		t.Error("renderRightPanel() missing command hint")
	}
	if !strings.Contains(output, "Description") {
		t.Error("renderRightPanel() missing description")
	}

	// Case 2: Input Mode
	m.inputMode = true
	m.textInput = textinput.New()
	m.inputPrompt = "Enter value:"

	output = m.renderRightPanel(20)
	if !strings.Contains(output, "Enter value:") {
		t.Error("renderRightPanel() missing input prompt")
	}
	if !strings.Contains(output, "Input Required") {
		t.Error("renderRightPanel() missing Input Required title")
	}

	// Case 3: Running
	m.inputMode = false
	m.isRunning = true
	m.currentTask = "Running things..."

	output = m.renderRightPanel(20)
	if !strings.Contains(output, "Running things...") {
		t.Error("renderRightPanel() missing current task")
	}
}

func TestModel_getScrollInfo(t *testing.T) {
	m := newTestModel()
	m.viewport = viewport.New(10, 2)
	m.viewport.SetContent("line1\nline2\nline3") // Content larger than height

	// Scroll to bottom to get %
	m.viewport.GotoBottom()

	info := m.getScrollInfo()
	if !strings.Contains(info, "%") {
		t.Errorf("getScrollInfo() = %q, expected percentage", info)
	}

	m.viewport.SetContent("short")
	info = m.getScrollInfo()
	if info != "" {
		t.Errorf("getScrollInfo() = %q, expected empty for short content", info)
	}
}

func TestModel_renderFooter(t *testing.T) {
	m := newTestModel()
	output := m.renderFooter()
	if !strings.Contains(output, "Navigate") || !strings.Contains(output, "Quit") {
		t.Error("renderFooter() missing key hints")
	}
}
