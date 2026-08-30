package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/config"
)

// RenderFetchCard generates a high-impact TUI system information summary
func RenderFetchCard(cfg *config.Config) string {
	home := os.Getenv("HOME")
	_ = home

	// 1. Gather System Metrics
	kernel := getKernelVersion()
	sessionType := getSessionType()
	shellName := getShellName()
	aiModel := cfg.AI.DefaultModel
	themeName := strings.Title(cfg.Theme.Palette)

	// 2. Format Info Lines
	infoStyle := lipgloss.NewStyle().Foreground(TextLightColor)
	labelStyle := lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
	accentValStyle := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)

	lines := []string{
		fmt.Sprintf("%s %s", labelStyle.Render("OS:          "), infoStyle.Render("Ubuntu 24.04 LTS (Noble Numbat)")),
		fmt.Sprintf("%s %s", labelStyle.Render("Kernel:      "), infoStyle.Render(kernel)),
		fmt.Sprintf("%s %s", labelStyle.Render("Desktop:     "), infoStyle.Render(fmt.Sprintf("GNOME Shell 46 (%s)", sessionType))),
		fmt.Sprintf("%s %s", labelStyle.Render("Shell:       "), infoStyle.Render(shellName)),
		fmt.Sprintf("%s %s", labelStyle.Render("Vula Engine: "), accentValStyle.Render(fmt.Sprintf("Go %s | Version %s", runtime.Version(), Version))),
		fmt.Sprintf("%s %s", labelStyle.Render("Active Theme:"), infoStyle.Render(themeName)),
		fmt.Sprintf("%s %s", labelStyle.Render("Local AI:    "), infoStyle.Render(fmt.Sprintf("Ollama (%s)", aiModel))),
		fmt.Sprintf("%s %s", labelStyle.Render("Voice AI:    "), infoStyle.Render("Whisper STT + Piper Neural TTS")),
	}

	// 3. Color Blocks Preview
	colorBlocks := renderColorBlocks()

	infoBlock := strings.Join(lines, "\n") + "\n\n" + colorBlocks

	// 4. Combine Logo and Info Side-by-Side
	logo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		MarginRight(3).
		Render(LogoASCII)

	content := lipgloss.JoinHorizontal(lipgloss.Center, logo, infoBlock)

	return CardStyle.Width(82).Render(content) + "\n"
}

func getKernelVersion() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "Linux"
	}
	return strings.TrimSpace(string(out))
}

func getSessionType() string {
	session := os.Getenv("XDG_SESSION_TYPE")
	if session == "" {
		return "Wayland/X11"
	}
	return strings.Title(session)
}

func getShellName() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}
	return "fish"
}

func renderColorBlocks() string {
	colors := []lipgloss.Color{
		lipgloss.Color("#1E1E2E"), // Dark
		lipgloss.Color("#EF4444"), // Red
		lipgloss.Color("#10B981"), // Green
		lipgloss.Color("#F59E0B"), // Yellow
		lipgloss.Color("#3B82F6"), // Blue
		lipgloss.Color("#7C3AED"), // Purple
		lipgloss.Color("#06B6D4"), // Cyan
		lipgloss.Color("#F8FAFC"), // White
	}

	var blocks strings.Builder
	for _, c := range colors {
		blocks.WriteString(lipgloss.NewStyle().Background(c).Render("   "))
	}
	return blocks.String()
}
