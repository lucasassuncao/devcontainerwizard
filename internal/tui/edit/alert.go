package edit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucasassuncao/devcontainerwizard/internal/tui/theme"
)

// AlertDismissedMsg is sent when the user closes the alert overlay.
type AlertDismissedMsg struct{}

type alertKind int

const (
	alertError alertKind = iota
	alertSuccess
	alertConfirm
)

// AlertModel is a simple modal that shows a message with an OK button,
// or a Yes/No pair when kind is alertConfirm.
// confirmYes tracks which button is focused (true = Yes, false = No).
// confirmCmd is the command executed when the user confirms (confirm kind only).
type AlertModel struct {
	title      string
	lines      []string
	kind       alertKind
	confirmYes bool    // only meaningful when kind == alertConfirm
	confirmCmd tea.Cmd // command to run when user picks Yes
	totalW     int
	totalH     int
}

// NewAlert creates an informational alert with a single OK button.
func NewAlert(title, message string, kind alertKind, totalW, totalH int) AlertModel {
	return AlertModel{
		title:  title,
		lines:  strings.Split(message, "\n"),
		kind:   kind,
		totalW: totalW,
		totalH: totalH,
	}
}

// NewConfirmAlert creates a yes/no alert that runs confirmCmd when the user confirms.
func NewConfirmAlert(title, message string, confirmCmd tea.Cmd, totalW, totalH int) AlertModel {
	return AlertModel{
		title:      title,
		lines:      strings.Split(message, "\n"),
		kind:       alertConfirm,
		confirmYes: true,
		confirmCmd: confirmCmd,
		totalW:     totalW,
		totalH:     totalH,
	}
}

func (a AlertModel) accentColor() lipgloss.Color {
	switch a.kind {
	case alertSuccess:
		return theme.Success
	case alertConfirm:
		return theme.Accent
	default:
		return theme.Danger
	}
}

func (a AlertModel) Update(msg tea.KeyMsg) (AlertModel, tea.Cmd) {
	if a.kind == alertConfirm {
		switch msg.String() {
		case "left", "h", "right", "l", "tab":
			a.confirmYes = !a.confirmYes
		case "y", "Y":
			return a, a.confirmCmd
		case "n", "N", "esc", "q":
			return a, func() tea.Msg { return AlertDismissedMsg{} }
		case "enter", " ":
			if a.confirmYes {
				return a, a.confirmCmd
			}
			return a, func() tea.Msg { return AlertDismissedMsg{} }
		}
		return a, nil
	}
	switch msg.String() {
	case " ", "enter", "esc", "q":
		return a, func() tea.Msg { return AlertDismissedMsg{} }
	}
	return a, nil
}

func (a AlertModel) View() string {
	color := a.accentColor()

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(color)
	title := titleStyle.Render(a.title)

	maxW := 0
	for _, l := range a.lines {
		if len(l) > maxW {
			maxW = len(l)
		}
	}

	body := strings.Join(a.lines, "\n")

	btnStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Background(color).
		Padding(0, 2)

	var buttons string
	if a.kind == alertConfirm {
		inactiveStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("231")).
			Background(theme.Muted).
			Padding(0, 2)
		yesStyle, noStyle := inactiveStyle, inactiveStyle
		if a.confirmYes {
			yesStyle = btnStyle
		} else {
			noStyle = btnStyle
		}
		yes := yesStyle.Render("  Yes  ")
		no := noStyle.Render("  No  ")
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, yes, "  ", no)
	} else {
		ok := btnStyle.Render("  OK  ")
		buttons = lipgloss.NewStyle().Width(maxW).Align(lipgloss.Center).Render(ok)
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(color).
		Padding(1, 2)

	box := border.Render(strings.Join([]string{title, "", body, "", buttons}, "\n"))

	return centerBox(box, a.totalW, a.totalH)
}
