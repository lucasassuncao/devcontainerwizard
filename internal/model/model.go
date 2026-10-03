// Package model defines the data structures used to represent
// dev container configuration and related types. These structs are
// used for parsing, validation and documentation generation.
package model

// DevContainer is the whole configuration document: the YAML the user edits and
// the devcontainer.json generated from it. Fields are declared in the order they
// should appear in the generated JSON, which is the order encoding/json emits
// them.
type DevContainer struct {
	Schema            string       `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	Name              string       `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
	Image             string       `json:"image,omitempty" yaml:"image,omitempty"`
	Build             *BuildConfig `json:"build,omitempty" yaml:"build,omitempty" validate:"omitempty"`
	DockerFile        string       `json:"dockerFile,omitempty" yaml:"dockerFile,omitempty"`
	DockerComposeFile []string     `json:"dockerComposeFile,omitempty" yaml:"dockerComposeFile,omitempty"`
	Service           string       `json:"service,omitempty" yaml:"service,omitempty" validate:"required_with=DockerComposeFile"`
	RunServices       []string     `json:"runServices,omitempty" yaml:"runServices,omitempty"`

	WorkspaceFolder     string `json:"workspaceFolder,omitempty" yaml:"workspaceFolder,omitempty"`
	WorkspaceMount      string `json:"workspaceMount,omitempty" yaml:"workspaceMount,omitempty"`
	RemoteUser          string `json:"remoteUser,omitempty" yaml:"remoteUser,omitempty"`
	ContainerUser       string `json:"containerUser,omitempty" yaml:"containerUser,omitempty"`
	UpdateRemoteUserUID bool   `json:"updateRemoteUserUID,omitempty" yaml:"updateRemoteUserUID,omitempty"`
	UserEnvProbe        string `json:"userEnvProbe,omitempty" yaml:"userEnvProbe,omitempty" validate:"omitempty,oneof=none loginShell loginInteractiveShell interactiveShell"`

	// Environment variables
	ContainerEnv map[string]string `json:"containerEnv,omitempty" yaml:"containerEnv,omitempty"`
	RemoteEnv    map[string]string `json:"remoteEnv,omitempty" yaml:"remoteEnv,omitempty"`

	ForwardPorts         []any                      `json:"forwardPorts,omitempty" yaml:"forwardPorts,omitempty"`
	AppPort              []any                      `json:"appPort,omitempty" yaml:"appPort,omitempty"`
	PortsAttributes      map[string]*PortAttributes `json:"portsAttributes,omitempty" yaml:"portsAttributes,omitempty" validate:"omitempty,dive"`
	OtherPortsAttributes *PortAttributes            `json:"otherPortsAttributes,omitempty" yaml:"otherPortsAttributes,omitempty" validate:"omitempty"`
	Mounts               []MountOrString            `json:"mounts,omitempty" yaml:"mounts,omitempty" validate:"omitempty,dive"`

	RunArgs         []string `json:"runArgs,omitempty" yaml:"runArgs,omitempty"`
	OverrideCommand bool     `json:"overrideCommand,omitempty" yaml:"overrideCommand,omitempty"`

	Init       bool     `json:"init,omitempty" yaml:"init,omitempty"`
	Privileged bool     `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	CapAdd     []string `json:"capAdd,omitempty" yaml:"capAdd,omitempty"`

	SecurityOpt []string `json:"securityOpt,omitempty" yaml:"securityOpt,omitempty"`
	Devices     []string `json:"devices,omitempty" yaml:"devices,omitempty"`

	HostRequirements            *HostRequirements         `json:"hostRequirements,omitempty" yaml:"hostRequirements,omitempty" validate:"omitempty"`
	OverrideFeatureInstallOrder []string                  `json:"overrideFeatureInstallOrder,omitempty" yaml:"overrideFeatureInstallOrder,omitempty"`
	Features                    map[string]map[string]any `json:"features,omitempty" yaml:"features,omitempty"`

	InitializeCommand    *CommandValue `json:"initializeCommand,omitempty" yaml:"initializeCommand,omitempty"`
	OnCreateCommand      *CommandValue `json:"onCreateCommand,omitempty" yaml:"onCreateCommand,omitempty"`
	UpdateContentCommand *CommandValue `json:"updateContentCommand,omitempty" yaml:"updateContentCommand,omitempty"`
	PostCreateCommand    *CommandValue `json:"postCreateCommand,omitempty" yaml:"postCreateCommand,omitempty"`
	PostStartCommand     *CommandValue `json:"postStartCommand,omitempty" yaml:"postStartCommand,omitempty"`
	PostAttachCommand    *CommandValue `json:"postAttachCommand,omitempty" yaml:"postAttachCommand,omitempty"`
	WaitFor              string        `json:"waitFor,omitempty" yaml:"waitFor,omitempty" validate:"omitempty,oneof=initializeCommand onCreateCommand updateContentCommand postCreateCommand postStartCommand"`

	Watch          *WatchConfig    `json:"watch,omitempty" yaml:"watch,omitempty" validate:"omitempty"`
	Customizations *Customizations `json:"customizations,omitempty" yaml:"customizations,omitempty" validate:"omitempty"`

	Secrets map[string]Secret `json:"secrets,omitempty" yaml:"secrets,omitempty" validate:"omitempty,dive"`

	ShutdownAction string `json:"shutdownAction,omitempty" yaml:"shutdownAction,omitempty" validate:"omitempty,oneof=none stopContainer stopCompose"`
}

// BuildConfig defines parameters for building a dev container image.
type BuildConfig struct {
	Dockerfile string            `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty" validate:"required"`
	Context    string            `json:"context,omitempty" yaml:"context,omitempty" validate:"omitempty"`
	Args       map[string]string `json:"args,omitempty" yaml:"args,omitempty" validate:"omitempty"`
	Target     string            `json:"target,omitempty" yaml:"target,omitempty" validate:"omitempty"`
	CacheFrom  []string          `json:"cacheFrom,omitempty" yaml:"cacheFrom,omitempty" validate:"omitempty"`
	Options    []string          `json:"options,omitempty" yaml:"options,omitempty" validate:"omitempty"`
}

