# mounts Examples

## Preset: base

```yaml
mounts:
  - type: bind
    source: ${localWorkspaceFolder}/.cache
    target: /home/vscode/.cache
```

## Preset: data-volume

```yaml
mounts:
  - type: volume
    source: app-data
    target: /var/lib/app
```

## Preset: string-form

```yaml
mounts:
  - source=${localWorkspaceFolder}/.cache,target=/home/vscode/.cache,type=bind,consistency=cached
```

