package apps

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/packages"
)

type AppRecipe struct {
	ID          string
	Name        string
	Category    string // CLI, Editor, Browser, Database, Productivity
	Description string
	InstallType string // apt, curl-sh, flatpak, github-tar
	PackageName string
}

var Catalog = []AppRecipe{
	// CLI Essentials
	{ID: "starship", Name: "Starship Prompt", Category: "CLI", Description: "Blazing-fast, customizable cross-shell prompt", InstallType: "curl-sh", PackageName: "https://starship.rs/install.sh"},
	{ID: "zoxide", Name: "Zoxide", Category: "CLI", Description: "Smarter cd command inspired by z and autojump", InstallType: "apt", PackageName: "zoxide"},
	{ID: "btop", Name: "Btop", Category: "CLI", Description: "Modern resource monitor with graphs and process manager", InstallType: "apt", PackageName: "btop"},
	{ID: "fzf", Name: "FZF", Category: "CLI", Description: "Command-line fuzzy finder", InstallType: "apt", PackageName: "fzf"},
	{ID: "ripgrep", Name: "Ripgrep", Category: "CLI", Description: "Ultra-fast line-oriented search tool", InstallType: "apt", PackageName: "ripgrep"},
	{ID: "neovim", Name: "Neovim", Category: "CLI", Description: "Vim-fork focused on extensibility and usability", InstallType: "apt", PackageName: "neovim"},
	{ID: "fish", Name: "Fish Shell", Category: "CLI", Description: "Smart and user-friendly command-line shell", InstallType: "apt", PackageName: "fish"},
	{ID: "tmux", Name: "Tmux", Category: "CLI", Description: "Terminal multiplexer with custom Vula theme", InstallType: "apt", PackageName: "tmux"},

	// Developer GUI Apps
	{ID: "vscode", Name: "Visual Studio Code", Category: "Editor", Description: "Code editing redefined", InstallType: "snap", PackageName: "code --classic"},
	{ID: "obsidian", Name: "Obsidian", Category: "Productivity", Description: "Knowledge base and markdown notes", InstallType: "snap", PackageName: "obsidian --classic"},
	{ID: "dbeaver", Name: "DBeaver Community", Category: "Database", Description: "Universal database tool and SQL client", InstallType: "snap", PackageName: "dbeaver-ce"},
	{ID: "brave", Name: "Brave Browser", Category: "Browser", Description: "Privacy-focused web browser", InstallType: "snap", PackageName: "brave"},
}

type Manager struct {
	cfg *config.Config
	pkg *packages.Manager
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
		pkg: packages.NewManager(),
	}
}

// InstallCLIStack installs the complete modern Unix developer CLI stack
func (m *Manager) InstallCLIStack() error {
	aptPkgs := []string{"zoxide", "btop", "fzf", "ripgrep", "fd-find", "bat", "neovim", "fish", "tmux", "xclip", "wl-clipboard"}
	if err := m.pkg.InstallAptPackages(aptPkgs); err != nil {
		return fmt.Errorf("failed installing APT packages: %w", err)
	}

	// Install Starship if missing
	if _, err := exec.LookPath("starship"); err != nil {
		cmd := exec.Command("sh", "-c", "curl -sS https://starship.rs/install.sh | sh -s -- -y -b "+os.Getenv("HOME")+"/.local/bin")
		_ = cmd.Run()
	}

	// Install Eza (modern ls) if missing
	if _, err := exec.LookPath("eza"); err != nil {
		installEza()
	}

	// Install Lazygit if missing
	if _, err := exec.LookPath("lazygit"); err != nil {
		installLazygit()
	}

	return nil
}

func installEza() {
	cmd := exec.Command("sh", "-c", `
		sudo mkdir -p /etc/apt/keyrings
		wget -qO- https://raw.githubusercontent.com/eza-community/eza/main/deb.asc | sudo gpg --dearmor -o /etc/apt/keyrings/gierens.gpg 2>/dev/null || true
		echo "deb [signed-by=/etc/apt/keyrings/gierens.gpg] http://deb.gierens.de stable main" | sudo tee /etc/apt/sources.list.d/gierens.list >/dev/null
		sudo apt-get update -qq && sudo apt-get install -y -qq eza
	`)
	_ = cmd.Run()
}

func installLazygit() {
	home := os.Getenv("HOME")
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
		cd /tmp && \
		LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*') && \
		curl -sLo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz" && \
		tar -xzf lazygit.tar.gz lazygit && \
		install lazygit %s/.local/bin/ && \
		rm -f lazygit.tar.gz lazygit
	`, home))
	_ = cmd.Run()
}
