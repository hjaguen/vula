package ai

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vula-os/vula/internal/config"
)

// ExplainSelection reads active window/clipboard content and provides an AI explanation & refactoring
func ExplainSelection(ctx context.Context, cfg *config.Config) (string, error) {
	ctxInfo := CaptureActiveContext()
	if ctxInfo == "" {
		_ = exec.Command("notify-send", "-a", "Vula AI", "Vula AI", "No hay texto o código en el portapapeles.").Start()
		return "", fmt.Errorf("no active window title or clipboard text found")
	}

	_ = exec.Command("notify-send", "-a", "Vula AI", "-i", "dialog-information", "⚡ Vula AI Explicando...", "Analizando código/texto del portapapeles").Start()

	client := NewClient(cfg)
	prompt := fmt.Sprintf(
		"Analiza el siguiente texto o código seleccionado por el usuario. Explica de forma muy clara, concisa y profesional en español en 2-3 oraciones qué hace, y si es código, proporciona la versión optimizada o refactorizada.\n\nContexto:\n%s",
		ctxInfo,
	)

	resp, err := client.Ask(ctx, prompt, nil)
	if err != nil {
		_ = exec.Command("notify-send", "-a", "Vula AI", "Error Vula AI", err.Error()).Start()
		return "", fmt.Errorf("AI explanation failed: %w", err)
	}

	cleaned := strings.TrimSpace(resp)
	summaryNotif := cleaned
	if len([]rune(summaryNotif)) > 140 {
		summaryNotif = string([]rune(summaryNotif)[:140]) + "..."
	}

	_ = exec.Command("notify-send", "-a", "Vula AI", "-i", "accessories-text-editor", "⚡ Vula AI Explicación", summaryNotif).Start()

	return cleaned, nil
}
