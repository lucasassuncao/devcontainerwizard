package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func forwardPortsPresetsMap() map[string][]any {
	return map[string][]any{
		"base":      {3000},
		"web-stack": {3000, 5432, 6379},
		"node-dev":  {3000, 9229},
	}
}

func ForwardPortsPreset(name string) []any { return forwardPortsPresetsMap()[name] }
func ListForwardPortsPresets() []string    { return sortedKeys(forwardPortsPresetsMap()) }

func appPortPresetsMap() map[string][]any {
	return map[string][]any{
		"base": {3000},
	}
}

func AppPortPreset(name string) []any { return appPortPresetsMap()[name] }
func ListAppPortPresets() []string    { return sortedKeys(appPortPresetsMap()) }

func portsAttributesPresetsMap() map[string]map[string]*model.PortAttributes {
	return map[string]map[string]*model.PortAttributes{
		"base": {
			"3000": {
				Label:         "Web App",
				OnAutoForward: "notify",
				Protocol:      "http",
			},
		},
		"web-stack": {
			"3000": {Label: "App", OnAutoForward: "openBrowser", Protocol: "http"},
			"5432": {Label: "PostgreSQL", OnAutoForward: "silent"},
			"6379": {Label: "Redis", OnAutoForward: "silent"},
		},
	}
}

func PortsAttributesPreset(name string) map[string]*model.PortAttributes {
	return portsAttributesPresetsMap()[name]
}
func ListPortsAttributesPresets() []string { return sortedKeys(portsAttributesPresetsMap()) }

func otherPortsAttributesPresetsMap() map[string]*model.PortAttributes {
	return map[string]*model.PortAttributes{
		"base": {
			OnAutoForward: "silent",
		},
		"ignore": {
			OnAutoForward: "ignore",
		},
	}
}

func OtherPortsAttributesPreset(name string) *model.PortAttributes {
	return otherPortsAttributesPresetsMap()[name]
}
func ListOtherPortsAttributesPresets() []string { return sortedKeys(otherPortsAttributesPresetsMap()) }
