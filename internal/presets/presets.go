package presets

import "fmt"

// stringPreset wraps a string-valued preset getter so it returns an error on
// empty (which signals "not found" for enum/string fields).
func stringPreset(field string, getter func(string) string) func(string) (string, error) {
	return func(name string) (string, error) {
		v := getter(name)
		if v == "" {
			return "", fmt.Errorf("preset %q not found for field %q", name, field)
		}
		return marshalAsBlock(field, v)
	}
}

// boolPreset wraps a bool-valued preset getter that uses a presence map.
func boolPreset(field string, mapFn func() map[string]bool, getter func(string) bool) func(string) (string, error) {
	return func(name string) (string, error) {
		if _, ok := mapFn()[name]; !ok {
			return "", fmt.Errorf("preset %q not found for field %q", name, field)
		}
		return marshalAsBlock(field, getter(name))
	}
}

// directPreset wraps a getter whose zero value is always valid (slices, maps, structs).
func directPreset(field string, getter func(string) any) func(string) (string, error) {
	return func(name string) (string, error) { return marshalAsBlock(field, getter(name)) }
}

// presetYAMLDispatch maps each field name to its YAML-producing function.
var presetYAMLDispatch = map[string]func(string) (string, error){
	"name":            stringPreset("name", NamePreset),
	"image":           stringPreset("image", ImagePreset),
	"service":         stringPreset("service", ServicePreset),
	"workspaceFolder": stringPreset("workspaceFolder", WorkspaceFolderPreset),
	"workspaceMount":  stringPreset("workspaceMount", WorkspaceMountPreset),
	"remoteUser":      stringPreset("remoteUser", RemoteUserPreset),
	"containerUser":   stringPreset("containerUser", ContainerUserPreset),
	"userEnvProbe":    stringPreset("userEnvProbe", UserEnvProbePreset),
	"startupCommand":  stringPreset("startupCommand", StartupCommandPreset),
	"command":         stringPreset("command", CommandPreset),
	"entrypoint":      stringPreset("entrypoint", EntrypointPreset),
	"waitFor":         stringPreset("waitFor", WaitForPreset),
	"shutdownAction":  stringPreset("shutdownAction", ShutdownActionPreset),

	"updateRemoteUserUID": boolPreset("updateRemoteUserUID", updateRemoteUserUIDPresetsMap, UpdateRemoteUserUIDPreset),
	"overrideCommand":     boolPreset("overrideCommand", overrideCommandPresetsMap, OverrideCommandPreset),
	"init":                boolPreset("init", initPresetsMap, InitPreset),
	"privileged":          boolPreset("privileged", privilegedPresetsMap, PrivilegedPreset),

	"dockerComposeFile":           directPreset("dockerComposeFile", func(n string) any { return DockerComposeFilePreset(n) }),
	"runServices":                 directPreset("runServices", func(n string) any { return RunServicesPreset(n) }),
	"runArgs":                     directPreset("runArgs", func(n string) any { return RunArgsPreset(n) }),
	"capAdd":                      directPreset("capAdd", func(n string) any { return CapAddPreset(n) }),
	"capDrop":                     directPreset("capDrop", func(n string) any { return CapDropPreset(n) }),
	"securityOpt":                 directPreset("securityOpt", func(n string) any { return SecurityOptPreset(n) }),
	"devices":                     directPreset("devices", func(n string) any { return DevicesPreset(n) }),
	"overrideFeatureInstallOrder": directPreset("overrideFeatureInstallOrder", func(n string) any { return OverrideFeatureInstallOrderPreset(n) }),
	"forwardPorts":                directPreset("forwardPorts", func(n string) any { return ForwardPortsPreset(n) }),
	"appPort":                     directPreset("appPort", func(n string) any { return AppPortPreset(n) }),
	"containerEnv":                directPreset("containerEnv", func(n string) any { return ContainerEnvPreset(n) }),
	"remoteEnv":                   directPreset("remoteEnv", func(n string) any { return RemoteEnvPreset(n) }),
	"localEnv":                    directPreset("localEnv", func(n string) any { return LocalEnvPreset(n) }),
	"build":                       directPreset("build", func(n string) any { return BuildPreset(n) }),
	"hostRequirements":            directPreset("hostRequirements", func(n string) any { return HostRequirementsPreset(n) }),
	"watch":                       directPreset("watch", func(n string) any { return WatchPreset(n) }),
	"mounts":                      directPreset("mounts", func(n string) any { return MountsPreset(n) }),
	"portsAttributes":             directPreset("portsAttributes", func(n string) any { return PortsAttributesPreset(n) }),
	"otherPortsAttributes":        directPreset("otherPortsAttributes", func(n string) any { return OtherPortsAttributesPreset(n) }),
	"secrets":                     directPreset("secrets", func(n string) any { return SecretsPreset(n) }),
	"features":                    directPreset("features", func(n string) any { return FeaturesPreset(n) }),
	"initializeCommand":           directPreset("initializeCommand", func(n string) any { return InitializeCommandPreset(n) }),
	"onCreateCommand":             directPreset("onCreateCommand", func(n string) any { return OnCreateCommandPreset(n) }),
	"updateContentCommand":        directPreset("updateContentCommand", func(n string) any { return UpdateContentCommandPreset(n) }),
	"postCreateCommand":           directPreset("postCreateCommand", func(n string) any { return PostCreateCommandPreset(n) }),
	"postStartCommand":            directPreset("postStartCommand", func(n string) any { return PostStartCommandPreset(n) }),
	"postAttachCommand":           directPreset("postAttachCommand", func(n string) any { return PostAttachCommandPreset(n) }),
	"customizations":              directPreset("customizations", func(n string) any { return CustomizationsPreset(n) }),
}

