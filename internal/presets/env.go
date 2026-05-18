package presets

func containerEnvPresetsMap() map[string]map[string]string {
	return map[string]map[string]string{
		"base": {
			"NODE_ENV": "development",
		},
		"verbose-logging": {
			"LOG_LEVEL": "debug",
			"NODE_ENV":  "development",
		},
	}
}

func ContainerEnvPreset(name string) map[string]string { return containerEnvPresetsMap()[name] }
func ListContainerEnvPresets() []string                { return sortedKeys(containerEnvPresetsMap()) }

func remoteEnvPresetsMap() map[string]map[string]string {
	return map[string]map[string]string{
		"base": {
			"PATH": "${containerEnv:PATH}:/usr/local/bin",
		},
	}
}

func RemoteEnvPreset(name string) map[string]string { return remoteEnvPresetsMap()[name] }
func ListRemoteEnvPresets() []string                { return sortedKeys(remoteEnvPresetsMap()) }
