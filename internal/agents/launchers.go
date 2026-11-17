package agents

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AgentLauncher defines a supported AI coding agent harness
type AgentLauncher struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Package     string   `json:"package"`
	InstallCmd  string   `json:"install_cmd"`
	RunCmd      string   `json:"run_cmd"`
}

// BuiltinAgents returns the curated registry of AI agent harnesses in Vula
var BuiltinAgents = []AgentLauncher{
	{
		ID:          "claude",
		Name:        "Claude Code",
		Command:     "claude",
		Aliases:     []string{"cx", "claude-code"},
		Description: "Anthropic's official agentic terminal coding assistant",
		Package:     "@anthropic-ai/claude-code",
		InstallCmd:  "npm install -g @anthropic-ai/claude-code",
		RunCmd:      "claude",
	},
	{
		ID:          "opencode",
		Name:        "OpenCode",
		Command:     "opencode",
		Aliases:     []string{"c"},
		Description: "AI coding agent supporting Gemini, OpenAI, Claude and local LLMs",
		Package:     "opencode-ai",
		InstallCmd:  "npm install -g opencode-ai",
		RunCmd:      "opencode",
	},
	{
		ID:          "gemini",
		Name:        "Gemini CLI",
		Command:     "gemini",
		Aliases:     []string{"gemini-cli"},
		Description: "Official Google Gemini terminal assistant & code generator",
		Package:     "@google/gemini-cli",
		InstallCmd:  "npm install -g @google/gemini-cli",
		RunCmd:      "gemini",
	},
	{
		ID:          "aider",
		Name:        "Aider Chat",
		Command:     "aider",
		Aliases:     []string{"ai-coder"},
		Description: "Pair-programming AI agent in terminal with Git integration",
		Package:     "aider-chat",
		InstallCmd:  "pipx install aider-chat || python3 -m pip install -U aider-chat",
		RunCmd:      "aider",
	},
	{
		ID:          "codex",
		Name:        "Codex CLI",
		Command:     "codex",
		Aliases:     []string{"openai-codex"},
		Description: "OpenAI developer CLI agent for repository operations",
		Package:     "@openai/codex",
		InstallCmd:  "npm install -g @openai/codex",
		RunCmd:      "codex",
	},
}

// FindAgent locates an agent by ID or alias
func FindAgent(query string) *AgentLauncher {
	q := strings.ToLower(strings.TrimSpace(query))
	for _, a := range BuiltinAgents {
		if strings.ToLower(a.ID) == q || strings.ToLower(a.Command) == q {
			return &a
		}
		for _, alias := range a.Aliases {
			if strings.ToLower(alias) == q {
				return &a
			}
		}
	}
	return nil
}

// EnsureInstalled checks if the agent binary exists; if not, attempts automatic installation
func (a *AgentLauncher) EnsureInstalled() error {
	if _, err := exec.LookPath(a.Command); err == nil {
		return nil
	}

	// Not installed: attempt install command
	fmt.Printf("📦 Agent '%s' (%s) not found. Installing via: %s...\n", a.Name, a.Command, a.InstallCmd)
	cmd := exec.Command("sh", "-c", a.InstallCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed installing %s: %w", a.Name, err)
	}

	return nil
}

// GenerateLazyStub creates an executable stub in ~/.local/bin if the binary is missing
func (a *AgentLauncher) GenerateLazyStub(binDir string) error {
	target := filepath.Join(binDir, a.Command)
	if _, err := os.Stat(target); err == nil {
		// Already exists
		return nil
	}

	stubContent := fmt.Sprintf(`#!/usr/bin/env bash
# Vula Lazy-Loaded AI Agent Launcher: %s
if ! command -v %s &>/dev/null; then
    echo "⚡ [Vula] Initializing %s (%s)..."
    %s
fi
exec %s "$@"
`, a.Name, a.Command, a.Name, a.Package, a.InstallCmd, a.RunCmd)

	if err := os.WriteFile(target, []byte(stubContent), 0755); err != nil {
		return err
	}

	// Also link primary aliases (e.g. 'c' for opencode, 'cx' for claude)
	for _, alias := range a.Aliases {
		aliasTarget := filepath.Join(binDir, alias)
		if _, err := os.Stat(aliasTarget); os.IsNotExist(err) {
			_ = os.Symlink(target, aliasTarget)
		}
	}

	return nil
}

// InstallAllLazyStubs installs lightweight stubs in ~/.local/bin for all supported agents
func InstallAllLazyStubs() error {
	home := os.Getenv("HOME")
	binDir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	for _, a := range BuiltinAgents {
		_ = a.GenerateLazyStub(binDir)
	}
	return nil
}
