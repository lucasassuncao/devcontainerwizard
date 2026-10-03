# build Examples

## Preset: base

```yaml
build:
  dockerfile: Dockerfile
  context: .
```

## Preset: multi-stage-dev

```yaml
build:
  dockerfile: Dockerfile
  context: .
  target: dev
```

## Preset: with-args

```yaml
build:
  dockerfile: Dockerfile
  context: .
  args:
    VARIANT: latest
```

## Preset: with-options

```yaml
build:
  dockerfile: Dockerfile
  context: .
  options:
    - --no-cache
```

