package edit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/presets"
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

	// Preset state.
	currentPreset string             // "base" by default; "custom" in edit mode
	presetPicker  *PresetPickerModel // non-nil while the picker popover is open
}

// NewOverlay builds an overlay for the given key.
// All overlays use two-panel mode: left panel shows field toggles (or a
// "(no sub-fields)" hint for simple fields), right panel is the YAML editor.
func NewOverlay(key, initialContent string, totalW, totalH int) OverlayModel {
	defs := FieldsForKey(key)

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
		key:           key,
		twoPanel:      true,
		totalW:        totalW,
		totalH:        totalH,
		currentPreset: "custom",
	}

	// Prefer the "base" preset when opening a new (empty) block.
	trivial := key + ":\n"
	if initialContent == "" || initialContent == trivial {
		if y, err := presets.PresetYAML(key, "base"); err == nil {
			initialContent = y
			om.currentPreset = "base"
		} else {
			// Fall back to the static template so the textarea is never empty.
			initialContent = Template(key)
		}
	}

	om.initTwoPanel(defs, contentW, panelH, initialContent)
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

	// Start in the YAML panel when there are no field toggles to interact with.
	if len(defs) == 0 {
		om.active = overlayPanelYAML
		om.yamlEditor.Focus()
	} else {
		om.active = overlayPanelFields
	}

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

func (om OverlayModel) Init() tea.Cmd { return textarea.Blink }

func (om OverlayModel) Update(msg tea.Msg) (OverlayModel, tea.Cmd) {
	// Preset messages are handled at the top level, regardless of picker state.
	switch m := msg.(type) {
	case PresetSelectedMsg:
		return om.applyPreset(m.Name), nil
	case PresetPickerCancelledMsg:
		om.presetPicker = nil
		return om, nil
	}

	// Picker owns key input while open.
	if om.presetPicker != nil {
		if key, ok := msg.(tea.KeyMsg); ok {
			updated, cmd := om.presetPicker.Update(key)
			om.presetPicker = &updated
			return om, cmd
		}
		return om, nil
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return om.updateYAMLEditor(msg)
	}
	return om.updateKey(key)
}

func (om OverlayModel) updateKey(msg tea.KeyMsg) (OverlayModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		return om, func() tea.Msg { return OverlayCancelledMsg{} }
	case tea.KeyCtrlS:
		return om.confirm()
	case tea.KeyTab:
		return om.switchPanel(), nil
	}

	// Left panel: p opens picker, other keys navigate field toggles.
	if om.active == overlayPanelFields {
		if msg.String() == "p" {
			names := presets.ListPresets(om.key)
			if len(names) > 0 {
				picker := NewPresetPicker(names, om.currentPreset, om.totalW, om.totalH)
				om.presetPicker = &picker
			}
			return om, nil
		}
		return om.updateFieldPanel(msg), nil
	}

	return om.updateYAMLEditor(msg)
}

// updateYAMLEditor forwards a message to the textarea when the YAML panel is
// active. Field-panel mode has no use for textarea ticks.
func (om OverlayModel) updateYAMLEditor(msg tea.Msg) (OverlayModel, tea.Cmd) {
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

// applyPreset replaces the textarea content with the selected preset and
// closes the picker. In two-panel mode, refreshes the field-toggle state
// from the new YAML.
func (om OverlayModel) applyPreset(name string) OverlayModel {
	y, err := presets.PresetYAML(om.key, name)
	if err != nil {
		om.errMsg = fmt.Sprintf("preset error: %v", err)
		om.presetPicker = nil
		return om
	}
	om.yamlEditor.SetValue(y)
	om.currentPreset = name
	om.errMsg = ""
	if om.twoPanel {
		om.fieldList.SetFields(syncFieldsFromYAML(om.key, om.fieldList.Fields(), y))
	}
	om.presetPicker = nil
	return om
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
	titleText := fmt.Sprintf(" %s [%s · preset: %s] ", om.key, action, om.currentPreset)
	title := overlayTitleStyle.Render(titleText)

	content := om.viewTwoPanel()

	var hintText string
	if om.active == overlayPanelFields {
		if len(om.fieldList.Fields()) > 0 {
			hintText = "[Tab] switch panel • [Space] toggle • [p] preset • [ctrl+s] apply • [Esc] cancel"
		} else {
			hintText = "[Tab] switch panel • [p] preset • [ctrl+s] apply • [Esc] cancel"
		}
	} else {
		hintText = "[Tab] switch panel • [ctrl+s] apply • [Esc] cancel"
	}
	hint := statusStyle.Render(hintText)

	parts := []string{title, content}
	if om.errMsg != "" {
		parts = append(parts, lipgloss.NewStyle().Foreground(theme.Danger).Render(om.errMsg))
	}
	parts = append(parts, hint)

	box := overlayBorderStyle.Render(strings.Join(parts, "\n"))

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
	overlay := lipgloss.NewStyle().PaddingLeft(lp).PaddingTop(tp).Render(box)

	// Layer the picker over the overlay if open.
	if om.presetPicker != nil {
		return om.presetPicker.View()
	}
	return overlay
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
