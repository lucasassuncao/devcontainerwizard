# Getting Started

## Install

```bash
go install github.com/lucasassuncao/devcontainerwizard@latest
```

Or grab a binary from the [releases page](https://github.com/lucasassuncao/devcontainerwizard/releases).

Or build from source:

```bash
git clone https://github.com/lucasassuncao/devcontainerwizard.git
cd devcontainerwizard
go install
```

Once installed, keep it current with:

```bash
devcontainerwizard self-update
```

---

## 1. Create a config

Every project starts from a template. List them first:

```bash
devcontainerwizard init --list
```

Then pick one:

```bash
devcontainerwizard init --template golang
```

That writes `config.yaml` in the current directory. Use `--output` to put it
somewhere else, and `--force` to overwrite an existing file.

See [Templates](/TEMPLATES.md) for what each one contains.

---

## 2. Edit it

```bash
devcontainerwizard edit
```

The left panel lists the top-level blocks. **Enter** opens one; inside, **space**
checks a field on and inserts a working snippet for it. **h** shows the hint
panel with that field's description, default, and allowed values. **Ctrl+S**
saves.

You never have to use the editor — `config.yaml` is ordinary YAML and any text
editor will do. The TUI exists so you do not have to remember the field names.

See [Editor](/EDIT.md) for the full keybinding reference.

---

## 3. Convert

```bash
devcontainerwizard convert
```

This validates the config and writes `.devcontainer/devcontainer.json`. Open the
folder in VS Code and it offers to reopen in the container.

If validation fails, nothing is written and the errors name the offending
fields. See [Converting](/CONVERT.md).

---

## Editor completion for config.yaml

`devcontainerwizard` ships a JSON Schema for `config.yaml`. Point your editor at
it and you get completion and inline validation while typing:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/lucasassuncao/devcontainerwizard/main/docs/devcontainerwizard/schema/devcontainer.schema.json
name: my-project
image: mcr.microsoft.com/devcontainers/base:ubuntu
```

---

## Where to go next

- [Configuration](/CONFIGURATION.md) — the shape of `config.yaml`, block by block
- [Cookbook](/COOKBOOK.md) — worked examples for common stacks
- [Commands](/COMMANDS.md) — every flag of every command
