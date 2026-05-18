# PortAttributes

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | Human-readable label for the port. | string | No | - |
| onAutoForward | Behavior when the port is auto-forwarded (notify, openBrowser, ignore). | string | No | - |
| protocol | Network protocol (http/https) for the port. | string | No | - |
| elevateIfNeeded | Prompt for elevated privileges if the port requires it (e.g. ports below 1024). | boolean | No | - |
| requireLocalPort | Require the local port to match the remote port. Shows a modal if not available. | boolean | No | - |

