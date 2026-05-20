package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// setupInitCmd creates a fresh initCmd with output buffers wired for capture.
func setupInitCmd(t *testing.T, args []string) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	c := newInitCmd()
	c.SetOut(outBuf)
	c.SetErr(errBuf)
	c.SetArgs(args)
	return c, outBuf, errBuf
}

func TestInitNoTemplate(t *testing.T) {
	c, _, errOut := setupInitCmd(t, []string{})

	err := c.Execute()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	got := errOut.String()
	for _, want := range []string{"image", "dockerfile", "dockercompose", "full", "golang"} {
		if !strings.Contains(got, want) {
			t.Errorf("stderr missing %q\nfull output: %s", want, got)
		}
	}
}

func TestInitList(t *testing.T) {
	c, out, _ := setupInitCmd(t, []string{"--list"})

	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.String()
	for _, want := range []string{"image", "dockerfile", "dockercompose", "full", "golang"} {
		if !strings.Contains(got, want) {
			t.Errorf("stdout missing %q\nfull output: %s", want, got)
		}
	}
}

func TestInitTemplate(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	c, out, _ := setupInitCmd(t, []string{"--template", "golang"})

	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err != nil {
		t.Error("config.yaml was not created")
	}
	if !strings.Contains(out.String(), "Next:") {
		t.Errorf("expected next-step hint in output, got: %s", out.String())
	}
}

func TestInitInvalidTemplate(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	c, _, _ := setupInitCmd(t, []string{"--template", "nonexistent"})

	if err := c.Execute(); err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}
