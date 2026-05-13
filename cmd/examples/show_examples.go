package examples

import (
	"log"

	"github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"github.com/lucasassuncao/devcontainerwizard/internal/tui/examples"

	"github.com/spf13/cobra"
)

// ShowCmd is the cobra entrypoint for browsing built-in devcontainer presets.
var ShowCmd = &cobra.Command{
	Use:   "show-examples",
	Short: "Browse built-in devcontainer presets in YAML",
	Run:   runShow,
}

func runShow(cmd *cobra.Command, args []string) {
	err := examples.Run(presets.ListFields(), presets.ListPresets, presets.PresetYAML)
	if err != nil {
		log.Fatalf("examples TUI: %v", err)
	}
}
