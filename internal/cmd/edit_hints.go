package cmd

import (
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
	"github.com/lucasassuncao/yedit/metadata"
	"github.com/lucasassuncao/yedit/spec"
)

func buildDevContainerHints() (spec.MetadataSource, error) {
	return metadata.New(model.DevContainer{})
}
