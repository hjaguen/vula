package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string          `yaml:"version"`
	Theme    ThemeConfig     `yaml:"theme"`
	AI       AIConfig        `yaml:"ai"`
	Voice    VoiceConfig     `yaml:"voice"`
	Desktop  DesktopConfig   `yaml:"desktop"`
	DevTools DevToolsConfig  `yaml:"devtools"`
	Security SecurityConfig  `yaml:"security"`
}

type ThemeConfig struct {
	Palette     string `yaml:"palette"`     // catppuccin, tokyonight, nord, rose-pine
	DarkMode    bool   `yaml:"dark_mode"`
	FontFamily  string `yaml:"font_family"`  // JetBrains Mono, Fira Code, GeekMono
	FontSize    int    `yaml:"font_size"`
	AccentColor string `yaml:"accent_color"` // hex code or preset name
}

type AIConfig struct {
	DefaultProvider string                 `yaml:"default_provider"` // ollama, gemini, openai, anthropic
	OllamaHost      string                 `yaml:"ollama_host"`      // http://localhost:11434
	DefaultModel    string                 `yaml:"default_model"`    // qwen2.5:1.5b, qwen2.5-coder:1.5b, llama3.2
	NumThreads      int                    `yaml:"num_threads"`      // CPU inference threads (e.g. 4)
	ContextLength   int                    `yaml:"context_length"`   // Context window tokens (e.g. 4096)
	ContextEnabled  bool                   `yaml:"context_enabled"`  // captures active window & clipboard
	Streaming       bool                   `yaml:"streaming"`
	SystemPrompt    string                 `yaml:"system_prompt"`
	APIKeys         map[string]string      `yaml:"api_keys,omitempty"`
}

type VoiceConfig struct {
	Enabled         bool    `yaml:"enabled"`
	Engine          string  `yaml:"engine"`          // whisper.cpp, faster-whisper, sherpa-onnx
	WakeWord        string  `yaml:"wake_word"`       // "hey vula", "vula"
	WakeWordEnabled bool    `yaml:"wakeword_enabled"`
	PushToTalkKey   string  `yaml:"push_to_talk_key"`// <Super>Space, <Super><Alt>v
	STTModel        string  `yaml:"stt_model"`       // tiny, base, small
	TTSVoice        string  `yaml:"tts_voice"`       // es_ES-davefx-medium, en_US-lessac-medium
	AudioDevice     string  `yaml:"audio_device"`    // default or specific pulse/pipewire device
	EnergyThreshold float64 `yaml:"energy_threshold"`
}

type DesktopConfig struct {
	TilingEngine    string `yaml:"tiling_engine"`    // forge, pop-shell, native
	GapsInner       int    `yaml:"gaps_inner"`
	GapsOuter       int    `yaml:"gaps_outer"`
	TopBarCompact   bool   `yaml:"top_bar_compact"`
	BlurEnabled     bool   `yaml:"blur_enabled"`
	GlobalHUDHotkey string `yaml:"global_hud_hotkey"`// <Super>space
}

type DevToolsConfig struct {
	DefaultTerminal string   `yaml:"default_terminal"` // ghostty, kitty, alacritty
	DefaultShell    string   `yaml:"default_shell"`    // fish, zsh, bash
	DefaultEditor   string   `yaml:"default_editor"`   // nvim, code, zed
	MisePlugins     []string `yaml:"mise_plugins"`     // node, python, go, rust, ruby
}

type SecurityConfig struct {
	SandboxAI           bool `yaml:"sandbox_ai"`
	ConfirmDestructive  bool `yaml:"confirm_destructive_cmds"` // require prompt before executing rm, sudo, etc.
	AuditLogging        bool `yaml:"audit_logging"`
	AllowClipboardRead  bool `yaml:"allow_clipboard_read"`
}

func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Theme: ThemeConfig{
			Palette:     "tokyonight",
			DarkMode:    true,
			FontFamily:  "JetBrainsMono Nerd Font",
			FontSize:    12,
			AccentColor: "#7C3AED",
		},
		AI: AIConfig{
			DefaultProvider: "ollama",
			OllamaHost:      "http://localhost:11434",
			DefaultModel:    "qwen2.5-coder:1.5b",
			NumThreads:      4,
			ContextLength:   4096,
			ContextEnabled:  true,
			Streaming:       true,
			SystemPrompt:    "You are Vula, an intelligent, concise and hyper-competent developer OS assistant integrated into Ubuntu. Help with code, terminal commands, and system operations safely.",
			APIKeys:         make(map[string]string),
		},
		Voice: VoiceConfig{
			Enabled:         true,
			Engine:          "whisper.cpp",
			WakeWord:        "hey vula",
			WakeWordEnabled: false, // push-to-talk default for battery & privacy
			PushToTalkKey:   "<Super><Alt>v",
			STTModel:        "base",
			TTSVoice:        "es_ES-davefx-medium",
			AudioDevice:     "default",
			EnergyThreshold: 0.05,
		},
		Desktop: DesktopConfig{
			TilingEngine:    "forge",
			GapsInner:       8,
			GapsOuter:       8,
			TopBarCompact:   true,
			BlurEnabled:     true,
			GlobalHUDHotkey: "<Super>space",
		},
		DevTools: DevToolsConfig{
			DefaultTerminal: "ghostty",
			DefaultShell:    "fish",
			DefaultEditor:   "nvim",
			MisePlugins:     []string{"node@lts", "python@latest", "go@latest", "rust@latest"},
		},
		Security: SecurityConfig{
			SandboxAI:          true,
			ConfirmDestructive: true,
			AuditLogging:       true,
			AllowClipboardRead: true,
		},
	}
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "vula"), nil
}

func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func LoadConfig() (*Config, error) {
	cfgFile, err := ConfigFile()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			_ = SaveConfig(cfg)
			return cfg, nil
		}
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", cfgFile, err)
	}

	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	cfgFile := filepath.Join(dir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
