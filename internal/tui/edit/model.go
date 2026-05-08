package edit

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane int

const (
	paneList pane = iota
	panePreview
	paneOverlay
)

type quitState int

const (
	quitNone quitState = iota
	quitAsking
)

// Model is the root Bubble Tea model for the edit TUI.
type Model struct {
	filePath string

	rawYAML []byte
	blocks  []Block

	list    ListModel
	preview PreviewModel
	overlay *OverlayModel

	active    pane
	dirty     bool
	quitting  quitState
	statusMsg string

	width  int
	height int
	listW  int // derived by relayout(); read by View()
	innerH int // derived by relayout(); read by View()
}

// New loads the YAML file and initialises the model.
func New(filePath string) (Model, error) {
	raw, err := os.ReadFile(filePath) // #nosec G304 -- path is user-supplied via CLI arg
	if err != nil && !os.IsNotExist(err) {
		return Model{}, fmt.Errorf("reading %s: %w", filePath, err)
	}
	if raw == nil {
		raw = []byte{}
	}

	blocks, err := ParseBlocksFromBytes(raw)
	if err != nil {
		return Model{}, fmt.Errorf("parsing YAML: %w", err)
	}

	list := NewListModel(blocks, 0)
	preview := NewPreviewModel(0, 0)
	preview.SetContent(string(raw))

	return Model{
		filePath:  filePath,
		rawYAML:   raw,
		blocks:    blocks,
		list:      list,
		preview:   preview,
		active:    paneList,
		statusMsg: "",
	}, nil
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
		return m.handleKey(msg)

	case SpaceOnItemMsg:
		return m.handleSpace(msg.Item, msg.Guided)

	case OverlayConfirmedMsg:
		return m.handleOverlayConfirmed(msg.Snippet)

	case OverlayCancelledMsg:
		m.overlay = nil
		m.active = paneList
		m.statusMsg = "Cancelled."
		return m, nil

	case DeleteItemMsg:
		return m.handleDelete(msg.Key)
	}

	return m.forwardToActive(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Quit-confirm prompt intercepts all keys.
	if m.quitting == quitAsking {
		switch msg.String() {
		case "y", "Y":
			return m, tea.Quit
		default:
			m.quitting = quitNone
			m.statusMsg = "Quit cancelled."
			return m, nil
		}
	}

	// Overlay active — forward to it.
	if m.active == paneOverlay && m.overlay != nil {
		ov, cmd := m.overlay.Update(msg)
		m.overlay = &ov
		return m, cmd
	}

	switch msg.String() {
	case "ctrl+s":
		return m.save()
	case "q", "ctrl+c":
		if m.active == panePreview {
			break // let q through to the textarea
		}
		if m.dirty {
			m.quitting = quitAsking
			m.statusMsg = "Unsaved changes. Quit without saving? (y/N)"
			return m, nil
		}
		return m, tea.Quit
	case "tab":
		return m.togglePreviewPane()
	case "esc":
		if m.active == panePreview {
			return m.togglePreviewPane()
		}
	}

	if m.active == panePreview {
		return m.updatePreviewEditor(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if it := m.list.SelectedItem(); it != nil {
		m.preview.ScrollToKey(it.Key)
	}
	return m, cmd
}

func (m Model) togglePreviewPane() (tea.Model, tea.Cmd) {
	if m.active == panePreview {
		m.active = paneList
		m.preview.Blur()
		m.statusMsg = ""
		return m, nil
	}
	m.active = panePreview
	cmd := m.preview.Focus()
	m.statusMsg = "Editing YAML directly — Tab/Esc back to list."
	return m, cmd
}

func (m Model) updatePreviewEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prev, cmd := m.preview.Update(msg)
	m.preview = prev
	raw := []byte(m.preview.Value())
	if blocks, err := ParseBlocksFromBytes(raw); err == nil {
		m.rawYAML = raw
		m.blocks = blocks
		m.list.Rebuild(m.blocks)
		m.dirty = true
	}
	return m, cmd
}

func (m Model) handleSpace(it ListItem, guided bool) (tea.Model, tea.Cmd) {
	var initial string
	if it.Existing {
		// For existing items, pre-fill overlay with the current block content.
		current, err := BlockContent(m.rawYAML, m.blocks, it.Key)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Error reading %s: %v", it.Key, err)
			return m, nil
		}
		initial = current
	} else {
		initial = it.Key + ":\n"
		if guided && len(FieldsForKey(it.Key)) == 0 {
			// Only use the guided template for single-textarea blocks (no field defs).
			// Two-panel blocks initialise from rebuildYAML() when content is trivial.
			initial = GuidedTemplate(it.Key)
		}
	}

	mode := "free"
	if guided {
		mode = "guided"
	}

	ov := NewOverlay(it.Key, initial, guided, m.width, m.height)
	if it.Existing {
		ov.isEdit = true
		ov.editKey = it.Key
	}
	m.overlay = &ov
	m.active = paneOverlay
	m.statusMsg = fmt.Sprintf("Editing %q [%s] — Tab painel, Ctrl+S confirma, Esc cancela.", it.Key, mode)
	return m, nil
}

