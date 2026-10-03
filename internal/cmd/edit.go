package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/lucasassuncao/bezel/theme"
	"github.com/lucasassuncao/bezel/themebrowser"
	"github.com/lucasassuncao/devcontainerwizard/internal/model"
	"github.com/lucasassuncao/yedit/editor"
)

// editCmd returns the "edit" command, which opens an interactive TUI editor
// for the devcontainerwizard configuration file.
func editCmd() *cobra.Command {
	var output string
	var themeName string
	var listThemes bool
	var noSaveConfirm bool
	var noDeleteConfirm bool
	var noValidateOnSave bool
	var dump bool
	var dumpPath string

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a devcontainer config file in an interactive TUI",
		Long: `Open a devcontainer configuration file in an interactive two-panel TUI editor.

The left panel lists top-level configuration keys; pressing Enter opens the
block editor where sub-fields can be toggled and edited. Ctrl+S writes the
file; Ctrl+E writes it and converts it to .devcontainer/devcontainer.json, as
convert --force would; Ctrl+U undoes the last change; Ctrl+Y redoes it; Esc quits.

Pressing p on the root list opens the whole-document template picker, backed
by the same templates the init command writes.

Use --output to write to a different file than the one loaded (e.g. to
produce a new config from an existing template).`,
		Example: `  # Edit the default configuration file
  devcontainerwizard edit

  # Edit with the Grape theme
  devcontainerwizard edit --theme grape

  # Browse available themes in an interactive, tabbed terminal UI
  devcontainerwizard edit --list-themes

  # Load from --config but save to a new file
  devcontainerwizard edit --output /path/to/new.yaml

  # Record a session trace to attach to a bug report
  devcontainerwizard edit --dump

  # Record the trace to a specific file instead of the OS temp dir
  devcontainerwizard edit --dump-path ./trace.jsonl`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if listThemes {
				return themebrowser.BrowseInTerminal()
			}

			selectedTheme, err := theme.Lookup(themeName)
			if err != nil {
				return fmt.Errorf("%w: run 'devcontainerwizard edit --list-themes' to see available themes", err)
			}

			devContainerHints, err := buildDevContainerHints()
			if err != nil {
				return fmt.Errorf("building hint source: %w", err)
			}

			loadPath := configPath(cmd)

			res, err := editor.Run(editor.Config{
				Path:              loadPath,
				SavePath:          output,
				Schema:            &model.DevContainer{},
				Title:             "devcontainerwizard",
				BlockPresets:      DevContainerBlockPresets,
				DocPresets:        DevContainerDocPresets,
				EnableHints:       true,
				Metadata:          devContainerHints,
				Validators:        DevContainerValidators,
				Hidden:            []string{"dockerFile"},
				Theme:             selectedTheme,
				NoSaveConfirm:     noSaveConfirm,
				NoDeleteConfirm:   noDeleteConfirm,
				NoValidateOnSave:  noValidateOnSave,
				AnimationDuration: 600 * time.Millisecond,
				Trace: editor.Trace{
					Dump:     dump || dumpPath != "",
					DumpPath: dumpPath,
				},
				Actions: []editor.Action{saveAndConvert},
			})
			if err != nil {
				return err
			}
			if res.Saved {
				savedTo := loadPath
				if output != "" {
					savedTo = output
				}
				fmt.Println("configuration saved to", savedTo)
			}
			if res.DumpPath != "" {
				fmt.Println("session trace written to", res.DumpPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Save to this file instead of the loaded config (load path is unchanged)")
	cmd.Flags().StringVar(&themeName, "theme", "plain", "Theme name (run --list-themes to see options)")
	cmd.Flags().BoolVar(&listThemes, "list-themes", false, "Browse available themes in an interactive, tabbed terminal UI")
	cmd.Flags().BoolVar(&noSaveConfirm, "no-save-confirm", false, "Skip the 'Save changes?' confirmation dialog")
	cmd.Flags().BoolVar(&noDeleteConfirm, "no-delete-confirm", false, "Skip the 'Remove block?' confirmation dialog")
	cmd.Flags().BoolVar(&noValidateOnSave, "no-validate-on-save", false, "Allow saving even when validators report errors (a warning is shown)")
	cmd.Flags().BoolVar(&dump, "dump", false, "Record every editor action to a JSONL trace file for bug reports (path is printed on exit)")
	cmd.Flags().StringVar(&dumpPath, "dump-path", "", "Write the session trace to this file instead of a temp file (implies --dump)")

	return cmd
}

// saveAndConvert is ctrl+e in the editor: the saved config goes through the
// convert pipeline to the canonical path. The JSON is generated, so it is
// overwritten, as convert --force would.
var saveAndConvert = editor.Action{
	Key:       "ctrl+e",
	Help:      "save & convert",
	SaveFirst: true,
	Run: func(ctx editor.ActionContext) (editor.ActionResult, error) {
		out, err := convertFile(ctx.Path, canonicalOutputPath, true)
		if err != nil {
			return editor.ActionResult{}, err
		}
		return editor.ActionResult{Message: fmt.Sprintf("Wrote %s from %s.", out, ctx.Path)}, nil
	},
}
