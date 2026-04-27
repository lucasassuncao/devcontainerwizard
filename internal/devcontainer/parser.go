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
				stringOrSliceDecodeHook,
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

// stringOrSliceDecodeHook converts string or []interface{} values to model.StringOrSlice
// so that koanf/mapstructure can bind both YAML scalars and sequences to the type.
func stringOrSliceDecodeHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(model.StringOrSlice{}) {
		return data, nil
	}
	switch v := data.(type) {
	case string:
		return model.StringOrSlice{v}, nil
	case []interface{}:
		result := make(model.StringOrSlice, len(v))
		for i, item := range v {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("expected string at index %d, got %T", i, item)
			}
			result[i] = str
		}
		return result, nil
	case []string:
		return model.StringOrSlice(v), nil
	}
	return data, nil
}
