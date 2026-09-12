package keys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/gnome"
	"github.com/vula-os/vula/internal/ui"
)

type Keybinding struct {
	ID          string `json:"id" yaml:"id"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"`
	Shortcut    string `json:"shortcut" yaml:"shortcut"` // e.g. <Super>space, <Super><Alt>v
	Command     string `json:"command" yaml:"command"`
}

var DefaultKeybindings = []Keybinding{
	// Launchers & System
	{
		ID:          "hud",
		Title:       "Vula HUD Launcher",
		Description: "Toggle Spotlight-style application & command launcher",
		Category:    "System & Launchers",
		Shortcut:    "<Super>space",
		Command:     "~/.local/bin/vula-hud-launch",
	},
	{
		ID:          "keys_map",
		Title:       "Keybindings Cheat-Sheet",
		Description: "Open interactive keyboard shortcuts map",
		Category:    "System & Launchers",
		Shortcut:    "<Super><Shift>k",
		Command:     "vula keys",
	},

	// Voice & AI
	{
		ID:          "voice_dictate",
		Title:       "Voice Dictation",
		Description: "Record voice & type directly into active window",
		Category:    "Voice & AI",
		Shortcut:    "<Super><Alt>v",
		Command:     "~/.local/bin/vula voice dictate",
	},
	{
		ID:          "voice_assistant",
		Title:       "Active Voice AI Assistant",
		Description: "Conversational voice AI with Whisper STT & Piper TTS",
		Category:    "Voice & AI",
		Shortcut:    "<Super><Alt>a",
		Command:     "~/.local/bin/vula listen",
	},
	{
		ID:          "ai_explain",
		Title:       "AI Selection Explainer",
		Description: "Explain clipboard text selection with notification",
		Category:    "Voice & AI",
		Shortcut:    "<Super><Alt>c",
		Command:     "~/.local/bin/vula ai explain",
	},

	// Window & Tiling Management
	{
		ID:          "tactile_grid",
		Title:       "Tactile Window Grid",
		Description: "Open lettered grid matrix to place & resize windows",
		Category:    "Window Management",
		Shortcut:    "<Super>t",
		Command:     "gnome-shell extension (tactile)",
	},
	{
		ID:          "tiling_toggle",
		Title:       "Toggle Auto-Tiling",
		Description: "Toggle automatic window tiling gaps & layout",
		Category:    "Window Management",
		Shortcut:    "<Super><Alt>t",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_left_half",
		Title:       "Snap Window Left Half",
		Description: "Snap window to left 50% screen column",
		Category:    "Window Management",
		Shortcut:    "<Super>Left",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_right_half",
		Title:       "Snap Window Right Half",
		Description: "Snap window to right 50% screen column",
		Category:    "Window Management",
		Shortcut:    "<Super>Right",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_top_half",
		Title:       "Snap Top / Maximize",
		Description: "Snap window to top half or maximize screen space",
		Category:    "Window Management",
		Shortcut:    "<Super>Up",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_bottom_half",
		Title:       "Snap Bottom / Restore",
		Description: "Snap window to bottom half or restore size",
		Category:    "Window Management",
		Shortcut:    "<Super>Down",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_quarter_topleft",
		Title:       "Snap Top-Left Quarter",
		Description: "Snap window to top-left 25% quadrant",
		Category:    "Window Management",
		Shortcut:    "<Super><Alt>u",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_quarter_topright",
		Title:       "Snap Top-Right Quarter",
		Description: "Snap window to top-right 25% quadrant",
		Category:    "Window Management",
		Shortcut:    "<Super><Alt>i",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_quarter_botleft",
		Title:       "Snap Bottom-Left Quarter",
		Description: "Snap window to bottom-left 25% quadrant",
		Category:    "Window Management",
		Shortcut:    "<Super><Alt>j",
		Command:     "gnome-shell extension (tiling-assistant)",
	},
	{
		ID:          "tile_quarter_botright",
		Title:       "Snap Bottom-Right Quarter",
		Description: "Snap window to bottom-right 25% quadrant",
		Category:    "Window Management",
		Shortcut:    "<Super><Alt>k",
		Command:     "gnome-shell extension (tiling-assistant)",
	},

	// Window State Controls
	{
		ID:          "window_close",
		Title:       "Close Active Window",
		Description: "Close currently focused desktop window",
		Category:    "Window Controls",
		Shortcut:    "<Super>q",
		Command:     "org.gnome.desktop.wm.keybindings close",
	},
	{
		ID:          "window_minimize",
		Title:       "Minimize Active Window",
		Description: "Minimize active window to panel",
		Category:    "Window Controls",
		Shortcut:    "<Super>h",
		Command:     "org.gnome.desktop.wm.keybindings minimize",
	},
	{
		ID:          "window_maximize",
		Title:       "Maximize / Restore Window",
		Description: "Toggle fullscreen maximized state",
		Category:    "Window Controls",
		Shortcut:    "<Super>m",
		Command:     "org.gnome.desktop.wm.keybindings toggle-maximized",
	},
	{
		ID:          "show_desktop",
		Title:       "Show Desktop",
		Description: "Hide all windows to view clean desktop",
		Category:    "Window Controls",
		Shortcut:    "<Super>d",
		Command:     "org.gnome.desktop.wm.keybindings show-desktop",
	},
	{
		ID:          "window_drag",
		Title:       "Drag Window Anywhere",
		Description: "Hold Super and click anywhere inside window to drag",
		Category:    "Window Controls",
		Shortcut:    "<Super>Click",
		Command:     "org.gnome.desktop.wm.preferences mouse-button-modifier",
	},

	// Workspaces & Navigation
	{
		ID:          "workspace_switch",
		Title:       "Switch Workspace 1-6",
		Description: "Directly jump to workspace 1 through 6",
		Category:    "Workspaces & Navigation",
		Shortcut:    "<Super>1..6",
		Command:     "org.gnome.desktop.wm.keybindings switch-to-workspace-N",
	},
	{
		ID:          "workspace_move_window",
		Title:       "Move Window to Workspace 1-6",
		Description: "Move focused window to workspace 1 through 6",
		Category:    "Workspaces & Navigation",
		Shortcut:    "<Shift><Super>1..6",
		Command:     "org.gnome.desktop.wm.keybindings move-to-workspace-N",
	},
	{
		ID:          "workspace_navigate",
		Title:       "Navigate Workspaces",
		Description: "Switch to adjacent workspace (Vim HJKL or Arrows)",
		Category:    "Workspaces & Navigation",
		Shortcut:    "<Alt><Super>Arrows",
		Command:     "org.gnome.desktop.wm.keybindings switch-to-workspace-dir",
	},

	// Dock Quick Jump
	{
		ID:          "dock_jump",
		Title:       "Dock App Quick Jump 1-9",
		Description: "Launch or switch to pinned dock app position 1 through 9",
		Category:    "Dock Quick Jump",
		Shortcut:    "<Alt>1..9",
		Command:     "org.gnome.shell.keybindings app-shift-N",
	},
}

