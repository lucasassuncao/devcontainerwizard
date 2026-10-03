package cmd

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/lucasassuncao/devcontainerwizard/internal/devcontainer"
	"github.com/lucasassuncao/yedit/spec"
	"github.com/lucasassuncao/yedit/validate"
)

// DevContainerValidators is the rule set enforced by the edit command at
// validate/save time.
//
// Per-field constraints (required, allowed values, ranges, counts,
// uniqueness) are declared once in the hint tree (internal/hints) and enforced
// by the FromMetadata family — hints are the single source of field metadata.
// Only cross-field rules, which cannot live in per-field metadata, are
// declared here explicitly.
var DevContainerValidators = []spec.Validator{
	// Enforce everything the metadata declares.
	validate.RequiredFromMetadata(),
	validate.OneOfFromMetadata(),
	validate.RangeFromMetadata(),
	validate.PatternFromMetadata(),
	validate.CountFromMetadata(),
	validate.UniqueFromMetadata(),
	validate.DeprecatedFromMetadata(),
	validate.FormatFromMetadata(),
	validate.LengthFromMetadata(),
	validate.NotOneOfFromMetadata(),

	// Exactly one of image / build / dockerComposeFile must be present.
	validate.MutuallyExclusive("image", "build", "dockerComposeFile"),
	// service and runServices are only valid when dockerComposeFile is set.
	validate.RequiredWith("service", "dockerComposeFile"),
	validate.RequiredWith("runServices", "dockerComposeFile"),
	// RequiredWith reads "key is only valid alongside parent", so it cannot
	// express the rules below, which run the other way: the presence of
	// dockerComposeFile constrains another field.
	composeRequires("service"),
	composeRequires("workspaceFolder"),
	shutdownActionMatchesVariant(),

	// Whatever convert would still refuse after all of the above.
	convertRules(),
}

// composeRequires enforces one key of the specification's composeContainer
// variant, whose required set is dockerComposeFile, service, and
// workspaceFolder: without them the IDE has no container to attach to or path to open.
func composeRequires(key string) spec.Validator {
	return spec.ValidatorFunc(func(in spec.ValidationInput) []spec.Violation {
		var doc map[string]any
		if err := yaml.Unmarshal(in.Raw, &doc); err != nil {
			return nil
		}
		if doc["dockerComposeFile"] == nil {
			return nil
		}
		if s, _ := doc[key].(string); s != "" {
			return nil
		}
		return []spec.Violation{{
			Path:    key,
			Message: `required when "dockerComposeFile" is set`,
		}}
	})
}

// convertRules runs the pipeline convert writes through, so the editor refuses
// what convert would: a value that does not decode, an empty list or map entry,
// a malformed port. The rules the metadata already enforces are not repeated.
func convertRules() spec.Validator {
	return spec.ValidatorFunc(func(in spec.ValidationInput) []spec.Violation {
		k, err := devcontainer.LoadYAML(bytes.NewReader(in.Raw))
		if err != nil {
			return nil // malformed YAML is reported by the editor itself
		}
		dc, err := devcontainer.Parse(k)
		if err != nil {
			return []spec.Violation{{Message: "convert cannot read this document: " + err.Error()}}
		}
		var errs validator.ValidationErrors
		if !errors.As(devcontainer.Validate(dc), &errs) {
			return nil
		}
		var out []spec.Violation
		for _, e := range errs {
			if v, ok := convertOnlyViolation(e); ok {
				out = append(out, v)
			}
		}
		return out
	})
}

// convertOnlyViolation translates the failures only convert checks; every
// other tag has a metadata or cross-field rule above that already reports it.
func convertOnlyViolation(e validator.FieldError) (spec.Violation, bool) {
	switch {
	case e.Tag() == "no_empty_entry":
		return spec.Violation{Path: convertPaths[e.Field()], Message: "has an empty entry; give it a value or remove it"}, true
	case e.Tag() == "invalid_port":
		return spec.Violation{Path: convertPaths[e.Field()], Message: `must hold TCP ports between 1 and 65535, or "host:port" strings`}, true
	}
	return spec.Violation{}, false
}

// convertPaths maps the struct fields convertOnlyViolation reports to their YAML keys.
var convertPaths = map[string]string{
	"Mounts":          "mounts",
	"PortsAttributes": "portsAttributes",
	"ForwardPorts":    "forwardPorts",
	"AppPort":         "appPort",
}

// shutdownActionMatchesVariant enforces that shutdownAction uses the value of
// the container variant in play. The specification declares the key twice with
// a different enum in each: none|stopContainer for an image or Dockerfile
// config, none|stopCompose for a Docker Compose one. A single OneOf cannot say
// that, so the union is declared in the metadata and the variant is narrowed
// here.
func shutdownActionMatchesVariant() spec.Validator {
	return spec.ValidatorFunc(func(in spec.ValidationInput) []spec.Violation {
		var doc struct {
			DockerComposeFile any    `yaml:"dockerComposeFile"`
			ShutdownAction    string `yaml:"shutdownAction"`
		}
		if err := yaml.Unmarshal(in.Raw, &doc); err != nil {
			return nil
		}
		if doc.ShutdownAction == "" || doc.ShutdownAction == "none" {
			return nil
		}
		compose := doc.DockerComposeFile != nil
		want, other := "stopContainer", "stopCompose"
		if compose {
			want, other = "stopCompose", "stopContainer"
		}
		if doc.ShutdownAction != other {
			return nil
		}
		variant := "an image or Dockerfile"
		if compose {
			variant = "a Docker Compose"
		}
		return []spec.Violation{{
			Path:    "shutdownAction",
			Message: fmt.Sprintf("%q is not valid for %s config; use %q or \"none\"", other, variant, want),
		}}
	})
}
