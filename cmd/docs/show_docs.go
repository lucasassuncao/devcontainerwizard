package docs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	if err := showDocs(); err != nil {
		log.Fatalf("%v", err)
	}
}

func showDocs() error {
	tmpDir, err := os.MkdirTemp("", "devcontainerwizard-docs-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	markdownDir := filepath.Join(tmpDir, "markdown")
	schemaDir := filepath.Join(tmpDir, "schema")

	gen, err := docgenerator.NewSchemaGenerator(markdownDir, schemaDir, docgenerator.WithCleanupSchemas())
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	docs, err := gen.GenerateSchemaAndDocsInMemory(model.GetAllTypes())
	if err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}

	if err := docgenerator.RenderMarkdownDocsInTerminal(docs); err != nil {
		return fmt.Errorf("failed to render docs: %w", err)
	}
	return nil
}
