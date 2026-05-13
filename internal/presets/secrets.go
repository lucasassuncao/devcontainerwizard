package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func secretsPresetsMap() map[string]map[string]model.Secret {
	return map[string]map[string]model.Secret{
		"base": {
			"MY_SECRET": {
				Description: "Description of the secret",
				Default:     "",
			},
		},
		"git-tokens": {
			"GITHUB_TOKEN": {Description: "GitHub personal access token", Default: ""},
			"NPM_TOKEN":    {Description: "NPM authentication token", Default: ""},
		},
	}
}

func SecretsPreset(name string) map[string]model.Secret { return secretsPresetsMap()[name] }
func ListSecretsPresets() []string                      { return sortedKeys(secretsPresetsMap()) }
