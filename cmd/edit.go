package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
	dcpresets "github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"github.com/lucasassuncao/yedit/editor"
)

var editConfig string

var editCmd = &cobra.Command{
	Use:           "edit",
	Short:         "Interactively edit a devcontainer config YAML file",
	Long:          "Opens a two-panel TUI to add, remove, and edit top-level blocks in a config.yaml file.",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runEditE,
}

func init() {
	editCmd.Flags().StringVarP(&editConfig, "config", "c", "config.yaml", "Path to the config file")
}

func runEditE(cmd *cobra.Command, _ []string) error {
	err := editor.Run(editor.Config{
		Path:    editConfig,
		Schema:  &model.DevContainer{},
		Title:   "devcontainer wizard",
		Presets: dcpresets.Source(),
		Validators: []editor.Validator{
			editor.MutuallyExclusive("image", "build", "dockerComposeFile"),
			editor.RequiredWith("service", "dockerComposeFile"),
			editor.RequiredWith("runServices", "dockerComposeFile"),
		},
		PreCheckedFields: devcontainer.PreCheckedFields(),
		FieldSnippets:    devcontainer.FieldSnippets(),
		Hidden:           []string{"dockerFile"}, // legacy alias, prefer build.dockerfile
	})
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
	}
	return err
}
