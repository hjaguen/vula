package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/vula-os/vula/internal/actions"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/apps"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/doctor"
	"github.com/vula-os/vula/internal/dotfiles"
	"github.com/vula-os/vula/internal/font"
	"github.com/vula-os/vula/internal/gnome"
	"github.com/vula-os/vula/internal/hud"
	"github.com/vula-os/vula/internal/installer"
	"github.com/vula-os/vula/internal/keys"
	"github.com/vula-os/vula/internal/media"
	"github.com/vula-os/vula/internal/theme"
	"github.com/vula-os/vula/internal/ui"
	"github.com/vula-os/vula/internal/voice"
)

var rootCmd = &cobra.Command{
	Use:   "vula",
	Short: "Vula — Next-Gen AI & Voice Developer OS for Ubuntu 24.04 LTS",
	Long:  ui.RenderBanner(),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		// Default action when running 'vula' without subcommands is launching the HUD
		if err := hud.RunHUD(cfg); err != nil {
			log.Error("HUD exited with error", "error", err)
		}
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run full system diagnostic checks for Ubuntu, GNOME, Audio, AI and tools",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		report := doctor.RunDiagnostics(cfg)
		fmt.Println(report.Render())
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run interactive installer to configure Vula modules on Ubuntu",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		if err := installer.RunInteractiveInstaller(cfg); err != nil {
			log.Error("Installation failed", "error", err)
			os.Exit(1)
		}
	},
}

var hudCmd = &cobra.Command{
	Use:   "hud",
	Short: "Launch the floating Raycast-style TUI launcher & AI assistant",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		if err := hud.RunHUD(cfg); err != nil {
			log.Error("HUD error", "error", err)
			os.Exit(1)
		}
	},
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Interact with the local AI intelligence engine",
}

var aiAskCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Ask a question to the AI assistant with active desktop context",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		client := ai.NewClient(cfg)
		prompt := args[0]

		fmt.Printf("\n%s\n\n", ui.InfoStyle.Render("⚡ Consulting Vula AI..."))
		_, err := client.Ask(context.Background(), prompt, func(chunk string) {
			fmt.Print(chunk)
		})
		fmt.Println()
		if err != nil {
			fmt.Printf("\n%s %v\n", ui.ErrorStyle.Render("Error:"), err)
		}
	},
}

var aiCmdSuggest = &cobra.Command{
	Use:   "cmd [task description]",
	Short: "Translate natural language to a shell command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		client := ai.NewClient(cfg)
		task := args[0]

		suggested, err := client.SuggestCommand(context.Background(), task)
		if err != nil {
			log.Error("Failed to generate command", "error", err)
			os.Exit(1)
		}
		fmt.Printf("\n%s\n  %s\n\n", ui.SubtitleStyle.Render("Suggested command:"), ui.SuccessStyle.Render(suggested))
	},
}

var aiModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List installed Ollama local models",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		client := ai.NewClient(cfg)
		models, err := client.ListLocalModels(context.Background())
		if err != nil {
			log.Error("Failed to list models", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.RenderHeader("Local AI Models", fmt.Sprintf("Found %d models in Ollama", len(models))))
		for _, m := range models {
			fmt.Printf("  • %-24s (Size: %.2f GB)\n", ui.InfoStyle.Render(m.Name), float64(m.Size)/(1024*1024*1024))
		}
		fmt.Println()
	},
}

var voiceCmd = &cobra.Command{
	Use:   "voice",
	Short: "Voice subsystem commands (dictation, speech synthesis, wake word)",
}

var voiceRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record speech and type transcription into active window",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		engine := voice.NewEngine(cfg)

		tempFile := "/tmp/vula_voice_record.wav"
		fmt.Println(ui.InfoStyle.Render("🎙 Recording audio (speak now)..."))
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := engine.RecordAudio(ctx, tempFile, 4); err != nil {
			log.Error("Recording failed", "error", err)
			os.Exit(1)
		}

		fmt.Println(ui.InfoStyle.Render("Transcribing with Whisper..."))
		text, err := engine.Transcribe(context.Background(), tempFile)
		if err != nil {
			log.Error("Transcription failed", "error", err)
			os.Exit(1)
		}

		fmt.Printf("\n%s %s\n", ui.SuccessStyle.Render("Result:"), text)
		_ = voice.TypeIntoActiveWindow(text)
	},
}

