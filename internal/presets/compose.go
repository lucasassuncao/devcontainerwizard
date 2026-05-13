package presets

func dockerComposeFilePresetsMap() map[string][]string {
	return map[string][]string{
		"base":     {"docker-compose.yml"},
		"with-dev": {"docker-compose.yml", "docker-compose.dev.yml"},
	}
}

func DockerComposeFilePreset(name string) []string { return dockerComposeFilePresetsMap()[name] }
func ListDockerComposeFilePresets() []string       { return sortedKeys(dockerComposeFilePresetsMap()) }

func servicePresetsMap() map[string]string {
	return map[string]string{
		"base": "app",
	}
}

func ServicePreset(name string) string { return servicePresetsMap()[name] }
func ListServicePresets() []string     { return sortedKeys(servicePresetsMap()) }

func runServicesPresetsMap() map[string][]string {
	return map[string][]string{
		"base":      {"db"},
		"web-stack": {"db", "redis"},
	}
}

func RunServicesPreset(name string) []string { return runServicesPresetsMap()[name] }
func ListRunServicesPresets() []string       { return sortedKeys(runServicesPresetsMap()) }
