// Package tui provides the interactive Text User Interface for elmos.
// This file contains the main Update function and core TUI orchestration.
package terminal

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// refreshViewport updates the viewport content with log lines.
func (m *Model) refreshViewport() {
	m.viewport.SetContent(strings.Join(m.logLines, "\n"))
	m.viewport.GotoBottom()
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model and handles all input events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.inputMode {
		return m.handleInputMode(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case CommandDoneMsg:
		m.handleCommandDone(msg)
	case tea.KeyMsg:
		return m.handleKeyMsg(msg, cmds)
	}
	return m, tea.Batch(cmds...)
}
