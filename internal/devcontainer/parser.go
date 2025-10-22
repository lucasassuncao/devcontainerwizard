package devcontainer

import (
	"fmt"

	"github.com/lucasassuncao/devcontainerwizard/internal/model"

	koanf "github.com/knadh/koanf/v2"
)

func Parse(k *koanf.Koanf) (model.DevContainer, error) {
	var dc model.DevContainer
	if err := k.Unmarshal("", &dc); err != nil {
		return model.DevContainer{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return dc, nil
}
