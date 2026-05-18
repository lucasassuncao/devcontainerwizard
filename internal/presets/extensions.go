package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func featuresPresetsMap() map[string]map[string]map[string]any {
	return map[string]map[string]map[string]any{
		"base": {
			"ghcr.io/devcontainers/features/git:1": {},
		},
		"common-utils": {
			"ghcr.io/devcontainers/features/common-utils:2": {
				"installZsh":      true,
				"installOhMyZsh":  true,
				"upgradePackages": true,
			},
		},
		"docker-in-docker": {
			"ghcr.io/devcontainers/features/docker-in-docker:2": {
				"version": "latest",
				"moby":    true,
			},
		},
		"go-toolchain": {
			"ghcr.io/devcontainers/features/go:1": {
				"version": "1.25",
			},
		},
	}
}

func FeaturesPreset(name string) map[string]map[string]any { return featuresPresetsMap()[name] }
func ListFeaturesPresets() []string                        { return sortedKeys(featuresPresetsMap()) }

func overrideFeatureInstallOrderPresetsMap() map[string][]string {
	return map[string][]string{
		"base": {
			"ghcr.io/devcontainers/features/common-utils",
			"ghcr.io/devcontainers/features/git",
		},
	}
}

func OverrideFeatureInstallOrderPreset(name string) []string {
	return overrideFeatureInstallOrderPresetsMap()[name]
}
func ListOverrideFeatureInstallOrderPresets() []string {
	return sortedKeys(overrideFeatureInstallOrderPresetsMap())
}

func customizationsPresetsMap() map[string]*model.Customizations {
	return map[string]*model.Customizations{
		"base": {
			VSCode: &model.VSCodeCustomization{
				Extensions: []string{"editorconfig.editorconfig"},
				Settings: map[string]any{
					"editor.formatOnSave":      true,
					"files.insertFinalNewline": true,
				},
			},
		},
		"vscode-go": {
			VSCode: &model.VSCodeCustomization{
				Extensions: []string{"golang.go"},
				Settings: map[string]any{
					"go.useLanguageServer": true,
					"go.lintTool":          "golangci-lint",
					"go.formatTool":        "goimports",
					"editor.formatOnSave":  true,
				},
			},
		},
		"vscode-node": {
			VSCode: &model.VSCodeCustomization{
				Extensions: []string{
					"dbaeumer.vscode-eslint",
					"esbenp.prettier-vscode",
				},
				Settings: map[string]any{
					"editor.formatOnSave":     true,
					"editor.defaultFormatter": "esbenp.prettier-vscode",
				},
			},
		},
		"vscode-python": {
			VSCode: &model.VSCodeCustomization{
				Extensions: []string{
					"ms-python.python",
					"ms-python.vscode-pylance",
				},
				Settings: map[string]any{
					"python.linting.enabled":     true,
					"python.formatting.provider": "black",
				},
			},
		},
		"jetbrains-default": {
			JetBrains: &model.JetBrainsCustomization{
				Plugins: []string{"com.intellij.plugins.github"},
			},
		},
	}
}

func CustomizationsPreset(name string) *model.Customizations {
	return customizationsPresetsMap()[name]
}
func ListCustomizationsPresets() []string { return sortedKeys(customizationsPresetsMap()) }
