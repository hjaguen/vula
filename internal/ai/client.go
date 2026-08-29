package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	cfg        *config.Config
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Ask sends a prompt and streams back chunks via the callback
func (c *Client) Ask(ctx context.Context, prompt string, streamHandler func(chunk string)) (string, error) {
	host := c.cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}

	model := c.cfg.AI.DefaultModel
	if model == "" {
		model = "qwen2.5-coder:7b"
	}

	// Try resolving model if possible
	if models, err := c.ListLocalModels(ctx); err == nil && len(models) > 0 {
		hasModel := false
		for _, m := range models {
			if m.Name == model || strings.HasPrefix(m.Name, model) {
				hasModel = true
				break
			}
		}
		if !hasModel {
			// Fallback to first available model
			model = models[0].Name
		}
	}

	// Capture optional context if enabled
	var systemContent = c.cfg.AI.SystemPrompt
	if c.cfg.AI.ContextEnabled {
		ctxInfo := CaptureActiveContext()
		if ctxInfo != "" {
			systemContent += fmt.Sprintf("\n[Current User Desktop Context:\n%s]", ctxInfo)
		}
	}

	var chatOpts *ChatOptions
	if c.cfg.AI.NumThreads > 0 || c.cfg.AI.ContextLength > 0 {
		chatOpts = &ChatOptions{
			NumThread: c.cfg.AI.NumThreads,
			NumCtx:    c.cfg.AI.ContextLength,
		}
	}

	reqBody := ChatRequest{
		Model:    model,
		Messages: []Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: prompt},
		},
		Options:  chatOpts,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", host+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to reach Ollama at %s: %w", host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var chunk ChatResponseChunk
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		if chunk.Message.Content != "" {
			fullResponse.WriteString(chunk.Message.Content)
			if streamHandler != nil {
				streamHandler(chunk.Message.Content)
			}
		}

		if chunk.Done {
			break
		}
	}

	return fullResponse.String(), scanner.Err()
}

// ListLocalModels returns all available local models from Ollama
func (c *Client) ListLocalModels(ctx context.Context) ([]ModelInfo, error) {
	host := c.cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", host+"/api/tags", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d", resp.StatusCode)
	}

	var tags TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags.Models, nil
}

// SuggestCommand translates a natural language request into a safe bash command
func (c *Client) SuggestCommand(ctx context.Context, task string) (string, error) {
	prompt := fmt.Sprintf("Return ONLY the precise single bash command line to achieve this task: \"%s\". Do NOT include markdown code blocks or explanations, just the command itself.", task)
	resp, err := c.Ask(ctx, prompt, nil)
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
