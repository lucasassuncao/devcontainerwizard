package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"

	"github.com/spf13/cobra"
)

var convertCmd = newConvertCmd()

func newConvertCmd() *cobra.Command {
	var (
		configFile string
		output     string
		force      bool
	)
	cmd := &cobra.Command{
		Use:           "convert",
		Short:         "Convert config.yaml to a devcontainer.json file",
		Long:          "Reads config.yaml (or the file given by --config) and writes a devcontainer.json to the path given by --output.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConvertE(cmd, configFile, output, force)
		},
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	cmd.Flags().StringVarP(&output, "output", "o", ".devcontainer/devcontainer.json", "Output devcontainer.json file path")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing output file")
	return cmd
}

func runConvertE(cmd *cobra.Command, configFile, output string, force bool) error {
	k, err := devcontainer.LoadYAMLFile(configFile)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: failed to load config: %v\n", err)
		return err
	}

	dc, err := devcontainer.Parse(k)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: failed to parse config: %v\n", err)
		return err
	}

	if err := devcontainer.Validate(dc); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Invalid devcontainer config:\n%s\n", devcontainer.HumanizeValidationError(err))
		return err
	}

	path, err := devcontainer.WriteFile(dc, output, force)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: failed to write devcontainer: %v\n", err)
		return err
	}

	clean := filepath.ToSlash(filepath.Clean(path))
	canonical := clean == ".devcontainer/devcontainer.json" ||
		clean == ".devcontainer.json" ||
		strings.HasSuffix(clean, "/.devcontainer/devcontainer.json")
	if !canonical {
		fmt.Fprintln(cmd.ErrOrStderr(), "Warning: output path is not a canonical devcontainer location.")
		fmt.Fprintln(cmd.ErrOrStderr(), "VS Code and the devcontainer CLI won't auto-detect this file.")
		fmt.Fprintln(cmd.ErrOrStderr(), "Canonical paths: .devcontainer/devcontainer.json or .devcontainer.json")
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Saved devcontainer to %s\n", path)
	return nil
}
