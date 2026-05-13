package edit

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type fieldState struct {
	Def     FieldDef
	Checked bool
}

// FieldListModel is the left panel of the guided overlay: a scrollable,
// toggleable list of sub-field definitions for a given YAML key.
type FieldListModel struct {
	fields []fieldState
	cursor int
	offset int
	width  int
	height int
}

// NewFieldListModel initialises the model from a slice of FieldDefs.
// Required fields start checked.
func NewFieldListModel(defs []FieldDef, w, h int) FieldListModel {
	fields := make([]fieldState, len(defs))
	for i, d := range defs {
		fields[i] = fieldState{Def: d, Checked: d.Required}
	}
	return FieldListModel{fields: fields, width: w, height: h}
}

// SetSize updates the visible panel dimensions.
func (fl *FieldListModel) SetSize(w, h int) {
	fl.width = w
	fl.height = h
}

// Fields returns the current field states.
func (fl FieldListModel) Fields() []fieldState {
	return fl.fields
}

// SetFields replaces the field states (used after syncFieldsFromYAML).
func (fl *FieldListModel) SetFields(fields []fieldState) {
	fl.fields = fields
}

// Update handles keyboard input. Returns the updated model and a bool
// indicating whether a toggle occurred (Space was pressed on a field).
func (fl FieldListModel) Update(msg tea.KeyMsg) (FieldListModel, bool) {
	n := len(fl.fields)
	switch msg.String() {
	case "up", "k":
		if fl.cursor > 0 {
			fl.cursor--
			if fl.cursor < fl.offset {
				fl.offset = fl.cursor
			}
		}
	case "down", "j":
		if fl.cursor < n-1 {
			fl.cursor++
			if fl.cursor >= fl.offset+fl.height {
				fl.offset = fl.cursor - fl.height + 1
			}
		}
	case " ":
		if fl.cursor < n {
			fl.fields[fl.cursor].Checked = !fl.fields[fl.cursor].Checked
			return fl, true
		}
	}
	return fl, false
}

// ToggledField returns the field at the cursor. Call after Update returns toggled=true.
func (fl FieldListModel) ToggledField() fieldState {
	return fl.fields[fl.cursor]
}

// View renders the visible slice of the field list.
func (fl FieldListModel) View() string {
	if len(fl.fields) == 0 {
		return availableItemStyle.Render("  (no sub-fields)")
	}

	var sb strings.Builder
	end := fl.offset + fl.height
	if end > len(fl.fields) {
		end = len(fl.fields)
	}

	for i := fl.offset; i < end; i++ {
		fs := fl.fields[i]
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
		case i == fl.cursor:
			line = selectedItemStyle.Render("▶ " + label)
		case fs.Checked:
			line = existingItemStyle.Render("  " + label)
		default:
			line = availableItemStyle.Render("  " + label)
		}
		sb.WriteString(line + "\n")
	}

	// Pad remaining rows so the panel height stays stable.
	rendered := end - fl.offset
	for i := rendered; i < fl.height; i++ {
		sb.WriteByte('\n')
	}

	return sb.String()
}
