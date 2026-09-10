# ⚡ Vula (v0.1.0-alpha)

> **"Vula"** *(verb, Zulu & Xhosa)*: **To open, to start, to unlock.**

**Vula** is an opinionated, keyboard-first developer operating environment crafted for **Ubuntu 24.04+ LTS (Noble Numbat & future LTS releases like 26.04)**. It combines the rock-solid hardware stability of Ubuntu and GNOME Shell with the elegance of a **100% Go + Charm TUI** orchestrator, native local/hybrid AI intelligence, conversational real-time voice assistant, OS action control engine (`vula do`), unified theming, and modern developer dotfiles.

[![CI Pipeline](https://github.com/hjaguen/vula/actions/workflows/ci.yml/badge.svg)](https://github.com/hjaguen/vula/actions/workflows/ci.yml)
[![Security Scan](https://github.com/hjaguen/vula/actions/workflows/security.yml/badge.svg)](https://github.com/hjaguen/vula/actions/workflows/security.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Ubuntu](https://img.shields.io/badge/Ubuntu-24.04%2B%20LTS-E95420?style=flat&logo=ubuntu)](https://ubuntu.com/)
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

**Vula delivers this through six core pillars:**

1. **Rock-Solid LTS Stability:** Built natively for Ubuntu 24.04+ LTS (Noble Numbat and forward-compatible with 26.04+ LTS) and GNOME Shell 46+. Zero kernel panics, full GPU acceleration support, and production-grade stability.
2. **100% Go & Charm TUI:** No fragile 5,000-line bash scripts. Vula is powered by a compiled, type-safe Go binary utilizing the [Charm](https://charm.sh) ecosystem (`bubbletea`, `lipgloss`, `huh`, `bubbles`).
3. **OS Control Engine (`vula do`):** Safe, natural language OS control to launch apps, adjust brightness/volume, manage files, and control processes with multi-tier risk classification and TUI confirmation prompts.
4. **Conversational Local & Hybrid AI & Voice:** Integrated directly into the OS with desktop context awareness (window title & clipboard) powered by Ollama, Cloud AI Providers (Groq, Gemini, OpenAI), Whisper STT, and Piper neural TTS.
5. **Unified Aesthetics & Theme Studio (`vula theme`):** System-wide theme switcher and interactive palette creator (TUI + Local AI) across GNOME Shell, terminal emulators (Ghostty/Kitty), Neovim, and Starship prompt.
6. **Turnkey Dotfiles & App Catalog:** One-command installation of modern Unix CLI tools (`eza`, `bat`, `lazygit`, `zoxide`, `btop`, `fzf`) and pre-configured dotfiles (Fish, Starship, Neovim Lua IDE, Tmux).

---

## 🏗 System Architecture

```text
vula/
├── cmd/vula/              # Cobra CLI Entrypoint & Subcommands
├── internal/
│   ├── actions/           # OS Actions Engine, Risk Guardrails & Audit Logger (vula do)
│   ├── ai/                # Context-Aware Hybrid AI Engine (Ollama, Gemini, Groq, OpenAI)
│   ├── apps/              # Curated Developer CLI & GUI Apps Catalog
│   ├── config/            # Declarative YAML Config (~/.config/vula/config.yaml)
│   ├── doctor/            # Self-Healing System Diagnostic Engine
│   ├── dotfiles/          # Fish, Starship, Neovim Lua, and Tmux Managers
│   ├── gnome/             # Dconf, Keybindings, Tiling & Extensions Manager
│   ├── hud/               # Floating Bubbletea Launcher & Multi-view HUD
│   ├── installer/         # Interactive Huh-based Idempotent Installer
│   ├── packages/          # Idempotent APT, Snap & Binary Toolchain Manager
│   ├── theme/             # Unified Global Theme Switcher & AI Theme Studio
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
* **OS:** Ubuntu 24.04 LTS or newer (24.04+, 26.04 LTS compatible)
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
  [  OK  ] Ubuntu Version           Ubuntu 24.04.4 LTS detected (Ubuntu 24.04+ LTS compatible)
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
  [  OK  ] Cloud AI Key Validation  Active Cloud AI provider configured and verified

◆ Voice Subsystem
  [  OK  ] Microphone Capture       ALSA/PipeWire audio capture utility available (@DEFAULT_SOURCE@)

Diagnostic Summary: 13 Passed | 0 Failures
```

---

## ⌨️ Daily Workflow & Keybindings

| Shortcut | Action | Description |
| :--- | :--- | :--- |
| **`Super + Space`** | **Vula Floating HUD** | Raycast-style launcher, action palette, OS control (`do`), and AI assistant |
| **`Super + Alt + A`** | **Active Voice AI & OS Control** | Conversational assistant: speak questions or OS commands, receive spoken response |
| **`Super + Alt + V`** | **Voice Dictation** | Transcribe spoken words directly into the active editor or input field |
| **`Super + Alt + C`** | **AI Selection Explain** | Analyze & refactor active selection/clipboard with local AI & desktop notification |
| **`Super + Left / Right / Up / Down`** | **Half-Screen Snap** | Snap window to left, right, top, or bottom screen half |
| **`Super + Alt + U / I / J / K`** | **Quarter-Screen Snap** | Snap window to Top-Left, Top-Right, Bottom-Left, or Bottom-Right |
| **`Super + T`** | **Toggle Auto-Tile** | Enable / disable automatic window tiling mode |
| **`Super + Alt + H / J / K / L`** | **Workspace Navigation** | Switch workspaces dynamically using Vim keys |
| **`Super + Shift + H / J / K / L`** | **Move Window** | Move focused window to adjacent workspace |
| **`Super + H`** | **Minimize Window** | Hide/minimize currently focused window |
| **`Super + D`** | **Toggle Desktop** | Show or hide all open desktop windows |
| **`Super + Q`** | **Close Window** | Close focused application window |
| **`Super + M`** | **Maximize** | Toggle window maximization |

---

## ⚡ Vula OS Actions Engine (`vula do`)

Control your desktop environment natively with natural language instructions:

```bash
# Adjust screen brightness and audio volume:
vula do "ajusta el brillo al 40%"
vula do "pon el volumen al 70%"

# Launch developer applications:
vula do "abre VS Code"
vula do "abre Obsidian y el navegador"

# Manage files & directories securely:
vula do "mueve ~/Descargas/reporte.pdf a ~/Documentos/"

# Change visual themes & desktop settings:
vula do "cambia al tema tokyonight"
vula do "activa la luz nocturna"

# Inspect audit logs of all executed actions:
vula actions history
```

### 🛡️ Multi-Tier Risk Guardrails & Safety
* **`SAFE` (Auto-execute):** Launching apps, searching, adjusting brightness, volume, and theme.
* **`SENSITIVE` (TUI Confirmation):** Moving/copying files, creating directories, and terminating processes (`pkill`). Displays an interactive Charm `huh` prompt before proceeding.
* **`HIGH` (Strict Confirmation & Path Restriction):** File deletions (`rm`), shutdown, reboot, or operations targeting system root directories (`/`, `/boot`, `/etc`, `/usr`, etc.).

---

## 🖥️ Floating HUD Overlay (`Super + Space`)

The Vula HUD (`Super + Space`) is an interactive floating launcher with dynamic mode badges:

* **Tab Navigation:** Press **`Tab`** to cycle between modes: `[COMANDOS]`, `[⚡ ACCIÓN OS]`, `[🤖 IA CHAT]`, `[🎙 VOZ]`, and `[🎨 TEMA]`.
* **Direct OS Control:** Type `do [instruction]` directly in the HUD input line or select **`⚡ Execute OS Action (do)`** from the menu.
* **Instant Mode Switching:** Type `do ...` for OS control or `?...` for general AI chat.

---

## 🎙️ Active Voice AI & Audio Subsystem

### 1. Active Conversational Voice & OS Control Assistant
Press **`Super + Alt + A`** or run:
```bash
vula listen          # Default 4 seconds
vula listen 6        # Listen for 6 seconds
```
* **Audio Capture:** PipeWire `@DEFAULT_SOURCE@` native recording.
* **STT:** Transcribes speech locally via Whisper.cpp in Spanish/English (with silence hallucination filtering).
* **OS Actions Integration:** Speaks commands directly (*"Vula, ajusta el brillo al 40%"* or *"Abre VS Code"*) to execute system actions hands-free.
* **Dual Output:** Speaks answers back through your speakers via **Piper TTS** (<0.8s) and displays a desktop notification.

### 2. Microphone Diagnostic Test
Test live microphone input signal levels (in dB):
```bash
vula voice test
```

### 3. Voice Dictation into Any App
Press **`Super + Alt + V`** or run:
```bash
vula voice record
```

---

## 🤖 Hybrid AI Engine & Cloud Offloading

Vula features a hybrid AI engine combining local Ollama models with Cloud AI providers:

### 1. Configure Cloud Providers with Live Key Validation
Configure API keys with instant HTTP live validation checks:
```bash
# Configure Groq:
vula ai config --key=groq:gsk_YOUR_KEY

# Configure Google Gemini:
vula ai config --key=gemini:YOUR_KEY --cloud-model=gemini-2.0-flash

# Set hybrid mode and cloud token threshold:
vula ai config --mode=hybrid --threshold=1200
```

### 2. Live AI Status & Diagnostics
```bash
vula ai status
```

### 3. Developer Workflow Commands
```bash
# Smart Git Commit Generator:
vula ai commit

# Terminal Error Diagnoser:
vula ai fix "command not found: zoxide"

# Selection & Clipboard AI Explain (Super + Alt + C):
vula ai explain

# Ask AI with Desktop Context:
vula ai ask "How do I optimize this SQL query in my clipboard?"
```

---

## 🎨 Global Theme Engine & Theme Studio (`vula theme`)

### 1. Switch Themes via HUD (`Super + Space`) or CLI
```bash
# List available palettes:
vula theme list

# Switch theme instantly:
vula theme set tokyonight    # Tokyo Night (Default)
vula theme set catppuccin    # Catppuccin Mocha
vula theme set nord          # Nord Arctic
vula theme set rose-pine     # Rosé Pine
```

### 2. Generate Themes with Local AI
```bash
vula theme generate "cyberpunk neon obsidian with emerald and violet accents"
```

---

## 🖼 HD Wallpaper Engine (`vula wallpaper`)

```bash
vula wallpaper next          # Rotate wallpaper
vula wallpaper fetch         # Download 4K aesthetic wallpapers
```

---

## 🪟 Tiling Manager & Extensions (`vula desktop`)

```bash
vula desktop tiling          # Configure 6px gaps, active border, & snapping
vula desktop extensions install-curated   # Blur my Shell & Just Perfection
```

---

## 📁 Developer Dotfiles (`vula dotfiles`)

```bash
vula dotfiles install
```
* **Fish Shell (`~/.config/fish/config.fish`):** Aliases (`ls -> eza`, `cat -> bat`, `lg -> lazygit`, `ask -> vula ai ask`, `do -> vula do`).
* **Starship Prompt (`~/.config/starship.toml`)**
* **Neovim Lua IDE (`~/.config/nvim/init.lua`)**
* **Tmux (`~/.tmux.conf`)**

---

## 📦 Developer App Store (`vula apps`)

```bash
vula apps ui                 # Visual Charm Huh App Store
vula apps install-cli        # Install CLI stack (eza, bat, lazygit, starship, zoxide, btop, fzf)
```

---

## 🛡 Security Guarantees

* **Zero Root Daemons:** AI, Voice, and OS action engines run strictly in user space (`$USER`).
* **Human-in-the-Loop:** Sensitive & High-risk OS actions (`rm`, file moves, process termination, shutdown) require explicit TUI user confirmation.
* **System Boundaries:** Path sanitization prevents accidental modification of critical system root directories (`/`, `/boot`, `/etc`, `/usr`).
* **Local-First Privacy:** Prompts and audio streams never leave your device unless you explicitly configure an external API provider.

---

## 🗺 Roadmap

- [x] **v0.1.0-alpha:** Core Go + Charm CLI, `vula doctor`, Declarative YAML configs, Floating HUD, Local Ollama integration, Whisper STT + Piper TTS pipeline.
- [x] **v0.2.0-alpha:** Conversational Active Voice AI (`vula listen`), `Super + Alt + A` global shortcut, System notification integration.
- [x] **v0.3.0-alpha:** Global Theme Engine (`tokyonight`, `catppuccin`, `nord`, `rose-pine`), Dotfiles manager (Fish, Starship, Neovim, Tmux), Developer CLI app catalog.
- [x] **v0.4.0-alpha:** Tiling Assistant with custom gaps & active border highlight, Curated GNOME Extensions API installer (Blur my Shell, Just Perfection), Interactive Theme Studio, AI Palette Generator & HUD Theme Selector.
- [x] **v0.5.0-alpha:** Vula OS Actions Engine (`vula do`), Multi-tier risk guardrails & audit logging, Hybrid AI Engine (Groq, Gemini, OpenAI) with Live Key Validation, PipeWire `@DEFAULT_SOURCE@` microphone diagnostics (`vula voice test`), and HUD Tab mode cycling.
- [ ] **v1.0.0:** Production Release, `.deb` packaging, PPA & Sigstore release signing.

---

## 📄 License & Contributing

Distributed under the **MIT License**. See [LICENSE](LICENSE) and [SECURITY.md](SECURITY.md) for more details. Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).
