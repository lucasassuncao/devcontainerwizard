package devcontainer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

func mustParse(t *testing.T, body string) model.DevContainer {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "c.yaml")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	k, err := LoadYAMLFile(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	dc, err := Parse(k)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return dc
}

// The decode chain used to split every scalar on commas. A comma is legal
// inside a run argument, a device path and a Docker --mount string, so that
// turned one value into several broken ones.
func TestScalarPromotionKeepsCommas(t *testing.T) {
	t.Run("run argument stays whole", func(t *testing.T) {
		dc := mustParse(t, "name: x\nimage: y\nrunArgs: \"--env=A,B\"")
		if len(dc.RunArgs) != 1 || dc.RunArgs[0] != "--env=A,B" {
			t.Fatalf("runArgs = %#v, want [\"--env=A,B\"]", dc.RunArgs)
		}
	})

	t.Run("mount string is not chopped up", func(t *testing.T) {
		dc := mustParse(t, "name: x\nimage: y\nmounts: \"source=/a,target=/b,type=bind\"")
		if len(dc.Mounts) != 1 || dc.Mounts[0].Str != "source=/a,target=/b,type=bind" {
			t.Fatalf("mounts = %#v, want one intact mount string", dc.Mounts)
		}
	})

	t.Run("scalar still promotes to a one-element list", func(t *testing.T) {
		dc := mustParse(t, "name: x\ndockerComposeFile: docker-compose.yml\nservice: s\nworkspaceFolder: /w")
		if len(dc.DockerComposeFile) != 1 || dc.DockerComposeFile[0] != "docker-compose.yml" {
			t.Fatalf("dockerComposeFile = %#v", dc.DockerComposeFile)
		}
	})

	t.Run("list form is untouched", func(t *testing.T) {
		dc := mustParse(t, "name: x\ndockerComposeFile: [a.yml, b.yml]\nservice: s\nworkspaceFolder: /w")
		if len(dc.DockerComposeFile) != 2 {
			t.Fatalf("dockerComposeFile = %#v", dc.DockerComposeFile)
		}
	})
}

// forwardPorts and appPort are []any, which no struct tag can constrain.
func TestPortListValidation(t *testing.T) {
	bad := map[string]string{
		"above range":  "name: x\nimage: y\nforwardPorts: [70000]",
		"zero":         "name: x\nimage: y\nforwardPorts: [0]",
		"negative":     "name: x\nimage: y\nforwardPorts: [-1]",
		"huge":         "name: x\nimage: y\nforwardPorts: [99999999999999999999]",
		"fraction":     "name: x\nimage: y\nforwardPorts: [3000.5]",
		"boolean":      "name: x\nimage: y\nforwardPorts: [true]",
		"empty string": "name: x\nimage: y\nforwardPorts: [\"\"]",
		"appPort too":  "name: x\nimage: y\nappPort: [70000]",
	}
	for name, body := range bad {
		t.Run("rejects "+name, func(t *testing.T) {
			got := accepts(t, body)
			if got == "" {
				t.Fatal("invalid port list was accepted")
			}
			if !strings.Contains(got, "65535") {
				t.Errorf("unexpected message: %q", got)
			}
		})
	}

	good := map[string]string{
		"integer":    "name: x\nimage: y\nforwardPorts: [3000]",
		"boundaries": "name: x\nimage: y\nforwardPorts: [1, 65535]",
		"host:port":  "name: x\nimage: y\nforwardPorts: [\"db:5432\"]",
		"mixed":      "name: x\nimage: y\nforwardPorts: [3000, \"db:5432\"]",
	}
	for name, body := range good {
		t.Run("accepts "+name, func(t *testing.T) {
			if got := accepts(t, body); got != "" {
				t.Fatalf("valid port list rejected: %s", got)
			}
		})
	}
}

func TestWriteFileOverwriteRules(t *testing.T) {
	dc := mustParse(t, "name: x\nimage: y")
	dir := t.TempDir()
	target := filepath.Join(dir, ".devcontainer", "devcontainer.json")

	if _, err := WriteFile(dc, target, false); err != nil {
		t.Fatalf("first write: %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("parent directory was not created: %v", err)
	}
	if _, err := WriteFile(dc, target, false); err == nil {
		t.Fatal("overwrite without force was allowed")
	}
	if _, err := WriteFile(dc, target, true); err != nil {
		t.Fatalf("overwrite with force: %v", err)
	}
}

// LoadYAML is the reader-based seam: parsing must not require a file on disk.
func TestLoadYAMLFromReader(t *testing.T) {
	k, err := LoadYAML(strings.NewReader("name: from-memory\nimage: ubuntu\n"))
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	dc, err := Parse(k)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if dc.Name != "from-memory" || dc.Image != "ubuntu" {
		t.Fatalf("decoded %+v", dc)
	}
	if err := Validate(dc); err != nil {
		t.Fatalf("Validate: %s", HumanizeValidationError(err))
	}
}

func TestLoadYAMLRejectsMalformedInput(t *testing.T) {
	if _, err := LoadYAML(strings.NewReader("name: x\n\timage: y")); err == nil {
		t.Fatal("malformed YAML was accepted")
	}
}

func TestLoadYAMLFileReportsMissingFile(t *testing.T) {
	if _, err := LoadYAMLFile(filepath.Join(t.TempDir(), "absent.yaml")); err == nil {
		t.Fatal("missing file was accepted")
	}
}
