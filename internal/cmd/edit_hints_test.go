package cmd

import "testing"

// TestBuildDevContainerHints_metadataReachesFieldMeta guards the contract the
// FromMetadata validators depend on: the metadata tree must be accepted by
// metadata.New, and key fields must carry the correct metadata.
func TestBuildDevContainerHints_metadataReachesFieldMeta(t *testing.T) {
	src, err := buildDevContainerHints()
	if err != nil {
		t.Fatalf("buildDevContainerHints: %v", err)
	}

	required := []struct{ block, path string }{
		{"name", ""},
		{"build", "dockerfile"},
		{"mounts", "type"},
		{"mounts", "target"},
	}
	for _, f := range required {
		if !src.FieldMeta(f.block, f.path).Required {
			t.Errorf("FieldMeta(%q, %q).Required = false, want true", f.block, f.path)
		}
	}

	notRequired := []struct{ block, path string }{
		{"image", ""},
		{"build", "args"},
		{"build", "target"},
		{"secrets", "documentationUrl"},
		{"secrets", "description"},
		{"build", "context"},
		{"watch", "waitFor"},
	}
	for _, f := range notRequired {
		if src.FieldMeta(f.block, f.path).Required {
			t.Errorf("FieldMeta(%q, %q).Required = true, want false", f.block, f.path)
		}
	}

	oneOf := []struct {
		block, path string
		count       int
	}{
		{"userEnvProbe", "", 4},
		{"waitFor", "", 5},
		{"shutdownAction", "", 3},
		{"portsAttributes", "onAutoForward", 6},
		{"portsAttributes", "protocol", 2},
		{"otherPortsAttributes", "onAutoForward", 5}, // no openBrowserOnce per the spec
		{"otherPortsAttributes", "protocol", 2},
	}
	for _, f := range oneOf {
		if got := len(src.FieldMeta(f.block, f.path).OneOf); got != f.count {
			t.Errorf("FieldMeta(%q, %q).OneOf has %d values, want %d", f.block, f.path, got, f.count)
		}
	}

	described := []struct{ block, path string }{
		{"name", ""},
		{"image", ""},
		{"build", "dockerfile"},
		{"dockerComposeFile", ""},
		{"mounts", "type"},
		{"hostRequirements", "cpus"},
		{"customizations", "vscode"},
		{"customizations", "vscode.extensions"},
		{"customizations", "jetbrains.plugins"},
		{"initializeCommand", ""},
		{"postCreateCommand", ""},
	}
	for _, f := range described {
		if src.FieldMeta(f.block, f.path).Description == "" {
			t.Errorf("FieldMeta(%q, %q).Description is empty", f.block, f.path)
		}
	}

	if src.FieldMeta("hostRequirements", "cpus").Min == "" {
		t.Error("hostRequirements.cpus should declare Min")
	}
}
