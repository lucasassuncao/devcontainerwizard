package presets

func workspaceFolderPresetsMap() map[string]string {
	return map[string]string{
		"base":    "/workspace",
		"by-name": "/workspaces/${localWorkspaceFolderBasename}",
	}
}

func WorkspaceFolderPreset(name string) string { return workspaceFolderPresetsMap()[name] }
func ListWorkspaceFolderPresets() []string     { return sortedKeys(workspaceFolderPresetsMap()) }

func workspaceMountPresetsMap() map[string]string {
	return map[string]string{
		"base": "source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached",
	}
}

func WorkspaceMountPreset(name string) string { return workspaceMountPresetsMap()[name] }
func ListWorkspaceMountPresets() []string     { return sortedKeys(workspaceMountPresetsMap()) }