var voiceSpeakCmd = &cobra.Command{
	Use:   "speak [text]",
	Short: "Synthesize text into speech locally using Piper TTS",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		engine := voice.NewEngine(cfg)
		text := args[0]
		fmt.Printf("%s Synthesizing speech with Piper...\n", ui.InfoStyle.Render("⚡"))
		if err := engine.Speak(context.Background(), text); err != nil {
			log.Error("Speech synthesis failed", "error", err)
			os.Exit(1)
		}
	},
}

var voiceTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test default microphone audio recording level and diagnostics",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		engine := voice.NewEngine(cfg)

		fmt.Println(ui.InfoStyle.Render("🎙 Probando micrófono durante 3 segundos... (¡Habla ahora!)"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		maxVol, err := engine.TestMicrophone(ctx)
		if err != nil {
			log.Error("Prueba de micrófono fallida", "error", err)
			os.Exit(1)
		}

		if maxVol > -45.0 {
			fmt.Printf("\n%s Se detectó entrada de audio activa (Nivel máximo: %.1f dB)\n", ui.SuccessStyle.Render("✓"), maxVol)
			fmt.Println("  El micrófono está capturando audio correctamente.")
		} else {
			fmt.Printf("\n%s Se detectó silencio o nivel de audio muy bajo (Nivel máximo: %.1f dB)\n", ui.WarnStyle.Render("⚠️"), maxVol)
			fmt.Println("  Revisa que el micrófono interno no tenga silenciador activado o revisa el volumen en Configuración -> Sonido.")
		}
	},
}

var listenCmd = &cobra.Command{
	Use:   "listen [seconds]",
	Short: "Start active listening AI assistant (Whisper STT -> Ollama -> Piper TTS)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		assistant := voice.NewAssistant(cfg)
		duration := 4
		if len(args) > 0 {
			if d, err := strconv.Atoi(args[0]); err == nil && d > 0 {
				duration = d
			}
		}
		fmt.Printf("%s Vula AI Escuchando durante %d segundos... (¡Habla ahora al micrófono!)\n", ui.InfoStyle.Render("⚡"), duration)
		q, ans, err := assistant.ListenAndRespond(context.Background(), duration)
		if err != nil {
			log.Error("Escucha de voz finalizada", "info", err)
			os.Exit(1)
		}
		fmt.Printf("\n%s %s\n", ui.SubtitleStyle.Render("Voz detectada:"), q)
		fmt.Printf("%s %s\n\n", ui.SuccessStyle.Render("Respuesta Vula:"), ans)
	},
}

var desktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Configure GNOME Shell desktop environment and shortcuts",
}

var desktopSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Apply GNOME Shell optimizations, themes, and developer keybindings",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := gnome.NewManager(cfg)
		if err := mgr.ApplyDesktopOptimizations(); err != nil {
			log.Error("Desktop configuration failed", "error", err)
			os.Exit(1)
		}
		_ = mgr.ConfigureTilingKeybindings()
		fmt.Println(ui.SuccessStyle.Render("✓ GNOME Shell optimizations and global keybindings applied successfully!"))
		fmt.Println("  • [Super + Space] -> Vula HUD")
		fmt.Println("  • [Super + Alt + V] -> Voice Dictation")
		fmt.Println("  • [Super + Alt + A] -> Active Voice AI Assistant")
		fmt.Println("  • [Super + T] -> Toggle Auto-Tile mode")
	},
}

var desktopTilingCmd = &cobra.Command{
	Use:   "tiling",
	Short: "Configure GNOME 46 Tiling Assistant (gaps, quarter snaps, active border hint)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := gnome.NewManager(cfg)
		if err := mgr.ConfigureTilingKeybindings(); err != nil {
			log.Error("Tiling configuration failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.SuccessStyle.Render("✓ Tiling Assistant configured with custom gaps and active border highlight!"))
		fmt.Println("  • [Super + Left/Right/Up/Down] -> Half-screen snap")
		fmt.Println("  • [Super + Alt + U/I/J/K] -> Quarter-screen snap")
		fmt.Println("  • [Super + T] -> Toggle Auto-Tile mode")
	},
}

var desktopExtensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: "Manage and install curated GNOME Shell 46 visual extensions",
}

var extensionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List curated GNOME extensions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.RenderHeader("Curated GNOME Extensions", "Visual & functional enhancements for GNOME 46"))
		for _, e := range gnome.CuratedExtensions {
			fmt.Printf("  • %-36s %s\n", ui.InfoStyle.Render(e.Name), e.Description)
		}
		fmt.Println()
	},
}

