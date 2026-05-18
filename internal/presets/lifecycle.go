package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

// cmdPtr is a convenience helper that wraps a single-string CommandValue in a pointer.
func cmdPtr(s string) *model.CommandValue { v := model.CommandString(s); return &v }

// cmdMapPtr wraps a named parallel-command map in a pointer.
func cmdMapPtr(m map[string][]string) *model.CommandValue { v := model.CommandMap(m); return &v }

func initializeCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base": cmdPtr("echo 'Initializing on host'"),
	}
}

func InitializeCommandPreset(name string) *model.CommandValue {
	return initializeCommandPresetsMap()[name]
}
func ListInitializeCommandPresets() []string { return sortedKeys(initializeCommandPresetsMap()) }

func onCreateCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base":     cmdPtr("echo 'Container created'"),
		"setup":    cmdPtr("apt-get update && apt-get install -y curl"),
		"parallel": cmdMapPtr(map[string][]string{"apt": {"apt-get update && apt-get install -y curl"}, "pip": {"pip install -r requirements.txt"}}),
	}
}

func OnCreateCommandPreset(name string) *model.CommandValue {
	return onCreateCommandPresetsMap()[name]
}
func ListOnCreateCommandPresets() []string { return sortedKeys(onCreateCommandPresetsMap()) }

func updateContentCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base":        cmdPtr("echo 'Content updated'"),
		"npm-install": cmdPtr("npm install"),
		"go-mod-tidy": cmdPtr("go mod tidy"),
	}
}

func UpdateContentCommandPreset(name string) *model.CommandValue {
	return updateContentCommandPresetsMap()[name]
}
func ListUpdateContentCommandPresets() []string { return sortedKeys(updateContentCommandPresetsMap()) }

func postCreateCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base":     cmdPtr("echo 'Container ready'"),
		"npm-deps": cmdPtr("npm install"),
		"pip-deps": cmdPtr("pip install -r requirements.txt"),
		"go-deps":  cmdPtr("go mod download"),
		"parallel": cmdMapPtr(map[string][]string{"deps": {"go mod download"}, "tools": {"go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}}),
	}
}

func PostCreateCommandPreset(name string) *model.CommandValue {
	return postCreateCommandPresetsMap()[name]
}
func ListPostCreateCommandPresets() []string { return sortedKeys(postCreateCommandPresetsMap()) }

func postStartCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base": cmdPtr("echo 'Container started'"),
	}
}

func PostStartCommandPreset(name string) *model.CommandValue {
	return postStartCommandPresetsMap()[name]
}
func ListPostStartCommandPresets() []string { return sortedKeys(postStartCommandPresetsMap()) }

func postAttachCommandPresetsMap() map[string]*model.CommandValue {
	return map[string]*model.CommandValue{
		"base": cmdPtr("echo 'Attached to container'"),
	}
}

func PostAttachCommandPreset(name string) *model.CommandValue {
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
		"base": "stopContainer",
		"none": "none",
	}
}

func ShutdownActionPreset(name string) string { return shutdownActionPresetsMap()[name] }
func ListShutdownActionPresets() []string     { return sortedKeys(shutdownActionPresetsMap()) }
