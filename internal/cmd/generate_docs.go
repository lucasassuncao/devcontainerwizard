package cmd

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
	"github.com/lucasassuncao/docgen"
	"github.com/lucasassuncao/yedit/presets"

	"github.com/spf13/cobra"
)

func generateDocsCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "generate-docs",
		Short:  "Generate documentation for devcontainerwizard",
		RunE:   runGenerate,
		Hidden: true,
	}
}

func runGenerate(cmd *cobra.Command, _ []string) error {
	return generateDocs(cmd.OutOrStdout())
}

func generateDocs(w io.Writer) error {
	fmt.Fprintln(w, "Generating documentation...")

	docsDir := "docs/devcontainerwizard"
	attributesDir := filepath.Join(docsDir, "attributes")
	examplesDir := filepath.Join(docsDir, "examples")
	schemaDir := filepath.Join(docsDir, "schema")

	// The config file is one document with one root type, so a single entry
	// covers it; FilePerField gives each nested struct its own reference page.
	entries := []docgen.Entry{{Config: model.DevContainer{}}}

	if _, err := docgen.Generate(entries,
		docgen.WithMarkdown(attributesDir, docgen.Layout{Files: docgen.FilePerField}),
		// The templates fill the DevContainer page; each block with its own
		// attributes page gets its block presets. Titles match the attributes
		// page names, so each one links to its examples.
		docgen.WithExamples(presets.Combine(DevContainerDocPresets, DevContainerBlockPresets), examplesDir, map[string]string{
			"":                     "DevContainer",
			"build":                "build",
			"portsAttributes":      "portsAttributes",
			"otherPortsAttributes": "otherPortsAttributes",
			"mounts":               "mounts",
			"hostRequirements":     "hostRequirements",
			"watch":                "watch",
			"customizations":       "customizations",
			"secrets":              "secrets",
		}),
		docgen.WithIndex(docsDir),
	); err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}

	// The JSON Schema describes the config file as a whole, so a language
	// server can complete and validate config.yaml directly.
	if _, err := docgen.Generate(
		[]docgen.Entry{{Config: model.DevContainer{}}},
		docgen.WithJSONSchema(schemaDir),
	); err != nil {
		return fmt.Errorf("failed to generate json schema: %w", err)
	}

	fmt.Fprintf(w, "Documentation generated in '%s' directory.\n", docsDir)
	return nil
}