// Mount represents a filesystem or volume mount for the container.
type Mount struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty" validate:"required,oneof=bind volume tmpfs"`
	Source string `json:"source,omitempty" yaml:"source,omitempty"`
	Target string `json:"target,omitempty" yaml:"target,omitempty" validate:"required"`
}

// PortAttributes defines additional metadata for a forwarded port.
type PortAttributes struct {
	Label            string `json:"label,omitempty" yaml:"label,omitempty" validate:"omitempty"`
	OnAutoForward    string `json:"onAutoForward,omitempty" yaml:"onAutoForward,omitempty" validate:"omitempty,oneof=notify openBrowser openBrowserOnce openPreview silent ignore"`
	Protocol         string `json:"protocol,omitempty" yaml:"protocol,omitempty" validate:"omitempty,oneof=http https"`
	ElevateIfNeeded  bool   `json:"elevateIfNeeded,omitempty" yaml:"elevateIfNeeded,omitempty"`
	RequireLocalPort bool   `json:"requireLocalPort,omitempty" yaml:"requireLocalPort,omitempty"`
}

// WatchConfig controls which files or processes trigger restarts.
type WatchConfig struct {
	WaitFor []string `json:"waitFor,omitempty" yaml:"waitFor,omitempty" validate:"omitempty"`
	Restart []string `json:"restart,omitempty" yaml:"restart,omitempty" validate:"omitempty"`
}

// Customizations defines editor or IDE specific configurations.
type Customizations struct {
	VSCode     *VSCodeCustomization     `json:"vscode,omitempty" yaml:"vscode,omitempty" validate:"omitempty"`
	Codespaces *CodespacesCustomization `json:"codespaces,omitempty" yaml:"codespaces,omitempty" validate:"omitempty"`
	JetBrains  *JetBrainsCustomization  `json:"jetbrains,omitempty" yaml:"jetbrains,omitempty" validate:"omitempty"`
}

// VSCodeCustomization defines VS Code-specific settings.
type VSCodeCustomization struct {
	Settings   map[string]any `json:"settings,omitempty" yaml:"settings,omitempty" validate:"omitempty"`
	Extensions []string       `json:"extensions,omitempty" yaml:"extensions,omitempty" validate:"omitempty"`
}

// CodespacesCustomization defines GitHub Codespaces-specific settings.
type CodespacesCustomization struct {
	Settings   map[string]any `json:"settings,omitempty" yaml:"settings,omitempty" validate:"omitempty"`
	Extensions []string       `json:"extensions,omitempty" yaml:"extensions,omitempty" validate:"omitempty"`
}

// JetBrainsCustomization defines JetBrains IDE configuration.
type JetBrainsCustomization struct {
	Plugins []string `json:"plugins,omitempty" yaml:"plugins,omitempty" validate:"omitempty"`
}

// Secret defines a reusable secret for builds or runtime.
type Secret struct {
	Description      string `json:"description,omitempty" yaml:"description,omitempty" validate:"omitempty"`
	DocumentationURL string `json:"documentationUrl,omitempty" yaml:"documentationUrl,omitempty" validate:"omitempty,url"`
}

// HostRequirements defines minimum hardware resources the host must provide.
type HostRequirements struct {
	CPUs    int       `json:"cpus,omitempty" yaml:"cpus,omitempty" validate:"omitempty,min=1"`
	Memory  string    `json:"memory,omitempty" yaml:"memory,omitempty"`
	Storage string    `json:"storage,omitempty" yaml:"storage,omitempty"`
	GPU     *GPUValue `json:"gpu,omitempty" yaml:"gpu,omitempty" validate:"omitempty"`
}

// GPURequirement describes GPU resource needs within HostRequirements.
type GPURequirement struct {
	Cores  int    `json:"cores,omitempty" yaml:"cores,omitempty" validate:"omitempty,min=1"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

// TopLevelKeys is the single source of truth for every recognised top-level
// DevContainer field, in canonical display/insertion order.
// Both the edit TUI (allKnownKeys) and the presets package (ListFields) derive
// their lists from this slice so they can never diverge.
var TopLevelKeys = []string{
	"name", "image", "build", "dockerComposeFile", "service", "runServices",
	"workspaceFolder", "workspaceMount", "remoteUser", "containerUser",
	"updateRemoteUserUID", "userEnvProbe",
	"containerEnv", "remoteEnv",
	"forwardPorts", "appPort", "portsAttributes", "otherPortsAttributes",
	"mounts", "runArgs", "overrideCommand",
	"init", "privileged", "capAdd", "securityOpt", "devices",
	"hostRequirements", "overrideFeatureInstallOrder", "features",
	"initializeCommand", "onCreateCommand", "updateContentCommand",
	"postCreateCommand", "postStartCommand", "postAttachCommand", "waitFor",
	"watch", "customizations", "secrets", "shutdownAction",
}
