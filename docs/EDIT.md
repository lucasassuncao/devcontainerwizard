# Interactive Config Editor

`devcontainerwizard edit` opens a two-panel TUI editor for your config file. It
is the fastest way to build a devcontainer config without memorising the
specification, and it validates on every save.

```bash
devcontainerwizard edit
```

---

## Layout

**Left panel — block list**\
Every top-level key of the devcontainer schema: `name`, `image`, `build`,
`forwardPorts`, `customizations`, and so on. Keys already present in your file
are marked; the rest are available to add. Navigate with **↑ / ↓**, open one
with **Enter**.

**Right panel — YAML preview**\
The document as it stands, syntax-highlighted, updated on every edit. **Tab**
moves focus here so you can scroll it.

**Block editor (Enter)**\
Opens over the list. The left side is the field tree for that block, the right
side is the YAML for the block alone. Checking a field on with **space**
inserts a working snippet for it — already indented, with a plausible value —
rather than an empty key.

**Hint panel (h)**\
Slides open next to the field tree with that field's description, type,
default, allowed values, and a YAML example. This is the same metadata the
generated [attribute reference](/devcontainerwizard/attributes/devcontainer.md)
is built from, so the two can never disagree.

---

## Keybindings

### Root list

| Key | Action |
|---|---|
| **↑ / ↓** | Move between blocks |
| **Enter** | Open the selected block |
| **p** | Open the whole-document template picker |
| **/** | Filter the list |
| **Tab** | Focus the YAML preview |
| **h** | Toggle the hint panel |
| **Ctrl+D** | Delete the selected block |
| **Ctrl+L** | Validate now |
| **Ctrl+R** | Reload the file from disk |
| **Ctrl+S** | Save |
| **Ctrl+E** | Save, then convert to `.devcontainer/devcontainer.json` |
| **Ctrl+U** / **Ctrl+Y** | Undo / redo |
| **Esc** | Quit |

### Inside a block

| Key | Action |
|---|---|
| **↑ / ↓** | Move between fields |
| **→ / ←** | Expand / collapse a nested field |
| **space** | Check a field on or off |
| **Enter** | Edit the selected value |
| **p** | Open the preset picker for this field |
| **a** | Append an entry to a list |
| **h** | Toggle the hint panel |
| **Tab** | Switch to the block's YAML pane |
| **Ctrl+S** | Confirm and close |
| **Esc** | Cancel |

---

## Presets and templates

**p inside a block** opens the preset picker for the field under the cursor:
ready-made values drawn from the same set `devcontainerwizard show-examples`
browses. Applying one replaces the field's current value.

**p on the root list** opens the whole-document template picker, backed by the
same templates `init --template` writes. Applying one replaces the entire
document, so it is most useful on a fresh or throwaway file.

---

## Saving and validation

**Ctrl+S** runs every validator before writing. Failures are listed in a modal
and the file is left untouched. The rules are:

- **Per-field constraints** — required fields, allowed values, minimums, string
  formats, and deprecation warnings, all declared alongside the schema.
- **Cross-field rules** — `image`, `build`, and `dockerComposeFile` are mutually
  exclusive; `service` and `runServices` require `dockerComposeFile`.

Pass `--no-validate-on-save` to write anyway. A warning is shown; the errors are
not silently dropped.

**Ctrl+E** saves the same way, then runs the saved file through `convert` and
writes `.devcontainer/devcontainer.json`, the path VS Code looks for first. The
JSON is generated, so an existing one is overwritten, as `convert --force` would.
Nothing is converted when the save is cancelled or validation fails, and a file
with no unsaved changes is converted without being written again.

---

## Writing to a different file

`--output` loads one file and saves to another:

```bash
devcontainerwizard edit --config config.yaml --output ../other/config.yaml
```

Useful for deriving a new project's config from an existing one without
touching the original.

---

## Themes

```bash
devcontainerwizard edit --theme grape
devcontainerwizard edit --list-themes   # browse them all, interactively
```

`--list-themes` opens a tabbed browser showing each palette against real UI
chrome. The default theme is `plain`.

---

## Reporting a bug in the editor

`--dump` records every keystroke and every editor action to a JSONL trace and
prints the path on exit:

```bash
devcontainerwizard edit --dump
```

Attach that file to the issue. Use `--dump-path ./trace.jsonl` to choose where
it lands.

---

See [Commands](/COMMANDS.md) for the full flag reference.
