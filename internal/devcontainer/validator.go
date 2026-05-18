package devcontainer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"

	"github.com/go-playground/validator/v10"
)

func Validate(dc model.DevContainer) error {
	v := validator.New()
	v.RegisterStructValidation(DevContainerStructLevelValidation, model.DevContainer{})
	return v.Struct(dc)
}

// DevContainerStructLevelValidation enforces that exactly one of Image, Build
// or DockerComposeFile is set on a DevContainer.
func DevContainerStructLevelValidation(sl validator.StructLevel) {
	dc := sl.Current().Interface().(model.DevContainer)

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
}

// validationMessages maps a validator tag to a fmt template. Placeholders:
//
//	%[1]s → field name (empty for struct-level errors)
//	%[2]s → tag param (e.g. allowed values for oneof, threshold for gt/lt)
//
// Tags absent from the map fall through to a generic "failed validation" line.
var validationMessages = map[string]string{
	"required":           "Field '%[1]s' is required.",
	"one_required":       "At least one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' must be set.",
	"mutually_exclusive": "Only one of the fields 'Image', 'Build', 'DockerFile' or 'DockerComposeFile' can be set at a time.",
	"file":               "Field '%[1]s' must point to a valid file path.",
	"dir":                "Field '%[1]s' must point to a valid directory path.",
	"oneof":              "Field '%[1]s' must be one of the following values: %[2]s.",
	"gt":                 "Field '%[1]s' must be greater than %[2]s.",
	"lt":                 "Field '%[1]s' must be less than %[2]s.",
	"dive":               "Field '%[1]s' contains invalid nested elements.",
	"keys":               "Field '%[1]s' has invalid map keys.",
	"endkeys":            "Field '%[1]s' has invalid map values.",
	"omitempty":          "Field '%[1]s' is optional but invalid when provided.",
}

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
			fmt.Fprintf(&sb, tmpl+"\n", e.Field(), e.Param())
		}
	} else {
		sb.WriteString(err.Error())
	}

	return strings.TrimSpace(sb.String())
}