// FormatKeycap transforms <Super><Alt>v into [Super] + [Alt] + [V]
func FormatKeycap(shortcut string) string {
	if shortcut == "" {
		return "[Unset]"
	}
	s := shortcut
	s = strings.ReplaceAll(s, "<Super>", "[Super] ")
	s = strings.ReplaceAll(s, "<Alt>", "[Alt] ")
	s = strings.ReplaceAll(s, "<Shift>", "[Shift] ")
	s = strings.ReplaceAll(s, "<Ctrl>", "[Ctrl] ")
	s = strings.ReplaceAll(s, "<Primary>", "[Ctrl] ")
	
	parts := strings.Fields(s)
	for i, p := range parts {
		if !strings.HasPrefix(p, "[") {
			parts[i] = "[" + strings.ToUpper(p) + "]"
		}
	}
	return strings.Join(parts, " + ")
}

// GetKeybindings reads active GNOME shortcuts or returns defaults
func GetKeybindings() []Keybinding {
	kbList := make([]Keybinding, len(DefaultKeybindings))
	copy(kbList, DefaultKeybindings)

	// Attempt reading current dconf custom keybinding values
	for i, kb := range kbList {
		path := fmt.Sprintf("/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom-vula-%s/", kb.ID)
		out, err := gnome.GetDconfKeyPath(path, "binding")
		if err == nil && out != "" {
			val := strings.Trim(strings.TrimSpace(out), "'")
			if val != "" {
				kbList[i].Shortcut = val
			}
		}
	}

	return kbList
}

