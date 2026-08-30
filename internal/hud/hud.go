package hud

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/doctor"
	"github.com/vula-os/vula/internal/theme"
	"github.com/vula-os/vula/internal/ui"
	"github.com/vula-os/vula/internal/voice"
	"gopkg.in/yaml.v3"
)

type Mode int

const (
	ModeCommands Mode = iota
	ModeAI
	ModeDoctor
	ModeVoice
	ModeConfig
	ModeTheme
)

type ActionItem struct {
	Title       string
	Description string
	Category    string
	Command     string
	IsAI        bool
}

type Model struct {
	cfg          *config.Config
	aiClient     *ai.Client
	voiceEngine  *voice.Engine
	themeMgr     *theme.Manager
	themesList   []theme.ThemePalette
	themeIdx     int
	input        textinput.Model
	spinner      spinner.Model
	mode         Mode
	actions      []ActionItem
	filtered     []ActionItem
	selectedIdx  int
	aiResponse   strings.Builder
	doctorOutput string
	configOutput string
	loading      bool
	recording    bool
	statusMsg    string
	width        int
	height       int
	quitting     bool
}

type aiChunkMsg string
type aiDoneMsg struct{ err error }
type doctorDoneMsg struct{ output string }
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

	themeMgr := theme.NewManager(cfg)
	themesList := themeMgr.ListThemes()

	allActions := []ActionItem{
		{Title: "⚡ Voice AI Assistant", Description: "Listen by voice, query local AI & speak answer back", Category: "Voice", Command: "voice_assistant"},
		{Title: "Vula Doctor", Description: "Run full system, audio, GPU & AI diagnostics", Category: "System", Command: "doctor"},
		{Title: "Switch System Theme", Description: "Change global system palette across GNOME & terminal", Category: "Theme", Command: "theme"},
		{Title: "Open Terminal", Description: "Launch modern terminal window (Ghostty/GNOME)", Category: "Apps", Command: "terminal"},
		{Title: "Ask Vula AI", Description: "Query local LLM with active desktop context", Category: "AI", IsAI: true},
		{Title: "Voice Dictation", Description: "Transcribe voice directly into focused window", Category: "Voice", Command: "voice_dictate"},
		{Title: "Synthesize Voice", Description: "Test Piper local neural text-to-speech", Category: "Voice", Command: "voice_speak"},
		{Title: "Vula Settings", Description: "Inspect declarative YAML configuration", Category: "System", Command: "config"},
		{Title: "Lock Screen", Description: "Lock the current session safely", Category: "System", Command: "lock"},
	}

	currentThemeIdx := 0
	for i, t := range themesList {
		if t.Name == cfg.Theme.Palette {
			currentThemeIdx = i
			break
		}
	}

	return Model{
		cfg:         cfg,
		aiClient:    ai.NewClient(cfg),
		voiceEngine: voice.NewEngine(cfg),
		themeMgr:    themeMgr,
		themesList:  themesList,
		themeIdx:    currentThemeIdx,
		input:       ti,
		spinner:     sp,
		mode:        ModeCommands,
		actions:     allActions,
		filtered:    allActions,
		selectedIdx: 0,
		width:       76,
		height:      24,
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
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.mode != ModeCommands {
				// Return to main commands list
				m.mode = ModeCommands
				m.input.SetValue("")
				m.input.Placeholder = "Search actions, apps, or type '?' for AI..."
				m.statusMsg = ""
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "tab":
			// Cycle modes: Commands -> AI -> Voice -> Theme -> Commands
			if m.mode == ModeCommands {
				m.mode = ModeAI
				m.input.Placeholder = "Ask Vula AI anything (uses active desktop context)..."
			} else if m.mode == ModeAI {
				m.mode = ModeVoice
				m.input.Placeholder = "Voice mode: Press Enter to dictate..."
			} else if m.mode == ModeVoice {
				m.mode = ModeTheme
				m.input.Placeholder = "Theme mode: Select theme and press Enter..."
			} else {
				m.mode = ModeCommands
				m.input.Placeholder = "Search actions, apps, or type '?' for AI..."
			}
			return m, nil

		case "up", "ctrl+p":
			if m.mode == ModeCommands && m.selectedIdx > 0 {
				m.selectedIdx--
			} else if m.mode == ModeTheme && m.themeIdx > 0 {
				m.themeIdx--
			}
			return m, nil

		case "down", "ctrl+n":
			if m.mode == ModeCommands && m.selectedIdx < len(m.filtered)-1 {
				m.selectedIdx++
			} else if m.mode == ModeTheme && m.themeIdx < len(m.themesList)-1 {
				m.themeIdx++
			}
			return m, nil

		case "enter":
			val := strings.TrimSpace(m.input.Value())

			// Theme Mode Enter: Apply selected theme
			if m.mode == ModeTheme && len(m.themesList) > 0 {
				selectedTheme := m.themesList[m.themeIdx]
				if err := m.themeMgr.ApplyTheme(selectedTheme.Name); err != nil {
					m.statusMsg = fmt.Sprintf("Theme error: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Theme switched to '%s' across GNOME & terminal!", selectedTheme.DisplayName)
				}
				return m, nil
			}

			// Voice Mode Enter: Record audio and transcribe
			if m.mode == ModeVoice {
				m.loading = true
				m.recording = true
				m.statusMsg = "Recording 4 seconds of audio..."
				return m, func() tea.Msg {
					tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("vula_rec_%d.wav", time.Now().UnixNano()))
					ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancel()

					if err := m.voiceEngine.RecordAudio(ctx, tempFile, 4); err != nil {
						return voiceDoneMsg{err: err}
					}
					text, err := m.voiceEngine.Transcribe(ctx, tempFile)
					_ = os.Remove(tempFile)
					return voiceDoneMsg{text: text, err: err}
				}
			}

			// AI Mode Enter: Query model
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
					resp, err := m.aiClient.Ask(ctx, query, nil)
					if err == nil {
						return aiChunkMsg(resp)
					}
					return aiDoneMsg{err: err}
				}
			}

			// Run selected action from commands list
			if m.mode == ModeCommands && len(m.filtered) > 0 {
				action := m.filtered[m.selectedIdx]
				if action.IsAI {
					m.mode = ModeAI
					m.input.SetValue("")
					m.input.Placeholder = "Ask Vula AI anything..."
					return m, nil
				}
				return m.handleActionSelection(action)
			}
		}

	case aiChunkMsg:
		m.loading = false
		m.aiResponse.Reset()
		m.aiResponse.WriteString(string(msg))
		m.statusMsg = "Answer generated"

	case aiDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.aiResponse.WriteString(fmt.Sprintf("\n[Error: %v]", msg.err))
		}
		m.statusMsg = "Done"

	case doctorDoneMsg:
		m.loading = false
		m.doctorOutput = msg.output
		m.statusMsg = "Diagnostics complete"

	case voiceDoneMsg:
		m.recording = false
		m.loading = false
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Voice error: %v", msg.err)
		} else {
			m.statusMsg = fmt.Sprintf("Transcribed: %s", msg.text)
			_ = voice.TypeIntoActiveWindow(msg.text)
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

func (m *Model) handleActionSelection(action ActionItem) (Model, tea.Cmd) {
	switch action.Command {
	case "doctor":
		m.mode = ModeDoctor
		m.loading = true
		m.statusMsg = "Running system diagnostics..."
		return *m, func() tea.Msg {
			report := doctor.RunDiagnostics(m.cfg)
			return doctorDoneMsg{output: report.Render()}
		}

	case "theme":
		m.mode = ModeTheme
		m.themesList = m.themeMgr.ListThemes()
		for i, t := range m.themesList {
			if t.Name == m.cfg.Theme.Palette {
				m.themeIdx = i
				break
			}
		}
		m.statusMsg = "Select a theme and press Enter to apply."
		return *m, nil

	case "terminal":
		// Launch detached terminal so HUD can exit or stay responsive
		go launchTerminal()
		m.quitting = true
		return *m, tea.Quit

	case "voice_assistant":
		m.mode = ModeAI
		m.loading = true
		m.statusMsg = "Escuchando tu voz..."
		return *m, func() tea.Msg {
			assistant := voice.NewAssistant(m.cfg)
			q, ans, err := assistant.ListenAndRespond(context.Background(), 4)
			if err != nil {
				return aiDoneMsg{err: err}
			}
			return aiChunkMsg(fmt.Sprintf("🎙 Tú: %s\n\n⚡ Vula AI: %s", q, ans))
		}

	case "voice_dictate":
		m.mode = ModeVoice
		m.statusMsg = "Voice dictation mode ready. Press Enter to record."
		return *m, nil

	case "voice_speak":
		m.statusMsg = "Speaking test phrase with Piper..."
		go func() {
			_ = m.voiceEngine.Speak(context.Background(), "Hola Mauricio, el motor de voz de Vula está activo y funcionando.")
		}()
		return *m, nil

	case "config":
		m.mode = ModeConfig
		data, _ := yaml.Marshal(m.cfg)
		m.configOutput = string(data)
		return *m, nil

	case "lock":
		_ = exec.Command("loginctl", "lock-session").Run()
		m.quitting = true
		return *m, tea.Quit
	}

	return *m, nil
}

func launchTerminal() {
	// Try preferred terminal emulators in order
	terminals := []string{"gnome-terminal", "ghostty", "kitty", "alacritty", "xterm"}
	for _, term := range terminals {
		if path, err := exec.LookPath(term); err == nil {
			cmd := exec.Command(path)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			_ = cmd.Start()
			return
		}
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
	case ModeDoctor:
		modeBadge = ui.InfoStyle.Render(" DOCTOR ")
	case ModeVoice:
		modeBadge = ui.WarnBadge.Render(" VOICE ")
	case ModeConfig:
		modeBadge = ui.BadgeStyle.Render(" CONFIG ")
	case ModeTheme:
		modeBadge = ui.SuccessBadge.Render(" THEME ")
	}

	topBar := lipgloss.JoinHorizontal(lipgloss.Center, headerLeft, " ", modeBadge)
	b.WriteString(topBar)
	b.WriteString("\n\n")

	// Search Input Box (only shown in Commands, AI, Voice, or Theme modes)
	if m.mode == ModeCommands || m.mode == ModeAI || m.mode == ModeVoice || m.mode == ModeTheme {
		b.WriteString(m.input.View())
		b.WriteString("\n\n")
	}

	// Content area depending on mode
	switch m.mode {
	case ModeCommands:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
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

	case ModeDoctor:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Running full hardware and OS diagnostics...")))
		} else {
			b.WriteString(m.doctorOutput)
		}

	case ModeAI:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Consulting local Ollama with desktop context...")))
		}
		if m.aiResponse.Len() > 0 {
			content := m.aiResponse.String()
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Padding(0, 1).Render(content))
			b.WriteString("\n")
		} else if !m.loading {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  Type your question and press Enter to query Ollama.\n"))
		}

	case ModeVoice:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎙 Voice Subsystem Active (Whisper + Piper)\n\n"))
		if m.recording {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), ui.WarnStyle.Render("Recording microphone audio... Speak now!")))
		} else if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Render("  Press Enter to record another voice sample.\n"))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render("  Press [Enter] to record 4s of voice and type into active window.\n"))
		}

	case ModeTheme:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎨 System Theme Selector (GNOME & Terminal)\n\n"))
		if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
		}
		for i, t := range m.themesList {
			cursor := "  "
			titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor)
			accentBlock := lipgloss.NewStyle().Foreground(lipgloss.Color(t.AccentColor)).Bold(true).Render("■■")
			bgHex := lipgloss.NewStyle().Foreground(ui.MutedColor).Render(fmt.Sprintf("[%s]", t.Background))

			activeTag := ""
			if t.Name == m.cfg.Theme.Palette {
				activeTag = "  " + ui.SuccessBadge.Render(" ACTIVE ")
			}

			if i == m.themeIdx {
				cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ ")
				titleStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
			}

			line := fmt.Sprintf("%s%s %-18s %s%s", cursor, accentBlock, titleStyle.Render(t.DisplayName), bgHex, activeTag)
			b.WriteString(line)
			b.WriteString("\n")
		}

	case ModeConfig:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 70)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render(m.configOutput))
	}

	// Footer Hints
	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Render("  [Tab] Switch Mode  •  [Enter] Select  •  [Esc] Back/Close  •  [?] AI")
	b.WriteString(footer)

	// Wrap in HUD Card
	return ui.CardStyle.Width(76).Render(b.String()) + "\n"
}

func RunHUD(cfg *config.Config) error {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
