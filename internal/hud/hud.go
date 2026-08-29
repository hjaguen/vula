package hud

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/ui"
	"github.com/vula-os/vula/internal/voice"
)

type Mode int

const (
	ModeCommands Mode = iota
	ModeAI
	ModeVoice
)

type ActionItem struct {
	Title       string
	Description string
	Category    string
	Command     string
	IsAI        bool
}

type Model struct {
	cfg         *config.Config
	aiClient    *ai.Client
	voiceEngine *voice.Engine
	input       textinput.Model
	spinner     spinner.Model
	mode        Mode
	actions     []ActionItem
	filtered    []ActionItem
	selectedIdx int
	aiResponse  strings.Builder
	loading     bool
	recording   bool
	statusMsg   string
	width       int
	height      int
	quitting    bool
}

type aiChunkMsg string
type aiDoneMsg struct{ err error }
type voiceDoneMsg struct {
	text string
	err  error
}

func InitialModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "Search actions, apps, or type '?' for AI..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(ui.TextLightColor)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(ui.SecondaryColor)

	allActions := []ActionItem{
		{Title: "Vula Doctor", Description: "Run full system & hardware diagnostics", Category: "System", Command: "doctor"},
		{Title: "Open Terminal", Description: "Launch modern terminal (Ghostty/Kitty)", Category: "Apps", Command: "terminal"},
		{Title: "Ask Vula AI", Description: "Ask local LLM with active desktop context", Category: "AI", IsAI: true},
		{Title: "Voice Dictation", Description: "Transcribe voice directly into active window", Category: "Voice", Command: "voice_dictate"},
		{Title: "Toggle Tiling Mode", Description: "Enable/disable automatic window tiling", Category: "Desktop", Command: "toggle_tiling"},
		{Title: "Vula Settings", Description: "Open Vula declarative configuration", Category: "System", Command: "config"},
		{Title: "Lock Screen", Description: "Lock the current session safely", Category: "System", Command: "lock"},
	}

	return Model{
		cfg:         cfg,
		aiClient:    ai.NewClient(cfg),
		voiceEngine: voice.NewEngine(cfg),
		input:       ti,
		spinner:     sp,
		mode:        ModeCommands,
		actions:     allActions,
		filtered:    allActions,
		selectedIdx: 0,
		width:       72,
		height:      20,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > 76 {
			m.input.Width = 66
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			// Cycle modes: Commands -> AI -> Voice
			if m.mode == ModeCommands {
				m.mode = ModeAI
				m.input.Placeholder = "Ask Vula AI anything (uses active desktop context)..."
			} else if m.mode == ModeAI {
				m.mode = ModeVoice
				m.input.Placeholder = "Voice mode: Press Enter to dictate..."
			} else {
				m.mode = ModeCommands
				m.input.Placeholder = "Search actions, apps, or type '?' for AI..."
			}
			return m, nil

		case "up", "ctrl+p":
			if m.mode == ModeCommands && m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "ctrl+n":
			if m.mode == ModeCommands && m.selectedIdx < len(m.filtered)-1 {
				m.selectedIdx++
			}
			return m, nil

		case "enter":
			val := strings.TrimSpace(m.input.Value())

			// Check for AI prompt prefix "?"
			if strings.HasPrefix(val, "?") || m.mode == ModeAI {
				query := strings.TrimPrefix(val, "?")
				query = strings.TrimSpace(query)
				if query == "" {
					return m, nil
				}
				m.loading = true
				m.mode = ModeAI
				m.aiResponse.Reset()
				m.statusMsg = "Thinking..."

				return m, func() tea.Msg {
					ctx := context.Background()
					_, err := m.aiClient.Ask(ctx, query, func(chunk string) {
						// Streaming updates
					})
					return aiDoneMsg{err: err}
				}
			}

			// Run selected action
			if m.mode == ModeCommands && len(m.filtered) > 0 {
				action := m.filtered[m.selectedIdx]
				if action.IsAI {
					m.mode = ModeAI
					m.input.SetValue("")
					m.input.Placeholder = "Ask Vula AI anything..."
					return m, nil
				}
				return m, m.executeAction(action)
			}
		}

	case aiChunkMsg:
		m.aiResponse.WriteString(string(msg))

	case aiDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.aiResponse.WriteString(fmt.Sprintf("\n[Error: %v]", msg.err))
		}
		m.statusMsg = "Done"

	case voiceDoneMsg:
		m.recording = false
		m.loading = false
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Voice error: %v", msg.err)
		} else {
			m.input.SetValue(msg.text)
			m.statusMsg = "Transcribed speech!"
		}
	}

	// Update text input and filter actions
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	if m.mode == ModeCommands {
		val := strings.ToLower(m.input.Value())
		m.filtered = nil
		for _, a := range m.actions {
			if strings.Contains(strings.ToLower(a.Title), val) || strings.Contains(strings.ToLower(a.Description), val) {
				m.filtered = append(m.filtered, a)
			}
		}
		if m.selectedIdx >= len(m.filtered) {
			m.selectedIdx = 0
		}
	}

	if m.loading {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) executeAction(action ActionItem) tea.Cmd {
	return func() tea.Msg {
		switch action.Command {
		case "terminal":
			_ = exec.Command("gnome-terminal").Start()
		case "lock":
			_ = exec.Command("loginctl", "lock-session").Run()
		case "toggle_tiling":
			_ = exec.Command("vula", "extensions", "toggle-tiling").Run()
		}
		return tea.Quit()
	}
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	b := strings.Builder{}

	// Header Bar
	headerLeft := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.PrimaryColor).
		Render(" ⚡ VULA HUD ")

	modeBadge := ""
	switch m.mode {
	case ModeCommands:
		modeBadge = ui.SuccessBadge.Render(" ACTIONS ")
	case ModeAI:
		modeBadge = ui.BadgeStyle.Render(" AI ASSIST ")
	case ModeVoice:
		modeBadge = ui.WarnBadge.Render(" VOICE ")
	}

	topBar := lipgloss.JoinHorizontal(lipgloss.Center, headerLeft, " ", modeBadge)
	b.WriteString(topBar)
	b.WriteString("\n\n")

	// Search Input Box
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// Content area depending on mode
	switch m.mode {
	case ModeCommands:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 68)))
		b.WriteString("\n")
		if len(m.filtered) == 0 {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  No matching actions found. Type '?' to ask Vula AI.\n"))
		} else {
			for i, a := range m.filtered {
				cursor := "  "
				titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor)
				descStyle := lipgloss.NewStyle().Foreground(ui.MutedColor)

				if i == m.selectedIdx {
					cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ ")
					titleStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
					descStyle = lipgloss.NewStyle().Foreground(ui.TextLightColor)
				}

				line := fmt.Sprintf("%s%-24s %s", cursor, titleStyle.Render(a.Title), descStyle.Render(a.Description))
				b.WriteString(line)
				b.WriteString("\n")
			}
		}

	case ModeAI:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 68)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Consulting local AI with desktop context...")))
		}
		if m.aiResponse.Len() > 0 {
			content := m.aiResponse.String()
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Padding(0, 1).Render(content))
			b.WriteString("\n")
		} else if !m.loading {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  Type your question and press Enter to query Ollama/Cloud AI.\n"))
		}

	case ModeVoice:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 68)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎙 Voice Subsystem Active\n"))
		b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Render("  Press Enter to record 5 seconds of audio and transcribe with Whisper.\n"))
	}

	// Footer Hints
	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Render("  [Tab] Switch Mode  •  [Enter] Select  •  [Esc] Close  •  [?] AI Mode")
	b.WriteString(footer)

	// Wrap in HUD Card
	return ui.CardStyle.Width(74).Render(b.String()) + "\n"
}

func RunHUD(cfg *config.Config) error {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
