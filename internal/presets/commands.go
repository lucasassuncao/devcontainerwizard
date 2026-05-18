package presets

func overrideCommandPresetsMap() map[string]bool {
	return map[string]bool{
		"base": true,
	}
}

func OverrideCommandPreset(name string) bool { return overrideCommandPresetsMap()[name] }
func ListOverrideCommandPresets() []string   { return sortedKeys(overrideCommandPresetsMap()) }
