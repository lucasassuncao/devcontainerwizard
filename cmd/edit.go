package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/edit"
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

func runEditE(cmd *cobra.Command, args []string) error {
	m, err := edit.New(editConfig)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "TUI error: %v\n", err)
		return err
	}
	return nil
}
