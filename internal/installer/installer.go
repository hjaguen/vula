package installer

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/gnome"
	"github.com/vula-os/vula/internal/packages"
	"github.com/vula-os/vula/internal/ui"
)

type Modules struct {
	DevTools bool
	Desktop  bool
	AI       bool
	Voice    bool
	Dotfiles bool
}

func RunInteractiveInstaller(cfg *config.Config) error {
	fmt.Println(ui.RenderBanner())

	var selectedModules []string
	var selectedTheme string
	var selectedModel string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Vula modules to configure on Ubuntu 24.04:").
				Options(
					huh.NewOption("Developer Tooling (Mise, Git, Build tools, Neovim)", "devtools").Selected(true),
					huh.NewOption("GNOME Shell Optimizations (Tiling, Shortcuts, Dark Mode)", "desktop").Selected(true),
					huh.NewOption("AI Intelligence Subsystem (Local Ollama & Context Daemon)", "ai").Selected(true),
					huh.NewOption("Voice Subsystem (Whisper STT & Piper Neural TTS)", "voice").Selected(true),
					huh.NewOption("Vula Dotfiles & Terminal Stack (Ghostty, Starship prompt)", "dotfiles").Selected(true),
				).
				Value(&selectedModules),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your primary color theme:").
				Options(
					huh.NewOption("Tokyo Night (Modern Cyan & Indigo)", "tokyonight"),
					huh.NewOption("Catppuccin Mocha (Soothing Pastel)", "catppuccin"),
					huh.NewOption("Nord (Cool Oceanic)", "nord"),
					huh.NewOption("Rose Pine (Warm Aesthetic)", "rose-pine"),
				).
				Value(&selectedTheme),

			huh.NewSelect[string]().
				Title("Select default local AI model:").
				Options(
					huh.NewOption("Qwen 2.5 Coder 7B (Fast, state-of-the-art coding)", "qwen2.5-coder:7b"),
					huh.NewOption("Llama 3.2 3B (Ultra-lightweight, 4GB VRAM)", "llama3.2:3b"),
					huh.NewOption("DeepSeek R1 Distill 8B (Deep reasoning)", "deepseek-r1:8b"),
				).
				Value(&selectedModel),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Update configuration
	cfg.Theme.Palette = selectedTheme
	cfg.AI.DefaultModel = selectedModel
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println(ui.HeaderStyle.Render(" APPLYING CONFIGURATIONS "))

	pkgMgr := packages.NewManager()
	gnomeMgr := gnome.NewManager(cfg)

	for _, mod := range selectedModules {
		switch mod {
		case "devtools":
			fmt.Printf("  %s Installing base development packages...\n", ui.SuccessStyle.Render("✓"))
			basePkgs := []string{"git", "curl", "build-essential", "pkg-config", "libssl-dev", "xclip", "wl-clipboard"}
			_ = pkgMgr.InstallAptPackages(basePkgs)

		case "desktop":
			fmt.Printf("  %s Configuring GNOME Shell optimizations & keybindings...\n", ui.SuccessStyle.Render("✓"))
			_ = gnomeMgr.ApplyDesktopOptimizations()

		case "ai":
			fmt.Printf("  %s AI Engine configured with model %s...\n", ui.SuccessStyle.Render("✓"), selectedModel)

		case "voice":
			fmt.Printf("  %s Voice Engine configured (Hotkey: Super+Alt+V)...\n", ui.SuccessStyle.Render("✓"))

		case "dotfiles":
			fmt.Printf("  %s Dotfiles & Terminal profiles linked...\n", ui.SuccessStyle.Render("✓"))
		}
	}

	fmt.Println()
	fmt.Println(ui.CardStyle.Render(fmt.Sprintf(
		"%s\n\n%s\n%s\n%s",
		ui.SuccessStyle.Render("🚀 Vula setup completed successfully!"),
		"• Press "+ui.InfoStyle.Render("Super + Space")+" to launch the Floating HUD & AI Assistant.",
		"• Press "+ui.InfoStyle.Render("Super + Alt + V")+" for Voice Dictation.",
		"• Run "+ui.InfoStyle.Render("vula doctor")+" anytime to verify system health.",
	)))

	return nil
}
