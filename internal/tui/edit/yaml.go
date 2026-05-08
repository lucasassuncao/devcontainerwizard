package edit

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Block represents a top-level YAML key with its line range (1-based).
type Block struct {
	Key     string
	Line    int // line of the key node
	EndLine int // last line occupied by this block (exclusive of next key)
}

// ParseBlocks reads path and returns its top-level blocks.
func ParseBlocks(path string) ([]Block, error) {
	raw, err := os.ReadFile(path) // #nosec G304 -- path is user-supplied via CLI arg
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return ParseBlocksFromBytes(raw)
}

// ParseBlocksFromBytes parses raw YAML bytes and returns top-level blocks.
func ParseBlocksFromBytes(raw []byte) ([]Block, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}
	if doc.Kind == 0 || len(doc.Content) == 0 {
		return nil, nil
	}
	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping at root, got kind %d", mapping.Kind)
	}

	totalLines := bytes.Count(raw, []byte("\n")) + 1
	blocks := make([]Block, 0, len(mapping.Content)/2)

	for i := 0; i < len(mapping.Content)-1; i += 2 {
		keyNode := mapping.Content[i]
		blocks = append(blocks, Block{
			Key:  keyNode.Value,
			Line: keyNode.Line,
		})
	}

	// Fill EndLine: each block ends one line before the next key starts.
	for i := range blocks {
		if i+1 < len(blocks) {
			blocks[i].EndLine = blocks[i+1].Line - 1
		} else {
			blocks[i].EndLine = totalLines
		}
	}

	return blocks, nil
}

