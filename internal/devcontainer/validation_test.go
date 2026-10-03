package devcontainer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// accepts runs the full convert pipeline over body and returns the humanized
// validation error, or "" when the config is accepted.
func accepts(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	cfg := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfg, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	k, err := LoadYAMLFile(cfg)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	dc, err := Parse(k)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := Validate(dc); err != nil {
		return HumanizeValidationError(err)
	}
	return ""
}

// portsAttributes and secrets are maps: without "dive" on the field, the
// validator never descends into their values and every per-field rule declared
// on PortAttributes and Secret was dead. An invalid onAutoForward then reached
// devcontainer.json, where the IDE - not this tool - reported it.
func TestValidationDescendsIntoMaps(t *testing.T) {
	cases := map[string]struct{ body, want string }{
		"portsAttributes.onAutoForward": {
			"name: x\nimage: y\nportsAttributes:\n  \"3000\":\n    onAutoForward: bogus",
			"OnAutoForward",
		},
		"portsAttributes.protocol": {
			"name: x\nimage: y\nportsAttributes:\n  \"3000\":\n    protocol: ftp",
			"Protocol",
		},
		"secrets.documentationUrl": {
			"name: x\nimage: y\nsecrets:\n  A:\n    documentationUrl: not-a-url",
			"DocumentationURL",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := accepts(t, tc.body)
			if got == "" {
				t.Fatalf("config was accepted; expected a complaint about %s", tc.want)
			}
			if !strings.Contains(got, tc.want) {
				t.Errorf("error %q does not mention %s", got, tc.want)
			}
		})
	}
}

// A YAML null in either container decodes to an empty entry that dive skips and
// json.Marshal renders as null - a devcontainer.json the IDE refuses.
func TestEmptyEntriesAreRejected(t *testing.T) {
	cases := map[string]string{
		"mounts":          "name: x\nimage: y\nmounts:\n  - null",
		"portsAttributes": "name: x\nimage: y\nportsAttributes:\n  \"3000\": null",
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			got := accepts(t, body)
			if got == "" {
				t.Fatal("a null entry was accepted and would be written as JSON null")
			}
			if !strings.Contains(got, "empty entry") {
				t.Errorf("unexpected message: %q", got)
			}
		})
	}
}

// Every tag the model actually uses needs an entry in validationMessages;
// without one the user sees the raw tag name.
func TestValidationMessagesCoverModelTags(t *testing.T) {
	cases := map[string]string{
		"min":           "name: x\nimage: y\nhostRequirements:\n  cpus: -5",
		"url":           "name: x\nimage: y\nsecrets:\n  A:\n    documentationUrl: not-a-url",
		"required_with": "name: x\ndockerComposeFile: [a.yml]\nworkspaceFolder: /w",
		"oneof":         "name: x\nimage: y\nuserEnvProbe: bogus",
		"required":      "image: y",
	}
	for tag, body := range cases {
		t.Run(tag, func(t *testing.T) {
			got := accepts(t, body)
			if got == "" {
				t.Fatalf("expected a %s violation", tag)
			}
			if strings.Contains(got, "failed validation") {
				t.Errorf("tag %q has no humanized message: %q", tag, got)
			}
			if strings.Contains(got, "%!") {
				t.Errorf("message template mis-formatted: %q", got)
			}
		})
	}
}
