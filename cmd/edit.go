package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/edit"
)

var editCmd = &cobra.Command{
	Use:   "edit [file]",
	Short: "Interactively edit a devcontainer config YAML file",
	Long:  "Opens a two-panel TUI to add, remove, and edit top-level blocks in a config.yaml file.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runEdit,
}

func runEdit(cmd *cobra.Command, args []string) {
	file := "config.yaml"
	if len(args) == 1 {
		file = args[0]
	}

	m, err := edit.New(file)
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
