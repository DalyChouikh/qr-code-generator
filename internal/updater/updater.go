// Package updater provides self-update functionality for the CLI.
//
// It checks the GitHub Releases API for newer versions and can download
// and replace the running binary. Compatible with macOS, Linux, and Windows.
package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
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

const (
	repoOwner  = "DalyChouikh"
	repoName   = "qr-code-generator"
	binaryName = "qrgen"

	// GitHub API endpoint for the latest release.
	latestReleaseURL = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases/latest"
)

// githubRelease represents the relevant fields from the GitHub Releases API.
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

// githubAsset represents a downloadable file attached to a release.
type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateResult contains the outcome of an update check.
type UpdateResult struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
}

// CheckForUpdate queries the GitHub API to see if a newer version is available.
// The currentVersion should be the semver string without the "v" prefix (e.g., "1.0.3").
func CheckForUpdate(currentVersion string) (*UpdateResult, error) {
	release, err := fetchLatestRelease()
	if err != nil {
		return nil, err
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	return &UpdateResult{
		CurrentVersion:  current,
		LatestVersion:   latest,
		UpdateAvailable: latest != current && current != "dev",
	}, nil
}

// SelfUpdate downloads the latest release and replaces the current binary.
// It returns the new version string on success.
func SelfUpdate(currentVersion string) (string, error) {
	release, err := fetchLatestRelease()
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	if current == "dev" {
		return "", fmt.Errorf("cannot update a development build — install a released version first")
	}

	if latest == current {
		return current, fmt.Errorf("already up to date (v%s)", current)
	}

	// Find the correct asset for this OS/arch
	asset, err := findAsset(release, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return "", err
	}

	// Download the archive
	archiveData, err := downloadAsset(asset.BrowserDownloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download update: %w", err)
	}

	// Extract the binary from the archive
	binaryData, err := extractBinary(archiveData, asset.Name)
	if err != nil {
		return "", fmt.Errorf("failed to extract update: %w", err)
	}

	// Replace the running binary
	if err := replaceBinary(binaryData); err != nil {
		return "", fmt.Errorf("failed to install update: %w", err)
	}

	return latest, nil
}

// fetchLatestRelease queries the GitHub Releases API.
func fetchLatestRelease() (*githubRelease, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", latestReleaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", binaryName+"-updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("GitHub API rate limit exceeded — try again later")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	return &release, nil
}

// findAsset locates the correct archive for the given OS and architecture.
func findAsset(release *githubRelease, goos, goarch string) (*githubAsset, error) {
	// Build expected archive name suffix based on goreleaser naming convention:
	// qrgen_<version>_<os>_<arch>.tar.gz (or .zip for windows)
	var ext string
	if goos == "windows" {
		ext = ".zip"
	} else {
		ext = ".tar.gz"
	}

	target := fmt.Sprintf("_%s_%s%s", goos, goarch, ext)

	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, binaryName+"_") && strings.HasSuffix(asset.Name, target) {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no release found for %s/%s — you may need to update manually", goos, goarch)
}

// downloadAsset fetches the archive from the given URL.
func downloadAsset(url string) ([]byte, error) {
	client := &http.Client{Timeout: 120 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Limit to 50MB to prevent abuse
	data, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20))
	if err != nil {
		return nil, err
	}

	return data, nil
}

// extractBinary pulls the qrgen binary out of a .tar.gz or .zip archive.
func extractBinary(archiveData []byte, archiveName string) ([]byte, error) {
	if strings.HasSuffix(archiveName, ".zip") {
		return extractFromZip(archiveData)
	}
	return extractFromTarGz(archiveData)
}

// extractFromTarGz extracts the binary from a tar.gz archive.
func extractFromTarGz(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := filepath.Base(hdr.Name)
		if name == binaryName {
			return io.ReadAll(tr)
		}
	}

	return nil, fmt.Errorf("binary %q not found in archive", binaryName)
}

// extractFromZip extracts the binary from a zip archive.
func extractFromZip(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	targetNames := []string{binaryName, binaryName + ".exe"}

	for _, f := range zr.File {
		name := filepath.Base(f.Name)
		for _, target := range targetNames {
			if name == target {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				return io.ReadAll(rc)
			}
		}
	}

	return nil, fmt.Errorf("binary %q not found in archive", binaryName)
}

// replaceBinary atomically replaces the running binary with new content.
// It writes to a temp file first, then renames over the old binary.
// On Windows, where you can't overwrite a running exe, it renames the
// old binary out of the way first.
func replaceBinary(newBinary []byte) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine executable path: %w", err)
	}

	// Resolve symlinks to get the real path
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("cannot resolve executable path: %w", err)
	}

	dir := filepath.Dir(execPath)
	base := filepath.Base(execPath)

	// Write new binary to a temp file in the same directory (same filesystem for rename)
	tmpFile, err := os.CreateTemp(dir, base+".update-*")
	if err != nil {
		return fmt.Errorf("cannot create temp file (do you have write permission to %s?): %w", dir, err)
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on any error
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	if _, err = tmpFile.Write(newBinary); err != nil {
		tmpFile.Close()
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	// Make the temp file executable
	if err = os.Chmod(tmpPath, 0755); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		// Windows: can't overwrite a running exe, so rename old one first
		oldPath := execPath + ".old"
		os.Remove(oldPath) // ignore error if it doesn't exist
		if err = os.Rename(execPath, oldPath); err != nil {
			return fmt.Errorf("cannot move old binary: %w", err)
		}
		if err = os.Rename(tmpPath, execPath); err != nil {
			// Try to restore the old binary
			os.Rename(oldPath, execPath)
			return fmt.Errorf("cannot move new binary into place: %w", err)
		}
		// Clean up old binary (may fail if still running, that's OK)
		os.Remove(oldPath)
	} else {
		// Unix: atomic rename
		if err = os.Rename(tmpPath, execPath); err != nil {
			return fmt.Errorf("cannot replace binary: %w", err)
		}
	}

	return nil
}
