// Package theme centralises the colours, base lipgloss styles, and shared
// layout primitives used by all project TUIs.
package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

var (
	headerTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(AccentBright).PaddingLeft(1)
	headerInfoStyle  = lipgloss.NewStyle().Foreground(Dim).PaddingRight(1)
)

// RenderHeader returns the single-line app header used across all TUIs.
// subtitle, if non-empty, is rendered next to the title on the left (e.g. "presets", "docs").
// right is optional contextual info on the right side (e.g. filename).
func RenderHeader(subtitle, right string, width int) string {
	left := headerTitleStyle.Render("devcontainer wizard")
	if subtitle != "" {
		left += headerInfoStyle.Render(" · " + subtitle)
	}
	rightRendered := ""
	if right != "" {
		rightRendered = headerInfoStyle.Render(right)
	}
	spacerW := width - lipgloss.Width(left) - lipgloss.Width(rightRendered)
	if spacerW < 1 {
		spacerW = 1
	}
	return left + strings.Repeat(" ", spacerW) + rightRendered
}

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

// RenderTitledPanel renders a rounded-border panel with the title embedded in
// the top edge: ╭─ Title ──────╮. width and height are OUTER dimensions
// (including the border rows/cols). Same visual as the main edit TUI panels.
func RenderTitledPanel(title string, width, height int, active bool, content string) string {
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}

	borderColor := Muted
	titleColor := Dim
	if active {
		borderColor = Accent
		titleColor = AccentBright
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
