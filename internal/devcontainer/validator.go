package devcontainer

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"

	"github.com/go-playground/validator/v10"
)

// Validate checks dc against the struct tags declared on the model and the
// cross-field rules in DevContainerStructLevelValidation. Pass the result to
// HumanizeValidationError to render it for a user.
func Validate(dc model.DevContainer) error {
	v := validator.New()
	v.RegisterStructValidation(DevContainerStructLevelValidation, model.DevContainer{})
	return v.Struct(dc)
}

// DevContainerStructLevelValidation enforces that exactly one of Image, Build
// or DockerComposeFile is set on a DevContainer.
func DevContainerStructLevelValidation(sl validator.StructLevel) {
	dc, ok := sl.Current().Interface().(model.DevContainer)
	if !ok {
		return
	}

	count := 0
	if dc.Image != "" {
		count++
	}
	if dc.Build != nil {
		count++
	}
	if dc.DockerFile != "" {
		count++
	}
	if len(dc.DockerComposeFile) > 0 {
		count++
	}

	if count == 0 {
		sl.ReportError(dc, "DevContainer", "", "one_required", "")
	}
	if count > 1 {
		sl.ReportError(dc, "DevContainer", "", "mutually_exclusive", "")
	}

	compose := len(dc.DockerComposeFile) > 0

	// The specification's composeContainer variant requires workspaceFolder:
	// a Compose config has no automatic source-code mount to fall back on, so
	// without it the IDE has no path to open.
	if compose && dc.WorkspaceFolder == "" {
		sl.ReportError(dc.WorkspaceFolder, "WorkspaceFolder", "", "required_with_compose", "")
	}

	// shutdownAction is declared twice in the specification, with a different
	// enum per variant: none|stopContainer for image/Dockerfile,
	// none|stopCompose for Compose. The struct tag accepts the union; this
	// narrows it to the variant actually in play.
	switch {
	case !compose && dc.ShutdownAction == "stopCompose":
		sl.ReportError(dc.ShutdownAction, "ShutdownAction", "", "compose_only", "")
	case compose && dc.ShutdownAction == "stopContainer":
		sl.ReportError(dc.ShutdownAction, "ShutdownAction", "", "non_compose_only", "")
	}

	// A YAML null inside the mounts list or the portsAttributes map decodes to
	// an empty entry that dive skips - it validates values, and there is no
	// value here. Left alone it marshals straight back out as JSON null, and
	// the specification permits none in either place, so the run would write a
	// devcontainer.json the IDE refuses to open.
	for i := range dc.Mounts {
		if dc.Mounts[i].Mount == nil && dc.Mounts[i].Str == "" {
			sl.ReportError(dc.Mounts, "Mounts", "", "no_empty_entry", "")
			break
		}
	}
	for _, attrs := range dc.PortsAttributes {
		if attrs == nil {
			sl.ReportError(dc.PortsAttributes, "PortsAttributes", "", "no_empty_entry", "")
			break
		}
	}

	// forwardPorts and appPort are []any because the specification allows either
	// an integer port or a "host:port" string. No struct tag can express that,
	// so without this a bool, a fraction, or a port above 65535 travelled all
	// the way into devcontainer.json unchallenged.
	validatePortList(sl, dc.ForwardPorts, "ForwardPorts")
	validatePortList(sl, dc.AppPort, "AppPort")
}

// validatePortList reports the field when any entry is neither a valid TCP port
// number nor a non-empty host:port string.
func validatePortList(sl validator.StructLevel, ports []any, field string) {
	for _, p := range ports {
		if !validPort(p) {
			sl.ReportError(ports, field, "", "invalid_port", "")
			return
		}
	}
}

// validPort accepts the two forms the specification allows: an integer TCP port,
// or a non-empty string such as "8080" or "db:5432". YAML decoding can hand us
// any of the numeric kinds, so all of them are checked against the same range.
func validPort(p any) bool {
	switch v := p.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case int:
		return inPortRange(float64(v))
	case int32:
		return inPortRange(float64(v))
	case int64:
		return inPortRange(float64(v))
	case uint:
		return inPortRange(float64(v))
	case uint32:
		return inPortRange(float64(v))
	case uint64:
		return inPortRange(float64(v))
	case float32:
		return inPortRange(float64(v))
	case float64:
		return inPortRange(v)
	}
	return false
}

// inPortRange reports whether n is a whole number in the TCP port range.
func inPortRange(n float64) bool {
	return n == math.Trunc(n) && n >= 1 && n <= 65535
}

// validationMessages maps a validator tag to a fmt template. Placeholders:
//
//	%[1]s → field name (empty for struct-level errors)
//	%[2]s → tag param (e.g. allowed values for oneof, threshold for gt/lt)
//
// Tags absent from the map fall through to a generic "failed validation" line.
var validationMessages = map[string]string{
	"required":              "Field '%[1]s' is required.",
	"no_empty_entry":        "Field '%[1]s' has an empty entry; give it a value or remove it.",
	"invalid_port":          "Field '%[1]s' must hold TCP port numbers between 1 and 65535, or \"host:port\" strings.",
	"one_required":          "At least one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' must be set.",
	"mutually_exclusive":    "Only one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' can be set at a time.",
	"required_with_compose": "Field 'WorkspaceFolder' is required when 'DockerComposeFile' is set.",
	"compose_only":          "Field 'ShutdownAction' value \"stopCompose\" is only valid for a Docker Compose config; use \"stopContainer\" or \"none\".",
	"non_compose_only":      "Field 'ShutdownAction' value \"stopContainer\" is only valid for an image or Dockerfile config; use \"stopCompose\" or \"none\".",
	"file":                  "Field '%[1]s' must point to a valid file path.",
	"dir":                   "Field '%[1]s' must point to a valid directory path.",
	"oneof":                 "Field '%[1]s' must be one of the following values: %[2]s.",
	"required_with":         "Field '%[1]s' is required when the field it depends on is set.",
	"url":                   "Field '%[1]s' must be a valid URL.",
	"min":                   "Field '%[1]s' must be at least %[2]s.",
	"max":                   "Field '%[1]s' must be at most %[2]s.",
	"gt":                    "Field '%[1]s' must be greater than %[2]s.",
	"lt":                    "Field '%[1]s' must be less than %[2]s.",
	"dive":                  "Field '%[1]s' contains invalid nested elements.",
	"keys":                  "Field '%[1]s' has invalid map keys.",
	"endkeys":               "Field '%[1]s' has invalid map values.",
	"omitempty":             "Field '%[1]s' is optional but invalid when provided.",
}

// HumanizeValidationError renders a validation failure as one plain-language
// line per violation, looked up in validationMessages. An error that is not a
// validator.ValidationErrors is returned as its own text, and a nil error as the
// empty string.
func HumanizeValidationError(err error) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			tmpl, ok := validationMessages[e.Tag()]
			if !ok {
				tmpl = "Field '%[1]s' failed validation '" + e.Tag() + "'."
			}
			// Struct-level messages name their field inline and take no
			// placeholders. Passing args to those anyway made fmt append
			// "%!(EXTRA string=..., string=)" to the line the user reads.
			if strings.ContainsRune(tmpl, '%') {
				fmt.Fprintf(&sb, tmpl+"\n", e.Field(), e.Param())
			} else {
				fmt.Fprintln(&sb, tmpl)
			}
		}
	} else {
		sb.WriteString(err.Error())
	}

	return strings.TrimSpace(sb.String())
}