var extensionsInstallCuratedCmd = &cobra.Command{
	Use:   "install-curated",
	Short: "Download and install curated visual extensions (Blur my Shell, Just Perfection)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := gnome.NewManager(cfg)
		fmt.Println(ui.InfoStyle.Render("⚡ Installing curated GNOME Shell 46 extensions..."))
		if err := mgr.InstallCuratedExtensions(); err != nil {
			log.Error("Extension installation failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.SuccessStyle.Render("✓ Curated GNOME extensions installed and enabled!"))
	},
}

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage and switch global system and terminal themes",
}

var themeSetCmd = &cobra.Command{
	Use:   "set [theme-name]",
	Short: "Apply a global theme (built-in or custom)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := theme.NewManager(cfg)
		themeName := args[0]
		if err := mgr.ApplyTheme(themeName); err != nil {
			log.Error("Failed to apply theme", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Theme switched to %s across GNOME and desktop!\n", ui.SuccessStyle.Render("✓"), themeName)
	},
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available global themes (built-in and custom)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := theme.NewManager(cfg)
		fmt.Println(ui.RenderHeader("Available Themes", "Unified aesthetic palettes for Vula"))
		for _, t := range mgr.ListThemes() {
			tag := ""
			if t.IsCustom {
				tag = " [custom]"
			}
			fmt.Printf("  • %-18s %-22s %s (Accent: %s)\n", ui.InfoStyle.Render(t.Name), t.DisplayName+tag, t.Background, t.AccentColor)
		}
		fmt.Println()
	},
}

var themeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactive TUI theme creator (Charm Huh form with live preview)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		fmt.Println(ui.RenderHeader("Vula Theme Studio", "Create a cohesive system-wide color palette"))
		palette, err := theme.RunInteractiveCreator(cfg)
		if err != nil {
			log.Error("Theme creation cancelled or failed", "error", err)
			os.Exit(1)
		}
		fmt.Println("\n" + theme.RenderPreview(*palette))
		fmt.Printf("\n%s Custom theme '%s' saved to ~/.config/vula/themes/%s.yaml and applied!\n\n",
			ui.SuccessStyle.Render("✓"), palette.DisplayName, palette.Name)
	},
}

var themeGenerateCmd = &cobra.Command{
	Use:   "generate [aesthetic-description]",
	Short: "Generate an accessible, cohesive palette using local AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		prompt := strings.Join(args, " ")
		fmt.Printf("%s Generating AI theme for '%s'...\n", ui.InfoStyle.Render("⚡"), prompt)
		palette, err := theme.GenerateAITheme(context.Background(), cfg, prompt)
		if err != nil {
			log.Error("AI theme generation failed", "error", err)
			os.Exit(1)
		}
		fmt.Println("\n" + theme.RenderPreview(*palette))
		fmt.Printf("\n%s AI Theme '%s' created and applied!\n\n", ui.SuccessStyle.Render("✓"), palette.DisplayName)
	},
}

var themePreviewCmd = &cobra.Command{
	Use:   "preview [theme-name]",
	Short: "Render terminal preview card for a theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := theme.NewManager(cfg)
		all := mgr.GetAllThemes()
		t, exists := all[args[0]]
		if !exists {
			log.Error("Theme not found", "name", args[0])
			os.Exit(1)
		}
		fmt.Println("\n" + theme.RenderPreview(t) + "\n")
	},
}

var dotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Manage curated developer dotfiles (Fish, Starship, Neovim, Tmux)",
}

var dotfilesInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and link Vula developer dotfiles into ~/.config",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := dotfiles.NewManager(cfg)
		if err := mgr.InstallAllDotfiles(); err != nil {
			log.Error("Failed to install dotfiles", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.SuccessStyle.Render("✓ Vula developer dotfiles installed successfully!"))
		fmt.Println("  • Fish Shell: ~/.config/fish/config.fish")
		fmt.Println("  • Starship Prompt: ~/.config/starship.toml")
		fmt.Println("  • Neovim: ~/.config/nvim/init.lua")
		fmt.Println("  • Tmux: ~/.tmux.conf")
	},
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage and install curated developer applications and CLI tooling",
}

