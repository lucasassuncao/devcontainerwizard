package edit

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane int

const (
	paneList pane = iota
	panePreview
	paneOverlay
	paneAlert
)

type quitState int

const (
	quitNone quitState = iota
	quitAsking
)

// Model is the root Bubble Tea model for the edit TUI.
//
// The active pane is derived from state, not tracked explicitly:
//   - alert != nil       → paneAlert
//   - overlay != nil     → paneOverlay
//   - previewFocused     → panePreview
//   - otherwise          → paneList
type Model struct {
	filePath string

	rawYAML []byte
	blocks  []Block

	list    ListModel
	preview textarea.Model
	overlay *OverlayModel
	alert   *AlertModel

	previewFocused bool
	dirty          bool
	quitting       quitState
	statusMsg      string

	width  int
	height int
	listW  int // derived by relayout(); read by View()
	innerH int // derived by relayout(); read by View()
}

// activePane reports which pane currently owns input/rendering. Derived from
// state so the four indicators can never disagree.
func (m Model) activePane() pane {
	switch {
	case m.alert != nil:
		return paneAlert
	case m.overlay != nil:
		return paneOverlay
	case m.previewFocused:
		return panePreview
	default:
		return paneList
	}
}

// scrollPreviewToKey moves the preview cursor to the line where key starts.
// A no-op if the key is empty or not present.
func (m *Model) scrollPreviewToKey(key string) {
	if key == "" {
		return
	}
	target := key + ":"
	for i, l := range strings.Split(string(m.rawYAML), "\n") {
		if strings.HasPrefix(l, target) {
			m.preview.SetCursor(i)
			return
		}
	}
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
	preview := textarea.New()
	preview.CharLimit = 0
	preview.ShowLineNumbers = false
	preview.Blur()
	preview.SetValue(strings.ReplaceAll(string(raw), "\r\n", "\n"))

	return Model{
		filePath:  filePath,
		rawYAML:   raw,
		blocks:    blocks,
		list:      list,
		preview:   preview,
		statusMsg: "",
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Root-level messages handled regardless of the active pane.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()
		return m, nil
	case SpaceOnItemMsg:
		return m.handleSpace(msg.Item)
	case OverlayConfirmedMsg:
		return m.handleOverlayConfirmed(msg.Snippet)
	case OverlayCancelledMsg:
		m.overlay = nil
		m.statusMsg = "Cancelled."
		return m, nil
	case DeleteItemMsg:
		return m.handleDelete(msg.Key)
	case AlertDismissedMsg:
		m.alert = nil
		return m, nil
	}

	// 2. Quit-confirm prompt intercepts keys before pane dispatch.
	if key, ok := msg.(tea.KeyMsg); ok && m.quitting == quitAsking {
		if s := key.String(); s == "y" || s == "Y" {
			return m, tea.Quit
		}
		m.quitting = quitNone
		m.statusMsg = "Quit cancelled."
		return m, nil
	}

	// 3. Delegate to the active pane.
	switch m.activePane() {
	case paneAlert:
		if key, ok := msg.(tea.KeyMsg); ok {
			al, cmd := m.alert.Update(key)
			m.alert = &al
			return m, cmd
		}
	case paneOverlay:
		ov, cmd := m.overlay.Update(msg)
		m.overlay = &ov
		return m, cmd
	case panePreview:
		if key, ok := msg.(tea.KeyMsg); ok {
			return m.handlePreviewKey(key)
		}
		var cmd tea.Cmd
		m.preview, cmd = m.preview.Update(msg)
		return m, cmd
	case paneList:
		if key, ok := msg.(tea.KeyMsg); ok {
			return m.handleListKey(key)
		}
	}
	return m, nil
}

// handleListKey processes keys while the list pane has focus.
func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		return m.save()
	case "ctrl+l":
		return m.validateKeys()
	case "tab":
		return m.togglePreviewPane()
	case "q", "ctrl+c":
		if m.dirty {
			m.quitting = quitAsking
			m.statusMsg = "Unsaved changes. Quit without saving? (y/N)"
			return m, nil
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if it := m.list.SelectedItem(); it != nil {
		m.scrollPreviewToKey(it.Key)
	}
	return m, cmd
}

// handlePreviewKey processes keys while the preview pane has focus.
// q/ctrl+c are NOT quit shortcuts here — they go to the textarea as input.
func (m Model) handlePreviewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		return m.save()
	case "ctrl+l":
		return m.validateKeys()
	case "tab", "esc":
		return m.togglePreviewPane()
	}
	return m.updatePreviewEditor(msg)
}

func (m Model) togglePreviewPane() (tea.Model, tea.Cmd) {
	if m.previewFocused {
		m.previewFocused = false
		m.preview.Blur()
		m.statusMsg = ""
		return m, nil
	}
	m.previewFocused = true
	cmd := m.preview.Focus()
	m.statusMsg = "Editing YAML directly — Tab/Esc back to list."
	return m, cmd
}

