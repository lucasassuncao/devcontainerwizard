package cmd

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/lucasassuncao/bezel/table"
	"github.com/lucasassuncao/devcontainerwizard/internal/updater"
	"github.com/spf13/cobra"
)

// DefaultRepo is set at build time via ldflags.
var DefaultRepo = ""

func selfUpdateCmd(currentVersion string) *cobra.Command {
	var (
		repo       string
		version    string
		list       bool
		prerelease bool
		limit      int
	)

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update devcontainerwizard to a GitHub release",
		Long: `Downloads a release of devcontainerwizard from GitHub and replaces the current binary.
The old binary is kept as devcontainerwizard.old until the next run.

With no flags, installs the latest stable release.
Use --list to see available versions, --version to pick a specific one,
and --prerelease to include rc/beta/alpha releases.`,
		Example: `  devcontainerwizard self-update
  devcontainerwizard self-update --list
  devcontainerwizard self-update --list --prerelease
  devcontainerwizard self-update --version v1.2.0
  devcontainerwizard self-update --repo lucasassuncao/devcontainerwizard`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if list {
				return runList(cmd.Context(), cmd.OutOrStdout(), repo, prerelease, limit, currentVersion)
			}
			return updater.SelfUpdate(cmd.Context(), updater.Options{
				Repo:              repo,
				CurrentVersion:    currentVersion,
				Version:           version,
				IncludePrerelease: prerelease,
				Out:               cmd.OutOrStdout(),
				ErrOut:            cmd.ErrOrStderr(),
			})
		},
	}

	cmd.Flags().StringVar(&repo, "repo", DefaultRepo, `GitHub repository in "owner/repo" format`)
	cmd.Flags().StringVar(&version, "version", "", "Install this specific release tag (e.g. v1.2.0) instead of the latest")
	cmd.Flags().BoolVar(&list, "list", false, "List available releases and exit")
	cmd.Flags().BoolVar(&prerelease, "prerelease", false, "Include prereleases (rc/beta/alpha) in --list, or as the latest target when no --version is given")
	cmd.Flags().IntVar(&limit, "limit", 20, "Maximum number of releases to show with --list (max 100)")

	return cmd
}

func runList(ctx context.Context, w io.Writer, repo string, includePrerelease bool, limit int, currentVersion string) error {
	if repo == "" {
		return fmt.Errorf("--repo is required (e.g. --repo lucasassuncao/devcontainerwizard)")
	}

	releases, err := updater.ListReleases(ctx, repo, "", includePrerelease, limit)
	if err != nil {
		return err
	}
	if len(releases) == 0 {
		fmt.Fprintf(w, "No releases found for %s.\n", repo)
		return nil
	}

	current := updater.NormalizeVersion(currentVersion)
	rows := make([][]string, len(releases))
	for i, r := range releases {
		installed := updater.NormalizeVersion(r.Tag) == current
		var tags []string
		if i == 0 && !r.Prerelease {
			tags = append(tags, "latest")
		}
		if r.Prerelease {
			tags = append(tags, "prerelease")
		}
		if installed {
			tags = append(tags, "installed")
		}
		label := ""
		if len(tags) > 0 {
			labelStyle := th.Dim
			if installed {
				labelStyle = th.Success
			}
			label = labelStyle.Render("(" + strings.Join(tags, ", ") + ")")
		}
		published := ""
		if !r.PublishedAt.IsZero() {
			published = r.PublishedAt.Format(time.DateOnly)
		}
		rows[i] = []string{th.Bold.Render(r.Tag), label, th.Dim.Render(published)}
	}
	widths := table.Fit(make([]table.Column, 3), rows, math.MaxInt)
	for _, r := range rows {
		fmt.Fprintln(w, "  "+table.Row(widths, r))
	}
	return nil
}
