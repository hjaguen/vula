package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/gnome"
	"gopkg.in/yaml.v3"
)

type ThemePalette struct {
	Name           string `yaml:"name"`
	DisplayName    string `yaml:"display_name"`
	AccentColor    string `yaml:"accent_color"`
	SecondaryColor string `yaml:"secondary_color,omitempty"`
	GtkTheme       string `yaml:"gtk_theme"`
	GnomeAccent    string `yaml:"gnome_accent"` // orange, blue, teal, purple, red, green, yellow, slate
	Background     string `yaml:"background"`
	Foreground     string `yaml:"foreground"`
	HeaderBg       string `yaml:"header_bg,omitempty"`
	TitlebarBg     string `yaml:"titlebar_bg,omitempty"`
	TitlebarFg     string `yaml:"titlebar_fg,omitempty"`
	IsCustom       bool   `yaml:"is_custom,omitempty"`
}

var BuiltInThemes = map[string]ThemePalette{
	"tokyonight": {
		Name:           "tokyonight",
		DisplayName:    "Tokyo Night",
		AccentColor:    "#7AA2F7",
		SecondaryColor: "#BB9AF7",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "purple",
		Background:     "#1A1B26",
		Foreground:     "#C0CAF5",
		HeaderBg:       "#16161E",
		TitlebarBg:     "#1F2335",
		TitlebarFg:     "#C0CAF5",
	},
	"catppuccin": {
		Name:           "catppuccin",
		DisplayName:    "Catppuccin Mocha",
		AccentColor:    "#CBA6F7",
		SecondaryColor: "#89B4FA",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "purple",
		Background:     "#1E1E2E",
		Foreground:     "#CDD6F4",
		HeaderBg:       "#11111B",
		TitlebarBg:     "#181825",
		TitlebarFg:     "#CDD6F4",
	},
	"nord": {
		Name:           "nord",
		DisplayName:    "Nord Arctic",
		AccentColor:    "#88C0D0",
		SecondaryColor: "#81A1C1",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "blue",
		Background:     "#2E3440",
		Foreground:     "#D8DEE9",
		HeaderBg:       "#242933",
		TitlebarBg:     "#2E3440",
		TitlebarFg:     "#D8DEE9",
	},
	"rose-pine": {
		Name:           "rose-pine",
		DisplayName:    "Rosé Pine",
		AccentColor:    "#EBBCBA",
		SecondaryColor: "#31748F",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "red",
		Background:     "#191724",
		Foreground:     "#E0DEF4",
		HeaderBg:       "#16141F",
		TitlebarBg:     "#1F1D2E",
		TitlebarFg:     "#E0DEF4",
	},
}

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// GetThemesDir returns the path to custom user themes (~/.config/vula/themes)
func GetThemesDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "vula", "themes")
}

// GetAllThemes returns both built-in and user-created custom themes
func (m *Manager) GetAllThemes() map[string]ThemePalette {
	themes := make(map[string]ThemePalette)
	for k, v := range BuiltInThemes {
		themes[k] = v
	}

	// Load custom user themes from disk
	customDir := GetThemesDir()
	files, err := os.ReadDir(customDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				data, err := os.ReadFile(filepath.Join(customDir, f.Name()))
				if err == nil {
					var p ThemePalette
					if err := yaml.Unmarshal(data, &p); err == nil && p.Name != "" {
						p.IsCustom = true
						themes[p.Name] = p
					}
				}
			}
		}
	}

	return themes
}

// SaveCustomTheme writes a custom theme to ~/.config/vula/themes/<name>.yaml
func (m *Manager) SaveCustomTheme(p ThemePalette) error {
	dir := GetThemesDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	p.IsCustom = true
	if p.GtkTheme == "" {
		p.GtkTheme = "Yaru-dark"
	}
	if p.GnomeAccent == "" {
		p.GnomeAccent = "purple"
	}

	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, p.Name+".yaml")
	return os.WriteFile(filePath, data, 0644)
}

