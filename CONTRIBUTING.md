# Contributing to devcontainerwizard

Thank you for your interest in contributing. This document covers how to report issues, propose changes, and submit pull requests.

---

## Reporting issues

Open a GitHub issue and include:

- **What you expected** to happen
- **What actually happened** (paste the full error output)
- **Steps to reproduce** — the minimal `config.yaml` and command that triggers the issue
- **Environment:** OS, architecture, and version (`devcontainerwizard --version`)

For a bug in the TUI editor, run it with `--dump` and attach the trace file whose path is printed on exit. It records every keystroke and action of the session — check it for anything sensitive first.

For security vulnerabilities, do **not** open a public issue. Email the maintainer directly.

---

## Proposing changes

For non-trivial changes, open an issue or start a discussion before writing code. This avoids investing time in an implementation that conflicts with the project direction.

Small fixes (typos, one-line corrections) can go straight to a pull request.

---

## Submitting a pull request

1. Fork the repository and create a branch from `main`.
2. Follow the [Development Guide](docs/DEVELOPMENT.md) for setup, workflow, and conventions.
3. Run the full check suite before pushing:
   ```bash
   make all
   ```
4. Open a pull request against `main`.

### PR checklist

- [ ] `make all` passes locally (deps, fmt, docs, lint, security, tests)
- [ ] New behaviour is covered by tests
- [ ] Generated docs under `docs/devcontainerwizard/` are committed if the schema changed — CI fails on a stale diff
- [ ] A new schema field has a metadata node, a `TopLevelKeys` entry, and a preset entry
- [ ] `docs/ARCHITECTURE.md` updated if a new architectural pattern or decision was introduced
- [ ] Commit messages follow the project style (see below)

### Commit message style

One short subject line in imperative mood, 50 characters or fewer. Add up to 3–5 terse bullet points in the body only when context is genuinely needed.

```
feat: add hostRequirements.gpu to the schema

- add GPUValue union with bool/string/object forms
- declare gpu metadata under HostRequirements
- cover all three YAML shapes in round-trip tests
```

No trailing period on the subject line.

---

## Code of conduct

Be direct and constructive. Critique code, not people. Assume good intent.
