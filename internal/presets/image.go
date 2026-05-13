package presets

import "github.com/lucasassuncao/devcontainerwizard/internal/model"

func imagePresetsMap() map[string]string {
	return map[string]string{
		"base":   "ubuntu:22.04",
		"golang": "mcr.microsoft.com/devcontainers/go:latest",
		"node":   "mcr.microsoft.com/devcontainers/typescript-node:latest",
		"python": "mcr.microsoft.com/devcontainers/python:latest",
		"rust":   "mcr.microsoft.com/devcontainers/rust:latest",
		"java":   "mcr.microsoft.com/devcontainers/java:latest",
	}
}

func ImagePreset(name string) string { return imagePresetsMap()[name] }
func ListImagePresets() []string     { return sortedKeys(imagePresetsMap()) }

func buildPresetsMap() map[string]*model.BuildConfig {
	return map[string]*model.BuildConfig{
		"base": {
			Dockerfile: "Dockerfile",
			Context:    ".",
		},
		"with-args": {
			Dockerfile: "Dockerfile",
			Context:    ".",
			Args: map[string]string{
				"VARIANT": "latest",
			},
		},
		"multi-stage-dev": {
			Dockerfile: "Dockerfile",
			Context:    ".",
			Target:     "dev",
		},
	}
}

func BuildPreset(name string) *model.BuildConfig { return buildPresetsMap()[name] }
func ListBuildPresets() []string                 { return sortedKeys(buildPresetsMap()) }
