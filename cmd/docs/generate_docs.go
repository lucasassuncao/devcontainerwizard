// Package docs ...
package docs

import (
	"fmt"
	"log"

	"github.com/lucasassuncao/devcontainerwizard/internal/docgenerator"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"

	"path/filepath"

	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:    "generate-docs",
	Short:  "Generate documentation for devcontainer",
	Run:    runGenerate,
	Hidden: true,
}

func runGenerate(cmd *cobra.Command, args []string) {
	fmt.Println("Generating documentation...")

	docsDir := "docs"
	markdownDir := filepath.Join(docsDir, "markdown")
	schemaDir := filepath.Join(docsDir, "schema")

	gen, err := docgenerator.NewSchemaGenerator(markdownDir, schemaDir, false)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	types := model.GetAllTypes()

	for _, t := range types {
		if err := gen.GenerateSchemaAndDocs(t); err != nil {
			log.Fatalf("Failed to generate docs for %T: %v", t, err)
		}
	}

	if err := docgenerator.GenerateIndex(docsDir, types); err != nil {
		log.Fatalf("Failed to generate index: %v", err)
	}

	fmt.Println("Documentation generated in 'docs' directory.")
}
