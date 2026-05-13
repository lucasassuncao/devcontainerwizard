// Package theme centralises the colours and base lipgloss styles shared by
// the project's TUIs (edit, show-docs). Layout helpers stay in each TUI
// because they differ; only the palette and a few primitive styles are common.
package theme

import "github.com/charmbracelet/lipgloss"

// Palette — kept narrow on purpose. Add a colour here only when at least two
// call sites need it; otherwise define it locally.
var (
	Accent       = lipgloss.Color("63")  // blue — active borders, primary highlight
	AccentBright = lipgloss.Color("212") // pink — titles, selection
	Muted        = lipgloss.Color("240") // grey — inactive borders, status hints
	Dim          = lipgloss.Color("245") // light grey — secondary text
	Success      = lipgloss.Color("82")  // green — existing/added items, success alerts
	Warning      = lipgloss.Color("214") // orange — dirty marker
	Danger       = lipgloss.Color("196") // red — error alerts
)

// Common item styles. Each TUI is free to compose its own variants on top.
var (
	SelectedItem  = lipgloss.NewStyle().Bold(true).Foreground(AccentBright)
	ExistingItem  = lipgloss.NewStyle().Foreground(Success)
	AvailableItem = lipgloss.NewStyle().Foreground(Dim)
	StatusBar     = lipgloss.NewStyle().Foreground(Muted).PaddingLeft(1)
)

// PanelBorder returns a rounded-border style coloured for the active/inactive
// state. Width/Height are left to the caller because layout differs per TUI.
func PanelBorder(active bool) lipgloss.Style {
	colour := Muted
	if active {
		colour = Accent
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colour)
}
