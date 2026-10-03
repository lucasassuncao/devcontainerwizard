package cmd

import (
	"strings"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
	dcpresets "github.com/lucasassuncao/devcontainerwizard/internal/presets"
	"github.com/lucasassuncao/yedit/document"
	"github.com/lucasassuncao/yedit/schema"
	"github.com/lucasassuncao/yedit/validate"
)

// wireValidators builds the exact validator set the edit command runs, against
// the same schema tree and metadata source the TUI uses.
func wireValidators(t *testing.T) validate.WiredValidators {
	t.Helper()
	src, err := buildDevContainerHints()
	if err != nil {
		t.Fatalf("buildDevContainerHints: %v", err)
	}
	return validate.WireWithSchema(DevContainerValidators, schema.Discover(&model.DevContainer{}), src)
}

func validateYAML(t *testing.T, raw string) []string {
	t.Helper()
	blocks, err := document.ParseBlocks([]byte(raw))
	if err != nil {
		t.Fatalf("ParseBlocks: %v", err)
	}
	var msgs []string
	for _, v := range validate.RunAll(wireValidators(t), []byte(raw), blocks) {
		msgs = append(msgs, v.Path+": "+v.Message)
	}
	return msgs
}

// TestShippedTemplatesPassValidators is the guard that keeps metadata
// constraints honest: every template the init command writes, and the sample
// config at the repository root, must survive the same rules the editor runs
// on save. A constraint declared too tightly fails here rather than in a
// user's terminal.
func TestShippedTemplatesPassValidators(t *testing.T) {
	names, err := listAvailableTemplates()
	if err != nil {
		t.Fatalf("listAvailableTemplates: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("no templates found")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			raw, err := getTemplateContent(name)
			if err != nil {
				t.Fatalf("getTemplateContent(%q): %v", name, err)
			}
			if msgs := validateYAML(t, raw); len(msgs) > 0 {
				t.Errorf("template %q fails validation:\n  %s", name, strings.Join(msgs, "\n  "))
			}
		})
	}
}

// presetShell returns the surrounding document a preset for field is spliced
// into. Presets are fragments, so the shell has to supply a container source -
// except for the fields that are one, or that only mean anything with Compose,
// where the wrong shell would manufacture a conflict that is not the preset's
// fault.
func presetShell(field string) string {
	const compose = "name: smoke\ndockerComposeFile:\n  - docker-compose.yml\nworkspaceFolder: /w\n"
	switch field {
	case "name":
		return "image: ubuntu:22.04\n"
	case "image", "build":
		return "name: smoke\n"
	case "dockerComposeFile":
		// The preset supplies dockerComposeFile itself, and the specification
		// requires service and workspaceFolder to travel with it.
		return "name: smoke\nservice: app\nworkspaceFolder: /w\n"
	case "service":
		return compose
	case "runServices":
		return compose + "service: app\n"
	default:
		return "name: smoke\nimage: ubuntu:22.04\n"
	}
}

// TestPresetsPassValidators is the preset counterpart of
// TestShippedTemplatesPassValidators: every value the editor's preset picker
// and show-examples offer must survive the rules the editor runs on save.
// Without it a preset could insert a value the same session then refuses to
// save - and the enum tightening for otherPortsAttributes is exactly the kind
// of change that would cause it.
func TestPresetsPassValidators(t *testing.T) {
	fields := dcpresets.ListFields()
	if len(fields) == 0 {
		t.Fatal("no preset fields found")
	}
	for _, field := range fields {
		for _, name := range dcpresets.ListPresets(field) {
			t.Run(field+"/"+name, func(t *testing.T) {
				body, err := dcpresets.PresetYAML(field, name)
				if err != nil {
					t.Fatalf("PresetYAML: %v", err)
				}
				if msgs := validateYAML(t, presetShell(field)+body); len(msgs) > 0 {
					t.Errorf("preset fails validation:\n  %s", strings.Join(msgs, "\n  "))
				}
			})
		}
	}
}

