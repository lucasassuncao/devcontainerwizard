package edit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
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

type fieldState struct {
	Def     FieldDef
	Checked bool
}

// OverlayModel is the floating overlay for adding a YAML block.
//
// Two-panel mode (guided + complex block): left field-toggle list + right YAML editor.
// Single mode  (free  or simple block):   just the YAML textarea.
type OverlayModel struct {
	key      string
	guided   bool
	twoPanel bool

	// Left panel — two-panel mode only
	fields      []fieldState
	fieldCursor int
	fieldOffset int
	fieldPanelW int // content width passed to panelStyle.Width()
	fieldPanelH int // content height passed to panelStyle.Height()

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
//   - guided=true + block has field defs → two-panel mode
//   - guided=true + simple block         → single textarea with guided template
//   - guided=false                        → single blank textarea
func NewOverlay(key, initialContent string, guided bool, totalW, totalH int) OverlayModel {
	defs := FieldsForKey(key)
	twoPanel := guided && len(defs) > 0

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
		guided:   guided,
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

	fields := make([]fieldState, len(defs))
	for i, d := range defs {
		fields[i] = fieldState{Def: d, Checked: d.Required}
	}

	ta := textarea.New()
	ta.SetWidth(yamlPanelW - 2) // 1-char margin on each side inside the panel
	ta.SetHeight(panelH - 1)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Blur()

	om.fields = fields
	om.fieldPanelW = fieldPanelW
	om.fieldPanelH = panelH
	om.yamlPanelW = yamlPanelW
	om.yamlPanelH = panelH
	om.yamlEditor = ta
	om.active = overlayPanelFields

	// When editing an existing block, seed the textarea with the current
	// content and derive toggle states from it; otherwise build from defaults.
	trivial := om.key + ":\n"
	if initialContent != "" && initialContent != trivial {
		om.yamlEditor.SetValue(strings.ReplaceAll(initialContent, "\r\n", "\n"))
		om.syncFieldsFromYAML()
	} else {
		om.rebuildYAML()
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

// syncFieldsFromYAML parses the current textarea value and updates Checked
// on each field to reflect what is actually present in the YAML.
func (om *OverlayModel) syncFieldsFromYAML() {
	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(om.yamlEditor.Value()), &doc); err != nil {
		return
	}
	sub, _ := doc[om.key].(map[string]interface{})
	for i := range om.fields {
		_, om.fields[i].Checked = sub[om.fields[i].Def.Key]
	}
}

// rebuildYAML reconstructs the textarea value from the checked field states
// in canonical definition order. Any previous manual edits are overwritten.
func (om *OverlayModel) rebuildYAML() {
	var sb strings.Builder
	sb.WriteString(om.key + ":\n")
	for _, fs := range om.fields {
		if fs.Checked {
			sb.WriteString(fs.Def.YAML)
		}
	}
	om.yamlEditor.SetValue(sb.String())
	om.errMsg = ""
}

func removeFieldNode(valueNode *yaml.Node, idx int) {
	if idx >= 0 {
		valueNode.Content = append(valueNode.Content[:idx], valueNode.Content[idx+2:]...)
	}
}

func addFieldNode(valueNode *yaml.Node, idx int, parentKey string, def FieldDef) {
	if idx >= 0 {
		return // already present
	}
	var templateRoot yaml.Node
	if err := yaml.Unmarshal([]byte(parentKey+":\n"+def.YAML), &templateRoot); err != nil {
		return
	}
	if templateRoot.Kind == 0 || len(templateRoot.Content) == 0 {
		return
	}
	tMapping := templateRoot.Content[0]
	if tMapping.Kind != yaml.MappingNode || len(tMapping.Content) < 2 {
		return
	}
	tValue := tMapping.Content[1]
	if tValue.Kind == yaml.MappingNode && len(tValue.Content) >= 2 {
		valueNode.Content = append(valueNode.Content, tValue.Content[0], tValue.Content[1])
	}
}

// applyFieldToggle surgically adds or removes a single sub-field from the
// current yamlEditor value, preserving any edits the user made to other fields.
// Used when isEdit=true so that existing content is not overwritten by defaults.
func (om *OverlayModel) applyFieldToggle(def FieldDef) {
	current := om.yamlEditor.Value()

	var root yaml.Node
	if err := yaml.Unmarshal([]byte(current), &root); err != nil || root.Kind == 0 || len(root.Content) == 0 {
		om.rebuildYAML()
		return
	}
	mapping := root.Content[0]
	if mapping.Kind != yaml.MappingNode || len(mapping.Content) < 2 {
		om.rebuildYAML()
		return
	}
	// mapping.Content[0] is the top-level key node, mapping.Content[1] is the value.
	valueNode := mapping.Content[1]
	if valueNode.Kind != yaml.MappingNode {
		om.rebuildYAML()
		return
	}

	// Find whether the sub-field already exists in the value mapping.
	idx := -1
	for i := 0; i < len(valueNode.Content)-1; i += 2 {
		if valueNode.Content[i].Value == def.Key {
			idx = i
			break
		}
	}

	checked := false
	for _, fs := range om.fields {
		if fs.Def.Key == def.Key {
			checked = fs.Checked
			break
		}
	}

	if !checked {
		removeFieldNode(valueNode, idx)
	} else {
		addFieldNode(valueNode, idx, om.key, def)
	}

	var buf strings.Builder
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&root); err != nil {
		om.rebuildYAML()
		return
	}
	om.yamlEditor.SetValue(strings.TrimRight(buf.String(), "\n") + "\n")
	om.errMsg = ""
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
			om.syncFieldsFromYAML()
		}
		return om, cmd
	}
	return om, nil
}

