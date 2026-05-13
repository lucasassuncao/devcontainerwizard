package cmd

import (
	"fmt"

	"github.com/lucasassuncao/devcontainerwizard/cmd/docs"
	"github.com/lucasassuncao/devcontainerwizard/cmd/examples"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcontainerwizard",
	Short: "Manage DevContainer configurations",
	Long:  "A CLI to create DevContainer configuration files",
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.AddCommand(
		convertCmd,
		docs.GenerateCmd,
		docs.ShowCmd,
		examples.ShowCmd,
		initCmd,
		selfUpdateCmd(version),
		editCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