// SetKeybinding updates a keybinding shortcut in GNOME dconf
func SetKeybinding(cfg *config.Config, id, newShortcut string) error {
	var target *Keybinding
	for _, kb := range DefaultKeybindings {
		if kb.ID == id {
			t := kb
			target = &t
			break
		}
	}

	if target == nil {
		return fmt.Errorf("unknown keybinding ID '%s'", id)
	}

	target.Shortcut = newShortcut

	home := os.Getenv("HOME")
	vulaBin := filepath.Join(home, ".local", "bin", "vula")

	cmdStr := target.Command
	cmdStr = strings.ReplaceAll(cmdStr, "~/.local/bin/vula", vulaBin)
	if target.ID == "hud" {
		cmdStr = filepath.Join(home, ".local", "bin", "vula-hud-launch")
	} else if target.ID == "keys_map" {
		cmdStr = fmt.Sprintf("gnome-terminal --title=\"Vula Keybindings\" -- %s keys", vulaBin)
	}

	return gnome.RegisterCustomKeybinding(target.ID, target.Title, cmdStr, newShortcut)
}

// RenderTable outputs a clean terminal table of all keybindings
func RenderTable() string {
	b := strings.Builder{}
	b.WriteString(ui.RenderHeader("Vula Keybindings Map", "Global desktop hotkeys registered in GNOME"))
	b.WriteString("\n")

	kbList := GetKeybindings()

	hdrStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.SecondaryColor)
	b.WriteString(hdrStyle.Render(fmt.Sprintf("  %-26s %-28s %s\n", "SHORTCUT", "ACTION", "DESCRIPTION")))
	b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render("  "+strings.Repeat("─", 84)+"\n"))

	keyStyle := lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
	titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(ui.MutedColor)

	for _, kb := range kbList {
		keycap := FormatKeycap(kb.Shortcut)
		line := fmt.Sprintf("  %-26s %-28s %s\n", keyStyle.Render(keycap), titleStyle.Render(kb.Title), descStyle.Render(kb.Description))
		b.WriteString(line)
	}

	b.WriteString("\n")
	b.WriteString(ui.CardStyle.Render("Tip: Press [Super + Shift + K] anytime or run 'vula keys' to interactively edit hotkeys."))
	b.WriteString("\n")

	return b.String()
}

// Model handles the interactive Bubbletea Keybinding Visualizer and Editor
type Model struct {
	cfg         *config.Config
	bindings    []Keybinding
	selectedIdx int
	editing     bool
	input       textinput.Model
	statusMsg   string
	quitting    bool
}

func InitialModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "e.g. <Super><Shift>x or <Alt>F1"
	ti.CharLimit = 64
	ti.Width = 30
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)

	return Model{
		cfg:         cfg,
		bindings:    GetKeybindings(),
		selectedIdx: 0,
		editing:     false,
		input:       ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.input.SetValue("")
				m.statusMsg = ""
				return m, nil

			case "enter":
				val := strings.TrimSpace(m.input.Value())
				if val != "" {
					target := m.bindings[m.selectedIdx]
					if err := SetKeybinding(m.cfg, target.ID, val); err != nil {
						m.statusMsg = fmt.Sprintf("Error: %v", err)
					} else {
						m.bindings[m.selectedIdx].Shortcut = val
						m.statusMsg = fmt.Sprintf("Updated %s -> %s", target.Title, FormatKeycap(val))
					}
				}
				m.editing = false
				m.input.SetValue("")
				return m, nil
			}

			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "j":
			if m.selectedIdx < len(m.bindings)-1 {
				m.selectedIdx++
			}
			return m, nil

		case "enter", "e":
			m.editing = true
			m.input.SetValue(m.bindings[m.selectedIdx].Shortcut)
			m.input.Focus()
			m.statusMsg = "Type new shortcut format (e.g. <Super><Alt>x)"
			return m, textinput.Blink
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(ui.RenderHeader("Vula Keybindings Studio", "Interactive hotkey visualizer & editor"))
	b.WriteString("\n")

	if m.statusMsg != "" {
		b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.InfoStyle.Render("ℹ"), m.statusMsg))
	}

	for i, kb := range m.bindings {
		cursor := "  "
		keyStyle := lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
		titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor)

		if i == m.selectedIdx {
			cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ ")
			keyStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
			titleStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
		}

		keycap := FormatKeycap(kb.Shortcut)
		line := fmt.Sprintf("%s%-26s %-28s %s", cursor, keyStyle.Render(keycap), titleStyle.Render(kb.Title), lipgloss.NewStyle().Foreground(ui.MutedColor).Render(kb.Category))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.editing {
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("Edit Shortcut: "))
		b.WriteString(m.input.View())
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("Press Enter to save, Esc to cancel"))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Render(" [↑/↓] Select  •  [Enter/E] Edit Shortcut  •  [Q/Esc] Quit"))
	}
	b.WriteString("\n")

	return ui.CardStyle.Width(86).Render(b.String()) + "\n"
}

func RunInteractiveKeymap(cfg *config.Config) error {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
