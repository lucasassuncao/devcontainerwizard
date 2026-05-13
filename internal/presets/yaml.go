// Package presets provides preset value maps for each top-level field of the
// DevContainer model. Presets are exposed via a string-keyed dispatcher used by
// the edit TUI overlay and the show-examples command.
package presets

import (
	"fmt"
	"reflect"
	"sort"
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

	out, err := yaml.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal %s: %w", field, err)
	}
	body := strings.TrimRight(string(out), "\n")

	// Single-line scalar (not a list item, not a mapping) → render inline.
	if !strings.Contains(body, "\n") && !strings.Contains(body, ":") && !strings.HasPrefix(body, "- ") {
		return fmt.Sprintf("%s: %s\n", field, body), nil
	}

	lines := strings.Split(body, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return field + ":\n" + strings.Join(lines, "\n") + "\n", nil
}

// isNil reports whether v is a typed or untyped nil. Catches *T(nil),
// map(nil), []T(nil), and the bare nil interface.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
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
	sort.Strings(keys)
	for i, k := range keys {
		if k == "base" {
			keys[0], keys[i] = keys[i], keys[0]
			break
		}
	}
	return keys
}
