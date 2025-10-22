package devcontainer

import (
	"fmt"

	kyaml "github.com/knadh/koanf/parsers/yaml"
	kfile "github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
)

func LoadYAMLFile(path string) (*koanf.Koanf, error) {
	k := koanf.New(".")
	if err := k.Load(kfile.Provider(path), kyaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading file: %w", err)
	}
	return k, nil
}
