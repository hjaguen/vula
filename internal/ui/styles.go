package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Vula Color Palette
var (
	PrimaryColor   = lipgloss.Color("#7C3AED") // Deep Violet / Royal Purple
	SecondaryColor = lipgloss.Color("#06B6D4") // Cyan
	AccentColor    = lipgloss.Color("#10B981") // Emerald Green
	WarningColor   = lipgloss.Color("#F59E0B") // Amber
	DangerColor    = lipgloss.Color("#EF4444") // Red
	MutedColor     = lipgloss.Color("#64748B") // Slate Muted
	BgDarkColor    = lipgloss.Color("#0F172A") // Deep Slate Background
	TextLightColor = lipgloss.Color("#F8FAFC") // Off-white
	BorderColor    = lipgloss.Color("#334155") // Subtle Slate Border
)

// Lipgloss Styles
var (
	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextLightColor).
			Background(PrimaryColor).
			Padding(0, 1).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

	WarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(WarningColor)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(DangerColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	// Container & Box Styles
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2).
			Margin(0, 1)

	SubtleBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Margin(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextLightColor).
			Background(PrimaryColor).
			Padding(0, 1)

	SuccessBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(AccentColor).
			Padding(0, 1)

	WarnBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(WarningColor).
			Padding(0, 1)

	ErrorBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(DangerColor).
			Padding(0, 1)
)
