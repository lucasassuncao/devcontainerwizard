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

func mountsPresetsMap() map[string][]model.MountOrString {
	return map[string][]model.MountOrString{
		"base": {
			model.MountObject(model.Mount{
				Type:   "bind",
				Source: "${localWorkspaceFolder}/.cache",
				Target: "/home/vscode/.cache",
			}),
		},
		"data-volume": {
			model.MountObject(model.Mount{
				Type:   "volume",
				Source: "app-data",
				Target: "/var/lib/app",
			}),
		},
		"string-form": {
			model.MountString("source=${localWorkspaceFolder}/.cache,target=/home/vscode/.cache,type=bind,consistency=cached"),
		},
	}
}

func MountsPreset(name string) []model.MountOrString { return mountsPresetsMap()[name] }
func ListMountsPresets() []string                    { return sortedKeys(mountsPresetsMap()) }
