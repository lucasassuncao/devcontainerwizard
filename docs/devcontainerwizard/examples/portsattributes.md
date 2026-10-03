# portsAttributes Examples

## Preset: base

```yaml
portsAttributes:
  "3000":
    label: Web App
    onAutoForward: notify
    protocol: http
```

## Preset: privileged-port

```yaml
portsAttributes:
  "80":
    label: HTTP
    onAutoForward: notify
    protocol: http
    elevateIfNeeded: true
    requireLocalPort: true
```

## Preset: web-stack

```yaml
portsAttributes:
  "3000":
    label: App
    onAutoForward: openBrowser
    protocol: http
  "5432":
    label: PostgreSQL
    onAutoForward: silent
  "6379":
    label: Redis
    onAutoForward: silent
```

