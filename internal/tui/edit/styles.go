package edit

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Palette ───────────────────────────────────────────────────────────────────
var (
	colorAccent       = lipgloss.Color("63")  // blue — active borders
	colorAccentBright = lipgloss.Color("212") // pink — titles, selection
	colorMuted        = lipgloss.Color("240") // grey — inactive borders, hints
	colorDim          = lipgloss.Color("245") // lighter grey — secondary text
	colorSuccess      = lipgloss.Color("82")  // green — added items
	colorWarning      = lipgloss.Color("214") // orange — dirty marker
)

var (
	// List items
	existingItemStyle  = lipgloss.NewStyle().Foreground(colorSuccess)
	availableItemStyle = lipgloss.NewStyle().Foreground(colorDim)
	selectedItemStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorAccentBright)
	sectionLabelStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorDim).PaddingLeft(1)

	// Status bar
	statusStyle = lipgloss.NewStyle().Foreground(colorMuted).PaddingLeft(1)
	dirtyStyle  = lipgloss.NewStyle().Foreground(colorWarning)

	// Header bar
	headerTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorAccentBright).PaddingLeft(1)
	headerInfoStyle  = lipgloss.NewStyle().Foreground(colorDim).PaddingRight(1)

	// Overlay
	overlayBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)
	overlayTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorAccentBright)

	// Panel borders (used by the two-panel overlay; root panels use renderTitledPanel)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted)
	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent)
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

	borderColor := colorMuted
	titleColor := colorDim
	if active {
		borderColor = colorAccent
		titleColor = colorAccentBright
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
