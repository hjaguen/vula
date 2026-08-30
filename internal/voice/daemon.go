package voice

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/vula-os/vula/internal/config"
)

// StartVoiceDaemon launches continuous hands-free voice assistant loop
func StartVoiceDaemon(ctx context.Context, cfg *config.Config) error {
	assistant := NewAssistant(cfg)

	_ = exec.Command("notify-send", "-a", "Vula Voice Daemon", "-i", "audio-input-microphone", "⚡ Vula Voice Daemon Activo", "Escuchando en segundo plano...").Start()
	fmt.Println("⚡ Vula Voice Daemon running in background... Press Ctrl+C to stop.")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Listen in 4-second audio frames
			q, ans, err := assistant.ListenAndRespond(ctx, 4)
			if err == nil && q != "" && ans != "" {
				fmt.Printf("\n[Daemon] 🎙 Question: %s\n[Daemon] ⚡ Answer: %s\n", q, ans)
			}
			// Brief sleep to avoid CPU spinning
			time.Sleep(500 * time.Millisecond)
		}
	}
}
