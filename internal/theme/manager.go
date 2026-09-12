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
		AccentColor:    "#89B4FA",
		SecondaryColor: "#CBA6F7",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "purple",
		Background:     "#1E1E2E",
		Foreground:     "#CDD6F4",
		HeaderBg:       "#11111B",
		TitlebarBg:     "#181825",
		TitlebarFg:     "#CDD6F4",
	},
	"catppuccin-latte": {
		Name:           "catppuccin-latte",
		DisplayName:    "Catppuccin Latte",
		AccentColor:    "#1E66F5",
		SecondaryColor: "#EA76CB",
		GtkTheme:       "Yaru",
		GnomeAccent:    "blue",
		Background:     "#EFF1F5",
		Foreground:     "#4C4F69",
		HeaderBg:       "#E6E9EF",
		TitlebarBg:     "#DCE0E8",
		TitlebarFg:     "#4C4F69",
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
		DisplayName:    "Rosé Pine Dawn",
		AccentColor:    "#56949F",
		SecondaryColor: "#907AA9",
		GtkTheme:       "Yaru",
		GnomeAccent:    "teal",
		Background:     "#FAF4ED",
		Foreground:     "#575279",
		HeaderBg:       "#F2E9E1",
		TitlebarBg:     "#FFFFFF",
		TitlebarFg:     "#575279",
	},
	"kanagawa": {
		Name:           "kanagawa",
		DisplayName:    "Kanagawa Wave",
		AccentColor:    "#7E9CD8",
		SecondaryColor: "#957FB8",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "blue",
		Background:     "#1F1F28",
		Foreground:     "#DCD7BA",
		HeaderBg:       "#16161D",
		TitlebarBg:     "#2A2A37",
		TitlebarFg:     "#DCD7BA",
	},
	"gruvbox": {
		Name:           "gruvbox",
		DisplayName:    "Gruvbox Dark",
		AccentColor:    "#7DAEA3",
		SecondaryColor: "#D8A657",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "yellow",
		Background:     "#282828",
		Foreground:     "#D4BE98",
		HeaderBg:       "#1D2021",
		TitlebarBg:     "#3C3836",
		TitlebarFg:     "#D4BE98",
	},
	"everforest": {
		Name:           "everforest",
		DisplayName:    "Everforest Dark",
		AccentColor:    "#7FBBB3",
		SecondaryColor: "#A7C080",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "green",
		Background:     "#2D353B",
		Foreground:     "#D3C6AA",
		HeaderBg:       "#232A2E",
		TitlebarBg:     "#343F44",
		TitlebarFg:     "#D3C6AA",
	},
	"ethereal": {
		Name:           "ethereal",
		DisplayName:    "Ethereal Glow",
		AccentColor:    "#7D82D9",
		SecondaryColor: "#FFCEAD",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "purple",
		Background:     "#060B1E",
		Foreground:     "#FFCEAD",
		HeaderBg:       "#040714",
		TitlebarBg:     "#0C1433",
		TitlebarFg:     "#FFCEAD",
	},
	"hackerman": {
		Name:           "hackerman",
		DisplayName:    "Hackerman Cyberpunk",
		AccentColor:    "#82FB9C",
		SecondaryColor: "#50F7D4",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "green",
		Background:     "#0B0C16",
		Foreground:     "#DDF7FF",
		HeaderBg:       "#06070D",
		TitlebarBg:     "#121424",
		TitlebarFg:     "#DDF7FF",
	},
	"lumon": {
		Name:           "lumon",
		DisplayName:    "Lumon Slate",
		AccentColor:    "#8BC9EB",
		SecondaryColor: "#6FB8E3",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "blue",
		Background:     "#16242D",
		Foreground:     "#D6E2EE",
		HeaderBg:       "#0F1920",
		TitlebarBg:     "#1D303C",
		TitlebarFg:     "#D6E2EE",
	},
	"matte-black": {
		Name:           "matte-black",
		DisplayName:    "Matte Black",
		AccentColor:    "#E68E0D",
		SecondaryColor: "#FFC107",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "orange",
		Background:     "#121212",
		Foreground:     "#BEBEBE",
		HeaderBg:       "#0A0A0A",
		TitlebarBg:     "#1E1E1E",
		TitlebarFg:     "#BEBEBE",
	},
	"miasma": {
		Name:           "miasma",
		DisplayName:    "Miasma Forest",
		AccentColor:    "#78824B",
		SecondaryColor: "#C9A554",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "green",
		Background:     "#222222",
		Foreground:     "#C2C2B0",
		HeaderBg:       "#181818",
		TitlebarBg:     "#2C2C2C",
		TitlebarFg:     "#C2C2B0",
	},
	"osaka-jade": {
		Name:           "osaka-jade",
		DisplayName:    "Osaka Jade",
		AccentColor:    "#509475",
		SecondaryColor: "#2DD5B7",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "teal",
		Background:     "#111C18",
		Foreground:     "#C1C497",
		HeaderBg:       "#0B1310",
		TitlebarBg:     "#192923",
		TitlebarFg:     "#C1C497",
	},
	"retro-82": {
		Name:           "retro-82",
		DisplayName:    "Retro 82 Synthwave",
		AccentColor:    "#FAA968",
		SecondaryColor: "#028391",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "orange",
		Background:     "#05182E",
		Foreground:     "#F6DCAC",
		HeaderBg:       "#030F1D",
		TitlebarBg:     "#092444",
		TitlebarFg:     "#F6DCAC",
	},
	"ristretto": {
		Name:           "ristretto",
		DisplayName:    "Ristretto Espresso",
		AccentColor:    "#F38D70",
		SecondaryColor: "#A8A9EB",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "red",
		Background:     "#2C2525",
		Foreground:     "#E6D9DB",
		HeaderBg:       "#1E1919",
		TitlebarBg:     "#3A3131",
		TitlebarFg:     "#E6D9DB",
	},
	"vantablack": {
		Name:           "vantablack",
		DisplayName:    "Vantablack OLED",
		AccentColor:    "#8D8D8D",
		SecondaryColor: "#ECECEC",
		GtkTheme:       "Yaru-dark",
		GnomeAccent:    "slate",
		Background:     "#000000",
		Foreground:     "#FFFFFF",
		HeaderBg:       "#000000",
		TitlebarBg:     "#111111",
		TitlebarFg:     "#FFFFFF",
	},
	"flexoki-light": {
		Name:           "flexoki-light",
		DisplayName:    "Flexoki Paper Light",
		AccentColor:    "#205EA6",
		SecondaryColor: "#CE5D97",
		GtkTheme:       "Yaru",
		GnomeAccent:    "blue",
		Background:     "#FFFCF0",
		Foreground:     "#100F0F",
		HeaderBg:       "#F2EFE9",
		TitlebarBg:     "#E6E4D9",
		TitlebarFg:     "#100F0F",
	},
	"white": {
		Name:           "white",
		DisplayName:    "Pure Crisp Light",
		AccentColor:    "#6E6E6E",
		SecondaryColor: "#1A1A1A",
		GtkTheme:       "Yaru",
		GnomeAccent:    "slate",
		Background:     "#FFFFFF",
		Foreground:     "#000000",
		HeaderBg:       "#F5F5F5",
		TitlebarBg:     "#EAEAEA",
		TitlebarFg:     "#000000",
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
	colorScheme := "'prefer-dark'"
	if palette.Name == "catppuccin-latte" || palette.Name == "flexoki-light" || palette.Name == "white" || palette.Name == "rose-pine" {
		colorScheme = "'prefer-light'"
	}
	_ = gnome.SetDconfKey("org.gnome.desktop.interface", "color-scheme", colorScheme)
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

/* Hide GTK HeaderBar / Window Titlebar specifically for Vula HUD floating window */
window.vula-hud headerbar,
window.vula-hud .titlebar,
window.vula-hud headerbar *,
window.vula-hud .titlebar *,
.vula-hud headerbar,
.vula-hud .titlebar,
.vula-hud headerbar * {
    min-height: 0px !important;
    height: 0px !important;
    padding: 0px !important;
    margin: 0px !important;
    border: none !important;
    background: transparent !important;
    box-shadow: none !important;
    font-size: 0px !important;
    opacity: 0 !important;
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

var ThemePreferredWallpapers = map[string]string{
	"tokyonight":       "tokyo-night_0-tokyo-night.jpg",
	"catppuccin":       "catppuccin_2-waves.png",
	"catppuccin-latte": "catppuccin-latte_1-color-fade.png",
	"nord":             "nord_1-city-view.png",
	"rose-pine":        "rose-pine_1-funky-shapes.jpg",
	"kanagawa":         "kanagawa_1-kanagawa.jpg",
	"gruvbox":          "gruvbox_1-the-backwater.jpg",
	"everforest":       "everforest_1-tree-tops.jpg",
	"ethereal":         "ethereal_1-cosmic.jpg",
	"hackerman":        "hackerman_1-synth-scape.jpg",
	"lumon":            "lumon_01-united-in-severance.jpg",
	"matte-black":      "matte-black_0-ship-at-sea.jpg",
	"miasma":           "miasma_01-nature-of-fe.png",
	"osaka-jade":       "osaka-jade_1-growing-city.jpg",
	"retro-82":         "retro-82_1-launch.png",
	"ristretto":        "ristretto_0-launch.png",
	"vantablack":       "vantablack_0-dot-hands.jpg",
	"flexoki-light":    "flexoki-light_1-orb.png",
	"white":            "white_1-white.jpg",
}

// EnsureThemeWallpaper creates or resolves the theme wallpaper path
func EnsureThemeWallpaper(p ThemePalette) (string, error) {
	home := os.Getenv("HOME")
	wpDir := filepath.Join(home, ".config", "vula", "wallpapers")
	if err := os.MkdirAll(wpDir, 0755); err != nil {
		return "", err
	}

	// 1. Check for preferred wallpaper image in ~/.config/vula/wallpapers/
	if prefImg, exists := ThemePreferredWallpapers[p.Name]; exists {
		imgPath := filepath.Join(wpDir, prefImg)
		if _, err := os.Stat(imgPath); err == nil {
			return imgPath, nil
		}
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
