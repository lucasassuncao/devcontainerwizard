# DevContainer Wizard 🧙‍♂️

A powerful CLI tool to generate and manage DevContainer configurations with ease.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

## ✨ Features

- 🚀 **Quick Init** - Generate config.yaml with interactive prompts or pre-built templates
- 📝 **YAML to JSON** - Convert user-friendly YAML configs to devcontainer.json
- ✅ **Validation** - Catch configuration errors before building containers
- 📚 **Documentation** - Auto-generate markdown docs from your models
- 🎨 **Templates** - Choose from basic, docker, compose, or full templates
- 🔧 **Customizable** - Full support for VS Code, Codespaces, JetBrains, and Neovim

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/devcontainer-wizard.git
cd devcontainer-wizard

# Build and install
go build -o devcontainer
sudo mv devcontainer /usr/local/bin/

# Or install directly
go install
```

### Using Go Install

```bash
go install github.com/yourusername/devcontainer-wizard@latest
```

## 🚀 Quick Start

### 1. Initialize a New Configuration

```bash
# Interactive mode (recommended for first time)
devcontainer init -i

# Or use a template
devcontainer init -t image
devcontainer init -t dockerfile
devcontainer init -t dockercompose
devcontainer init -t full
```

### 2. Edit config.yaml (Optional)

```yaml
name: my-project
image: mcr.microsoft.com/devcontainers/base:ubuntu

containerEnv:
  NODE_ENV: development

forwardPorts:
  - 3000

customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
```

### 3. Generate DevContainer

```bash
# Generate .devcontainer/devcontainer.json
devcontainer

# Or specify custom paths
devcontainer -c my-config.yaml -o .devcontainer
```

### 4. Open in VS Code

```bash
# Open your project in VS Code
code .

# VS Code will prompt to reopen in container
```

## 📖 Commands

### `init` - Create Configuration

Create a new `config.yaml` file with interactive prompts or templates.

```bash
# Interactive mode
devcontainer init -i

# Use specific template
devcontainer init -t docker

# Force overwrite existing config
devcontainer init -f

# Combine flags
devcontainer init -i -f
```

**Available Templates:**
- `image` - Minimal configuration with Docker image
- `dockerfile` - Custom Dockerfile with build configuration
- `dockercompose` - Docker Compose multi-service setup
- `full` - Complete example with all options

### `show-docs` - View Documentation

View documentation interactively in your terminal.

```bash
devcontainer show-docs
```

Features:
- 📜 Scrollable with arrow keys or j/k
- 🎨 Syntax-highlighted markdown
- 📱 Responsive to terminal size
- ⌨️  Press 'q' to quit

### Default Command - Convert

Convert `config.yaml` to `.devcontainer/devcontainer.json`.

```bash
# Use default config.yaml
devcontainer

# Specify custom config file
devcontainer -c path/to/config.yaml

# Specify custom output directory
devcontainer -o .devcontainer-custom
```

## 📋 Configuration Reference

### Environment Variables

DevContainer supports three types of environment variables:

#### `containerEnv` (Recommended)
Variables available inside the container.

```yaml
containerEnv:
  NODE_ENV: development
  API_URL: http://localhost:3000
```

#### `remoteEnv`
Variables for remote connections (SSH, VS Code Remote).

```yaml
remoteEnv:
  PATH: "${containerEnv:PATH}:/custom/bin"
```

#### `localEnv`
Variables for local processing only (not exported to JSON).

```yaml
localEnv:
  MY_SECRET: ${env:MY_SECRET}
```

### Port Forwarding

```yaml
forwardPorts:
  - 3000
  - 5432

portsAttributes:
  "3000":
    label: "Web Server"
    onAutoForward: openBrowser
    protocol: http
  "5432":
    label: "PostgreSQL"
    onAutoForward: silent
```

### Features

Add pre-built features from the [Dev Containers Features](https://containers.dev/features) catalog:

```yaml
features:
  ghcr.io/devcontainers/features/git:1: {}
  ghcr.io/devcontainers/features/github-cli:1:
    version: latest
  ghcr.io/devcontainers/features/docker-in-docker:2:
    moby: true
```

### Lifecycle Hooks

```yaml
onCreateCommand: echo "Container created"
postCreateCommand: npm install
postStartCommand: npm run dev
postAttachCommand: echo "Welcome!"
```

### Mounts

```yaml
mounts:
  - type: bind
    source: ./data
    target: /data
  - type: volume
    source: node_modules
    target: /workspace/node_modules
```

### VS Code Customizations

```yaml
customizations:
  vscode:
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
      - eamodio.gitlens
    settings:
      editor.formatOnSave: true
      editor.defaultFormatter: esbenp.prettier-vscode
```

## 🔧 Development

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Building

```bash
# Build binary
go build -o devcontainer

# Build with optimizations
go build -ldflags="-s -w" -o devcontainer

# Install locally
go install
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/devcontainer/...
```

### Linting

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

### Using Makefile

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run linter
make clean          # Clean build artifacts
make install        # Install globally
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go best practices and idioms
- Update documentation for user-facing changes
- Keep commits atomic and well-described

## 🐛 Troubleshooting

### Config file already exists

```bash
# Use --force to overwrite
devcontainer init -f
```

### Validation errors

```bash
# Check your config.yaml syntax
yamllint config.yaml

# Try generating to see detailed errors
devcontainer -c config.yaml
```

### Container won't start

1. Check Docker is running: `docker ps`
2. Validate your config: `devcontainer -c config.yaml`
3. Check VS Code DevContainer logs: View → Output → Dev Containers

### Port already in use

```yaml
# Change the port in config.yaml
forwardPorts:
  - 3001  # Changed from 3000
```

## 📚 Resources

- [Dev Containers Documentation](https://containers.dev/)
- [Dev Containers Specification](https://containers.dev/implementors/spec/)
- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
- [Dev Container Features](https://containers.dev/features)
- [Dev Container Templates](https://containers.dev/templates)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Dev Containers Specification](https://containers.dev/)
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Koanf](https://github.com/knadh/koanf) - Configuration management
- [Validator](https://github.com/go-playground/validator) - Struct validation
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI

## 📮 Support

- 🐛 [Report a bug](https://github.com/lucasassuncao/devcontainerwizard/issues)
- 💡 [Request a feature](https://github.com/lucasassuncao/devcontainerwizard/issues)
- 💬 [Discussions](https://github.com/lucasassuncao/devcontainerwizard/discussions)

---

Made with ❤️ by [Lucas Assunção da Silva](https://github.com/lucasassuncao)

**⭐ If you find this project useful, please consider giving it a star!**