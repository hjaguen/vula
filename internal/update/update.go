package update

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const GitHubRepo = "hjaguen/vula"

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type ReleaseInfo struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	PublishedAt time.Time      `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

// FetchLatestRelease queries the GitHub API for the latest published release
func FetchLatestRelease(ctx context.Context) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", GitHubRepo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "vula-updater")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found for %s yet", GitHubRepo)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var rel ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("failed to decode release payload: %w", err)
	}
	return &rel, nil
}

// IsNewer compares two semver tags (e.g., "v0.6.0" vs "v1.0.0")
func IsNewer(latest, current string) bool {
	lParts := parseSemver(latest)
	cParts := parseSemver(current)

	for i := 0; i < 3; i++ {
		if lParts[i] > cParts[i] {
			return true
		}
		if lParts[i] < cParts[i] {
			return false
		}
	}
	return false
}

func parseSemver(ver string) [3]int {
	clean := strings.TrimPrefix(strings.TrimSpace(ver), "v")
	parts := strings.Split(clean, ".")
	var res [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		numStr := parts[i]
		if dashIdx := strings.Index(numStr, "-"); dashIdx != -1 {
			numStr = numStr[:dashIdx]
		}
		if val, err := strconv.Atoi(numStr); err == nil {
			res[i] = val
		}
	}
	return res
}

// SelfUpdate updates the Vula CLI binary in ~/.local/bin/vula
func SelfUpdate(ctx context.Context, currentVersion string) (string, error) {
	home := os.Getenv("HOME")
	targetBin := filepath.Join(home, ".local", "bin", "vula")

	// 1. Check GitHub Releases
	rel, err := FetchLatestRelease(ctx)
	if err == nil && rel != nil {
		if !IsNewer(rel.TagName, currentVersion) {
			return fmt.Sprintf("Vula is already up to date (%s)", currentVersion), nil
		}

		// Look for architecture asset
		arch := runtime.GOARCH
		targetAssetSubstr := fmt.Sprintf("linux_%s", arch)

		var matchingAsset *ReleaseAsset
		for _, asset := range rel.Assets {
			if strings.Contains(asset.Name, targetAssetSubstr) && (strings.HasSuffix(asset.Name, ".tar.gz") || strings.HasSuffix(asset.Name, ".tar.xz")) {
				matchingAsset = &asset
				break
			}
		}

		if matchingAsset != nil {
			err := downloadAndExtractBinary(ctx, matchingAsset.BrowserDownloadURL, targetBin)
			if err == nil {
				return fmt.Sprintf("Successfully upgraded Vula from %s to %s via GitHub Release", currentVersion, rel.TagName), nil
			}
		}
	}

	// 2. Fallback: Local git repository build if ~/repos/vula exists
	vulaRepoDir := filepath.Join(home, "repos", "vula")
	if fi, err := os.Stat(vulaRepoDir); err == nil && fi.IsDir() {
		// Pull latest main
		pullCmd := exec.CommandContext(ctx, "git", "-C", vulaRepoDir, "pull", "origin", "main")
		_ = pullCmd.Run()

		// Build
		buildCmd := exec.CommandContext(ctx, "go", "build", "-ldflags=-s -w", "-o", targetBin, "./cmd/vula")
		buildCmd.Dir = vulaRepoDir
		if out, err := buildCmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("local source build failed: %w (output: %s)", err, string(out))
		}
		return "Successfully updated Vula from source repository (~/repos/vula)", nil
	}

	return "", fmt.Errorf("no release asset found and no local git repo available to update")
}

func downloadAndExtractBinary(ctx context.Context, downloadURL, targetBin string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("asset download failed: %s", resp.Status)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	tmpBin := targetBin + ".tmp"

	found := false
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg && (header.Name == "vula" || filepath.Base(header.Name) == "vula") {
			_ = os.MkdirAll(filepath.Dir(targetBin), 0755)
			outFile, err := os.OpenFile(tmpBin, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("binary 'vula' not found inside archive")
	}

	return os.Rename(tmpBin, targetBin)
}
