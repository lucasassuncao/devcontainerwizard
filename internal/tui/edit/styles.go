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
	dirtyStyle  = lipgloss.NewStyle().Foreground(theme.Warning)

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

// renderTitledPanel delegates to the shared theme implementation.
func renderTitledPanel(title string, width, height int, active bool, content string) string {
	return theme.RenderTitledPanel(title, width, height, active, content)
}

// renderHeader returns the single-line app header.
func renderHeader(file string, dirty bool, width int) string {
	info := file
	if dirty {
		info = file + " ● modified"
	}
	return theme.RenderHeader(info, "", width)
}
