package edit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AlertDismissedMsg is sent when the user closes the alert overlay.
type AlertDismissedMsg struct{}

type alertKind int

const (
	alertError alertKind = iota
	alertSuccess
)

// AlertModel is a simple modal that shows a message with an OK button.
type AlertModel struct {
	title  string
	lines  []string
	kind   alertKind
	totalW int
	totalH int
}

func NewAlert(title, message string, kind alertKind, totalW, totalH int) AlertModel {
	return AlertModel{
		title:  title,
		lines:  strings.Split(message, "\n"),
		kind:   kind,
		totalW: totalW,
		totalH: totalH,
	}
}

func (a AlertModel) borderColor() lipgloss.Color {
	if a.kind == alertSuccess {
		return lipgloss.Color("82")
	}
	return lipgloss.Color("196")
}

func (a AlertModel) buttonBG() lipgloss.Color {
	if a.kind == alertSuccess {
		return lipgloss.Color("82")
	}
	return lipgloss.Color("196")
}

func (a AlertModel) Update(msg tea.KeyMsg) (AlertModel, tea.Cmd) {
	switch msg.String() {
	case " ", "enter", "esc", "q":
		return a, func() tea.Msg { return AlertDismissedMsg{} }
	}
	return a, nil
}

func (a AlertModel) View() string {
	color := a.borderColor()

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(color)
	title := titleStyle.Render(a.title)

	maxW := 0
	for _, l := range a.lines {
		if len(l) > maxW {
			maxW = len(l)
		}
	}

	body := strings.Join(a.lines, "\n")
	ok := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Background(a.buttonBG()).
		Padding(0, 2).
		Render("  OK  ")

	// Centre the OK button under the message.
	okLine := lipgloss.NewStyle().Width(maxW).Align(lipgloss.Center).Render(ok)

	border := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(color).
		Padding(1, 2)

	box := border.Render(strings.Join([]string{title, "", body, "", okLine}, "\n"))

	bw := lipgloss.Width(box)
	bh := lipgloss.Height(box)
	lp := (a.totalW - bw) / 2
	tp := (a.totalH - bh) / 2
	if lp < 0 {
		lp = 0
	}
	if tp < 0 {
		tp = 0
	}
	return lipgloss.NewStyle().PaddingLeft(lp).PaddingTop(tp).Render(box)
}
