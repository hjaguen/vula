package gnome

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExtensionRecipe struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var CuratedExtensions = []ExtensionRecipe{
	{
		UUID:        "blur-my-shell@aunetx",
		Name:        "Blur my Shell",
		Description: "Adds a sleek glassmorphism blur effect to GNOME panel, overview, and dash",
	},
	{
		UUID:        "just-perfection-desktop@just-perfection",
		Name:        "Just Perfection",
		Description: "Fine-tune GNOME Shell UI, remove visual clutter, and optimize desktop animations",
	},
	{
		UUID:        "tiling-assistant@ubuntu.com",
		Name:        "Tiling Assistant",
		Description: "Built-in Ubuntu 24.04 tiling manager with quarter-snapping, gaps, and border highlights",
	},
}

type ExtensionInfoResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

// InstallExtension downloads and installs a GNOME extension from extensions.gnome.org
func (m *Manager) InstallExtension(uuid string) error {
	// If it's already a system extension (like tiling-assistant), just enable it
	if uuid == "tiling-assistant@ubuntu.com" {
		_ = exec.Command("gnome-extensions", "enable", uuid).Run()
		return nil
	}

	home := os.Getenv("HOME")
	extDir := filepath.Join(home, ".local/share/gnome-shell/extensions", uuid)
	_ = os.MkdirAll(extDir, 0755)

	// 1. Query GNOME Extensions API for GNOME 46
	apiURL := fmt.Sprintf("https://extensions.gnome.org/extension-info/?uuid=%s&shell_version=46", uuid)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed querying extension API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("extension '%s' not available for GNOME 46 (status: %d)", uuid, resp.StatusCode)
	}

	var extInfo ExtensionInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&extInfo); err != nil {
		return fmt.Errorf("failed decoding extension metadata: %w", err)
	}

	if extInfo.DownloadURL == "" {
		return fmt.Errorf("no download URL found for extension '%s'", uuid)
	}

	downloadFullURL := "https://extensions.gnome.org" + extInfo.DownloadURL
	if strings.HasPrefix(extInfo.DownloadURL, "http") {
		downloadFullURL = extInfo.DownloadURL
	}

	// 2. Download ZIP
	zipResp, err := client.Get(downloadFullURL)
	if err != nil {
		return fmt.Errorf("failed downloading extension zip: %w", err)
	}
	defer zipResp.Body.Close()

	tempZip := filepath.Join(os.TempDir(), uuid+".zip")
	out, err := os.Create(tempZip)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, zipResp.Body)
	out.Close()
	if err != nil {
		return err
	}
	defer os.Remove(tempZip)

	// 3. Unzip into target directory
	if err := unzip(tempZip, extDir); err != nil {
		return fmt.Errorf("failed unzipping extension: %w", err)
	}

	// 4. Enable extension via CLI
	_ = exec.Command("gnome-extensions", "enable", uuid).Run()
	return nil
}

// InstallCuratedExtensions installs all recommended visual extensions
func (m *Manager) InstallCuratedExtensions() error {
	for _, ext := range CuratedExtensions {
		if err := m.InstallExtension(ext.UUID); err != nil {
			// Log but don't halt entire loop
			continue
		}
	}
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
