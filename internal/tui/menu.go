package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a single menu item
type MenuItem struct {
	Label       string
	Category    string
	Action      string
	Status      string // "ready", "pending", "error", or ""
	Description string // Description shown in right panel
}

// Category represents a menu category
type Category struct {
	Name     string
	Expanded bool
	Items    []MenuItem
}

// Styles
var (
	// Colors
	accentColor    = lipgloss.Color("#7C3AED") // Purple
	successColor   = lipgloss.Color("#10B981") // Green
	warningColor   = lipgloss.Color("#F59E0B") // Amber
	errorColor     = lipgloss.Color("#EF4444") // Red
	dimColor       = lipgloss.Color("#6B7280") // Gray
	highlightColor = lipgloss.Color("#A78BFA") // Light purple
	bgColor        = lipgloss.Color("#1F2937") // Dark background

	// Box styles for left panel
	leftPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	// Box styles for right panel
	rightPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dimColor).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(accentColor).
			Padding(0, 2)

	categoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			PaddingLeft(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlightColor).
			PaddingLeft(2).
			PaddingRight(2)

	statusReadyStyle   = lipgloss.NewStyle().Foreground(successColor)
	statusPendingStyle = lipgloss.NewStyle().Foreground(warningColor)
	statusErrorStyle   = lipgloss.NewStyle().Foreground(errorColor)

	footerStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Right panel styles
	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			MarginBottom(1)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9CA3AF"))

	outputSuccessStyle = lipgloss.NewStyle().Foreground(successColor)
	outputInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
	outputWarnStyle    = lipgloss.NewStyle().Foreground(warningColor)
	outputErrorStyle   = lipgloss.NewStyle().Foreground(errorColor)
)

// FlatItem represents a flattened menu item for navigation
type FlatItem struct {
	Category    *Category
	Item        *MenuItem
	IsCategory  bool
	CategoryIdx int
	ItemIdx     int
}

// MenuModel represents the TUI menu state
type MenuModel struct {
	categories     []Category
	flatItems      []FlatItem
	cursor         int
	choice         string
	quitting       bool
	width          int
	height         int
	showHelp       bool
	outputLines    []string
	outputViewport viewport.Model
	ready          bool
}

