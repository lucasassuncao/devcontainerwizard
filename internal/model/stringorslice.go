package model

import (
	"encoding/json"
	"fmt"
)

// StringOrSlice holds a value that the devcontainer spec allows as either a single
// string or an array of strings (lifecycle commands such as onCreateCommand).
// JSON output: one element → string, multiple → array.
type StringOrSlice []string

func (s StringOrSlice) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}
	return json.Marshal([]string(s))
}

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringOrSlice{str}
		return nil
	}
	var slice []string
	if err := json.Unmarshal(data, &slice); err != nil {
		return fmt.Errorf("must be a string or array of strings: %w", err)
	}
	*s = StringOrSlice(slice)
	return nil
}