// ApplyTheme sets the theme across GNOME Shell, desktop wallpaper, GTK titlebars, Starship prompt, and terminal
func (m *Manager) ApplyTheme(themeName string) error {
	allThemes := m.GetAllThemes()
	palette, exists := allThemes[themeName]
	if !exists {
		return fmt.Errorf("unknown theme '%s'. Run 'vula theme list' to see available themes", themeName)
	}

	// 1. GNOME Shell & Accent Color
	_ = gnome.SetDconfKey("org.gnome.desktop.interface", "color-scheme", "'prefer-dark'")
	_ = gnome.SetDconfKey("org.gnome.desktop.interface", "accent-color", fmt.Sprintf("'%s'", palette.GnomeAccent))

	// 2. GTK 3 & GTK 4 Custom CSS Override for Window Titlebars & Top Bar
	syncGtkCustomCss(palette)

	// 3. Generate and Set Dynamic High-Res Theme Wallpaper
	wpPath, err := EnsureThemeWallpaper(palette)
	if err == nil && wpPath != "" {
		_ = gnome.SetDconfKey("org.gnome.desktop.background", "picture-uri", fmt.Sprintf("'file://%s'", wpPath))
		_ = gnome.SetDconfKey("org.gnome.desktop.background", "picture-uri-dark", fmt.Sprintf("'file://%s'", wpPath))
		_ = gnome.SetDconfKey("org.gnome.desktop.background", "picture-options", "'zoom'")
	}

	// 4. Update Vula Config & Tiling Assistant Border Highlight
	m.cfg.Theme.Palette = themeName
	m.cfg.Theme.AccentColor = palette.AccentColor
	_ = config.SaveConfig(m.cfg)

	gnomeMgr := gnome.NewManager(m.cfg)
	_ = gnomeMgr.ConfigureTilingAssistant(m.cfg.Desktop.GapsInner, palette.AccentColor)

	// 5. Update Ghostty / Terminal configs if present
	home := os.Getenv("HOME")
	ghosttyDir := filepath.Join(home, ".config", "ghostty")
	if err := os.MkdirAll(ghosttyDir, 0755); err == nil {
		ghosttyCfg := fmt.Sprintf("# Auto-generated by Vula Theme Engine\ntheme = %s\nbackground = %s\nforeground = %s\nfont-family = JetBrains Mono\nfont-size = 12\nbackground-blur-radius = 20\n",
			palette.Name, palette.Background, palette.Foreground)
		_ = os.WriteFile(filepath.Join(ghosttyDir, "config"), []byte(ghosttyCfg), 0644)
	}

	// 6. Update Tmux Status Bar Accent Color
	tmuxFile := filepath.Join(home, ".tmux.conf")
	if _, err := os.Stat(tmuxFile); err == nil {
		syncTmuxTheme(tmuxFile, palette)
	}

	return nil
}

