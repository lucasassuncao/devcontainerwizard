package cmd

import (
	"github.com/lucasassuncao/bezel/theme"
	"github.com/spf13/cobra"
)

// defaultConfigPath is the config file every command falls back to. It is
// declared once so changing it does not mean hunting through three commands.
const defaultConfigPath = "config.yaml"

// th is the terminal palette the editor draws with, so plain command output
// reads in the same colours.
var th = theme.Resolve(theme.Terminal, true)

// RootCmd builds the command tree. version is injected at build time and
// reported by --version and self-update.
func RootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:     "devcontainerwizard",
		Short:   "Manage DevContainer configurations",
		Long:    "A CLI to create DevContainer configuration files",
		Version: version,
		// main prints the error and sets the exit code, so cobra must not do
		// it a second time. Usage is suppressed because almost every failure
		// here is a runtime one - a missing file, an invalid config - and
		// answering that with thirty lines of flag help buries the message.
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringP("config", "c", defaultConfigPath, "Path to the config file")

	root.AddCommand(
		convertCmd(),
		editCmd(),
		generateDocsCmd(),
		initCmd(),
		selfUpdateCmd(version),
		showExamplesCmd(),
	)

	return root
}

// configPath returns the config file to read, from the root's persistent
// --config flag. A command executed detached from the root has no such flag,
// so the default stands in rather than an empty path.
func configPath(cmd *cobra.Command) string {
	if v, err := cmd.Flags().GetString("config"); err == nil && v != "" {
		return v
	}
	return defaultConfigPath
}
