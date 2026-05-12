package edit_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasassuncao/devcontainerwizard/internal/tui/edit"
)

func TestBuildListItems(t *testing.T) {
	existing := []edit.Block{
		{Key: "name"}, {Key: "image"}, {Key: "features"},
	}
	items := edit.BuildListItems(existing)

	// items[0] is the "── Added ──" separator.
	if !items[0].Separator {
		t.Errorf("items[0] should be the 'Added' separator, got %+v", items[0])
	}
	if items[1].Key != "name" || !items[1].Existing {
		t.Errorf("items[1] = %+v, want {Key:name, Existing:true}", items[1])
	}
	if items[2].Key != "image" || !items[2].Existing {
		t.Errorf("items[2] = %+v, want {Key:image, Existing:true}", items[2])
	}
	if items[3].Key != "features" || !items[3].Existing {
		t.Errorf("items[3] = %+v, want {Key:features, Existing:true}", items[3])
	}

	foundAvailable := false
	for _, it := range items {
		if !it.Existing && !it.Separator {
			foundAvailable = true
			break
		}
	}
	if !foundAvailable {
		t.Error("no available (non-existing) items in list")
	}
}

func TestListModelNavigation(t *testing.T) {
	existing := []edit.Block{{Key: "name"}, {Key: "image"}}
	lm := edit.NewListModel(existing, 10)

	// Should start on first non-separator item.
	it := lm.SelectedItem()
	if it == nil || it.Key != "name" {
		t.Fatalf("expected cursor on 'name', got %v", it)
	}

	// Move down once.
	lm, _ = lm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	it = lm.SelectedItem()
	if it == nil || it.Key != "image" {
		t.Errorf("expected cursor on 'image' after j, got %v", it)
	}
}