var appsInstallCLICmd = &cobra.Command{
	Use:   "install-cli",
	Short: "Install the complete modern developer CLI stack (eza, bat, lazygit, starship, btop, fzf)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := apps.NewManager(cfg)
		fmt.Println(ui.InfoStyle.Render("⚡ Installing modern developer CLI toolchain..."))
		if err := mgr.InstallCLIStack(); err != nil {
			log.Error("CLI installation error", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.SuccessStyle.Render("✓ Developer CLI stack installed successfully!"))
	},
}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List curated applications available in Vula",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.RenderHeader("Developer App Recipes", "Curated developer software catalog"))
		for _, a := range apps.Catalog {
			fmt.Printf("  • %-14s [%-8s] %s\n", ui.InfoStyle.Render(a.ID), a.Category, a.Description)
		}
		fmt.Println()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Vula version and build info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Vula %s (Ubuntu 24.04 LTS Noble | Go 1.24+ | Charm Engine)\n", ui.Version)
	},
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Display system metrics, hardware specs, and Vula status card",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		fmt.Println(ui.RenderFetchCard(cfg))
	},
}

var aiCommitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"git-commit"},
	Short:   "Generate conventional git commit message from git diff",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		fmt.Printf("\n%s\n\n", ui.InfoStyle.Render("⚡ Analyzing git diff with local AI..."))
		msg, err := ai.GenerateCommitMessage(context.Background(), cfg)
		if err != nil {
			log.Error("Failed generating commit message", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n  %s\n\n", ui.SubtitleStyle.Render("Suggested commit message:"), ui.SuccessStyle.Render(msg))

		fmt.Print("Commit with this message? [Y/n]: ")
		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" || input == "y" || input == "yes" {
			if err := ai.ExecuteGitCommit(msg); err != nil {
				log.Error("Git commit failed", "error", err)
				os.Exit(1)
			}
			fmt.Println(ui.SuccessStyle.Render("✓ Git commit created successfully!"))
		} else {
			fmt.Println(ui.WarnStyle.Render("Commit cancelled."))
		}
	},
}

var aiFixCmd = &cobra.Command{
	Use:   "fix [error text]",
	Short: "Diagnose terminal error and suggest immediate corrective command",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		inputText := ""
		if len(args) > 0 {
			inputText = strings.Join(args, " ")
		}
		fmt.Printf("\n%s\n\n", ui.InfoStyle.Render("⚡ Diagnosing terminal error..."))
		diag, err := ai.DiagnoseTerminalError(context.Background(), cfg, inputText)
		if err != nil {
			log.Error("Diagnosis failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.RenderHeader("Terminal Error Diagnosis", "AI Root Cause & Fix Recommendation"))
		fmt.Println(diag)
		fmt.Println()
	},
}

var aiExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain or refactor selected text/code from active clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		fmt.Printf("\n%s\n\n", ui.InfoStyle.Render("⚡ Explaining selection..."))
		exp, err := ai.ExplainSelection(context.Background(), cfg)
		if err != nil {
			log.Error("Explanation failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.RenderHeader("AI Code & Context Explanation", "Active Clipboard Analysis"))
		fmt.Println(exp)
		fmt.Println()
	},
}

var voiceDaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Launch continuous hands-free voice assistant background daemon",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		fmt.Println(ui.InfoStyle.Render("⚡ Starting Vula Hands-Free Voice Daemon..."))
		if err := voice.StartVoiceDaemon(context.Background(), cfg); err != nil {
			log.Error("Voice daemon failed", "error", err)
			os.Exit(1)
		}
	},
}

var appsUICmd = &cobra.Command{
	Use:     "ui",
	Aliases: []string{"store"},
	Short:   "Interactive visual TUI developer app store (Charm Huh form)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		mgr := apps.NewManager(cfg)
		if err := mgr.RunInteractiveAppStore(); err != nil {
			log.Error("App store failed", "error", err)
			os.Exit(1)
		}
	},
}

var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Manage and rotate high-definition theme wallpapers",
}

var wallpaperNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Rotate to the next theme wallpaper automatically",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		wm := theme.NewWallpaperManager(cfg)
		path, err := wm.RotateNextWallpaper()
		if err != nil {
			log.Error("Wallpaper rotation failed", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Switched wallpaper to %s\n", ui.SuccessStyle.Render("✓"), filepath.Base(path))
	},
}

