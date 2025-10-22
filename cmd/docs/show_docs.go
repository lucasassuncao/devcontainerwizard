// Package docs ...
package docs

import (
	"log"

	"github.com/lucasassuncao/devcontainerwizard/internal/docgenerator"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show-docs",
	Short: "Show documentation in terminal",
	Run:   runShow,
}

func runShow(cmd *cobra.Command, args []string) {
	gen, err := docgenerator.NewSchemaGenerator("docs/markdown", "docs/schema", true)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	types := model.GetAllTypes()

	docs, err := gen.GenerateSchemaAndDocsInMemory(types)
	if err != nil {
		log.Fatalf("Failed to generate docs: %v", err)
	}

	if err := docgenerator.RenderMarkdownDocsInTerminal(docs); err != nil {
		log.Fatalf("Failed to render docs: %v", err)
	}
}
