// Package updater implements the self-update mechanism for devcontainerwizard.
package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// apiClient handles short metadata requests against the GitHub API.
// downloadClient covers the larger binary fetch and uses a longer ceiling.
// Both have timeouts so a slow or hung server cannot stall the CLI indefinitely.
var (
	apiClient      = &http.Client{Timeout: 30 * time.Second}
	downloadClient = &http.Client{Timeout: 5 * time.Minute}
)

// osExecutable is a variable so tests can override os.Executable.
var osExecutable = os.Executable

// SelfUpdate downloads the latest release of devcontainerwizard from GitHub and
// replaces the current binary. The old binary is kept as <name>.old until the
// next run, when it is cleaned up automatically.
//
// repo must be in "owner/repo" format, e.g. "lucasassuncao/devcontainerwizard".
// currentVersion is the running binary's version (e.g. "1.0.0" or "v1.0.0");
// the update is skipped when it matches the latest release tag.
func SelfUpdate(repo, token, currentVersion string) error {
	if repo == "" {
		return fmt.Errorf("--repo is required (e.g. --repo lucasassuncao/devcontainerwizard)")
	}

	// Clean up any leftover .old binary from a previous update.
	CleanOldBinary()

	fmt.Printf("Checking latest release of %s...\n", repo)

	rel, err := fetchLatestRelease(repo, token)
	if err != nil {
		return err
	}

	// Normalise both versions to a bare "X.Y.Z" form before comparing.
	if normalizeVersion(currentVersion) == normalizeVersion(rel.TagName) {
		fmt.Printf("Already up to date (%s).\n", rel.TagName)
		return nil
	}

	asset := selectAsset(rel.Assets)
	if asset == nil {
		return fmt.Errorf("no compatible binary found in release %s", rel.TagName)
	}

	fmt.Printf("Found %s → %s (%.1f MB)\n", rel.TagName, asset.Name, float64(asset.Size)/1e6)

	exePath, err := osExecutable()
	if err != nil {
		return fmt.Errorf("cannot determine current executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("resolving executable path: %w", err)
	}

	fmt.Printf("Downloading new binary...\n")
	tmpPath := exePath + ".new"
	if err := download(asset.BrowserDownloadURL, tmpPath, token, asset.Size); err != nil {
		return err
	}

	// On Windows we cannot delete the running binary, but we can rename it.
	oldPath := exePath + ".old"
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: could not remove stale binary %s: %v\n", oldPath, err)
	}

	fmt.Printf("Replacing binary...\n")
	if err := os.Rename(exePath, oldPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming current binary: %w", err)
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		// Rollback: try to restore the previous binary. Surface any failure so
		// the user knows the install is in a broken state and where to recover.
		if rbErr := os.Rename(oldPath, exePath); rbErr != nil {
			fmt.Fprintf(os.Stderr,
				"CRITICAL: install failed and rollback failed too — original binary is at %s, downloaded binary at %s. Restore manually. Rollback error: %v\n",
				oldPath, tmpPath, rbErr)
		} else {
			os.Remove(tmpPath)
		}
		return fmt.Errorf("installing new binary: %w", err)
	}

	fmt.Printf("✓ Updated to %s  (old binary saved as %s.old)\n", rel.TagName, filepath.Base(exePath))
	return nil
}

// CleanOldBinary removes a <exe>.old file left by a previous self-update.
// Call this from main() at startup.
func CleanOldBinary() {
	exe, err := osExecutable()
	if err != nil {
		return
	}
	old := exe + ".old"
	if _, err := os.Stat(old); err == nil {
		if err := os.Remove(old); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not remove old binary %s: %v\n", old, err)
		}
	}
}

// normalizeVersion strips a leading "v" so that "v1.2.3" and "1.2.3" compare equal.
func normalizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

// ── GitHub API ────────────────────────────────────────────────────────────────

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

func fetchLatestRelease(repo, token string) (*ghRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "github.com/lucasassuncao/devcontainerwizard")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found for %s", repo)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decoding release: %w", err)
	}
	return &rel, nil
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
var osAliases = map[string][]string{
	"windows": {"windows", "win64", "win32", "win"},
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

// selectAsset picks the best asset for the current platform.
// Each candidate is scored: OS name match outweighs arch match, which outweighs
// extension. The highest-scored asset wins; ties keep the first candidate.
// Checksums, signatures, and plain-text files are excluded before scoring.
func selectAsset(assets []ghAsset) *ghAsset {
	skip := []string{".sha256", ".sha512", ".sig", ".asc", "checksums", ".txt"}

	osPatterns := osAliases[runtime.GOOS]
	archPatterns := archAliases[runtime.GOARCH]

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
		for _, w := range osPatterns {
			if strings.Contains(lower, w) {
				score += scoreOS
				break
			}
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

// download fetches url into destPath. expectedSize is the asset size reported by
// the release metadata; the response body is capped at expectedSize +
// maxDownloadOverhead to prevent a misconfigured or hostile server from filling
// the disk.
func download(url, destPath, token string, expectedSize int64) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
	defer f.Close()

	limit := expectedSize + maxDownloadOverhead
	if expectedSize <= 0 {
		limit = 256 << 20 // 256 MiB hard cap when no size is advertised
	}
	limited := io.LimitReader(resp.Body, limit+1) // +1 so we can detect overflow

	written, err := io.Copy(f, limited)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("writing binary: %w", err)
	}
	if written > limit {
		os.Remove(destPath)
		return fmt.Errorf("download exceeded expected size (%d bytes, cap %d)", written, limit)
	}
	if expectedSize > 0 && written < expectedSize {
		os.Remove(destPath)
		return fmt.Errorf("download truncated: got %d bytes, expected %d", written, expectedSize)
	}
	if err := f.Sync(); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("syncing binary: %w", err)
	}
	return nil
}
