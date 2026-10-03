# Troubleshooting

---

## `convert` refuses to write

### "At least one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' must be set."

The config never says where the container comes from. Add exactly one of:

```yaml
image: mcr.microsoft.com/devcontainers/base:ubuntu
```

```yaml
build:
  dockerfile: Dockerfile
  context: .
```

```yaml
dockerComposeFile: docker-compose.yml
service: app
```

### "Only one of the fields ... can be set at a time."

Two or more sources are present. This usually happens after switching approaches
and leaving the old block behind — check for a stray `image:` above a `build:`
block. The `full` template ships with `image` active and `build` commented out
for exactly this reason.

`dockerFile` counts as a source too. It is the deprecated spelling of
`build.dockerfile`; delete it and use the `build` block.

### "Field 'Service' is required."

`dockerComposeFile` is set but nothing names the container to attach to. Add the
service name as it appears in the Compose file:

```yaml
service: app
```

### "Field 'X' must be one of the following values: ..."

An enum got a value outside its set. The message lists the legal ones. Press
**h** on the field in the editor to see the same list plus what each value does.

---

## "a mount object has no readonly property"

The Dev Containers specification defines a mount object as exactly `type`,
`source`, and `target`, and declares it closed to anything else. `readonly` was
never part of it, so a `devcontainer.json` carrying it is flagged by the schema.

Write the mount as a Docker `--mount` string instead — `mounts` accepts both
forms, and the string is handed to Docker unchanged:

```yaml
mounts:
  - source=./configs,target=/etc/app/config,type=bind,readonly
```

The same applies to any other Docker mount option the object form has no
property for, such as `consistency` or `bind-propagation`.

Earlier versions accepted `readonly` in the object form and then dropped it,
which turned a read-only mount into a writable one with no warning. It is now
an error naming the replacement.

---

## VS Code does not offer to reopen in a container

The file is somewhere VS Code does not look. Only these are auto-detected:

- `.devcontainer/devcontainer.json`
- `.devcontainer.json`
- `.devcontainer/<name>/devcontainer.json`

`convert` warns on stderr when the output path is not one of them. Re-run
without `--output`, or move the file.

---

## `init` says the file already exists

By design — `init` will not clobber a config you may have edited. Pass `--force`
to overwrite, or `--output` to write elsewhere.

---

## The editor shows no hints

The hint panel is toggled with **h**. If it opens but is empty for a field, that
field has no metadata declared — worth an issue, since every schema field is
supposed to carry a description.

---

## Ports are forwarded but `portsAttributes` is ignored

The keys must be strings:

```yaml
portsAttributes:
  "3000":       # quoted
    label: Web
```

Unquoted, YAML parses `3000` as an integer and the entry never matches the
forwarded port.

---

## A `${localEnv:VAR}` reference comes through empty

`${localEnv:...}` is resolved by the IDE at attach time, not by
`devcontainerwizard` — the variable has to be set in the environment that
launched the IDE, which is not necessarily the shell you ran `convert` in. On
Linux and macOS, an IDE started from the desktop rather than a terminal often
has a much smaller environment.

There is no top-level `localEnv` block; use `remoteEnv` with a
`${localEnv:VAR}` value.

---

## A secret is never requested

`secrets` is a hint to the IDE, and support varies: Codespaces prompts, other
tools may not. Where it is not supported, pass the value through `remoteEnv` and
`${localEnv:...}` instead.

---

## `make lint` fails on a clean checkout, on Windows

Almost always CRLF. `.gitattributes` pins the tree to LF (`* text=auto eol=lf`)
precisely because gofmt treats a CRLF file as unformatted. If you have an old
checkout from before that file existed:

```bash
git rm --cached -r .
git reset --hard
```

---

## The editor renders as boxes or the wrong colours

The TUI needs a terminal with a Nerd-capable font and true colour. Try another
theme (`--list-themes`), and on Windows prefer Windows Terminal over the legacy
console host.

---

## Reporting an editor bug

Record a trace and attach it to the issue:

```bash
devcontainerwizard edit --dump
```

The path is printed on exit. It contains every keystroke and every editor action
of the session — check it for anything sensitive before uploading.
