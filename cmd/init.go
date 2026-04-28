// Package cmd ...
package cmd

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	// Get available templates dynamically from embedded filesystem
	templates, err := listAvailableTemplates()
	if err != nil {
		// If we can't read templates, something is seriously wrong
		panic(fmt.Sprintf("failed to load templates: %v", err))
	}

	templateList := strings.Join(templates, ", ")

	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config.yaml")
	initCmd.Flags().BoolVarP(&initInteractive, "interactive", "i", false, "Interactive mode with prompts")
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "image", fmt.Sprintf("Template to use (%s)", templateList))
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

	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		fmt.Printf("❌ Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("✅ Created '%s' successfully!\n", configPath)
	fmt.Printf("📝 Edit the file and run 'devcontainerwizard' to generate your devcontainer.json\n")
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
			name := strings.TrimSuffix(entry.Name(), ".yaml")
			templates = append(templates, name)
		}
	}

	return templates, nil
}

// generateInteractiveConfig creates a config.yaml through interactive prompts.
// It parses the selected template into a yaml.Node tree so that modifications
// (name, features, ports) are applied structurally — preserving all comments
// and avoiding the fragility of raw string replacement.
func generateInteractiveConfig() (string, error) {
	fmt.Println("🚀 DevContainer Configuration Generator")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	templates, err := listAvailableTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to list templates: %w", err)
	}

	templateDescriptions := map[string]string{
		"image":         "Image - Minimal config with Docker image",
		"dockerfile":    "Docker - Custom Dockerfile with build config",
		"dockercompose": "Compose - Docker Compose multi-service",
		"full":          "Full - Complete example with all options",
		"golang":        "Golang - Optimized setup for Go development",
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

	// Parse into a yaml.Node tree — this lets us modify fields structurally
	// while keeping every comment from the original template intact.
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}
	root := doc.Content[0] // DocumentNode wraps a MappingNode

	// 1. Customize container name
	if askYesNo("Customize container name?") {
		namePrompt := promptui.Prompt{
			Label:   "Container Name",
			Default: "my-devcontainer",
		}
		name, _ := namePrompt.Run()
		yamlSetScalar(root, "name", name)
	}

	// 2. Add extra features (skipped for the full template which already has them)
	if selectedTemplate != "full" && askYesNo("Add extra features? (docker-in-docker, aws-cli)") {
		yamlAddFeatures(root, []string{
			"ghcr.io/devcontainers/features/docker-in-docker:2",
			"ghcr.io/devcontainers/features/aws-cli:1",
		})
	}

	// 3. Change forwarded port
	if askYesNo("Change default port?") {
		portPrompt := promptui.Prompt{
			Label:   "Port to forward",
			Default: "3000",
		}
		port, _ := portPrompt.Run()
		yamlSetFirstPort(root, port)
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return "", fmt.Errorf("serializing config: %w", err)
	}
	// yaml.Marshal prepends "---\n" for document nodes; strip it for cleaner output.
	return strings.TrimPrefix(string(out), "---\n"), nil
}

// yamlSetScalar finds key in a YAML mapping node and replaces its scalar value.
// No-op if the key is not present.
func yamlSetScalar(mapping *yaml.Node, key, value string) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Value = value
			return
		}
	}
}

// yamlSetFirstPort replaces the first element of the forwardPorts sequence,
// or appends the port if the sequence is empty, or creates the section if absent.
func yamlSetFirstPort(mapping *yaml.Node, port string) {
	portNode := &yaml.Node{Kind: yaml.ScalarNode, Value: port, Tag: "!!int"}

	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == "forwardPorts" {
			seq := mapping.Content[i+1]
			if seq.Kind == yaml.SequenceNode {
				if len(seq.Content) > 0 {
					seq.Content[0] = portNode
				} else {
					seq.Content = append(seq.Content, portNode)
				}
			}
			return
		}
	}

	// Section absent: append forwardPorts with the requested port.
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "forwardPorts", Tag: "!!str"}
	seqNode := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{portNode}}
	mapping.Content = append(mapping.Content, keyNode, seqNode)
}

// yamlAddFeatures adds feature keys (with empty-map values) to the features
// mapping, creating the section if it does not exist.
func yamlAddFeatures(mapping *yaml.Node, featureKeys []string) {
	emptyMap := func() *yaml.Node {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Style: yaml.FlowStyle}
	}

	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == "features" {
			featMap := mapping.Content[i+1]
			if featMap.Kind == yaml.MappingNode {
				for _, k := range featureKeys {
					featMap.Content = append(featMap.Content,
						&yaml.Node{Kind: yaml.ScalarNode, Value: k, Tag: "!!str"},
						emptyMap(),
					)
				}
			}
			return
		}
	}

	// Section absent: create features mapping and append to root.
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "features", Tag: "!!str"}
	featMapNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, k := range featureKeys {
		featMapNode.Content = append(featMapNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: k, Tag: "!!str"},
			emptyMap(),
		)
	}
	mapping.Content = append(mapping.Content, keyNode, featMapNode)
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
