package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/edit"
)

var editConfig string

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Interactively edit a devcontainer config YAML file",
	Long:  "Opens a two-panel TUI to add, remove, and edit top-level blocks in a config.yaml file.",
	Run:   runEdit,
}

func init() {
	editCmd.Flags().StringVarP(&editConfig, "config", "c", "config.yaml", "Path to the config file")
}

func runEdit(cmd *cobra.Command, args []string) {
	m, err := edit.New(editConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
