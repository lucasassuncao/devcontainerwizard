package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func hostRequirementsPresetsMap() map[string]*model.HostRequirements {
	return map[string]*model.HostRequirements{
		"base": {
			CPUs:   2,
			Memory: "4gb",
		},
		"heavy": {
			CPUs:    8,
			Memory:  "16gb",
			Storage: "64gb",
		},
		"gpu": {
			CPUs:   4,
			Memory: "16gb",
			GPU: &model.GPURequirement{
				Cores:  1,
				Memory: "4gb",
			},
		},
	}
}

func HostRequirementsPreset(name string) *model.HostRequirements {
	return hostRequirementsPresetsMap()[name]
}
func ListHostRequirementsPresets() []string { return sortedKeys(hostRequirementsPresetsMap()) }
