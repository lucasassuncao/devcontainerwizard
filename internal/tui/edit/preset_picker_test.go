package edit

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPresetPickerInitialState(t *testing.T) {
	p := NewPresetPicker([]string{"base", "vscode-go"}, "base", 80, 24)
	if p.SelectedName() != "base" {
		t.Errorf("initial selection = %q, want \"base\"", p.SelectedName())
	}
}

func TestPresetPickerArrowDownMoves(t *testing.T) {
	p := NewPresetPicker([]string{"base", "vscode-go", "vscode-node"}, "base", 80, 24)
	updated, _ := p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if updated.SelectedName() != "vscode-go" {
		t.Errorf("after Down, selection = %q, want \"vscode-go\"", updated.SelectedName())
	}
}

func TestPresetPickerEnterEmitsSelectedMsg(t *testing.T) {
	p := NewPresetPicker([]string{"base", "vscode-go"}, "base", 80, 24)
	_, cmd := p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter should emit a command")
	}
	msg := cmd()
	sel, ok := msg.(PresetSelectedMsg)
	if !ok {
		t.Fatalf("got %T, want PresetSelectedMsg", msg)
	}
	if sel.Name != "base" {
		t.Errorf("selected name = %q, want \"base\"", sel.Name)
	}
}

func TestPresetPickerEscEmitsCancelMsg(t *testing.T) {
	p := NewPresetPicker([]string{"base"}, "base", 80, 24)
	_, cmd := p.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Esc should emit a command")
	}
	if _, ok := cmd().(PresetPickerCancelledMsg); !ok {
		t.Fatalf("got %T, want PresetPickerCancelledMsg", cmd())
	}
}