// RemoveBlock deletes the lines belonging to key from raw YAML bytes.
func RemoveBlock(raw []byte, blocks []Block, key string) ([]byte, error) {
	var target *Block
	for i := range blocks {
		if blocks[i].Key == key {
			target = &blocks[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("key %q not found in blocks", key)
	}

	lines := splitLines(raw)
	// Lines are 1-based; slice indices are 0-based.
	start := target.Line - 1
	end := target.EndLine // exclusive upper bound (0-based = EndLine)
	lines = append(lines[:start:start], lines[end:]...)
	return joinLines(lines), nil
}

// InsertBlock inserts a YAML snippet into raw, respecting the canonical key
// order defined by allKnownKeys. The snippet is placed before the first
// existing block whose key follows the new key in that order. If the new key
// is unknown or no later block exists, the snippet is appended at the end.
func InsertBlock(raw []byte, snippet string) ([]byte, error) {
	if err := ValidateSnippet(snippet); err != nil {
		return nil, err
	}

	// Determine which key the snippet introduces.
	snippetBlocks, err := ParseBlocksFromBytes([]byte(snippet))
	if err != nil || len(snippetBlocks) == 0 {
		return appendBlock(raw, snippet), nil
	}
	newKey := snippetBlocks[0].Key

	// Build a rank map from allKnownKeys.
	rank := make(map[string]int, len(allKnownKeys))
	for i, k := range allKnownKeys {
		rank[k] = i
	}
	newRank, known := rank[newKey]
	if !known {
		return appendBlock(raw, snippet), nil
	}

	// Parse existing blocks to find the insertion line.
	blocks, err := ParseBlocksFromBytes(raw)
	if err != nil || len(blocks) == 0 {
		return appendBlock(raw, snippet), nil
	}

	// Find the first existing block that should come after newKey.
	insertBeforeLine := -1
	for _, b := range blocks {
		if r, ok := rank[b.Key]; ok && r > newRank {
			insertBeforeLine = b.Line // 1-based
			break
		}
	}

	if insertBeforeLine == -1 {
		return appendBlock(raw, snippet), nil
	}

	// Insert snippet lines before insertBeforeLine (convert to 0-based index).
	lines := splitLines(raw)
	idx := insertBeforeLine - 1
	snippetLines := splitLines([]byte(snippet))
	// Remove trailing empty string that splitLines adds for a newline-terminated snippet.
	if len(snippetLines) > 0 && snippetLines[len(snippetLines)-1] == "" {
		snippetLines = snippetLines[:len(snippetLines)-1]
	}
	merged := make([]string, 0, len(lines)+len(snippetLines))
	merged = append(merged, lines[:idx]...)
	merged = append(merged, snippetLines...)
	merged = append(merged, lines[idx:]...)
	return joinLines(merged), nil
}

// appendBlock adds snippet after the last non-empty line of raw.
func appendBlock(raw []byte, snippet string) []byte {
	trimmed := bytes.TrimRight(raw, "\n")
	if len(trimmed) == 0 {
		return []byte(snippet)
	}
	return append(trimmed, append([]byte("\n"), []byte(snippet)...)...)
}

// BlockContent returns the raw lines for a given block key.
func BlockContent(raw []byte, blocks []Block, key string) (string, error) {
	for _, b := range blocks {
		if b.Key == key {
			lines := splitLines(raw)
			start := b.Line - 1
			end := b.EndLine
			if end > len(lines) {
				end = len(lines)
			}
			return string(joinLines(lines[start:end])), nil
		}
	}
	return "", fmt.Errorf("key %q not found", key)
}

// ValidateSnippet returns an error if the YAML text is not parseable.
func ValidateSnippet(text string) error {
	var check interface{}
	return yaml.Unmarshal([]byte(text), &check)
}

func splitLines(raw []byte) []string {
	return strings.Split(string(raw), "\n")
}

func joinLines(lines []string) []byte {
	return []byte(strings.Join(lines, "\n"))
}

// rebuildYAML constructs the YAML content for key from the checked field states.
func rebuildYAML(key string, fields []fieldState) string {
	var sb strings.Builder
	sb.WriteString(key + ":\n")
	for _, fs := range fields {
		if fs.Checked {
			sb.WriteString(fs.Def.YAML)
		}
	}
	return sb.String()
}

// syncFieldsFromYAML updates Checked on each field to reflect what is present
// in content (the current textarea value for the given key).
func syncFieldsFromYAML(key string, fields []fieldState, content string) []fieldState {
	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return fields
	}
	sub, _ := doc[key].(map[string]interface{})
	out := make([]fieldState, len(fields))
	copy(out, fields)
	for i := range out {
		_, out[i].Checked = sub[out[i].Def.Key]
	}
	return out
}

// applyFieldToggle surgically adds or removes a single sub-field from current
// (the textarea YAML value), preserving edits to other fields.
// checked is the NEW desired state of def.
func applyFieldToggle(key string, fields []fieldState, def FieldDef, checked bool, current string) string {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(current), &root); err != nil || root.Kind == 0 || len(root.Content) == 0 {
		return rebuildYAML(key, fields)
	}
	mapping := root.Content[0]
	if mapping.Kind != yaml.MappingNode || len(mapping.Content) < 2 {
		return rebuildYAML(key, fields)
	}
	valueNode := mapping.Content[1]
	if valueNode.Kind != yaml.MappingNode {
		return rebuildYAML(key, fields)
	}

	idx := -1
	for i := 0; i < len(valueNode.Content)-1; i += 2 {
		if valueNode.Content[i].Value == def.Key {
			idx = i
			break
		}
	}

	if !checked {
		removeFieldNode(valueNode, idx)
	} else {
		addFieldNode(valueNode, idx, key, def)
	}

	var buf strings.Builder
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&root); err != nil {
		return rebuildYAML(key, fields)
	}
	return strings.TrimRight(buf.String(), "\n") + "\n"
}

func removeFieldNode(valueNode *yaml.Node, idx int) {
	if idx >= 0 {
		valueNode.Content = append(valueNode.Content[:idx], valueNode.Content[idx+2:]...)
	}
}

func addFieldNode(valueNode *yaml.Node, idx int, parentKey string, def FieldDef) {
	if idx >= 0 {
		return // already present
	}
	var templateRoot yaml.Node
	if err := yaml.Unmarshal([]byte(parentKey+":\n"+def.YAML), &templateRoot); err != nil {
		return
	}
	if templateRoot.Kind == 0 || len(templateRoot.Content) == 0 {
		return
	}
	tMapping := templateRoot.Content[0]
	if tMapping.Kind != yaml.MappingNode || len(tMapping.Content) < 2 {
		return
	}
	tValue := tMapping.Content[1]
	if tValue.Kind == yaml.MappingNode && len(tValue.Content) >= 2 {
		valueNode.Content = append(valueNode.Content, tValue.Content[0], tValue.Content[1])
	}
}
