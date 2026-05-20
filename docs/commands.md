# Command Reference

## init

Create a new `config.yaml` file from a template. Running `init` without flags prints an error with the available templates. Use `--list` to browse them without creating a file.

```bash
devcontainerwizard init [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--template` | `-t` | — | Template to use (required). See `--list` for available names. |
| `--output` | `-o` | `config.yaml` | Output file path |
| `--list` | `-l` | false | Print available templates and exit |
| `--force` | `-f` | false | Overwrite an existing output file |

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
devcontainerwizard edit [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `config.yaml` | Path to the config file |

**Key bindings:**

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Move cursor |
| `Space` | Add or edit a block |
| `d` | Delete the selected block |
| `Tab` | Toggle focus to the YAML preview panel |
| `Ctrl+S` | Save to file |
| `q` | Quit (prompts if there are unsaved changes) |

**Inside a block overlay (fields panel):**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move field cursor |
| `Space` | Toggle a field on/off |
| `p` | Open the preset picker |
| `Tab` | Switch to the YAML editor panel |
| `Ctrl+S` | Confirm and close |
| `Esc` | Cancel |

**Preset picker (`p`):**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate presets |
| `Enter` | Apply selected preset |
| `Esc` | Close without applying |

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

**Key bindings:**

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Navigate topics / scroll content |
| `PgUp` / `PgDn` | Half-page scroll in the content panel |
| `Tab` | Switch between the Topics and content panels |
| `q` | Quit |

---

## show-examples

Browse built-in YAML presets for every devcontainer config field in a two-panel TUI.
Use it to discover ready-made values before opening the editor.

```bash
devcontainerwizard show-examples
```

**Navigation:**

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Move cursor |
| `Enter` / `→` / `l` | Drill into a field's preset list |
| `Esc` / `←` / `h` | Go back to the fields list |
| `Tab` | Switch focus to the YAML preview panel |
| `q` | Quit |

The left panel lists all 45 config fields. Selecting a field shows its available presets; selecting a preset renders the corresponding YAML block on the right with syntax highlighting.

---

## self-update

Update `devcontainerwizard` to the latest GitHub release.

```bash
devcontainerwizard self-update
```
