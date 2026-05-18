package model

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// GPUValue represents the gpu field in hostRequirements.
// The devcontainer spec allows a boolean, the string "optional", or a
// detailed GPURequirement object with cores/memory constraints.
type GPUValue struct {
	Bool        *bool
	StringVal   string
	Requirement *GPURequirement
}

// GPUBool returns a GPUValue representing a boolean GPU requirement.
func GPUBool(v bool) GPUValue { return GPUValue{Bool: &v} }

// GPUOptional returns a GPUValue representing an optional GPU.
func GPUOptional() GPUValue { return GPUValue{StringVal: "optional"} }

// GPURequire returns a GPUValue with detailed hardware requirements.
func GPURequire(r GPURequirement) GPUValue { return GPUValue{Requirement: &r} }

// GPUBoolPtr is a convenience constructor returning a pointer.
func GPUBoolPtr(v bool) *GPUValue { r := GPUBool(v); return &r }

// GPUOptionalPtr is a convenience constructor returning a pointer.
func GPUOptionalPtr() *GPUValue { r := GPUOptional(); return &r }

// GPURequirePtr is a convenience constructor returning a pointer.
func GPURequirePtr(r GPURequirement) *GPUValue { v := GPURequire(r); return &v }

func (g GPUValue) MarshalJSON() ([]byte, error) {
	if g.Bool != nil {
		return json.Marshal(*g.Bool)
	}
	if g.StringVal != "" {
		return json.Marshal(g.StringVal)
	}
	if g.Requirement != nil {
		return json.Marshal(g.Requirement)
	}
	return json.Marshal(nil)
}

func (g *GPUValue) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		g.Bool = &b
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		g.StringVal = s
		return nil
	}
	var r GPURequirement
	if err := json.Unmarshal(data, &r); err != nil {
		return fmt.Errorf("gpu must be a boolean, \"optional\", or an object with cores/memory: %w", err)
	}
	g.Requirement = &r
	return nil
}

// MarshalYAML implements yaml.Marshaler so yaml.Marshal produces the correct output.
func (g GPUValue) MarshalYAML() (any, error) {
	if g.Bool != nil {
		return *g.Bool, nil
	}
	if g.StringVal != "" {
		return g.StringVal, nil
	}
	return g.Requirement, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for direct yaml.v3 decoding.
func (g *GPUValue) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		if value.Tag == "!!bool" {
			var b bool
			if err := value.Decode(&b); err != nil {
				return err
			}
			g.Bool = &b
			return nil
		}
		g.StringVal = value.Value
		return nil
	}
	if value.Kind == yaml.MappingNode {
		var r GPURequirement
		if err := value.Decode(&r); err != nil {
			return err
		}
		g.Requirement = &r
		return nil
	}
	return fmt.Errorf("gpu must be a boolean, \"optional\", or an object")
}
