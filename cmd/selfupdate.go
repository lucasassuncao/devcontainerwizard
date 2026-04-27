package cmd

import (
	"github.com/lucasassuncao/devcontainerwizard/internal/updater"
	"github.com/spf13/cobra"
)

// DefaultRepo is set at build time via ldflags.
var DefaultRepo = ""

func selfUpdateCmd(currentVersion string) *cobra.Command {
	var repo string

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update devcontainerwizard to the latest GitHub release",
		Long: `Downloads the latest devcontainerwizard release from GitHub and replaces the current binary.
The old binary is kept as devcontainerwizard.old until the next run.`,
		Example: `  devcontainerwizard self-update
  devcontainerwizard self-update --repo lucasassuncao/devcontainerwizard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updater.SelfUpdate(repo, "", currentVersion)
		},
	}

	cmd.Flags().StringVar(&repo, "repo", DefaultRepo, `GitHub repository in "owner/repo" format`)

	return cmd
}
