package devcontainer

// FieldSnippets returns the per-parent → per-child YAML snippet map used by
// the editor TUI when the user toggles a sub-field on. Keeping this in the
// devcontainer package keeps yedit free of devcontainer-specific knowledge.
//
// Each snippet is the indented YAML chunk inserted under the parent key,
// already including the child key, a placeholder value, and a trailing newline.
func FieldSnippets() map[string]map[string]string {
	return map[string]map[string]string{
		"build": {
			"dockerfile": "  dockerfile: Dockerfile\n",
			"context":    "  context: .\n",
			"args":       "  args:\n    MY_ARG: value\n",
			"target":     "  target: dev\n",
			"cacheFrom":  "  cacheFrom:\n    - myregistry/image:cache\n",
			"options":    "  options:\n    - --no-cache\n",
		},
		"customizations": {
			"vscode":     "  vscode:\n    extensions:\n      - ms-python.python\n    settings:\n      editor.formatOnSave: true\n",
			"jetbrains":  "  jetbrains:\n    plugins:\n      - org.rust.lang\n",
			"codespaces": "  codespaces:\n    extensions:\n      - github.copilot\n    settings:\n      editor.formatOnSave: true\n",
		},
		"watch": {
			"waitFor": "  waitFor:\n    - postCreateCommand\n",
			"restart": "  restart:\n    - '**/*.go'\n",
		},
		"hostRequirements": {
			"cpus":    "  cpus: 4\n",
			"memory":  "  memory: 8gb\n",
			"storage": "  storage: 32gb\n",
			"gpu":     "  gpu: true\n",
		},
		"mounts": {
			"type":     "  - type: bind\n",
			"source":   "    source: ${localWorkspaceFolder}/.cache\n",
			"target":   "    target: /home/vscode/.cache\n",
			"readonly": "    readonly: false\n",
		},
		"portsAttributes": {
			"label":            "  \"3000\":\n    label: Web App\n",
			"onAutoForward":    "    onAutoForward: notify\n",
			"protocol":         "    protocol: http\n",
			"elevateIfNeeded":  "    elevateIfNeeded: false\n",
			"requireLocalPort": "    requireLocalPort: false\n",
		},
		"secrets": {
			"description":      "  MY_SECRET:\n    description: Description of the secret\n",
			"documentationUrl": "    documentationUrl: https://example.com/docs/secrets\n",
		},
		"otherPortsAttributes": {
			"onAutoForward":    "  onAutoForward: silent\n",
			"label":            "  label: Other Port\n",
			"protocol":         "  protocol: http\n",
			"elevateIfNeeded":  "  elevateIfNeeded: false\n",
			"requireLocalPort": "  requireLocalPort: false\n",
		},
	}
}

// PreCheckedFields returns the per-parent list of sub-fields that should be
// pre-checked when the editor opens an overlay for that parent. These are UX
// hints — they are *not* the same as validation-required fields.
func PreCheckedFields() map[string][]string {
	return map[string][]string{
		"build":                {"dockerfile", "context"},
		"watch":                {"waitFor"},
		"hostRequirements":     {"cpus", "memory"},
		"mounts":               {"type", "target"},
		"portsAttributes":      {"label", "onAutoForward"},
		"otherPortsAttributes": {"onAutoForward"},
		"secrets":              {"description"},
	}
}
