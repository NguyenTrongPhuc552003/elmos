// Package tui provides the interactive Text User Interface for elmos.
// This file contains type definitions and key bindings.
package terminal

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// MenuItem represents a menu entry in the TUI.
type MenuItem struct {
	Label, Desc, Action, Command string
	Interactive                  bool     // Whether this runs an interactive command
	NeedsInput                   bool     // Whether this option needs text input
	InputPrompt                  string   // Prompt to show for text input
	InputPlaceholder             string   // Placeholder text
	Args                         []string // Args for interactive commands
	Children                     []MenuItem
}

// Model represents the TUI application state.
type Model struct {
	menuStack             [][]MenuItem
	currentMenu           []MenuItem
	cursor                int
	parentTitle           string
	viewport              viewport.Model
	logLines              []string
	spinner               spinner.Model
	isRunning             bool
	currentTask           string
	width, height         int
	leftWidth, rightWidth int
	quitting              bool
	execPath              string

	// Text input state
	textInput   textinput.Model
	inputMode   bool
	inputAction string
	inputPrompt string
}

// CommandDoneMsg is sent when a command finishes execution.
type CommandDoneMsg struct {
	Action string
	Err    error
	Output string
}

// keyMap defines keyboard shortcuts for the TUI.
type keyMap struct {
	Up, Down, Enter, Back, Quit, Clear     key.Binding
	PageUp, PageDown, ScrollUp, ScrollDown key.Binding
}

// keys contains the default key bindings.
var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k")),
	Down:       key.NewBinding(key.WithKeys("down", "j")),
	Enter:      key.NewBinding(key.WithKeys("enter")),
	Back:       key.NewBinding(key.WithKeys("esc", "backspace")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c")),
	Clear:      key.NewBinding(key.WithKeys("c")),
	PageUp:     key.NewBinding(key.WithKeys("[", "ctrl+u")),
	PageDown:   key.NewBinding(key.WithKeys("]", "ctrl+d")),
	ScrollUp:   key.NewBinding(key.WithKeys("{")),
	ScrollDown: key.NewBinding(key.WithKeys("}")),
}
