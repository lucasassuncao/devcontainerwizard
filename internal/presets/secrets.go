package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func secretsPresetsMap() map[string]map[string]model.Secret {
	return map[string]map[string]model.Secret{
		"base": {
			"MY_SECRET": {
				Description: "Description of the secret",
			},
		},
		"git-tokens": {
			"GITHUB_TOKEN": {
				Description:      "GitHub personal access token",
				DocumentationURL: "https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens",
			},
			"NPM_TOKEN": {
				Description:      "NPM authentication token",
				DocumentationURL: "https://docs.npmjs.com/creating-and-viewing-access-tokens",
			},
		},
	}
}

func SecretsPreset(name string) map[string]model.Secret { return secretsPresetsMap()[name] }
func ListSecretsPresets() []string                      { return sortedKeys(secretsPresetsMap()) }
