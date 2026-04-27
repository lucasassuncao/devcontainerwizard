package devcontainer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

func WriteFile(dc model.DevContainer, outputDir string) error {
	jsonBytes, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "devcontainer.json")
	if err := os.WriteFile(filePath, jsonBytes, 0600); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fmt.Printf("Saved devcontainer to %s\n", filePath)
	return nil
}
