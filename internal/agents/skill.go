package agents

import (
	"os"
	"path/filepath"
	"strings"
)

// GenerateSkillMarkdown returns the comprehensive Agent Skill for Vula OS
func GenerateSkillMarkdown() string {
	const raw = `---
name: vula
description: Expert AI Agent control and inspection system for Vula OS on Ubuntu 24.04+ LTS
tools:
  - run_command
  - view_file
---

# Vula OS Agent Skill (The Malleable Developer OS)

You are interacting with **Vula OS**, a keyboard-first, AI & voice-powered developer operating environment crafted for Ubuntu 24.04+ LTS, built 100% in Go + Charm TUI.

As an AI agent, you can inspect, modify, and control the user's desktop environment safely using the official Vula CLI tools instead of inventing risky ad-hoc shell commands.

---

## ⚡ 1. Natural Language OS Control (@vula do@)

Whenever the user asks you to perform an OS-level task (adjusting hardware, launching apps, converting media, changing themes), use @vula do@:

@@@bash
# Hardware & Sound
vula do "ajusta el brillo al 50%"
vula do "sube el volumen al 80%"
vula do "bloquea la pantalla"

# Web Applications & Windows
vula do "abre Visual Studio Code"
vula do "crea una app web para ChatGPT en https://chatgpt.com"
vula do "convierte la última grabación de pantalla a mp4"

# Safe & Guardrailed
vula do "<instrucción en lenguaje natural>"
@@@

> **Note:** Actions are evaluated in 3 safety tiers (Safe, Sensitive, High-Risk) and audited in @~/.config/vula/audit.log@.

---

## 🤖 2. Herdr AI Agent Multiplexer & Fleet (@vula agents@)

Vula runs Herdr ("tmux for AI agents") to maintain persistent background worker panes:

@@@bash
# Check running agent fleet status
vula agents status

# List active agent sessions/panes
vula agents list

# Offload long-running task to Herdr
vula do --agent "analiza la base de código y ejecuta las pruebas"

# Attach interactive terminal to agent pane
vula agents attach <session-id>
@@@

---

## 🎨 3. Themes & Aesthetics (@vula theme@ & @vula font@)

Vula harmonizes themes across GNOME, terminal (Ghostty/Kitty), Neovim, and Starship:

@@@bash
# Available palettes: tokyonight, catppuccin, nord, rose-pine, everforest, dracula, gruvbox, etc.
vula theme set tokyonight
vula theme list

# Developer Nerd Fonts (jetbrains, cascadia, firacode, meslo)
vula font set jetbrains
vula font size 12
@@@

---

## 🪟 4. Tiling & Focus Mode (@vula desktop@)

@@@bash
# Toggle Zen/Focus mode (hide distraction bars & enable DND)
vula desktop focus

# Toggle Automatic Window Tiling
vula desktop toggle-tile
@@@

---

## 🩺 5. System Health & Diagnostics (@vula doctor@)

Always run diagnostics to verify if system audio, Whisper STT, Ollama, or GNOME extensions are healthy:

@@@bash
vula doctor
@@@
`
	return strings.ReplaceAll(raw, "@", "`")
}

// WriteSkillFile generates ~/.config/vula/SKILL.md
func WriteSkillFile() (string, error) {
	home := os.Getenv("HOME")
	targetDir := filepath.Join(home, ".config", "vula")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	targetPath := filepath.Join(targetDir, "SKILL.md")
	if err := os.WriteFile(targetPath, []byte(GenerateSkillMarkdown()), 0644); err != nil {
		return "", err
	}
	return targetPath, nil
}
