package devcontainer

import (
	"os"
	"regexp"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

var (
	// matches ${env:VAR_NAME} — reads from host OS environment
	envPattern = regexp.MustCompile(`\$\{env:([^}]+)\}`)
	// matches ${localEnv:KEY} — references a key defined in localEnv
	localEnvPattern = regexp.MustCompile(`\$\{localEnv:([^}]+)\}`)
)

// ExpandLocalEnv resolves the localEnv section and substitutes its values into
// containerEnv and remoteEnv.
//
// Two-pass expansion:
//  1. Each localEnv value may contain ${env:VAR} — expanded from the host OS.
//  2. Values in containerEnv and remoteEnv may contain ${localEnv:KEY} —
//     replaced with the resolved localEnv value. Unresolved references are left as-is.
//
// localEnv itself is never written to the JSON output (json:"-" on the field).
func ExpandLocalEnv(dc *model.DevContainer) {
	if len(dc.LocalEnv) == 0 {
		return
	}

	resolved := resolveLocalEnv(dc.LocalEnv)
	dc.ContainerEnv = substituteLocalEnv(dc.ContainerEnv, resolved)
	dc.RemoteEnv = substituteLocalEnv(dc.RemoteEnv, resolved)
}

// resolveLocalEnv expands ${env:VAR} patterns in each localEnv value.
func resolveLocalEnv(local map[string]string) map[string]string {
	resolved := make(map[string]string, len(local))
	for k, v := range local {
		resolved[k] = envPattern.ReplaceAllStringFunc(v, func(match string) string {
			varName := envPattern.FindStringSubmatch(match)[1]
			return os.Getenv(varName)
		})
	}
	return resolved
}

// substituteLocalEnv replaces ${localEnv:KEY} in each map value with the
// corresponding resolved value. Unrecognised keys are left unexpanded.
func substituteLocalEnv(m map[string]string, resolved map[string]string) map[string]string {
	if len(m) == 0 {
		return m
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = localEnvPattern.ReplaceAllStringFunc(v, func(match string) string {
			key := localEnvPattern.FindStringSubmatch(match)[1]
			if val, ok := resolved[key]; ok {
				return val
			}
			return match
		})
	}
	return result
}
