package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func watchPresetsMap() map[string]*model.WatchConfig {
	return map[string]*model.WatchConfig{
		"base": {
			WaitFor: []string{"postCreateCommand"},
			Restart: []string{".devcontainer/devcontainer.json"},
		},
	}
}

func WatchPreset(name string) *model.WatchConfig { return watchPresetsMap()[name] }
func ListWatchPresets() []string                 { return sortedKeys(watchPresetsMap()) }
