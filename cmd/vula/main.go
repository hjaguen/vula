package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/apps"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/doctor"
	"github.com/vula-os/vula/internal/dotfiles"
	"github.com/vula-os/vula/internal/gnome"
	"github.com/vula-os/vula/internal/hud"
	"github.com/vula-os/vula/internal/installer"
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
		fmt.Println(ui.SuccessStyle.Render("✓ GNOME Shell optimizations and global keybindings applied successfully!"))
		fmt.Println("  • [Super + Space] -> Vula HUD")
		fmt.Println("  • [Super + Alt + V] -> Voice Dictation")
	},
}

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage and switch global system and terminal themes",
}

var themeSetCmd = &cobra.Command{
	Use:   "set [theme-name]",
	Short: "Apply a global theme (tokyonight, catppuccin, nord, rose-pine)",
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
	Short: "List all available global themes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.RenderHeader("Available Themes", "Unified aesthetic palettes for Vula"))
		for _, t := range theme.ListThemes() {
			fmt.Printf("  • %-16s %s (Accent: %s)\n", ui.InfoStyle.Render(t.Name), t.DisplayName, t.AccentColor)
		}
		fmt.Println()
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

func init() {
	aiCmd.AddCommand(aiAskCmd)
	aiCmd.AddCommand(aiCmdSuggest)
	aiCmd.AddCommand(aiModelsCmd)

	voiceCmd.AddCommand(voiceRecordCmd)
	voiceCmd.AddCommand(voiceSpeakCmd)

	desktopCmd.AddCommand(desktopSetupCmd)

	themeCmd.AddCommand(themeSetCmd)
	themeCmd.AddCommand(themeListCmd)

	dotfilesCmd.AddCommand(dotfilesInstallCmd)

	appsCmd.AddCommand(appsInstallCLICmd)
	appsCmd.AddCommand(appsListCmd)

	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(hudCmd)
	rootCmd.AddCommand(aiCmd)
	rootCmd.AddCommand(voiceCmd)
	rootCmd.AddCommand(desktopCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(dotfilesCmd)
	rootCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
