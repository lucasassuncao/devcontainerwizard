<!-- markdownlint-disable MD033 -->
<p align="center">
  <img src="../devcontainerwizard.png" alt="devcontainerwizard logo" width="300" height="300">
</p>

<h1 align="center">devcontainerwizard</h1>

<p align="center">
  <strong>A short, human-friendly YAML in; a valid <code>devcontainer.json</code> out.</strong><br>
  Every field documented, validated before it is written.
</p>
<!-- markdownlint-enable MD033 -->

You write this:

```yaml
name: my-project
image: mcr.microsoft.com/devcontainers/base:ubuntu
remoteUser: vscode
forwardPorts:
  - 3000
```

You run `devcontainerwizard convert`, and VS Code offers to reopen the folder in a container.

## What it does

|  | |
|---|---|
| 📦 **Templates** | `init --template` starts from a minimal image, a Dockerfile, Docker Compose, Go, Python, Rust, DevOps, or a fully annotated example. |
| ✍️ **Two-panel editor** | Check a field on and it is inserted with a working snippet, already indented, with its description, default and allowed values in the hint panel. |
| 🧩 **Presets** | Ready-made values for any field, or a whole template for the document. |
| ✅ **Validated twice** | The editor validates on every save and `convert` again before writing, including the cross-field rules a schema cannot express. |
| 🧠 **Editor completion** | A generated JSON Schema for `config.yaml`, so VS Code completes and validates it as you type. |
| 📚 **One source of truth** | The hint panel, the validators, the reference pages and the JSON Schema are all generated from the same field descriptions. |

## Install

```bash
go install github.com/lucasassuncao/devcontainerwizard@latest
```

Or grab a binary from [Releases](https://github.com/lucasassuncao/devcontainerwizard/releases). `devcontainerwizard self-update` upgrades it in place.

## First run

```bash
devcontainerwizard init -t golang       # write config.yaml from a template
devcontainerwizard edit                 # optional: refine it in the TUI
devcontainerwizard convert              # write .devcontainer/devcontainer.json
```

See [Getting Started](/GETTING-STARTED.md) for the full walkthrough.

## Where to go next

| Page | For |
| --- | --- |
| [Templates](/TEMPLATES.md) | the built-in starting points |
| [Editor](/EDIT.md) | the TUI, its panels and every key |
| [Converting](/CONVERT.md) | producing `devcontainer.json`, and where to put it |
| [Commands](/COMMANDS.md) | every command and flag |
| [Configuration](/CONFIGURATION.md) | what goes in `config.yaml` |
| [Cookbook](/COOKBOOK.md) | recipes for common setups |
| [Troubleshooting](/TROUBLESHOOTING.md) and [FAQ](/FAQ.md) | when something does not work |
| [Attributes](/devcontainerwizard/attributes/devcontainer.md) | the generated reference for every field |
| [Development](/DEVELOPMENT.md) | building, testing, adding a field or a template |
| [Architecture](/ARCHITECTURE.md) | the single source of truth, and where a change belongs |
| [Releasing](/RELEASING.md) | tags, ldflags, and the self-update coupling |
