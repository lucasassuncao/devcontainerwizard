package cmd

import (
	"embed"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasassuncao/bezel/table"
	"github.com/spf13/cobra"
)

//go:embed templates/*.yaml
var templatesFS embed.FS

func initCmd() *cobra.Command {
	var (
		force    bool
		list     bool
		output   string
		template string
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new config.yaml file",
		Long:  "Create a new config.yaml from a template. Run with --list to see available templates.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runInitE(cmd, list, force, output, template)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing config.yaml")
	cmd.Flags().BoolVarP(&list, "list", "l", false, "List available templates")
	cmd.Flags().StringVarP(&output, "output", "o", defaultConfigPath, "Output file path")
	cmd.Flags().StringVarP(&template, "template", "t", "", "Template to use")
	return cmd
}

func runInitE(cmd *cobra.Command, list, force bool, output, template string) error {
	if list {
		printTemplateList(cmd.OutOrStdout())
		return nil
	}

	if template == "" {
		fmt.Fprintln(cmd.ErrOrStderr(), "Use: devcontainerwizard init --template <template>")
		fmt.Fprintln(cmd.ErrOrStderr())
		printTemplateList(cmd.ErrOrStderr())
		return fmt.Errorf("no template specified")
	}

	if _, err := os.Stat(output); err == nil && !force {
		return fmt.Errorf("file %q already exists — use --force to overwrite", output)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0o750); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	content, err := getTemplateContent(template)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, []byte(content), 0o600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s Created %q from template %q.\n\n%s devcontainerwizard edit -c %s\n",
		th.Success.Render("✓"), output, template, th.Dim.Render("Next:"), output)
	return nil
}

// printTemplateList writes the formatted template table to w.
// Used by both the no-template error path and --list.
func printTemplateList(w io.Writer) {
	templates, err := listAvailableTemplates()
	if err != nil {
		fmt.Fprintf(w, "  (could not list templates: %v)\n", err)
		return
	}
	rows := make([][]string, len(templates))
	for i, name := range templates {
		rows[i] = []string{th.Bold.Render(name), th.Dim.Render(readTemplateDescription(name))}
	}
	widths := table.Fit(make([]table.Column, 2), rows, math.MaxInt)
	fmt.Fprintln(w, "Available templates:")
	for _, r := range rows {
		fmt.Fprintln(w, "  "+table.Row(widths, r))
	}
}

// getTemplateContent returns the YAML of the named template with LF endings.
// The embed keeps the checkout's bytes, and a CR left in a CRLF checkout sends
// the terminal cursor to column 0 in the editor's preview, wrecking the layout.
func getTemplateContent(template string) (string, error) {
	data, err := templatesFS.ReadFile("templates/" + template + ".yaml")
	if err != nil {
		return "", fmt.Errorf("unknown template %q — run 'devcontainerwizard init --list' to see available templates", template)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

// listAvailableTemplates returns template names from the embedded filesystem,
// sorted alphabetically (the order ReadDir guarantees).
func listAvailableTemplates() ([]string, error) {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, err
	}
	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			templates = append(templates, strings.TrimSuffix(entry.Name(), ".yaml"))
		}
	}
	return templates, nil
}

// readTemplateDescription reads the "# description: " header from a template file.
// Returns an empty string if the header is absent.
func readTemplateDescription(name string) string {
	data, err := getTemplateContent(name)
	if err != nil {
		return ""
	}
	line, _, _ := strings.Cut(data, "\n")
	desc, found := strings.CutPrefix(line, "# description: ")
	if !found {
		return ""
	}
	return desc
}
