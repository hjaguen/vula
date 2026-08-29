package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.Theme.Palette != "tokyonight" {
		t.Errorf("expected default theme tokyonight, got %s", cfg.Theme.Palette)
	}

	if cfg.AI.DefaultProvider != "ollama" {
		t.Errorf("expected default AI provider ollama, got %s", cfg.AI.DefaultProvider)
	}

	if !cfg.Security.ConfirmDestructive {
		t.Error("expected ConfirmDestructive to be true by default for security")
	}
}
