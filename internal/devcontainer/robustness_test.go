package devcontainer

import (
	"os"
	"path/filepath"
	"testing"
)

// Hostile / malformed configs. Any panic here is a crash the user can trigger
// with a hand-written config.yaml.
var probes = map[string]string{
	"empty":                  ``,
	"only_comment":           `# nothing`,
	"null_doc":               `null`,
	"scalar_doc":             `just-a-string`,
	"seq_doc":                "- a\n- b",
	"name_null":              "name: null",
	"name_int":               "name: 123",
	"name_map":               "name:\n  a: b",
	"mounts_null_item":       "name: x\nimage: y\nmounts:\n  - null",
	"mounts_int_item":        "name: x\nimage: y\nmounts:\n  - 123",
	"mounts_seq_item":        "name: x\nimage: y\nmounts:\n  - [1, 2]",
	"mounts_scalar":          "name: x\nimage: y\nmounts: nope",
	"mounts_readonly":        "name: x\nimage: y\nmounts:\n  - {type: bind, source: a, target: b, readonly: true}",
	"mounts_unknown_key":     "name: x\nimage: y\nmounts:\n  - {type: bind, target: b, bogus: 1}",
	"cmd_int":                "name: x\nimage: y\npostCreateCommand: 123",
	"cmd_seq_of_int":         "name: x\nimage: y\npostCreateCommand: [1, 2]",
	"cmd_map_of_int":         "name: x\nimage: y\npostCreateCommand:\n  a: 123",
	"cmd_map_of_seq_int":     "name: x\nimage: y\npostCreateCommand:\n  a: [1, 2]",
	"cmd_map_of_map":         "name: x\nimage: y\npostCreateCommand:\n  a: {b: c}",
	"cmd_null":               "name: x\nimage: y\npostCreateCommand: null",
	"cmd_nested_seq":         "name: x\nimage: y\npostCreateCommand: [[1]]",
	"gpu_seq":                "name: x\nimage: y\nhostRequirements:\n  gpu: [1, 2]",
	"gpu_int":                "name: x\nimage: y\nhostRequirements:\n  gpu: 5",
	"gpu_null":               "name: x\nimage: y\nhostRequirements:\n  gpu: null",
	"gpu_bad_map":            "name: x\nimage: y\nhostRequirements:\n  gpu:\n    cores: notanint",
	"hostreq_scalar":         "name: x\nimage: y\nhostRequirements: nope",
	"build_scalar":           "name: x\nbuild: nope",
	"build_seq":              "name: x\nbuild: [1]",
	"features_scalar":        "name: x\nimage: y\nfeatures: nope",
	"features_deep_null":     "name: x\nimage: y\nfeatures:\n  a:\n    b: null",
	"forwardports_map":       "name: x\nimage: y\nforwardPorts:\n  a: b",
	"portsattrs_scalar":      "name: x\nimage: y\nportsAttributes: nope",
	"portsattrs_null_val":    "name: x\nimage: y\nportsAttributes:\n  \"3000\": null",
	"otherportsattrs_scalar": "name: x\nimage: y\notherPortsAttributes: nope",
	"customizations_seq":     "name: x\nimage: y\ncustomizations: [1]",
	"secrets_seq":            "name: x\nimage: y\nsecrets: [1]",
	"watch_scalar":           "name: x\nimage: y\nwatch: nope",
	"compose_no_ws":          "name: x\ndockerComposeFile: [a.yml]\nservice: s",
	"compose_scalar":         "name: x\ndockerComposeFile: a.yml\nservice: s\nworkspaceFolder: /w",
	"deep_alias_bomb":        "name: &a x\nimage: *a",
	"tab_indent":             "name: x\n\timage: y",
	"dup_keys":               "name: x\nname: z\nimage: y",
	"huge_int":               "name: x\nimage: y\nforwardPorts: [99999999999999999999]",
	"neg_cpus":               "name: x\nimage: y\nhostRequirements:\n  cpus: -5",
}

func TestConvertPipelineNeverPanics(t *testing.T) {
	for name, body := range probes {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			cfg := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(cfg, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}

			// A panic anywhere in load → parse → validate → write fails the test
			// with the offending input named.
			k, err := LoadYAMLFile(cfg)
			if err != nil {
				return
			}
			dc, err := Parse(k)
			if err != nil {
				return
			}
			if err := Validate(dc); err != nil {
				_ = HumanizeValidationError(err)
				return
			}
			if _, err := WriteFile(dc, filepath.Join(dir, "out.json"), true); err != nil {
				return
			}
		})
	}
}
