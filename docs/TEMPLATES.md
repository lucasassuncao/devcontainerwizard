# Templates

Templates are complete, working `config.yaml` files shipped inside the binary.
They are the starting point for `init`, and the same set backs the
whole-document picker (**p** on the root list) in the editor.

```bash
devcontainerwizard init --list                 # see what is available
devcontainerwizard init --template golang      # write config.yaml
devcontainerwizard init -t rust -o api.yaml    # write somewhere else
```

`init` refuses to overwrite an existing file unless you pass `--force`.

---

## Available templates

| Name | What you get |
|------|--------------|
| [`image`](#image) | A prebuilt Docker image, and nothing you do not need |
| [`dockerfile`](#dockerfile) | A build from your own Dockerfile |
| [`dockercompose`](#dockercompose) | A multi-service Compose stack |
| [`golang`](#golang) | Go toolchain, gopls, golangci-lint |
| [`python`](#python) | Python 3.12, Pylance, black, isort, Jupyter |
| [`rust`](#rust) | Rust stable, rust-analyzer, clippy, rustfmt |
| [`devops`](#devops) | Terraform, kubectl, Helm, and the AWS/GCP/Azure CLIs |
| [`full`](#full) | Every option in the specification, annotated and mostly commented out |

---

## image

The smallest thing that works: a Microsoft base image, a workspace folder, a
non-root user, one forwarded port, and a couple of VS Code extensions.

Start here when your project needs a Linux box with a language runtime already
on it and nothing more.

## dockerfile

Builds from a `Dockerfile` in the repository, passing a build arg. Adds port
attributes, a bind mount for `./data`, and lifecycle commands.

Start here when the image needs system packages, a compiler toolchain, or
anything else you would install with `apt`.

## dockercompose

Attaches the IDE to the `app` service of a `docker-compose.yml` and forwards the
app and database ports. Adds the git, GitHub CLI, and docker-in-docker features.

Start here when the dev environment is more than one container — an app plus a
database, a queue, or a cache.

## golang

`mcr.microsoft.com/devcontainers/go`, plus the `go` and `git` features and a VS
Code setup wired for gopls: `golangci-lint` on save, `goimports` as formatter,
race-enabled tests, and coverage on save.

## python

`mcr.microsoft.com/devcontainers/python:3.12` with the `python` feature,
unbuffered output, and a `postCreateCommand` that installs `requirements.txt`
when one exists. Extensions cover Pylance, Jupyter, black, and isort.

## rust

`mcr.microsoft.com/devcontainers/rust` with the `rust` feature on stable,
`RUST_BACKTRACE` on, and a `postCreateCommand` adding clippy and rustfmt.
Extensions cover rust-analyzer, LLDB debugging, TOML, and dependency inlays.

## devops

An Ubuntu base carrying the whole toolbox: docker-outside-of-docker, kubectl,
Helm, Terraform with tflint, and the AWS, Azure, and Google Cloud CLIs, with
`KUBECONFIG` already pointed at the right place.

## full

Every property the specification defines, each with a comment explaining what it
does, almost all commented out. It is a reference you edit down rather than a
config you use as-is.

Reach for it when you want to see the shape of the whole schema in one file. For
looking a single field up, the [attribute pages](/devcontainerwizard/attributes/devcontainer.md)
or the editor's hint panel are quicker.

---

## Using a template from inside the editor

Press **p** on the root list to open the template picker. Applying one replaces
the whole document, so do it on a fresh or throwaway file — not on a config you
have already filled in.

See [Editor](/EDIT.md).
