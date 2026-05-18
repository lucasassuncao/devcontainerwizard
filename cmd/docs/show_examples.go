package docs

import (
	"log"

	"github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"github.com/lucasassuncao/devcontainerwizard/internal/tui/examples"

	"github.com/spf13/cobra"
)

var ShowExamplesCmd = &cobra.Command{
	Use:   "show-examples",
	Short: "Browse built-in devcontainer presets in YAML",
	Run:   runShowExamples,
}

func runShowExamples(cmd *cobra.Command, args []string) {
	err := examples.Run(presets.ListFields(), presets.ListPresets, presets.PresetYAML)
	if err != nil {
		log.Fatalf("examples TUI: %v", err)
	}
}
