package presets

import yeditpresets "github.com/lucasassuncao/yedit/presets"

// Source returns this package's preset registry as a yedit/presets.Source so
// the editor TUI can consume it without depending on devcontainerwizard.
func Source() yeditpresets.Source {
	return source{}
}

type source struct{}

func (source) ListFields() []string                          { return ListFields() }
func (source) ListPresets(field string) []string             { return ListPresets(field) }
func (source) PresetYAML(field, name string) (string, error) { return PresetYAML(field, name) }
