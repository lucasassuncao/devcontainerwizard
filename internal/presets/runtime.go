package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func runArgsPresetsMap() map[string][]string {
	return map[string][]string{
		"base":  {"--network=host"},
		"debug": {"--cap-add=SYS_PTRACE", "--security-opt=seccomp=unconfined"},
	}
}

func RunArgsPreset(name string) []string { return runArgsPresetsMap()[name] }
func ListRunArgsPresets() []string       { return sortedKeys(runArgsPresetsMap()) }

func mountsPresetsMap() map[string][]model.Mount {
	return map[string][]model.Mount{
		"base": {
			{
				Type:        "bind",
				Source:      "${localWorkspaceFolder}/.cache",
				Target:      "/home/vscode/.cache",
				Consistency: "cached",
			},
		},
		"data-volume": {
			{
				Type:   "volume",
				Source: "app-data",
				Target: "/var/lib/app",
			},
		},
	}
}

func MountsPreset(name string) []model.Mount { return mountsPresetsMap()[name] }
func ListMountsPresets() []string            { return sortedKeys(mountsPresetsMap()) }
