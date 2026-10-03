package presets

import (
	"strings"
	"testing"
)

// The inline/block decision used to be made by scanning the encoded text for a
// colon, which treated every ordinary image tag as if it were a mapping.
func TestScalarsRenderInline(t *testing.T) {
	cases := []struct {
		field, value, want string
	}{
		{"image", "ubuntu:22.04", "image: ubuntu:22.04\n"},
		{"image", "mcr.microsoft.com/devcontainers/go:latest", "image: mcr.microsoft.com/devcontainers/go:latest\n"},
		{"name", "my-devcontainer", "name: my-devcontainer\n"},
		{"workspaceFolder", "/workspace", "workspaceFolder: /workspace\n"},
		{"updateRemoteUserUID", "true", "updateRemoteUserUID: \"true\"\n"},
	}
	for _, tc := range cases {
		got, err := marshalAsBlock(tc.field, tc.value)
		if err != nil {
			t.Errorf("%s: %v", tc.field, err)
			continue
		}
		if got != tc.want {
			t.Errorf("marshalAsBlock(%q, %q) = %q, want %q", tc.field, tc.value, got, tc.want)
		}
	}
}

// Mappings and sequences must still be indented under the field key.
func TestCompositesRenderAsBlocks(t *testing.T) {
	got, err := marshalAsBlock("containerEnv", map[string]string{"NODE_ENV": "development"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "containerEnv:\n  NODE_ENV: development\n" {
		t.Errorf("mapping rendered as %q", got)
	}

	got, err = marshalAsBlock("capAdd", []string{"SYS_PTRACE"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "capAdd:\n  - SYS_PTRACE") {
		t.Errorf("sequence rendered as %q", got)
	}
}

// Every real preset must render as a block whose first line opens the field.
func TestEveryPresetOpensItsField(t *testing.T) {
	for _, field := range ListFields() {
		for _, name := range ListPresets(field) {
			body, err := PresetYAML(field, name)
			if err != nil {
				t.Errorf("%s/%s: %v", field, name, err)
				continue
			}
			if !strings.HasPrefix(body, field+":") {
				t.Errorf("%s/%s does not start with %q:\n%s", field, name, field+":", body)
			}
			if !strings.HasSuffix(body, "\n") {
				t.Errorf("%s/%s does not end in a newline", field, name)
			}
		}
	}
}
