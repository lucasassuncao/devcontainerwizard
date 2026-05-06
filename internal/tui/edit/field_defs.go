package edit

// FieldDef describes one configurable sub-field within a complex YAML block.
type FieldDef struct {
	Key      string // sub-field name shown in the list
	Desc     string // one-line description
	YAML     string // indented YAML contribution (2-space indent, trailing \n)
	Required bool   // pre-checked when the overlay opens
}

// blockFields maps top-level keys to their toggleable sub-fields.
// Only complex blocks with a fixed, known schema are listed here.
// Simple scalars (name, image, remoteUser, etc.) get a plain textarea instead.
var blockFields = map[string][]FieldDef{
	"build": {
		{
			Key:      "dockerfile",
			Desc:     "Path to Dockerfile",
			YAML:     "  dockerfile: Dockerfile\n",
			Required: true,
		},
		{
			Key:      "context",
			Desc:     "Build context directory",
			YAML:     "  context: .\n",
			Required: true,
		},
		{
			Key:  "args",
			Desc: "Build-time arguments",
			YAML: "  args:\n    MY_ARG: value\n",
		},
		{
			Key:  "target",
			Desc: "Multi-stage build target",
			YAML: "  target: dev\n",
		},
		{
			Key:  "cacheFrom",
			Desc: "Images to use as cache source",
			YAML: "  cacheFrom:\n    - myregistry/image:cache\n",
		},
		{
			Key:  "output",
			Desc: "Build output destination",
			YAML: "  output: type=local,dest=./out\n",
		},
		{
			Key:  "ssh",
			Desc: "SSH agent socket / key mounts",
			YAML: "  ssh:\n    - default\n",
		},
	},

	"customizations": {
		{
			Key:  "vscode",
			Desc: "VS Code extensions and settings",
			YAML: "  vscode:\n    extensions:\n      - ms-python.python\n    settings:\n      editor.formatOnSave: true\n",
		},
		{
			Key:  "jetbrains",
			Desc: "JetBrains IDE plugins",
			YAML: "  jetbrains:\n    plugins:\n      - org.rust.lang\n",
		},
		{
			Key:  "codespaces",
			Desc: "GitHub Codespaces options",
			YAML: "  codespaces:\n    openFiles:\n      - README.md\n",
		},
	},

	"watch": {
		{
			Key:      "waitFor",
			Desc:     "Lifecycle hook to wait for",
			YAML:     "  waitFor:\n    - postCreateCommand\n",
			Required: true,
		},
		{
			Key:  "restart",
			Desc: "Glob patterns that trigger restart",
			YAML: "  restart:\n    - '**/*.go'\n",
		},
	},

	"hostRequirements": {
		{
			Key:      "cpus",
			Desc:     "Minimum number of CPUs",
			YAML:     "  cpus: 4\n",
			Required: true,
		},
		{
			Key:      "memory",
			Desc:     "Minimum memory (e.g. 8gb)",
			YAML:     "  memory: 8gb\n",
			Required: true,
		},
		{
			Key:  "storage",
			Desc: "Minimum disk storage (e.g. 32gb)",
			YAML: "  storage: 32gb\n",
		},
		{
			Key:  "gpu",
			Desc: "GPU requirement (true or object)",
			YAML: "  gpu: true\n",
		},
	},

	"otherPortsAttributes": {
		{
			Key:      "onAutoForward",
			Desc:     "Behavior when port is auto-forwarded",
			YAML:     "  onAutoForward: silent\n",
			Required: true,
		},
		{
			Key:  "label",
			Desc: "Human-readable label",
			YAML: "  label: Other Port\n",
		},
		{
			Key:  "protocol",
			Desc: "Network protocol (http / https)",
			YAML: "  protocol: http\n",
		},
	},
}

// FieldsForKey returns sub-field definitions for a top-level key, or nil when
// the block is a simple scalar and should fall back to a plain textarea.
func FieldsForKey(key string) []FieldDef {
	return blockFields[key]
}
