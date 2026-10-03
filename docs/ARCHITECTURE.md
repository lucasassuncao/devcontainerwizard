# Architecture

## The one idea

Every field of the devcontainer schema is described exactly once, in
`internal/model`, next to the Go struct that defines it. A `Metadata()` method
on each struct returns that field's description, default, allowed values,
constraints, editor snippet, and example.

Everything else is a projection of that:

```
internal/model  ── Metadata() ──┬─→ editor hint panel        (TUI)
      │                         ├─→ FromMetadata validators  (save / convert)
      │                         ├─→ attribute reference       (docs/…/attributes)
      │                         └─→ JSON Schema               (docs/…/schema)
      │
      └─ struct tags ──→ yaml decode, json encode, go-playground/validator
```

Adding a field means adding a struct field and a metadata node. The hint panel,
the validators, the reference table, and the schema all pick it up with no
further work — and `TestMetadataCoversTopLevelKeys` fails if the metadata node
is forgotten.

---

## Packages

| Package | Responsibility |
|---------|----------------|
| `internal/cmd` | Cobra command tree: `init`, `edit`, `convert`, `show-examples`, `self-update`, `generate-docs`. Also the editor's wiring — metadata source, validators, block and document presets. |
| `internal/model` | The schema. Structs, their yaml/json tags, the union types (`MountOrString`, `GPUValue`, `CommandValue`), and every field's `Metadata()`. |
| `internal/devcontainer` | The convert pipeline: load YAML, parse into the model, validate, marshal JSON. |
| `internal/presets` | Ready-made YAML values per top-level field, exposed as a `yedit/presets.Source`. |
| `internal/updater` | `self-update`: release listing, download, binary replacement. |

Each package carries a generated `README.md` (`make docs`, via gomarkdoc).

---

## The convert pipeline

`internal/devcontainer` is four small steps, one file each:

1. `loader.go` — read the YAML file into a koanf instance.
2. `parser.go` — decode into `model.DevContainer`, handling the union types.
3. `validator.go` — `go-playground/validator` over the struct tags, plus a
   struct-level rule for the mutually exclusive container sources.
   `HumanizeValidationError` turns tag failures into sentences.
4. `writer.go` — marshal to indented JSON, inject `$schema`, create the parent
   directory, refuse to overwrite without `--force`.

Only `convert` uses this path. The editor validates through yedit's rule set
instead, which sees the raw YAML rather than a decoded struct — that is what
lets it report an error against a line while you type.

---

## The editor

The TUI itself is [yedit](https://github.com/lucasassuncao/yedit), a generic
struct-schema-driven YAML editor. `devcontainerwizard` supplies:

| `editor.Config` field | Comes from |
|---|---|
| `Schema` | `&model.DevContainer{}` |
| `Metadata` | `metadata.New(model.DevContainer{})` — the tree above |
| `Validators` | `internal/cmd/edit_validators.go` |
| `BlockPresets` | `internal/presets` |
| `DocPresets` | the embedded `internal/cmd/templates/*.yaml` |
| `Hidden` | `dockerFile`, the deprecated source key |

Validators split in two. The `FromMetadata` family — required, one-of, range,
pattern, count, unique, deprecated, format, length, not-one-of — needs no
per-project code: it reads the same metadata tree the hint panel does. Only
cross-field rules that cannot live on a single field are written by hand:

```go
validate.MutuallyExclusive("image", "build", "dockerComposeFile")
validate.RequiredWith("service", "dockerComposeFile")
validate.RequiredWith("runServices", "dockerComposeFile")
```

---

## Union types

Three schema properties accept more than one shape, and each is a struct that
marshals to whichever form it holds:

- `MountOrString` — a `Mount` object or a Docker `--mount` string. It implements
  `yedit/schema.Provider`, declaring the field list reflection cannot read.
- `GPUValue` — `true`, `false`, `"optional"`, or a `GPURequirement` object.
- `CommandValue` — a string, a string list, or a map of named parallel commands.

Because reflection cannot see into them, their metadata is attached by hand:
`mountChildren()` hangs off the `mounts` node, and the other two stay leaves.

---

## Documentation generation

`generate-docs` (hidden; run it via `make docs`) calls
[docgen](https://github.com/lucasassuncao/docgen) with the same metadata and
presets the editor uses and writes, under `docs/devcontainerwizard/`:

- `attributes/*.md` — one reference page per struct, cross-linked
- `examples/devcontainer.md` — the whole-document templates as YAML
- `examples/*.md` — the block presets, one page per block that has its own
  attribute page, linked from it
- `schema/devcontainer.schema.json` — JSON Schema for editor completion
- `README.md` — the index over both

CI runs `make docs` and fails on any diff, so the committed output cannot go
stale.

---

## Command wiring

`main.go` installs a signal-cancelled context and hands it to
`cmd.RootCmd(version)`; commands return errors rather than printing and exiting,
and `main` prints once at the top. `version` and `cmd.DefaultRepo` are injected
by goreleaser via `-ldflags`.
