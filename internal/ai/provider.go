package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/config"
)

type AIProvider interface {
	Ask(ctx context.Context, prompt string, systemPrompt string, streamHandler func(chunk string)) (string, error)
	Name() string
	IsAvailable(ctx context.Context) bool
}

// ----------------------------------------------------------------------------
// 1. Local Ollama Provider
// ----------------------------------------------------------------------------

type OllamaProvider struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewOllamaProvider(cfg *config.Config) *OllamaProvider {
	return &OllamaProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 0, // context-managed
		},
	}
}

func (p *OllamaProvider) Name() string {
	return "Ollama (Local)"
}

func (p *OllamaProvider) IsAvailable(ctx context.Context) bool {
	host := p.cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", host+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (p *OllamaProvider) Ask(ctx context.Context, prompt string, systemPrompt string, streamHandler func(chunk string)) (string, error) {
	host := p.cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}

	model := p.cfg.AI.DefaultModel
	if model == "" {
		model = "qwen2.5-coder:1.5b"
	}

	// Auto-resolve to an installed local model if configured model is not pulled
	if tagsReq, err := p.httpClient.Get(host + "/api/tags"); err == nil {
		var tags TagsResponse
		if json.NewDecoder(tagsReq.Body).Decode(&tags) == nil && len(tags.Models) > 0 {
			hasModel := false
			var firstValidModel string
			for _, m := range tags.Models {
				// Skip embedding models
				if strings.Contains(m.Name, "embed") {
					continue
				}
				if firstValidModel == "" {
					firstValidModel = m.Name
				}
				if m.Name == model || strings.HasPrefix(m.Name, model) || strings.HasPrefix(model, strings.Split(m.Name, ":")[0]) {
					hasModel = true
					model = m.Name
					break
				}
			}
			if !hasModel && firstValidModel != "" {
				model = firstValidModel
			}
		}
		_ = tagsReq.Body.Close()
	}

	var chatOpts *ChatOptions
	if p.cfg.AI.NumThreads > 0 || p.cfg.AI.ContextLength > 0 {
		chatOpts = &ChatOptions{
			NumThread: p.cfg.AI.NumThreads,
			NumCtx:    p.cfg.AI.ContextLength,
		}
	}

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Options: chatOpts,
		Stream:  true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", host+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("local Ollama unreachable at %s: %w", host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama HTTP %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var chunk ChatResponseChunk
		if err := json.Unmarshal(line, &chunk); err == nil && chunk.Message.Content != "" {
			fullResponse.WriteString(chunk.Message.Content)
			if streamHandler != nil {
				streamHandler(chunk.Message.Content)
			}
		}
	}

	return fullResponse.String(), scanner.Err()
}

// ----------------------------------------------------------------------------
// 2. Google Gemini Cloud Provider
// ----------------------------------------------------------------------------

type GeminiProvider struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewGeminiProvider(cfg *config.Config) *GeminiProvider {
	return &GeminiProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 0,
		},
	}
}

func (p *GeminiProvider) Name() string {
	return "Google Gemini (Cloud)"
}

func (p *GeminiProvider) getAPIKey() string {
	if key, ok := p.cfg.AI.APIKeys["gemini"]; ok && key != "" {
		return key
	}
	return os.Getenv("GEMINI_API_KEY")
}

func (p *GeminiProvider) IsAvailable(ctx context.Context) bool {
	return p.getAPIKey() != ""
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiReq struct {
	Contents          []geminiContent `json:"contents"`
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiRespChunk struct {
	Candidates []geminiCandidate `json:"candidates"`
}

func (p *GeminiProvider) Ask(ctx context.Context, prompt string, systemPrompt string, streamHandler func(chunk string)) (string, error) {
	apiKey := p.getAPIKey()
	if apiKey == "" {
		return "", fmt.Errorf("Gemini API key missing. Set GEMINI_API_KEY environment variable or run 'vula ai config --key=gemini:YOUR_KEY'")
	}

	model := p.cfg.AI.CloudModel
	if model == "" {
		model = "gemini-2.0-flash"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?alt=sse&key=%s", model, apiKey)

	bodyObj := geminiReq{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}
	if systemPrompt != "" {
		bodyObj.SystemInstruction = &geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		}
	}

	jsonData, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini Cloud API connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API HTTP %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonStr := strings.TrimPrefix(line, "data: ")
			var chunk geminiRespChunk
			if err := json.Unmarshal([]byte(jsonStr), &chunk); err == nil && len(chunk.Candidates) > 0 {
				for _, part := range chunk.Candidates[0].Content.Parts {
					if part.Text != "" {
						fullResponse.WriteString(part.Text)
						if streamHandler != nil {
							streamHandler(part.Text)
						}
					}
				}
			}
		}
	}

	return fullResponse.String(), scanner.Err()
}

// ----------------------------------------------------------------------------
// 3. OpenAI / Groq / Ollama Cloud Provider
// ----------------------------------------------------------------------------

type OpenAIProvider struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewOpenAIProvider(cfg *config.Config) *OpenAIProvider {
	return &OpenAIProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 0,
		},
	}
}

func (p *OpenAIProvider) Name() string {
	switch p.cfg.AI.CloudProvider {
	case "groq":
		return "Groq (Cloud)"
	case "ollama-cloud":
		return "Ollama Cloud"
	default:
		return "OpenAI / Cloud"
	}
}

func (p *OpenAIProvider) getAPIKey() string {
	provider := p.cfg.AI.CloudProvider
	if key, ok := p.cfg.AI.APIKeys[provider]; ok && key != "" {
		return key
	}
	switch provider {
	case "groq":
		return os.Getenv("GROQ_API_KEY")
	case "ollama-cloud":
		return os.Getenv("OLLAMA_API_KEY")
	default:
		return os.Getenv("OPENAI_API_KEY")
	}
}

