package edit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

// OverlayConfirmedMsg is sent when the user confirms with Ctrl+S.
type OverlayConfirmedMsg struct{ Snippet string }

// OverlayCancelledMsg is sent when the user presses Esc.
type OverlayCancelledMsg struct{}

type overlayPanel int

const (
	overlayPanelFields overlayPanel = iota
	overlayPanelYAML
)

// OverlayModel is the floating overlay for adding or editing a YAML block.
//
// Two-panel mode (complex block): left field-toggle list + right YAML editor.
// Single mode (simple block):     just the YAML textarea.
type OverlayModel struct {
	key      string
	twoPanel bool

	// Left panel — two-panel mode only
	fieldList   FieldListModel
	fieldPanelW int // border sizing for viewTwoPanel
	fieldPanelH int // border sizing for viewTwoPanel

	// Right / only YAML panel
	yamlEditor textarea.Model
	yamlPanelW int
	yamlPanelH int

	active overlayPanel
	errMsg string

	isEdit  bool
	editKey string

	totalW int
	totalH int
}

// NewOverlay builds an overlay for the given key.
// Keys with field definitions open in two-panel mode; all others use a single textarea.
func NewOverlay(key, initialContent string, totalW, totalH int) OverlayModel {
	defs := FieldsForKey(key)
	twoPanel := len(defs) > 0

	// ── Outer box dimensions (including double border + padding) ──────────────
	//
	// overlayBorderStyle has Border(DoubleBorder()) + Padding(0,1).
	// In lipgloss, Width(n) means: n = content + padding (border is extra).
	// Total rendered width  = Width(n) + border-left(1) + border-right(1).
	// Total rendered height = Height(n) + border-top(1) + border-bottom(1).
	//
	// We calculate sizes from the OUTSIDE in:
	//   boxW    = desired total outer width  (border + padding + content)
	//   contentW = boxW  - border(2) - padding(2) = boxW  - 4   ← available content width
	//   boxH    = desired total outer height (border + content, padding top/bot = 0)
	//   contentH = boxH  - border(2)                             ← available content lines
	boxW := totalW - 4
	boxH := totalH - 4
	if boxW > 120 {
		boxW = 120
	}
	if boxH > 36 {
		boxH = 36
	}
	if boxW < 60 {
		boxW = 60
	}
	if boxH < 16 {
		boxH = 16
	}

	// Content area inside the outer border+padding.
	contentW := boxW - 4 // border L+R (2) + padding L+R (2)

	// The rendered height occupied by fixed rows inside the box:
	//   title (1) + sep-\n (1) + panels + sep-\n (1) + hint (1) = 4 fixed rows
	//   + panel border top+bot (2) = 6 rows overhead
	// → panelH (content rows inside panel) = boxH - border(2) - 6 overhead
	panelH := boxH - 8
	if panelH < 4 {
		panelH = 4
	}

	om := OverlayModel{
		key:      key,
		twoPanel: twoPanel,
		totalW:   totalW,
		totalH:   totalH,
	}

	if twoPanel {
		om.initTwoPanel(defs, contentW, panelH, initialContent)
	} else {
		om.initSinglePanel(contentW, panelH, initialContent)
	}

	return om
}

func (om *OverlayModel) initTwoPanel(defs []FieldDef, contentW, panelH int, initialContent string) {
	// Each inner panel has a rounded border (1 char each side = 2 overhead).
	// Total panels rendered width = fieldPanelW+2 + yamlPanelW+2 = contentW
	// → fieldPanelW + yamlPanelW = contentW - 4
	panelSpace := contentW - 4
	if panelSpace < 24 {
		panelSpace = 24
	}
	fieldPanelW := panelSpace / 3
	if fieldPanelW < 18 {
		fieldPanelW = 18
	}
	yamlPanelW := panelSpace - fieldPanelW

	om.fieldList = NewFieldListModel(defs, panelH)
	om.fieldPanelW = fieldPanelW
	om.fieldPanelH = panelH

	ta := textarea.New()
	ta.SetWidth(yamlPanelW - 2) // 1-char margin on each side inside the panel
	ta.SetHeight(panelH - 1)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Blur()

	om.yamlPanelW = yamlPanelW
	om.yamlPanelH = panelH
	om.yamlEditor = ta
	om.active = overlayPanelFields

	// When editing an existing block, seed the textarea with the current
	// content and derive toggle states from it; otherwise build from defaults.
	trivial := om.key + ":\n"
	if initialContent != "" && initialContent != trivial {
		om.yamlEditor.SetValue(strings.ReplaceAll(initialContent, "\r\n", "\n"))
		om.fieldList.SetFields(syncFieldsFromYAML(om.key, om.fieldList.Fields(), om.yamlEditor.Value()))
	} else {
		om.yamlEditor.SetValue(rebuildYAML(om.key, om.fieldList.Fields()))
		om.errMsg = ""
	}
}

