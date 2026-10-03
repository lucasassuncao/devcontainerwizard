package devcontainer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
)

func parse(t *testing.T, yaml string) error {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	k, err := devcontainer.LoadYAMLFile(path)
	if err != nil {
		t.Fatalf("LoadYAMLFile: %v", err)
	}
	_, parseErr := devcontainer.Parse(k)
	return parseErr
}

// TestMountObjectRejectsUnknownKeys pins the specification's Mount definition,
// which is "additionalProperties": false over exactly type, source and target.
// Before this was enforced, an unrecognised key was dropped in silence - and
// for "readonly" that turned a mount the user wrote as read-only into a
// read-write one, with nothing printed.
func TestMountObjectRejectsUnknownKeys(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want string
	}{{
		name: "readonly names the string-form replacement",
		yaml: "name: x\nimage: alpine\nmounts:\n  - type: bind\n    source: ./c\n    target: /c\n    readonly: true\n",
		want: "--mount string",
	}, {
		name: "a misspelled key is rejected rather than dropped",
		yaml: "name: x\nimage: alpine\nmounts:\n  - type: bind\n    source: ./c\n    taget: /c\n",
		want: "taget",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parse(t, tt.yaml)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error does not mention %q: %v", tt.want, err)
			}
		})
	}
}

// TestMountStringFormIsAccepted covers the replacement the readonly error
// points at: the specification allows a raw Docker --mount string alongside
// the object form, and that string carries the options the object cannot.
func TestMountStringFormIsAccepted(t *testing.T) {
	yaml := "name: x\nimage: alpine\nmounts:\n  - source=./configs,target=/etc/app/config,type=bind,readonly\n"
	if err := parse(t, yaml); err != nil {
		t.Fatalf("string-form mount rejected: %v", err)
	}
}

// TestMountObjectRoundTrips guards the fields that did survive.
func TestMountObjectRoundTrips(t *testing.T) {
	yaml := "name: x\nimage: alpine\nmounts:\n  - type: volume\n    source: cache\n    target: /cache\n"
	if err := parse(t, yaml); err != nil {
		t.Fatalf("valid mount rejected: %v", err)
	}
}

// TestBareSizesDecodeAsStrings pins hostRequirements sizes without a unit: YAML
// reads them as ints, and convert refused a byte count the specification allows.
func TestBareSizesDecodeAsStrings(t *testing.T) {
	k, err := devcontainer.LoadYAML(strings.NewReader("name: x\nimage: alpine\nhostRequirements:\n  memory: 512\n  storage: 64\n"))
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	dc, err := devcontainer.Parse(k)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := dc.HostRequirements.Memory; got != "512" {
		t.Errorf("memory = %q, want %q", got, "512")
	}
	if got := dc.HostRequirements.Storage; got != "64" {
		t.Errorf("storage = %q, want %q", got, "64")
	}
}
