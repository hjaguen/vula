package apps

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type WebAppInfo struct {
	Name        string
	URL         string
	IconPath    string
	DesktopPath string
}

// AddWebApp creates a desktop launcher entry for a web URL with custom icon
func AddWebApp(name, targetURL, customIconURL string) (*WebAppInfo, error) {
	name = strings.TrimSpace(name)
	targetURL = strings.TrimSpace(targetURL)
	if name == "" || targetURL == "" {
		return nil, fmt.Errorf("web app name and URL cannot be empty")
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	slug := sanitizeSlug(name)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	iconsDir := filepath.Join(homeDir, ".local", "share", "icons")
	appsDir := filepath.Join(homeDir, ".local", "share", "applications")
	_ = os.MkdirAll(iconsDir, 0755)
	_ = os.MkdirAll(appsDir, 0755)

	iconPath := filepath.Join(iconsDir, fmt.Sprintf("vula-webapp-%s.png", slug))
	desktopPath := filepath.Join(appsDir, fmt.Sprintf("vula-webapp-%s.desktop", slug))

	// Fetch Icon
	iconFetchURL := customIconURL
	if iconFetchURL == "" {
		iconFetchURL = fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=256", parsedURL.Host)
	}

	if err := downloadIcon(iconFetchURL, iconPath); err != nil {
		// Fallback to web icon default
		_ = downloadIcon(fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=128", parsedURL.Host), iconPath)
	}

	// Detect Browser
	execCmd := detectBrowserExec(targetURL)

	// Build .desktop content
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Comment=Vula Web Application for %s
Exec=%s
Icon=%s
Terminal=false
Categories=Network;WebBrowser;VulaApp;
StartupWMClass=%s
Actions=NewWindow;
`, name, targetURL, execCmd, iconPath, parsedURL.Host)

	if err := os.WriteFile(desktopPath, []byte(desktopContent), 0755); err != nil {
		return nil, fmt.Errorf("failed to write desktop file: %w", err)
	}

	// Update desktop database
	_ = exec.Command("update-desktop-database", appsDir).Run()

	return &WebAppInfo{
		Name:        name,
		URL:         targetURL,
		IconPath:    iconPath,
		DesktopPath: desktopPath,
	}, nil
}

// RemoveWebApp deletes a generated web app launcher and icon
func RemoveWebApp(name string) error {
	slug := sanitizeSlug(name)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appsDir := filepath.Join(homeDir, ".local", "share", "applications")
	iconsDir := filepath.Join(homeDir, ".local", "share", "icons")

	desktopPath := filepath.Join(appsDir, fmt.Sprintf("vula-webapp-%s.desktop", slug))
	iconPath := filepath.Join(iconsDir, fmt.Sprintf("vula-webapp-%s.png", slug))

	_ = os.Remove(desktopPath)
	_ = os.Remove(iconPath)

	_ = exec.Command("update-desktop-database", appsDir).Run()
	return nil
}

// ListWebApps returns all Vula-managed web app launchers
func ListWebApps() ([]WebAppInfo, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appsDir := filepath.Join(homeDir, ".local", "share", "applications")
	files, err := os.ReadDir(appsDir)
	if err != nil {
		return nil, nil
	}

	var webApps []WebAppInfo
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "vula-webapp-") && strings.HasSuffix(f.Name(), ".desktop") {
			fullPath := filepath.Join(appsDir, f.Name())
			data, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}

			lines := strings.Split(string(data), "\n")
			var name, appURL, iconPath string
			for _, line := range lines {
				if strings.HasPrefix(line, "Name=") {
					name = strings.TrimPrefix(line, "Name=")
				} else if strings.HasPrefix(line, "Comment=Vula Web Application for ") {
					appURL = strings.TrimPrefix(line, "Comment=Vula Web Application for ")
				} else if strings.HasPrefix(line, "Icon=") {
					iconPath = strings.TrimPrefix(line, "Icon=")
				}
			}

			if name != "" {
				webApps = append(webApps, WebAppInfo{
					Name:        name,
					URL:         appURL,
					IconPath:    iconPath,
					DesktopPath: fullPath,
				})
			}
		}
	}

	return webApps, nil
}

func sanitizeSlug(input string) string {
	cleaned := strings.ToLower(strings.TrimSpace(input))
	cleaned = strings.ReplaceAll(cleaned, " ", "-")
	var sb strings.Builder
	for _, r := range cleaned {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func detectBrowserExec(targetURL string) string {
	browsers := []string{"brave-browser", "google-chrome", "chromium-browser", "chromium"}
	for _, b := range browsers {
		if _, err := exec.LookPath(b); err == nil {
			return fmt.Sprintf("%s --app=\"%s\"", b, targetURL)
		}
	}
	return fmt.Sprintf("gio open \"%s\"", targetURL)
}

func downloadIcon(rawURL, destPath string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
