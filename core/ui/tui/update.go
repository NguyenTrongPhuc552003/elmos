// Package tui provides the interactive Text User Interface for elmos.
// This file contains the Update function and message handling logic.
package tui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// If in input mode, handle text input
	if m.inputMode {
		return m.handleInputMode(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.leftWidth = maxInt(25, int(float64(m.width)*0.30))
		m.rightWidth = m.width - m.leftWidth - 3
		m.viewport.Width = m.rightWidth - 4
		m.viewport.Height = m.height - 10
		m.refreshViewport()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case CommandDoneMsg:
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

	case tea.KeyMsg:
		if m.isRunning && !key.Matches(msg, keys.Quit) {
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			if len(m.menuStack) > 0 {
				m.currentMenu = m.menuStack[len(m.menuStack)-1]
				m.menuStack = m.menuStack[:len(m.menuStack)-1]
				m.cursor, m.parentTitle = 0, ""
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		case key.Matches(msg, keys.Back):
			if len(m.menuStack) > 0 {
				m.currentMenu = m.menuStack[len(m.menuStack)-1]
				m.menuStack = m.menuStack[:len(m.menuStack)-1]
				m.cursor, m.parentTitle = 0, ""
			}
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.currentMenu)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.PageUp):
			m.viewport.PageUp()
		case key.Matches(msg, keys.PageDown):
			m.viewport.PageDown()
		case key.Matches(msg, keys.ScrollUp):
			m.viewport.ScrollUp(1)
		case key.Matches(msg, keys.ScrollDown):
			m.viewport.ScrollDown(1)
		case key.Matches(msg, keys.Clear):
			m.logLines = make([]string, 0)
			m.refreshViewport()
		case key.Matches(msg, keys.Enter):
			if m.cursor < len(m.currentMenu) {
				item := m.currentMenu[m.cursor]
				if len(item.Children) > 0 {
					m.menuStack = append(m.menuStack, m.currentMenu)
					m.parentTitle = item.Label
					m.currentMenu = item.Children
					m.cursor = 0
					return m, nil
				}

				// If item needs input, enter input mode
				if item.NeedsInput {
					m.inputMode = true
					m.inputAction = item.Action
					m.inputPrompt = item.InputPrompt
					m.textInput.Placeholder = item.InputPlaceholder
					m.textInput.SetValue("")
					m.textInput.Focus()
					return m, textinput.Blink
				}

				if item.Interactive {
					m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+item.Command))
					m.refreshViewport()
					c := exec.Command(m.execPath, item.Args...)
					c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return CommandDoneMsg{Action: item.Action, Err: err}
					})
				}
				if item.Action != "" {
					m.isRunning = true
					m.currentTask = item.Label
					m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+item.Command))
					m.refreshViewport()
					return m, m.runCommand(item.Action, "")
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
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
			if value == "" {
				return m, nil
			}
			m.inputMode = false
			m.textInput.Blur()

			// Build command with the input value
			cmdStr := m.getCommandWithInput(m.inputAction, value)
			m.logLines = append(m.logLines, lipgloss.NewStyle().Foreground(cyan).Render("  ▶ "+cmdStr))
			m.refreshViewport()

			m.isRunning = true
			m.currentTask = m.inputPrompt + " " + value
			return m, m.runCommand(m.inputAction, value)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// getCommandWithInput returns the display command string for a given action and input.
func (m *Model) getCommandWithInput(action, value string) string {
	switch action {
	case "module:new":
		return "elmos module new " + value
	case "module:build:one":
		return "elmos module build " + value
	case "app:new":
		return "elmos app new " + value
	case "app:build:one":
		return "elmos app build " + value
	case "config:arch":
		return "elmos config set arch " + value
	case "config:jobs":
		return "elmos config set jobs " + value
	case "config:memory":
		return "elmos config set memory " + value
	case "rootfs:create:custom":
		return "elmos rootfs create -s " + value
	default:
		return "elmos " + action
	}
}

// runCommand executes a command asynchronously and returns the result.
func (m *Model) runCommand(action, inputValue string) tea.Cmd {
	return func() tea.Msg {
		args := m.actionToArgs(action, inputValue)
		cmd := exec.Command(m.execPath, args...)
		var output bytes.Buffer
		cmd.Stdout, cmd.Stderr = &output, &output
		err := cmd.Run()
		return CommandDoneMsg{Action: action, Err: err, Output: output.String()}
	}
}

// actionToArgs converts an action identifier to CLI arguments.
func (m *Model) actionToArgs(action, inputValue string) []string {
	switch action {
	case "init:workspace":
		return []string{"init"}
	case "workspace:status":
		return []string{"status"}
	case "workspace:exit":
		return []string{"exit"}
	case "gdb:connect":
		return []string{"gdb"}
	case "arch:show":
		return []string{"arch", "show"}
	case "arch:set":
		return []string{"arch", inputValue}
	case "kernel:status":
		return []string{"kernel", "status"}
	case "kernel:clone":
		return []string{"kernel", "clone"}
	case "kernel:pull":
		return []string{"kernel", "pull"}
	case "kernel:branch:list":
		return []string{"kernel", "branch"}
	case "kernel:branch:switch":
		return []string{"kernel", "branch", inputValue}
	case "kernel:reset":
		return []string{"kernel", "reset"}
	case "kernel:config":
		if inputValue == "" || inputValue == "defconfig" {
			return []string{"kernel", "config"}
		}
		return []string{"kernel", "config", inputValue}
	case "kernel:build":
		return []string{"kernel", "build"}
	case "kernel:clean":
		return []string{"kernel", "clean"}
	case "module:list":
		return []string{"module", "list"}
	case "module:build":
		if inputValue == "" {
			return []string{"module", "build"}
		}
		return []string{"module", "build", inputValue}
	case "module:new":
		return []string{"module", "new", inputValue}
	case "module:clean":
		return []string{"module", "clean"}
	case "app:list":
		return []string{"app", "list"}
	case "app:build":
		if inputValue == "" {
			return []string{"app", "build"}
		}
		return []string{"app", "build", inputValue}
	case "app:new":
		return []string{"app", "new", inputValue}
	case "app:clean":
		return []string{"app", "clean"}
	case "rootfs:create":
		return []string{"rootfs", "create"}
	case "rootfs:create:custom":
		return []string{"rootfs", "create", "-s", inputValue}
	case "config:show":
		return []string{"config", "show"}
	case "config:arch":
		return []string{"config", "set", "arch", inputValue}
	case "config:jobs":
		return []string{"config", "set", "jobs", inputValue}
	case "config:memory":
		return []string{"config", "set", "memory", inputValue}
	case "doctor:check":
		return []string{"doctor"}
	default:
		return []string{}
	}
}