var wallpaperSetCmd = &cobra.Command{
	Use:   "set [path-or-filename]",
	Short: "Apply a specific wallpaper file to GNOME desktop",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		wm := theme.NewWallpaperManager(cfg)
		target := args[0]

		if !filepath.IsAbs(target) {
			target = filepath.Join(wm.GetWallpaperDir(), target)
		}

		if err := wm.SetWallpaper(target); err != nil {
			log.Error("Failed to set wallpaper", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Set desktop wallpaper: %s\n", ui.SuccessStyle.Render("✓"), filepath.Base(target))
	},
}

var wallpaperListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all downloaded and generated wallpapers in ~/.config/vula/wallpapers",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		wm := theme.NewWallpaperManager(cfg)
		files, err := wm.ListWallpapers()
		if err != nil {
			log.Error("Failed listing wallpapers", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.RenderHeader("Wallpaper Collection", fmt.Sprintf("Found %d wallpapers in ~/.config/vula/wallpapers", len(files))))
		for _, f := range files {
			fmt.Printf("  • %-28s (%s)\n", ui.InfoStyle.Render(filepath.Base(f)), f)
		}
		fmt.Println()
	},
}

var wallpaperFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Download curated 4K aesthetic wallpapers matching Vula theme palettes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		wm := theme.NewWallpaperManager(cfg)
		fmt.Printf("%s Downloading curated theme wallpapers...\n", ui.InfoStyle.Render("⚡"))
		if err := wm.FetchPresetWallpapers(); err != nil {
			log.Error("Wallpaper fetch failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.SuccessStyle.Render("✓ Curated HD wallpapers downloaded to ~/.config/vula/wallpapers!"))
	},
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Interactive TUI visualizer & customizer for Vula desktop hotkeys",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		if err := keys.RunInteractiveKeymap(cfg); err != nil {
			log.Error("Keybindings editor error", "error", err)
			os.Exit(1)
		}
	},
}

var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered desktop hotkeys in terminal table",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(keys.RenderTable())
	},
}

var keysSetCmd = &cobra.Command{
	Use:   "set [id] [shortcut]",
	Short: "Update keybinding shortcut in GNOME dconf (e.g. vula keys set hud '<Super>space')",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		id := args[0]
		shortcut := args[1]
		if err := keys.SetKeybinding(cfg, id, shortcut); err != nil {
			log.Error("Failed to update keybinding", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Updated keybinding '%s' to %s in GNOME!\n", ui.SuccessStyle.Render("✓"), id, keys.FormatKeycap(shortcut))
	},
}

var aiStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show AI provider health, current mode (hybrid/local/cloud), and active models",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		client := ai.NewClient(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		localModels, err := client.ListLocalModels(ctx)
		localOk := err == nil

		fmt.Println(ui.RenderHeader("Vula AI Provider Status", "Hybrid AI Engine Diagnostic"))
		fmt.Printf("  • Execution Mode:        %s\n", ui.InfoStyle.Render(cfg.AI.Mode))
		fmt.Printf("  • Heavy Token Threshold: %d tokens\n\n", cfg.AI.HeavyTokenThreshold)

		fmt.Println(ui.SubtitleStyle.Render("Local AI (Ollama):"))
		fmt.Printf("  • Host:                 %s\n", cfg.AI.OllamaHost)
		fmt.Printf("  • Default Model:        %s\n", cfg.AI.DefaultModel)
		if localOk {
			fmt.Printf("  • Status:               %s (%d local models installed)\n", ui.SuccessStyle.Render("✓ Online"), len(localModels))
		} else {
			fmt.Printf("  • Status:               %s (Ollama not running or unreachable)\n", ui.ErrorStyle.Render("✗ Offline"))
		}

		fmt.Println("\n" + ui.SubtitleStyle.Render("Cloud AI Provider:"))
		fmt.Printf("  • Active Provider:      %s\n", cfg.AI.CloudProvider)
		fmt.Printf("  • Cloud Model:          %s\n", cfg.AI.CloudModel)

		keyName := cfg.AI.CloudProvider
		keyVal := cfg.AI.APIKeys[keyName]
		if keyVal == "" {
			switch keyName {
			case "gemini":
				keyVal = os.Getenv("GEMINI_API_KEY")
			case "groq":
				keyVal = os.Getenv("GROQ_API_KEY")
			case "openai":
				keyVal = os.Getenv("OPENAI_API_KEY")
			}
		}

		if keyVal != "" {
			masked := keyVal
			if len(keyVal) > 8 {
				masked = keyVal[:4] + "..." + keyVal[len(keyVal)-4:]
			}
			valErr := client.ValidateActiveCloudKey(ctx)
			if valErr == nil {
				fmt.Printf("  • API Key:              %s (%s - Valid & Reachable)\n", ui.SuccessStyle.Render("✓ Configured"), masked)
			} else {
				fmt.Printf("  • API Key:              %s (%s - %v)\n", ui.WarnStyle.Render("⚠️ Invalid/Reachable issue"), masked, valErr)
			}
		} else {
			fmt.Printf("  • API Key:              %s (Set with 'vula ai config --key=%s:YOUR_KEY')\n", ui.WarnStyle.Render("⚠️ Missing"), keyName)
		}
		fmt.Println()
	},
}

