// Package presets provides preset value maps for each top-level field of the
// DevContainer model. Presets are exposed via a string-keyed dispatcher used by
// the edit TUI overlay and the show-examples command.
//
// Every field follows one naming convention, so the accessors are not
// documented individually: XPreset(name) returns the value of preset name for
// field X, or the zero value when there is no such preset, and ListXPresets()
// returns that field's preset names sorted alphabetically with "base" first.
// PresetYAML and ListPresets dispatch to them by field name.
package presets

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// marshalAsBlock wraps value as a top-level YAML block: "field: <yaml>".
// Scalars render inline ("name: foo\n"); complex values are indented two spaces
// under the field key. Returns an error when value is a typed or untyped nil
// (signals "preset not found").
func marshalAsBlock(field string, value any) (string, error) {
	if isNil(value) {
		return "", fmt.Errorf("preset not found for %q", field)
	}

	// Use an explicit encoder with two-space indentation. yaml.Marshal defaults
	// to four spaces, which — combined with the two-space prefix below — produced
	// a mixed 2/4-space style that read inconsistently in the editor preview.
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(value); err != nil {
		return "", fmt.Errorf("marshal %s: %w", field, err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("marshal %s: %w", field, err)
	}
	body := strings.TrimRight(buf.String(), "\n")

	// Inline or indented is decided by what the encoder actually produced, read
	// back as a node. Testing the raw text instead meant any scalar containing a
	// colon was treated as a mapping, so "ubuntu:22.04" - a perfectly ordinary
	// image tag - came out as an indented block under "image:".
	//
	// A scalar always goes inline, including a multi-line one: the encoder's
	// output is the scalar's own representation, so "field: " prefixes it
	// correctly whether it is plain, quoted, or a "|-" block.
	if isScalarYAML(body) {
		return fmt.Sprintf("%s: %s\n", field, body), nil
	}

	lines := strings.Split(body, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return field + ":\n" + strings.Join(lines, "\n") + "\n", nil
}

// isScalarYAML reports whether the encoded document is a single scalar rather
// than a mapping or a sequence.
func isScalarYAML(body string) bool {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
		return false
	}
	return len(doc.Content) == 1 && doc.Content[0].Kind == yaml.ScalarNode
}

// isNil reports whether v is a typed or untyped nil. Catches *T(nil),
// map(nil), []T(nil), and the bare nil interface.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

// sortedKeys returns map keys sorted alphabetically with "base" always first.
// Used by every ListXPresets function so the picker UI always presents "base"
// at the top of the list.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for i, k := range keys {
		if k == "base" {
			keys[0], keys[i] = keys[i], keys[0]
			break
		}
	}
	return keys
}
