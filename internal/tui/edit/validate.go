package edit

import "fmt"

// mutuallyExclusiveGroups lists sets of fields where at most one may be present.
// These map to the three ways to define the container in the DevContainer spec.
var mutuallyExclusiveGroups = [][]string{
	{"image", "build", "dockerComposeFile"},
}

// dockerComposeOnly lists fields that are only meaningful when dockerComposeFile
// is present. Having them without it is a configuration error.
var dockerComposeOnly = []string{"service", "runServices"}

// ValidateMutualExclusions checks the active blocks for mutual-exclusion
// violations and returns one human-readable message per violation.
func ValidateMutualExclusions(blocks []Block) []string {
	present := make(map[string]bool, len(blocks))
	for _, b := range blocks {
		present[b.Key] = true
	}

	var violations []string

	for _, group := range mutuallyExclusiveGroups {
		var found []string
		for _, k := range group {
			if present[k] {
				found = append(found, k)
			}
		}
		if len(found) > 1 {
			violations = append(violations, fmt.Sprintf(
				"mutually exclusive — use only one of: %s",
				joinQuoted(found),
			))
		}
	}

	if !present["dockerComposeFile"] {
		for _, k := range dockerComposeOnly {
			if present[k] {
				violations = append(violations, fmt.Sprintf(
					"%q requires \"dockerComposeFile\" to be set", k,
				))
			}
		}
	}

	return violations
}

func joinQuoted(ss []string) string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = `"` + s + `"`
	}
	result := ""
	for i, s := range out {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
