# customizations Examples

## Preset: base

```yaml
customizations:
  vscode:
    settings:
      editor.formatOnSave: true
      files.insertFinalNewline: true
    extensions:
      - editorconfig.editorconfig
```

## Preset: jetbrains-default

```yaml
customizations:
  jetbrains:
    plugins:
      - com.intellij.plugins.github
```

## Preset: vscode-go

```yaml
customizations:
  vscode:
    settings:
      editor.formatOnSave: true
      go.formatTool: goimports
      go.lintTool: golangci-lint
      go.useLanguageServer: true
    extensions:
      - golang.go
```

## Preset: vscode-node

```yaml
customizations:
  vscode:
    settings:
      editor.defaultFormatter: esbenp.prettier-vscode
      editor.formatOnSave: true
    extensions:
      - dbaeumer.vscode-eslint
      - esbenp.prettier-vscode
```

## Preset: vscode-python

```yaml
customizations:
  vscode:
    settings:
      python.formatting.provider: black
      python.linting.enabled: true
    extensions:
      - ms-python.python
      - ms-python.vscode-pylance
```

