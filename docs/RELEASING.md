# Releasing

A release is a git tag. Everything else follows from it.

```bash
make tag VERSION=v0.2.0
```

That target refuses to run with a dirty tree (`git diff --exit-code`), creates a tag annotated with the release notes and pushes it. The notes are every commit since the previous tag, subject and body, grouped by type, written by [git-cliff](https://git-cliff.org) from [`cliff.toml`](https://github.com/lucasassuncao/devcontainerwizard/blob/main/cliff.toml); the target installs git-cliff into `./.gobin` like the other tools. The push matches `v*.*.*` in [`.github/workflows/release.yml`](https://github.com/lucasassuncao/devcontainerwizard/blob/main/.github/workflows/release.yml), which runs GoReleaser, which builds six binaries (linux, darwin and windows, each amd64 and arm64) and attaches them to a GitHub release, with the tag's notes as its description.

Run `make all` first. The tag workflow does not re-run the test suite; it trusts that the commit it points at was already green on CI.

## What the build injects

The binary knows two things that are not in the source, both set by `ldflags` in [`.goreleaser.yaml`](https://github.com/lucasassuncao/devcontainerwizard/blob/main/.goreleaser.yaml):

| Symbol | Becomes |
| --- | --- |
| `main.version` | the tag being built |
| `internal/cmd.DefaultRepo` | `lucasassuncao/devcontainerwizard`, the repo `self-update` checks |

A binary built any other way (`go build`, `go install`, `make run`) gets neither. `--version` reports `dev`, so a working-tree build never claims to be a release, and `self-update` has no repository to check until it is given one with `--repo`.

## The coupling worth knowing about

`self-update` finds its download by **matching the asset filename**. The `name_template` in `.goreleaser.yaml` produces names like `devcontainerwizard_0.2.0_Windows_x86_64`, and `selectAsset` in [`internal/updater/selfupdate.go`](https://github.com/lucasassuncao/devcontainerwizard/blob/main/internal/updater/selfupdate.go) scores them against its own `osAliases` and `archAliases` tables: `Windows` and `x86_64` are entries in those tables, not coincidences.

Change the name template without changing the tables and every already-installed devcontainerwizard stops finding its upgrade. It fails quietly, and only for users on the old binary, never on your machine, where you build from source.

Naming the OS is a hard requirement. An asset that matches on architecture alone is a binary for another operating system, and `selectAsset` skips it rather than install something that cannot run. Checksum and signature files are skipped too; `self-update` does not verify the download against them.

## The token

`release.yml` reads `secrets.GORELEASER_TOKEN`, and that secret has to exist on the repository or the job publishes nothing: an unset secret is an empty string, not a fallback. Publishing here does not strictly need a PAT, since the workflow declares `contents: write` and the `GITHUB_TOKEN` Actions mints per run covers it; one is required only when GoReleaser writes to a *different* repository, such as a Homebrew tap.
