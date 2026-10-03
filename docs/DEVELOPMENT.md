# Development

## Requirements

Go as pinned in `go.mod`. Nothing else: every tool the Makefile uses is pinned
in `scripts/make/tools.mk` and installed into `./.gobin` on first use, so a
target runs the same version everywhere rather than whatever is on PATH.

```bash
git clone https://github.com/lucasassuncao/devcontainerwizard.git
cd devcontainerwizard
make all
```

`make all` runs the whole pipeline: `deps`, `fmt`, `docs`, `lint`, `security`,
`test-coverage`, then the supply-chain scans `govulncheck`, `sbom` and `vuln`,
which go last because they are the only steps that need the network.

Targets live in `scripts/make/*.mk`, one file per area. `make clean-tools`
forces a reinstall after a version bump in `tools.mk`.

---

## Makefile targets

| Target | What it does |
|---|---|
| `make help` | List every target |
| `make build` | Build for the current platform (goreleaser snapshot) |
| `make build-all` | Build for every target platform |
| `make run` | `go run main.go` |
| `make install` | `go install` |
| `make fmt` | `go fmt ./...` |
| `make lint` | golangci-lint |
| `make security` | gosec, medium severity and above |
| `make test` | gotestsum, testdox format, race detector on, shuffled order |
| `make test-watch` | The same, rerunning on file changes |
| `make test-coverage` | Tests plus HTML and Cobertura reports in `coverage/` |
| `make fuzz` | Fuzz the config parse pipeline (`FUZZTIME=5m` to override) |
| `make docs` | Regenerate the reference docs and the package READMEs |
| `make generate-docs` | Just the reference docs, examples, and JSON Schema |
| `make show-examples` | Browse the presets in the terminal |
| `make govulncheck` | Advisories reachable from this code |
| `make sbom` | CycloneDX SBOM of the dependency tree, in `sbom.json` |
| `make vuln` | The SBOM, scanned with grype |
| `make deps` | `go mod download && go mod tidy` |
| `make tools` | Install every pinned tool into `./.gobin` |
| `make tag VERSION=v1.2.3` | Tag and push a release |
| `make clean` | Remove `bin/`, `dist/`, `coverage/` and `sbom.json` |
| `make clean-all` | The same, plus `./.gobin` and the Go caches |

---

## Working against a local yedit

`go.work` is git-ignored, so a workspace file pointing at a checkout of
[yedit](https://github.com/lucasassuncao/yedit) lets you develop both at once
without touching `go.mod`:

```
go 1.27.0

use (
	.
	../yedit
)
```

Remove it (or `GOWORK=off`) to build against the pinned release. Never commit a
`replace` directive for this — CI has no sibling checkout and the release build
would pick up whatever is on the developer's disk.

---

## Adding a schema field

1. Add the struct field in `internal/model/model.go`, with `json` and `yaml`
   tags and any `validate` tag `convert` should enforce.
2. Add a node for it in the matching `Metadata()` method in
   `internal/model/metadata.go`: at minimum a `Description`, plus `Default`,
   `OneOf`, `Required`, `Min`/`Max`, `Formats`, `Snippet`, and `PreChecked`
   where they apply.
3. If it is a top-level field, add it to `model.TopLevelKeys` and give it an
   entry in the `internal/presets` registry — both have tests that fail
   otherwise.
4. `make docs` and commit the regenerated pages.

Metadata keys are validated against the struct's yaml tags at startup, so a typo
fails the `internal/cmd` tests rather than becoming dead metadata.

Cross-field rules go in `internal/cmd/edit_validators.go` (for the editor) and
in `internal/devcontainer/validator.go` (for `convert`). Anything expressible on
a single field belongs in the metadata instead — the `FromMetadata` validators
already enforce it in both places.

A rule only `convert` can check (a decode failure, an empty entry, a malformed
port) reaches the editor through `convertRules`. A new struct-level tag in
`validator.go` needs a case in `convertOnlyViolation`, or `ctrl+l` will not see it;
`TestEditorRefusesWhatConvertRefuses` is where to prove it does.

---

## Adding a template

Drop a `.yaml` file into `internal/cmd/templates/`. The first line must be a
description header, which is what `init --list` prints:

```yaml
# description: Ruby - Optimized setup for Ruby on Rails
```

It is picked up by `init --template <name>` and by the editor's document picker
automatically — both read the same embedded filesystem. Add a section to
`docs/TEMPLATES.md`.

---

## Tests

```bash
make test           # race detector on, testdox output
make test-coverage  # + coverage/coverage.html and coverage.xml
```

The tests that guard the invariants above:

- `internal/model` — every key in `TopLevelKeys` has a described metadata node.
- `internal/cmd` — the metadata tree is accepted by `metadata.New` (this is what
  catches yaml-tag typos), and key fields carry the metadata the validators need.
- `internal/presets` — every top-level key has a preset entry.
- `internal/cmd` — `TestEditorRefusesWhatConvertRefuses`: a document `convert`
  would refuse is also refused by `ctrl+l` and on save.

---

## CI

`.github/workflows/ci.yml` has four jobs:

- **checks**, on Linux only: `go mod tidy`, `make fmt`, and `make docs` each
  followed by `git diff --exit-code`. These are repository-hygiene checks whose
  diffs would be pure line-ending noise on Windows.
- **test**, on Linux, Windows, and macOS: `make lint`, `make security`,
  `make test-coverage`. Lint and security run inside the matrix rather than once
  on Linux, because Go only compiles the files matching the target platform.
- **supply-chain**, on Linux: `make govulncheck`, then `make vuln`. It can go red
  with nothing in the repository changed, when an advisory is published.
- **pages**, after checks and never on a pull request: copies `docs/` into
  `docs/pages/` and deploys the site.

Every job caches `./.gobin`, keyed on `scripts/make/tools.mk`.

`.gitattributes` pins the whole tree to LF. gofmt treats a CRLF file as
unformatted, so without it a Windows checkout fails `make lint` on files that
are perfectly fine in the repository.

---

## Releasing

`make tag VERSION=v1.2.3` tags and pushes; the rest is in
[Releasing](/RELEASING.md).
