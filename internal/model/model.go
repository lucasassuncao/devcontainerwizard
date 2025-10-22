// Package model defines the data structures used to represent
// dev container configuration and related types. These structs are
// used for parsing, validation and documentation generation.
package model

type DevContainer struct {
	Schema            string       `json:"$schema,omitempty" yaml:"$schema,omitempty" jsonschema_description:"URL of the JSON schema that describes the format of this file."`
	Name              string       `json:"name,omitempty" yaml:"name,omitempty" validate:"required" jsonschema:"required" jsonschema_description:"Name of the dev container."`
	Image             string       `json:"image,omitempty" yaml:"image,omitempty" validate:"required_without=Build" jsonschema_description:"Docker image to use for the dev container."`
	Build             *BuildConfig `json:"build,omitempty" yaml:"build,omitempty" validate:"omitempty" jsonschema_description:"Configuration for building the image."`
	DockerComposeFile []string     `json:"dockerComposeFile,omitempty" yaml:"dockerComposeFile,omitempty" jsonschema_description:"List of Docker Compose files to use."`
	Service           string       `json:"service,omitempty" yaml:"service,omitempty" validate:"required_with=DockerComposeFile" jsonschema_description:"Specific service to run from Docker Compose."`

	WorkspaceFolder string `json:"workspaceFolder,omitempty" yaml:"workspaceFolder,omitempty" jsonschema_description:"Path to the workspace folder inside the container."`
	WorkspaceMount  string `json:"workspaceMount,omitempty" yaml:"workspaceMount,omitempty" jsonschema_description:"Mount type for the workspace folder."`
	RemoteUser      string `json:"remoteUser,omitempty" yaml:"remoteUser,omitempty" jsonschema_description:"User to use inside the container."`
	UserEnvProbe    string `json:"userEnvProbe,omitempty" yaml:"userEnvProbe,omitempty" validate:"omitempty,oneof=none loginShell loginInteractiveShell interactiveShell" jsonschema_description:"Command to detect the default user inside the container."`

	// Environment variables
	ContainerEnv map[string]string `json:"containerEnv,omitempty" yaml:"containerEnv,omitempty" jsonschema_description:"Environment variables to set in the container."`
	RemoteEnv    map[string]string `json:"remoteEnv,omitempty" yaml:"remoteEnv,omitempty" jsonschema_description:"Environment variables for remote connections (like SSH)."`
	LocalEnv     map[string]string `json:"-" yaml:"localEnv,omitempty" jsonschema_description:"Environment variables local to the host (not exported to JSON)."`

	ForwardPorts         []interface{}              `json:"forwardPorts,omitempty" yaml:"forwardPorts,omitempty" jsonschema_description:"Ports that are forwarded from the container to the local machine. Can be an integer port number, or a string of the format \"host:port_number\""`
	PortsAttributes      map[string]*PortAttributes `json:"portsAttributes,omitempty" yaml:"portsAttributes,omitempty" validate:"omitempty" jsonschema_description:"Additional attributes for forwarded ports."`
	OtherPortsAttributes *PortAttributes            `json:"otherPortsAttributes,omitempty" yaml:"otherPortsAttributes,omitempty" validate:"omitempty" jsonschema_description:"Default attributes applied to all forwarded ports not defined in portsAttributes."`
	Mounts               []Mount                    `json:"mounts,omitempty" yaml:"mounts,omitempty" validate:"omitempty,dive" jsonschema_description:"Mount points inside the container."`

	RunArgs        []string `json:"runArgs,omitempty" yaml:"runArgs,omitempty" jsonschema_description:"Additional arguments to pass to 'docker run'."`
	StartupCommand string   `json:"startupCommand,omitempty" yaml:"startupCommand,omitempty" jsonschema_description:"Command to run on container startup."`
	Command        string   `json:"command,omitempty" yaml:"command,omitempty" jsonschema_description:"Command to run inside the container instead of the default CMD."`
	Entrypoint     string   `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty" jsonschema_description:"Entrypoint to override in the container."`

	Init       bool     `json:"init,omitempty" yaml:"init,omitempty" jsonschema_description:"Whether to run an init process inside the container."`
	Privileged bool     `json:"privileged,omitempty" yaml:"privileged,omitempty" jsonschema_description:"Run the container in privileged mode."`
	CapAdd     []string `json:"capAdd,omitempty" yaml:"capAdd,omitempty" jsonschema_description:"Linux capabilities to add to the container."`
	CapDrop    []string `json:"capDrop,omitempty" yaml:"capDrop,omitempty" jsonschema_description:"Linux capabilities to drop from the container."`

	SecurityOpt []string `json:"securityOpt,omitempty" yaml:"securityOpt,omitempty" jsonschema_description:"Security options for the container."`
	Devices     []string `json:"devices,omitempty" yaml:"devices,omitempty" jsonschema_description:"Devices to expose to the container."`

	OverrideFeatureInstallOrder []string                  `json:"overrideFeatureInstallOrder,omitempty" yaml:"overrideFeatureInstallOrder,omitempty" jsonschema_description:"Order to install features inside the container, overriding defaults."`
	Features                    map[string]map[string]any `json:"features,omitempty" yaml:"features,omitempty" jsonschema_description:"Features to install in the container and their options."`

	OnCreateCommand   string `json:"onCreateCommand,omitempty" yaml:"onCreateCommand,omitempty" jsonschema_description:"Command to run after the container is created."`
	PostCreateCommand string `json:"postCreateCommand,omitempty" yaml:"postCreateCommand,omitempty" jsonschema_description:"Command to run after the container is created and initialized."`
	PostStartCommand  string `json:"postStartCommand,omitempty" yaml:"postStartCommand,omitempty" jsonschema_description:"Command to run after the container starts."`
	PostAttachCommand string `json:"postAttachCommand,omitempty" yaml:"postAttachCommand,omitempty" jsonschema_description:"Command to run after attaching to the container."`

	Watch          *WatchConfig    `json:"watch,omitempty" yaml:"watch,omitempty" validate:"omitempty" jsonschema_description:"Configuration for files/processes to watch for restarts."`
	Customizations *Customizations `json:"customizations,omitempty" yaml:"customizations,omitempty" validate:"omitempty" jsonschema_description:"Editor/IDE customizations inside the container."`

	Secrets map[string]Secret `json:"secrets,omitempty" yaml:"secrets,omitempty" validate:"omitempty" jsonschema_description:"Secrets to pass to the container."`

	ShutdownAction string `json:"shutdownAction,omitempty" yaml:"shutdownAction,omitempty" validate:"omitempty,oneof=none stopContainer stopCompose" jsonschema_description:"Action to take when the container is stopped."`
}