func (m Model) handleDelete(key string) (tea.Model, tea.Cmd) {
	newRaw, err := RemoveBlock(m.rawYAML, m.blocks, key)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Error removing %s: %v", key, err)
		return m, nil
	}
	m.applyRaw(newRaw)
	m.statusMsg = fmt.Sprintf("Removed %q (not saved yet).", key)
	return m, nil
}

func (m Model) handleOverlayConfirmed(snippet string) (tea.Model, tea.Cmd) {
	raw := m.rawYAML

	if m.overlay != nil && m.overlay.isEdit {
		// Replace: remove the old block first, then re-parse before inserting.
		removed, err := RemoveBlock(raw, m.blocks, m.overlay.editKey)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Remove error: %v", err)
			m.overlay = nil
			m.active = paneList
			return m, nil
		}
		raw = removed
		blocks, err := ParseBlocksFromBytes(raw)
		if err == nil {
			m.blocks = blocks
		}
	}

	isEdit := m.overlay != nil && m.overlay.isEdit

	newRaw, err := InsertBlock(raw, snippet)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Insert error: %v", err)
		m.overlay = nil
		m.active = paneList
		return m, nil
	}
	m.applyRaw(newRaw)
	m.overlay = nil
	m.active = paneList
	if isEdit {
		m.statusMsg = "Block updated (not saved yet)."
	} else {
		m.statusMsg = "Block added (not saved yet)."
	}
	return m, nil
}

func (m *Model) applyRaw(raw []byte) {
	m.rawYAML = raw
	blocks, err := ParseBlocksFromBytes(raw)
	if err == nil {
		m.blocks = blocks
	}
	m.list.Rebuild(m.blocks)
	m.preview.SetContent(string(raw))
	if it := m.list.SelectedItem(); it != nil {
		m.preview.ScrollToKey(it.Key)
	}
	m.dirty = true
}

func (m Model) save() (tea.Model, tea.Cmd) {
	if err := os.WriteFile(m.filePath, m.rawYAML, 0o600); err != nil {
		m.statusMsg = fmt.Sprintf("Save failed: %v", err)
		return m, nil
	}
	m.dirty = false
	m.statusMsg = fmt.Sprintf("Saved to %s.", m.filePath)
	return m, nil
}

const statusBarLines = 2 // feedback line + hint line

func (m *Model) relayout() {
	m.listW = m.width / 3
	previewW := m.width - m.listW - 4
	m.innerH = m.height - statusBarLines - 2 // 2 panel borders (top+bottom)

	if m.innerH < 1 {
		m.innerH = 1
	}
	if previewW < 10 {
		previewW = 10
	}

	m.list.height = m.innerH
	m.list.clampScroll()
	m.preview.Resize(previewW-2, m.innerH)
}

func (m Model) forwardToActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.active == paneOverlay && m.overlay != nil {
		ov, cmd := m.overlay.Update(msg)
		m.overlay = &ov
		return m, cmd
	}
	if m.active == panePreview {
		var cmd tea.Cmd
		m.preview, cmd = m.preview.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.active == paneOverlay && m.overlay != nil {
		return m.overlay.View()
	}

	listContent := m.list.View()
	listBorder := panelStyle
	if m.active == paneList {
		listBorder = activePanelStyle
	}
	leftPanel := listBorder.
		Width(m.listW - 2).
		Height(m.innerH).
		Render(listContent)

	previewW := m.width - m.listW - 4
	rightBorder := panelStyle
	if m.active == panePreview {
		rightBorder = activePanelStyle
	}
	rightPanel := rightBorder.
		Width(previewW - 2).
		Height(m.innerH).
		Render(m.preview.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// ── Feedback line (dynamic) ──────────────────────────────────────────────
	dirtyMarker := ""
	if m.dirty {
		dirtyMarker = dirtyStyle.Render(" [modified]")
	}

	feedback := statusStyle.Render(m.statusMsg) + dirtyMarker
	if m.quitting == quitAsking {
		feedback = dirtyStyle.Render(" " + m.statusMsg)
	}

	// ── Hint line (always visible) ────────────────────────────────────────────
	var hintText string
	switch m.active {
	case panePreview:
		hintText = "[Tab]/[Esc] back to list • [ctrl+s] save"
	default:
		if it := m.list.SelectedItem(); it != nil && it.Existing {
			hintText = "[↑/↓] navigate • [Space] guided edit • [e] free edit • [d] delete • [Tab] edit YAML • [ctrl+s] save • [q] quit"
		} else {
			hintText = "[↑/↓] navigate • [Space] guided add • [e] free add • [Tab] edit YAML • [ctrl+s] save • [q] quit"
		}
	}
	hint := statusStyle.Render(hintText)

	feedbackLine := lipgloss.NewStyle().Width(m.width).Render(feedback)
	hintLine := lipgloss.NewStyle().Width(m.width).Render(hint)

	return strings.Join([]string{body, feedbackLine, hintLine}, "\n")
}
