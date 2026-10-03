package devcontainer

import (
	"fmt"
	"reflect"

	mapstructure "github.com/go-viper/mapstructure/v2"
	koanf "github.com/knadh/koanf/v2"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

// Parse decodes a loaded configuration into a DevContainer.
//
// The decode hooks exist for the specification's union types, none of which
// mapstructure can bind on its own: a lifecycle command that may be a string, a
// list or a named map; a gpu requirement that may be a bool, a string or an
// object; and a mount that may be a string or an object.
func Parse(k *koanf.Koanf) (model.DevContainer, error) {
	return parse(k, false)
}

// ParseStrict is Parse with unknown keys rejected instead of dropped.
//
// mapstructure discards a key it cannot match, so a config saying
// "remoteUsr: vscode" decodes without complaint and the setting silently does
// nothing. Parse stays lenient on purpose: the specification's customizations
// block is open-ended, so a real config may carry keys for a tool this model
// does not describe, and refusing those would reject working configs.
//
// The templates this project ships have no such excuse - every key in them is
// one this model claims to support - so the template tests decode through here.
func ParseStrict(k *koanf.Koanf) (model.DevContainer, error) {
	return parse(k, true)
}

func parse(k *koanf.Koanf, strict bool) (model.DevContainer, error) {
	var dc model.DevContainer
	err := k.UnmarshalWithConf("", &dc, koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				commandValueDecodeHook,
				gpuValueDecodeHook,
				mountOrStringDecodeHook,
				stringToSingletonSliceHook,
				intToStringHook,
			),
			ErrorUnused: strict,
			Result:      &dc,
		},
	})
	if err != nil {
		return model.DevContainer{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return dc, nil
}

// stringToSingletonSliceHook promotes a bare scalar to a one-element list, so
// "dockerComposeFile: docker-compose.yml" is accepted alongside the array form
// the specification also allows.
//
// It replaces mapstructure.StringToSliceHookFunc(","), which split on commas.
// A comma is legal inside these values, so that turned
// runArgs: "--env=A,B" into two arguments and a single Docker --mount string
// into one broken mount per comma. Promoting without splitting keeps the
// convenience and drops the corruption.
func stringToSingletonSliceHook(f, t reflect.Kind, data any) (any, error) {
	if f != reflect.String || t != reflect.Slice {
		return data, nil
	}
	s, ok := data.(string)
	if !ok {
		return data, nil
	}
	if s == "" {
		return []string{}, nil
	}
	return []string{s}, nil
}

// intToStringHook lets a whole number fill a string field. YAML reads
// "storage: 64" as an int, while the specification types the size as a string
// with an optional unit, so a bare byte count was refused.
func intToStringHook(f, t reflect.Kind, data any) (any, error) {
	if t != reflect.String {
		return data, nil
	}
	switch f {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprint(data), nil
	}
	return data, nil
}

var commandValueType = reflect.TypeOf(model.CommandValue{})

// commandValueDecodeHook converts string, []any, or map values into CommandValue
// so koanf/mapstructure can bind YAML scalars, sequences, and mappings.
func commandValueDecodeHook(_, t reflect.Type, data any) (any, error) {
	if t != commandValueType {
		return data, nil
	}
	switch v := data.(type) {
	case string:
		return model.CommandString(v), nil
	case []any:
		items := make([]string, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("expected string at index %d, got %T", i, item)
			}
			items[i] = s
		}
		return model.CommandSlice(items), nil
	case []string:
		return model.CommandSlice(v), nil
	case map[string]any:
		named := make(map[string][]string, len(v))
		for k, val := range v {
			switch s := val.(type) {
			case string:
				named[k] = []string{s}
			case []any:
				sl := make([]string, len(s))
				for i, item := range s {
					str, ok := item.(string)
					if !ok {
						return nil, fmt.Errorf("expected string at index %d of key %q, got %T", i, k, item)
					}
					sl[i] = str
				}
				named[k] = sl
			case []string:
				named[k] = s
			default:
				return nil, fmt.Errorf("expected string or []string value for key %q, got %T", k, val)
			}
		}
		return model.CommandMap(named), nil
	}
	return data, nil
}

var gpuValueType = reflect.TypeOf(model.GPUValue{})

// gpuValueDecodeHook converts bool, string, or map values into GPUValue.
func gpuValueDecodeHook(_, t reflect.Type, data any) (any, error) {
	if t != gpuValueType {
		return data, nil
	}
	switch v := data.(type) {
	case bool:
		return model.GPUBool(v), nil
	case string:
		return model.GPUValue{StringVal: v}, nil
	case map[string]any:
		var r model.GPURequirement
		dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &r})
		if err != nil {
			return nil, err
		}
		if err := dec.Decode(v); err != nil {
			return nil, fmt.Errorf("decoding gpu requirement: %w", err)
		}
		return model.GPURequire(r), nil
	}
	return data, nil
}

var mountOrStringType = reflect.TypeOf(model.MountOrString{})

// mountOrStringDecodeHook converts string or map values into MountOrString.
func mountOrStringDecodeHook(_, t reflect.Type, data any) (any, error) {
	if t != mountOrStringType {
		return data, nil
	}
	switch v := data.(type) {
	case string:
		return model.MountString(v), nil
	case map[string]any:
		var m model.Mount
		// ErrorUnused mirrors the specification, whose Mount definition is
		// "additionalProperties": false with exactly type, source and target.
		// Without it an unrecognised key is dropped in silence - which for
		// "readonly" meant a mount the user wrote as read-only was emitted
		// read-write.
		dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &m, ErrorUnused: true})
		if err != nil {
			return nil, err
		}
		if err := dec.Decode(v); err != nil {
			if _, hasReadonly := v["readonly"]; hasReadonly {
				return nil, fmt.Errorf("a mount object has no %q property in the devcontainer specification; "+
					"write the mount as a Docker --mount string instead, e.g. "+
					"%q", "readonly", "source=./configs,target=/etc/app/config,type=bind,readonly")
			}
			return nil, fmt.Errorf("decoding mount: %w", err)
		}
		return model.MountObject(m), nil
	}
	return data, nil
}
