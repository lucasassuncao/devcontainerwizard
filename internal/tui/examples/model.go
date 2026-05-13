package examples

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

type pane int

const (
	paneList pane = iota
	paneViewport
)

// PresetYAMLFn returns the YAML body for a (field, preset) pair.
type PresetYAMLFn func(field, name string) (string, error)

// ListPresetsFn returns preset names for a field.
type ListPresetsFn func(field string) []string

// Model is the Bubble Tea root for the show-examples TUI.
type Model struct {
	fields []string
	list   listModel
	yamlFn PresetYAMLFn
	listFn ListPresetsFn

	width  int
	height int
	listW  int
	vpW    int

	vpScroll int

	active pane
}

// NewModel constructs the TUI given the fields and accessor functions.
func NewModel(fields []string, listFn ListPresetsFn, yamlFn PresetYAMLFn) Model {
	presetsByField := make(map[string][]string, len(fields))
	for _, f := range fields {
		presetsByField[f] = listFn(f)
	}
	return Model{
		fields: fields,
		list:   newListModel(fields, presetsByField),
		listFn: listFn,
		yamlFn: yamlFn,
		active: paneList,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.active == paneList {
				m.active = paneViewport
			} else {
				m.active = paneList
			}
			return m, nil
		case "up", "k":
			if m.active == paneList {
				m.list.MoveUp()
			} else if m.vpScroll > 0 {
				m.vpScroll--
			}
			return m, nil
		case "down", "j":
			if m.active == paneList {
				m.list.MoveDown()
			} else {
				m.vpScroll++
			}
			return m, nil
		case "enter", "l", "right":
			if m.active == paneList {
				if m.list.Mode() == modeFields {
					m.list.DrillIn()
				}
			}
			return m, nil
		case "esc", "h", "left":
			if m.active == paneList {
				if m.list.Mode() == modePresets {
					m.list.Back()
				}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) relayout() {
	m.listW = m.width / 6
	if m.listW < 20 {
		m.listW = 20
	}
	m.vpW = m.width - m.listW - 4
	if m.vpW < 20 {
		m.vpW = 20
	}
	innerH := m.height - 4 // 1 header + 1 hint + 2 panel borders
	if innerH < 3 {
		innerH = 3
	}
	m.list.SetSize(m.listW-2, innerH)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	if len(m.fields) == 0 {
		return "No presets available."
	}

	field, preset := m.list.Selected()
	yaml := ""
	if preset != "" {
		if y, err := m.yamlFn(field, preset); err == nil {
			yaml = y
		} else {
			yaml = "# error: " + err.Error()
		}
	}

	rightTitle := "Preset"
	if field != "" && preset != "" {
		rightTitle = fmt.Sprintf("%s · %s", field, preset)
	}

	innerH := m.height - 4
	if innerH < 3 {
		innerH = 3
	}

	leftPanel := theme.RenderTitledPanel("Fields", m.listW, innerH+2, m.active == paneList, m.list.View())
	rightPanel := theme.RenderTitledPanel(rightTitle, m.vpW, innerH+2, m.active == paneViewport,
		renderYAML(yaml, m.vpW-2))

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	hintText := "[↑/↓] navigate  [Enter/→] open  [Esc/←] back  [Tab] panel  [q] quit"
	if m.list.Mode() == modePresets {
		hintText = "[↑/↓] navigate  [Esc/←] back to fields  [Tab] panel  [q] quit"
	}
	hint := lipgloss.NewStyle().Faint(true).Render(hintText)
	header := theme.RenderHeader("presets", "", m.width)
	return strings.Join([]string{header, body, hint}, "\n")
}

// Run starts the show-examples TUI as a blocking call.
func Run(fields []string, listFn ListPresetsFn, yamlFn PresetYAMLFn) error {
	m := NewModel(fields, listFn, yamlFn)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
