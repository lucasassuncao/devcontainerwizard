package presets

import (
	"fmt"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

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

// presetEntry pairs the YAML-producing function with the preset-name lister for
// a single field. Both live in one place so adding a field means one map entry.
type presetEntry struct {
	yaml func(string) (string, error)
	list func() []string
}

// presetRegistry is the single map consulted by PresetYAML and ListPresets.
// Every field in model.TopLevelKeys must have an entry here — the test
// TestPresetRegistryCoverageAllTopLevelKeys enforces this at test time.
var presetRegistry = map[string]presetEntry{
	"name":            {stringPreset("name", NamePreset), ListNamePresets},
	"image":           {stringPreset("image", ImagePreset), ListImagePresets},
	"service":         {stringPreset("service", ServicePreset), ListServicePresets},
	"workspaceFolder": {stringPreset("workspaceFolder", WorkspaceFolderPreset), ListWorkspaceFolderPresets},
	"workspaceMount":  {stringPreset("workspaceMount", WorkspaceMountPreset), ListWorkspaceMountPresets},
	"remoteUser":      {stringPreset("remoteUser", RemoteUserPreset), ListRemoteUserPresets},
	"containerUser":   {stringPreset("containerUser", ContainerUserPreset), ListContainerUserPresets},
	"userEnvProbe":    {stringPreset("userEnvProbe", UserEnvProbePreset), ListUserEnvProbePresets},
	"startupCommand":  {stringPreset("startupCommand", StartupCommandPreset), ListStartupCommandPresets},
	"command":         {stringPreset("command", CommandPreset), ListCommandPresets},
	"entrypoint":      {stringPreset("entrypoint", EntrypointPreset), ListEntrypointPresets},
	"waitFor":         {stringPreset("waitFor", WaitForPreset), ListWaitForPresets},
	"shutdownAction":  {stringPreset("shutdownAction", ShutdownActionPreset), ListShutdownActionPresets},

	"updateRemoteUserUID": {boolPreset("updateRemoteUserUID", updateRemoteUserUIDPresetsMap, UpdateRemoteUserUIDPreset), ListUpdateRemoteUserUIDPresets},
	"overrideCommand":     {boolPreset("overrideCommand", overrideCommandPresetsMap, OverrideCommandPreset), ListOverrideCommandPresets},
	"init":                {boolPreset("init", initPresetsMap, InitPreset), ListInitPresets},
	"privileged":          {boolPreset("privileged", privilegedPresetsMap, PrivilegedPreset), ListPrivilegedPresets},

	"dockerComposeFile":           {directPreset("dockerComposeFile", func(n string) any { return DockerComposeFilePreset(n) }), ListDockerComposeFilePresets},
	"runServices":                 {directPreset("runServices", func(n string) any { return RunServicesPreset(n) }), ListRunServicesPresets},
	"runArgs":                     {directPreset("runArgs", func(n string) any { return RunArgsPreset(n) }), ListRunArgsPresets},
	"capAdd":                      {directPreset("capAdd", func(n string) any { return CapAddPreset(n) }), ListCapAddPresets},
	"capDrop":                     {directPreset("capDrop", func(n string) any { return CapDropPreset(n) }), ListCapDropPresets},
	"securityOpt":                 {directPreset("securityOpt", func(n string) any { return SecurityOptPreset(n) }), ListSecurityOptPresets},
	"devices":                     {directPreset("devices", func(n string) any { return DevicesPreset(n) }), ListDevicesPresets},
	"overrideFeatureInstallOrder": {directPreset("overrideFeatureInstallOrder", func(n string) any { return OverrideFeatureInstallOrderPreset(n) }), ListOverrideFeatureInstallOrderPresets},
	"forwardPorts":                {directPreset("forwardPorts", func(n string) any { return ForwardPortsPreset(n) }), ListForwardPortsPresets},
	"appPort":                     {directPreset("appPort", func(n string) any { return AppPortPreset(n) }), ListAppPortPresets},
	"containerEnv":                {directPreset("containerEnv", func(n string) any { return ContainerEnvPreset(n) }), ListContainerEnvPresets},
	"remoteEnv":                   {directPreset("remoteEnv", func(n string) any { return RemoteEnvPreset(n) }), ListRemoteEnvPresets},
	"localEnv":                    {directPreset("localEnv", func(n string) any { return LocalEnvPreset(n) }), ListLocalEnvPresets},
	"build":                       {directPreset("build", func(n string) any { return BuildPreset(n) }), ListBuildPresets},
	"hostRequirements":            {directPreset("hostRequirements", func(n string) any { return HostRequirementsPreset(n) }), ListHostRequirementsPresets},
	"watch":                       {directPreset("watch", func(n string) any { return WatchPreset(n) }), ListWatchPresets},
	"mounts":                      {directPreset("mounts", func(n string) any { return MountsPreset(n) }), ListMountsPresets},
	"portsAttributes":             {directPreset("portsAttributes", func(n string) any { return PortsAttributesPreset(n) }), ListPortsAttributesPresets},
	"otherPortsAttributes":        {directPreset("otherPortsAttributes", func(n string) any { return OtherPortsAttributesPreset(n) }), ListOtherPortsAttributesPresets},
	"secrets":                     {directPreset("secrets", func(n string) any { return SecretsPreset(n) }), ListSecretsPresets},
	"features":                    {directPreset("features", func(n string) any { return FeaturesPreset(n) }), ListFeaturesPresets},
	"initializeCommand":           {directPreset("initializeCommand", func(n string) any { return InitializeCommandPreset(n) }), ListInitializeCommandPresets},
	"onCreateCommand":             {directPreset("onCreateCommand", func(n string) any { return OnCreateCommandPreset(n) }), ListOnCreateCommandPresets},
	"updateContentCommand":        {directPreset("updateContentCommand", func(n string) any { return UpdateContentCommandPreset(n) }), ListUpdateContentCommandPresets},
	"postCreateCommand":           {directPreset("postCreateCommand", func(n string) any { return PostCreateCommandPreset(n) }), ListPostCreateCommandPresets},
	"postStartCommand":            {directPreset("postStartCommand", func(n string) any { return PostStartCommandPreset(n) }), ListPostStartCommandPresets},
	"postAttachCommand":           {directPreset("postAttachCommand", func(n string) any { return PostAttachCommandPreset(n) }), ListPostAttachCommandPresets},
	"customizations":              {directPreset("customizations", func(n string) any { return CustomizationsPreset(n) }), ListCustomizationsPresets},
}

// PresetYAML returns the YAML block for a (field, preset) pair, ready to be
// inserted into the overlay textarea or rendered by show-examples.
// Returns an error if the field is unknown or the preset is not found.
func PresetYAML(field, name string) (string, error) {
	e, ok := presetRegistry[field]
	if !ok {
		return "", fmt.Errorf("unknown field: %s", field)
	}
	return e.yaml(name)
}

// ListPresets returns the preset names for a field, sorted with "base" first.
// Returns nil for unknown fields.
func ListPresets(field string) []string {
	if e, ok := presetRegistry[field]; ok {
		return e.list()
	}
	return nil
}

// ListFields returns the canonical ordering of all top-level DevContainer fields.
// Delegates to model.TopLevelKeys — the single source of truth.
func ListFields() []string {
	return model.TopLevelKeys
}
