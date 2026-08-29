package gnome

import (
	"os/exec"
	"strings"

	"github.com/vula-os/vula/internal/config"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

type ExtensionInfo struct {
	UUID        string
	Name        string
	Description string
	Essential   bool
}

var RecommendedExtensions = []ExtensionInfo{
	{
		UUID:        "forge@jmmaranan.com",
		Name:        "Forge Tiling",
		Description: "Tiling window management with keyboard shortcuts",
		Essential:   true,
	},
	{
		UUID:        "blur-my-shell@aunetx",
		Name:        "Blur my Shell",
		Description: "Aesthetic blur effect on panel, dash, and overview",
		Essential:   false,
	},
	{
		UUID:        "appindicatorsupport@rgcjonas.gmail.com",
		Name:        "AppIndicator Support",
		Description: "System tray icons for background services",
		Essential:   true,
	},
	{
		UUID:        "just-perfection-desktop@just-perfection",
		Name:        "Just Perfection",
		Description: "Fine-grained control to minimize GNOME Shell UI clutter",
		Essential:   false,
	},
}

// SetDconfKey applies a dconf setting via dconf write or gsettings
func SetDconfKey(schema, key, value string) error {
	cmd := exec.Command("gsettings", "set", schema, key, value)
	return cmd.Run()
}

// GetDconfKey retrieves a dconf setting
func GetDconfKey(schema, key string) (string, error) {
	out, err := exec.Command("gsettings", "get", schema, key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ApplyDesktopOptimizations tunes GNOME Shell for high performance and developer ergonomics
func (m *Manager) ApplyDesktopOptimizations() error {
	// 1. Dark Mode
	_ = SetDconfKey("org.gnome.desktop.interface", "color-scheme", "'prefer-dark'")
	_ = SetDconfKey("org.gnome.desktop.interface", "gtk-theme", "'Yaru-dark'")

	// 2. Window management tweaks
	_ = SetDconfKey("org.gnome.desktop.wm.preferences", "button-layout", "'appmenu:close'")
	_ = SetDconfKey("org.gnome.mutter", "center-new-windows", "true")
	_ = SetDconfKey("org.gnome.mutter", "attach-modal-dialogs", "true")

	// 3. Fonts
	_ = SetDconfKey("org.gnome.desktop.interface", "font-name", "'Ubuntu Sans 11'")
	_ = SetDconfKey("org.gnome.desktop.interface", "monospace-font-name", "'JetBrains Mono 11'")

	// 4. Configure Developer Keybindings
	return m.ConfigureKeybindings()
}

// ConfigureKeybindings sets up keyboard-centric shortcuts
func (m *Manager) ConfigureKeybindings() error {
	// Custom Keybinding 1: Vula Floating HUD
	bindingPath := "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings"
	custom0 := bindingPath + "/custom0/"

	_ = exec.Command("dconf", "write", custom0+"name", "'Vula HUD'").Run()
	_ = exec.Command("dconf", "write", custom0+"command", "'vula hud'").Run()
	_ = exec.Command("dconf", "write", custom0+"binding", "'<Super>space'").Run()

	// Custom Keybinding 2: Voice Dictation
	custom1 := bindingPath + "/custom1/"
	_ = exec.Command("dconf", "write", custom1+"name", "'Vula Voice Dictate'").Run()
	_ = exec.Command("dconf", "write", custom1+"command", "'vula voice record'").Run()
	_ = exec.Command("dconf", "write", custom1+"binding", "'<Super><Alt>v'").Run()

	// Register custom keybinding list
	_ = exec.Command("dconf", "write", "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings",
		"['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/', '/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/']").Run()

	// Core navigation shortcuts
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "close", "['<Super>q', '<Alt>F4']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "toggle-maximized", "['<Super>m']")

	return nil
}

// ListInstalledExtensions returns enabled extension UUIDs
func ListInstalledExtensions() ([]string, error) {
	out, err := exec.Command("gnome-extensions", "list", "--enabled").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result, nil
}

// EnableExtension enables a GNOME shell extension by UUID
func EnableExtension(uuid string) error {
	cmd := exec.Command("gnome-extensions", "enable", uuid)
	return cmd.Run()
}
