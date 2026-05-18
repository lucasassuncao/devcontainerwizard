package model

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// CommandValue represents a lifecycle-command field that the devcontainer spec
// allows as a string, an array of strings, or a named-command object
// (map[string]string|[]string). Named commands run in parallel; string/slice
// forms run sequentially in a shell.
type CommandValue struct {
	Items []string            `json:"-" yaml:"-"`
	Named map[string][]string `json:"-" yaml:"-"`
}

// CommandString returns a CommandValue for a single shell command.
func CommandString(s string) CommandValue { return CommandValue{Items: []string{s}} }

// CommandSlice returns a CommandValue for an argument list (no shell expansion).
func CommandSlice(s []string) CommandValue { return CommandValue{Items: s} }

// CommandMap returns a CommandValue for a set of named parallel commands.
func CommandMap(m map[string][]string) CommandValue { return CommandValue{Named: m} }

func (c CommandValue) MarshalJSON() ([]byte, error) {
	if c.Named != nil {
		// Collapse single-element slices to plain strings for cleaner JSON output.
		out := make(map[string]any, len(c.Named))
		for k, v := range c.Named {
			if len(v) == 1 {
				out[k] = v[0]
			} else {
				out[k] = v
			}
		}
		return json.Marshal(out)
	}
	if len(c.Items) == 1 {
		return json.Marshal(c.Items[0])
	}
	return json.Marshal(c.Items)
}

func (c *CommandValue) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		c.Items = []string{str}
		return nil
	}
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		c.Items = slice
		return nil
	}
	// Try map[string][]string first, then map[string]string for backwards compat.
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("must be a string, array of strings, or named command object: %w", err)
	}
	named := make(map[string][]string, len(obj))
	for k, raw := range obj {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			named[k] = []string{s}
			continue
		}
		var sl []string
		if err := json.Unmarshal(raw, &sl); err != nil {
			return fmt.Errorf("named command value for %q must be a string or array of strings: %w", k, err)
		}
		named[k] = sl
	}
	c.Named = named
	return nil
}

// MarshalYAML implements yaml.Marshaler so yaml.Marshal produces the correct output.
func (c CommandValue) MarshalYAML() (any, error) {
	if c.Named != nil {
		// Collapse single-element slices to plain strings for cleaner YAML output.
		out := make(map[string]any, len(c.Named))
		for k, v := range c.Named {
			if len(v) == 1 {
				out[k] = v[0]
			} else {
				out[k] = v
			}
		}
		return out, nil
	}
	if len(c.Items) == 1 {
		return c.Items[0], nil
	}
	return c.Items, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for direct yaml.v3 decoding.
func (c *CommandValue) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		c.Items = []string{value.Value}
	case yaml.SequenceNode:
		return value.Decode(&c.Items)
	case yaml.MappingNode:
		// Each map value can be a string or a sequence of strings.
		raw := make(map[string]yaml.Node)
		if err := value.Decode(&raw); err != nil {
			return err
		}
		named := make(map[string][]string, len(raw))
		for k, n := range raw {
			switch n.Kind {
			case yaml.ScalarNode:
				named[k] = []string{n.Value}
			case yaml.SequenceNode:
				var sl []string
				if err := n.Decode(&sl); err != nil {
					return fmt.Errorf("named command value for %q must be a string or array of strings: %w", k, err)
				}
				named[k] = sl
			default:
				return fmt.Errorf("named command value for %q must be a string or array of strings", k)
			}
		}
		c.Named = named
	default:
		return fmt.Errorf("command value must be a string, array of strings, or named command object")
	}
	return nil
}
