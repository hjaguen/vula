package ai

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vula-os/vula/internal/config"
)

// DiagnoseTerminalError analyzes terminal error tracebacks and recommends corrective commands
func DiagnoseTerminalError(ctx context.Context, cfg *config.Config, input string) (string, error) {
	errorText := strings.TrimSpace(input)

	// If no arg provided, try reading from stdin (piped input)
	if errorText == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err == nil {
				errorText = strings.TrimSpace(string(data))
			}
		}
	}

	// Fallback to clipboard if stdin is empty
	if errorText == "" {
		ctxInfo := CaptureActiveContext()
		if ctxInfo != "" {
			errorText = ctxInfo
		}
	}

	if errorText == "" {
		return "", fmt.Errorf("no error message provided via argument, pipe, or clipboard")
	}

	client := NewClient(cfg)
	prompt := fmt.Sprintf(
		"You are a Linux and terminal debugging expert. Analyze the following terminal command failure or error trace.\n"+
			"1. Briefly state the ROOT CAUSE in 1 sentence.\n"+
			"2. Provide the EXACT SHELL COMMAND line to resolve it under 'Suggested Fix:'.\n\nError output:\n%s",
		errorText,
	)

	resp, err := client.Ask(ctx, prompt, nil)
	if err != nil {
		return "", fmt.Errorf("AI diagnosis failed: %w", err)
	}

	return strings.TrimSpace(resp), nil
}
