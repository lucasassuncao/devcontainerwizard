# TUI Edit — Developer Reference

The `internal/tui/edit` package implements the interactive YAML editor launched by `devcontainerwizard edit [file]`. It is built on [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss).

---

## Table of Contents

1. [High-level architecture](#1-high-level-architecture)
2. [File map](#2-file-map)
3. [Data flow](#3-data-flow)
4. [Key bindings reference](#4-key-bindings-reference)
5. [YAML manipulation](#5-yaml-manipulation)
6. [Left panel — list model](#6-left-panel--list-model)
7. [Right panel — preview / editor](#7-right-panel--preview--editor)
8. [Overlay](#8-overlay)
9. [Guided mode — field definitions](#9-guided-mode--field-definitions)
10. [Guided mode — templates](#10-guided-mode--templates)
11. [Layout and sizing](#11-layout-and-sizing)
12. [Styles](#12-styles)
13. [How to extend](#13-how-to-extend)

---

## 1. High-level architecture

```
cmd/edit.go
    └─ tea.NewProgram(edit.New(file))
           └─ edit.Model   (root Bubble Tea model)
                 ├─ ListModel    (left panel — key list)
                 ├─ PreviewModel (right panel — YAML textarea)
                 └─ *OverlayModel (floating overlay, nil when hidden)
```

The root model follows the standard Bubble Tea pattern: all state is carried in value types, `Update` returns a new copy of the model, and pointer receivers are only used for internal sub-model mutations that would otherwise require an extra return value.

**Panes** (`pane` enum):

| Constant | Meaning |
|----------|---------|
| `paneList` | Default — keyboard navigates the left list |
| `panePreview` | Tab was pressed — keyboard types into the YAML textarea |
| `paneOverlay` | An overlay is open — all input goes to `OverlayModel` |

---

## 2. File map

| File | Responsibility |
|------|---------------|
| `model.go` | Root Bubble Tea model; orchestrates all sub-models, routing, save/quit logic, layout |
| `list.go` | Left-panel list model; canonical key ordering, item rendering, message types |
| `preview.go` | Right-panel textarea; always visible, focus toggled by Tab |
| `overlay.go` | Floating overlay for adding/editing a block; single or two-panel layout |
| `yaml.go` | Pure YAML manipulation: parse blocks, insert, remove, validate |
| `field_defs.go` | Sub-field registry for guided two-panel overlay mode |
| `templates.go` | Guided YAML templates for every known top-level key |
| `styles.go` | All Lipgloss style constants |

---

## 3. Data flow

### Reading

```
os.ReadFile  →  rawYAML []byte
             →  ParseBlocksFromBytes  →  []Block  →  ListModel (Rebuild)
             →  PreviewModel.SetContent(string(rawYAML))
```

### Editing via list actions (add / edit / remove)

```
Key press  →  ListModel emits SpaceOnItemMsg / DeleteItemMsg
           →  Model.handleSpace / handleDelete
           →  InsertBlock / RemoveBlock  →  new rawYAML
           →  Model.applyRaw(newRaw)
                 ├─ m.rawYAML = newRaw
                 ├─ ParseBlocksFromBytes  →  m.blocks
                 ├─ ListModel.Rebuild(m.blocks)
                 ├─ PreviewModel.SetContent(string(newRaw))
                 └─ m.dirty = true
```

### Editing YAML directly (Tab → textarea)

```
Tab        →  panePreview; PreviewModel.Focus()
Typing     →  PreviewModel.Update(msg)
           →  Model.updatePreviewEditor:
                 raw = []byte(preview.Value())
                 if ParseBlocksFromBytes(raw) succeeds:
                     m.rawYAML = raw
                     m.blocks  = blocks
                     ListModel.Rebuild(blocks)   ← list syncs in real-time
                     m.dirty = true
Tab / Esc  →  paneList; PreviewModel.Blur()
```

Only syntactically valid YAML triggers a list rebuild — the list does not flicker during half-typed edits.

### Saving

```
ctrl+s  →  os.WriteFile(m.filePath, m.rawYAML, 0o644)
        →  m.dirty = false
```

---

## 4. Key bindings reference

### List pane (default)

| Key | Item type | Action |
|-----|-----------|--------|
| `↑` / `k` | any | Move cursor up |
| `↓` / `j` | any | Move cursor down |
| `Space` | available | Open guided overlay (template pre-filled) |
| `Space` | existing | Open guided overlay (current content pre-filled) |
| `e` | available | Open free overlay (blank textarea) |
| `e` | existing | Open free overlay (current content pre-filled) |
| `d` | existing | Delete block immediately (reversible until save) |
| `Tab` | — | Switch focus to YAML textarea (right panel) |
| `ctrl+s` | — | Save to disk |
| `q` / `ctrl+c` | — | Quit (prompts if dirty) |

### Preview pane (YAML editor)

| Key | Action |
|-----|--------|
| any printable | Edit YAML; list syncs in real-time on valid YAML |
| `Tab` / `Esc` | Return focus to list |
| `ctrl+s` | Save to disk |

### Overlay — single panel (free / simple guided)

| Key | Action |
|-----|--------|
| type | Edit YAML in textarea |
| `ctrl+s` | Validate and confirm; returns `OverlayConfirmedMsg` |
| `Esc` | Cancel; returns `OverlayCancelledMsg` |

### Overlay — two panels (guided + complex block)

| Key | Panel | Action |
|-----|-------|--------|
| `↑` / `k` | fields | Move field cursor up |
| `↓` / `j` | fields | Move field cursor down |
| `Space` | fields | Toggle field; rebuilds YAML textarea |
| `Tab` | either | Switch between field list and YAML editor |
| typing | yaml | Edit YAML; field toggles sync in real-time |
| `ctrl+s` | either | Validate and confirm |
| `Esc` | either | Cancel |

---

## 5. YAML manipulation

All YAML operations live in `yaml.go` and work on raw `[]byte` to preserve formatting, comments, and key ordering.

### `Block`

```go
type Block struct {
    Key     string
    Line    int     // 1-based line of the key node
    EndLine int     // last line of this block (exclusive — next key starts here)
}
```

Blocks are parsed using `gopkg.in/yaml.v3`'s `yaml.Node` tree to get accurate line numbers, then converted to `[]Block` for range operations.

### `ParseBlocksFromBytes(raw []byte) ([]Block, error)`

Unmarshals the document, walks the root mapping node, and records `(key, line)` for each top-level entry. `EndLine` is filled by looking at the next block's `Line - 1`; the last block extends to the total line count.

### `RemoveBlock(raw []byte, blocks []Block, key string) ([]byte, error)`

Splits `raw` into lines, removes the slice `[block.Line-1 : block.EndLine]`, and rejoins. Pure string operation — no re-serialisation, no formatting loss.

### `InsertBlock(raw []byte, snippet string) ([]byte, error)`

1. Validates the snippet with `ValidateSnippet`.
2. Parses the snippet to discover its key.
3. Looks up the key's rank in `allKnownKeys` (defined in `list.go`).
4. Scans existing blocks for the first one with a higher rank.
5. Splices the snippet lines in before that block's `Line`. Falls back to `appendBlock` if the key is unknown or no later block exists.

This keeps the file in the canonical DevContainer key order without re-serialising the entire document.

### `BlockContent(raw []byte, blocks []Block, key string) (string, error)`

Extracts the raw lines for a block as a string. Used to pre-fill the overlay when editing an existing block.

### `ValidateSnippet(text string) error`

Calls `yaml.Unmarshal` into `interface{}`. Returns an error for any parse failure.

---

## 6. Left panel — list model

**File:** `list.go`

### Canonical key order

`allKnownKeys []string` lists all 36 recognised DevContainer top-level keys in display/insertion order. This slice drives two things:

1. **List display** — keys are shown in this order in the "Available to add" section.
2. **Insertion rank** — `InsertBlock` uses the index in this slice to find the correct insertion point.

If you add a new key, add it to `allKnownKeys` at the correct position.

### `ListItem`

```go
type ListItem struct {
    Key       string
    Existing  bool   // true if the key is present in the current YAML
    Separator bool   // visual divider row — not selectable, emits no messages
}
```

### `BuildListItems(existing []Block) []ListItem`

Produces the merged item list:

1. Existing blocks in their file order (marked `Existing: true`).
2. A separator row.
3. Available keys (those in `allKnownKeys` but absent from the file) in alphabetical order.

### Messages emitted by `ListModel`

| Type | When |
|------|------|
| `SpaceOnItemMsg{Item, Guided bool}` | `Space` (Guided=true) or `e` (Guided=false) on a selectable item |
| `DeleteItemMsg{Key string}` | `d` on an existing item |

### Scrolling

`ListModel` maintains `cursor` and `offset` (scroll position). `clampScroll()` keeps the cursor visible within `height` rows. `moveCursor` skips separator rows automatically.

---

## 7. Right panel — preview / editor

**File:** `preview.go`

`PreviewModel` wraps a `bubbles/textarea` that is **always rendered** in the right panel. It is blurred (read-only appearance) when the list pane is active and focused (editable) when `panePreview` is active.

### Key methods

| Method | Description |
|--------|-------------|
| `SetContent(yaml string)` | Replaces textarea value; normalises `\r\n` → `\n` to prevent blank-line artefacts on Windows |
| `Value() string` | Returns current textarea content (may differ from last `SetContent` if user has typed) |
| `Focus() tea.Cmd` | Activates the textarea for typing |
| `Blur()` | Deactivates the textarea |
| `ScrollToKey(key string)` | Moves the textarea cursor to the line of `key:` (best-effort; textarea does not expose `SetYOffset`) |
| `Resize(width, height int)` | Updates textarea dimensions — called from `Model.relayout()` |

### Width invariant

`PreviewModel.Resize` receives `previewW - 2` (not `previewW`) so the textarea content fits exactly inside the panel border without triggering word-wrap. See [Layout and sizing](#11-layout-and-sizing) for the full calculation.

---

## 8. Overlay

**File:** `overlay.go`

The overlay is a floating box rendered over the main layout when `m.active == paneOverlay`. It is stored as `*OverlayModel` (nil when hidden) on the root model.

### Construction — `NewOverlay(key, initialContent string, guided bool, totalW, totalH int)`

The overlay decides its layout mode at construction time:

- **Two-panel** (`twoPanel = true`): `guided == true` AND `FieldsForKey(key)` returns at least one definition. Shows a field-toggle list on the left and a YAML textarea on the right.
- **Single textarea**: all other cases. Shows only a YAML textarea, pre-filled with `initialContent`.

### Edit vs. Add

`OverlayModel` carries two fields that control whether a confirm triggers an insert or a replace:

```go
isEdit  bool   // true when editing an existing block
editKey string // key of the block being replaced
```

When `isEdit` is true, `Model.handleOverlayConfirmed` removes the old block first, re-parses `m.blocks`, then calls `InsertBlock` with the new snippet.

### Two-panel overlay — field sync

When the user edits the YAML textarea manually, `syncFieldsFromYAML()` is called after every keystroke. It unmarshals the current textarea value and updates `Checked` on each `fieldState` to reflect what sub-keys are actually present. This keeps the toggle list in sync with manual edits in real-time.

When the user toggles a field in the left panel, `rebuildYAML()` reconstructs the textarea value from scratch in canonical field definition order — discarding any manual edits in the textarea. This is intentional: the left panel is the source of truth when you use it.

### Sizing

The overlay sizes itself from the outside in:

```
boxW  = min(totalW - 4, 120)   clamped to [60, 120]
boxH  = min(totalH - 4, 36)    clamped to [16, 36]

contentW = boxW - 4            (border L+R=2, padding L+R=2 from overlayBorderStyle)
panelH   = boxH - 8            (border 2 + 6 fixed rows: title, sep, hint, panel borders)

Two-panel widths:
  panelSpace  = contentW - 4   (two inner panel borders: 2×2)
  fieldPanelW = panelSpace / 3
  yamlPanelW  = panelSpace - fieldPanelW
```

The outer box uses `overlayBorderStyle.Render(...)` **without** an explicit `.Width()` call. Adding `.Width()` to the outer box caused a layout overflow because `lipgloss.Width(n)` includes padding but excludes the border — the box rendered 2 chars wider than expected. Auto-sizing from content avoids this.

### Messages

| Type | Sent when |
|------|-----------|
| `OverlayConfirmedMsg{Snippet string}` | `ctrl+s` and YAML is valid |
| `OverlayCancelledMsg{}` | `Esc` |

---

## 9. Guided mode — field definitions

**File:** `field_defs.go`

`blockFields map[string][]FieldDef` maps top-level keys to their toggleable sub-fields. Only keys with a fixed, well-known schema are listed. Simple scalar keys (e.g. `name`, `image`, `remoteUser`) are absent and fall back to a plain textarea.

```go
type FieldDef struct {
    Key      string // sub-field name shown in the toggle list
    Desc     string // one-line description (displayed next to the toggle)
    YAML     string // 2-space indented snippet, trailing \n (e.g. "  dockerfile: Dockerfile\n")
    Required bool   // pre-checked when the overlay opens
}
```

`FieldsForKey(key string) []FieldDef` is the only public accessor.

### Currently registered complex keys

| Key | Fields |
|-----|--------|
| `build` | `dockerfile`\*, `context`\*, `args`, `target`, `cacheFrom`, `output`, `ssh` |
| `customizations` | `vscode`, `jetbrains`, `codespaces` |
| `watch` | `waitFor`\*, `restart` |
| `otherPortsAttributes` | `onAutoForward`\*, `label`, `protocol` |

(\* = Required — pre-checked)

---

## 10. Guided mode — templates

**File:** `templates.go`

`guidedTemplates map[string]string` provides a full YAML template for every key in `allKnownKeys`. Templates include commented-out optional fields so the user can see what is available and uncomment as needed.

`GuidedTemplate(key string) string` returns the template for the key, or `key + ":\n"` as a safe fallback.

Templates are used only in single-panel guided mode (i.e. keys that have no `FieldDef` entries in `blockFields`). Complex keys use the two-panel overlay instead and build their YAML from `FieldDef.YAML` snippets.

---

## 11. Layout and sizing

### Main layout

```
Terminal width  = m.width
Terminal height = m.height

listW    = m.width / 3
previewW = m.width - listW - 4      (4 = two panel border pairs: 2×2)
innerH   = m.height - statusBarLines - 2   (statusBarLines=2; 2 = panel top+bottom border)

Panel rendering (Lipgloss):
  leftPanel  = panelStyle.Width(listW - 2).Height(innerH).Render(...)
  rightPanel = panelStyle.Width(previewW - 2).Height(innerH).Render(...)
```

### Lipgloss `Width(n)` semantics

> **Important:** In Lipgloss, `Width(n)` sets the *content + padding* width. Border characters are added on top. A `RoundedBorder` adds 1 char on each side, so the total rendered width is `n + 2`.

Consequence: if you want a panel that occupies `W` terminal columns, use `.Width(W - 2)`. The `PreviewModel` textarea is initialised with `previewW - 2` for this reason.

### Status bar

`const statusBarLines = 2` reserves two rows at the bottom:

1. **Feedback line** — current action message + `[modified]` marker if dirty.
2. **Hint line** — context-sensitive keyboard shortcut reference.

---

## 12. Styles

**File:** `styles.go`

| Variable | Used for |
|----------|----------|
| `panelStyle` | Inactive panel border (grey rounded) |
| `activePanelStyle` | Focused panel border (blue/violet rounded) |
| `existingItemStyle` | Existing keys in the list (green) |
| `availableItemStyle` | Available-to-add keys (grey) |
| `selectedItemStyle` | Cursor row in the list (pink, bold) |
| `separatorStyle` | Divider row between existing and available sections |
| `statusStyle` | Both status bar lines and overlay hint text |
| `dirtyStyle` | `[modified]` marker and unsaved-changes prompt (orange) |
| `overlayBorderStyle` | Outer overlay box (double border, blue/violet, padding 0 1) |
| `overlayTitleStyle` | Overlay title text (pink, bold) |

---

## 13. How to extend

### Add a new known top-level key

1. Add the key to `allKnownKeys` in `list.go` at the correct canonical position.
2. Add a template entry in `guidedTemplates` in `templates.go`.
3. Optionally add a `[]FieldDef` entry in `blockFields` in `field_defs.go` if the key has a fixed sub-schema worth exposing in guided two-panel mode.

### Add sub-fields to an existing complex key

Add a `FieldDef` to the relevant slice in `blockFields`. Set `Required: true` if the field should be pre-checked when the overlay opens. The `YAML` string must be 2-space indented and end with `\n`.

### Change the canonical key order

Edit `allKnownKeys` in `list.go`. The order affects both the display order in the "Available to add" section and the insertion position when calling `InsertBlock`.

### Add a new overlay mode

`NewOverlay` decides between two-panel and single-panel at construction time based on `guided && len(FieldsForKey(key)) > 0`. To add a third mode (e.g. a form-based editor), extend `overlayPanel` and add a branch in `NewOverlay`, `Update`, and `View`.

### Modify the status bar

- To add a third line, increment `statusBarLines` and update `relayout()` accordingly.
- To make the hint line context-aware in new ways, extend the `switch m.active` block in `View()`.