var (
	cfgMode       string
	cfgProvider   string
	cfgKey        string
	cfgModel      string
	cfgCloudModel string
	cfgThreshold  int
)

var aiConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI execution mode, cloud providers, API keys, and models",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		updated := false

		if cfgMode != "" {
			cfg.AI.Mode = cfgMode
			updated = true
		}
		if cfgProvider != "" {
			cfg.AI.CloudProvider = cfgProvider
			updated = true
		}
		if cfgModel != "" {
			cfg.AI.DefaultModel = cfgModel
			updated = true
		}
		if cfgCloudModel != "" {
			cfg.AI.CloudModel = cfgCloudModel
			updated = true
		}
		if cfgThreshold > 0 {
			cfg.AI.HeavyTokenThreshold = cfgThreshold
			updated = true
		}
		var keyProviderName, keyProviderVal string
		if cfgKey != "" {
			parts := strings.SplitN(cfgKey, ":", 2)
			if len(parts) == 2 {
				if cfg.AI.APIKeys == nil {
					cfg.AI.APIKeys = make(map[string]string)
				}
				keyProviderName = parts[0]
				keyProviderVal = parts[1]
				cfg.AI.APIKeys[keyProviderName] = keyProviderVal
				updated = true
			} else {
				log.Error("Invalid key format. Use --key=provider:API_KEY (e.g. --key=groq:gsk_...)")
				os.Exit(1)
			}
		}

		if !updated {
			fmt.Println(ui.WarnStyle.Render("No configuration flags passed."))
			fmt.Println("Usage examples:")
			fmt.Println("  vula ai config --mode=hybrid --provider=groq")
			fmt.Println("  vula ai config --key=groq:gsk_YourGroqKeyHere")
			fmt.Println("  vula ai config --cloud-model=llama-3.3-70b-versatile")
			return
		}

		if err := config.SaveConfig(cfg); err != nil {
			log.Error("Failed to save config", "error", err)
			os.Exit(1)
		}

		fmt.Println(ui.SuccessStyle.Render("✓ Vula AI configuration updated and saved successfully!"))

		// Perform live validation check if a key was configured
		if keyProviderVal != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			client := ai.NewClient(cfg)
			fmt.Printf("%s Validating %s API key live with remote endpoint...\n", ui.InfoStyle.Render("⚡"), keyProviderName)
			if err := client.ValidateActiveCloudKey(ctx); err != nil {
				fmt.Printf("  %s Key saved, but live validation check failed: %v\n", ui.WarnStyle.Render("⚠️ Warning:"), err)
			} else {
				fmt.Printf("  %s Live validation succeeded! '%s' API key is active and ready.\n", ui.SuccessStyle.Render("✓"), keyProviderName)
			}
		}
	},
}

var doCmd = &cobra.Command{
	Use:   "do [instruction]",
	Short: "Execute OS control actions using natural language AI intent parser",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		client := ai.NewClient(cfg)
		userPrompt := strings.Join(args, " ")

		fmt.Printf("\n%s Parsing OS action intent...\n\n", ui.InfoStyle.Render("⚡"))
		plan, err := actions.ParseIntent(context.Background(), client, userPrompt)
		if err != nil {
			log.Error("Failed to parse OS action plan", "error", err)
			os.Exit(1)
		}

		if err := actions.ExecutePlan(context.Background(), plan, cfg); err != nil {
			log.Error("Failed executing OS action plan", "error", err)
			os.Exit(1)
		}
	},
}

var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Inspect OS actions history and audit logs",
}

var actionsHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show audit log entries for all OS actions executed by Vula",
	Run: func(cmd *cobra.Command, args []string) {
		logs, err := actions.GetAuditLogs()
		if err != nil {
			log.Error("Failed to read audit logs", "error", err)
			os.Exit(1)
		}

		fmt.Println(ui.RenderHeader("OS Actions Audit Log", fmt.Sprintf("Found %d audit entries in ~/.config/vula/audit.log", len(logs))))
		if len(logs) == 0 {
			fmt.Println(ui.WarnStyle.Render("No audit log entries found."))
			return
		}

		for _, entry := range logs {
			timeStr := entry.Timestamp.Format("2006-01-02 15:04:05")
			statusStyle := ui.SuccessStyle
			if entry.ExecutedStatus != "SUCCESS" {
				statusStyle = ui.WarnStyle
			}
			fmt.Printf("  • [%s] %-9s %-12s | Prompt: \"%s\" | %s\n",
				timeStr,
				ui.InfoStyle.Render(string(entry.Risk)),
				statusStyle.Render(entry.ExecutedStatus),
				entry.UserPrompt,
				entry.Action.Description,
			)
		}
		fmt.Println()
	},
}

var webappCmd = &cobra.Command{
	Use:   "webapp",
	Short: "Create and manage desktop launcher entries for web applications",
}

var webappAddCmd = &cobra.Command{
	Use:   "add [Name] [URL] [IconURL]",
	Short: "Create a desktop launcher for a web URL with auto-fetched icon",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		targetURL := args[1]
		customIcon := ""
		if len(args) > 2 {
			customIcon = args[2]
		}
		info, err := apps.AddWebApp(name, targetURL, customIcon)
		if err != nil {
			log.Error("Failed to add web app", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Web app '%s' created successfully!\n", ui.SuccessStyle.Render("✓"), info.Name)
		fmt.Printf("  • URL:          %s\n", info.URL)
		fmt.Printf("  • Desktop File: %s\n", info.DesktopPath)
		fmt.Printf("  • Icon Path:    %s\n", info.IconPath)
	},
}

var webappRemoveCmd = &cobra.Command{
	Use:   "remove [Name]",
	Short: "Remove a Vula-managed web app launcher",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := apps.RemoveWebApp(name); err != nil {
			log.Error("Failed to remove web app", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Web app '%s' launcher removed.\n", ui.SuccessStyle.Render("✓"), name)
	},
}

var webappListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Vula-created web application launchers",
	Run: func(cmd *cobra.Command, args []string) {
		webApps, err := apps.ListWebApps()
		if err != nil {
			log.Error("Failed to list web apps", "error", err)
			os.Exit(1)
		}
		fmt.Println(ui.RenderHeader("Vula Web Applications", fmt.Sprintf("Found %d web apps in ~/.local/share/applications", len(webApps))))
		for _, app := range webApps {
			fmt.Printf("  • %-20s %s\n", ui.InfoStyle.Render(app.Name), app.URL)
		}
		fmt.Println()
	},
}

var fontCmd = &cobra.Command{
	Use:   "font",
	Short: "Manage and synchronize developer Nerd Fonts across Ghostty, Alacritty, VS Code, and GNOME",
}

var fontSetCmd = &cobra.Command{
	Use:   "set [cascadia|firacode|jetbrains|meslo]",
	Short: "Download and set developer font family",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := font.NewManager()
		if err != nil {
			log.Error("Font manager init failed", "error", err)
			os.Exit(1)
		}
		recipe, err := mgr.SetFont(args[0])
		if err != nil {
			log.Error("Failed setting font", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Font set to %s (%s) across Ghostty, Alacritty, VS Code & GNOME!\n",
			ui.SuccessStyle.Render("✓"), recipe.Name, recipe.FontFamily)
	},
}

var fontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available developer Nerd Fonts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.RenderHeader("Developer Nerd Fonts", "Curated coding fonts with icon glyphs"))
		for _, f := range font.AvailableFonts {
			fmt.Printf("  • %-12s %-20s (%s)\n", ui.InfoStyle.Render(f.ID), f.Name, f.FontFamily)
		}
		fmt.Println()
	},
}

var fontSizeCmd = &cobra.Command{
	Use:   "size [pt]",
	Short: "Update font size across Ghostty, Alacritty, and VS Code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		size, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error("Invalid font size. Expected integer pt value (e.g. 11)")
			os.Exit(1)
		}
		mgr, _ := font.NewManager()
		if err := mgr.SetFontSize(size); err != nil {
			log.Error("Failed setting font size", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Font size set to %d pt across developer tools!\n", ui.SuccessStyle.Render("✓"), size)
	},
}

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Media processing and screen recording conversion utilities",
}