func (p *OpenAIProvider) getEndpoint() string {
	if p.cfg.AI.CloudHost != "" {
		return p.cfg.AI.CloudHost
	}
	switch p.cfg.AI.CloudProvider {
	case "groq":
		return "https://api.groq.com/openai/v1/chat/completions"
	case "ollama-cloud":
		return "https://api.ollama.com/v1/chat/completions"
	default:
		return "https://api.openai.com/v1/chat/completions"
	}
}

func (p *OpenAIProvider) IsAvailable(ctx context.Context) bool {
	return p.getAPIKey() != "" || p.cfg.AI.CloudHost != ""
}

type openAIReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type openAIChoice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
}

type openAIChunk struct {
	Choices []openAIChoice `json:"choices"`
}

func (p *OpenAIProvider) Ask(ctx context.Context, prompt string, systemPrompt string, streamHandler func(chunk string)) (string, error) {
	apiKey := p.getAPIKey()
	endpoint := p.getEndpoint()

	model := p.cfg.AI.CloudModel
	if model == "" {
		if p.cfg.AI.CloudProvider == "groq" {
			model = "llama-3.3-70b-versatile"
		} else {
			model = "gpt-4o-mini"
		}
	}

	bodyObj := openAIReq{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Cloud AI connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Cloud API HTTP %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonStr := strings.TrimPrefix(line, "data: ")
			if jsonStr == "[DONE]" {
				break
			}
			var chunk openAIChunk
			if err := json.Unmarshal([]byte(jsonStr), &chunk); err == nil && len(chunk.Choices) > 0 {
				txt := chunk.Choices[0].Delta.Content
				if txt != "" {
					fullResponse.WriteString(txt)
					if streamHandler != nil {
						streamHandler(txt)
					}
				}
			}
		}
	}

	return fullResponse.String(), scanner.Err()
}

// ----------------------------------------------------------------------------
// 4. Hybrid Provider (Smart Router with Complexity Offloading & Fallback)
// ----------------------------------------------------------------------------

type HybridProvider struct {
	cfg           *config.Config
	localProvider AIProvider
	cloudProvider AIProvider
}

func NewHybridProvider(cfg *config.Config) *HybridProvider {
	hp := &HybridProvider{
		cfg:           cfg,
		localProvider: NewOllamaProvider(cfg),
	}

	switch cfg.AI.CloudProvider {
	case "gemini":
		hp.cloudProvider = NewGeminiProvider(cfg)
	case "groq", "ollama-cloud", "openai":
		hp.cloudProvider = NewOpenAIProvider(cfg)
	default:
		hp.cloudProvider = NewGeminiProvider(cfg)
	}

	return hp
}

func (h *HybridProvider) Name() string {
	return "Hybrid (Local + Cloud Router)"
}

func (h *HybridProvider) IsAvailable(ctx context.Context) bool {
	return h.localProvider.IsAvailable(ctx) || (h.cloudProvider != nil && h.cloudProvider.IsAvailable(ctx))
}

func (h *HybridProvider) AskWithTaskType(ctx context.Context, prompt string, systemPrompt string, taskType string, streamHandler func(chunk string)) (string, error) {
	mode := h.cfg.AI.Mode
	if mode == "" {
		mode = "hybrid"
	}

	// 1. Force Local Mode
	if mode == "local" {
		return h.localProvider.Ask(ctx, prompt, systemPrompt, streamHandler)
	}

	// 2. Force Cloud Mode
	if mode == "cloud" {
		if h.cloudProvider == nil || !h.cloudProvider.IsAvailable(ctx) {
			return "", fmt.Errorf("Cloud AI mode configured but no valid API key found. Set GEMINI_API_KEY or run 'vula ai config --key=gemini:YOUR_KEY'")
		}
		notifyCloudUse("☁️ Vula AI (Cloud)")
		return h.cloudProvider.Ask(ctx, prompt, systemPrompt, streamHandler)
	}

	// 3. Hybrid Mode: Complexity-based Offloading & Fallback
	complexity := EvaluateComplexity(prompt, taskType, h.cfg.AI.HeavyTokenThreshold)

	// Offload heavy/complex tasks to Cloud AI if available
	if complexity == ComplexityHeavy && h.cloudProvider != nil && h.cloudProvider.IsAvailable(ctx) {
		notifyCloudUse("☁️ Vula AI (Cloud Offload - Heavy Task)")
		return h.cloudProvider.Ask(ctx, prompt, systemPrompt, streamHandler)
	}

	// Try Local AI for lightweight or standard tasks
	if h.localProvider.IsAvailable(ctx) {
		resp, err := h.localProvider.Ask(ctx, prompt, systemPrompt, streamHandler)
		if err == nil {
			return resp, nil
		}
		// If local execution failed mid-request, fallback to cloud if available
	}

	// Fallback to Cloud AI if local Ollama is offline or fails
	if h.cloudProvider != nil && h.cloudProvider.IsAvailable(ctx) {
		notifyCloudUse("⚡ Vula AI (Cloud Fallback - Local Down)")
		return h.cloudProvider.Ask(ctx, prompt, systemPrompt, streamHandler)
	}

	return "", fmt.Errorf("neither local Ollama nor Cloud AI is available. Start Ollama locally ('ollama serve') or add a cloud API key ('vula ai config --key=gemini:YOUR_KEY')")
}

func notifyCloudUse(title string) {
	_ = exec.Command("notify-send", "-a", "Vula AI", "-i", "network-transmit-receive", title, "Enviando consulta al modelo de la nube...").Start()
}
