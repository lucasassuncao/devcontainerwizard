# FAQ

## Why YAML? `devcontainer.json` already exists.

Because JSON has no comments and unforgiving quoting, and a devcontainer config
is exactly the kind of file you want to annotate. YAML also lets the wizard
validate and complete as you type, which JSON-with-comments tooling does
unevenly.

The output is still plain `devcontainer.json` — nothing downstream knows
devcontainerwizard was involved.

## Do I have to keep `config.yaml` in the repository?

No, but it is the more useful of the two files to review: it is shorter, it has
your comments in it, and a diff on it is readable. Committing both means anyone
who clones the repository can open it in a container without installing
anything.

## Can I edit `devcontainer.json` directly afterwards?

Yes, but the next `convert` overwrites it. Treat `config.yaml` as the source and
`devcontainer.json` as build output.

## Is there an import from an existing `devcontainer.json`?

Not yet. Start from the `full` template and copy your values across — it lists
every property with a comment explaining it.

## Does it support Dev Container Features?

Yes. `features` is passed through as-is, so anything in the
[catalog](https://containers.dev/features) works, and
`overrideFeatureInstallOrder` controls the order. See
[Configuration](/CONFIGURATION.md#features).

## Which spec version does it target?

The properties in the current
[Dev Containers specification](https://containers.dev/implementors/spec/).
Anything the model does not know about is rejected rather than passed through,
so the config cannot silently contain a typo. If a property is missing, that is
a bug worth reporting.

## Why does `dockerFile` warn?

It is the deprecated shorthand for `build.dockerfile`. It still converts, but
the editor hides it from the block list and marks it deprecated in the hint
panel. New configs should use the `build` block.

## Can I use it without the TUI?

Entirely. `init` writes a file, `convert` reads one, and both are plain
non-interactive commands. The editor is a convenience, not a required step.

## Does the editor need a special terminal?

It needs true colour and a reasonably modern terminal. Windows Terminal, iTerm2,
Alacritty, WezTerm, and the built-in terminals of VS Code and JetBrains IDEs all
work. The legacy Windows console host does not render it well.

## How do I get completion for `config.yaml` in my editor?

Reference the generated JSON Schema from the top of the file. See
[Getting Started](/GETTING-STARTED.md#editor-completion-for-configyaml).

## Where does `self-update` put the old binary?

Next to the new one, as `devcontainerwizard.old`. It is deleted on the following
run. Nothing else on the system is touched.

## Are the docs on this site hand-written?

Partly. The guides are; the
[attribute pages](/devcontainerwizard/attributes/devcontainer.md), the
[example pages](/devcontainerwizard/examples/devcontainer.md), and the JSON
Schema are generated from the same field metadata the editor's hint panel reads,
so they cannot drift from the tool's actual behaviour. CI fails if the committed
output is stale.

## How do I contribute?

See [CONTRIBUTING.md](https://github.com/lucasassuncao/devcontainerwizard/blob/main/CONTRIBUTING.md)
and [Development](/DEVELOPMENT.md).
