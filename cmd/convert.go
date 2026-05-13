package cmd

import (
	"fmt"
	"os"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"

	"github.com/spf13/cobra"
)

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

var (
	configFile string
	outputDir  string

	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert config.yaml to .devcontainer/devcontainer.json",
		Long:  "Reads config.yaml (or the file given by --config) and writes a devcontainer.json to the output directory.",
		Run:   runConvert,
	}
)

func init() {
	convertCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	convertCmd.Flags().StringVarP(&outputDir, "output", "o", ".devcontainer", "Output directory")
}

func runConvert(cmd *cobra.Command, args []string) {
	// Load YAML file
	k, err := devcontainer.LoadYAMLFile(configFile)
	if err != nil {
		fatal("Failed to load config: %v", err)
	}

	// Parse to struct
	dc, err := devcontainer.Parse(k)
	if err != nil {
		fatal("Failed to parse config: %v", err)
	}

	// Expand localEnv into containerEnv and remoteEnv
	devcontainer.ExpandLocalEnv(&dc)

	// Validate struct
	if err := devcontainer.Validate(dc); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid devcontainer config:\n%s\n", devcontainer.HumanizeValidationError(err))
		os.Exit(1)
	}

	// Write devcontainer files
	if err := devcontainer.WriteFile(dc, outputDir); err != nil {
		fatal("Failed to write devcontainer: %v", err)
	}
}
