# Converting to devcontainer.json

`convert` is the step that produces the file VS Code and the devcontainer CLI
actually read.

```bash
devcontainerwizard convert
```

By default it reads `config.yaml` and writes `.devcontainer/devcontainer.json`,
creating the directory if needed.

```bash
devcontainerwizard convert -c staging.yaml -o .devcontainer/devcontainer.json --force
```

---

## What happens

1. **Load** — `config.yaml` is parsed as YAML.
2. **Validate** — every rule below is checked. On failure nothing is written and
   the errors name the offending fields.
3. **Write** — the document is marshalled to indented JSON, with a `$schema`
   pointing at `https://containers.dev/implementors/json_schema` if you did not
   set one yourself.

`convert` never overwrites an existing output file unless you pass `--force`.

---

## Validation rules

- `name` is required.
- Exactly one of `image`, `build`, `dockerFile`, or `dockerComposeFile` must be
  set. Zero is an error; so is more than one.
- `service` is required when `dockerComposeFile` is set.
- Enum fields must hold a listed value: `userEnvProbe`, `waitFor`,
  `shutdownAction` (`none`, `stopContainer`, `stopCompose`), `onAutoForward`,
  `protocol`, and mount `type`.
- `build.dockerfile` is required whenever `build` is present. `build.context`
  is optional and defaults to the folder holding `devcontainer.json`.
- `secrets.*.documentationUrl` must be a URL when set.
- `hostRequirements.cpus` must be at least 1.
- `mounts[].type` and `mounts[].target` are required.

The editor enforces more than `convert` does, because it validates the YAML
directly instead of a decoded struct — it also checks the size format of
`hostRequirements.memory` and `storage`, rejects duplicate entries in
`forwardPorts` and the other set-like lists, and warns on `dockerFile`.

Errors are reported one per line, naming the field and what it expected:

```
Field 'Name' is required.
Only one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' can be set at a time.
Field 'UserEnvProbe' must be one of the following values: none loginShell loginInteractiveShell interactiveShell.
```

The editor runs the equivalent checks on every save, so a config that came out
of `devcontainerwizard edit` normally converts on the first try.

---

## Output paths

VS Code and the devcontainer CLI only auto-detect two locations:

- `.devcontainer/devcontainer.json`
- `.devcontainer.json`

Write anywhere else and `convert` still does it, but prints a warning to stderr
so a config that will never be picked up does not pass unnoticed.

Multi-root setups (`.devcontainer/<name>/devcontainer.json`) are detected by
VS Code too, and `convert` treats a path ending in
`/.devcontainer/devcontainer.json` as canonical.

---

## In CI

`convert` exits non-zero on any validation failure, which makes it a usable
check on its own:

```yaml
- name: Validate devcontainer config
  run: devcontainerwizard convert --output /tmp/devcontainer.json --force
```

To assert the committed `devcontainer.json` is current, regenerate and diff:

```bash
devcontainerwizard convert --force
git diff --exit-code .devcontainer/devcontainer.json
```
