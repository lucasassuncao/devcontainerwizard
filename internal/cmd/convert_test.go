package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasassuncao/yedit/editor"
	"github.com/spf13/cobra"
)

func setupConvertCmd(t *testing.T, args []string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	errBuf := new(bytes.Buffer)
	c := RootCmd("test")
	c.SetOut(new(bytes.Buffer))
	c.SetErr(errBuf)
	c.SetArgs(append([]string{"convert"}, args...))
	return c, errBuf
}

func writeMinimalConfig(t *testing.T, dir string) {
	t.Helper()
	body := "name: t\nimage: ubuntu:22.04\n"
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestConvertCanonicalNoWarning(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeMinimalConfig(t, dir)

	c, errOut := setupConvertCmd(t, []string{"-o", ".devcontainer/devcontainer.json"})
	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(errOut.String(), "Warning") {
		t.Errorf("expected no warning for canonical path, got: %s", errOut.String())
	}
	if _, err := os.Stat(filepath.Join(dir, ".devcontainer", "devcontainer.json")); err != nil {
		t.Error("canonical output file not created")
	}
}

func TestConvertNonCanonicalWarning(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeMinimalConfig(t, dir)

	c, errOut := setupConvertCmd(t, []string{"-o", "foo/devcontainer.json"})
	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(errOut.String(), "Warning") {
		t.Errorf("expected warning for non-canonical path, got: %s", errOut.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "foo", "devcontainer.json")); err != nil {
		t.Error("non-canonical output file not created")
	}
}

func TestConvertCreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeMinimalConfig(t, dir)

	c, _ := setupConvertCmd(t, []string{"-o", "deeply/nested/path/devcontainer.json"})
	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "deeply", "nested", "path", "devcontainer.json")); err != nil {
		t.Error("output file not created in nested dir")
	}
}

func TestConvertMissingConfig(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	c, _ := setupConvertCmd(t, []string{"-o", ".devcontainer/devcontainer.json"})
	if err := c.Execute(); err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

// ctrl+e in the editor writes the canonical file and keeps writing it: the JSON
// is generated, so a second run overwrites the first instead of failing.
func TestSaveAndConvertWritesTheCanonicalFile(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeMinimalConfig(t, dir)

	if !saveAndConvert.SaveFirst {
		t.Fatal("ctrl+e must save before it converts")
	}
	for range 2 {
		res, err := saveAndConvert.Run(editor.ActionContext{Path: "config.yaml"})
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if !strings.Contains(res.Message, canonicalOutputPath) {
			t.Errorf("message does not name the output: %q", res.Message)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".devcontainer", "devcontainer.json")); err != nil {
		t.Fatalf("canonical file not written: %v", err)
	}
}

func TestSaveAndConvertReportsAnInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("image: ubuntu\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := saveAndConvert.Run(editor.ActionContext{Path: "config.yaml"}); err == nil {
		t.Fatal("a config convert refuses must fail the action")
	}
}
