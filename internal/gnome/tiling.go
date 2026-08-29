package gnome

import (
	"fmt"
	"os/exec"
)

// ConfigureTilingAssistant sets up window snapping, quarter tiling, custom gaps, and active border highlights
func (m *Manager) ConfigureTilingAssistant(gaps int, accentColor string) error {
	schema := "org.gnome.shell.extensions.tiling-assistant"

	if gaps <= 0 {
		gaps = 6
	}
	if accentColor == "" {
		accentColor = "#7AA2F7"
	}

	// 1. Ensure extension is enabled
	_ = exec.Command("gnome-extensions", "enable", "tiling-assistant@ubuntu.com").Run()

	// 2. Configure Window Gaps & Padding
	_ = SetDconfKey(schema, "window-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "single-screen-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "screen-top-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "screen-bottom-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "screen-left-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "screen-right-gap", fmt.Sprintf("%d", gaps))
	_ = SetDconfKey(schema, "maximize-with-gap", "true")

	// 3. Configure Active Window Border Highlight
	_ = SetDconfKey(schema, "active-window-hint", "1")
	_ = SetDconfKey(schema, "active-window-hint-color", fmt.Sprintf("'%s'", accentColor))
	_ = SetDconfKey(schema, "active-window-hint-border-size", "2")
	_ = SetDconfKey(schema, "active-window-hint-inner-border-size", "0")

	// 4. Snapping & Tiling Shortcuts
	_ = SetDconfKey(schema, "tile-left-half", "['<Super>Left']")
	_ = SetDconfKey(schema, "tile-right-half", "['<Super>Right']")
	_ = SetDconfKey(schema, "tile-top-half", "['<Super>Up']")
	_ = SetDconfKey(schema, "tile-bottom-half", "['<Super>Down']")

	// Quarter-tiling shortcuts
	_ = SetDconfKey(schema, "tile-topleft-quarter", "['<Super><Alt>u']")
	_ = SetDconfKey(schema, "tile-topright-quarter", "['<Super><Alt>i']")
	_ = SetDconfKey(schema, "tile-bottomleft-quarter", "['<Super><Alt>j']")
	_ = SetDconfKey(schema, "tile-bottomright-quarter", "['<Super><Alt>k']")

	// Auto-tile toggle shortcut
	_ = SetDconfKey(schema, "auto-tile", "['<Super>t']")
	_ = SetDconfKey(schema, "enable-tiling-popup", "false")
	_ = SetDconfKey(schema, "dynamic-keybinding-behavior", "0")

	return nil
}

// ConfigureTilingKeybindings sets up vim-style window navigation and tiling hotkeys
func (m *Manager) ConfigureTilingKeybindings() error {
	// Focus window shortcuts (Vim keys: H, J, K, L)
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-left", "['<Shift><Super>Left', '<Shift><Super>h']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-right", "['<Shift><Super>Right', '<Shift><Super>l']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-up", "['<Shift><Super>Up', '<Shift><Super>k']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-down", "['<Shift><Super>Down', '<Shift><Super>j']")

	// Switch workspace shortcuts
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-left", "['<Alt><Super>Left', '<Alt><Super>h']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-right", "['<Alt><Super>Right', '<Alt><Super>l']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-up", "['<Alt><Super>Up', '<Alt><Super>k']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-down", "['<Alt><Super>Down', '<Alt><Super>j']")

	// Window state controls
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "minimize", "['<Super>h']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "show-desktop", "['<Super>d']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "toggle-maximized", "['<Super>m']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "close", "['<Super>q', '<Alt>F4']")

	// Mutter snapping & tiling ergonomics
	_ = SetDconfKey("org.gnome.mutter", "edge-tiling", "true")
	_ = SetDconfKey("org.gnome.mutter", "dynamic-workspaces", "true")
	_ = SetDconfKey("org.gnome.mutter", "workspaces-only-on-primary", "false")

	// Apply Tiling Assistant Gaps & Border Highlight matching current theme
	_ = m.ConfigureTilingAssistant(m.cfg.Desktop.GapsInner, m.cfg.Theme.AccentColor)

	return nil
}
