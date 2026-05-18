package edit

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

// Reused directly from the shared palette. Defined here only as short aliases
// to keep the call sites in this package tidy.
var (
	existingItemStyle  = theme.ExistingItem
	availableItemStyle = theme.AvailableItem
	selectedItemStyle  = theme.SelectedItem
	sectionLabelStyle  = lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).PaddingLeft(1)

	// Status bar
	statusStyle = theme.StatusBar

	// Filter prompt (shown at the bottom of the left panel in filter mode)
	filterPromptStyle = lipgloss.NewStyle().Bold(true).Foreground(theme.AccentBright)

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

// centerBox positions box at the centre of a totalW×totalH terminal region
// by adding padding. Used by all floating overlay/alert/picker views.
func centerBox(box string, totalW, totalH int) string {
	bw := lipgloss.Width(box)
	bh := lipgloss.Height(box)
	lp := (totalW - bw) / 2
	tp := (totalH - bh) / 2
	if lp < 0 {
		lp = 0
	}
	if tp < 0 {
		tp = 0
	}
	return lipgloss.NewStyle().PaddingLeft(lp).PaddingTop(tp).Render(box)
}

// renderHeader returns the single-line app header.
func renderHeader(file string, dirty bool, width int) string {
	info := file
	if dirty {
		info = file + " ● modified"
	}
	return theme.RenderHeader(info, "", width)
}
