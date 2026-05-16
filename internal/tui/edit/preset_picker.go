package edit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

// PresetSelectedMsg fires when the user picks a preset and confirms with Enter.
type PresetSelectedMsg struct{ Name string }

// PresetPickerCancelledMsg fires when the user dismisses the picker with Esc.
type PresetPickerCancelledMsg struct{}

// PresetPickerModel is a compact list popover for preset selection.
type PresetPickerModel struct {
	names  []string
	cursor int
	totalW int
	totalH int
}

// NewPresetPicker creates a picker preselecting current if present in names.
func NewPresetPicker(names []string, current string, totalW, totalH int) PresetPickerModel {
	cursor := 0
	for i, n := range names {
		if n == current {
			cursor = i
			break
		}
	}
	return PresetPickerModel{
		names:  names,
		cursor: cursor,
		totalW: totalW,
		totalH: totalH,
	}
}

// SelectedName returns the name of the currently-highlighted preset.
func (p PresetPickerModel) SelectedName() string {
	if p.cursor < 0 || p.cursor >= len(p.names) {
		return ""
	}
	return p.names[p.cursor]
}

func (p PresetPickerModel) Update(msg tea.Msg) (PresetPickerModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return p, nil
	}
	switch key.Type {
	case tea.KeyEsc:
		return p, func() tea.Msg { return PresetPickerCancelledMsg{} }
	case tea.KeyEnter:
		name := p.SelectedName()
		return p, func() tea.Msg { return PresetSelectedMsg{Name: name} }
	case tea.KeyUp:
		if p.cursor > 0 {
			p.cursor--
		}
	case tea.KeyDown:
		if p.cursor < len(p.names)-1 {
			p.cursor++
		}
	}
	return p, nil
}

func (p PresetPickerModel) View() string {
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render(" Preset "))
	lines = append(lines, strings.Repeat("─", 20))
	for i, n := range p.names {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == p.cursor {
			prefix = "▸ "
			style = style.Foreground(theme.Accent).Bold(true)
		}
		lines = append(lines, style.Render(prefix+n))
	}
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Faint(true).Render("[Enter] select  [Esc] cancel"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(0, 1).
		Render(strings.Join(lines, "\n"))

	return centerBox(box, p.totalW, p.totalH)
}
