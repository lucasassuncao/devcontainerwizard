package updater

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	cases := map[string]string{
		"v1.2.3":  "1.2.3",
		"V1.2.3":  "1.2.3",
		"1.2.3":   "1.2.3",
		"":        "",
		"version": "ersion", // only the leading character is stripped
	}
	for in, want := range cases {
		if got := NormalizeVersion(in); got != want {
			t.Errorf("NormalizeVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

// asset name for the platform the test is running on.
func nativeAsset() string {
	name := fmt.Sprintf("devcontainerwizard_1.2.3_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// foreignOS returns a GOOS value that is not the one under test.
func foreignOS() string {
	for _, o := range []string{"linux", "darwin", "windows"} {
		if o != runtime.GOOS {
			return o
		}
	}
	return "plan9"
}

func TestSelectAssetPicksNativeAsset(t *testing.T) {
	assets := []ghAsset{
		{Name: fmt.Sprintf("devcontainerwizard_1.2.3_%s_%s", foreignOS(), runtime.GOARCH)},
		{Name: nativeAsset()},
	}
	got := selectAsset(assets)
	if got == nil || got.Name != nativeAsset() {
		t.Fatalf("selectAsset = %v, want %s", got, nativeAsset())
	}
}

// The regression this guards: an asset matching only on architecture is a
// binary for another operating system, and installing it bricks the CLI.
func TestSelectAssetRejectsForeignOS(t *testing.T) {
	assets := []ghAsset{
		{Name: fmt.Sprintf("devcontainerwizard_1.2.3_%s_%s", foreignOS(), runtime.GOARCH)},
	}
	if got := selectAsset(assets); got != nil {
		t.Fatalf("selectAsset = %q, want nil for a foreign-OS-only release", got.Name)
	}
}

// "darwin" contains "win", so with a bare "win" alias a macOS asset counted as an
// OS match on Windows. When a real Windows build is present the .exe bonus hides
// that, but on a release with no Windows build it defeats the OS gate entirely
// and hands the user a macOS binary instead of a clean "no compatible binary".
func TestSelectAssetDarwinDoesNotMatchWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("collision only manifests on windows, running on %s", runtime.GOOS)
	}
	assets := []ghAsset{
		{Name: "devcontainerwizard_1.2.3_Darwin_" + runtime.GOARCH},
	}
	if got := selectAsset(assets); got != nil {
		t.Fatalf("selectAsset = %q, want nil: darwin must not satisfy the windows OS match", got.Name)
	}
}

func TestSelectAssetSkipsChecksums(t *testing.T) {
	assets := []ghAsset{
		{Name: nativeAsset() + ".sha256"},
		{Name: "checksums.txt"},
	}
	if got := selectAsset(assets); got != nil {
		t.Fatalf("selectAsset = %q, want nil when only checksums are present", got.Name)
	}
}

func TestSelfUpdateReportsMissingRepo(t *testing.T) {
	var out, errOut bytes.Buffer
	err := SelfUpdate(t.Context(), Options{Out: &out, ErrOut: &errOut})
	if err == nil {
		t.Fatal("SelfUpdate with no Repo returned nil error")
	}
	if out.Len() != 0 || errOut.Len() != 0 {
		t.Errorf("expected no output before validation, got out=%q errOut=%q", out.String(), errOut.String())
	}
}

func TestOptionsWritersDefaultToProcessStreams(t *testing.T) {
	out, errOut := Options{}.writers()
	if out != os.Stdout || errOut != os.Stderr {
		t.Errorf("nil writers = (%v, %v), want (os.Stdout, os.Stderr)", out, errOut)
	}

	var a, b bytes.Buffer
	out, errOut = Options{Out: &a, ErrOut: &b}.writers()
	if out != &a || errOut != &b {
		t.Error("configured writers were not returned unchanged")
	}
}

// cleanOldBinary is the seam SelfUpdate uses to route its cleanup warning
// through Options.ErrOut rather than straight to os.Stderr.
func TestCleanOldBinaryRemovesStaleBinary(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "devcontainerwizard")
	old := exe + ".old"
	if err := os.WriteFile(old, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}

	orig := osExecutable
	osExecutable = func() (string, error) { return exe, nil }
	t.Cleanup(func() { osExecutable = orig })

	var w bytes.Buffer
	cleanOldBinary(&w)

	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("stale binary still present: %v", err)
	}
	if w.Len() != 0 {
		t.Errorf("unexpected warning on the success path: %q", w.String())
	}
}

// An unopenable destination must be reported, not silently treated as success.
// This one fails in os.OpenFile, before the cleanup defer is registered.
func TestDownloadReportsUnopenableDestination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("binary-payload"))
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	blocked := filepath.Join(dir, "sub")
	if err := os.MkdirAll(blocked, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := download(t.Context(), srv.URL, blocked, "", 14); err == nil {
		t.Fatal("writing over a directory was reported as success")
	}
}

// Every failure path must take the partial file with it, so nothing half-written
// can be renamed over the running executable. This is the observable half of the
// deferred cleanup; the Close error it also propagates needs a filesystem that
// defers write errors to Close, which no portable test can force.
func TestDownloadRemovesPartialFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("short"))
	}))
	t.Cleanup(srv.Close)

	dest := filepath.Join(t.TempDir(), "dcw.new")
	err := download(t.Context(), srv.URL, dest, "", 9999)
	if err == nil {
		t.Fatal("truncated download was accepted")
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Errorf("partial download left behind: %v", statErr)
	}
}

func TestDownloadWritesCompleteFile(t *testing.T) {
	payload := []byte("binary-payload")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(payload)
	}))
	t.Cleanup(srv.Close)

	dest := filepath.Join(t.TempDir(), "dcw.new")
	if err := download(t.Context(), srv.URL, dest, "", int64(len(payload))); err != nil {
		t.Fatalf("download: %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Errorf("wrote %q, want %q", got, payload)
	}
}
