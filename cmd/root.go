// Package cmd ...
package cmd

import (
	"fmt"

	"github.com/lucasassuncao/devcontainerwizard/cmd/docs"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcontainerwizard",
	Short: "Manage DevContainer configurations",
	Long:  "A CLI to create DevContainer configuration files",
	Run:   runConvert,
}

func Execute() {
	rootCmd.AddCommand(docs.GenerateCmd, docs.ShowCmd, initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
