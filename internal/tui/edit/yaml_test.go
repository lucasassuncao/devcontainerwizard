package edit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/edit"
)

const sampleYAML = `name: mydev
image: ubuntu:22.04
features:
  ghcr.io/devcontainers/features/git:1: {}
forwardPorts:
  - 3000
`

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseBlocks(t *testing.T) {
	p := writeTmp(t, sampleYAML)
	blocks, err := edit.ParseBlocks(p)
	if err != nil {
		t.Fatalf("ParseBlocks: %v", err)
	}
	if len(blocks) != 4 {
		t.Fatalf("expected 4 blocks, got %d", len(blocks))
	}
	if blocks[0].Key != "name" {
		t.Errorf("first block key = %q, want \"name\"", blocks[0].Key)
	}
	if blocks[2].Key != "features" {
		t.Errorf("third block key = %q, want \"features\"", blocks[2].Key)
	}
}

func TestRemoveBlock(t *testing.T) {
	raw := []byte(sampleYAML)
	blocks, _ := edit.ParseBlocksFromBytes(raw)
	result, err := edit.RemoveBlock(raw, blocks, "image")
	if err != nil {
		t.Fatalf("RemoveBlock: %v", err)
	}
	remaining, _ := edit.ParseBlocksFromBytes(result)
	for _, b := range remaining {
		if b.Key == "image" {
			t.Error("image block still present after removal")
		}
	}
}

func TestInsertBlock(t *testing.T) {
	raw := []byte(sampleYAML)
	snippet := "remoteUser: vscode\n"
	result, err := edit.InsertBlock(raw, snippet)
	if err != nil {
		t.Fatalf("InsertBlock: %v", err)
	}
	blocks, _ := edit.ParseBlocksFromBytes(result)
	found := false
	for _, b := range blocks {
		if b.Key == "remoteUser" {
			found = true
		}
	}
	if !found {
		t.Error("remoteUser block not found after insert")
	}
}

func TestInsertBlockOrdered(t *testing.T) {
	// File has name + forwardPorts. Inserting "image" should land between them.
	base := "name: mydev\nforwardPorts:\n  - 3000\n"
	snippet := "image: ubuntu:22.04\n"

	result, err := edit.InsertBlock([]byte(base), snippet)
	if err != nil {
		t.Fatalf("InsertBlock: %v", err)
	}

	blocks, _ := edit.ParseBlocksFromBytes(result)
	order := make([]string, 0, len(blocks))
	for _, b := range blocks {
		order = append(order, b.Key)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 blocks, got %d: %v", len(order), order)
	}
	if order[0] != "name" || order[1] != "image" || order[2] != "forwardPorts" {
		t.Errorf("wrong order: %v, want [name image forwardPorts]", order)
	}
}

func TestValidateSnippet(t *testing.T) {
	valid := "remoteUser: vscode\n"
	if err := edit.ValidateSnippet(valid); err != nil {
		t.Errorf("expected valid snippet to pass: %v", err)
	}

	invalid := "remoteUser: :\n  broken:"
	if err := edit.ValidateSnippet(invalid); err == nil {
		t.Error("expected invalid snippet to fail validation")
	}
}
