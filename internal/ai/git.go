package ai

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vula-os/vula/internal/config"
)

// GenerateCommitMessage analyzes git diff and queries local AI for a conventional commit message
func GenerateCommitMessage(ctx context.Context, cfg *config.Config) (string, error) {
	// 1. Check staged changes first
	out, err := exec.CommandContext(ctx, "git", "diff", "--staged").Output()
	diff := strings.TrimSpace(string(out))

	// If no staged changes, fallback to unstaged diff
	if diff == "" {
		out, err = exec.CommandContext(ctx, "git", "diff").Output()
		diff = strings.TrimSpace(string(out))
	}

	if diff == "" {
		return "", fmt.Errorf("no staged or unstaged changes found in current repository")
	}

	// Truncate ultra-large diffs for prompt window safety
	runes := []rune(diff)
	if len(runes) > 3000 {
		diff = string(runes[:3000]) + "\n... [diff truncated]"
	}

	client := NewClient(cfg)
	prompt := fmt.Sprintf(
		"You are a git expert. Write a concise, precise conventional commit message (e.g. 'feat(scope): add feature' or 'fix(scope): resolve bug') for the following git diff.\nOUTPUT ONLY THE SINGLE COMMIT MESSAGE LINE. DO NOT INCLUDE MARKDOWN OR EXPLANATIONS.\n\nDiff:\n%s",
		diff,
	)

	resp, err := client.Ask(ctx, prompt, nil)
	if err != nil {
		return "", fmt.Errorf("AI model query failed: %w", err)
	}

	cleaned := strings.TrimSpace(resp)
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.Trim(cleaned, "\"")
	cleaned = strings.TrimSpace(cleaned)

	// Ensure single line
	lines := strings.Split(cleaned, "\n")
	if len(lines) > 0 {
		cleaned = strings.TrimSpace(lines[0])
	}

	return cleaned, nil
}

// ExecuteGitCommit runs git commit -m "message"
func ExecuteGitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}