// Key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Toggle key.Binding
	Quit   key.Binding
	Help   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("tab", " "),
		key.WithHelp("tab", "toggle"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// NewMenuModel creates a new menu model with categories
func NewMenuModel() MenuModel {
	categories := []Category{
		{
			Name:     "Setup",
			Expanded: true,
			Items: []MenuItem{
				{
					Label:       "Doctor",
					Action:      "Doctor (Check Environment)",
					Status:      "ready",
					Description: "Check all dependencies and environment setup.\nVerifies: Homebrew, LLVM, QEMU, cross-compilers, headers.",
				},
				{
					Label:       "Init Workspace",
					Action:      "Init Workspace",
					Status:      "pending",
					Description: "Mount sparse image and clone Linux kernel.\nCreates a case-sensitive volume for kernel development.",
				},
				{
					Label:       "Configure",
					Action:      "Configure (Arch, Jobs...)",
					Description: "View and modify build configuration.\nSet architecture, parallel jobs, and other options.",
				},
			},
		},
		{
			Name:     "Build",
			Expanded: true,
			Items: []MenuItem{
				{
					Label:       "Kernel Config (defconfig)",
					Action:      "Kernel Config (defconfig)",
					Description: "Generate default kernel configuration.\nRuns 'make defconfig' for the target architecture.",
				},
				{
					Label:       "Kernel Menuconfig",
					Action:      "Kernel Menuconfig (UI)",
					Description: "Interactive kernel configuration.\nOpens the ncurses-based menuconfig interface.",
				},
				{
					Label:       "Build Kernel",
					Action:      "Build Kernel",
					Description: "Build kernel image, device trees, and modules.\nTargets: Image, dtbs, modules",
				},
				{
					Label:       "Build Modules",
					Action:      "Build Modules",
					Description: "Build out-of-tree kernel modules.\nCompiles modules in the modules/ directory.",
				},
				{
					Label:       "Build Apps",
					Action:      "Build Apps",
					Description: "Build userspace applications.\nCompiles apps in the apps/ directory for target arch.",
				},
			},
		},
		{
			Name:     "Run",
			Expanded: true,
			Items: []MenuItem{
				{
					Label:       "Run QEMU",
					Action:      "Run QEMU",
					Description: "Launch kernel in QEMU emulator.\nBoots the built kernel with the Debian rootfs.",
				},
				{
					Label:       "Run QEMU (Debug)",
					Action:      "Run QEMU (Debug Mode)",
					Description: "Launch QEMU with GDB stub enabled.\nConnects debugger on port 1234 for kernel debugging.",
				},
			},
		},
	}

	vp := viewport.New(40, 15)
	vp.Style = lipgloss.NewStyle()

	m := MenuModel{
		categories:     categories,
		width:          100,
		height:         24,
		outputViewport: vp,
		outputLines:    []string{},
	}
	m.buildFlatItems()
	m.updateOutputPanel()
	return m
}

// buildFlatItems creates a flat list for navigation
func (m *MenuModel) buildFlatItems() {
	m.flatItems = nil
	for catIdx := range m.categories {
		cat := &m.categories[catIdx]
		m.flatItems = append(m.flatItems, FlatItem{
			Category:    cat,
			IsCategory:  true,
			CategoryIdx: catIdx,
		})
		if cat.Expanded {
			for itemIdx := range cat.Items {
				m.flatItems = append(m.flatItems, FlatItem{
					Category:    cat,
					Item:        &cat.Items[itemIdx],
					IsCategory:  false,
					CategoryIdx: catIdx,
					ItemIdx:     itemIdx,
				})
			}
		}
	}
}

// updateOutputPanel updates the right panel content based on current selection
func (m *MenuModel) updateOutputPanel() {
	if m.cursor >= len(m.flatItems) {
		return
	}

	item := m.flatItems[m.cursor]
	var content strings.Builder

	if item.IsCategory {
		// Show category info
		content.WriteString(panelTitleStyle.Render(fmt.Sprintf("📁 %s", item.Category.Name)))
		content.WriteString("\n\n")
		content.WriteString(descriptionStyle.Render(fmt.Sprintf("Contains %d items.\nPress Enter or Tab to expand/collapse.", len(item.Category.Items))))
	} else if item.Item != nil {
		// Show item info
		content.WriteString(panelTitleStyle.Render(fmt.Sprintf("▶ %s", item.Item.Label)))
		content.WriteString("\n\n")
		if item.Item.Description != "" {
			content.WriteString(descriptionStyle.Render(item.Item.Description))
		}
		content.WriteString("\n\n")
		content.WriteString(footerStyle.Render("Press Enter to execute"))
	}

	m.outputViewport.SetContent(content.String())
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport size
		leftWidth := 45
		rightWidth := m.width - leftWidth - 6
		if rightWidth < 30 {
			rightWidth = 30
		}
		m.outputViewport.Width = rightWidth
		m.outputViewport.Height = m.height - 8
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.updateOutputPanel()
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.flatItems)-1 {
				m.cursor++
				m.updateOutputPanel()
			}

		case key.Matches(msg, keys.Toggle):
			if m.cursor < len(m.flatItems) {
				item := m.flatItems[m.cursor]
				if item.IsCategory {
					m.categories[item.CategoryIdx].Expanded = !m.categories[item.CategoryIdx].Expanded
					m.buildFlatItems()
					if m.cursor >= len(m.flatItems) {
						m.cursor = len(m.flatItems) - 1
					}
					m.updateOutputPanel()
				}
			}

		case key.Matches(msg, keys.Enter):
			if m.cursor < len(m.flatItems) {
				item := m.flatItems[m.cursor]
				if item.IsCategory {
					m.categories[item.CategoryIdx].Expanded = !m.categories[item.CategoryIdx].Expanded
					m.buildFlatItems()
					m.updateOutputPanel()
				} else if item.Item != nil {
					m.choice = item.Item.Action
					return m, tea.Quit
				}
			}

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
		}
	}

	m.outputViewport, cmd = m.outputViewport.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	if m.quitting {
		return ""
	}

	// Calculate panel widths
	leftWidth := 45
	rightWidth := m.width - leftWidth - 4
	if rightWidth < 30 {
		rightWidth = 30
	}
	panelHeight := m.height - 4

	// Build left panel (menu)
	var leftContent strings.Builder

	// Title
	leftContent.WriteString(titleStyle.Render("ELMOS"))
	leftContent.WriteString("\n\n")

	// Menu items
	for i, flatItem := range m.flatItems {
		isSelected := i == m.cursor

		if flatItem.IsCategory {
			arrow := "▼"
			if !flatItem.Category.Expanded {
				arrow = "▶"
			}
			catText := fmt.Sprintf("%s %s", arrow, flatItem.Category.Name)
			if isSelected {
				leftContent.WriteString(selectedStyle.Render(catText))
			} else {
				leftContent.WriteString(categoryStyle.Render(catText))
			}
			leftContent.WriteString("\n")
		} else if flatItem.Item != nil {
			status := m.renderStatus(flatItem.Item.Status)
			label := flatItem.Item.Label

			// Truncate if too long
			maxLen := leftWidth - 15
			if len(label) > maxLen {
				label = label[:maxLen-3] + "..."
			}

			if isSelected {
				line := fmt.Sprintf("  ▶ %s", label)
				if status != "" {
					padding := leftWidth - 12 - len(line)
					if padding < 1 {
						padding = 1
					}
					line += strings.Repeat(" ", padding) + status
				}
				leftContent.WriteString(selectedStyle.Render(line))
			} else {
				line := fmt.Sprintf("    %s", label)
				if status != "" {
					padding := leftWidth - 12 - len(line)
					if padding < 1 {
						padding = 1
					}
					line += strings.Repeat(" ", padding) + status
				}
				leftContent.WriteString(itemStyle.Render(line))
			}
			leftContent.WriteString("\n")
		}
	}

	// Pad to fill height
	menuLines := strings.Count(leftContent.String(), "\n")
	for i := menuLines; i < panelHeight-3; i++ {
		leftContent.WriteString("\n")
	}

	// Footer
	leftContent.WriteString(footerStyle.Render("↑↓:Nav  ⏎:Select  q:Quit"))

	// Build right panel (description/output)
	var rightContent strings.Builder

	// Right panel title
	rightContent.WriteString(panelTitleStyle.Render("📋 Details"))
	rightContent.WriteString("\n")
	rightContent.WriteString(strings.Repeat("─", rightWidth-4))
	rightContent.WriteString("\n\n")

	// Get current item description
	if m.cursor < len(m.flatItems) {
		item := m.flatItems[m.cursor]
		if item.IsCategory {
			rightContent.WriteString(outputInfoStyle.Render(fmt.Sprintf("📁 Category: %s", item.Category.Name)))
			rightContent.WriteString("\n\n")
			rightContent.WriteString(descriptionStyle.Render(fmt.Sprintf("Contains %d items.", len(item.Category.Items))))
			rightContent.WriteString("\n")
			rightContent.WriteString(descriptionStyle.Render("Use Tab or Enter to expand/collapse."))
		} else if item.Item != nil {
			rightContent.WriteString(outputInfoStyle.Render(fmt.Sprintf("▶ %s", item.Item.Label)))
			rightContent.WriteString("\n\n")
			if item.Item.Description != "" {
				// Word wrap description
				lines := strings.Split(item.Item.Description, "\n")
				for _, line := range lines {
					rightContent.WriteString(descriptionStyle.Render(line))
					rightContent.WriteString("\n")
				}
			}
			rightContent.WriteString("\n")
			rightContent.WriteString(footerStyle.Render("Press Enter to execute this action."))
		}
	}

	// Pad right panel
	rightLines := strings.Count(rightContent.String(), "\n")
	for i := rightLines; i < panelHeight-2; i++ {
		rightContent.WriteString("\n")
	}

	// Apply panel styles
	left := leftPanelStyle.Width(leftWidth).Height(panelHeight).Render(leftContent.String())
	right := rightPanelStyle.Width(rightWidth).Height(panelHeight).Render(rightContent.String())

	// Join panels horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m MenuModel) renderStatus(status string) string {
	switch status {
	case "ready":
		return statusReadyStyle.Render("[✓]")
	case "pending":
		return statusPendingStyle.Render("[○]")
	case "error":
		return statusErrorStyle.Render("[✗]")
	default:
		return ""
	}
}

// Choice returns the selected action
func (m MenuModel) Choice() string {
	return m.choice
}
