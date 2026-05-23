package docs

import (
	"log"

	"github.com/lucasassuncao/yedit/viewer"

	dcpresets "github.com/lucasassuncao/devcontainerwizard/internal/presets"

	"github.com/spf13/cobra"
)

var ShowExamplesCmd = &cobra.Command{
	Use:   "show-examples",
	Short: "Browse built-in devcontainer presets in YAML",
	Run:   runShowExamples,
}

func runShowExamples(cmd *cobra.Command, args []string) {
	err := viewer.Run(dcpresets.Source())
	if err != nil {
		log.Fatalf("viewer TUI: %v", err)
	}
}
