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
	Title    string
	Icon     string
	Category string
	Command  string
	IsAI     bool
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
	ti.Placeholder = "Search or '?' for AI..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 28
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(ui.TextLightColor)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(ui.SecondaryColor)

	themeMgr := theme.NewManager(cfg)
	themesList := themeMgr.ListThemes()

	allActions := []ActionItem{
		{Title: "Voice AI Assistant", Icon: "⚡", Category: "Voice", Command: "voice_assistant"},
		{Title: "Voice Listen AI", Icon: "🎧", Category: "Voice", Command: "voice_listen"},
		{Title: "App Store TUI", Icon: "📦", Category: "Apps", Command: "apps_store"},
		{Title: "Vula Doctor", Icon: "🩺", Category: "System", Command: "doctor"},
		{Title: "Switch Theme", Icon: "🎨", Category: "Theme", Command: "theme"},
		{Title: "Open Terminal", Icon: "💻", Category: "Apps", Command: "terminal"},
		{Title: "Ask Vula AI", Icon: "🤖", Category: "AI", IsAI: true},
		{Title: "Voice Dictation", Icon: "🎙", Category: "Voice", Command: "voice_dictate"},
		{Title: "Synthesize Voice", Icon: "🔊", Category: "Voice", Command: "voice_speak"},
		{Title: "Vula Settings", Icon: "⚙", Category: "System", Command: "config"},
		{Title: "Lock Screen", Icon: "🔒", Category: "System", Command: "lock"},
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
		width:       38,
		height:      22,
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
		if m.width > 38 {
			m.input.Width = 28
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.mode != ModeCommands {
				m.mode = ModeCommands
				m.input.SetValue("")
				m.input.Placeholder = "Search or '?' for AI..."
				m.statusMsg = ""
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "tab":
			if m.mode == ModeCommands {
				m.mode = ModeAI
				m.input.Placeholder = "Ask Vula AI..."
			} else if m.mode == ModeAI {
				m.mode = ModeVoice
				m.input.Placeholder = "Enter to record..."
			} else if m.mode == ModeVoice {
				m.mode = ModeTheme
				m.input.Placeholder = "Select theme..."
			} else {
				m.mode = ModeCommands
				m.input.Placeholder = "Search or '?' for AI..."
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

			if m.mode == ModeTheme && len(m.themesList) > 0 {
				selectedTheme := m.themesList[m.themeIdx]
				if err := m.themeMgr.ApplyTheme(selectedTheme.Name); err != nil {
					m.statusMsg = fmt.Sprintf("Error: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Theme: %s", selectedTheme.DisplayName)
				}
				return m, nil
			}

			if m.mode == ModeVoice {
				m.loading = true
				m.recording = true
				m.statusMsg = "Recording 4s..."
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

			if m.mode == ModeCommands && len(m.filtered) > 0 {
				action := m.filtered[m.selectedIdx]
				if action.IsAI {
					m.mode = ModeAI
					m.input.SetValue("")
					m.input.Placeholder = "Ask Vula AI..."
					return m, nil
				}
				return m.handleActionSelection(action)
			}
		}

	case aiChunkMsg:
		m.loading = false
		m.aiResponse.Reset()
		m.aiResponse.WriteString(string(msg))
		m.statusMsg = "Done"

	case aiDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.aiResponse.WriteString(fmt.Sprintf("\n[Error: %v]", msg.err))
		}
		m.statusMsg = "Done"

	case doctorDoneMsg:
		m.loading = false
		m.doctorOutput = msg.output
		m.statusMsg = "Complete"

	case voiceDoneMsg:
		m.recording = false
		m.loading = false
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.statusMsg = fmt.Sprintf("Typed: %s", msg.text)
			_ = voice.TypeIntoActiveWindow(msg.text)
		}
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	if m.mode == ModeCommands {
		val := strings.ToLower(m.input.Value())
		m.filtered = nil
		for _, a := range m.actions {
			if strings.Contains(strings.ToLower(a.Title), val) {
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
	case "apps_store":
		go func() {
			home := os.Getenv("HOME")
			vulaBin := filepath.Join(home, ".local", "bin", "vula")
			cmd := exec.Command("gnome-terminal", "--title=Vula App Store", "--", vulaBin, "apps", "ui")
			_ = cmd.Start()
		}()
		m.quitting = true
		return *m, tea.Quit

	case "voice_listen":
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

	case "doctor":
		m.mode = ModeDoctor
		m.loading = true
		m.statusMsg = "Running..."
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
		m.statusMsg = "Select theme"
		return *m, nil

	case "terminal":
		go launchTerminal()
		m.quitting = true
		return *m, tea.Quit

	case "voice_assistant":
		m.mode = ModeAI
		m.loading = true
		m.statusMsg = "Escuchando..."
		return *m, func() tea.Msg {
			assistant := voice.NewAssistant(m.cfg)
			q, ans, err := assistant.ListenAndRespond(context.Background(), 4)
			if err != nil {
				return aiDoneMsg{err: err}
			}
			return aiChunkMsg(fmt.Sprintf("🎙 Tú: %s\n\n⚡ AI: %s", q, ans))
		}

	case "voice_dictate":
		m.mode = ModeVoice
		m.statusMsg = "Enter to record"
		return *m, nil

	case "voice_speak":
		m.statusMsg = "Speaking..."
		go func() {
			_ = m.voiceEngine.Speak(context.Background(), "Hola Mauricio, el motor de voz de Vula está activo.")
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
		Render("⚡ VULA HUD")

	modeBadge := ""
	switch m.mode {
	case ModeCommands:
		modeBadge = ui.SuccessBadge.Render("ACTIONS")
	case ModeAI:
		modeBadge = ui.BadgeStyle.Render("AI")
	case ModeDoctor:
		modeBadge = ui.InfoStyle.Render("DOCTOR")
	case ModeVoice:
		modeBadge = ui.WarnBadge.Render("VOICE")
	case ModeConfig:
		modeBadge = ui.BadgeStyle.Render("CONFIG")
	case ModeTheme:
		modeBadge = ui.SuccessBadge.Render("THEME")
	}

	topBar := lipgloss.JoinHorizontal(lipgloss.Center, headerLeft, "  ", modeBadge)
	b.WriteString(topBar)
	b.WriteString("\n\n")

	if m.mode == ModeCommands || m.mode == ModeAI || m.mode == ModeVoice || m.mode == ModeTheme {
		b.WriteString(m.input.View())
		b.WriteString("\n\n")
	}

	switch m.mode {
	case ModeCommands:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		if len(m.filtered) == 0 {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  No actions found.\n"))
		} else {
			for i, a := range m.filtered {
				cursor := "  "
				iconStyle := lipgloss.NewStyle().Foreground(ui.PrimaryColor)
				titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor)

				if i == m.selectedIdx {
					cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ ")
					iconStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
					titleStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
				}

				line := fmt.Sprintf("%s%s  %s", cursor, iconStyle.Render(a.Icon), titleStyle.Render(a.Title))
				b.WriteString(line)
				b.WriteString("\n")
			}
		}

	case ModeDoctor:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Running...")))
		} else {
			b.WriteString(m.doctorOutput)
		}

	case ModeAI:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Consulting AI...")))
		}
		if m.aiResponse.Len() > 0 {
			content := m.aiResponse.String()
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Padding(0, 1).Render(content))
			b.WriteString("\n")
		} else if !m.loading {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  Type & press Enter\n"))
		}

	case ModeVoice:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎙 Voice Subsystem\n\n"))
		if m.recording {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), ui.WarnStyle.Render("Recording...")))
		} else if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render("  Enter to record.\n"))
		}

	case ModeTheme:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎨 Theme Selector\n\n"))
		if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
		}
		for i, t := range m.themesList {
			cursor := "  "
			titleStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor)
			accentBlock := lipgloss.NewStyle().Foreground(lipgloss.Color(t.AccentColor)).Bold(true).Render("■")

			activeTag := ""
			if t.Name == m.cfg.Theme.Palette {
				activeTag = " " + ui.SuccessBadge.Render("✓")
			}

			if i == m.themeIdx {
				cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ ")
				titleStyle = lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true)
			}

			line := fmt.Sprintf("%s%s %-12s%s", cursor, accentBlock, titleStyle.Render(t.DisplayName), activeTag)
			b.WriteString(line)
			b.WriteString("\n")
		}

	case ModeConfig:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 32)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render(m.configOutput))
	}

	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Render(" [Tab] Mode  •  [Esc] Close")
	b.WriteString(footer)

	return ui.CardStyle.Width(38).Render(b.String()) + "\n"
}

func RunHUD(cfg *config.Config) error {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
