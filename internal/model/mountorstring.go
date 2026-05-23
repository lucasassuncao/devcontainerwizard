package model

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/lucasassuncao/yedit/schema"
)

// MountOrString represents an element of the mounts array.
// The devcontainer spec allows each mount to be either a Mount object or a
// Docker --mount string (e.g. "source=/path,target=/container,type=bind").
type MountOrString struct {
	Mount *Mount
	Str   string
}

// MountObject returns a MountOrString backed by a structured Mount.
func MountObject(m Mount) MountOrString { return MountOrString{Mount: &m} }

// MountString returns a MountOrString backed by a raw Docker mount string.
func MountString(s string) MountOrString { return MountOrString{Str: s} }

func (m MountOrString) MarshalJSON() ([]byte, error) {
	if m.Str != "" {
		return json.Marshal(m.Str)
	}
	return json.Marshal(m.Mount)
}

func (m *MountOrString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.Str = s
		return nil
	}
	var mount Mount
	if err := json.Unmarshal(data, &mount); err != nil {
		return fmt.Errorf("mount must be a string or an object: %w", err)
	}
	m.Mount = &mount
	return nil
}

// MarshalYAML implements yaml.Marshaler so yaml.Marshal produces the correct output.
func (m MountOrString) MarshalYAML() (any, error) {
	if m.Str != "" {
		return m.Str, nil
	}
	return m.Mount, nil
}

// YeditSchema implements yedit/schema.Provider, declaring the editor view of
// a mount item as the structured Mount form (string mounts edit through the
// raw YAML pane).
func (MountOrString) YeditSchema() []schema.FieldDef {
	return []schema.FieldDef{
		{YAMLName: "type", Kind: schema.KindScalar, Required: true, OneOf: []string{"bind", "volume", "tmpfs"}, Description: "Mount type."},
		{YAMLName: "source", Kind: schema.KindScalar, Description: "Source path or volume name."},
		{YAMLName: "target", Kind: schema.KindScalar, Required: true, Description: "Target path inside the container."},
		{YAMLName: "readonly", Kind: schema.KindScalar, Description: "Mount as read-only."},
	}
}

// UnmarshalYAML implements yaml.Unmarshaler for direct yaml.v3 decoding.
func (m *MountOrString) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		m.Str = value.Value
		return nil
	}
	if value.Kind == yaml.MappingNode {
		var mount Mount
		if err := value.Decode(&mount); err != nil {
			return err
		}
		m.Mount = &mount
		return nil
	}
	return fmt.Errorf("mount must be a string or an object")
}
