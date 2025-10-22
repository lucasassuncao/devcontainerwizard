// Package docgenerator ...
package docgenerator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/invopop/jsonschema"
)

type SchemaGenerator struct {
	reflector       *jsonschema.Reflector
	docsDir         string
	schemasDir      string
	visitedSections map[string]bool
	CleanupSchemas  bool
}

type SectionArgs struct {
	Builder      *strings.Builder
	Name         string
	Schema       *jsonschema.Schema
	Visited      map[*jsonschema.Schema]bool
	HeadingLevel int
	NameSuffix   string
}

// NewSchemaGenerator creates a new SchemaGenerator
func NewSchemaGenerator(docsDir, schemasDir string, cleanupSchemas bool) (*SchemaGenerator, error) {
	// Create directories if they do not exist
	for _, dir := range []string{docsDir, schemasDir} {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &SchemaGenerator{
		reflector: &jsonschema.Reflector{
			RequiredFromJSONSchemaTags: true,
			DoNotReference:             true,
		},
		docsDir:         docsDir,
		schemasDir:      schemasDir,
		visitedSections: make(map[string]bool),
		CleanupSchemas:  cleanupSchemas,
	}, nil
}

// GenerateSchemaAndDocs generates JSON schema and markdown docs for the given type
func (g *SchemaGenerator) GenerateSchemaAndDocs(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	typeName := t.Name()
	schema := g.reflector.Reflect(v)

	schemaPath := filepath.Join(g.schemasDir, fmt.Sprintf("%s.json", strings.ToLower(typeName)))

	if err := g.saveJSONSchema(schema, schemaPath); err != nil {
		return fmt.Errorf("error saving schema for %s: %w", typeName, err)
	}

	docsPath := filepath.Join(g.docsDir, fmt.Sprintf("%s.md", strings.ToLower(typeName)))
	if err := g.saveMarkdownDocs(schema, typeName, docsPath); err != nil {
		return fmt.Errorf("error saving docs for %s: %w", typeName, err)
	}

	if g.CleanupSchemas {
		if err := os.Remove(schemaPath); err != nil {
			return fmt.Errorf("error removing schema file %s: %w", schemaPath, err)
		}
	}

	return nil
}

// GenerateSchemaAndDocsInMemory generates Markdown docs in memory
// while respecting the schema generation and cleanup logic.
func (g *SchemaGenerator) GenerateSchemaAndDocsInMemory(types []interface{}) (map[string]string, error) {
	result := make(map[string]string)

	for _, v := range types {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		typeName := t.Name()
		schema := g.reflector.Reflect(v)

		// Caminho do schema temporário
		schemaPath := filepath.Join(g.schemasDir, fmt.Sprintf("%s.json", strings.ToLower(typeName)))

		// 1️⃣ Gera o schema JSON
		if err := g.saveJSONSchema(schema, schemaPath); err != nil {
			return nil, fmt.Errorf("error saving schema for %s: %w", typeName, err)
		}

		// 2️⃣ Gera o Markdown (em memória)
		markdown := g.generateMarkdownDocs(schema, typeName)
		result[typeName] = markdown

		// 3️⃣ Cleanup opcional dos schemas
		if g.CleanupSchemas {
			if err := os.Remove(schemaPath); err != nil {
				return nil, fmt.Errorf("error removing schema file %s: %w", schemaPath, err)
			}
		}
	}

	return result, nil
}

// saveJSONSchema saves the JSON schema to a file
func (g *SchemaGenerator) saveJSONSchema(schema *jsonschema.Schema, filename string) error {
	validatedPath, ok := validatePathWithinBase(g.schemasDir, filename)
	if !ok {
		return fmt.Errorf("invalid schema path: %s", filename)
	}

	file, err := os.OpenFile(validatedPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(schema)
}

// saveMarkdownDocs saves the markdown documentation to a file
func (g *SchemaGenerator) saveMarkdownDocs(schema *jsonschema.Schema, typeName, filename string) error {
	docs := g.generateMarkdownDocs(schema, typeName)
	validatedPath, ok := validatePathWithinBase(g.docsDir, filename)
	if !ok {
		return fmt.Errorf("invalid docs path: %s", filename)
	}
	return os.WriteFile(validatedPath, []byte(docs), 0600)
}

// getSectionKey generates a unique key for a schema section based on its properties
func (g *SchemaGenerator) getSectionKey(schema *jsonschema.Schema) string {
	if schema == nil || schema.Properties == nil {
		return ""
	}

	var props []string
	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		props = append(props, fmt.Sprintf("%s:%s", pair.Key, pair.Value.Type))
	}
	sort.Strings(props)
	return strings.Join(props, "|")
}

// hasVisitedSection checks if a section has been visited and marks it as visited
func (g *SchemaGenerator) hasVisitedSection(name string, schema *jsonschema.Schema) bool {
	key := fmt.Sprintf("%s:%s", name, g.getSectionKey(schema))
	visited := g.visitedSections[key]
	if !visited {
		g.visitedSections[key] = true
	}
	return visited
}

// getArrayItemSchema retrieves the schema of items in an array schema
func (g *SchemaGenerator) getArrayItemSchema(s *jsonschema.Schema) *jsonschema.Schema {
	if s == nil {
		return nil
	}
	v := reflect.ValueOf(s).Elem()
	itemsField := v.FieldByName("Items")
	if !itemsField.IsValid() || itemsField.IsZero() {
		return nil
	}

	return findFirstSchema(itemsField)
}

// getMapValueSchema retrieves the schema of values in a map schema
func (g *SchemaGenerator) getMapValueSchema(s *jsonschema.Schema) *jsonschema.Schema {
	if s == nil {
		return nil
	}

	if s.Properties != nil && s.Properties.Len() > 0 {
		return nil
	}

	v := reflect.ValueOf(s).Elem()
	apField := v.FieldByName("AdditionalProperties")
	if !apField.IsValid() || apField.IsZero() {
		return nil
	}

	valueSchema := findFirstSchema(apField)
	if valueSchema == nil {
		return nil
	}
	return valueSchema
}

// findFirstSchema recursively searches for the first *jsonschema.Schema in a reflect.Value
func findFirstSchema(val reflect.Value) *jsonschema.Schema {
	if !val.IsValid() {
		return nil
	}

	if val.CanInterface() {
		if s, ok := val.Interface().(*jsonschema.Schema); ok {
			return s
		}
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return findFirstSchema(val.Elem())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if s := findFirstSchema(val.Field(i)); s != nil {
				return s
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if s := findFirstSchema(val.Index(i)); s != nil {
				return s
			}
		}
	case reflect.Map:
		keys := val.MapKeys()
		if len(keys) > 0 {
			v := val.MapIndex(keys[0])
			if s := findFirstSchema(v); s != nil {
				return s
			}
		}
	}

	return nil
}

// validatePathWithinBase ensures the targetPath is within baseDir
func validatePathWithinBase(baseDir, targetPath string) (string, bool) {
	cleanBase := filepath.Clean(baseDir)
	cleanTarget := filepath.Clean(targetPath)

	rel, err := filepath.Rel(cleanBase, cleanTarget)
	if err != nil {
		return "", false
	}
	if rel == "." {
		return cleanBase, true
	}
	if strings.HasPrefix(rel, "..") {
		return "", false
	}
	return cleanTarget, true
}
