package gnome

import (
	"fmt"
	"os/exec"
)

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

	return nil
}

// EnableForgeTiling configures Forge extension settings if installed
func (m *Manager) EnableForgeTiling() error {
	schema := "org.gnome.shell.extensions.forge"
	// Check if schema is installed
	if err := exec.Command("gsettings", "describe", schema, "tiling-mode-enabled").Run(); err == nil {
		_ = SetDconfKey(schema, "tiling-mode-enabled", "true")
		_ = SetDconfKey(schema, "window-gap-size", fmt.Sprintf("%d", m.cfg.Desktop.GapsInner))
		_ = SetDconfKey(schema, "dnd-center-layout", "true")
	}
	return nil
}
