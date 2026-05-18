package devcontainer

import (
	"fmt"
	"reflect"

	mapstructure "github.com/go-viper/mapstructure/v2"
	koanf "github.com/knadh/koanf/v2"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
)

func Parse(k *koanf.Koanf) (model.DevContainer, error) {
	var dc model.DevContainer
	err := k.UnmarshalWithConf("", &dc, koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				commandValueDecodeHook,
				gpuValueDecodeHook,
				mountOrStringDecodeHook,
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
			),
			Result: &dc,
		},
	})
	if err != nil {
		return model.DevContainer{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return dc, nil
}

var commandValueType = reflect.TypeOf(model.CommandValue{})

// commandValueDecodeHook converts string, []any, or map values into CommandValue
// so koanf/mapstructure can bind YAML scalars, sequences, and mappings.
func commandValueDecodeHook(f, t reflect.Type, data any) (any, error) {
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
func gpuValueDecodeHook(f, t reflect.Type, data any) (any, error) {
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
func mountOrStringDecodeHook(f, t reflect.Type, data any) (any, error) {
	if t != mountOrStringType {
		return data, nil
	}
	switch v := data.(type) {
	case string:
		return model.MountString(v), nil
	case map[string]any:
		var m model.Mount
		dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &m})
		if err != nil {
			return nil, err
		}
		if err := dec.Decode(v); err != nil {
			return nil, fmt.Errorf("decoding mount: %w", err)
		}
		return model.MountObject(m), nil
	}
	return data, nil
}
