// Package updater implements the self-update mechanism for devcontainerwizard.
package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lucasassuncao/bezel/theme"
)

// th is the palette the rest of the CLI prints with, so progress and warnings
// read the same here.
var th = theme.Resolve(theme.ThemeTerminal, true)

// apiClient handles short metadata requests against the GitHub API.
// downloadClient covers the larger binary fetch and uses a longer ceiling.
// Both have timeouts so a slow or hung server cannot stall the CLI indefinitely.
var (
	apiClient      = &http.Client{Timeout: 30 * time.Second}
	downloadClient = &http.Client{Timeout: 5 * time.Minute}
)

// osExecutable is a variable so tests can override os.Executable.
var osExecutable = os.Executable

// Release is the public, presentation-friendly view of a GitHub release
// returned by ListReleases. It hides API-specific fields.
type Release struct {
	Tag         string
	Prerelease  bool
	PublishedAt time.Time
}

// Options configures a SelfUpdate run.
//
// Repo must be in "owner/repo" format, e.g. "lucasassuncao/devcontainerwizard".
// CurrentVersion is the running binary's version (e.g. "1.0.0" or "v1.0.0");
// the update is skipped when it matches the resolved release tag.
//
// Version selects the release to install: empty means "latest". IncludePrerelease
// only affects the empty-version path: when true, the most recent release wins
// even if it is a prerelease; otherwise the latest stable is used. When Version
// is non-empty, IncludePrerelease is ignored — the explicit tag is honored.
type Options struct {
	Repo              string
	Token             string
	CurrentVersion    string
	Version           string
	IncludePrerelease bool

	// Out receives progress, ErrOut receives warnings. Both may be nil, in
	// which case they fall back to os.Stdout and os.Stderr; the CLI passes
	// cobra's streams and tests pass buffers.
	Out    io.Writer
	ErrOut io.Writer
}

// writers returns the configured streams, substituting the process defaults for
// any left nil.
func (o Options) writers() (out, errOut io.Writer) {
	out, errOut = o.Out, o.ErrOut
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return out, errOut
}

// SelfUpdate downloads a release of devcontainerwizard from GitHub and replaces
// the current binary. The old binary is kept as <name>.old until the next run,
// when it is cleaned up automatically.
//
// ctx cancels the metadata and download requests, so an interrupted run stops
// before the binary is swapped rather than after.
//
// The run is configured by opts; see Options.
func SelfUpdate(ctx context.Context, opts Options) error {
	if opts.Repo == "" {
		return fmt.Errorf("--repo is required (e.g. --repo lucasassuncao/devcontainerwizard)")
	}
	out, errOut := opts.writers()

	// Clean up any leftover .old binary from a previous update.
	cleanOldBinary(errOut)

	rel, err := resolveRelease(ctx, out, opts.Repo, opts.Token, opts.Version, opts.IncludePrerelease)
	if err != nil {
		return err
	}

	// Normalise both versions to a bare "X.Y.Z" form before comparing.
	if NormalizeVersion(opts.CurrentVersion) == NormalizeVersion(rel.TagName) {
		fmt.Fprintf(out, "%s Already on %s.\n", th.Success.Render("✓"), rel.TagName)
		return nil
	}

	asset := selectAsset(rel.Assets)
	if asset == nil {
		return fmt.Errorf("no compatible binary found in release %s", rel.TagName)
	}

	fmt.Fprintf(out, "Found %s → %s (%.1f MB)\n", th.Success.Render(rel.TagName), asset.Name, float64(asset.Size)/1e6)

	exePath, err := osExecutable()
	if err != nil {
		return fmt.Errorf("cannot determine current executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("resolving executable path: %w", err)
	}

	fmt.Fprintln(out, "Downloading new binary...")
	tmpPath := exePath + ".new"
	if err := download(ctx, asset.BrowserDownloadURL, tmpPath, opts.Token, asset.Size); err != nil {
		return err
	}

	// On Windows we cannot delete the running binary, but we can rename it.
	oldPath := exePath + ".old"
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(errOut, "%s could not remove stale binary %s: %v\n", th.Warning.Render("warning:"), oldPath, err)
	}

	fmt.Fprintln(out, "Replacing binary...")
	if err := os.Rename(exePath, oldPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming current binary: %w", err)
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		// Rollback: try to restore the previous binary. Surface any failure so
		// the user knows the install is in a broken state and where to recover.
		if rbErr := os.Rename(oldPath, exePath); rbErr != nil {
			fmt.Fprintf(errOut,
				"%s install failed and rollback failed too — original binary is at %s, downloaded binary at %s. Restore manually. Rollback error: %v\n",
				th.Danger.Render("CRITICAL:"), oldPath, tmpPath, rbErr)
		} else {
			_ = os.Remove(tmpPath)
		}
		return fmt.Errorf("installing new binary: %w", err)
	}

	fmt.Fprintf(out, "%s  (old binary saved as %s.old)\n", th.Success.Render("✓ Installed "+rel.TagName), filepath.Base(exePath))
	return nil
}

