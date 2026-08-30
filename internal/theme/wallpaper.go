package theme

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/gnome"
)

type WallpaperManager struct {
	cfg *config.Config
}

func NewWallpaperManager(cfg *config.Config) *WallpaperManager {
	return &WallpaperManager{cfg: cfg}
}

func (wm *WallpaperManager) GetWallpaperDir() string {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".config", "vula", "wallpapers")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// ListWallpapers returns all wallpaper files in ~/.config/vula/wallpapers/
func (wm *WallpaperManager) ListWallpapers() ([]string, error) {
	dir := wm.GetWallpaperDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".svg" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				files = append(files, filepath.Join(dir, e.Name()))
			}
		}
	}

	return files, nil
}

// SetWallpaper applies a specific wallpaper path to GNOME desktop
func (wm *WallpaperManager) SetWallpaper(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("wallpaper file not found: %s", path)
	}

	uri := "file://" + path
	if err := gnome.SetDconfKey("org.gnome.desktop.background", "picture-uri", fmt.Sprintf("'%s'", uri)); err != nil {
		return err
	}
	return gnome.SetDconfKey("org.gnome.desktop.background", "picture-uri-dark", fmt.Sprintf("'%s'", uri))
}

// RotateNextWallpaper switches to the next wallpaper in the folder
func (wm *WallpaperManager) RotateNextWallpaper() (string, error) {
	files, err := wm.ListWallpapers()
	if err != nil || len(files) == 0 {
		// Fallback: Generate SVG wallpaper matching active theme
		mgr := NewManager(wm.cfg)
		p, exists := mgr.GetAllThemes()[wm.cfg.Theme.Palette]
		if !exists {
			p = BuiltInThemes["tokyonight"]
		}
		themeSvg, _ := EnsureThemeWallpaper(p)
		_ = wm.SetWallpaper(themeSvg)
		return themeSvg, nil
	}

	// Read current wallpaper
	current, _ := gnome.GetDconfKey("org.gnome.desktop.background", "picture-uri-dark")
	current = strings.Trim(current, "'")
	current = strings.TrimPrefix(current, "file://")

	nextIdx := 0
	for i, f := range files {
		if f == current {
			nextIdx = (i + 1) % len(files)
			break
		}
	}

	nextPath := files[nextIdx]
	if err := wm.SetWallpaper(nextPath); err != nil {
		return "", err
	}
	return nextPath, nil
}

// FetchPresetWallpapers downloads curated theme wallpapers into ~/.config/vula/wallpapers/
func (wm *WallpaperManager) FetchPresetWallpapers() error {
	dir := wm.GetWallpaperDir()
	mgr := NewManager(wm.cfg)
	allThemes := mgr.GetAllThemes()

	// 1. Always generate SVGs for built-in themes
	for themeName, palette := range allThemes {
		_, _ = EnsureThemeWallpaper(palette)
		_ = themeName
	}

	// 2. Curated external HD wallpapers (raw URLs)
	presets := map[string]string{
		"catppuccin_waves.png": "https://raw.githubusercontent.com/zhichaoh/catppuccin-wallpapers/main/waves.png",
		"tokyonight_night.png": "https://raw.githubusercontent.com/zatchyzh/tokyo-night-wallpapers/main/wallpapers/tokyo_night_mountain.png",
	}

	for filename, url := range presets {
		target := filepath.Join(dir, filename)
		if _, err := os.Stat(target); os.IsNotExist(err) {
			_ = downloadFile(url, target)
		}
	}

	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
