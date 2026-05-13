package presets

func remoteUserPresetsMap() map[string]string {
	return map[string]string{
		"base": "vscode",
		"root": "root",
		"node": "node",
	}
}

func RemoteUserPreset(name string) string { return remoteUserPresetsMap()[name] }
func ListRemoteUserPresets() []string     { return sortedKeys(remoteUserPresetsMap()) }

func containerUserPresetsMap() map[string]string {
	return map[string]string{
		"base": "vscode",
		"root": "root",
	}
}

func ContainerUserPreset(name string) string { return containerUserPresetsMap()[name] }
func ListContainerUserPresets() []string     { return sortedKeys(containerUserPresetsMap()) }

func updateRemoteUserUIDPresetsMap() map[string]bool {
	return map[string]bool{
		"base": true,
	}
}

func UpdateRemoteUserUIDPreset(name string) bool { return updateRemoteUserUIDPresetsMap()[name] }
func ListUpdateRemoteUserUIDPresets() []string   { return sortedKeys(updateRemoteUserUIDPresetsMap()) }

func userEnvProbePresetsMap() map[string]string {
	return map[string]string{
		"base":              "loginInteractiveShell",
		"none":              "none",
		"login-shell":       "loginShell",
		"interactive-shell": "interactiveShell",
	}
}

func UserEnvProbePreset(name string) string { return userEnvProbePresetsMap()[name] }
func ListUserEnvProbePresets() []string     { return sortedKeys(userEnvProbePresetsMap()) }
