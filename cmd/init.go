// Package cmd ...
package cmd

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

//go:embed templates/*.yaml
var templatesFS embed.FS

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create a new config.yaml file",
		Long:  "Interactive command to create a new config.yaml file with common configurations",
		Run:   runInit,
	}

	// Flags for init command
	initForce       bool
	initInteractive bool
	initTemplate    string
)

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config.yaml")
	initCmd.Flags().BoolVarP(&initInteractive, "interactive", "i", false, "Interactive mode with prompts")
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "image", "Template to use (image, dockerfile, dockercompose, full)")
}

func runInit(cmd *cobra.Command, args []string) {
	configPath := "config.yaml"

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil && !initForce {
		fmt.Printf("❌ File '%s' already exists. Use --force to overwrite.\n", configPath)
		return
	}

	var content string
	var err error

	if initInteractive {
		content, err = generateInteractiveConfig()
		if err != nil {
			fmt.Printf("❌ Error generating config: %v\n", err)
			return
		}
	} else {
		content, err = getTemplateContent(initTemplate)
		if err != nil {
			fmt.Printf("❌ Error loading template: %v\n", err)
			return
		}
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		fmt.Printf("❌ Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("✅ Created '%s' successfully!\n", configPath)
	fmt.Printf("📝 Edit the file and run 'devcontainer' to generate your devcontainer.json\n")
}

// getTemplateContent retrieves the content of the specified template from the embedded filesystem
func getTemplateContent(template string) (string, error) {
	filename := fmt.Sprintf("templates/%s.yaml", template)

	data, err := templatesFS.ReadFile(filename)
	if err != nil {
		// Fallback para image se template não encontrado
		if template != "image" {
			fmt.Printf("⚠️  Template '%s' not found, using 'image'\n", template)
			return getTemplateContent("image")
		}
		return "", fmt.Errorf("template not found: %s", template)
	}

	return string(data), nil
}

// listAvailableTemplates lists all available templates in the embedded filesystem
func listAvailableTemplates() ([]string, error) {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			// Remove a extensão .yaml
			name := strings.TrimSuffix(entry.Name(), ".yaml")
			templates = append(templates, name)
		}
	}

	return templates, nil
}

// generateInteractiveConfig creates a config.yaml content through interactive prompts
func generateInteractiveConfig() (string, error) {
	fmt.Println("🚀 DevContainer Configuration Generator")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 1. Select template
	templates, err := listAvailableTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to list templates: %w", err)
	}

	templateDescriptions := map[string]string{
		"image":         "Image - Minimal config with Docker image",
		"dockerfile":    "Docker - Custom Dockerfile with build config",
		"dockercompose": "Compose - Docker Compose multi-service",
		"full":          "Full - Complete example with all options",
	}

	var templateItems []string
	for _, t := range templates {
		desc := templateDescriptions[t]
		if desc == "" {
			desc = t
		}
		templateItems = append(templateItems, desc)
	}

	templatePrompt := promptui.Select{
		Label: "Choose a base template",
		Items: templateItems,
		Size:  len(templateItems),
	}

	idx, _, err := templatePrompt.Run()
	if err != nil {
		return "", err
	}

	selectedTemplate := templates[idx]
	content, err := getTemplateContent(selectedTemplate)
	if err != nil {
		return "", err
	}

	// 2. Customize devcontainer name
	if askYesNo("Customize container name?") {
		namePrompt := promptui.Prompt{
			Label:   "Container Name",
			Default: "my-devcontainer",
		}
		name, _ := namePrompt.Run()
		content = strings.Replace(content, "name: my-devcontainer", fmt.Sprintf("name: %s", name), 1)
		content = strings.Replace(content, "name: full-devcontainer", fmt.Sprintf("name: %s", name), 1)
	}

	// 3. Add extra features (if not full template)
	if selectedTemplate != "full" && askYesNo("Add extra features? (docker-in-docker, aws-cli)") {
		content = addExtraFeatures(content)
	}

	// 4. Customize port
	if askYesNo("Change default port?") {
		portPrompt := promptui.Prompt{
			Label:   "Port to forward",
			Default: "3000",
		}
		port, _ := portPrompt.Run()
		content = strings.Replace(content, "- 3000", fmt.Sprintf("- %s", port), 1)
	}

	return content, nil
}

func addExtraFeatures(content string) string {
	extraFeatures := `
  ghcr.io/devcontainers/features/docker-in-docker:2: {}
  ghcr.io/devcontainers/features/aws-cli:1: {}`

	// Searches for the features section and adds to it
	if strings.Contains(content, "features:") {
		content = strings.Replace(content,
			"ghcr.io/devcontainers/features/github-cli:1: {}",
			"ghcr.io/devcontainers/features/github-cli:1: {}"+extraFeatures,
			1)
	} else {
		// If not present, add the features section
		content += "\n# Additional features\nfeatures:" + extraFeatures + "\n"
	}

	return content
}

func askYesNo(question string) bool {
	prompt := promptui.Select{
		Label: question,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false
	}
	return result == "Yes"
}
