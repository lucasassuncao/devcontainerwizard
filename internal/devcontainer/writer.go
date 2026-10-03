package devcontainer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

// WriteFile serialises dc as JSON to outputPath. Parent directories are created
// as needed. When force is false, returns an error if outputPath already exists.
// Returns the cleaned absolute-or-relative path actually written.
func WriteFile(dc model.DevContainer, outputPath string, force bool) (string, error) {
	if dc.Schema == "" {
		dc.Schema = "https://containers.dev/implementors/json_schema"
	}

	jsonBytes, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Check the refusal before creating anything. An existing file implies an
	// existing parent, so this order changes no outcome today - it is here so a
	// later edit cannot introduce a side effect ahead of the guard. It also
	// matches the order the init command already uses.
	if _, err := os.Stat(outputPath); err == nil && !force {
		return "", fmt.Errorf("file '%s' already exists — use --force to overwrite", outputPath)
	}

	parent := filepath.Dir(outputPath)
	if err := os.MkdirAll(parent, 0o750); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outputPath, jsonBytes, 0o600); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}

	return outputPath, nil
}
