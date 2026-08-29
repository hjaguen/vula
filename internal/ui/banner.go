package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const LogoASCII = `
 ██▒   █▓ █▒   █▓ ██▓    ▄▄▄      
▓██░   █▒▓██░   █▒▓██▒   ▒████▄    
 ▓██  █▒░ ▓██  █▒░▒██░   ▒██  ▀█▄  
  ▒██ █░░  ▒██ █░░▒██░   ░██▄▄▄▄██ 
   ▒▀█░     ▒▀█░  ░██████▒▓█   ▓██▒
   ░ ▐░     ░ ▐░  ░ ▒░▓  ░▒▒   ▓▒█░
   ░ ░░     ░ ░░  ░ ░ ▒  ░ ▒   ▒▒ ░
     ░░       ░░    ░ ░    ░   ▒   
      ░        ░      ░  ░     ░  ░
     ░        ░                    
`

const Version = "v0.1.0-alpha"

func RenderBanner() string {
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Render(LogoASCII)

	tagline := lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true).
		Render("  ⚡ VULA — Next-Gen AI & Voice Powered Developer OS for Ubuntu")

	meta := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render(fmt.Sprintf("     Version: %s | Ubuntu 24.04 LTS (Noble) | 100%% Go Engine\n", Version))

	return lipgloss.JoinVertical(lipgloss.Left, logo, tagline, meta)
}

func RenderHeader(title, subtitle string) string {
	b := strings.Builder{}
	b.WriteString(HeaderStyle.Render(fmt.Sprintf(" VULA // %s ", strings.ToUpper(title))))
	b.WriteString("\n")
	if subtitle != "" {
		b.WriteString(SubtitleStyle.Render(subtitle))
		b.WriteString("\n\n")
	}
	return b.String()
}
