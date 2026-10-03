package devcontainer

import (
	"fmt"
	"io"
	"os"

	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	koanf "github.com/knadh/koanf/v2"
)

// LoadYAML parses a devcontainerwizard configuration from r.
//
// Taking a reader rather than only a path keeps parsing independent of the
// filesystem: the editor can validate a buffer it has not written yet, and
// tests can exercise the pipeline without a temp directory.
func LoadYAML(r io.Reader) (*koanf.Koanf, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	k := koanf.New(".")
	if err := k.Load(rawbytes.Provider(data), kyaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading file: %w", err)
	}
	return k, nil
}

// LoadYAMLFile parses the configuration stored at path. It is a convenience
// wrapper around LoadYAML for the common case of reading from disk.
func LoadYAMLFile(path string) (*koanf.Koanf, error) {
	f, err := os.Open(path) //#nosec G304 -- the path comes from the user's own --config flag
	if err != nil {
		return nil, fmt.Errorf("error loading file: %w", err)
	}
	defer f.Close()
	return LoadYAML(f)
}