// TestSpecConstraintsAreEnforced pins the constraints taken from the Dev
// Containers specification. Each case is a document that is invalid for
// exactly one reason, so a constraint silently dropped from the metadata
// fails here.
func TestSpecConstraintsAreEnforced(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want string
	}{{
		name: "hostRequirements.memory rejects an uppercase unit",
		yaml: "name: x\nimage: alpine\nhostRequirements:\n  memory: 8GB\n",
		want: "memory",
	}, {
		name: "hostRequirements.memory rejects a non-numeric size",
		yaml: "name: x\nimage: alpine\nhostRequirements:\n  memory: lots\n",
		want: "memory",
	}, {
		name: "hostRequirements.cpus rejects zero",
		yaml: "name: x\nimage: alpine\nhostRequirements:\n  cpus: 0\n",
		want: "cpus",
	}, {
		name: "secrets.documentationUrl rejects a non-URL",
		yaml: "name: x\nimage: alpine\nsecrets:\n  TOKEN:\n    documentationUrl: not-a-url\n",
		want: "documentationUrl",
	}, {
		name: "forwardPorts rejects duplicates",
		yaml: "name: x\nimage: alpine\nforwardPorts:\n  - 3000\n  - 3000\n",
		want: "forwardPorts",
	}, {
		name: "dockerComposeFile rejects an empty list",
		yaml: "name: x\ndockerComposeFile: []\nservice: app\nworkspaceFolder: /w\n",
		want: "dockerComposeFile",
	}, {
		name: "compose without workspaceFolder",
		yaml: "name: x\ndockerComposeFile:\n  - docker-compose.yml\nservice: app\n",
		want: "workspaceFolder",
	}, {
		name: "stopCompose on a non-compose config",
		yaml: "name: x\nimage: alpine\nshutdownAction: stopCompose\n",
		want: "shutdownAction",
	}, {
		name: "stopContainer on a compose config",
		yaml: "name: x\ndockerComposeFile:\n  - docker-compose.yml\nservice: app\nworkspaceFolder: /w\nshutdownAction: stopContainer\n",
		want: "shutdownAction",
	}, {
		name: "openBrowserOnce is not valid in otherPortsAttributes",
		yaml: "name: x\nimage: alpine\notherPortsAttributes:\n  onAutoForward: openBrowserOnce\n",
		want: "onAutoForward",
	}, {
		name: "name is required",
		yaml: "image: alpine\n",
		want: "name",
	}, {
		name: "image and build are mutually exclusive",
		yaml: "name: x\nimage: alpine\nbuild:\n  dockerfile: Dockerfile\n",
		want: "",
	}, {
		name: "service requires dockerComposeFile",
		yaml: "name: x\nimage: alpine\nservice: app\n",
		want: "",
	}, {
		name: "dockerFile is reported as deprecated",
		yaml: "name: x\ndockerFile: Dockerfile\n",
		want: "dockerFile",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgs := validateYAML(t, tt.yaml)
			if len(msgs) == 0 {
				t.Fatalf("expected a violation, got none")
			}
			if tt.want == "" {
				return // presence of any violation is the assertion
			}
			joined := strings.Join(msgs, "\n")
			if !strings.Contains(joined, tt.want) {
				t.Errorf("no violation mentioning %q:\n  %s", tt.want, joined)
			}
		})
	}
}

// TestEditorRefusesWhatConvertRefuses keeps ctrl+l and save honest: a document
// convert would refuse must not be saved as valid, and no field is reported twice.
func TestEditorRefusesWhatConvertRefuses(t *testing.T) {
	const compose = "name: x\ndockerComposeFile:\n  - docker-compose.yml\n"
	tests := []struct {
		name string
		yaml string
		want string
	}{
		{"compose without service", compose + "workspaceFolder: /w\n", "service"},
		{"compose without workspaceFolder", compose + "service: app\n", "workspaceFolder"},
		{"empty mount entry", "name: x\nimage: alpine\nmounts:\n  -\n", "mounts"},
		{"empty portsAttributes entry", "name: x\nimage: alpine\nportsAttributes:\n  \"3000\":\n", "portsAttributes"},
		{"port above 65535", "name: x\nimage: alpine\nforwardPorts:\n  - 70000\n", "forwardPorts"},
		{"port that is a bool", "name: x\nimage: alpine\nforwardPorts:\n  - true\n", "forwardPorts"},
		{"appPort zero", "name: x\nimage: alpine\nappPort:\n  - 0\n", "appPort"},
		{"value convert cannot decode", "name: x\nimage: alpine\ncapAdd:\n  key: value\n", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !convertRejects(t, tt.yaml) {
				t.Fatal("convert accepts this document; the case no longer tests anything")
			}
			msgs := validateYAML(t, tt.yaml)
			if len(msgs) == 0 {
				t.Fatal("convert refuses this document but the editor would save it")
			}
			joined := strings.Join(msgs, "\n")
			if tt.want != "" && !strings.Contains(joined, tt.want) {
				t.Errorf("no violation mentioning %q:\n  %s", tt.want, joined)
			}
			seen := map[string]bool{}
			for _, m := range msgs {
				path, _, _ := strings.Cut(m, ": ")
				if path != "" && seen[path] {
					t.Errorf("%q reported more than once:\n  %s", path, joined)
				}
				seen[path] = true
			}
		})
	}
}

// convertRejects reports whether convert's own pipeline refuses raw.
func convertRejects(t *testing.T, raw string) bool {
	t.Helper()
	k, err := devcontainer.LoadYAML(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	dc, err := devcontainer.Parse(k)
	if err != nil {
		return true
	}
	return devcontainer.Validate(dc) != nil
}

// TestSpecValidValuesAreAccepted covers the values the specification allows
// that an over-tight constraint would wrongly reject.
func TestSpecValidValuesAreAccepted(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{{
		name: "shutdownAction accepts stopCompose",
		yaml: "name: x\ndockerComposeFile:\n  - docker-compose.yml\nservice: app\nworkspaceFolder: /workspaces/app\nshutdownAction: stopCompose\n",
	}, {
		name: "openBrowserOnce is valid in portsAttributes",
		yaml: "name: x\nimage: alpine\nportsAttributes:\n  \"3000\":\n    onAutoForward: openBrowserOnce\n",
	}, {
		name: "build without context",
		yaml: "name: x\nbuild:\n  dockerfile: Dockerfile\n",
	}, {
		name: "secret without a description",
		yaml: "name: x\nimage: alpine\nsecrets:\n  TOKEN:\n    documentationUrl: https://example.com/tokens\n",
	}, {
		name: "hostRequirements sizes with and without units",
		yaml: "name: x\nimage: alpine\nhostRequirements:\n  cpus: 1\n  memory: 512mb\n  storage: 64\n",
	}, {
		name: "forwardPorts mixing numbers and host:port strings",
		yaml: "name: x\nimage: alpine\nforwardPorts:\n  - 3000\n  - \"db:5432\"\n",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if msgs := validateYAML(t, tt.yaml); len(msgs) > 0 {
				t.Errorf("expected no violations, got:\n  %s", strings.Join(msgs, "\n  "))
			}
		})
	}
}
