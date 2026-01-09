// Package tui provides the interactive Text User Interface for elmos.
// This file contains the View function for UI rendering.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model and renders the entire TUI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	panelHeight := m.height - 2

	// LEFT PANEL - Menu
	var left strings.Builder
	title := "ELMOS"
	if m.parentTitle != "" {
		title = m.parentTitle
	}
	left.WriteString(titleStyle.Render("─ "+title+" ─") + "\n\n")
	if len(m.menuStack) > 0 {
		left.WriteString(lipgloss.NewStyle().Foreground(darkGrey).Render("  ← Back (Esc)") + "\n\n")
	}
	for i, item := range m.currentMenu {
		var prefix string
		if len(item.Children) > 0 {
			prefix = "▸ "
		} else if item.Interactive {
			prefix = "⚡"
		} else if item.NeedsInput {
			prefix = "✎ "
		} else {
			prefix = "• "
		}
		label := prefix + item.Label
		if i == m.cursor {
			left.WriteString(selectedItemStyle.Render(" "+label+" ") + "\n")
		} else {
			left.WriteString(menuItemStyle.Render(label) + "\n")
		}
	}
	for i := strings.Count(left.String(), "\n"); i < panelHeight-4; i++ {
		left.WriteString("\n")
	}

	// RIGHT PANEL - Output and Input
	var right strings.Builder
	rightTitle := "Output"
	if m.isRunning {
		rightTitle = m.spinner.View() + " " + m.currentTask
	} else if m.inputMode {
		rightTitle = "📝 Input Required"
	}
	scrollInfo := ""
	if m.viewport.TotalLineCount() > m.viewport.Height {
		scrollInfo = fmt.Sprintf(" [%d%%]", int(m.viewport.ScrollPercent()*100))
	}
	right.WriteString(titleStyle.Render("─ "+rightTitle+scrollInfo+" ─") + "\n\n")

	// Show input field if in input mode
	if m.inputMode {
		right.WriteString(inputLabelStyle.Render(m.inputPrompt) + "\n\n")
		right.WriteString(inputStyle.Render(m.textInput.View()) + "\n\n")
		right.WriteString(descStyle.Render("  Press Enter to confirm, Esc to cancel") + "\n\n")
	} else if m.cursor < len(m.currentMenu) && !m.isRunning {
		item := m.currentMenu[m.cursor]
		if item.Command != "" {
			right.WriteString(hintStyle.Render(" $ "+item.Command+" ") + "\n")
			if item.Desc != "" {
				right.WriteString(descStyle.Render("  "+item.Desc) + "\n")
			}
			if item.NeedsInput {
				right.WriteString("\n" + inputLabelStyle.Render("  ✎ Press Enter to type: "+item.InputPrompt) + "\n")
			}
		} else if len(item.Children) > 0 {
			right.WriteString(descStyle.Render("  Press Enter to expand → "+item.Desc) + "\n")
		}
		right.WriteString("\n")
	}
	right.WriteString(m.viewport.View())

	// Combine panels
	leftPanel := leftPanelStyle.Width(m.leftWidth).Height(panelHeight).Render(left.String())
	rightPanel := rightPanelStyle.Width(m.rightWidth).Height(panelHeight).Render(right.String())
	main := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Footer with keybindings
	footer := lipgloss.NewStyle().Foreground(darkGrey).Render(
		lipgloss.NewStyle().Foreground(cyan).Render("↑↓") + " Navigate  " +
			lipgloss.NewStyle().Foreground(cyan).Render("⏎") + " Select  " +
			lipgloss.NewStyle().Foreground(cyan).Render("Esc") + " Back  " +
			lipgloss.NewStyle().Foreground(cyan).Render("[ ]") + " Scroll  " +
			lipgloss.NewStyle().Foreground(cyan).Render("c") + " Clear  " +
			lipgloss.NewStyle().Foreground(cyan).Render("q") + " Quit")

	return lipgloss.JoinVertical(lipgloss.Left, main, footer)
}
