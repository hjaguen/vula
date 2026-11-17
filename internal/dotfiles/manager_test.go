package dotfiles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vula-os/vula/internal/config"
)

func TestDotfileConfigs(t *testing.T) {
	// 1. Universal Shell
	u := UniversalShellConfig()
	if !strings.Contains(u, "alias do='vula do'") {
		t.Errorf("UniversalShellConfig missing alias 'do'")
	}
	if !strings.Contains(u, "alias c='vula agents run opencode'") {
		t.Errorf("UniversalShellConfig missing alias 'c'")
	}

	// 2. Fish
	f := FishConfig()
	if !strings.Contains(f, "fish_greeting") {
		t.Errorf("FishConfig missing fish_greeting")
	}

	// 3. Starship
	s := StarshipConfig()
	if !strings.Contains(s, "[character]") {
		t.Errorf("StarshipConfig missing character module")
	}

	// 4. Tmux
	m := TmuxConfig()
	if !strings.Contains(m, "prefix C-a") {
		t.Errorf("TmuxConfig missing prefix C-a")
	}
}

func TestInstallAllDotfiles(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "vula-dotfiles-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	// Create dummy .bashrc
	bashrcPath := filepath.Join(tmpHome, ".bashrc")
	_ = os.WriteFile(bashrcPath, []byte("# existing bashrc"), 0644)

	cfg := config.DefaultConfig()
	mgr := NewManager(cfg)

	if err := mgr.InstallAllDotfiles(); err != nil {
		t.Fatalf("InstallAllDotfiles failed: %v", err)
	}

	// Check shell.sh exists
	shellShPath := filepath.Join(tmpHome, ".config", "vula", "shell.sh")
	if _, err := os.Stat(shellShPath); err != nil {
		t.Errorf("expected %s to exist", shellShPath)
	}

	// Check bashrc was injected
	bashrcContent, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("failed to read .bashrc: %v", err)
	}
	if !strings.Contains(string(bashrcContent), "vula/shell.sh") {
		t.Errorf(".bashrc did not contain vula/shell.sh injection")
	}
}
