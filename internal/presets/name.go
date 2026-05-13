package presets

func namePresetsMap() map[string]string {
	return map[string]string{
		"base": "my-devcontainer",
	}
}

func NamePreset(name string) string { return namePresetsMap()[name] }
func ListNamePresets() []string     { return sortedKeys(namePresetsMap()) }
