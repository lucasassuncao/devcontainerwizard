package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"

	"github.com/spf13/cobra"
)

// canonicalOutputPath is where convert writes by default, and the location VS
// Code and the devcontainer CLI look for first.
const canonicalOutputPath = ".devcontainer/devcontainer.json"

func convertCmd() *cobra.Command {
	var (
		output string
		force  bool
	)
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert config.yaml to a devcontainer.json file",
		Long:  "Reads config.yaml (or the file given by --config) and writes a devcontainer.json to the path given by --output.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runConvertE(cmd, configPath(cmd), output, force)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", canonicalOutputPath, "Output devcontainer.json file path")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing output file")
	return cmd
}

func runConvertE(cmd *cobra.Command, configFile, output string, force bool) error {
	path, err := convertFile(configFile, output, force)
	if err != nil {
		return err
	}

	clean := filepath.ToSlash(filepath.Clean(path))
	canonical := clean == ".devcontainer/devcontainer.json" ||
		clean == ".devcontainer.json" ||
		strings.HasSuffix(clean, "/.devcontainer/devcontainer.json")
	if !canonical {
		fmt.Fprintf(cmd.ErrOrStderr(), "%s output path is not a canonical devcontainer location.\n", th.Warning.Render("Warning:"))
		fmt.Fprintln(cmd.ErrOrStderr(), th.Dim.Render("VS Code and the devcontainer CLI won't auto-detect this file."))
		fmt.Fprintln(cmd.ErrOrStderr(), th.Dim.Render("Canonical paths: .devcontainer/devcontainer.json or .devcontainer.json"))
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s Saved devcontainer to %s\n", th.Success.Render("✓"), path)
	return nil
}

// convertFile is the convert pipeline the command and the editor's save action
// share: load, parse, validate, write. It returns the path written.
func convertFile(configFile, output string, force bool) (string, error) {
	k, err := devcontainer.LoadYAMLFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	dc, err := devcontainer.Parse(k)
	if err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	if err := devcontainer.Validate(dc); err != nil {
		return "", fmt.Errorf("invalid devcontainer config:\n%s", devcontainer.HumanizeValidationError(err))
	}

	path, err := devcontainer.WriteFile(dc, output, force)
	if err != nil {
		return "", fmt.Errorf("failed to write devcontainer: %w", err)
	}
	return path, nil
}
