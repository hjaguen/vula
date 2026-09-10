package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatOptions struct {
	NumThread int `json:"num_thread,omitempty"`
	NumCtx    int `json:"num_ctx,omitempty"`
}

type ChatRequest struct {
	Model    string       `json:"model"`
	Messages []Message    `json:"messages"`
	Options  *ChatOptions `json:"options,omitempty"`
	Stream   bool         `json:"stream"`
}

type ChatResponseChunk struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

type ModelInfo struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
}

type TagsResponse struct {
	Models []ModelInfo `json:"models"`
}

type Client struct {
	cfg            *config.Config
	hybridProvider *HybridProvider
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:            cfg,
		hybridProvider: NewHybridProvider(cfg),
	}
}

// Ask sends a prompt using the Hybrid Provider system with general task type
func (c *Client) Ask(ctx context.Context, prompt string, streamHandler func(chunk string)) (string, error) {
	return c.AskTask(ctx, prompt, "general", streamHandler)
}

// AskTask sends a prompt with explicit task classification for complexity offloading
func (c *Client) AskTask(ctx context.Context, prompt string, taskType string, streamHandler func(chunk string)) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	var systemContent = c.cfg.AI.SystemPrompt
	if c.cfg.AI.ContextEnabled {
		ctxInfo := CaptureActiveContext()
		if ctxInfo != "" {
			systemContent += fmt.Sprintf("\n[Current User Desktop Context:\n%s]", ctxInfo)
		}
	}

	return c.hybridProvider.AskWithTaskType(ctx, prompt, systemContent, taskType, streamHandler)
}

// ValidateActiveCloudKey performs a live HTTP check against the configured Cloud AI provider
func (c *Client) ValidateActiveCloudKey(ctx context.Context) error {
	return c.hybridProvider.ValidateActiveCloudKey(ctx)
}

// ListLocalModels returns all available local models from Ollama
func (c *Client) ListLocalModels(ctx context.Context) ([]ModelInfo, error) {
	ollama := NewOllamaProvider(c.cfg)
	if !ollama.IsAvailable(ctx) {
		return nil, fmt.Errorf("local Ollama instance is not running")
	}

	host := c.cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}

	ollamaProvider := NewOllamaProvider(c.cfg)
	req, err := ollamaProvider.httpClient.Get(host + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var tags TagsResponse
	if err := json.NewDecoder(req.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags.Models, nil
}

// SuggestCommand translates a natural language request into a safe bash command
func (c *Client) SuggestCommand(ctx context.Context, task string) (string, error) {
	prompt := fmt.Sprintf("Return ONLY the precise single bash command line to achieve this task: \"%s\". Do NOT include markdown code blocks or explanations, just the command itself.", task)
	resp, err := c.AskTask(ctx, prompt, "suggest_command", nil)
	if err != nil {
		return "", err
	}
	cleaned := strings.TrimSpace(resp)
	cleaned = strings.TrimPrefix(cleaned, "```bash")
	cleaned = strings.TrimPrefix(cleaned, "```sh")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	return strings.TrimSpace(cleaned), nil
}

// CaptureActiveContext gets the clipboard and active window title
func CaptureActiveContext() string {
	var sb strings.Builder

	// 1. Get active window title
	if out, err := exec.Command("xdotool", "getactivewindow", "getwindowname").Output(); err == nil {
		title := strings.TrimSpace(string(out))
		if title != "" {
			sb.WriteString(fmt.Sprintf("- Active Window: %s\n", title))
		}
	}

	// 2. Get clipboard content (truncated for safety)
	var clipText string
	if out, err := exec.Command("wl-paste", "--no-newline").Output(); err == nil {
		clipText = string(out)
	} else if out, err := exec.Command("xclip", "-o", "-selection", "clipboard").Output(); err == nil {
		clipText = string(out)
	}

	if clipText != "" {
		runes := []rune(clipText)
		if len(runes) > 500 {
			clipText = string(runes[:500]) + "... [truncated]"
		}
		sb.WriteString(fmt.Sprintf("- Clipboard Content: %s\n", clipText))
	}

	return sb.String()
}
