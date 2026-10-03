package cmd

import (
	"fmt"
	"strings"

	dcpresets "github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"github.com/lucasassuncao/yedit/presets"
)

// DevContainerBlockPresets feeds the preset picker inside block editors: one
// preset set per top-level field, declared in internal/presets.
var DevContainerBlockPresets = dcpresets.Source()

// DevContainerDocPresets is a whole-document preset source for the root
// template picker (p on the root list). It is backed by the same embedded
// templates the init command writes, so the picker and `init --template`
// never drift apart.
var DevContainerDocPresets presets.Source = docPresets{}

// docPresets implements presets.Source over templatesFS. Whole-document
// sources use the empty field name: PresetYAML("", name) returns the full
// YAML document for the named template.
type docPresets struct{}

func (docPresets) ListFields() []string { return []string{""} }

func (docPresets) ListPresets(field string) []string {
	if field != "" {
		return nil
	}
	names, err := listAvailableTemplates()
	if err != nil {
		return nil
	}
	return names
}

func (docPresets) PresetYAML(field, name string) (string, error) {
	if field != "" {
		return "", fmt.Errorf("docPresets: unknown field %q", field)
	}
	content, err := getTemplateContent(name)
	if err != nil {
		return "", err
	}
	// The "# description: " header is metadata for `init --list`, not part of
	// the document; the picker shows the template body only.
	if _, body, found := strings.Cut(content, "\n"); found && strings.HasPrefix(content, "# description: ") {
		return body, nil
	}
	return content, nil
}
