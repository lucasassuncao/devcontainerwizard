package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
	dcpresets "github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"gopkg.in/yaml.v3"
)

// Every preset the TUI can insert must be valid YAML that sets the field it
// claims to.
func TestAllPresetsAreValidYAML(t *testing.T) {
	for _, field := range dcpresets.ListFields() {
		for _, name := range dcpresets.ListPresets(field) {
			body, err := dcpresets.PresetYAML(field, name)
			if err != nil {
				t.Errorf("%s/%s: %v", field, name, err)
				continue
			}
			var doc map[string]any
			if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
				t.Errorf("%s/%s: not valid YAML: %v\n%s", field, name, err, body)
				continue
			}
			if _, ok := doc[field]; !ok {
				t.Errorf("%s/%s: block does not set %q\n%s", field, name, field, body)
			}
		}
	}
}

// Presets must also be accepted by the parser they feed.
func TestAllPresetsParse(t *testing.T) {
	for _, field := range dcpresets.ListFields() {
		for _, name := range dcpresets.ListPresets(field) {
			body, err := dcpresets.PresetYAML(field, name)
			if err != nil {
				continue
			}
			dir := t.TempDir()
			cfg := filepath.Join(dir, "c.yaml")
			if err := os.WriteFile(cfg, []byte(scaffold(field)+body), 0o600); err != nil {
				t.Fatal(err)
			}
			k, err := devcontainer.LoadYAMLFile(cfg)
			if err != nil {
				t.Errorf("%s/%s: load: %v\n%s", field, name, err, body)
				continue
			}
			if _, err := devcontainer.Parse(k); err != nil {
				t.Errorf("%s/%s: parse: %v\n%s", field, name, err, body)
			}
		}
	}
}

// Unknown fields and presets must error, never panic.
func TestPresetLookupMisses(t *testing.T) {
	for _, f := range []string{"", "nope", "IMAGE", "name "} {
		if _, err := dcpresets.PresetYAML(f, "base"); err == nil {
			t.Errorf("PresetYAML(%q, base) returned nil error", f)
		}
		if got := dcpresets.ListPresets(f); got != nil {
			t.Errorf("ListPresets(%q) = %v, want nil", f, got)
		}
	}
	for _, field := range dcpresets.ListFields() {
		if _, err := dcpresets.PresetYAML(field, "definitely-not-a-preset"); err == nil {
			t.Errorf("PresetYAML(%s, bogus) returned nil error", field)
		}
	}
}

// The root template picker is backed by the same embedded templates init writes,
// minus the "# description: " header, which is metadata for `init --list`
// rather than part of the document.
func TestDocPresetsRoundTrip(t *testing.T) {
	for _, name := range DevContainerDocPresets.ListPresets("") {
		body, err := DevContainerDocPresets.PresetYAML("", name)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if strings.HasPrefix(body, "# description:") {
			t.Errorf("%s: description header leaked into the document body", name)
		}
		var doc map[string]any
		if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
			t.Errorf("%s: not valid YAML: %v", name, err)
		}
	}
	if _, err := DevContainerDocPresets.PresetYAML("", "nope"); err == nil {
		t.Error("unknown doc preset returned nil error")
	}
}

// scaffold supplies the required keys a preset block does not itself set, so the
// probe document is complete without duplicating the field under test.
func scaffold(field string) string {
	s := ""
	if field != "name" {
		s += "name: probe\n"
	}
	switch field {
	case "image", "build", "dockerComposeFile", "service", "runServices":
	default:
		s += "image: ubuntu\n"
	}
	return s
}
