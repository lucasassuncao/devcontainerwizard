package edit

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// PreviewModel wraps a textarea that is always visible as the right panel.
type PreviewModel struct {
	ta      textarea.Model
	content string
}

func NewPreviewModel(width, height int) PreviewModel {
	ta := textarea.New()
	ta.SetWidth(width)
	ta.SetHeight(height)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Blur()
	return PreviewModel{ta: ta}
}

func (pm *PreviewModel) SetContent(yaml string) {
	pm.content = strings.ReplaceAll(yaml, "\r\n", "\n")
	pm.ta.SetValue(pm.content)
}

func (pm PreviewModel) Value() string {
	return pm.ta.Value()
}

func (pm *PreviewModel) ScrollToKey(key string) {
	if key == "" {
		return
	}
	target := key + ":"
	lines := strings.Split(pm.content, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, target) {
			pm.ta.SetCursor(i)
			return
		}
	}
}

func (pm *PreviewModel) Focus() tea.Cmd {
	return pm.ta.Focus()
}

func (pm *PreviewModel) Blur() {
	pm.ta.Blur()
}

func (pm *PreviewModel) Resize(width, height int) {
	pm.ta.SetWidth(width)
	pm.ta.SetHeight(height)
}

func (pm PreviewModel) Update(msg tea.Msg) (PreviewModel, tea.Cmd) {
	var cmd tea.Cmd
	pm.ta, cmd = pm.ta.Update(msg)
	return pm, cmd
}

func (pm PreviewModel) View() string {
	return pm.ta.View()
}
