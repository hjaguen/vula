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
	"github.com/vula-os/vula/internal/actions"
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
	ModeDo
)

type ActionItem struct {
	Title       string
	Icon        string
	Category    string
	Command     string
	Description string
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
type appStoreDoneMsg struct{ err error }
type keysDoneMsg struct{ err error }
type voiceDoneMsg struct {
	text string
	err  error
}

func InitialModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "Buscar, 'do ...' para acción, '?' para IA..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 38
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(ui.TextLightColor)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(ui.SecondaryColor)

	themeMgr := theme.NewManager(cfg)
	themesList := themeMgr.ListThemes()

	allActions := []ActionItem{
		{Title: "Execute OS Action", Icon: "⚡", Category: "Sistema", Description: "Ejecutar comando del terminal o script por voz/texto", Command: "os_action"},
		{Title: "Workstation Profiles", Icon: "🚀", Category: "Sistema", Description: "Configurar perfil (Full-Stack, DevOps, Design, DBA, Minimal)", Command: "profile_selector"},
		{Title: "Keybindings Map & Shortcuts", Icon: "⌨️", Category: "Configuración", Description: "Ver y editar atajos de teclado del escritorio", Command: "keys_map"},
		{Title: "Tactile Grid Matrix (Super+T)", Icon: "🔲", Category: "Tiling", Description: "Malla interactiva de letras para mover y redimensionar", Command: "tactile_grid"},
		{Title: "Toggle Auto-Tiling (Super+Alt+T)", Icon: "📐", Category: "Tiling", Description: "Alternar distribución automática de ventanas", Command: "tiling_toggle"},
		{Title: "Voice AI Assistant", Icon: "⚡", Category: "Asistente IA", Description: "Conversar con el asistente de voz Vula", Command: "voice_assistant"},
		{Title: "Voice Listen AI", Icon: "🎧", Category: "Asistente IA", Description: "Escuchar instrucción por voz y responder", Command: "voice_listen"},
		{Title: "App Store TUI", Icon: "📦", Category: "Tienda", Description: "Instalar y gestionar aplicaciones y herramientas dev", Command: "apps_store"},
		{Title: "Developer Font Picker", Icon: "🔤", Category: "Sistema", Description: "Seleccionar fuente de programación Nerd Font", Command: "font_picker"},
		{Title: "WebApps Manager", Icon: "🌐", Category: "Apps", Description: "Crear y gestionar aplicaciones web de escritorio", Command: "webapp_mgr"},
		{Title: "Vula Doctor Diagnostics", Icon: "🩺", Category: "Sistema", Description: "Diagnóstico de salud de sonido, IA y escritorio", Command: "doctor"},
		{Title: "Switch Theme (19 Palettes)", Icon: "🎨", Category: "Apariencia", Description: "Cambiar entre temas visuales (Oscuro/Claro/etc.)", Command: "theme"},
		{Title: "Open Terminal", Icon: "💻", Category: "Herramientas", Description: "Abrir emulador de terminal moderno", Command: "terminal"},
		{Title: "Ask Vula AI", Icon: "🤖", Category: "IA", Description: "Chatear directamente con el motor Vula AI", IsAI: true},
		{Title: "Voice Dictation", Icon: "🎙", Category: "Voz & IA", Description: "Dictar por voz e insertar texto en ventana activa", Command: "voice_dictate"},
		{Title: "Synthesize Voice", Icon: "🔊", Category: "Voz & IA", Description: "Probar síntesis de voz neural Piper TTS", Command: "voice_speak"},
		{Title: "Vula Settings", Icon: "⚙", Category: "Sistema", Description: "Ver y modificar configuración de Vula", Command: "config"},
		{Title: "Lock Screen", Icon: "🔒", Category: "Sistema", Description: "Bloquear sesión de pantalla", Command: "lock"},
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
				m.input.Placeholder = "Buscar, 'do ...' para acción, '?' para IA..."
				m.statusMsg = ""
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "tab":
			if m.mode == ModeCommands {
				m.mode = ModeDo
				m.input.SetValue("do ")
				m.input.Placeholder = "Acción de sistema (ej: 'ajusta el brillo al 40%')..."
			} else if m.mode == ModeDo {
				m.mode = ModeAI
				m.input.SetValue("")
				m.input.Placeholder = "Pregunta a Vula AI..."
			} else if m.mode == ModeAI {
				m.mode = ModeVoice
				m.input.SetValue("")
				m.input.Placeholder = "Enter para grabar por voz..."
			} else if m.mode == ModeVoice {
				m.mode = ModeTheme
				m.input.SetValue("")
				m.input.Placeholder = "Selecciona un tema..."
			} else {
				m.mode = ModeCommands
				m.input.SetValue("")
				m.input.Placeholder = "Buscar, 'do ...' para acción, '?' para IA..."
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

			if strings.HasPrefix(strings.ToLower(val), "do ") || m.mode == ModeDo {
				query := strings.TrimPrefix(val, "do ")
				query = strings.TrimPrefix(query, "DO ")
				query = strings.TrimSpace(query)
				if query == "" {
					return m, nil
				}
				m.loading = true
				m.mode = ModeAI
				m.aiResponse.Reset()
				m.statusMsg = "Executing OS action..."

				return m, func() tea.Msg {
					ctx := context.Background()
					plan, err := actions.ParseIntent(ctx, m.aiClient, query)
					if err != nil {
						return aiDoneMsg{err: err}
					}
					if err := actions.ExecutePlan(ctx, plan, m.cfg); err != nil {
						return aiDoneMsg{err: err}
					}
					return aiChunkMsg(fmt.Sprintf("✓ OS Action Executed:\n\nInstrucción: %s\nResumen: %s", plan.OriginalPrompt, plan.Summary))
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

	case appStoreDoneMsg:
		m.mode = ModeCommands
		m.statusMsg = "App Store closed"

	case keysDoneMsg:
		m.mode = ModeCommands
		m.statusMsg = "Keybindings closed"

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
	case "os_action":
		m.mode = ModeDo
		m.input.SetValue("do ")
		m.input.Placeholder = "Escribe tu acción (ej: 'ajusta el brillo al 40%')..."
		return *m, nil

	case "apps_store":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		c := exec.Command(vulaBin, "apps", "ui")
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			return appStoreDoneMsg{err: err}
		})

	case "keys_map":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		c := exec.Command(vulaBin, "keys")
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			return keysDoneMsg{err: err}
		})

	case "tactile_grid":
		_ = exec.Command("gdbus", "call", "--session", "--dest", "org.gnome.Shell", "--object-path", "/org/gnome/Shell", "--method", "org.gnome.Shell.Extensions.Tactile.ShowTiles").Run()
		m.quitting = true
		return *m, tea.Quit

	case "tiling_toggle":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		_ = exec.Command(vulaBin, "desktop", "toggle-tile").Run()
		m.statusMsg = "Toggled Auto-Tiling"
		return *m, nil

	case "font_picker":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		c := exec.Command(vulaBin, "font", "list")
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			return keysDoneMsg{err: err}
		})

	case "webapp_mgr":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		c := exec.Command(vulaBin, "webapp", "list")
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			return keysDoneMsg{err: err}
		})

	case "profile_selector":
		home := os.Getenv("HOME")
		vulaBin := filepath.Join(home, ".local", "bin", "vula")
		c := exec.Command(vulaBin, "profile", "ui")
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			return keysDoneMsg{err: err}
		})

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
			return doctorDoneMsg{output: report.RenderCompact()}
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

	// 1. Centered Header Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.PrimaryColor).
		Width(52).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render("⚡ VULA HUD"))
	b.WriteString("\n\n")

	// 2. Sliding 3-Tab Carousel Row (3 spacious rounded-border badges shifting on Tab)
	type TabInfo struct {
		Label string
		Color lipgloss.Color
	}
	tabMap := map[Mode]TabInfo{
		ModeCommands: {Label: ">_ COMANDOS", Color: ui.PrimaryColor},
		ModeDo:       {Label: "⚡ ACCIÓN OS", Color: ui.SecondaryColor},
		ModeAI:       {Label: "🤖 IA CHAT", Color: ui.PrimaryColor},
		ModeVoice:    {Label: "🎙 VOZ", Color: ui.SecondaryColor},
		ModeTheme:    {Label: "🎨 TEMA", Color: ui.PrimaryColor},
	}

	var visibleModes []Mode
	switch m.mode {
	case ModeCommands:
		visibleModes = []Mode{ModeCommands, ModeDo, ModeAI}
	case ModeDo:
		visibleModes = []Mode{ModeCommands, ModeDo, ModeAI}
	case ModeAI:
		visibleModes = []Mode{ModeDo, ModeAI, ModeVoice}
	case ModeVoice:
		visibleModes = []Mode{ModeAI, ModeVoice, ModeTheme}
	case ModeTheme:
		visibleModes = []Mode{ModeVoice, ModeTheme, ModeCommands}
	default:
		visibleModes = []Mode{ModeCommands, ModeDo, ModeAI}
	}

	var renderedBadges []string
	for _, modeKey := range visibleModes {
		info := tabMap[modeKey]
		style := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder())
		if modeKey == m.mode {
			renderedBadges = append(renderedBadges, style.BorderForeground(info.Color).Foreground(info.Color).Bold(true).Render(info.Label))
		} else {
			renderedBadges = append(renderedBadges, style.BorderForeground(ui.MutedColor).Foreground(ui.MutedColor).Render(info.Label))
		}
	}

	badgeRow := lipgloss.JoinHorizontal(lipgloss.Center, renderedBadges[0], " ", renderedBadges[1], " ", renderedBadges[2])
	b.WriteString(lipgloss.NewStyle().Width(52).Align(lipgloss.Center).Render(badgeRow))
	b.WriteString("\n\n")

	// 3. Search Bar Container (always visible for visual consistency across all tabs)
	searchBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PrimaryColor).
		Padding(0, 1).
		Width(50).
		Render("🔍 " + m.input.View())
	b.WriteString(searchBox)
	b.WriteString("\n\n")

	switch m.mode {
	case ModeCommands:
		if len(m.filtered) == 0 {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  No se encontraron comandos.\n"))
		} else {
			maxVisible := 2
			start := 0
			if m.selectedIdx >= maxVisible {
				start = m.selectedIdx - maxVisible + 1
			}
			end := start + maxVisible
			if end > len(m.filtered) {
				end = len(m.filtered)
			}

			for i := start; i < end; i++ {
				a := m.filtered[i]
				if i == m.selectedIdx {
					cardStyle := lipgloss.NewStyle().
						Border(lipgloss.RoundedBorder()).
						BorderForeground(ui.PrimaryColor).
						Padding(0, 1).
						Width(50)

					iconBox := lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("▶ " + a.Icon)
					titleStr := lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render(a.Title)
					subStr := lipgloss.NewStyle().Foreground(ui.TextLightColor).Italic(true).Render(truncateText(fmt.Sprintf("%s • %s", a.Category, a.Description), 44))

					cardContent := fmt.Sprintf("%s  %s\n   %s", iconBox, titleStr, subStr)
					b.WriteString(cardStyle.Render(cardContent))
					b.WriteString("\n")
				} else {
					iconBox := lipgloss.NewStyle().Foreground(ui.MutedColor).Render("  " + a.Icon)
					titleStr := lipgloss.NewStyle().Foreground(ui.TextLightColor).Render(a.Title)
					subStr := lipgloss.NewStyle().Foreground(ui.MutedColor).Render(truncateText(fmt.Sprintf("%s • %s", a.Category, a.Description), 44))

					line1 := fmt.Sprintf("%s  %s", iconBox, titleStr)
					line2 := fmt.Sprintf("   %s", subStr)
					b.WriteString(line1 + "\n" + line2 + "\n\n")
				}
			}
		}

	case ModeDoctor:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 50)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Running...")))
		} else {
			b.WriteString(m.doctorOutput)
		}

	case ModeAI:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 50)))
		b.WriteString("\n")
		if m.loading {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), lipgloss.NewStyle().Foreground(ui.SecondaryColor).Render("Consulting AI...")))
		}
		if m.aiResponse.Len() > 0 {
			content := m.aiResponse.String()
			rendered := lipgloss.NewStyle().
				Foreground(ui.TextLightColor).
				Width(48).
				Render(content)
			b.WriteString(rendered)
			b.WriteString("\n")
		} else if !m.loading {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render("  Escribe tu consulta y presiona Enter\n"))
		}

	case ModeVoice:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 50)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎙 Subsistema de Voz\n\n"))
		if m.recording {
			b.WriteString(fmt.Sprintf("  %s %s\n", m.spinner.View(), ui.WarnStyle.Render("Grabando voz (4s)...")))
		} else if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render("  Presiona Enter para grabar audio.\n"))
		}

	case ModeTheme:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 50)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.SecondaryColor).Bold(true).Render("  🎨 Selector de Temas Estéticos (19 Paletas)\n\n"))
		if m.statusMsg != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n\n", ui.SuccessStyle.Render("✓"), m.statusMsg))
		}

		maxVisible := 4
		start := 0
		if m.themeIdx >= maxVisible {
			start = m.themeIdx - maxVisible + 1
		}
		end := start + maxVisible
		if end > len(m.themesList) {
			end = len(m.themesList)
		}

		for i := start; i < end; i++ {
			t := m.themesList[i]
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

			line := fmt.Sprintf("%s%s %-24s%s", cursor, accentBlock, titleStyle.Render(t.DisplayName), activeTag)
			b.WriteString(line)
			b.WriteString("\n")
		}

	case ModeConfig:
		b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render(strings.Repeat("─", 50)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ui.TextLightColor).Render(m.configOutput))
	}

	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Render(" [Tab] Cambiar Modo  •  [Esc] Cerrar HUD")
	b.WriteString(footer)

	return ui.CardStyle.Width(54).Render(b.String()) + "\n"
}

func truncateText(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func RunHUD(cfg *config.Config) error {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
