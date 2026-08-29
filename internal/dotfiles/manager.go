package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vula-os/vula/internal/config"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// FishConfig generates an ergonomic, modern fish configuration
func FishConfig() string {
	return `# ==============================================================================
# VULA FISH CONFIGURATION
# ==============================================================================

# Suppress fish greeting
set -g fish_greeting ""

# Environment Paths
set -gx PATH $HOME/.local/bin $HOME/.local/go/bin $HOME/.cargo/bin $PATH
set -gx LD_LIBRARY_PATH $HOME/.local/lib $LD_LIBRARY_PATH
set -gx EDITOR nvim

# Modern CLI Aliases
if type -q eza
    alias ls='eza --icons'
    alias ll='eza -la --icons --git'
    alias la='eza -a --icons'
    alias lt='eza --tree --level=2 --icons'
end

if type -q bat
    alias cat='bat --paging=never'
else if type -q batcat
    alias cat='batcat --paging=never'
end

if type -q lazygit
    alias lg='lazygit'
end

if type -q lazydocker
    alias lzd='lazydocker'
end

# Vula AI Aliases
alias ask='vula ai ask'
alias vcmd='vula ai cmd'

# Zoxide initialization
if type -q zoxide
    zoxide init fish | source
end

# Starship Prompt initialization
if type -q starship
    starship init fish | source
end
`
}

// StarshipConfig generates a fast, sleek prompt configuration
func StarshipConfig() string {
	return `format = """
[┌─](bold #7C3AED)$directory$git_branch$git_status$nodejs$golang$rust$package
[└─>](bold #06B6D4) """

add_newline = false

[directory]
style = "bold #F8FAFC"
truncation_length = 3
truncation_symbol = "…/"

[git_branch]
symbol = " "
style = "bold #10B981"
format = "on [$symbol$branch]($style) "

[git_status]
style = "bold #F59E0B"
format = "([$all_status$ahead_behind]($style) )"

[nodejs]
symbol = " "
format = "[$symbol($version )]($style)"
style = "bold #10B981"

[golang]
symbol = " "
format = "[$symbol($version )]($style)"
style = "bold #06B6D4"

[rust]
symbol = " "
format = "[$symbol($version )]($style)"
style = "bold #EF4444"

[character]
success_symbol = "[⚡](bold #10B981)"
error_symbol = "[✗](bold #EF4444)"
`
}

// NeovimConfig generates a clean, modern Lua configuration
func NeovimConfig() string {
	return `-- ==============================================================================
-- VULA NEOVIM CONFIGURATION
-- ==============================================================================

local opt = vim.opt

-- UI & Editor ergonomics
opt.number = true
opt.relativenumber = true
opt.cursorline = true
opt.termguicolors = true
opt.signcolumn = "yes"
opt.scrolloff = 8

-- Indentation & Tabs
opt.tabstop = 4
opt.shiftwidth = 4
opt.expandtab = true
opt.smartindent = true

-- Search & System
opt.ignorecase = true
opt.smartcase = true
opt.clipboard = "unnamedplus"
opt.undofile = true
opt.updatetime = 250

-- Keybindings
vim.g.mapleader = " "
local map = vim.keymap.set

map("n", "<leader>w", ":w<CR>", { desc = "Save file" })
map("n", "<leader>q", ":q<CR>", { desc = "Quit" })
map("n", "<leader>h", ":nohlsearch<CR>", { desc = "Clear search highlight" })

-- Window navigation
map("n", "<C-h>", "<C-w>h", { desc = "Move to left window" })
map("n", "<C-j>", "<C-w>j", { desc = "Move to bottom window" })
map("n", "<C-k>", "<C-w>k", { desc = "Move to top window" })
map("n", "<C-l>", "<C-w>l", { desc = "Move to right window" })
`
}

// TmuxConfig generates sensible tmux defaults
func TmuxConfig() string {
	return `# ==============================================================================
# VULA TMUX CONFIGURATION
# ==============================================================================

# Remap prefix to Ctrl+A
unbind C-b
set -g prefix C-a
bind C-a send-prefix

# Enable mouse support & True Color
set -g mouse on
set -g default-terminal "screen-256color"
set-option -sa terminal-overrides ",xterm*:Tc"

# Vi mode for copy
set-window-option -g mode-keys vi
bind -T copy-mode-vi v send-keys -X begin-selection
bind -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "xclip -in -selection clipboard"

# Easy pane splitting with current directory
bind | split-window -h -c "#{pane_current_path}"
bind - split-window -v -c "#{pane_current_path}"
unbind '"'
unbind %

# Minimalist status bar
set -g status-position top
set -g status-style bg='#1E1E2E',fg='#CDD6F4'
set -g status-left ' #[bold,fg=#7C3AED]⚡ VULA #[default]| '
set -g status-right ' %H:%M | %d-%b-%y '
`
}

// InstallAllDotfiles safely backs up and writes dotfiles
func (m *Manager) InstallAllDotfiles() error {
	home := os.Getenv("HOME")
	backupDir := filepath.Join(home, ".config", "vula", "backups", fmt.Sprintf("backup_%d", time.Now().Unix()))

	files := []struct {
		destDir  string
		fileName string
		content  string
	}{
		{filepath.Join(home, ".config", "fish"), "config.fish", FishConfig()},
		{filepath.Join(home, ".config"), "starship.toml", StarshipConfig()},
		{filepath.Join(home, ".config", "nvim"), "init.lua", NeovimConfig()},
		{home, ".tmux.conf", TmuxConfig()},
	}

	for _, f := range files {
		if err := os.MkdirAll(f.destDir, 0755); err != nil {
			return err
		}

		targetPath := filepath.Join(f.destDir, f.fileName)
		// Backup if already exists
		if _, err := os.Stat(targetPath); err == nil {
			_ = os.MkdirAll(backupDir, 0755)
			_ = os.Rename(targetPath, filepath.Join(backupDir, f.fileName))
		}

		if err := os.WriteFile(targetPath, []byte(f.content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
	}

	return nil
}