// BuildConfig defines parameters for building a dev container image.
type BuildConfig struct {
	Dockerfile string            `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty" validate:"required" jsonschema:"required,default=Dockerfile" jsonschema_description:"Path to the Dockerfile to use for building the image."`
	Context    string            `json:"context,omitempty" yaml:"context,omitempty" validate:"required" jsonschema:"required,default=." jsonschema_description:"Build context directory."`
	Args       map[string]string `json:"args,omitempty" yaml:"args,omitempty" validate:"omitempty" jsonschema_description:"Build arguments as key-value pairs."`
	Target     string            `json:"target,omitempty" yaml:"target,omitempty" validate:"omitempty" jsonschema_description:"Target stage for multi-stage Docker builds."`
	CacheFrom  []string          `json:"cacheFrom,omitempty" yaml:"cacheFrom,omitempty" validate:"omitempty" jsonschema_description:"List of images to cache from."`
	Output     string            `json:"output,omitempty" yaml:"output,omitempty" validate:"omitempty" jsonschema_description:"Output location of the build."`
	SSH        []string          `json:"ssh,omitempty" yaml:"ssh,omitempty" validate:"omitempty" jsonschema_description:"SSH mount sources to use during build."`
	Secrets    []BuildSecret     `json:"secrets,omitempty" yaml:"secrets,omitempty" validate:"omitempty,dive" jsonschema_description:"Secrets to pass to the build process."`
}

// BuildSecret represents a secret used during build.
type BuildSecret struct {
	ID  string `json:"id,omitempty" yaml:"id,omitempty" validate:"required" jsonschema_description:"Identifier for the secret."`
	Src string `json:"src,omitempty" yaml:"src,omitempty" validate:"required" jsonschema_description:"Path or source of the secret."`
}

// Mount represents a filesystem or volume mount for the container.
type Mount struct {
	Type        string `json:"type,omitempty" yaml:"type,omitempty" validate:"required,oneof=bind volume" jsonschema:"required" jsonschema_description:"Type of mount (e.g., bind, volume)."`
	Source      string `json:"source,omitempty" yaml:"source,omitempty" validate:"required" jsonschema:"required" jsonschema_description:"Source path of the mount."`
	Target      string `json:"target,omitempty" yaml:"target,omitempty" validate:"required" jsonschema:"required" jsonschema_description:"Target path inside the container."`
	Consistency string `json:"consistency,omitempty" yaml:"consistency,omitempty" validate:"omitempty,oneof=cached delegated consistent" jsonschema_description:"Consistency mode for the mount (e.g., cached, delegated, consistent)."`
	ReadOnly    bool   `json:"readonly,omitempty" yaml:"readonly,omitempty" jsonschema_description:"Whether the mount is read-only."`
}

