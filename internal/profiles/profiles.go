package profiles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/apps"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/font"
	"github.com/vula-os/vula/internal/packages"
	"github.com/vula-os/vula/internal/theme"
	"github.com/vula-os/vula/internal/ui"
)

type Profile struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Icon          string   `json:"icon"`
	Category      string   `json:"category"`
	AptPackages   []string `json:"apt_packages"`
	CLIStack      bool     `json:"cli_stack"`
	WebApps       []string `json:"web_apps"`
	DeveloperFont string   `json:"developer_font"`
	DefaultTheme  string   `json:"default_theme"`
}

var Catalog = []Profile{
	{
		ID:            "full-stack-dev",
		Name:          "Full-Stack Web & Backend Developer",
		Description:   "Complete stack: Node.js, Go, Python, Docker, Postman, VS Code, Lazygit & JetBrains Mono",
		Icon:          "💻",
		Category:      "Software Engineering",
		AptPackages:   []string{"git", "curl", "build-essential", "pkg-config", "libssl-dev", "xclip", "wl-clipboard"},
		CLIStack:      true,
		WebApps:       []string{"Postman", "Excalidraw"},
		DeveloperFont: "jetbrains",
		DefaultTheme:  "tokyonight",
	},
	{
		ID:            "devops",
		Name:          "DevOps & Cloud Infrastructure Engineer",
		Description:   "Cloud & containers: Docker, kubectl, Helm, k9s, Terraform, Ansible, Tmux & Hack font",
		Icon:          "☁️",
		Category:      "Infrastructure",
		AptPackages:   []string{"git", "curl", "tmux", "htop", "net-tools", "dnsutils", "xclip", "wl-clipboard"},
		CLIStack:      true,
		WebApps:       []string{"Datadog", "Excalidraw"},
		DeveloperFont: "cascadia",
		DefaultTheme:  "nord",
	},
	{
		ID:            "designer-creator",
		Name:          "UI/UX Designer & Digital Creator",
		Description:   "Creative suite: GIMP, Inkscape, Blender, OBS Studio, Figma WebApp & Fira Code",
		Icon:          "🎨",
		Category:      "Design & Media",
		AptPackages:   []string{"gimp", "inkscape", "obs-studio", "vlc", "ffmpeg"},
		CLIStack:      false,
		WebApps:       []string{"Figma", "Canva", "Excalidraw"},
		DeveloperFont: "firacode",
		DefaultTheme:  "catppuccin",
	},
	{
		ID:            "data-dba",
		Name:          "Database & Data Engineer",
		Description:   "Data & DBs: PostgreSQL Client, Redis CLI, DBeaver, Python Data tools & Cascadia Code",
		Icon:          "🗄️",
		Category:      "Data & Analytics",
		AptPackages:   []string{"postgresql-client", "redis-tools", "python3-pip", "python3-venv"},
		CLIStack:      true,
		WebApps:       []string{"MongoDB Compass", "Excalidraw"},
		DeveloperFont: "cascadia",
		DefaultTheme:  "dracula",
	},
	{
		ID:            "minimal",
		Name:          "Minimalist Ergonomic Workstation",
		Description:   "Core Vula environment: Tiling Assistant, Tactile Grid, Floating HUD, Voice AI & Starship",
		Icon:          "⚡",
		Category:      "Minimal",
		AptPackages:   []string{"git", "curl", "xclip", "wl-clipboard"},
		CLIStack:      false,
		WebApps:       []string{},
		DeveloperFont: "jetbrains",
		DefaultTheme:  "tokyonight",
	},
}

// GetProfileByID searches the profile catalog by ID
func GetProfileByID(id string) (*Profile, bool) {
	for _, p := range Catalog {
		if strings.EqualFold(p.ID, id) {
			return &p, true
		}
	}
	return nil, false
}

