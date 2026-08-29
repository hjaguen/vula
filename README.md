# ⚡ Vula (v0.1.0-alpha)

> **"Vula"** *(verb, Zulu & Xhosa)*: **To open, to start, to unlock.**

**Vula** is an opinionated, keyboard-first developer operating environment crafted for **Ubuntu 24.04 LTS (Noble Numbat)**. It combines the rock-solid hardware stability of Ubuntu and GNOME Shell with the elegance of a **100% Go + Charm TUI** orchestrator, native local AI intelligence, conversational real-time voice assistant, unified theming, and modern developer dotfiles.

[![CI Pipeline](https://github.com/hjaguen/vula/actions/workflows/ci.yml/badge.svg)](https://github.com/hjaguen/vula/actions/workflows/ci.yml)
[![Security Scan](https://github.com/hjaguen/vula/actions/workflows/security.yml/badge.svg)](https://github.com/hjaguen/vula/actions/workflows/security.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Ubuntu](https://img.shields.io/badge/Ubuntu-24.04%20LTS-E95420?style=flat&logo=ubuntu)](https://ubuntu.com/)
[![Charm](https://img.shields.io/badge/TUI-Charm_Ecosystem-FF5F87?style=flat)](https://charm.sh/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

```text
 ██▒   █▓ █▒   █▓ ██▓    ▄▄▄      
▓██░   █▒▓██░   █▒▓██▒   ▒████▄    
 ▓██  █▒░ ▓██  █▒░▒██░   ▒██  ▀█▄  
  ▒██ █░░  ▒██ █░░▒██░   ░██▄▄▄▄██ 
   ▒▀█░     ▒▀█░  ░██████▒▓█   ▓██▒
   ░ ▐░     ░ ▐░  ░ ▒░▓  ░▒▒   ▓▒█░
   ░ ░░     ░ ░░  ░ ░ ▒  ░ ▒   ▒▒ ░
     ░░       ░░    ░ ░    ░   ▒   
      ░        ░      ░  ░     ░  ░
```

---

## 🎯 Why Vula?

While projects like **Omakub** (Ubuntu/GNOME) and **Omarchy** (Arch/Hyprland) paved the way for curated developer workflows, developers running Ubuntu LTS were left without a modern, AI-native environment that retained LTS hardware reliability.

**Vula delivers this through five core pillars:**

1. **Rock-Solid LTS Stability:** Built natively on top of Ubuntu 24.04 LTS (Noble Numbat) and GNOME Shell 46+. Zero kernel panics, full GPU acceleration support, and production-grade stability.
2. **100% Go & Charm TUI:** No fragile 5,000-line bash scripts. Vula is powered by a compiled, type-safe Go binary utilizing the [Charm](https://charm.sh) ecosystem (`bubbletea`, `lipgloss`, `huh`, `bubbles`).
3. **Conversational Local AI & Voice:** Integrated directly into the OS with active desktop context awareness (reads focused window state and clipboard safely) powered by Ollama (`qwen2.5-coder`, `llama3.2`), local Whisper.cpp (STT), and Piper (neural TTS).
4. **Unified Aesthetics (`vula theme`):** Instant system-wide theme switcher across GNOME Shell, terminal emulators (Ghostty/Kitty), Neovim, and Starship prompt.
5. **Turnkey Dotfiles & App Catalog:** One-command installation of modern Unix CLI tools (`eza`, `bat`, `lazygit`, `zoxide`, `btop`, `fzf`) and pre-configured dotfiles (Fish, Starship, Neovim Lua IDE, Tmux).

---

## 🏗 System Architecture

```text
vula/
├── cmd/vula/              # Cobra CLI Entrypoint & Subcommands
├── internal/
│   ├── ai/                # Context-Aware Ollama Client & Shell Translator
│   ├── apps/              # Curated Developer CLI & GUI Apps Catalog
│   ├── config/            # Declarative YAML Config (~/.config/vula/config.yaml)
│   ├── doctor/            # Self-Healing System Diagnostic Engine
│   ├── dotfiles/          # Fish, Starship, Neovim Lua, and Tmux Managers
│   ├── gnome/             # Dconf, Keybindings, Tiling & Extensions Manager
│   ├── hud/               # Floating Bubbletea Launcher & Multi-view HUD
│   ├── installer/         # Interactive Huh-based Idempotent Installer
│   ├── packages/          # Idempotent APT, Snap & Binary Toolchain Manager
│   ├── theme/             # Unified Global Theme Switcher Engine
│   ├── ui/                # Charm Lipgloss Palettes & Banners
│   └── voice/             # STT (Whisper), Neural TTS (Piper) & Voice Assistant
├── dotfiles/              # Curated Dotfile Templates
├── scripts/
│   ├── bootstrap.sh       # Safe, Non-Destructive Bootstrap Script
│   └── vula-hud-launch    # Standalone Floating HUD Launcher
├── .github/workflows/     # CI Pipeline & Security/Vulnerability Scanners
├── Makefile               # Standard Development Targets
├── SECURITY.md            # Security & Vulnerability Policy
└── README.md
```

---

## 🚀 Quickstart

### 1. Requirements
* **OS:** Ubuntu 24.04 LTS (Noble Numbat)
* **Desktop:** GNOME Shell 46+ (Wayland or X11)
* **Toolchain:** Go 1.24+ (managed automatically)

### 2. Installation
```bash
git clone https://github.com/hjaguen/vula.git ~/repos/vula
cd ~/repos/vula
make install
vula desktop setup
```

### 3. Verify System Diagnostics
```bash
vula doctor
```

```text
  VULA // DOCTOR DIAGNOSTICS  
Checking system readiness, desktop environment, and AI stack

◆ Operating System
  [  OK  ] Ubuntu Version           Ubuntu 24.04 LTS (Noble Numbat) detected
  [  OK  ] Linux Kernel             6.8.0-138-generic

◆ Desktop Environment
  [  OK  ] GNOME Shell              GNOME Shell 46.0
  [  OK  ] Display Server           X11 / Wayland session active
  [  OK  ] dconf CLI                dconf binary available for declarative configuration

◆ Hardware & Audio
  [  OK  ] Audio Subsystem          PipeWire / PulseAudio detected (Ready for voice STT/TTS)
  [  OK  ] GPU Acceleration         Vulkan / OpenGL standard GPU inference mode

◆ AI Subsystem
  [  OK  ] Ollama Local Daemon      Connected to Ollama at http://localhost:11434

◆ Voice Subsystem
  [  OK  ] Microphone Capture       ALSA/PipeWire audio capture utility available

Diagnostic Summary: 12 Passed | 0 Failures
```

---

## ⌨️ Daily Workflow & Keybindings

| Shortcut | Action | Description |
| :--- | :--- | :--- |
| **`Super + Space`** | **Vula Floating HUD** | Raycast-style launcher, action palette, and AI assistant |
| **`Super + Alt + A`** | **Active Voice AI** | Conversational assistant: speak your question, receive spoken answer + desktop notification |
| **`Super + Alt + V`** | **Voice Dictation** | Transcribe spoken words directly into the active editor or input field |
| **`Super + Alt + H/J/K/L`** | **Workspace Navigation** | Switch workspaces dynamically using Vim keys |
| **`Super + Shift + H/J/K/L`** | **Move Window** | Move focused window to adjacent workspace |
| **`Super + H`** | **Minimize Window** | Hide/minimize currently focused window |
| **`Super + D`** | **Toggle Desktop** | Show or hide all open desktop windows |
| **`Super + Q`** | **Close Window** | Close focused application window |
| **`Super + M`** | **Maximize** | Toggle window maximization |

---

## 🎙️ Active Voice AI & Dictation

### 1. Active Conversational Voice Assistant
Press **`Super + Alt + A`** or run:
```bash
vula listen
```
* **Audio Capture:** Records your voice cleanly through PipeWire.
* **STT:** Transcribes speech locally via Whisper AVX2 in Spanish/English.
* **Contextual AI:** Queries the local LLM with active window and clipboard context.
* **Dual Output:** Speaks the answer back through your speakers with **Piper TTS** in <0.8s and displays a desktop notification.

### 2. Direct Voice Dictation into Any App
Press **`Super + Alt + V`** or run:
```bash
vula voice record
```

### 3. Neural Speech Synthesis
```bash
vula voice speak "Hola Mauricio, el motor de voz de Vula está listo y operativo."
```

---

## 🤖 Local AI & Shell Commands

### Ask AI with Desktop Context
```bash
# Query the local LLM with active window & clipboard context:
vula ai ask "How do I optimize this SQL query in my clipboard?"
```

### Natural Language Shell Translator
```bash
# Translate intent directly to a safe, executable bash command:
vula ai cmd "find all files larger than 100MB modified in the last 7 days"
```

### List Installed Local Models
```bash
vula ai models
```

---

## 🎨 Global Theme Engine (`vula theme`)

Switch system-wide themes with a single command across GNOME Shell, GTK accent colors, terminal emulators (Ghostty), and the HUD:

```bash
# List available palettes:
vula theme list

# Switch theme instantly:
vula theme set tokyonight    # Tokyo Night (Default)
vula theme set catppuccin    # Catppuccin Mocha
vula theme set nord          # Nord Arctic
vula theme set rose-pine     # Rosé Pine
```

---

## 📁 Developer Dotfiles (`vula dotfiles`)

Deploy opinionated, battle-tested configurations with automatic backups into `~/.config/vula/backups/`:

```bash
vula dotfiles install
```

* **Fish Shell (`~/.config/fish/config.fish`):** Preloaded with modern aliases (`ls -> eza`, `cat -> bat`, `lg -> lazygit`, `ask -> vula ai ask`, `vcmd -> vula ai cmd`).
* **Starship Prompt (`~/.config/starship.toml`):** Minimalist two-line prompt with Git status, execution time, and Node/Go/Rust toolchains.
* **Neovim (`~/.config/nvim/init.lua`):** Modern Lua configuration with relative line numbers, system clipboard sharing, and Vim window navigation.
* **Tmux (`~/.tmux.conf`):** `Ctrl+A` prefix, mouse scrolling, and Vula status bar.

---

## 📦 Curated App Catalog (`vula apps`)

```bash
# List curated developer software:
vula apps list

# Install complete modern CLI stack (eza, bat, lazygit, starship, zoxide, btop, fzf):
vula apps install-cli
```

---

## 🛡 Security Guarantees

* **Zero Root Daemons:** AI and Voice engines run strictly in user space (`$USER`).
* **Human-in-the-Loop:** Destructive commands (`rm`, `sudo`, `dd`) always require explicit user confirmation.
* **Supply-Chain Verification:** All downloaded binaries and neural models are verified against cryptographic SHA-256 checksums.
* **Local-First Privacy:** Prompts and audio streams never leave your device unless you explicitly configure an external API provider.

---

## 🗺 Roadmap

- [x] **v0.1.0-alpha:** Core Go + Charm CLI, `vula doctor`, Declarative YAML configs, Floating HUD, Local Ollama integration, Whisper STT + Piper TTS pipeline.
- [x] **v0.2.0-alpha:** Conversational Active Voice AI (`vula listen`), `Super + Alt + A` global shortcut, System notification integration.
- [x] **v0.3.0-alpha:** Global Theme Engine (`tokyonight`, `catppuccin`, `nord`, `rose-pine`), Dotfiles manager (Fish, Starship, Neovim, Tmux), Developer CLI app catalog.
- [ ] **v0.4.0:** Automatic Tiling Manager extension installer (Forge / Pop Shell for GNOME 46).
- [ ] **v1.0.0:** Production Release, `.deb` packaging, PPA & Sigstore release signing.

---

## 📄 License & Contributing

Distributed under the **MIT License**. See [LICENSE](LICENSE) and [SECURITY.md](SECURITY.md) for more details. Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).
