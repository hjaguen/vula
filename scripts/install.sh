#!/usr/bin/env bash
# ==============================================================================
# VULA OS — UNIVERSAL ONE-LINER INSTALLER
# Safe, idempotent setup for Ubuntu 24.04 LTS (x86_64 / arm64)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/hjaguen/vula/main/scripts/install.sh | bash
# ==============================================================================

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${PURPLE}"
cat << "EOF"
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
EOF
echo -e "${NC}"
echo -e "${CYAN}${BOLD}⚡ VULA OS — The Malleable, Keyboard-First AI Developer Operating Environment${NC}\n"

# 1. OS Validation
if [ ! -f /etc/os-release ]; then
    echo -e "${RED}Error: Cannot detect operating system release file (/etc/os-release).${NC}"
    exit 1
fi

. /etc/os-release
if [[ "${ID:-}" != "ubuntu" && "${ID_LIKE:-}" != *"ubuntu"* && "${ID_LIKE:-}" != *"debian"* ]]; then
    echo -e "${YELLOW}Warning: Vula is officially tailored for Ubuntu 24.04+ LTS. Proceeding anyway...${NC}"
fi

# 2. Architecture Detection
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  VULA_ARCH="amd64" ;;
    aarch64) VULA_ARCH="arm64" ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH (supported: x86_64, aarch64)${NC}"
        exit 1
        ;;
esac

# 3. Ensure Local Binaries in PATH
BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    export PATH="$BIN_DIR:$PATH"
    for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
        if [ -f "$rc" ] && ! grep -q 'PATH="$HOME/.local/bin:$PATH"' "$rc"; then
            echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$rc"
        fi
    done
fi

VULA_BIN="$BIN_DIR/vula"
REPO_OWNER="hjaguen"
REPO_NAME="vula"
INSTALLED_VIA_RELEASE=false

echo -e "${CYAN}→ Checking latest official Vula release on GitHub...${NC}"

# 4. Attempt download from GitHub Releases
LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)

if [ -n "$LATEST_TAG" ]; then
    echo -e "${GREEN}✓ Found release: ${LATEST_TAG} (${VULA_ARCH})${NC}"
    TAR_NAME="vula_${LATEST_TAG#v}_linux_${VULA_ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_TAG}/${TAR_NAME}"
    
    TMP_DIR=$(mktemp -d)
    if curl -sL --fail "$DOWNLOAD_URL" -o "$TMP_DIR/$TAR_NAME" 2>/dev/null; then
        echo -e "${CYAN}→ Extracting prebuilt binary...${NC}"
        tar -xzf "$TMP_DIR/$TAR_NAME" -C "$TMP_DIR"
        if [ -f "$TMP_DIR/vula" ]; then
            mv "$TMP_DIR/vula" "$VULA_BIN"
            chmod +x "$VULA_BIN"
            INSTALLED_VIA_RELEASE=true
        fi
    fi
    rm -rf "$TMP_DIR"
fi

# 5. Fallback: Build from Source if prebuilt release not available
if [ "$INSTALLED_VIA_RELEASE" = false ]; then
    echo -e "${YELLOW}Notice: Prebuilt asset not found for ${LATEST_TAG:-latest}. Falling back to source build...${NC}"
    
    # Check or install Go toolchain in user space
    if ! command -v go &> /dev/null; then
        echo -e "${CYAN}→ Installing Go toolchain locally into ~/.local/go...${NC}"
        GO_VERSION="1.24.0"
        GO_TAR="go${GO_VERSION}.linux-${VULA_ARCH}.tar.gz"
        mkdir -p "$HOME/.local"
        curl -fsSL "https://go.dev/dl/${GO_TAR}" | tar -xz -C "$HOME/.local"
        ln -sf "$HOME/.local/go/bin/go" "$BIN_DIR/go"
        ln -sf "$HOME/.local/go/bin/gofmt" "$BIN_DIR/gofmt"
        export PATH="$BIN_DIR:$HOME/.local/go/bin:$PATH"
    fi

    # Clone or update repository
    SRC_DIR="${VULA_SRC_DIR:-$HOME/repos/vula}"
    if [ ! -d "$SRC_DIR" ]; then
        echo -e "${CYAN}→ Cloning repository into $SRC_DIR...${NC}"
        mkdir -p "$(dirname "$SRC_DIR")"
        git clone "https://github.com/${REPO_OWNER}/${REPO_NAME}.git" "$SRC_DIR"
    fi

    echo -e "${CYAN}→ Compiling Vula binary from source...${NC}"
    cd "$SRC_DIR"
    go build -ldflags="-s -w" -o "$VULA_BIN" ./cmd/vula
    chmod +x "$VULA_BIN"
    echo -e "${GREEN}✓ Vula successfully compiled at $VULA_BIN${NC}"
fi

# 6. Safety Snapshot: Backup initial system state
echo -e "\n${CYAN}→ Creating initial safety snapshot with 'vula backup'...${NC}"
"$VULA_BIN" backup create "pre-vula-install" "Initial system snapshot prior to Vula OS configuration" || true

# 7. Install Universal Dotfiles & Shell Aliases
echo -e "${CYAN}→ Configuring universal shell dotfiles and aliases (do, c, cx, hud, fleet)...${NC}"
"$VULA_BIN" dotfiles install || true

# 8. Generate Agent Skill
echo -e "${CYAN}→ Generating AI Agent Skill (~/.config/vula/SKILL.md)...${NC}"
"$VULA_BIN" skill >/dev/null 2>&1 || true

# 9. Run Diagnostics
echo -e "\n${CYAN}→ Running Vula System Health Diagnostics...${NC}"
"$VULA_BIN" doctor || true

echo -e "\n${GREEN}${BOLD}══════════════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}🎉 VULA OS IS READY!${NC}"
echo -e "${GREEN}${BOLD}══════════════════════════════════════════════════════════════════════════${NC}"
echo -e "  • ${BOLD}HUD Launcher:${NC} Presiona ${CYAN}Super + Space${NC} o escribe ${CYAN}hud${NC}"
echo -e "  • ${BOLD}Acción Natural:${NC} ${CYAN}do \"ajusta el brillo al 50%\"${NC}"
echo -e "  • ${BOLD}Agentes IA:${NC} ${CYAN}c${NC} (OpenCode) o ${CYAN}cx${NC} (Claude Code)"
echo -e "  • ${BOLD}Modo Zen:${NC} ${CYAN}vula desktop focus${NC}"
echo -e "  • ${BOLD}Actualizar:${NC} ${CYAN}vula update${NC}"
echo -e "  • ${BOLD}Restaurar:${NC} ${CYAN}vula backup restore pre-vula-install${NC}"
echo -e "\n${PURPLE}Reinicia tu terminal o ejecuta: source ~/.bashrc${NC}\n"
