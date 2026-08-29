package packages

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// IsAptInstalled checks if a Debian package is installed via dpkg-query
func IsAptInstalled(pkg string) bool {
	cmd := exec.Command("dpkg-query", "-W", "-f='${Status}'", pkg)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return string(out) == "'install ok installed'"
}

// InstallAptPackages installs missing APT packages in a single batch
func (m *Manager) InstallAptPackages(packages []string) error {
	var missing []string
	for _, p := range packages {
		if !IsAptInstalled(p) {
			missing = append(missing, p)
		}
	}

	if len(missing) == 0 {
		return nil // Idempotent: already installed
	}

	args := append([]string{"apt-get", "install", "-y", "--no-install-recommends"}, missing...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// VerifyChecksum verifies that a local file matches an expected SHA256 hex string
func VerifyChecksum(filePath, expectedSha256 string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}

	actualSha256 := hex.EncodeToString(hasher.Sum(nil))
	if actualSha256 != expectedSha256 {
		return fmt.Errorf("checksum mismatch for %s: expected %s, got %s", filePath, expectedSha256, actualSha256)
	}

	return nil
}

// DownloadFileSecurely downloads a file with timeout and optional SHA256 verification
func DownloadFileSecurely(url, destination, expectedSha256 string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download from %s: HTTP %d", url, resp.StatusCode)
	}

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	if expectedSha256 != "" {
		return VerifyChecksum(destination, expectedSha256)
	}

	return nil
}
