# Command Reference

## init

Create a new `config.yaml` file.

```bash
devcontainerwizard init [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--template` | `-t` | `image` | Template to use: `image`, `dockerfile`, `dockercompose`, `full`, `golang` |
| `--interactive` | `-i` | false | Prompt for common fields instead of using a template |
| `--force` | `-f` | false | Overwrite an existing `config.yaml` |

**Templates:**

| Name | Description |
|------|-------------|
| `image` | Minimal config with a Docker image |
| `dockerfile` | Custom Dockerfile with build config |
| `dockercompose` | Docker Compose multi-service setup |
| `full` | Complete example with all options |
| `golang` | Optimised setup for Go development |

---

## edit

Open a two-panel TUI to add, remove, and edit top-level blocks in a config YAML file.

```bash
devcontainerwizard edit [file]
```

If `[file]` is omitted, opens `config.yaml` in the current directory.

**Key bindings:**

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Move cursor |
| `Space` | Add or edit a block |
| `d` | Delete the selected block |
| `Tab` | Toggle focus to the YAML preview panel |
| `Ctrl+S` | Save to file |
| `q` | Quit (prompts if there are unsaved changes) |

---

## convert

Convert `config.yaml` to `.devcontainer/devcontainer.json`.

```bash
devcontainerwizard convert [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `config.yaml` | Path to the config file |
| `--output` | `-o` | `.devcontainer` | Output directory |

---

## show-docs

Browse configuration documentation in the terminal with syntax-highlighted markdown.

```bash
devcontainerwizard show-docs
```

---

## self-update

Update `devcontainerwizard` to the latest GitHub release.

```bash
devcontainerwizard self-update
```
