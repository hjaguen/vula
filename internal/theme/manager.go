package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/gnome"
)

type ThemePalette struct {
	Name        string
	DisplayName string
	AccentColor string
	GtkTheme    string
	GnomeAccent string // orange, blue, teal, purple, red, etc.
	Background  string
	Foreground  string
}

var AvailableThemes = map[string]ThemePalette{
	"tokyonight": {
		Name:        "tokyonight",
		DisplayName: "Tokyo Night",
		AccentColor: "#7AA2F7",
		GtkTheme:    "Yaru-dark",
		GnomeAccent: "purple",
		Background:  "#1A1B26",
		Foreground:  "#C0CAF5",
	},
	"catppuccin": {
		Name:        "catppuccin",
		DisplayName: "Catppuccin Mocha",
		AccentColor: "#CBA6F7",
		GtkTheme:    "Yaru-dark",
		GnomeAccent: "purple",
		Background:  "#1E1E2E",
		Foreground:  "#CDD6F4",
	},
	"nord": {
		Name:        "nord",
		DisplayName: "Nord Arctic",
		AccentColor: "#88C0D0",
		GtkTheme:    "Yaru-dark",
		GnomeAccent: "blue",
		Background:  "#2E3440",
		Foreground:  "#D8DEE9",
	},
	"rose-pine": {
		Name:        "rose-pine",
		DisplayName: "Rosé Pine",
		AccentColor: "#EBBCBA",
		GtkTheme:    "Yaru-dark",
		GnomeAccent: "red",
		Background:  "#191724",
		Foreground:  "#E0DEF4",
	},
}

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// ApplyTheme sets the theme across GNOME Shell, Starship prompt, and terminal
func (m *Manager) ApplyTheme(themeName string) error {
	palette, exists := AvailableThemes[themeName]
	if !exists {
		return fmt.Errorf("unknown theme '%s'. Available: tokyonight, catppuccin, nord, rose-pine", themeName)
	}

	// 1. GNOME Shell & Accent Color
	_ = gnome.SetDconfKey("org.gnome.desktop.interface", "color-scheme", "'prefer-dark'")
	_ = gnome.SetDconfKey("org.gnome.desktop.interface", "accent-color", fmt.Sprintf("'%s'", palette.GnomeAccent))

	// 2. Update Vula Config
	m.cfg.Theme.Palette = themeName
	m.cfg.Theme.AccentColor = palette.AccentColor
	_ = config.SaveConfig(m.cfg)

	// 3. Update Ghostty / Terminal configs if present
	home := os.Getenv("HOME")
	ghosttyDir := filepath.Join(home, ".config", "ghostty")
	if err := os.MkdirAll(ghosttyDir, 0755); err == nil {
		ghosttyCfg := fmt.Sprintf("theme = %s\nfont-family = JetBrains Mono\nfont-size = 12\nbackground-blur-radius = 20\n", themeName)
		_ = os.WriteFile(filepath.Join(ghosttyDir, "config"), []byte(ghosttyCfg), 0644)
	}

	return nil
}

func ListThemes() []ThemePalette {
	list := make([]ThemePalette, 0, len(AvailableThemes))
	for _, t := range AvailableThemes {
		list = append(list, t)
	}
	return list
}
