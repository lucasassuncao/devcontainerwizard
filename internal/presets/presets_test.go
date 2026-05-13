package presets

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

func TestSortedKeysBaseFirst(t *testing.T) {
	in := map[string]string{"z": "z", "a": "a", "base": "b"}
	got := sortedKeys(in)
	want := []string{"base", "a", "z"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSortedKeysWithoutBase(t *testing.T) {
	in := map[string]string{"foo": "", "bar": ""}
	got := sortedKeys(in)
	want := []string{"bar", "foo"}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestMarshalAsBlockScalar(t *testing.T) {
	got, err := marshalAsBlock("name", "my-devcontainer")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	want := "name: my-devcontainer\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMarshalAsBlockBool(t *testing.T) {
	got, err := marshalAsBlock("privileged", false)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	want := "privileged: false\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMarshalAsBlockSlice(t *testing.T) {
	got, err := marshalAsBlock("capAdd", []string{"SYS_PTRACE", "NET_ADMIN"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	want := "capAdd:\n  - SYS_PTRACE\n  - NET_ADMIN\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMarshalAsBlockNilPointer(t *testing.T) {
	type Foo struct{ X int }
	var p *Foo
	_, err := marshalAsBlock("foo", p)
	if err == nil {
		t.Fatal("want error for nil pointer, got nil")
	}
}

func TestMarshalAsBlockUntypedNil(t *testing.T) {
	_, err := marshalAsBlock("foo", nil)
	if err == nil {
		t.Fatal("want error for nil, got nil")
	}
}

func TestPresetYAMLUnknownField(t *testing.T) {
	_, err := PresetYAML("not-a-field", "base")
	if err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestListPresetsUnknownField(t *testing.T) {
	got := ListPresets("not-a-field")
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestListFieldsNotEmpty(t *testing.T) {
	got := ListFields()
	if len(got) == 0 {
		t.Fatal("ListFields returned empty slice")
	}
}

func TestCustomizationsPresetBase(t *testing.T) {
	p := CustomizationsPreset("base")
	if p == nil {
		t.Fatal("base preset should exist")
	}
}

func TestCustomizationsPresetMissing(t *testing.T) {
	if got := CustomizationsPreset("does-not-exist"); got != nil {
		t.Errorf("missing preset should return nil, got %+v", got)
	}
}

func TestListCustomizationsPresetsHasBase(t *testing.T) {
	got := ListCustomizationsPresets()
	if len(got) == 0 || got[0] != "base" {
		t.Errorf("ListCustomizationsPresets()[0] = %q, want \"base\"", got)
	}
}

func TestPresetYAMLForCustomizationsBase(t *testing.T) {
	y, err := PresetYAML("customizations", "base")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.HasPrefix(y, "customizations:\n") {
		t.Errorf("yaml should start with \"customizations:\\n\", got %q", y)
	}
}

func TestPresetYAMLRoundtripAllFields(t *testing.T) {
	for _, field := range ListFields() {
		presetNames := ListPresets(field)
		if len(presetNames) == 0 {
			t.Errorf("field %q has no presets", field)
			continue
		}
		for _, name := range presetNames {
			y, err := PresetYAML(field, name)
			if err != nil {
				t.Errorf("PresetYAML(%q, %q) error: %v", field, name, err)
				continue
			}
			var dc model.DevContainer
			if err := yaml.Unmarshal([]byte(y), &dc); err != nil {
				t.Errorf("yaml.Unmarshal(%q.%q): %v\nyaml:\n%s", field, name, err, y)
			}
		}
	}
}

// TestPresetYAMLValidatesSchema sanity-checks that a representative preset per
// type-category produces YAML that validates as a complete DevContainer config.
func TestPresetYAMLValidatesSchema(t *testing.T) {
	cases := []struct {
		field, name string
	}{
		{"customizations", "base"},
		{"image", "golang"},
		{"capAdd", "base"},
		{"build", "base"},
		{"mounts", "base"},
		{"features", "base"},
		{"onCreateCommand", "base"},
	}
	shell := "name: smoke\nimage: ubuntu:22.04\n"
	for _, c := range cases {
		body, err := PresetYAML(c.field, c.name)
		if err != nil {
			t.Errorf("PresetYAML(%s, %s): %v", c.field, c.name, err)
			continue
		}
		full := shell + body
		if c.field == "name" || c.field == "image" {
			full = body
		}
		var dc model.DevContainer
		if err := yaml.Unmarshal([]byte(full), &dc); err != nil {
			t.Errorf("yaml.Unmarshal(%s, %s): %v\nyaml:\n%s", c.field, c.name, err, full)
		}
	}
}

func TestAllFieldsHaveBasePreset(t *testing.T) {
	for _, field := range ListFields() {
		names := ListPresets(field)
		if len(names) == 0 {
			t.Errorf("field %q has no presets — every field must have \"base\"", field)
			continue
		}
		found := false
		for _, n := range names {
			if n == "base" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("field %q is missing the mandatory \"base\" preset (got %v)", field, names)
		}
	}
}
