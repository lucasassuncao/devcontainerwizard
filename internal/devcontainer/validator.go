// Package devcontainer ...
package devcontainer

import (
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

// DevContainerStructLevelValidation garante exclusividade entre Image, Build e DockerComposeFile.
func DevContainerStructLevelValidation(sl validator.StructLevel) {
	dc := sl.Current().Interface().(model.DevContainer)

	count := 0
	if dc.Image != "" {
		count++
	}
	if dc.Build != nil {
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

func HumanizeValidationError(err error) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			switch tag {
			case "required":
				sb.WriteString(fmt.Sprintf("Field '%s' is required.\n", field))
			case "one_required":
				sb.WriteString("At least one of the fields 'Image', 'Build' or 'DockerComposeFile' must be set.\n")
			case "mutually_exclusive":
				sb.WriteString("Only one of the fields 'Image', 'Build' or 'DockerComposeFile' can be set at a time.\n")
			case "file":
				sb.WriteString(fmt.Sprintf("Field '%s' must point to a valid file path.\n", field))
			case "dir":
				sb.WriteString(fmt.Sprintf("Field '%s' must point to a valid directory path.\n", field))
			case "oneof":
				sb.WriteString(fmt.Sprintf("Field '%s' must be one of the following values: %s.\n", field, param))
			case "gt":
				sb.WriteString(fmt.Sprintf("Field '%s' must be greater than %s.\n", field, param))
			case "lt":
				sb.WriteString(fmt.Sprintf("Field '%s' must be less than %s.\n", field, param))
			case "dive":
				sb.WriteString(fmt.Sprintf("Field '%s' contains invalid nested elements.\n", field))
			case "keys":
				sb.WriteString(fmt.Sprintf("Field '%s' has invalid map keys.\n", field))
			case "endkeys":
				sb.WriteString(fmt.Sprintf("Field '%s' has invalid map values.\n", field))
			case "omitempty":
				// omitempty is a soft rule; usually doesn't fail, but we can note it for completeness
				sb.WriteString(fmt.Sprintf("Field '%s' is optional but invalid when provided.\n", field))
			default:
				sb.WriteString(fmt.Sprintf("Field '%s' failed validation '%s'.\n", field, tag))
			}
		}
	} else {
		sb.WriteString(err.Error())
	}

	return strings.TrimSpace(sb.String())
}
