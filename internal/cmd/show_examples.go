package cmd

import (
	"github.com/lucasassuncao/yedit/viewer"

	"github.com/spf13/cobra"
)

func showExamplesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-examples",
		Short: "Browse built-in devcontainer presets in YAML",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return viewer.Run(DevContainerBlockPresets)
		},
	}
}