var mediaWebm2Mp4Cmd = &cobra.Command{
	Use:   "webm2mp4 [filepath]",
	Short: "Convert GNOME screen recording .webm file to web-friendly .mp4 using ffmpeg",
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := ""
		if len(args) > 0 {
			inputPath = args[0]
		}
		fmt.Printf("%s Converting .webm screen recording to .mp4...\n", ui.InfoStyle.Render("⚡"))
		outPath, err := media.ConvertWebmToMp4(inputPath, "")
		if err != nil {
			log.Error("Media conversion failed", "error", err)
			os.Exit(1)
		}
		fmt.Printf("%s Conversion complete: %s\n", ui.SuccessStyle.Render("✓"), outPath)
	},
}

func init() {
	aiConfigCmd.Flags().StringVar(&cfgMode, "mode", "", "AI mode: hybrid, local, cloud")
	aiConfigCmd.Flags().StringVar(&cfgProvider, "provider", "", "Cloud provider: gemini, groq, ollama-cloud, openai")
	aiConfigCmd.Flags().StringVar(&cfgKey, "key", "", "API key in format provider:KEY (e.g. gemini:YOUR_KEY)")
	aiConfigCmd.Flags().StringVar(&cfgModel, "model", "", "Default local model (e.g. qwen2.5-coder:1.5b)")
	aiConfigCmd.Flags().StringVar(&cfgCloudModel, "cloud-model", "", "Default cloud model (e.g. gemini-2.0-flash)")
	aiConfigCmd.Flags().IntVar(&cfgThreshold, "threshold", 0, "Heavy token threshold for cloud offloading (e.g. 1200)")

	aiCmd.AddCommand(aiAskCmd)
	aiCmd.AddCommand(aiCmdSuggest)
	aiCmd.AddCommand(aiModelsCmd)
	aiCmd.AddCommand(aiCommitCmd)
	aiCmd.AddCommand(aiFixCmd)
	aiCmd.AddCommand(aiExplainCmd)
	aiCmd.AddCommand(aiStatusCmd)
	aiCmd.AddCommand(aiConfigCmd)

	actionsCmd.AddCommand(actionsHistoryCmd)

	voiceCmd.AddCommand(voiceRecordCmd)
	voiceCmd.AddCommand(voiceSpeakCmd)
	voiceCmd.AddCommand(voiceTestCmd)
	voiceCmd.AddCommand(voiceDaemonCmd)

	desktopCmd.AddCommand(desktopSetupCmd)
	desktopCmd.AddCommand(desktopTilingCmd)
	desktopCmd.AddCommand(desktopExtensionsCmd)
	desktopExtensionsCmd.AddCommand(extensionsListCmd)
	desktopExtensionsCmd.AddCommand(extensionsInstallCuratedCmd)

	themeCmd.AddCommand(themeSetCmd)
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeCreateCmd)
	themeCmd.AddCommand(themeGenerateCmd)
	themeCmd.AddCommand(themePreviewCmd)

	wallpaperCmd.AddCommand(wallpaperNextCmd)
	wallpaperCmd.AddCommand(wallpaperSetCmd)
	wallpaperCmd.AddCommand(wallpaperListCmd)
	wallpaperCmd.AddCommand(wallpaperFetchCmd)

	dotfilesCmd.AddCommand(dotfilesInstallCmd)

	appsCmd.AddCommand(appsInstallCLICmd)
	appsCmd.AddCommand(appsListCmd)
	appsCmd.AddCommand(appsUICmd)

	keysCmd.AddCommand(keysListCmd)
	keysCmd.AddCommand(keysSetCmd)

	webappCmd.AddCommand(webappAddCmd)
	webappCmd.AddCommand(webappRemoveCmd)
	webappCmd.AddCommand(webappListCmd)

	fontCmd.AddCommand(fontSetCmd)
	fontCmd.AddCommand(fontListCmd)
	fontCmd.AddCommand(fontSizeCmd)

	mediaCmd.AddCommand(mediaWebm2Mp4Cmd)

	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(actionsCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(hudCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(aiCmd)
	rootCmd.AddCommand(voiceCmd)
	rootCmd.AddCommand(desktopCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(wallpaperCmd)
	rootCmd.AddCommand(dotfilesCmd)
	rootCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(webappCmd)
	rootCmd.AddCommand(fontCmd)
	rootCmd.AddCommand(mediaCmd)
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