// PortAttributes defines additional metadata for a forwarded port.
type PortAttributes struct {
	Label         string `json:"label,omitempty" yaml:"label,omitempty" validate:"omitempty" jsonschema_description:"Human-readable label for the port."`
	OnAutoForward string `json:"onAutoForward,omitempty" yaml:"onAutoForward,omitempty" validate:"omitempty,oneof=notify openBrowser openBrowserOnce openPreview silent ignore" jsonschema_description:"Behavior when the port is auto-forwarded (notify, openBrowser, ignore)."`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty" validate:"omitempty,oneof=http https" jsonschema_description:"Network protocol (tcp/udp) for the port."`
}

// WatchConfig controls which files or processes trigger restarts.
type WatchConfig struct {
	WaitFor []string `json:"waitFor,omitempty" yaml:"waitFor,omitempty" validate:"omitempty" jsonschema_description:"List of processes or files to wait for before starting."`
	Restart []string `json:"restart,omitempty" yaml:"restart,omitempty" validate:"omitempty" jsonschema_description:"List of files or events that trigger a restart."`
}

// Customizations defines editor or IDE specific configurations.
type Customizations struct {
	VSCode     *VSCodeCustomization     `json:"vscode,omitempty" yaml:"vscode,omitempty" validate:"omitempty" jsonschema_description:"VS Code specific customizations."`
	Codespaces *CodespacesCustomization `json:"codespaces,omitempty" yaml:"codespaces,omitempty" validate:"omitempty" jsonschema_description:"Codespaces specific customizations."`
	JetBrains  *JetBrainsCustomization  `json:"jetbrains,omitempty" yaml:"jetbrains,omitempty" validate:"omitempty" jsonschema_description:"JetBrains IDE specific customizations."`
	Neovim     *NeovimCustomization     `json:"neovim,omitempty" yaml:"neovim,omitempty" validate:"omitempty" jsonschema_description:"Neovim specific customizations."`
}

// VSCodeCustomization defines VS Code-specific settings.
type VSCodeCustomization struct {
	Settings   map[string]any `json:"settings,omitempty" yaml:"settings,omitempty" validate:"omitempty" jsonschema_description:"Key-value settings for VS Code."`
	Extensions []string       `json:"extensions,omitempty" yaml:"extensions,omitempty" validate:"omitempty" jsonschema_description:"List of VS Code extensions to install."`
	RemoteUser string         `json:"remoteUser,omitempty" yaml:"remoteUser,omitempty" validate:"omitempty" jsonschema_description:"Remote user for VS Code container setup."`
}

// CodespacesCustomization defines GitHub Codespaces-specific settings.
type CodespacesCustomization struct {
	Settings   map[string]any `json:"settings,omitempty" yaml:"settings,omitempty" validate:"omitempty" jsonschema_description:"Key-value settings for Codespaces."`
	Extensions []string       `json:"extensions,omitempty" yaml:"extensions,omitempty" validate:"omitempty" jsonschema_description:"List of Codespaces extensions to install."`
}

// JetBrainsCustomization defines JetBrains IDE configuration.
type JetBrainsCustomization struct {
	Plugins []string `json:"plugins,omitempty" yaml:"plugins,omitempty" validate:"omitempty" jsonschema_description:"List of JetBrains plugins to install."`
}

// NeovimCustomization defines Neovim-specific configuration.
type NeovimCustomization struct {
	Plugins []string `json:"plugins,omitempty" yaml:"plugins,omitempty" validate:"omitempty" jsonschema_description:"List of Neovim plugins to install."`
}

// Secret defines a reusable secret for builds or runtime.
type Secret struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" validate:"required" jsonschema_description:"Human-readable description of the secret."`
	Default     string `json:"default,omitempty" yaml:"default,omitempty" validate:"omitempty" jsonschema_description:"Default value for the secret if none is provided."`
}

// GetAllTypes returns all model types for documentation generation
func GetAllTypes() []interface{} {
	return []interface{}{
		DevContainer{},
		BuildConfig{},
		Mount{},
		PortAttributes{},
		Secret{},
		WatchConfig{},
		Customizations{},
	}
}
