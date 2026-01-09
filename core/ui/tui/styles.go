// Package tui provides the interactive Text User Interface for elmos.
// This file contains color definitions and lipgloss styles.
package tui

import "github.com/charmbracelet/lipgloss"

// ANSI color palette
var (
	purple      = lipgloss.Color("141")
	cyan        = lipgloss.Color("51")
	green       = lipgloss.Color("120")
	orange      = lipgloss.Color("214")
	red         = lipgloss.Color("203")
	white       = lipgloss.Color("255")
	grey        = lipgloss.Color("245")
	darkGrey    = lipgloss.Color("238")
	borderColor = lipgloss.Color("240")
)

// Panel and component styles
var (
	leftPanelStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(purple)
	rightPanelStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(borderColor)
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(purple)
	menuItemStyle     = lipgloss.NewStyle().Foreground(grey)
	selectedItemStyle = lipgloss.NewStyle().Bold(true).Foreground(white).Background(purple)
	hintStyle         = lipgloss.NewStyle().Foreground(cyan).Border(lipgloss.RoundedBorder()).BorderForeground(cyan).Padding(0, 1)
	descStyle         = lipgloss.NewStyle().Foreground(darkGrey).Italic(true)
	inputStyle        = lipgloss.NewStyle().Foreground(white).Background(lipgloss.Color("236")).Padding(0, 1)
	inputLabelStyle   = lipgloss.NewStyle().Foreground(orange).Bold(true)
)
