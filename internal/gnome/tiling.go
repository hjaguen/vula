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

	// Auto-tile toggle shortcut for Tiling Assistant
	_ = SetDconfKey(schema, "auto-tile", "['<Super><Alt>t']")
	_ = SetDconfKey(schema, "enable-tiling-popup", "false")
	_ = SetDconfKey(schema, "dynamic-keybinding-behavior", "0")

	return nil
}

// ConfigureTactile sets up Tactile interactive grid overlay shortcuts (Omakub style)
func (m *Manager) ConfigureTactile() error {
	schema := "org.gnome.shell.extensions.tactile"

	// 1. Ensure Tactile extension is installed & enabled
	_ = m.InstallExtension("tactile@lundal.io")
	_ = exec.Command("gdbus", "call", "--session", "--dest", "org.gnome.Shell", "--object-path", "/org/gnome/Shell", "--method", "org.gnome.Shell.Extensions.EnableExtension", "tactile@lundal.io").Run()
	_ = exec.Command("gnome-extensions", "enable", "tactile@lundal.io").Run()

	// 2. Set Super+t to open Tactile grid overlay exclusively
	_ = SetDconfKey(schema, "show-tiles", "['<Super>t']")

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

	// Clear conflicting GNOME shell message tray shortcut for Super+m (keep Super+v)
	_ = SetDconfKey("org.gnome.shell.keybindings", "toggle-message-tray", "['<Super>v']")

	// Window state controls
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "minimize", "['<Super>h']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "show-desktop", "['<Super>d']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "toggle-maximized", "['<Super>m']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "close", "['<Super>q', '<Alt>F4']")

	// Mouse Window Modifier: Super + Left Click anywhere inside window to drag/move
	_ = SetDconfKey("org.gnome.desktop.wm.preferences", "mouse-button-modifier", "'<Super>'")

	// Mutter snapping & tiling ergonomics
	_ = SetDconfKey("org.gnome.mutter", "edge-tiling", "true")
	_ = SetDconfKey("org.gnome.mutter", "dynamic-workspaces", "true")
	_ = SetDconfKey("org.gnome.mutter", "workspaces-only-on-primary", "false")

	// Switch to direct workspace 1..6 (Super + 1..6)
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-1", "['<Super>1']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-2", "['<Super>2']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-3", "['<Super>3']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-4", "['<Super>4']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-5", "['<Super>5']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "switch-to-workspace-6", "['<Super>6']")

	// Move window to workspace 1..6 (Shift + Super + 1..6)
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-1", "['<Shift><Super>1']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-2", "['<Shift><Super>2']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-3", "['<Shift><Super>3']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-4", "['<Shift><Super>4']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-5", "['<Shift><Super>5']")
	_ = SetDconfKey("org.gnome.desktop.wm.keybindings", "move-to-workspace-6", "['<Shift><Super>6']")

	// Dock App Quick Jump (Alt + 1..9)
	for i := 1; i <= 9; i++ {
		_ = SetDconfKey("org.gnome.shell.keybindings", fmt.Sprintf("app-shift-%d", i), fmt.Sprintf("['<Alt>%d']", i))
	}

	// Apply Tiling Assistant Gaps & Border Highlight matching current theme
	_ = m.ConfigureTilingAssistant(m.cfg.Desktop.GapsInner, m.cfg.Theme.AccentColor)

	// Configure Tactile Grid Overlay (Omakub style)
	_ = m.ConfigureTactile()

	return nil
}
