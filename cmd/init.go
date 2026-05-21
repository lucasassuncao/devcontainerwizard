package cmd

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed templates/*.yaml
var templatesFS embed.FS

var initCmd = newInitCmd()

func newInitCmd() *cobra.Command {
	var (
		force    bool
		list     bool
		output   string
		template string
	)
	cmd := &cobra.Command{
		Use:           "init",
		Short:         "Create a new config.yaml file",
		Long:          "Create a new config.yaml from a template. Run with --list to see available templates.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitE(cmd, list, force, output, template)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing config.yaml")
	cmd.Flags().BoolVarP(&list, "list", "l", false, "List available templates")
	cmd.Flags().StringVarP(&output, "output", "o", "config.yaml", "Output file path")
	cmd.Flags().StringVarP(&template, "template", "t", "", "Template to use")
	return cmd
}

func runInitE(cmd *cobra.Command, list, force bool, output, template string) error {
	if list {
		printTemplateList(cmd.OutOrStdout())
		return nil
	}

	if template == "" {
		fmt.Fprintln(cmd.ErrOrStderr(), "Error: no template specified.")
		fmt.Fprintln(cmd.ErrOrStderr())
		fmt.Fprintln(cmd.ErrOrStderr(), "Use: devcontainerwizard init --template <template>")
		fmt.Fprintln(cmd.ErrOrStderr())
		printTemplateList(cmd.ErrOrStderr())
		return fmt.Errorf("no template specified")
	}

	if _, err := os.Stat(output); err == nil && !force {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: file '%s' already exists — use --force to overwrite\n", output)
		return fmt.Errorf("file exists")
	}

	if err := os.MkdirAll(filepath.Dir(output), 0750); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: creating output directory: %v\n", err)
		return err
	}

	content, err := getTemplateContent(template)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
		return err
	}

	if err := os.WriteFile(output, []byte(content), 0600); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error writing config file: %v\n", err)
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %q from template %q.\n\nNext: devcontainerwizard edit -c %s\n", output, template, output)
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
	fmt.Fprintln(w, "Available templates:")
	for _, name := range templates {
		desc := readTemplateDescription(name)
		fmt.Fprintf(w, "  %-14s %s\n", name, desc)
	}
}

// getTemplateContent returns the raw YAML of the named template.
func getTemplateContent(template string) (string, error) {
	data, err := templatesFS.ReadFile("templates/" + template + ".yaml")
	if err != nil {
		return "", fmt.Errorf("unknown template %q — run 'devcontainerwizard init --list' to see available templates", template)
	}
	return string(data), nil
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
	data, err := templatesFS.ReadFile("templates/" + name + ".yaml")
	if err != nil {
		return ""
	}
	line, _, _ := strings.Cut(string(data), "\n")
	desc, found := strings.CutPrefix(line, "# description: ")
	if !found {
		return ""
	}
	return desc
}
