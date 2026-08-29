#!/usr/bin/env bash
# ==============================================================================
# VULA BOOTSTRAP INSTALLER
# Safe, idempotent initialization script for Ubuntu 24.04 LTS
# ==============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
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
echo -e "${CYAN}⚡ VULA — Next-Gen AI & Voice Developer OS for Ubuntu 24.04 LTS${NC}\n"

# 1. Verify OS
if [ ! -f /etc/os-release ]; then
    echo -e "${RED}Error: Cannot detect operating system release file.${NC}"
    exit 1
fi

. /etc/os-release
if [[ "${ID:-}" != "ubuntu" && "${ID_LIKE:-}" != *"ubuntu"* && "${ID_LIKE:-}" != *"debian"* ]]; then
    echo -e "${RED}Warning: Vula is designed and optimized for Ubuntu 24.04 LTS.${NC}"
fi

# 2. Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${CYAN}Installing Go toolchain in user space (~/.local/go)...${NC}"
    mkdir -p "$HOME/.local/bin" "$HOME/.local/go"
    GO_VERSION="1.24.0"
    ARCH="amd64"
    if [ "$(uname -m)" = "aarch64" ]; then ARCH="arm64"; fi
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" | tar -xz -C "$HOME/.local/"
    ln -sf "$HOME/.local/go/bin/go" "$HOME/.local/bin/go"
    ln -sf "$HOME/.local/go/bin/gofmt" "$HOME/.local/bin/gofmt"
    export PATH="$HOME/.local/bin:$HOME/.local/go/bin:$PATH"
fi

# 3. Build & Launch Vula CLI
VULA_DIR="${VULA_DIR:-$HOME/repos/vula}"
if [ -d "$VULA_DIR" ]; then
    cd "$VULA_DIR"
    echo -e "${CYAN}Building Vula binary...${NC}"
    go build -ldflags="-s -w" -o "$HOME/.local/bin/vula" ./cmd/vula
    echo -e "${GREEN}✓ Vula CLI installed at ~/.local/bin/vula${NC}\n"
    
    # Launch doctor
    "$HOME/.local/bin/vula" doctor
else
    echo -e "${RED}Error: Repository path $VULA_DIR not found.${NC}"
    exit 1
fi
