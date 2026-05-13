package presets

func commandPresetsMap() map[string]string {
	return map[string]string{
		"base": "sleep infinity",
	}
}

func CommandPreset(name string) string { return commandPresetsMap()[name] }
func ListCommandPresets() []string     { return sortedKeys(commandPresetsMap()) }

func entrypointPresetsMap() map[string]string {
	return map[string]string{
		"base": "/usr/local/bin/docker-entrypoint.sh",
	}
}

func EntrypointPreset(name string) string { return entrypointPresetsMap()[name] }
func ListEntrypointPresets() []string     { return sortedKeys(entrypointPresetsMap()) }

func startupCommandPresetsMap() map[string]string {
	return map[string]string{
		"base": "echo 'Container started'",
	}
}

func StartupCommandPreset(name string) string { return startupCommandPresetsMap()[name] }
func ListStartupCommandPresets() []string     { return sortedKeys(startupCommandPresetsMap()) }

func overrideCommandPresetsMap() map[string]bool {
	return map[string]bool{
		"base": true,
	}
}

func OverrideCommandPreset(name string) bool { return overrideCommandPresetsMap()[name] }
func ListOverrideCommandPresets() []string   { return sortedKeys(overrideCommandPresetsMap()) }
