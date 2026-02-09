package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// handleWindowSize handles window resize events.
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width, m.height = msg.Width, msg.Height
	m.leftWidth = maxInt(25, int(float64(m.width)*0.30))
	m.rightWidth = m.width - m.leftWidth - 3
	m.viewport.Width = m.rightWidth - 4
	m.viewport.Height = m.height - 10
	m.refreshViewport()
	return m, nil
}

// handleCommandDone handles command completion messages.
func (m *Model) handleCommandDone(msg CommandDoneMsg) {
	m.isRunning = false
	if msg.Output != "" {
		for _, line := range strings.Split(strings.TrimSpace(msg.Output), "\n") {
			m.logLines = append(m.logLines, "  "+line)
		}
	}
	if msg.Err != nil {
		m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(red).Render(fmt.Sprintf("  ✗ Error: %v", msg.Err)))
	} else {
		m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(green).Render("  ✓ Completed"))
	}
	m.logLines = append(m.logLines, "")
	m.refreshViewport()
	m.currentTask = ""
}

// handleKeyMsg handles all keyboard input.
func (m Model) handleKeyMsg(msg tea.KeyMsg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	if m.isRunning && !key.Matches(msg, keys.Quit) {
		return m, tea.Batch(cmds...)
	}

	// Handle special keys that return early
	if key.Matches(msg, keys.Quit) {
		return m.handleQuit()
	}
	if key.Matches(msg, keys.Enter) {
		return m.handleEnterKey()
	}

	// Handle navigation and other keys
	m.handleNavigationKey(msg)
	return m, tea.Batch(cmds...)
}

// handleNavigationKey handles navigation and viewport keys.
func (m *Model) handleNavigationKey(msg tea.KeyMsg) {
	switch {
	case key.Matches(msg, keys.Back):
		m.popMenuStack()
	case key.Matches(msg, keys.Up), key.Matches(msg, keys.Down):
		m.handleCursorKey(msg)
	case key.Matches(msg, keys.PageUp), key.Matches(msg, keys.PageDown),
		key.Matches(msg, keys.ScrollUp), key.Matches(msg, keys.ScrollDown):
		m.handleViewportKey(msg)
	case key.Matches(msg, keys.Clear):
		m.logLines = make([]string, 0)
		m.refreshViewport()
	}
}

// handleCursorKey handles up/down cursor navigation.
func (m *Model) handleCursorKey(msg tea.KeyMsg) {
	if key.Matches(msg, keys.Up) && m.cursor > 0 {
		m.cursor--
	} else if key.Matches(msg, keys.Down) && m.cursor < len(m.currentMenu)-1 {
		m.cursor++
	}
}

// handleViewportKey handles viewport scrolling keys.
func (m *Model) handleViewportKey(msg tea.KeyMsg) {
	switch {
	case key.Matches(msg, keys.PageUp):
		m.viewport.PageUp()
	case key.Matches(msg, keys.PageDown):
		m.viewport.PageDown()
	case key.Matches(msg, keys.ScrollUp):
		m.viewport.ScrollUp(1)
	case key.Matches(msg, keys.ScrollDown):
		m.viewport.ScrollDown(1)
	}
}

// handleQuit handles quit/back navigation.
func (m Model) handleQuit() (tea.Model, tea.Cmd) {
	if len(m.menuStack) > 0 {
		m.currentMenu = m.menuStack[len(m.menuStack)-1]
		m.menuStack = m.menuStack[:len(m.menuStack)-1]
		m.cursor, m.parentTitle = 0, ""
		return m, nil
	}
	m.quitting = true
	return m, tea.Quit
}

// handleEnterKey handles Enter key press on menu items.
func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.currentMenu) {
		return m, nil
	}

	item := m.currentMenu[m.cursor]

	// Expand submenu
	if len(item.Children) > 0 {
		m.menuStack = append(m.menuStack, m.currentMenu)
		m.parentTitle = item.Label
		m.currentMenu = item.Children
		m.cursor = 0
		return m, nil
	}

	// Enter input mode
	if item.NeedsInput {
		m.inputMode = true
		m.inputAction = item.Action
		m.inputPrompt = item.InputPrompt
		m.textInput.Placeholder = item.InputPlaceholder
		m.textInput.SetValue("")
		m.textInput.Focus()
		return m, textinput.Blink
	}

	// Run interactive command
	if item.Interactive {
		m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+item.Command))
		m.refreshViewport()
		c := exec.Command(m.execPath, item.Args...)
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return CommandDoneMsg{Action: item.Action, Err: err}
		})
	}

	// Run background command
	if len(item.Args) > 0 {
		m.isRunning = true
		m.currentTask = item.Label
		m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+item.Command))
		m.refreshViewport()
		return m, m.runCommand(item.Action, item.Args)
	}

	// Fallback for dynamic actions that use runCommand
	if item.Action != "" {
		m.isRunning = true
		m.currentTask = item.Label
		m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+item.Command))
		m.refreshViewport()
		return m, m.runCommand(item.Action, []string{})
	}

	return m, nil
}

// handleInputMode handles keyboard input when in text input mode.
func (m Model) handleInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.inputMode = false
			m.textInput.Blur()
			return m, nil
		case "enter":
			value := strings.TrimSpace(m.textInput.Value())

			// Auto-list refs if user types '?' for switch command
			if m.inputAction == "kernel:switch" && value == "?" {
				m.inputMode = false
				m.textInput.Blur()
				m.isRunning = true
				m.currentTask = "Listing refs..."
				m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ elmos kernel switch"))
				m.refreshViewport()
				return m, m.runCommand("kernel:switch", []string{"kernel", "switch"})
			}

			if value == "" {
				return m, nil
			}
			m.inputMode = false
			m.textInput.Blur()

			// Build command with the input value
			cmdStr := m.getCommandWithInput(m.inputAction, value)
			m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+cmdStr))
			m.refreshViewport()

			// Check if the resulting command is interactive
			if m.isInteractiveCommand(m.inputAction, value) {
				args := m.actionToArgs(m.inputAction, value)
				c := exec.Command(m.execPath, args...)
				c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return CommandDoneMsg{Action: m.inputAction, Err: err}
				})
			}

			m.isRunning = true
			m.currentTask = m.inputPrompt + " " + value
			args := m.actionToArgs(m.inputAction, value)
			return m, m.runCommand(m.inputAction, args)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}
