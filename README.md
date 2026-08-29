# ⚡ Vula (v0.1.0-alpha)

> **"Vula"** *(verb, Zulu & Xhosa)*: **To open, to start, to unlock.**

**Vula** is an opinionated, keyboard-first developer operating environment crafted for **Ubuntu 24.04 LTS (Noble Numbat)**. It combines the rock-solid hardware stability of Ubuntu and GNOME Shell with the elegance of a **100% Go + Charm TUI** orchestrator, native local AI intelligence, real-time voice interaction, and tiling window management.

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

**Vula solves this with four core pillars:**

1. **Rock-Solid Foundation:** Built on top of Ubuntu 24.04 LTS and GNOME 46+. No kernel breakages, full NVIDIA/Intel/AMD driver support, and business-grade stability.
2. **100% Go & Charm TUI:** No fragile 5,000-line bash scripts. Vula is powered by a compiled, type-safe Go binary utilizing the [Charm](https://charm.sh) ecosystem (`bubbletea`, `lipgloss`, `huh`, `bubbles`).
3. **Local-First AI Intelligence:** Integrated directly into the OS with active desktop context awareness (reads window state and clipboard safely) powered by Ollama (`qwen2.5-coder`, `llama3.2`, `deepseek-r1`).
4. **Zero-Latency Voice Subsystem:** Push-to-talk voice dictation and command processing utilizing local `whisper.cpp` (STT) and `Piper` (neural TTS).

---

## 🏗 System Architecture

```text
vula/
├── cmd/vula/              # Cobra CLI Entrypoint & Subcommands
├── internal/
│   ├── ui/                # Charm Lipgloss Palettes & Banners
│   ├── config/            # Declarative YAML Config (~/.config/vula/config.yaml)
│   ├── doctor/            # Self-Healing System Diagnostic Engine
│   ├── installer/         # Interactive Huh-based Idempotent Installer
│   ├── hud/               # Floating Bubbletea Launcher & AI HUD
│   ├── ai/                # Context-Aware Ollama & Cloud LLM Router
│   ├── voice/             # STT (Whisper) & Neural TTS (Piper) Engine
│   ├── gnome/             # Dconf, Keybindings & Extensions Manager
│   └── packages/          # Idempotent APT, Flatpak, Mise & SHA-256 Verifier
├── scripts/
│   └── bootstrap.sh       # Safe, Non-Destructive Bootstrap Script
├── Makefile               # Standard Development Targets
├── SECURITY.md            # Security & Vulnerability Policy
└── README.md
```

---

## 🚀 Quickstart

### 1. Requirements
* **OS:** Ubuntu 24.04 LTS (Noble Numbat)
* **Desktop:** GNOME Shell (Wayland or X11)
* **Toolchain:** Go 1.24+ (auto-installed if missing)

### 2. One-Line Installation (Safe Bootstrap)
```bash
git clone https://github.com/vula-os/vula.git ~/repos/vula
cd ~/repos/vula
make install
```

### 3. Verify System Health
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
  [  OK  ] Display Server           Wayland / X11 active
  [  OK  ] dconf CLI                dconf binary available

◆ AI Subsystem
  [  OK  ] Ollama Local Daemon      Connected to Ollama at http://localhost:11434

◆ Voice Subsystem
  [  OK  ] Microphone Capture       Audio capture utility available
```

---

## ⌨️ Daily Workflow & Keybindings

| Shortcut | Action | Description |
| :--- | :--- | :--- |
| **`Super + Space`** | **Vula Floating HUD** | Raycast-style launcher, action palette, and AI assistant |
| **`Super + Alt + V`** | **Voice Dictate** | Transcribe voice directly into the active window |
| **`Super + Q`** | **Close Window** | Close focused application window |
| **`Super + M`** | **Maximize** | Toggle window maximization |
| **`Super + Return`**| **Terminal** | Open fast developer terminal (Ghostty / Kitty) |

---

## 🤖 AI & Voice Command Line

### Ask AI with Desktop Context
```bash
# Query the local LLM with active window & clipboard context:
vula ai ask "How do I optimize this SQL query in my clipboard?"
```

### Natural Language Shell Commands
```bash
# Translate intent directly to a safe bash command:
vula ai cmd "find all files larger than 100MB modified in the last 7 days"
```

### List Installed Local Models
```bash
vula ai models
```

### Voice Speech Synthesis & Dictation
```bash
# Synthesize neural speech locally:
vula voice speak "Compilación completada exitosamente."

# Record voice and inject into active window:
vula voice record
```

---

## 🛡 Security Guarantees

* **Zero Root Daemons:** AI and Voice engines run strictly in user space (`$USER`).
* **Human-in-the-Loop:** Destructive commands (`rm`, `sudo`, `dd`) always require explicit user confirmation.
* **Supply-Chain Verification:** All downloaded binaries and neural models are verified against cryptographic SHA-256 checksums.
* **Local-First Privacy:** Prompts and audio streams never leave your device unless you explicitly configure an external API provider.

---

## 🗺 Roadmap

- [x] **v0.1.0-alpha:** Core Go + Charm CLI, `vula doctor`, Declarative YAML configs, Floating HUD prototype, Local Ollama integration, Voice engine architecture.
- [ ] **v0.2.0:** Streaming Whisper.cpp ONNX real-time voice transcription, Piper Spanish/English models auto-downloader.
- [ ] **v0.3.0:** GNOME Shell custom extension for native Floating HUD hotkey overlay.
- [ ] **v0.4.0:** Dotfiles bundle (Ghostty + Neovim + Fish + Starship + Mise presets).
- [ ] **v1.0.0:** Production Release, `.deb` packaging, PPA & Sigstore release signing.

---

## 📄 License & Contributing

Distributed under the **MIT License**. See [LICENSE](LICENSE) and [SECURITY.md](SECURITY.md) for more details. Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).
