// Package docgenerator ...
package docgenerator

import (
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
)

// SchemaGenerator generates markdown documentation from JSON schemas
func (g *SchemaGenerator) generateMarkdownDocs(schema *jsonschema.Schema, typeName string) string {
	var sb strings.Builder

	g.visitedSections = make(map[string]bool)

	sb.WriteString(fmt.Sprintf("# %s\n\n", typeName))
	if schema.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", schema.Description))
	}

	sb.WriteString("## Arguments\n\n")
	sb.WriteString("The following arguments are supported:\n\n")
	g.writePropertiesTable(&sb, schema)

	visited := make(map[*jsonschema.Schema]bool)
	g.writeNestedObjectSections(&sb, schema, visited, 3)

	return sb.String()
}

// writePropertiesTable writes a markdown table of properties for the given schema
func (g *SchemaGenerator) writePropertiesTable(sb *strings.Builder, schema *jsonschema.Schema) {
	sb.WriteString("| Name | Type | Description | Required | Default |\n")
	sb.WriteString("|------|------|-------------|----------|---------|\n")

	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		propName := pair.Key
		prop := pair.Value

		description := strings.ReplaceAll(prop.Description, "|", "\\|")
		description = strings.Join(strings.Fields(description), " ")

		required := "No"
		if isRequired(schema, propName) {
			required = "Yes"
		}

		defaultValue := "-"
		if prop.Default != nil {
			defaultValue = fmt.Sprintf("%v", prop.Default)
		}

		displayName := fmt.Sprint(propName)
		if anchor := g.getAnchorForProperty(propName, prop); anchor != "" {
			displayName = fmt.Sprintf("[%s](#%s)", propName, anchor)
		}

		displayType := g.getDisplayType(prop)
		fmt.Fprintf(sb, "| %s | %s | %s | %s | %s |\n", displayName, description, displayType, required, defaultValue)
	}
	sb.WriteString("\n")
}

// writeNestedObjectSections writes sections for nested object properties
func (g *SchemaGenerator) writeNestedObjectSections(sb *strings.Builder, schema *jsonschema.Schema, visited map[*jsonschema.Schema]bool, headingLevel int) {
	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		child := pair.Value

		args := SectionArgs{
			Builder:      sb,
			Name:         name,
			Visited:      visited,
			HeadingLevel: headingLevel,
		}

		if child != nil && child.Properties != nil && child.Properties.Len() > 0 {
			if g.hasVisitedSection(name, child) {
				continue
			}
			args.Schema = child
			g.writeObjectSection(args)
		}

		if itemSchema := g.getArrayItemSchema(child); itemSchema != nil && itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
			if g.hasVisitedSection(name+"-item", itemSchema) {
				continue
			}
			args.Schema = itemSchema
			g.writeArrayItemSection(args)
		}

		if valueSchema := g.getMapValueSchema(child); valueSchema != nil && valueSchema.Properties != nil && valueSchema.Properties.Len() > 0 {
			if g.hasVisitedSection(name+"-value", valueSchema) {
				continue
			}
			args.Schema = valueSchema
			g.writeMapValueSection(args)
		}
	}
}

// writeSection writes a markdown section for the given schema
func (g *SchemaGenerator) writeSection(args SectionArgs) {
	hashes := strings.Repeat("#", args.HeadingLevel)
	title := args.Name
	if args.NameSuffix != "" {
		title += args.NameSuffix
	}
	fmt.Fprintf(args.Builder, "%s %s\n\n", hashes, title)
	if args.Schema.Description != "" {
		args.Builder.WriteString(fmt.Sprintf("%s\n\n", args.Schema.Description)) //nolint: staticcheck
	}
	args.Builder.WriteString("The following arguments are supported:\n\n")
	g.writePropertiesTable(args.Builder, args.Schema)
	g.writeNestedObjectSections(args.Builder, args.Schema, args.Visited, args.HeadingLevel+1)
}

// writeObjectSection writes a section for an object schema
func (g *SchemaGenerator) writeObjectSection(args SectionArgs) {
	if args.Visited[args.Schema] {
		return
	}
	args.Visited[args.Schema] = true
	g.writeSection(args)
}

// writeArrayItemSection writes a section for an array item schema
func (g *SchemaGenerator) writeArrayItemSection(args SectionArgs) {
	if args.Visited[args.Schema] {
		return
	}
	args.Visited[args.Schema] = true
	args.NameSuffix = " Item"
	g.writeSection(args)
}

// writeMapValueSection writes a section for a map value schema
func (g *SchemaGenerator) writeMapValueSection(args SectionArgs) {
	if args.Visited[args.Schema] {
		return
	}
	args.Visited[args.Schema] = true
	args.NameSuffix = " Value"
	g.writeSection(args)
}