func (om OverlayModel) confirm() (OverlayModel, tea.Cmd) {
	om.errMsg = ""
	snippet := om.yamlEditor.Value()
	if err := ValidateSnippet(snippet); err != nil {
		om.errMsg = fmt.Sprintf("YAML inválido: %v", err)
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
	n := len(om.fields)
	switch msg.String() {
	case "up", "k":
		if om.fieldCursor > 0 {
			om.fieldCursor--
			if om.fieldCursor < om.fieldOffset {
				om.fieldOffset = om.fieldCursor
			}
		}
	case "down", "j":
		if om.fieldCursor < n-1 {
			om.fieldCursor++
			if om.fieldCursor >= om.fieldOffset+om.fieldPanelH {
				om.fieldOffset = om.fieldCursor - om.fieldPanelH + 1
			}
		}
	case " ":
		if om.fieldCursor < n {
			om.fields[om.fieldCursor].Checked = !om.fields[om.fieldCursor].Checked
			if om.isEdit {
				om.applyFieldToggle(om.fields[om.fieldCursor].Def)
			} else {
				om.rebuildYAML()
			}
		}
	}
	return om
}

// ── View ──────────────────────────────────────────────────────────────────────

func (om OverlayModel) View() string {
	mode := "free"
	if om.guided {
		mode = "guided"
	}
	if om.isEdit {
		mode = "edit"
	}
	title := overlayTitleStyle.Render(fmt.Sprintf(" %s [%s] ", om.key, mode))

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
		parts = append(parts, lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(om.errMsg))
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
		Render(om.renderFieldList())

	rightPanel := rightBorder.
		Width(om.yamlPanelW).
		Height(om.yamlPanelH).
		Render(om.yamlEditor.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

func (om OverlayModel) renderFieldList() string {
	if len(om.fields) == 0 {
		return availableItemStyle.Render("  (sem sub-campos)")
	}

	var sb strings.Builder
	end := om.fieldOffset + om.fieldPanelH
	if end > len(om.fields) {
		end = len(om.fields)
	}

	for i := om.fieldOffset; i < end; i++ {
		fs := om.fields[i]
		mark := "○"
		if fs.Checked {
			mark = "●"
		}
		req := ""
		if fs.Def.Required {
			req = " *"
		}
		label := fmt.Sprintf("%s %-16s%s", mark, fs.Def.Key, req)

		var line string
		switch {
		case i == om.fieldCursor && om.active == overlayPanelFields:
			line = selectedItemStyle.Render("▶ " + label)
		case fs.Checked:
			line = existingItemStyle.Render("  " + label)
		default:
			line = availableItemStyle.Render("  " + label)
		}
		sb.WriteString(line + "\n")
	}

	// Pad remaining rows so the panel height stays stable when few fields exist.
	rendered := end - om.fieldOffset
	for i := rendered; i < om.fieldPanelH; i++ {
		sb.WriteByte('\n')
	}

	return sb.String()
}