func syncGtkCustomCss(p ThemePalette) {
	home := os.Getenv("HOME")
	hdrBg := p.HeaderBg
	if hdrBg == "" {
		hdrBg = p.Background
	}
	tbBg := p.TitlebarBg
	if tbBg == "" {
		tbBg = p.Background
	}
	tbFg := p.TitlebarFg
	if tbFg == "" {
		tbFg = p.Foreground
	}

	cssContent := fmt.Sprintf(`/* Auto-generated by Vula Theme Engine - %s */
@define-color accent_color %s;
@define-color accent_bg_color %s;
@define-color window_bg_color %s;
@define-color window_fg_color %s;

headerbar, .titlebar, windowtitle, .header-bar {
    background-color: %s !important;
    color: %s !important;
    border-bottom: 2px solid %s !important;
    box-shadow: none !important;
}

headerbar label, windowtitle label {
    color: %s !important;
    font-weight: bold !important;
}

.gnome-shell-panel, #panel, headerbar.flat {
    background-color: %s !important;
    color: %s !important;
}
`, p.DisplayName, p.AccentColor, p.AccentColor, p.Background, p.Foreground, tbBg, tbFg, p.AccentColor, tbFg, hdrBg, p.Foreground)

	dirs := []string{
		filepath.Join(home, ".config", "gtk-3.0"),
		filepath.Join(home, ".config", "gtk-4.0"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err == nil {
			_ = os.WriteFile(filepath.Join(d, "gtk.css"), []byte(cssContent), 0644)
		}
	}
}

// EnsureThemeWallpaper creates a high-definition 4K SVG art wallpaper matching the palette
func EnsureThemeWallpaper(p ThemePalette) (string, error) {
	home := os.Getenv("HOME")
	wpDir := filepath.Join(home, ".config", "vula", "wallpapers")
	if err := os.MkdirAll(wpDir, 0755); err != nil {
		return "", err
	}

	wpPath := filepath.Join(wpDir, fmt.Sprintf("%s.svg", p.Name))

	secColor := p.SecondaryColor
	if secColor == "" {
		secColor = p.AccentColor
	}

	svgContent := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="3840" height="2160" viewBox="0 0 3840 2160">
  <defs>
    <radialGradient id="bgGlow" cx="50%%" cy="35%%" r="75%%">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.35"/>
      <stop offset="40%%" stop-color="%s" stop-opacity="0.15"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="1"/>
    </radialGradient>
    <linearGradient id="accentGrad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="%s"/>
      <stop offset="100%%" stop-color="%s"/>
    </linearGradient>
    <linearGradient id="waveGrad1" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.2"/>
      <stop offset="50%%" stop-color="%s" stop-opacity="0.4"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.1"/>
    </linearGradient>
    <linearGradient id="waveGrad2" x1="0%%" y1="100%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.3"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.05"/>
    </linearGradient>
    <filter id="glow" x="-20%%" y="-20%%" width="140%%" height="140%%">
      <feGaussianBlur stdDeviation="15" result="blur"/>
      <feComposite in="SourceGraphic" in2="blur" operator="over"/>
    </filter>
  </defs>

  <!-- Deep Background Layer -->
  <rect width="100%%" height="100%%" fill="%s"/>
  <rect width="100%%" height="100%%" fill="url(#bgGlow)"/>

  <!-- Cyber Grid Lines -->
  <g stroke="%s" stroke-opacity="0.05" stroke-width="1">
    <path d="M0,270 H3840 M0,540 H3840 M0,810 H3840 M0,1080 H3840 M0,1350 H3840 M0,1620 H3840 M0,1890 H3840"/>
    <path d="M480,0 V2160 M960,0 V2160 M1440,0 V2160 M1920,0 V2160 M2400,0 V2160 M2880,0 V2160 M3360,0 V2160"/>
  </g>

  <!-- Organic Wave Flow Paths -->
  <path d="M-100,1400 Q800,900 1920,1300 T3940,1100 L3940,2260 L-100,2260 Z" fill="url(#waveGrad1)"/>
  <path d="M-100,1600 Q1100,1100 2200,1600 T3940,1400 L3940,2260 L-100,2260 Z" fill="url(#waveGrad2)"/>

  <!-- Concentric Geometric Rings -->
  <circle cx="1920" cy="1080" r="450" fill="none" stroke="url(#accentGrad)" stroke-width="3" opacity="0.35" filter="url(#glow)"/>
  <circle cx="1920" cy="1080" r="750" fill="none" stroke="url(#accentGrad)" stroke-width="1.5" opacity="0.2"/>
  <circle cx="1920" cy="1080" r="1050" fill="none" stroke="%s" stroke-width="1" stroke-dasharray="8,12" opacity="0.25"/>

  <!-- Centerpiece Vula Bolt Emblem -->
  <g transform="translate(1860, 1010) scale(1.6)" filter="url(#glow)">
    <polygon points="40,0 10,50 35,50 20,90 60,40 35,40" fill="url(#accentGrad)"/>
  </g>
</svg>`,
		p.AccentColor, secColor, p.Background,
		p.AccentColor, secColor,
		p.AccentColor, secColor, p.AccentColor,
		secColor, p.AccentColor,
		p.Background,
		p.Foreground,
		secColor,
	)

	if err := os.WriteFile(wpPath, []byte(svgContent), 0644); err != nil {
		return "", err
	}

	return wpPath, nil
}

func syncTmuxTheme(tmuxFile string, p ThemePalette) {
	data, err := os.ReadFile(tmuxFile)
	if err != nil {
		return
	}
	content := string(data)
	// Replace status bar color
	newLeft := fmt.Sprintf("set -g status-left ' #[bold,fg=%s]⚡ VULA #[default]| '", p.AccentColor)
	newStyle := fmt.Sprintf("set -g status-style bg='%s',fg='%s'", p.Background, p.Foreground)

	lines := strings.Split(content, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "set -g status-left") {
			lines[i] = newLeft
		} else if strings.HasPrefix(l, "set -g status-style") {
			lines[i] = newStyle
		}
	}
	_ = os.WriteFile(tmuxFile, []byte(strings.Join(lines, "\n")), 0644)
}

func (m *Manager) ListThemes() []ThemePalette {
	allThemes := m.GetAllThemes()
	list := make([]ThemePalette, 0, len(allThemes))
	for _, t := range allThemes {
		list = append(list, t)
	}
	return list
}
