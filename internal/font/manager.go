package font

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

type FontRecipe struct {
	ID          string
	Name        string
	FontFamily  string
	ZipFilename string
	DownloadURL string
}

var AvailableFonts = []FontRecipe{
	{
		ID:          "cascadia",
		Name:        "Cascadia Code",
		FontFamily:  "CaskaydiaCoded Nerd Font",
		ZipFilename: "CascadiaCode.zip",
		DownloadURL: "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/CascadiaCode.zip",
	},
	{
		ID:          "firacode",
		Name:        "Fira Code",
		FontFamily:  "FiraCode Nerd Font",
		ZipFilename: "FiraCode.zip",
		DownloadURL: "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/FiraCode.zip",
	},
	{
		ID:          "jetbrains",
		Name:        "JetBrains Mono",
		FontFamily:  "JetBrainsMono Nerd Font",
		ZipFilename: "JetBrainsMono.zip",
		DownloadURL: "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/JetBrainsMono.zip",
	},
	{
		ID:          "meslo",
		Name:        "Meslo LG",
		FontFamily:  "MesloLGS Nerd Font",
		ZipFilename: "Meslo.zip",
		DownloadURL: "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/Meslo.zip",
	},
}

type Manager struct {
	homeDir string
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Manager{homeDir: home}, nil
}

// SetFont installs and configures target font family across terminals & VSCode
func (m *Manager) SetFont(fontID string) (*FontRecipe, error) {
	recipe := findRecipe(fontID)
	if recipe == nil {
		return nil, fmt.Errorf("unknown font '%s'. Available: cascadia, firacode, jetbrains, meslo", fontID)
	}

	// 1. Ensure font is installed locally
	fontsDir := filepath.Join(m.homeDir, ".local", "share", "fonts")
	_ = os.MkdirAll(fontsDir, 0755)

	fontTargetFile := filepath.Join(fontsDir, recipe.ZipFilename)
	if _, err := os.Stat(fontTargetFile); os.IsNotExist(err) {
		fmt.Printf("Downloading %s Nerd Font...\n", recipe.Name)
		if err := downloadAndExtractFont(recipe.DownloadURL, recipe.ZipFilename, fontsDir); err != nil {
			return nil, fmt.Errorf("failed downloading font: %w", err)
		}
		_ = exec.Command("fc-cache", "-f", fontsDir).Run()
	}

	// 2. Update Ghostty config
	_ = m.updateGhosttyFont(recipe.FontFamily)

	// 3. Update Alacritty config
	_ = m.updateAlacrittyFont(recipe.FontFamily)

	// 4. Update VS Code settings
	_ = m.updateVSCodeFont(recipe.FontFamily)

	// 5. Update GNOME Monospace interface font
	_ = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", fmt.Sprintf("'%s 10'", recipe.FontFamily)).Run()

	return recipe, nil
}

// SetFontSize updates font size across developer tools
func (m *Manager) SetFontSize(size int) error {
	if size <= 6 || size >= 48 {
		return fmt.Errorf("font size out of bounds (7-47)")
	}

	_ = m.updateGhosttyFontSize(size)
	_ = m.updateAlacrittyFontSize(size)
	_ = m.updateVSCodeFontSize(size)
	return nil
}

func findRecipe(query string) *FontRecipe {
	q := strings.ToLower(strings.TrimSpace(query))
	for _, f := range AvailableFonts {
		if f.ID == q || strings.Contains(strings.ToLower(f.Name), q) {
			return &f
		}
	}
	return nil
}

func downloadAndExtractFont(rawURL, zipName, targetDir string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d", resp.StatusCode)
	}

	tempZip := filepath.Join(os.TempDir(), zipName)
	out, err := os.Create(tempZip)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}
	defer os.Remove(tempZip)

	r, err := zip.OpenReader(tempZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".ttf") || strings.HasSuffix(strings.ToLower(f.Name), ".otf") {
			fpath := filepath.Join(targetDir, filepath.Base(f.Name))
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				continue
			}
			_, _ = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
		}
	}

	// Mark zip landmark marker
	_ = os.WriteFile(filepath.Join(targetDir, zipName), []byte("INSTALLED"), 0644)
	return nil
}

func (m *Manager) updateGhosttyFont(fontFamily string) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "ghostty", "config")
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)

	data, err := os.ReadFile(cfgPath)
	content := string(data)
	if err != nil {
		content = ""
	}

	lines := strings.Split(content, "\n")
	found := false
	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "font-family =") {
			newLines = append(newLines, fmt.Sprintf("font-family = \"%s\"", fontFamily))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}
	if !found {
		newLines = append(newLines, fmt.Sprintf("font-family = \"%s\"", fontFamily))
	}

	return os.WriteFile(cfgPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func (m *Manager) updateGhosttyFontSize(size int) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "ghostty", "config")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	found := false
	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "font-size =") {
			newLines = append(newLines, fmt.Sprintf("font-size = %d", size))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}
	if !found {
		newLines = append(newLines, fmt.Sprintf("font-size = %d", size))
	}

	return os.WriteFile(cfgPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func (m *Manager) updateAlacrittyFont(fontFamily string) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "alacritty", "alacritty.toml")
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)

	data, err := os.ReadFile(cfgPath)
	content := string(data)
	if err != nil {
		content = ""
	}

	lines := strings.Split(content, "\n")
	found := false
	var newLines []string
	for _, line := range lines {
		if strings.Contains(line, "family =") {
			newLines = append(newLines, fmt.Sprintf("family = \"%s\"", fontFamily))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}
	if !found {
		newLines = append(newLines, "\n[font.normal]", fmt.Sprintf("family = \"%s\"", fontFamily))
	}

	return os.WriteFile(cfgPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func (m *Manager) updateAlacrittyFontSize(size int) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "alacritty", "alacritty.toml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	found := false
	var newLines []string
	for _, line := range lines {
		if strings.Contains(line, "size =") {
			newLines = append(newLines, fmt.Sprintf("size = %d", size))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}
	if !found {
		newLines = append(newLines, fmt.Sprintf("size = %d", size))
	}

	return os.WriteFile(cfgPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func (m *Manager) updateVSCodeFont(fontFamily string) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "Code", "User", "settings.json")
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)

	data, err := os.ReadFile(cfgPath)
	var settings map[string]interface{}
	if err == nil {
		_ = json.Unmarshal(data, &settings)
	}
	if settings == nil {
		settings = make(map[string]interface{})
	}

	settings["editor.fontFamily"] = fmt.Sprintf("'%s', monospace", fontFamily)
	updated, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, updated, 0644)
}

func (m *Manager) updateVSCodeFontSize(size int) error {
	cfgPath := filepath.Join(m.homeDir, ".config", "Code", "User", "settings.json")
	data, err := os.ReadFile(cfgPath)
	var settings map[string]interface{}
	if err == nil {
		_ = json.Unmarshal(data, &settings)
	}
	if settings == nil {
		return nil
	}

	settings["editor.fontSize"] = size
	updated, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, updated, 0644)
}
