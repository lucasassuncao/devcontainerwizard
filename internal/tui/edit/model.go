package edit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

type pane int

const (
	paneList pane = iota
	panePreview
	paneOverlay
	paneAlert
)

// Model is the root Bubble Tea model for the edit TUI.
//
// The active pane is derived from state, not tracked explicitly:
//   - alert != nil       → paneAlert
//   - overlay != nil     → paneOverlay
//   - previewFocused     → panePreview
//   - otherwise          → paneList
//
// All YAML state (raw bytes, parsed blocks, undo history, dirty flag, file
// path) lives on doc. Model is a pure UI orchestrator.
type Model struct {
	doc *Document

	list    ListModel
	preview textarea.Model
	overlay *OverlayModel
	alert   *AlertModel

	previewFocused bool
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
	for i, l := range strings.Split(string(m.doc.Raw()), "\n") {
		if strings.HasPrefix(l, target) {
			m.preview.SetCursor(i)
			return
		}
	}
}

// New loads the YAML file and initialises the model.
func New(filePath string) (Model, error) {
	doc, err := LoadDocument(filePath)
	if err != nil {
		return Model{}, fmt.Errorf("loading %s: %w", filePath, err)
	}

	list := NewListModel(doc.Blocks(), 0)
	preview := textarea.New()
	preview.CharLimit = 0
	preview.ShowLineNumbers = false
	preview.Blur()
	preview.SetValue(string(doc.Raw()))

	return Model{
		doc:       doc,
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
	case doSaveMsg:
		return m.execSave()
	}

	// 2. Delegate to the active pane.
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

// handleGlobalKey handles shortcuts that work in any pane.
// Returns the updated model, a command, and whether the key was consumed.
func (m Model) handleGlobalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+s":
		mo, cmd := m.save()
		return mo, cmd, true
	case "ctrl+l":
		mo, cmd := m.validateKeys()
		return mo, cmd, true
	case "ctrl+z":
		return m.undo(), nil, true
	}
	return m, nil, false
}