func (m Model) updatePreviewEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.preview, cmd = m.preview.Update(msg)
	raw := []byte(m.preview.Value())
	if blocks, err := ParseBlocksFromBytes(raw); err == nil {
		m.rawYAML = raw
		m.blocks = blocks
		m.list.Rebuild(m.blocks)
		m.dirty = true
	}
	return m, cmd
}

func (m Model) handleSpace(it ListItem) (tea.Model, tea.Cmd) {
	var initial string
	if it.Existing {
		current, err := BlockContent(m.rawYAML, m.blocks, it.Key)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Error reading %s: %v", it.Key, err)
			return m, nil
		}
		initial = current
	} else {
		initial = it.Key + ":\n"
		if len(FieldsForKey(it.Key)) == 0 {
			// Only use the template for single-textarea blocks (no field defs).
			// Two-panel blocks initialise from rebuildYAML() when content is trivial.
			initial = Template(it.Key)
		}
	}

	ov := NewOverlay(it.Key, initial, m.width, m.height)
	if it.Existing {
		ov.isEdit = true
		ov.editKey = it.Key
	}
	m.overlay = &ov
	if it.Existing {
		m.statusMsg = fmt.Sprintf("Edit block %q — Tab panel, Ctrl+S confirm, Esc cancel.", it.Key)
	} else {
		m.statusMsg = fmt.Sprintf("Add block %q — Tab panel, Ctrl+S confirm, Esc cancel.", it.Key)
	}
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
		return m, nil
	}
	m.applyRaw(newRaw)
	m.overlay = nil
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
	m.preview.SetValue(strings.ReplaceAll(string(raw), "\r\n", "\n"))
	if it := m.list.SelectedItem(); it != nil {
		m.scrollPreviewToKey(it.Key)
	}
	m.dirty = true
}

func (m Model) save() (tea.Model, tea.Cmd) {
	if unknown := ValidateKnownKeys(m.rawYAML); len(unknown) > 0 {
		msg := "The following key(s) are not valid devcontainer fields:\n\n  " +
			strings.Join(unknown, ", ") +
			"\n\nFix the YAML before saving."
		return m.showAlert("Invalid key(s) — cannot save", msg, alertError)
	}
	if err := os.WriteFile(m.filePath, m.rawYAML, 0o600); err != nil {
		return m.showAlert("Save failed", err.Error(), alertError)
	}
	m.dirty = false
	m.statusMsg = fmt.Sprintf("Saved to %s.", m.filePath)
	return m, nil
}

func (m Model) validateKeys() (tea.Model, tea.Cmd) {
	if unknown := ValidateKnownKeys(m.rawYAML); len(unknown) > 0 {
		msg := "The following key(s) are not valid devcontainer fields:\n\n  " +
			strings.Join(unknown, ", ")
		return m.showAlert("Validation — unknown key(s)", msg, alertError)
	}
	return m.showAlert("Validation passed", "All keys are recognized devcontainer fields.", alertSuccess)
}

func (m Model) showAlert(title, message string, kind alertKind) (tea.Model, tea.Cmd) {
	al := NewAlert(title, message, kind, m.width, m.height)
	m.alert = &al
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
	m.preview.SetWidth(previewW - 2)
	m.preview.SetHeight(m.innerH)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.activePane() {
	case paneAlert:
		return m.alert.View()
	case paneOverlay:
		return m.overlay.View()
	}

	listContent := m.list.View()
	listBorder := panelStyle
	if !m.previewFocused {
		listBorder = activePanelStyle
	}
	leftPanel := listBorder.
		Width(m.listW - 2).
		Height(m.innerH).
		Render(listContent)

	previewW := m.width - m.listW - 4
	rightBorder := panelStyle
	if m.previewFocused {
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
	if m.previewFocused {
		hintText = "[Tab]/[Esc] back to list • [ctrl+l] validate • [ctrl+s] save"
	} else if it := m.list.SelectedItem(); it != nil && it.Existing {
		hintText = "[↑/↓] navigate • [Space] edit block • [d] delete • [Tab] edit YAML • [ctrl+l] validate • [ctrl+s] save • [q] quit"
	} else {
		hintText = "[↑/↓] navigate • [Space] add block • [Tab] edit YAML • [ctrl+l] validate • [ctrl+s] save • [q] quit"
	}
	hint := statusStyle.Render(hintText)

	feedbackLine := lipgloss.NewStyle().Width(m.width).Render(feedback)
	hintLine := lipgloss.NewStyle().Width(m.width).Render(hint)

	return strings.Join([]string{body, feedbackLine, hintLine}, "\n")
}
