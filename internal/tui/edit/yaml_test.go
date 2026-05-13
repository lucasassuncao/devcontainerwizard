package edit_test

import (
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

func TestParseBlocksFromBytes(t *testing.T) {
	blocks, err := edit.ParseBlocksFromBytes([]byte(sampleYAML))
	if err != nil {
		t.Fatalf("ParseBlocksFromBytes: %v", err)
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

func TestValidateKnownKeys(t *testing.T) {
	valid := []byte("name: mydev\nimage: ubuntu:22.04\n")
	if unknown := edit.ValidateKnownKeys(valid); len(unknown) != 0 {
		t.Errorf("expected no unknown keys, got %v", unknown)
	}

	withUnknown := []byte("name: mydev\ncustomization: bad\nimage: ubuntu:22.04\n")
	unknown := edit.ValidateKnownKeys(withUnknown)
	if len(unknown) != 1 || unknown[0] != "customization" {
		t.Errorf("expected [customization], got %v", unknown)
	}

	// Sub-key validation: customizations.vscod is not a valid sub-key.
	withBadSub := []byte("customizations:\n  vscod:\n    extensions:\n      - foo.bar\n")
	unknown2 := edit.ValidateKnownKeys(withBadSub)
	if len(unknown2) != 1 || unknown2[0] != "customizations.vscod" {
		t.Errorf("expected [customizations.vscod], got %v", unknown2)
	}

	// Valid sub-key should pass.
	withGoodSub := []byte("customizations:\n  vscode:\n    extensions:\n      - foo.bar\n")
	if unknown3 := edit.ValidateKnownKeys(withGoodSub); len(unknown3) != 0 {
		t.Errorf("expected no unknown keys for valid sub-key, got %v", unknown3)
	}

	// Third level: customizations.vscode.extension is invalid (should be extensions).
	withBadNested := []byte("customizations:\n  vscode:\n    extension:\n      - foo.bar\n    settings:\n      editor.formatOnSave: true\n")
	unknown4 := edit.ValidateKnownKeys(withBadNested)
	if len(unknown4) != 1 || unknown4[0] != "customizations.vscode.extension" {
		t.Errorf("expected [customizations.vscode.extension], got %v", unknown4)
	}

	// settings under vscode is a free-form object — its keys must not be validated.
	withSettings := []byte("customizations:\n  vscode:\n    extensions:\n      - foo.bar\n    settings:\n      editor.formatOnSave: true\n      any.arbitrary.key: 42\n")
	if unknown5 := edit.ValidateKnownKeys(withSettings); len(unknown5) != 0 {
		t.Errorf("expected no errors for free-form settings, got %v", unknown5)
	}
}

func TestValidateKnownKeysDepth(t *testing.T) {
	cases := []struct {
		desc    string
		yaml    string
		wantErr bool
	}{
		{"typo top-level", "buiild:\n  dockerfile: Dockerfile\n", true},
		{"typo subfield", "build:\n  dockerfilee: Dockerfile\n  context: .\n", true},
		{"free-form args", "build:\n  dockerfile: Dockerfile\n  context: .\n  args:\n    MY_ARG: value\n    OTHER_ARG: x\n", false},
		{"list items in cacheFrom", "build:\n  dockerfile: Dockerfile\n  context: .\n  cacheFrom:\n    - myregistry/image:cache\n    - other:cache\n", false},
		{"fully populated build", "build:\n  dockerfile: Dockerfile\n  context: .\n  args:\n    MY_ARG: value\n  target: dev\n  cacheFrom:\n    - myregistry/image:cache\n  output: type=local,dest=./out\n  ssh:\n    - default\n", false},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			unknown := edit.ValidateKnownKeys([]byte(c.yaml))
			if c.wantErr && len(unknown) == 0 {
				t.Error("expected validation error, got none")
			}
			if !c.wantErr && len(unknown) != 0 {
				t.Errorf("expected no errors, got %v", unknown)
			}
		})
	}
}
