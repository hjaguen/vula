package theme

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/config"
)

// RenderPreview outputs a terminal preview of the palette
func RenderPreview(p ThemePalette) string {
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(p.Background)).Foreground(lipgloss.Color(p.Foreground))
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(p.AccentColor)).Bold(true)
	secStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(p.SecondaryColor)).Bold(true)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.AccentColor)).
		Padding(1, 2).
		Width(60)

	hdrBg := p.HeaderBg
	if hdrBg == "" {
		hdrBg = p.Background
	}
	tbBg := p.TitlebarBg
	if tbBg == "" {
		tbBg = p.Background
	}

	content := fmt.Sprintf(
		"%s\n%s\n\n  • Background:      %s\n  • Foreground:      %s\n  • Primary Accent:   %s\n  • Secondary:       %s\n  • Top Bar Color:   %s\n  • Titlebar Color:  %s\n  • GNOME Accent:    %s",
		accentStyle.Render("⚡ THEME PREVIEW: "+p.DisplayName),
		secStyle.Render("Palette ID: "+p.Name),
		p.Background,
		p.Foreground,
		accentStyle.Render(p.AccentColor+"  ■■■■■"),
		secStyle.Render(p.SecondaryColor+"  ■■■■■"),
		hdrBg,
		tbBg,
		p.GnomeAccent,
	)

	return card.Render(bgStyle.Render(content))
}

// RunInteractiveCreator guides the user through creating a theme with Charm Huh
func RunInteractiveCreator(cfg *config.Config) (*ThemePalette, error) {
	var (
		name        string
		displayName string
		bg          string = "#1E1E2E"
		fg          string = "#CDD6F4"
		accent      string = "#7AA2F7"
		sec         string = "#06B6D4"
		headerBg    string = "#11111B"
		titlebarBg  string = "#181825"
		gnomeAccent string = "purple"
		confirm     bool   = true
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Theme Identifier (slug)").
				Description("Lowercase, no spaces (e.g. cyber-emerald, obsidian-gold)").
				Value(&name).
				Validate(func(str string) error {
					if len(strings.TrimSpace(str)) < 3 {
						return fmt.Errorf("name must be at least 3 characters")
					}
					return nil
				}),

			huh.NewInput().
				Title("Display Name").
				Description("Human-readable title (e.g. Cyber Emerald)").
				Value(&displayName),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Background Color (Hex)").
				Value(&bg).
				Placeholder("#181825"),

			huh.NewInput().
				Title("Foreground / Text Color (Hex)").
				Value(&fg).
				Placeholder("#E0DEF4"),

			huh.NewInput().
				Title("Primary Accent Color (Hex)").
				Value(&accent).
				Placeholder("#10B981"),

			huh.NewInput().
				Title("Secondary Highlight Color (Hex)").
				Value(&sec).
				Placeholder("#06B6D4"),

			huh.NewInput().
				Title("Top Bar Background Color (Hex)").
				Value(&headerBg).
				Placeholder("#11111B"),

			huh.NewInput().
				Title("Window Titlebar Color (Hex)").
				Value(&titlebarBg).
				Placeholder("#181825"),

			huh.NewSelect[string]().
				Title("GNOME Desktop Accent Color").
				Options(
					huh.NewOption("Purple / Violet", "purple"),
					huh.NewOption("Blue / Cyan", "blue"),
					huh.NewOption("Teal / Mint", "teal"),
					huh.NewOption("Green / Emerald", "green"),
					huh.NewOption("Orange / Amber", "orange"),
					huh.NewOption("Red / Ruby", "red"),
					huh.NewOption("Slate / Neutral", "slate"),
				).
				Value(&gnomeAccent),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Save and apply this theme immediately?").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	if displayName == "" {
		displayName = strings.ReplaceAll(name, "-", " ")
	}

	palette := ThemePalette{
		Name:           strings.ToLower(strings.TrimSpace(name)),
		DisplayName:    displayName,
		Background:     bg,
		Foreground:     fg,
		AccentColor:    accent,
		SecondaryColor: sec,
		HeaderBg:       headerBg,
		TitlebarBg:     titlebarBg,
		TitlebarFg:     fg,
		GnomeAccent:    gnomeAccent,
		GtkTheme:       "Yaru-dark",
		IsCustom:       true,
	}

	mgr := NewManager(cfg)
	if err := mgr.SaveCustomTheme(palette); err != nil {
		return nil, fmt.Errorf("failed to save custom theme: %w", err)
	}

	if confirm {
		if err := mgr.ApplyTheme(palette.Name); err != nil {
			return nil, fmt.Errorf("failed applying theme: %w", err)
		}
	}

	return &palette, nil
}

// GenerateAITheme uses local Ollama AI to generate a cohesive color palette from prompt
func GenerateAITheme(ctx context.Context, cfg *config.Config, prompt string) (*ThemePalette, error) {
	aiClient := ai.NewClient(cfg)

	systemPrompt := `You are an expert color designer and theme creator for developer operating systems.
The user will give you a theme concept or aesthetic idea.
Generate a cohesive, accessible, high-contrast dark theme palette.
Respond ONLY with a JSON object strictly adhering to this schema:
{
  "name": "lowercase-slug-without-spaces",
  "display_name": "Title Name",
  "background": "#HEX",
  "foreground": "#HEX",
  "accent_color": "#HEX",
  "secondary_color": "#HEX",
  "header_bg": "#HEX",
  "titlebar_bg": "#HEX",
  "gnome_accent": "one of: purple, blue, teal, green, orange, red, slate"
}`

	userMessage := fmt.Sprintf("%s\nTheme Concept: %s", systemPrompt, prompt)
	rawResponse, err := aiClient.Ask(ctx, userMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("AI query failed: %w", err)
	}

	// Extract JSON block if surrounded by markdown code fences
	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonMatch := re.FindString(rawResponse)
	if jsonMatch == "" {
		jsonMatch = rawResponse
	}

	var palette ThemePalette
	if err := json.Unmarshal([]byte(jsonMatch), &palette); err != nil {
		// Fallback clean extraction
		palette = ThemePalette{
			Name:           "custom-ai",
			DisplayName:    "AI Generated Theme",
			Background:     "#121218",
			Foreground:     "#E2E8F0",
			AccentColor:    "#7C3AED",
			SecondaryColor: "#06B6D4",
			HeaderBg:       "#0D0D12",
			TitlebarBg:     "#1A1A24",
			GnomeAccent:    "purple",
		}
	}

	palette.Name = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(palette.Name), " ", "-"))
	if palette.Name == "" {
		palette.Name = "ai-theme"
	}
	if palette.DisplayName == "" {
		palette.DisplayName = "AI " + prompt
	}
	if palette.TitlebarFg == "" {
		palette.TitlebarFg = palette.Foreground
	}
	palette.GtkTheme = "Yaru-dark"
	palette.IsCustom = true

	mgr := NewManager(cfg)
	if err := mgr.SaveCustomTheme(palette); err != nil {
		return nil, err
	}
	_ = mgr.ApplyTheme(palette.Name)

	return &palette, nil
}
