# Command Reference

## Global flags

Available on every command:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `config.yaml` | Path to the config file that `edit` loads and `convert` reads |

## Commands

| Command | Description |
|---------|-------------|
| [`init`](#init) | Create a new `config.yaml` from a template |
| [`edit`](#edit) | Open the interactive TUI editor |
| [`convert`](#convert) | Convert `config.yaml` to `devcontainer.json` |
| [`show-examples`](#show-examples) | Browse the built-in presets in the terminal |
| [`self-update`](#self-update) | Update to a GitHub release |

---

## init

Create a new `config.yaml` from a template. Running `init` without `--template`
prints the available templates and exits non-zero. Use `--list` to browse them
without creating a file.

```bash
devcontainerwizard init --template golang
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--template` | `-t` | — | Template to use (required). See `--list` for available names. |
| `--output` | `-o` | `config.yaml` | Output file path |
| `--list` | `-l` | false | Print available templates and exit |
| `--force` | `-f` | false | Overwrite an existing output file |

See [Templates](/TEMPLATES.md) for what each template contains.

---

## edit

Open a two-panel TUI to add, remove, and edit blocks in a config file. See
[Editor](/EDIT.md) for the layout and keybindings.

```bash
devcontainerwizard edit
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | — | Save to this file instead of the loaded config (the load path is unchanged) |
| `--theme` | | `plain` | Theme name |
| `--list-themes` | | false | Browse available themes in an interactive, tabbed terminal UI |
| `--no-save-confirm` | | false | Skip the "Save changes?" confirmation dialog |
| `--no-delete-confirm` | | false | Skip the "Remove block?" confirmation dialog |
| `--no-validate-on-save` | | false | Allow saving even when validators report errors (a warning is shown) |
| `--dump` | | false | Record every editor action to a JSONL trace file for bug reports |
| `--dump-path` | | — | Write the trace to this file instead of a temp file (implies `--dump`) |

---

## convert

Convert `config.yaml` to a `devcontainer.json` file. By default it writes to
`.devcontainer/devcontainer.json`. When the output path is not a canonical
devcontainer location, a warning goes to stderr — VS Code and the devcontainer
CLI only auto-detect `.devcontainer/devcontainer.json` and `.devcontainer.json`.

```bash
devcontainerwizard convert
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `.devcontainer/devcontainer.json` | Output `devcontainer.json` file path |
| `--force` | `-f` | false | Overwrite an existing output file |

See [Converting](/CONVERT.md).

---

## show-examples

Browse the built-in YAML presets for every config field in a two-panel TUI. Use
it to discover ready-made values before opening the editor.

```bash
devcontainerwizard show-examples
```

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor |
| `Enter` / `→` | Drill into a field's preset list |
| `Esc` / `←` | Back to the field list |
| `Tab` | Focus the YAML preview panel |
| `q` | Quit |

The left panel lists every top-level config field. Selecting a field shows its
presets; selecting a preset renders the corresponding YAML on the right.

---

## self-update

Download a release of `devcontainerwizard` from GitHub and replace the running
binary. The old binary is kept as `devcontainerwizard.old` until the next run.

```bash
devcontainerwizard self-update
```

| Flag | Default | Description |
|------|---------|-------------|
| `--repo` | build-time default | GitHub repository in `owner/repo` format |
| `--version` | — | Install this specific release tag instead of the latest |
| `--list` | false | List available releases and exit |
| `--prerelease` | false | Include rc/beta/alpha releases |
| `--limit` | 20 | Maximum releases to show with `--list` (max 100) |

---

## generate-docs

Hidden maintenance command. Regenerates the reference pages, the example pages,
and the JSON Schema under `docs/devcontainerwizard/` from the schema metadata.
Run it through `make docs`, which also refreshes the Go package READMEs — CI
fails if the committed output is stale.

```bash
make docs
```