func (om *OverlayModel) initSinglePanel(contentW, panelH int, initialContent string) {
	ta := textarea.New()
	ta.SetWidth(contentW - 2) // small margin
	ta.SetHeight(panelH)
	ta.Placeholder = fmt.Sprintf("%s:\n  # your YAML here", om.key)
	ta.SetValue(initialContent)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = true

	om.yamlEditor = ta
	om.yamlPanelW = contentW - 2
	om.yamlPanelH = panelH
	om.active = overlayPanelYAML
}

func (om OverlayModel) Init() tea.Cmd { return textarea.Blink }

func (om OverlayModel) Update(msg tea.Msg) (OverlayModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Global shortcuts — work from any panel.
		switch msg.Type {
		case tea.KeyEsc:
			return om, func() tea.Msg { return OverlayCancelledMsg{} }
		case tea.KeyCtrlS:
			return om.confirm()
		case tea.KeyTab:
			if om.twoPanel {
				return om.switchPanel(), nil
			}
		}

		// Left-panel navigation when in two-panel mode.
		if om.twoPanel && om.active == overlayPanelFields {
			return om.updateFieldPanel(msg), nil
		}
	}

	// Forward non-key messages (e.g. textarea.Blink) only when the YAML panel is
	// active; field-panel mode has no use for textarea ticks.
	if !om.twoPanel || om.active == overlayPanelYAML {
		var cmd tea.Cmd
		om.yamlEditor, cmd = om.yamlEditor.Update(msg)
		if om.twoPanel {
			om.fieldList.SetFields(syncFieldsFromYAML(om.key, om.fieldList.Fields(), om.yamlEditor.Value()))
		}
		return om, cmd
	}
	return om, nil
}

func (om OverlayModel) confirm() (OverlayModel, tea.Cmd) {
	om.errMsg = ""
	snippet := om.yamlEditor.Value()
	if err := ValidateSnippet(snippet); err != nil {
		om.errMsg = fmt.Sprintf("Invalid YAML: %v", err)
		return om, nil
	}
	if !strings.HasSuffix(snippet, "\n") {
		snippet += "\n"
	}
	return om, func() tea.Msg { return OverlayConfirmedMsg{Snippet: snippet} }
}

func (om OverlayModel) switchPanel() OverlayModel {
	if om.active == overlayPanelFields {
		om.active = overlayPanelYAML
		om.yamlEditor.Focus()
	} else {
		om.active = overlayPanelFields
		om.yamlEditor.Blur()
	}
	return om
}

func (om OverlayModel) updateFieldPanel(msg tea.KeyMsg) OverlayModel {
	updated, toggled := om.fieldList.Update(msg)
	om.fieldList = updated
	if toggled {
		fs := om.fieldList.ToggledField()
		if om.isEdit {
			val := applyFieldToggle(om.key, om.fieldList.Fields(), fs.Def, fs.Checked, om.yamlEditor.Value())
			om.yamlEditor.SetValue(val)
			om.errMsg = ""
		} else {
			om.yamlEditor.SetValue(rebuildYAML(om.key, om.fieldList.Fields()))
			om.errMsg = ""
		}
	}
	return om
}

// ── View ──────────────────────────────────────────────────────────────────────

func (om OverlayModel) View() string {
	action := "add block"
	if om.isEdit {
		action = "edit block"
	}
	title := overlayTitleStyle.Render(fmt.Sprintf(" %s [%s] ", om.key, action))

	var content string
	if om.twoPanel {
		content = om.viewTwoPanel()
	} else {
		content = om.yamlEditor.View()
	}

	hint := statusStyle.Render("[Tab] switch panel • [Space] toggle • [ctrl+s] confirm • [Esc] cancel")
	if !om.twoPanel {
		hint = statusStyle.Render("[ctrl+s] confirm • [Esc] cancel")
	}

	parts := []string{title, content}
	if om.errMsg != "" {
		parts = append(parts, lipgloss.NewStyle().Foreground(theme.Danger).Render(om.errMsg))
	}
	parts = append(parts, hint)

	// Let the box auto-size to the content — no explicit Width() to avoid the
	// lipgloss Width-includes-padding gotcha that caused layout overflow.
	box := overlayBorderStyle.Render(strings.Join(parts, "\n"))

	// Centre within the terminal.
	bw := lipgloss.Width(box)
	bh := lipgloss.Height(box)
	lp := (om.totalW - bw) / 2
	tp := (om.totalH - bh) / 2
	if lp < 0 {
		lp = 0
	}
	if tp < 0 {
		tp = 0
	}
	return lipgloss.NewStyle().PaddingLeft(lp).PaddingTop(tp).Render(box)
}

func (om OverlayModel) viewTwoPanel() string {
	leftBorder := panelStyle
	rightBorder := panelStyle
	if om.active == overlayPanelFields {
		leftBorder = activePanelStyle
	} else {
		rightBorder = activePanelStyle
	}

	leftPanel := leftBorder.
		Width(om.fieldPanelW).
		Height(om.fieldPanelH).
		Render(om.fieldList.View())

	rightPanel := rightBorder.
		Width(om.yamlPanelW).
		Height(om.yamlPanelH).
		Render(om.yamlEditor.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}
