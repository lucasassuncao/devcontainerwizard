// Package docgenerator ...
package docgenerator

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/viewport"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/manifoldco/promptui"
)

type viewportModel struct {
	viewport         viewport.Model
	originalMarkdown string
	quitting         bool
}

func NewViewportModel(markdown string) viewportModel {
	vp := viewport.New(80, 24) // tamanho inicial padrão
	vp.SetContent(markdown)
	return viewportModel{
		viewport:         vp,
		originalMarkdown: markdown,
	}
}

func (m viewportModel) Init() bubbletea.Cmd { return nil }

func (m viewportModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, bubbletea.Quit
		case "up", "k":
			m.viewport.ScrollUp(1)
		case "down", "j":
			m.viewport.ScrollDown(1)
		case "pgup":
			m.viewport.ScrollUp(10)
		case "pgdn":
			m.viewport.ScrollDown(10)
		}
	case bubbletea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2

		// Re-renderiza o Markdown com largura correta
		renderer, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(msg.Width-4), // pequena margem
		)
		out, _ := renderer.Render(m.originalMarkdown)
		m.viewport.SetContent(out)
	}
	return m, nil
}

func (m viewportModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s\n(Use ↑/↓, PgUp/PgDn to scroll, q to quit)\n", m.viewport.View())
}

// RenderMarkdownDocsInTerminal renders the provided markdown documentation
// in an interactive terminal viewport using Bubble Tea.
func RenderMarkdownDocsInTerminal(docs map[string]string) error {
	if len(docs) == 0 {
		return fmt.Errorf("no documentation to display")
	}

	// Sort keys for consistent order
	var names []string
	for name := range docs {
		names = append(names, name)
	}
	sort.Strings(names)

	for {
		// List options using promptui
		prompt := promptui.Select{
			Label: "Select documentation to view",
			Items: append(names, "Exit"),
			Size:  len(names) + 1,
		}

		idx, choice, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		if choice == "Exit" || idx == len(names) {
			fmt.Println("Exiting documentation viewer.")
			return nil
		}

		// Selected markdown
		markdown := docs[choice]

		// Create and run Bubble Tea program
		p := bubbletea.NewProgram(NewViewportModel(markdown), bubbletea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to start viewport program: %w", err)
		}
	}
}