// handleListKey processes keys while the list pane has focus.
func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if mo, cmd, handled := m.handleGlobalKey(msg); handled {
		return mo, cmd
	}

	// tab and quit shortcuts are blocked in filter mode to avoid key conflicts.
	if !m.list.IsFiltering() {
		switch msg.String() {
		case "tab":
			return m.togglePreviewPane()
		case "q", "ctrl+c":
			if m.doc.Dirty() {
				return m.showConfirmAlert("Quit without saving?",
					"Unsaved changes will be lost.", tea.Quit)
			}
			return m, tea.Quit
		}
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
	if mo, cmd, handled := m.handleGlobalKey(msg); handled {
		return mo, cmd
	}
	switch msg.String() {
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
	// Best-effort apply. If the draft does not parse, leave the textarea as-is
	// (the user is mid-edit); the document only updates on parse success.
	if err := m.doc.ReplaceRaw([]byte(m.preview.Value())); err == nil {
		m.list.Rebuild(m.doc.Blocks())
	}
	return m, cmd
}

func (m Model) handleSpace(it ListItem) (tea.Model, tea.Cmd) {
	var initial string
	if it.Existing {
		current, err := m.doc.BlockContent(it.Key)
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
	action := "Add"
	if it.Existing {
		ov.isEdit = true
		ov.editKey = it.Key
		action = "Edit"
	}
	m.overlay = &ov
	m.statusMsg = fmt.Sprintf("%s block %q — Tab panel, Ctrl+S confirm, Esc cancel.", action, it.Key)
	return m, nil
}

func (m Model) handleDelete(key string) (tea.Model, tea.Cmd) {
	if err := m.doc.Remove(key); err != nil {
		m.statusMsg = fmt.Sprintf("Error removing %s: %v", key, err)
		return m, nil
	}
	m.syncView()
	m.statusMsg = fmt.Sprintf("Removed %q (not saved yet).", key)
	return m, nil
}

func (m Model) handleOverlayConfirmed(snippet string) (tea.Model, tea.Cmd) {
	isEdit := m.overlay != nil && m.overlay.isEdit
	editKey := ""
	if isEdit {
		editKey = m.overlay.editKey
	}

	var err error
	if isEdit {
		err = m.doc.Replace(editKey, snippet)
	} else {
		err = m.doc.Insert(snippet)
	}
	if err != nil {
		m.statusMsg = fmt.Sprintf("Apply error: %v", err)
		m.overlay = nil
		return m, nil
	}
	m.syncView()
	m.overlay = nil
	if isEdit {
		m.statusMsg = "Block updated (not saved yet)."
	} else {
		m.statusMsg = "Block added (not saved yet)."
	}
	return m, nil
}

// syncView propagates current Document state to the preview textarea and list.
// Call after any mutation through the Document.
func (m *Model) syncView() {
	m.preview.SetValue(string(m.doc.Raw()))
	m.list.Rebuild(m.doc.Blocks())
	if it := m.list.SelectedItem(); it != nil {
		m.scrollPreviewToKey(it.Key)
	}
}

func (m Model) undo() tea.Model {
	if !m.doc.Undo() {
		m.statusMsg = "Nothing to undo."
		return m
	}
	m.syncView()
	m.statusMsg = "Undone."
	return m
}

func formatErrors(errs []string) string {
	var sb strings.Builder
	for i, e := range errs {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString("• ")
		sb.WriteString(e)
	}
	return sb.String()
}

// collectErrors gathers all blocking validation errors for the current document.
// Used by both save() and validateKeys() to avoid divergence.
func (m Model) collectErrors() []string {
	var errs []string
	if u := m.doc.UnknownKeys(); len(u) > 0 {
		errs = append(errs, "Unknown key(s): "+strings.Join(u, ", "))
	}
	return append(errs, m.doc.Conflicts()...)
}

func (m Model) save() (tea.Model, tea.Cmd) {
	if errs := m.collectErrors(); len(errs) > 0 {
		return m.showAlert("Cannot save — fix errors first", formatErrors(errs), alertError)
	}
	doSave := func() tea.Msg { return doSaveMsg{} }
	return m.showConfirmAlert("Save changes?", fmt.Sprintf("Save to %s?", m.doc.Path()), doSave)
}

type doSaveMsg struct{}

func (m Model) execSave() (tea.Model, tea.Cmd) {
	if err := m.doc.Save(); err != nil {
		return m.showAlert("Save failed", err.Error(), alertError)
	}
	return m.showAlert("Saved", fmt.Sprintf("Saved to %s.", m.doc.Path()), alertSuccess)
}

func (m Model) validateKeys() (tea.Model, tea.Cmd) {
	if errs := m.collectErrors(); len(errs) > 0 {
		return m.showAlert("Validation errors", formatErrors(errs), alertError)
	}
	return m.showAlert("Validation passed", "All keys are valid devcontainer fields with no conflicts.", alertSuccess)
}

func (m Model) showAlert(title, message string, kind alertKind) (tea.Model, tea.Cmd) {
	al := NewAlert(title, message, kind, m.width, m.height)
	m.alert = &al
	return m, nil
}

func (m Model) showConfirmAlert(title, message string, confirmCmd tea.Cmd) (tea.Model, tea.Cmd) {
	al := NewConfirmAlert(title, message, confirmCmd, m.width, m.height)
	m.alert = &al
	return m, nil
}

const (
	headerLines    = 1
	statusBarLines = 2 // feedback line + hint line
)

func (m *Model) relayout() {
	var previewW int
	m.listW, previewW = theme.TwoColumnWidths(m.width)
	m.innerH = m.height - headerLines - statusBarLines - 2 // 2 panel borders (top+bottom)
	if m.innerH < 1 {
		m.innerH = 1
	}
	m.list.SetHeight(m.innerH)
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

	header := renderHeader(m.doc.Path(), m.doc.Dirty(), m.width)

	leftTitle := fmt.Sprintf("Blocks (%d/%d)", m.list.AddedCount(), len(allKnownKeys))
	leftPanel := theme.RenderTitledPanel(leftTitle, m.listW, m.innerH+2, !m.previewFocused, m.list.View())

	_, previewW := theme.TwoColumnWidths(m.width)
	rightPanel := theme.RenderTitledPanel("Preview", previewW, m.innerH+2, m.previewFocused, m.preview.View())

	var hintText string
	if m.previewFocused {
		hintText = "[Tab]/[Esc] back to list • [ctrl+l] validate • [ctrl+s] save"
	} else if m.list.IsFiltering() {
		hintText = "[type] filter • [↑/↓] navigate • [Enter] select • [Esc] clear filter"
	} else if it := m.list.SelectedItem(); it != nil && it.Existing {
		hintText = "[↑/↓] navigate • [Space] edit block • [d] delete • [/] filter • [Tab] edit YAML • [ctrl+z] undo • [ctrl+s] save • [q] quit"
	} else {
		hintText = "[↑/↓] navigate • [Space] add block • [/] filter • [Tab] edit YAML • [ctrl+z] undo • [ctrl+s] save • [q] quit"
	}

	feedback := lipgloss.NewStyle().Width(m.width).Render(statusStyle.Render(m.statusMsg))
	hint := lipgloss.NewStyle().Width(m.width).Render(statusStyle.Render(hintText))

	return theme.RenderTwoColumnView(header, leftPanel, rightPanel, feedback, hint)
}
