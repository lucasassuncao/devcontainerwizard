# hostRequirements Examples

## Preset: base

```yaml
hostRequirements:
  cpus: 2
  memory: 4gb
```

## Preset: gpu

```yaml
hostRequirements:
  cpus: 4
  memory: 16gb
  gpu:
    cores: 1
    memory: 4gb
```

## Preset: gpu-optional

```yaml
hostRequirements:
  cpus: 2
  memory: 8gb
  gpu: optional
```

## Preset: gpu-required

```yaml
hostRequirements:
  cpus: 4
  memory: 16gb
  gpu: true
```

## Preset: heavy

```yaml
hostRequirements:
  cpus: 8
  memory: 16gb
  storage: 64gb
```

