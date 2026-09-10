package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/actions"
	"github.com/vula-os/vula/internal/ai"
	"github.com/vula-os/vula/internal/config"
)

type Assistant struct {
	cfg         *config.Config
	voiceEngine *Engine
	aiClient    *ai.Client
}

func NewAssistant(cfg *config.Config) *Assistant {
	return &Assistant{
		cfg:         cfg,
		voiceEngine: NewEngine(cfg),
		aiClient:    ai.NewClient(cfg),
	}
}

// ListenAndRespond triggers active speech capture, queries local AI or OS Actions, and speaks back the answer
func (a *Assistant) ListenAndRespond(ctx context.Context, durationSec int) (string, string, error) {
	if durationSec <= 0 {
		durationSec = 4
	}

	tempAudio := filepath.Join(os.TempDir(), fmt.Sprintf("vula_listen_%d.wav", time.Now().UnixNano()))
	defer os.Remove(tempAudio)

	// Send visual system notification that listening has started
	_ = exec.Command("notify-send", "-a", "Vula AI", "-i", "audio-input-microphone", "⚡ Vula AI Escuchando...", "Habla ahora con tu asistente").Start()

	// 1. Record voice audio
	recordCtx, cancel := context.WithTimeout(ctx, time.Duration(durationSec+10)*time.Second)
	defer cancel()

	if err := a.voiceEngine.RecordAudio(recordCtx, tempAudio, time.Duration(durationSec)*time.Second); err != nil {
		return "", "", fmt.Errorf("error grabando audio: %w", err)
	}

	// 2. Transcribe voice audio with Whisper
	transcription, err := a.voiceEngine.Transcribe(ctx, tempAudio)
	if err != nil {
		return "", "", fmt.Errorf("error en transcripción: %w", err)
	}

	transcription = strings.TrimSpace(transcription)
	if transcription == "" {
		_ = exec.Command("notify-send", "-a", "Vula AI", "Vula AI", "No se detectó audio de voz en los segundos de grabación.").Start()
		return "", "", fmt.Errorf("no se detectó audio de voz hablada. Habla cerca al micrófono o especifica más tiempo (ej. 'vula listen 6')")
	}

	// Send visual notification that AI is computing response
	_ = exec.Command("notify-send", "-a", "Vula AI", fmt.Sprintf("🎙 \"%s\"", transcription), "Procesando intención...").Start()

	// 3. Check if transcription is an OS action request
	lowerText := strings.ToLower(transcription)
	isActionCandidate := strings.HasPrefix(lowerText, "vula do") ||
		strings.HasPrefix(lowerText, "haz ") ||
		strings.HasPrefix(lowerText, "hazme ") ||
		strings.Contains(lowerText, "brillo") ||
		strings.Contains(lowerText, "volumen") ||
		strings.Contains(lowerText, "abre ") ||
		strings.Contains(lowerText, "cierra ") ||
		strings.Contains(lowerText, "bloquea") ||
		strings.Contains(lowerText, "mueve ")

	if isActionCandidate {
		actionText := strings.TrimPrefix(transcription, "vula do")
		actionText = strings.TrimPrefix(actionText, "vula")
		actionText = strings.TrimSpace(actionText)
		if actionText == "" {
			actionText = transcription
		}

		plan, planErr := actions.ParseIntent(ctx, a.aiClient, actionText)
		if planErr == nil && plan != nil && len(plan.Actions) > 0 {
			_ = exec.Command("notify-send", "-a", "Vula OS Control", "-i", "dialog-information", fmt.Sprintf("⚡ Acción: %s", plan.Summary), "Ejecutando acción del sistema...").Start()
			execErr := actions.ExecutePlan(ctx, plan, a.cfg)
			if execErr == nil {
				respMsg := fmt.Sprintf("Listo, he ejecutado la acción: %s", plan.Summary)
				_ = exec.Command("notify-send", "-a", "Vula OS Control", "-i", "dialog-information", "✓ Acción Ejecutada", respMsg).Start()
				_ = a.voiceEngine.Speak(ctx, respMsg)
				return transcription, respMsg, nil
			}
		}
	}

	// 4. Query local AI model with desktop context for general questions
	aiPrompt := fmt.Sprintf("El usuario te dice o pregunta por voz: \"%s\". Responde de manera conversacional, muy concisa, precisa y amigable en 1 o 2 oraciones en español.", transcription)
	aiResponse, err := a.aiClient.Ask(ctx, aiPrompt, nil)
	if err != nil {
		return transcription, "", fmt.Errorf("error en respuesta de IA: %w", err)
	}

	aiResponse = strings.TrimSpace(aiResponse)

	// Send visual notification with the full generated answer
	_ = exec.Command("notify-send", "-a", "Vula AI", "-i", "dialog-information", fmt.Sprintf("⚡ Vula: %s", transcription), aiResponse).Start()

	// 5. Speak response back through speakers with Piper
	_ = a.voiceEngine.Speak(ctx, aiResponse)

	return transcription, aiResponse, nil
}
