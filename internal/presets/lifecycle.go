package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

// Lifecycle commands share the StringOrSlice type. Grouped in one file to
// avoid six near-identical files.

func initializeCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base": {"echo 'Initializing on host'"},
	}
}

func InitializeCommandPreset(name string) model.StringOrSlice {
	return initializeCommandPresetsMap()[name]
}
func ListInitializeCommandPresets() []string { return sortedKeys(initializeCommandPresetsMap()) }

func onCreateCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base":  {"echo 'Container created'"},
		"setup": {"apt-get update && apt-get install -y curl"},
	}
}

func OnCreateCommandPreset(name string) model.StringOrSlice { return onCreateCommandPresetsMap()[name] }
func ListOnCreateCommandPresets() []string                  { return sortedKeys(onCreateCommandPresetsMap()) }

func updateContentCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base":        {"echo 'Content updated'"},
		"npm-install": {"npm install"},
		"go-mod-tidy": {"go mod tidy"},
	}
}

func UpdateContentCommandPreset(name string) model.StringOrSlice {
	return updateContentCommandPresetsMap()[name]
}
func ListUpdateContentCommandPresets() []string { return sortedKeys(updateContentCommandPresetsMap()) }

func postCreateCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base":     {"echo 'Container ready'"},
		"npm-deps": {"npm install"},
		"pip-deps": {"pip install -r requirements.txt"},
		"go-deps":  {"go mod download"},
	}
}

func PostCreateCommandPreset(name string) model.StringOrSlice {
	return postCreateCommandPresetsMap()[name]
}
func ListPostCreateCommandPresets() []string { return sortedKeys(postCreateCommandPresetsMap()) }

func postStartCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base": {"echo 'Container started'"},
	}
}

func PostStartCommandPreset(name string) model.StringOrSlice {
	return postStartCommandPresetsMap()[name]
}
func ListPostStartCommandPresets() []string { return sortedKeys(postStartCommandPresetsMap()) }

func postAttachCommandPresetsMap() map[string]model.StringOrSlice {
	return map[string]model.StringOrSlice{
		"base": {"echo 'Attached to container'"},
	}
}

func PostAttachCommandPreset(name string) model.StringOrSlice {
	return postAttachCommandPresetsMap()[name]
}
func ListPostAttachCommandPresets() []string { return sortedKeys(postAttachCommandPresetsMap()) }

func waitForPresetsMap() map[string]string {
	return map[string]string{
		"base":        "updateContentCommand",
		"on-create":   "onCreateCommand",
		"post-create": "postCreateCommand",
		"post-start":  "postStartCommand",
		"initialize":  "initializeCommand",
	}
}

func WaitForPreset(name string) string { return waitForPresetsMap()[name] }
func ListWaitForPresets() []string     { return sortedKeys(waitForPresetsMap()) }

func shutdownActionPresetsMap() map[string]string {
	return map[string]string{
		"base":         "stopContainer",
		"none":         "none",
		"stop-compose": "stopCompose",
	}
}

func ShutdownActionPreset(name string) string { return shutdownActionPresetsMap()[name] }
func ListShutdownActionPresets() []string     { return sortedKeys(shutdownActionPresetsMap()) }