// ApplyProfile provisions tools, fonts, webapps, theme, and config matching the target profile
func ApplyProfile(cfg *config.Config, profileID string) error {
	p, found := GetProfileByID(profileID)
	if !found {
		return fmt.Errorf("unknown profile ID '%s'", profileID)
	}

	fmt.Printf("%s Applying Workstation Profile: %s (%s)\n\n", ui.InfoStyle.Render("⚡"), p.Name, p.ID)

	pkgMgr := packages.NewManager()
	fontMgr, _ := font.NewManager()
	themeMgr := theme.NewManager(cfg)
	appMgr := apps.NewManager(cfg)

	// 1. Install APT Packages
	if len(p.AptPackages) > 0 {
		fmt.Printf("  %s Installing Apt system packages: %s...\n", ui.SuccessStyle.Render("✓"), strings.Join(p.AptPackages, ", "))
		_ = pkgMgr.InstallAptPackages(p.AptPackages)
	}

	// 2. Install Developer CLI Toolchain if requested
	if p.CLIStack {
		fmt.Printf("  %s Provisioning modern CLI stack (eza, bat, lazygit, btop, fzf, starship)...\n", ui.SuccessStyle.Render("✓"))
		_ = appMgr.InstallCLIStack()
	}

	// 3. Register Preset WebApps
	webAppURLs := map[string]string{
		"Postman":         "https://identity.getpostman.com",
		"Excalidraw":      "https://excalidraw.com",
		"Figma":           "https://figma.com",
		"Canva":           "https://canva.com",
		"Datadog":         "https://app.datadoghq.com",
		"MongoDB Compass": "https://cloud.mongodb.com",
	}

	for _, webAppName := range p.WebApps {
		if targetURL, ok := webAppURLs[webAppName]; ok {
			fmt.Printf("  %s Registering WebApp: %s (%s)...\n", ui.SuccessStyle.Render("✓"), webAppName, targetURL)
			_, _ = apps.AddWebApp(webAppName, targetURL, "")
		}
	}

	// 4. Set Developer Font
	if p.DeveloperFont != "" && fontMgr != nil {
		fmt.Printf("  %s Setting developer font to %s...\n", ui.SuccessStyle.Render("✓"), p.DeveloperFont)
		_, _ = fontMgr.SetFont(strings.ToLower(p.DeveloperFont))
	}

	// 5. Apply Default Theme
	if p.DefaultTheme != "" {
		fmt.Printf("  %s Applying theme palette '%s'...\n", ui.SuccessStyle.Render("✓"), p.DefaultTheme)
		_ = themeMgr.ApplyTheme(p.DefaultTheme)
	}

	fmt.Println()
	fmt.Println(ui.CardStyle.Render(fmt.Sprintf(
		"%s\n\n%s\n%s",
		ui.SuccessStyle.Render(fmt.Sprintf("✓ Workstation Profile '%s' configured successfully!", p.Name)),
		"• Run "+ui.InfoStyle.Render("vula fetch")+" to view active workstation summary.",
		"• Press "+ui.InfoStyle.Render("Super + Space")+" anytime to launch Vula HUD.",
	)))

	return nil
}

// RenderTable returns a clean terminal table of all available profiles
func RenderTable() string {
	b := strings.Builder{}
	b.WriteString(ui.RenderHeader("Vula Workstation Profiles", "Preconfigured environment recipes for developer roles"))
	b.WriteString("\n")

	hdrStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.SecondaryColor)
	b.WriteString(hdrStyle.Render(fmt.Sprintf("  %-20s %-38s %s\n", "PROFILE ID", "NAME", "CATEGORY")))
	b.WriteString(lipgloss.NewStyle().Foreground(ui.BorderColor).Render("  "+strings.Repeat("─", 80)+"\n"))

	idStyle := lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(ui.TextLightColor).Bold(true)
	catStyle := lipgloss.NewStyle().Foreground(ui.MutedColor)

	for _, p := range Catalog {
		line := fmt.Sprintf("  %-20s %-38s %s\n", idStyle.Render(p.ID), nameStyle.Render(p.Icon+" "+p.Name), catStyle.Render(p.Category))
		b.WriteString(line)
		b.WriteString(fmt.Sprintf("                       %s\n", lipgloss.NewStyle().Foreground(ui.MutedColor).Italic(true).Render(p.Description)))
	}

	b.WriteString("\n")
	b.WriteString(ui.CardStyle.Render("Apply a profile: vula profile apply <id>  |  Interactive TUI: vula profile ui"))
	b.WriteString("\n")

	return b.String()
}

// RunInteractiveProfileSelector presents a Charm Huh form to pick a workstation profile
func RunInteractiveProfileSelector(cfg *config.Config) error {
	fmt.Println(ui.RenderHeader("Workstation Profile Selector", "Choose an environment profile matching your workflow"))

	var selectedID string

	options := make([]huh.Option[string], len(Catalog))
	for i, p := range Catalog {
		label := fmt.Sprintf("%s %s - %s", p.Icon, p.Name, p.Description)
		options[i] = huh.NewOption(label, p.ID)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select your Workstation Profile:").
				Options(options...).
				Value(&selectedID),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if selectedID == "" {
		return nil
	}

	return ApplyProfile(cfg, selectedID)
}