// ListReleases returns up to `limit` recent releases for the repo, newest first.
// Drafts are always excluded. When includePrerelease is false, prereleases are
// also excluded. limit <= 0 defaults to 20; the GitHub per-page cap is 100.
func ListReleases(ctx context.Context, repo, token string, includePrerelease bool, limit int) ([]Release, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo is required (e.g. lucasassuncao/devcontainerwizard)")
	}
	if limit <= 0 {
		limit = 20
	}

	// Over-fetch a little so filtering prereleases still gives us a usable list.
	perPage := limit
	if !includePrerelease {
		perPage *= 2
	}
	if perPage > 100 {
		perPage = 100
	}

	raw, err := fetchReleases(ctx, repo, token, perPage)
	if err != nil {
		return nil, err
	}

	out := make([]Release, 0, len(raw))
	for i := range raw {
		if raw[i].Draft {
			continue
		}
		if raw[i].Prerelease && !includePrerelease {
			continue
		}
		out = append(out, Release{
			Tag:         raw[i].TagName,
			Prerelease:  raw[i].Prerelease,
			PublishedAt: raw[i].PublishedAt,
		})
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// resolveRelease maps the (version, includePrerelease) inputs to a concrete release.
// Empty version → latest stable (or most recent including prereleases when flag set).
// Non-empty version → exact tag, with a fallback that retries with a leading "v".
func resolveRelease(ctx context.Context, out io.Writer, repo, token, version string, includePrerelease bool) (*ghRelease, error) {
	if version != "" {
		fmt.Fprintf(out, "Fetching release %s from %s...\n", version, repo)
		rel, err := fetchReleaseByTag(ctx, repo, token, version)
		if err == nil {
			return rel, nil
		}
		// Be forgiving: users often pass "1.2.3" when the actual tag is "v1.2.3".
		if !strings.HasPrefix(version, "v") {
			if alt, altErr := fetchReleaseByTag(ctx, repo, token, "v"+version); altErr == nil {
				return alt, nil
			}
		}
		return nil, err
	}

	if includePrerelease {
		fmt.Fprintf(out, "Checking most recent release of %s (including prereleases)...\n", repo)
		all, err := fetchReleases(ctx, repo, token, 10)
		if err != nil {
			return nil, err
		}
		for i := range all {
			if !all[i].Draft {
				return &all[i], nil
			}
		}
		return nil, fmt.Errorf("no releases found for %s", repo)
	}

	fmt.Fprintf(out, "Checking latest release of %s...\n", repo)
	return fetchLatestRelease(ctx, repo, token)
}

// CleanOldBinary removes a <exe>.old file left by a previous self-update,
// reporting any failure on stderr. Call this from main() at startup.
func CleanOldBinary() { cleanOldBinary(os.Stderr) }

// cleanOldBinary is CleanOldBinary with the warning stream injected, so a
// SelfUpdate run routes it through Options.ErrOut like every other message.
func cleanOldBinary(w io.Writer) {
	exe, err := osExecutable()
	if err != nil {
		return
	}
	old := exe + ".old"
	if _, err := os.Stat(old); err == nil {
		if err := os.Remove(old); err != nil {
			fmt.Fprintf(w, "%s could not remove old binary %s: %v\n", th.Warning.Render("warning:"), old, err)
		}
	}
}

// NormalizeVersion strips a leading "v" or "V" so that "v1.2.3" and "1.2.3"
// compare equal. It is exported because the CLI needs the same rule to mark the
// installed release in `self-update --list`.
func NormalizeVersion(v string) string {
	if len(v) > 0 && (v[0] == 'v' || v[0] == 'V') {
		return v[1:]
	}
	return v
}

// ── GitHub API ────────────────────────────────────────────────────────────────

type ghRelease struct {
	TagName     string    `json:"tag_name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

func newGitHubRequest(ctx context.Context, method, rawURL, token string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "github.com/lucasassuncao/devcontainerwizard")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

func fetchLatestRelease(ctx context.Context, repo, token string) (*ghRelease, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	var rel ghRelease
	if err := getJSON(ctx, u, token, &rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

func fetchReleaseByTag(ctx context.Context, repo, token, tag string) (*ghRelease, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, url.PathEscape(tag))
	var rel ghRelease
	if err := getJSON(ctx, u, token, &rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

func fetchReleases(ctx context.Context, repo, token string, perPage int) ([]ghRelease, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=%d", repo, perPage)
	var rels []ghRelease
	if err := getJSON(ctx, u, token, &rels); err != nil {
		return nil, err
	}
	return rels, nil
}

func getJSON(ctx context.Context, rawURL, token string, out any) error {
	req, err := newGitHubRequest(ctx, http.MethodGet, rawURL, token)
	if err != nil {
		return err
	}
	resp, err := apiClient.Do(req)
	if err != nil {
		return fmt.Errorf("github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found (404): %s", rawURL)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api returned %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// Scoring weights used by selectAsset to rank release assets.
// Higher weight = stronger signal of compatibility with the current platform.
const (
	scoreOS   = 5 // OS name match - strongest signal
	scoreArch = 3 // architecture match - secondary signal
	scoreExt  = 2 // platform-specific extension match (.exe on Windows)
)

// osAliases maps runtime.GOOS values to substring patterns that release asset
// names commonly use to indicate the same OS.
//
// A bare "win" is deliberately absent: "darwin" contains it, so on Windows it
// made every macOS asset an OS match. With ties resolved by asset order — and
// GitHub returning them alphabetically — "..._Darwin_arm64" outranked
// "..._Windows_arm64" on a Windows/arm64 host.
var osAliases = map[string][]string{
	"windows": {"windows", "win64", "win32"},
	"linux":   {"linux"},
	"darwin":  {"darwin", "macos", "mac", "osx"},
}

// archAliases maps runtime.GOARCH values to substring patterns commonly used
// in release asset names.
var archAliases = map[string][]string{
	"amd64": {"amd64", "x86_64", "x64"},
	"arm64": {"arm64", "aarch64"},
	"386":   {"i386", "x86", "386"},
	"arm":   {"armv7", "armhf", "arm"},
}

// selectAsset picks the best asset for the current platform, or nil when the
// release has none.
//
// Naming our OS is a hard requirement, not merely the heaviest score: an asset
// that matches on architecture alone is a binary for a different operating
// system, and installing it leaves the user with a CLI that will not run. Among
// the assets that clear that bar, arch match outweighs extension, the highest
// score wins, and ties keep the first candidate. Checksums, signatures, and
// plain-text files are excluded before scoring.
func selectAsset(assets []ghAsset) *ghAsset {
	skip := []string{".sha256", ".sha512", ".sig", ".asc", "checksums", ".txt"}

	// An unlisted GOOS/GOARCH still has a conventional spelling in asset names,
	// so fall back to the raw value rather than matching nothing at all.
	osPatterns := osAliases[runtime.GOOS]
	if len(osPatterns) == 0 {
		osPatterns = []string{runtime.GOOS}
	}
	archPatterns := archAliases[runtime.GOARCH]
	if len(archPatterns) == 0 {
		archPatterns = []string{runtime.GOARCH}
	}

	var best *ghAsset
	bestScore := -1

	for i := range assets {
		lower := strings.ToLower(assets[i].Name)
		excluded := false
		for _, s := range skip {
			if strings.Contains(lower, s) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		score := 0
		matchedOS := false
		for _, w := range osPatterns {
			if strings.Contains(lower, w) {
				score += scoreOS
				matchedOS = true
				break
			}
		}
		if !matchedOS {
			continue
		}
		for _, a := range archPatterns {
			if strings.Contains(lower, a) {
				score += scoreArch
				break
			}
		}
		if runtime.GOOS == "windows" && filepath.Ext(lower) == ".exe" {
			score += scoreExt
		}

		if score > bestScore {
			bestScore = score
			best = &assets[i]
		}
	}

	return best
}

// maxDownloadOverhead caps how much we read beyond the asset's advertised size.
// A small slack absorbs minor mismatches; anything past it indicates a server
// lying about Content-Length or serving the wrong asset, so we abort.
const maxDownloadOverhead = 1 << 20 // 1 MiB

// download fetches rawURL into destPath. The request honors ctx, so an interrupt
// during a large download aborts it rather than waiting out the client timeout.
// expectedSize is the asset size reported by the release metadata; the response
// body is capped at expectedSize + maxDownloadOverhead to prevent a
// misconfigured or hostile server from filling the disk.
func download(ctx context.Context, rawURL, destPath, token string, expectedSize int64) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "github.com/lucasassuncao/devcontainerwizard")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := downloadClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755) //#nosec G302 G304 -- binary must be executable; path is controlled by the updater
	if err != nil {
		return fmt.Errorf("creating temp binary: %w", err)
	}
	// Close is where a file system is entitled to report a write it had only
	// buffered, and this file is about to be renamed over the running
	// executable - discarding that error installs a truncated binary. The named
	// result lets the failure propagate, and any failure at all takes the
	// half-written file with it, so no partial download survives to be renamed.
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("closing binary: %w", closeErr)
		}
		if err != nil {
			_ = os.Remove(destPath)
		}
	}()

	limit := expectedSize + maxDownloadOverhead
	if expectedSize <= 0 {
		limit = 256 << 20 // 256 MiB hard cap when no size is advertised
	}
	limited := io.LimitReader(resp.Body, limit+1) // +1 so we can detect overflow

	written, err := io.Copy(f, limited)
	if err != nil {
		return fmt.Errorf("writing binary: %w", err)
	}
	if written > limit {
		return fmt.Errorf("download exceeded expected size (%d bytes, cap %d)", written, limit)
	}
	if expectedSize > 0 && written < expectedSize {
		return fmt.Errorf("download truncated: got %d bytes, expected %d", written, expectedSize)
	}
	if syncErr := f.Sync(); syncErr != nil {
		return fmt.Errorf("syncing binary: %w", syncErr)
	}
	return nil
}
