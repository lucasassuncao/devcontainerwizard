# otherPortsAttributes

## Examples

For usage examples, see [otherPortsAttributes presets](../examples/otherportsattributes.md).

## Arguments

The following arguments are supported:

| Name | Type | Description | Required | Default |
|------|------|-------------|----------|---------|
| label | string | Name shown next to the port in the IDE's forwarded-ports list. | No | Application |
| onAutoForward | string | What happens when a port is auto-forwarded: notify shows a toast, openBrowser opens the system browser, openPreview uses the IDE's built-in browser, silent forwards quietly, and ignore skips the port entirely. | No | notify |
| protocol | string | Scheme the IDE uses when it builds the URL for this port. | No | - |
| elevateIfNeeded | bool | Ask for elevated privileges when the local port needs them - ports below 1024 such as 22, 80 and 443. | No | false |
| requireLocalPort | bool | Insist on the same port number locally. When false the IDE silently maps to a free port instead; when true it warns that the port is taken. | No | false |

