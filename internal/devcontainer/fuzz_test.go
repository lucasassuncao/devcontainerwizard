package devcontainer

import (
	"encoding/json"
	"strings"
	"testing"
)

// FuzzConfigPipeline drives the whole convert path over generated input. The
// config file is the one thing a user hand-writes, so it is the untrusted
// surface of this tool, and a panic here is a crash they can trigger with a
// typo.
//
// Beyond "does not panic" it asserts the invariant the null-entry bugs broke: a
// config that Validate accepts must marshal to a devcontainer.json with no null
// where the specification allows none. Validate is what stands between a
// hand-written config and a file the IDE refuses to open, so any input that
// slips a null past it is a bug in Validate, not in the input.
func FuzzConfigPipeline(f *testing.F) {
	for _, body := range probes {
		f.Add(body)
	}

	f.Fuzz(func(t *testing.T, body string) {
		k, err := LoadYAML(strings.NewReader(body))
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

		encoded, err := json.Marshal(dc)
		if err != nil {
			t.Fatalf("validated config failed to marshal: %v\ninput: %q", err, body)
		}

		var out map[string]any
		if err := json.Unmarshal(encoded, &out); err != nil {
			t.Fatalf("emitted invalid JSON: %v\ninput: %q", err, body)
		}
		if mounts, ok := out["mounts"].([]any); ok {
			for i, m := range mounts {
				if m == nil {
					t.Fatalf("mounts[%d] emitted as null\ninput: %q", i, body)
				}
			}
		}
		if attrs, ok := out["portsAttributes"].(map[string]any); ok {
			for port, v := range attrs {
				if v == nil {
					t.Fatalf("portsAttributes[%q] emitted as null\ninput: %q", port, body)
				}
			}
		}
	})
}
