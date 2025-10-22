// Package docgenerator ...
package docgenerator

import (
	"strings"

	"github.com/invopop/jsonschema"
)

// getAnchorForProperty generates a markdown anchor for a property if it has nested structures
func (g *SchemaGenerator) getAnchorForProperty(propName string, prop *jsonschema.Schema) string {
	if prop == nil {
		return ""
	}

	if prop.Properties != nil && prop.Properties.Len() > 0 {
		return strings.ToLower(propName)
	}

	if itemSchema := g.getArrayItemSchema(prop); itemSchema != nil && itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
		return strings.ToLower(propName) + "-item"
	}

	if valueSchema := g.getMapValueSchema(prop); valueSchema != nil {
		if valueSchema.Properties != nil && valueSchema.Properties.Len() > 0 {
			return strings.ToLower(propName) + "-value"
		}

		if itemSchema := g.getArrayItemSchema(valueSchema); itemSchema != nil && itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
			return strings.ToLower(propName) + "-value-item"
		}
	}

	return ""
}

// getDisplayType returns a human-readable type description for a schema property
func (g *SchemaGenerator) getDisplayType(prop *jsonschema.Schema) string {
	return g.describeType(prop)
}

// describeType returns a string description of the type represented by the schema
func (g *SchemaGenerator) describeType(s *jsonschema.Schema) string {
	if s == nil {
		return "-"
	}

	if value := g.getMapValueSchema(s); value != nil {
		elem := g.describeType(value)
		if elem == "-" || elem == "object" || strings.HasPrefix(elem, "map[]") || strings.HasPrefix(elem, "array[") {
			elem = "object"
		}
		return "map[string]" + elem
	}

	if item := g.getArrayItemSchema(s); item != nil {
		elem := g.describeType(item)
		if elem == "-" || elem == "object" || strings.HasPrefix(elem, "map[]") || strings.HasPrefix(elem, "array[") {
			elem = "object"
		}
		return "array[" + elem + "]"
	}

	if s.Type != "" {
		return s.Type
	}

	if s.Properties != nil && s.Properties.Len() > 0 {
		return "object"
	}

	return "-"
}

// isRequired checks if a property is required in the schema
func isRequired(schema *jsonschema.Schema, propertyName string) bool {
	if schema.Required == nil {
		return false
	}
	for _, req := range schema.Required {
		if req == propertyName {
			return true
		}
	}
	return false
}
