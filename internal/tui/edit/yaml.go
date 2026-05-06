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
