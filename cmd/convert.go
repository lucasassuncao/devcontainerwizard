// Package cmd ...
package cmd

import (
	"fmt"
	"log"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"

	"github.com/spf13/cobra"
)

var (
	configFile string
	outputDir  string
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".devcontainer", "Output directory")
}

func runConvert(cmd *cobra.Command, args []string) {
	// Load YAML file
	k, err := devcontainer.LoadYAMLFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse to struct
	dc, err := devcontainer.Parse(k)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Validate struct
	if err := devcontainer.Validate(dc); err != nil {
		fmt.Printf("Invalid devcontainer config:\n%s\n", devcontainer.HumanizeValidationError(err))
		return
	}

	// Write devcontainer files
	if err := devcontainer.WriteFile(dc, outputDir); err != nil {
		log.Fatalf("Failed to write devcontainer: %v", err)
	}
}
