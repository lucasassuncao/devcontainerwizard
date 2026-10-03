package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
)

// Every embedded template must survive the exact pipeline `convert` runs.
// A template that cannot convert is shipped breakage.
func TestAllTemplatesConvert(t *testing.T) {
	for _, name := range templateNames(t) {
		t.Run(name, func(t *testing.T) {
			content, err := getTemplateContent(name)
			if err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			cfg := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(cfg, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			k, err := devcontainer.LoadYAMLFile(cfg)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			dc, err := devcontainer.Parse(k)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if err := devcontainer.Validate(dc); err != nil {
				t.Fatalf("validate: %s", devcontainer.HumanizeValidationError(err))
			}
			if _, err := devcontainer.WriteFile(dc, filepath.Join(dir, "o.json"), true); err != nil {
				t.Fatalf("write: %v", err)
			}
		})
	}
}

// Every key in every shipped template must be one the model actually supports.
//
// mapstructure drops keys it cannot match, so a typo like "remoteUsr" or
// "fowardPorts" parses clean, validates clean, and ships a template whose
// setting silently does nothing. TestAllTemplatesConvert cannot catch that -
// the document is still valid, just missing what the author meant. Decoding
// through ParseStrict is what ties the templates to the model.
func TestTemplatesUseOnlyKnownKeys(t *testing.T) {
	for _, name := range templateNames(t) {
		t.Run(name, func(t *testing.T) {
			content, err := getTemplateContent(name)
			if err != nil {
				t.Fatal(err)
			}
			k, err := devcontainer.LoadYAML(strings.NewReader(content))
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			if _, err := devcontainer.ParseStrict(k); err != nil {
				t.Errorf("template has keys the model does not define: %v", err)
			}
		})
	}
}

// The header feeds `init --list`; readTemplateDescription returns "" when it is
// missing, which would list a template with a blank description.
func TestTemplatesDeclareADescription(t *testing.T) {
	for _, name := range templateNames(t) {
		if desc := readTemplateDescription(name); desc == "" {
			t.Errorf("template %q has no %q header", name, "# description: ")
		}
	}
}

func templateNames(t *testing.T) []string {
	t.Helper()
	names, err := listAvailableTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatal("no templates found")
	}
	return names
}

// TestTemplatesReachTheEditorWithLFEndings guards CRLF checkouts: a CR in a
// template sends the cursor to column 0 in the preset preview and wrecks the layout.
func TestTemplatesReachTheEditorWithLFEndings(t *testing.T) {
	names, err := listAvailableTemplates()
	if err != nil {
		t.Fatalf("listAvailableTemplates: %v", err)
	}
	for _, name := range names {
		content, err := getTemplateContent(name)
		if err != nil {
			t.Fatalf("getTemplateContent(%q): %v", name, err)
		}
		if strings.Contains(content, "\r") {
			t.Errorf("template %q still carries a CR", name)
		}
		if strings.Contains(readTemplateDescription(name), "\r") {
			t.Errorf("description of %q still carries a CR", name)
		}
	}
}
