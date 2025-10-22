// Package docgenerator ...
package docgenerator

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// GenerateIndex creates an index markdown file listing all configuration structures.
func GenerateIndex(docsDir string, configs []interface{}) error {
	var sb strings.Builder

	sb.WriteString("# Documentation Index\n\n")
	sb.WriteString("This documentation describes all available configuration structurees.\n\n")
	sb.WriteString("## Available Configurations\n\n")

	for _, config := range configs {
		t := reflect.TypeOf(config)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		typeName := t.Name()
		fileName := strings.ToLower(typeName) + ".md"

		sb.WriteString(fmt.Sprintf("- [%s](./%s)\n", typeName, fileName))
	}

	return os.WriteFile(filepath.Join(docsDir, "index.md"), []byte(sb.String()), 0600)
}
