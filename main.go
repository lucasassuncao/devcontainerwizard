// Package main ...
package main

import (
	"github.com/lucasassuncao/devcontainerwizard/cmd"
	"github.com/lucasassuncao/devcontainerwizard/internal/updater"
)

// version is set at build time via -ldflags "-X main.version=<tag>"
var version = "dev"

func main() {
	updater.CleanOldBinary()
	cmd.Execute(version)
}