// PresetYAML returns the YAML block for a (field, preset) pair, ready to be
// inserted into the overlay textarea or rendered by show-examples.
// Returns an error if the field is unknown or the preset is not found.
func PresetYAML(field, name string) (string, error) {
	fn, ok := presetYAMLDispatch[field]
	if !ok {
		return "", fmt.Errorf("unknown field: %s", field)
	}
	return fn(name)
}

// listPresetsDispatch maps each field name to its List* function.
var listPresetsDispatch = map[string]func() []string{
	"name":                        ListNamePresets,
	"image":                       ListImagePresets,
	"service":                     ListServicePresets,
	"workspaceFolder":             ListWorkspaceFolderPresets,
	"workspaceMount":              ListWorkspaceMountPresets,
	"remoteUser":                  ListRemoteUserPresets,
	"containerUser":               ListContainerUserPresets,
	"userEnvProbe":                ListUserEnvProbePresets,
	"startupCommand":              ListStartupCommandPresets,
	"command":                     ListCommandPresets,
	"entrypoint":                  ListEntrypointPresets,
	"waitFor":                     ListWaitForPresets,
	"shutdownAction":              ListShutdownActionPresets,
	"updateRemoteUserUID":         ListUpdateRemoteUserUIDPresets,
	"overrideCommand":             ListOverrideCommandPresets,
	"init":                        ListInitPresets,
	"privileged":                  ListPrivilegedPresets,
	"dockerComposeFile":           ListDockerComposeFilePresets,
	"runServices":                 ListRunServicesPresets,
	"runArgs":                     ListRunArgsPresets,
	"capAdd":                      ListCapAddPresets,
	"capDrop":                     ListCapDropPresets,
	"securityOpt":                 ListSecurityOptPresets,
	"devices":                     ListDevicesPresets,
	"overrideFeatureInstallOrder": ListOverrideFeatureInstallOrderPresets,
	"forwardPorts":                ListForwardPortsPresets,
	"appPort":                     ListAppPortPresets,
	"containerEnv":                ListContainerEnvPresets,
	"remoteEnv":                   ListRemoteEnvPresets,
	"localEnv":                    ListLocalEnvPresets,
	"build":                       ListBuildPresets,
	"hostRequirements":            ListHostRequirementsPresets,
	"watch":                       ListWatchPresets,
	"mounts":                      ListMountsPresets,
	"portsAttributes":             ListPortsAttributesPresets,
	"otherPortsAttributes":        ListOtherPortsAttributesPresets,
	"secrets":                     ListSecretsPresets,
	"features":                    ListFeaturesPresets,
	"initializeCommand":           ListInitializeCommandPresets,
	"onCreateCommand":             ListOnCreateCommandPresets,
	"updateContentCommand":        ListUpdateContentCommandPresets,
	"postCreateCommand":           ListPostCreateCommandPresets,
	"postStartCommand":            ListPostStartCommandPresets,
	"postAttachCommand":           ListPostAttachCommandPresets,
	"customizations":              ListCustomizationsPresets,
}

// ListPresets returns the preset names for a field, sorted with "base" first.
// Returns nil for unknown fields.
func ListPresets(field string) []string {
	if fn, ok := listPresetsDispatch[field]; ok {
		return fn()
	}
	return nil
}

// ListFields returns the canonical ordering of preset-providing fields,
// matching the source order of model.DevContainer struct fields.
func ListFields() []string {
	return []string{
		"name", "image", "build", "dockerComposeFile", "service", "runServices",
		"workspaceFolder", "workspaceMount", "remoteUser", "containerUser",
		"updateRemoteUserUID", "userEnvProbe",
		"containerEnv", "remoteEnv", "localEnv",
		"forwardPorts", "appPort", "portsAttributes", "otherPortsAttributes",
		"mounts", "runArgs", "startupCommand", "overrideCommand",
		"command", "entrypoint",
		"init", "privileged", "capAdd", "capDrop", "securityOpt", "devices",
		"hostRequirements", "overrideFeatureInstallOrder", "features",
		"initializeCommand", "onCreateCommand", "updateContentCommand",
		"postCreateCommand", "postStartCommand", "postAttachCommand", "waitFor",
		"watch", "customizations", "secrets", "shutdownAction",
	}
}
