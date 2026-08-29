package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/config"
)

type Engine struct {
	cfg *config.Config
}

func NewEngine(cfg *config.Config) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

// RecordAudio captures audio from the default microphone for a given duration or until stopped
func (e *Engine) RecordAudio(ctx context.Context, outputFile string, maxDuration time.Duration) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 16kHz, 16-bit Mono WAV format optimal for Whisper
	var cmd *exec.Cmd
	if _, err := exec.LookPath("arecord"); err == nil {
		cmd = exec.CommandContext(ctx, "arecord", "-D", "default", "-f", "S16_LE", "-r", "16000", "-c", "1", "-d", fmt.Sprintf("%d", int(maxDuration.Seconds())), outputFile)
	} else if _, err := exec.LookPath("ffmpeg"); err == nil {
		cmd = exec.CommandContext(ctx, "ffmpeg", "-y", "-f", "pulse", "-i", "default", "-ar", "16000", "-ac", "1", "-t", fmt.Sprintf("%d", int(maxDuration.Seconds())), outputFile)
	} else {
		return fmt.Errorf("no audio recorder (arecord or ffmpeg) found in PATH")
	}

	return cmd.Run()
}

// Transcribe converts recorded WAV audio file into text using Whisper
func (e *Engine) Transcribe(ctx context.Context, audioFile string) (string, error) {
	// Check for local whisper.cpp / whisper binary
	if whisperPath, err := exec.LookPath("whisper-cpp"); err == nil {
		modelPath := filepath.Join(os.Getenv("HOME"), ".local/share/vula/models", "ggml-base.bin")
		cmd := exec.CommandContext(ctx, whisperPath, "-m", modelPath, "-f", audioFile, "--no-timestamps", "-l", "auto")
		home := os.Getenv("HOME")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("LD_LIBRARY_PATH=%s/.local/lib:%s", home, os.Getenv("LD_LIBRARY_PATH")),
			fmt.Sprintf("PATH=%s/.local/bin:%s", home, os.Getenv("PATH")),
		)
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("whisper-cpp failed: %w", err)
		}
		return strings.TrimSpace(string(out)), nil
	}

	// Fallback to python whisper or whisper CLI
	if whisperPath, err := exec.LookPath("whisper"); err == nil {
		cmd := exec.CommandContext(ctx, whisperPath, audioFile, "--model", e.cfg.Voice.STTModel, "--output_format", "txt")
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("whisper CLI failed: %w", err)
		}
		return strings.TrimSpace(string(out)), nil
	}

	return "", fmt.Errorf("whisper binary not found. Run 'vula install --module voice' to setup the voice stack")
}

// Speak synthesizes text to speech using Piper neural TTS
func (e *Engine) Speak(ctx context.Context, text string) error {
	piperPath, err := exec.LookPath("piper")
	if err != nil {
		// Fallback to speech-dispatcher / espeak if piper is not installed yet
		if espeak, err := exec.LookPath("espeak"); err == nil {
			cmd := exec.CommandContext(ctx, espeak, text)
			return cmd.Run()
		}
		return fmt.Errorf("piper TTS binary not found")
	}

	voiceModel := filepath.Join(os.Getenv("HOME"), ".local/share/vula/voices", e.cfg.Voice.TTSVoice+".onnx")

	// Piper pipeline: echo text | piper --model voice.onnx --output-raw | aplay -r 22050 -f S16_LE -t raw
	piperCmd := exec.CommandContext(ctx, piperPath, "--model", voiceModel, "--output-raw")
	home := os.Getenv("HOME")
	piperCmd.Env = append(os.Environ(),
		fmt.Sprintf("LD_LIBRARY_PATH=%s/.local/lib:%s", home, os.Getenv("LD_LIBRARY_PATH")),
		fmt.Sprintf("PATH=%s/.local/bin:%s", home, os.Getenv("PATH")),
	)
	piperCmd.Stdin = strings.NewReader(text)

	aplayCmd := exec.CommandContext(ctx, "aplay", "-r", "22050", "-f", "S16_LE", "-t", "raw", "-")
	aplayCmd.Stdin, _ = piperCmd.StdoutPipe()

	if err := aplayCmd.Start(); err != nil {
		return err
	}
	if err := piperCmd.Run(); err != nil {
		return err
	}
	return aplayCmd.Wait()
}

// TypeIntoActiveWindow types text into the currently focused window using ydotool or xdotool
func TypeIntoActiveWindow(text string) error {
	if _, err := exec.LookPath("ydotool"); err == nil {
		cmd := exec.Command("ydotool", "type", text)
		return cmd.Run()
	}
	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "type", "--delay", "10", text)
		return cmd.Run()
	}
	return fmt.Errorf("neither ydotool nor xdotool is installed")
}
