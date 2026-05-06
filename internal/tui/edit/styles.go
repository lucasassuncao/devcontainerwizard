package edit

import "github.com/charmbracelet/lipgloss"

var (
	// Panel borders
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63"))

	// List items
	existingItemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))  // green
	availableItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // grey
	selectedItemStyle  = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")) // pink

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	// Status bar
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	dirtyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")) // orange

	// Overlay
	overlayBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(0, 1)

	overlayTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212"))
)
