package model

import "testing"

// TestMetadataCoversTopLevelKeys guards against a field being added to
// TopLevelKeys (and so to the editor and the preset registry) without a
// description, which would show an empty hint panel and produce a blank cell
// in the generated reference table.
func TestMetadataCoversTopLevelKeys(t *testing.T) {
	tree := DevContainer{}.Metadata()
	for _, key := range TopLevelKeys {
		entry, ok := tree[key]
		if !ok {
			t.Errorf("no metadata for top-level key %q", key)
			continue
		}
		m, ok := entry.(meta)
		if !ok {
			t.Errorf("metadata for %q is %T, want meta", key, entry)
			continue
		}
		if m.Description == "" {
			t.Errorf("metadata for %q has an empty Description", key)
		}
	}
}
