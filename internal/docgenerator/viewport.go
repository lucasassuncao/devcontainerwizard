package docgenerator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

type docPane int

const (
	docPaneList docPane = iota
	docPaneView
)

const docStatusLines = 1

type docTUIModel struct {
	names    []string
	raw      map[string]string // raw markdown per name
	rendered map[string]string // glamour-rendered content per name (cache)

	// List panel
	cursor     int
	listOffset int
	listH      int // visible rows inside border
	listColW   int // total column width (content + 2 border chars)

	// Viewport panel
	vp     viewport.Model
	vpColW int // total column width (content + 2 border chars)
	vpH    int // content height inside border

	active docPane
	width  int
	height int
}

func newDocTUIModel(docs map[string]string) docTUIModel {
	names := make([]string, 0, len(docs))
	for k := range docs {
		names = append(names, k)
	}
	sort.Strings(names)

	return docTUIModel{
		names:    names,
		raw:      docs,
		rendered: make(map[string]string, len(docs)),
		active:   docPaneList,
	}
}

func (m docTUIModel) Init() tea.Cmd { return nil }

func (m docTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()
		m.invalidateRendered()
		m.loadCurrent()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.active == docPaneList {
				m.active = docPaneView
			} else {
				m.active = docPaneList
			}
			return m, nil
		}
		if m.active == docPaneList {
			m.handleListKey(msg.String())
		} else {
			m.handleViewportKey(msg.String())
		}
		return m, nil
	}
	return m, nil
}

func (m *docTUIModel) handleListKey(key string) {
	n := len(m.names)
	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.listOffset {
				m.listOffset = m.cursor
			}
			m.loadCurrent()
		}
	case "down", "j":
		if m.cursor < n-1 {
			m.cursor++
			if m.cursor >= m.listOffset+m.listH {
				m.listOffset = m.cursor - m.listH + 1
			}
			m.loadCurrent()
		}
	}
}

func (m *docTUIModel) handleViewportKey(key string) {
	switch key {
	case "up", "k":
		m.vp.ScrollUp(1)
	case "down", "j":
		m.vp.ScrollDown(1)
	case "pgup":
		m.vp.ScrollUp(m.vpH / 2)
	case "pgdn":
		m.vp.ScrollDown(m.vpH / 2)
	}
}

func (m *docTUIModel) relayout() {
	// Size the list column to the longest name + "▶ " prefix (2) + border (2) + 1 margin.
	maxName := 0
	for _, n := range m.names {
		if len(n) > maxName {
			maxName = len(n)
		}
	}
	m.listColW = maxName + 5 // "▶ " (2) + border L+R (2) + 1 trailing margin
	if m.listColW < 18 {
		m.listColW = 18
	}
	m.vpColW = m.width - m.listColW
	if m.vpColW < 22 {
		m.vpColW = 22
	}

	innerH := m.height - docStatusLines - 2 // 2 = top+bottom panel border
	if innerH < 1 {
		innerH = 1
	}
	m.listH = innerH
	m.vpH = innerH

	m.vp.Width = m.vpColW - 2 // subtract left+right border
	m.vp.Height = m.vpH

	if m.listOffset+m.listH <= m.cursor {
		m.listOffset = m.cursor - m.listH + 1
	}
	if m.listOffset < 0 {
		m.listOffset = 0
	}
}

func (m *docTUIModel) invalidateRendered() {
	m.rendered = make(map[string]string, len(m.names))
}

func (m *docTUIModel) renderDoc(name string) string {
	if r, ok := m.rendered[name]; ok {
		return r
	}
	raw := m.raw[name]
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.vp.Width),
	)
	if err != nil {
		m.rendered[name] = raw
		return raw
	}
	out, err := renderer.Render(raw)
	if err != nil {
		m.rendered[name] = raw
		return raw
	}
	m.rendered[name] = out
	return out
}

func (m *docTUIModel) loadCurrent() {
	if len(m.names) == 0 || m.vp.Width == 0 {
		return
	}
	m.vp.SetContent(m.renderDoc(m.names[m.cursor]))
	m.vp.GotoTop()
}

func (m docTUIModel) View() string {
	if m.width == 0 {
		return "Loading…"
	}

	// ── Left panel (topic list) ──────────────────────────────────────────────
	var listSB strings.Builder
	end := m.listOffset + m.listH
	if end > len(m.names) {
		end = len(m.names)
	}
	for i := m.listOffset; i < end; i++ {
		label := m.names[i]
		var line string
		if i == m.cursor {
			line = theme.SelectedItem.Render("▶ " + label)
		} else {
			line = theme.AvailableItem.Render("  " + label)
		}
		listSB.WriteString(line + "\n")
	}

	leftPanel := theme.PanelBorder(m.active == docPaneList).
		Width(m.listColW - 2).
		Height(m.listH).
		Render(listSB.String())

	// ── Right panel (viewport) ───────────────────────────────────────────────
	rightPanel := theme.PanelBorder(m.active == docPaneView).
		Width(m.vpColW - 2).
		Height(m.vpH).
		Render(m.vp.View())

	// ── Status bar ───────────────────────────────────────────────────────────
	status := theme.StatusBar.Render(
		"[Tab] switch panel  [↑/↓ j/k] navigate / scroll  [PgUp/PgDn] half-page  [q] quit",
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel) + "\n" + status
}

// RenderMarkdownDocsInTerminal launches the two-panel documentation TUI.
func RenderMarkdownDocsInTerminal(docs map[string]string) error {
	if len(docs) == 0 {
		return fmt.Errorf("no documentation to display")
	}
	p := tea.NewProgram(newDocTUIModel(docs), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run docs TUI: %w", err)
	}
	return nil
}
