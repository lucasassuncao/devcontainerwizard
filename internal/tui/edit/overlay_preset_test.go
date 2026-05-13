package edit

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewOverlayLoadsBasePresetForKnownField(t *testing.T) {
	ov := NewOverlay("customizations", "", 100, 30)
	if !strings.Contains(ov.yamlEditor.Value(), "vscode") {
		t.Errorf("expected base preset YAML to contain 'vscode', got:\n%s", ov.yamlEditor.Value())
	}
}

func TestNewOverlayFallsBackToTemplateWhenPresetMissing(t *testing.T) {
	ov := NewOverlay("$schema", "", 100, 30)
	if ov.yamlEditor.Value() == "" {
		t.Error("expected fallback content, got empty textarea")
	}
}

func TestNewOverlayEditModeKeepsExistingContent(t *testing.T) {
	existing := "customizations:\n  vscode:\n    extensions:\n      - some.ext\n"
	ov := NewOverlay("customizations", existing, 100, 30)
	if !strings.Contains(ov.yamlEditor.Value(), "some.ext") {
		t.Errorf("edit mode should preserve existing content; got:\n%s", ov.yamlEditor.Value())
	}
}

func TestOverlayCurrentPresetDefaultsToBase(t *testing.T) {
	ov := NewOverlay("customizations", "", 100, 30)
	if ov.currentPreset != "base" {
		t.Errorf("currentPreset = %q, want \"base\"", ov.currentPreset)
	}
}

func TestOverlayCurrentPresetCustomInEditMode(t *testing.T) {
	ov := NewOverlay("customizations", "customizations:\n  vscode: {}\n", 100, 30)
	if ov.currentPreset != "custom" {
		t.Errorf("currentPreset (edit mode) = %q, want \"custom\"", ov.currentPreset)
	}
}

func TestPOpensPresetPickerFromLeftPanel(t *testing.T) {
	ov := NewOverlay("customizations", "", 100, 30)
	updated, _ := ov.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if updated.presetPicker == nil {
		t.Fatal("p should open the preset picker when left panel is active")
	}
}

func TestPresetSelectedMsgUpdatesTextarea(t *testing.T) {
	ov := NewOverlay("customizations", "", 100, 30)
	before := ov.yamlEditor.Value()

	updated, _ := ov.Update(PresetSelectedMsg{Name: "vscode-go"})

	after := updated.yamlEditor.Value()
	if before == after {
		t.Error("textarea content should have changed when switching presets")
	}
	if !strings.Contains(after, "golang.go") {
		t.Errorf("expected vscode-go YAML to contain 'golang.go', got:\n%s", after)
	}
	if updated.currentPreset != "vscode-go" {
		t.Errorf("currentPreset = %q, want \"vscode-go\"", updated.currentPreset)
	}
	if updated.presetPicker != nil {
		t.Error("picker should be closed after selection")
	}
}

func TestPresetPickerCancelledKeepsTextarea(t *testing.T) {
	ov := NewOverlay("customizations", "", 100, 30)
	before := ov.yamlEditor.Value()

	withPicker, _ := ov.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	updated, _ := withPicker.Update(PresetPickerCancelledMsg{})

	if updated.yamlEditor.Value() != before {
		t.Error("textarea should be unchanged after picker cancellation")
	}
	if updated.presetPicker != nil {
		t.Error("picker should be closed after cancellation")
	}
}
