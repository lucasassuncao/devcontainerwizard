package edit

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

// Reused directly from the shared palette. Defined here only as short aliases
// to keep the call sites in this package tidy.
var (
	existingItemStyle  = theme.ExistingItem
	availableItemStyle = theme.AvailableItem
	selectedItemStyle  = theme.SelectedItem
	sectionLabelStyle  = lipgloss.NewStyle().Bold(true).Foreground(theme.Dim).PaddingLeft(1)

	// Status bar
	statusStyle = theme.StatusBar
	dirtyStyle  = lipgloss.NewStyle().Foreground(theme.Warning)

	// Header bar
	headerTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(theme.AccentBright).PaddingLeft(1)
	headerInfoStyle  = lipgloss.NewStyle().Foreground(theme.Dim).PaddingRight(1)

	// Overlay
	overlayBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(theme.Accent).
				Padding(0, 1)
	overlayTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(theme.AccentBright)

	// Panel borders (used by the two-panel overlay; root panels use renderTitledPanel).
	panelStyle       = theme.PanelBorder(false)
	activePanelStyle = theme.PanelBorder(true)
)

// renderTitledPanel renders a rounded-border panel with the title embedded in
// the top edge. width and height are the OUTER dimensions (including borders).
func renderTitledPanel(title string, width, height int, active bool, content string) string {
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}

	borderColor := theme.Muted
	titleColor := theme.Dim
	if active {
		borderColor = theme.Accent
		titleColor = theme.AccentBright
	}

	innerW := width - 2
	titleSegment := lipgloss.NewStyle().Bold(true).Foreground(titleColor).Render(" " + title + " ")
	fillLen := innerW - 1 - lipgloss.Width(titleSegment)
	if fillLen < 0 {
		fillLen = 0
	}

	borderInk := lipgloss.NewStyle().Foreground(borderColor)
	top := borderInk.Render("╭─") + titleSegment + borderInk.Render(strings.Repeat("─", fillLen)+"╮")

	body := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		BorderForeground(borderColor).
		Width(innerW).
		Height(height - 2).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, top, body)
}

// renderHeader returns the single-line app header.
func renderHeader(file string, dirty bool, width int) string {
	left := headerTitleStyle.Render("devcontainer wizard")

	info := file
	if dirty {
		info = file + " ● modified"
	}
	right := headerInfoStyle.Render(info)

	spacerW := width - lipgloss.Width(left) - lipgloss.Width(right)
	if spacerW < 1 {
		spacerW = 1
	}
	return left + strings.Repeat(" ", spacerW) + right
}
