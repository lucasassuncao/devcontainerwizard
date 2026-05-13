package examples

import "github.com/charmbracelet/glamour"

// renderYAML wraps yaml in a markdown code fence and runs it through glamour
// for syntax highlighting. Falls back to the raw YAML if rendering fails.
func renderYAML(yaml string, width int) string {
	if width < 10 {
		width = 80
	}
	md := "```yaml\n" + yaml + "\n```\n"

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return yaml
	}
	out, err := r.Render(md)
	if err != nil {
		return yaml
	}
	return out
}
