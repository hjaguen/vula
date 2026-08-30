package doctor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/ui"
)

type CheckStatus int

const (
	StatusOK CheckStatus = iota
	StatusWarn
	StatusFail
)

type CheckResult struct {
	Category string
	Name     string
	Status   CheckStatus
	Message  string
	Hint     string
}

type Report struct {
	Results   []CheckResult
	Passed    int
	Warnings  int
	Failures  int
	Timestamp time.Time
}

func RunDiagnostics(cfg *config.Config) *Report {
	report := &Report{
		Results:   make([]CheckResult, 0),
		Timestamp: time.Now(),
	}

	// 1. Operating System Checks
	checkOS(report)

	// 2. Desktop Environment Checks
	checkDesktop(report)

	// 3. Audio & Hardware Checks
	checkAudioAndHardware(report)

	// 4. Core Development Tooling
	checkDevTools(report)

	// 5. AI Subsystem Checks
	checkAI(report, cfg)

	// 6. Voice Subsystem Checks
	checkVoice(report, cfg)

	// Compute totals
	for _, res := range report.Results {
		switch res.Status {
		case StatusOK:
			report.Passed++
		case StatusWarn:
			report.Warnings++
		case StatusFail:
			report.Failures++
		}
	}

	return report
}

func checkOS(r *Report) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		r.Results = append(r.Results, CheckResult{
			Category: "Operating System",
			Name:     "Distribution",
			Status:   StatusFail,
			Message:  "Unable to read /etc/os-release",
			Hint:     "Ensure you are running on a Debian/Ubuntu-based distribution",
		})
		return
	}

	content := string(data)
	if strings.Contains(content, "Ubuntu") || strings.Contains(content, "ubuntu") || strings.Contains(content, "debian") {
		// Extract VERSION_ID if available
		var versionId, prettyName string
		for _, line := range strings.Split(content, "\n") {
			if strings.HasPrefix(line, "VERSION_ID=") {
				versionId = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				prettyName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
		}

		if prettyName == "" {
			prettyName = "Ubuntu Linux"
		}

		// Support Ubuntu 24.04 LTS, 24.10, 26.04 LTS and beyond (>= 24.04)
		if versionId >= "24.04" || versionId == "" {
			r.Results = append(r.Results, CheckResult{
				Category: "Operating System",
				Name:     "Ubuntu Version",
				Status:   StatusOK,
				Message:  fmt.Sprintf("%s detected (Ubuntu 24.04+ LTS compatible)", prettyName),
			})
		} else {
			r.Results = append(r.Results, CheckResult{
				Category: "Operating System",
				Name:     "Ubuntu Version",
				Status:   StatusWarn,
				Message:  fmt.Sprintf("Older release detected (%s)", prettyName),
				Hint:     "Vula requires Ubuntu >= 24.04 LTS for PipeWire and GNOME 46+ compatibility",
			})
		}
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Operating System",
			Name:     "Distribution",
			Status:   StatusWarn,
			Message:  "Non-Ubuntu Linux system detected",
			Hint:     "Vula is built for Ubuntu 24.04+ LTS; some packages may require custom config",
		})
	}

	// Kernel version
	out, err := exec.Command("uname", "-r").Output()
	if err == nil {
		r.Results = append(r.Results, CheckResult{
			Category: "Operating System",
			Name:     "Linux Kernel",
			Status:   StatusOK,
			Message:  strings.TrimSpace(string(out)),
		})
	}
}

func checkDesktop(r *Report) {
	// Check GNOME Shell
	out, err := exec.Command("gnome-shell", "--version").Output()
	if err == nil {
		version := strings.TrimSpace(string(out))
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "GNOME Shell",
			Status:   StatusOK,
			Message:  version,
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "GNOME Shell",
			Status:   StatusFail,
			Message:  "GNOME Shell is not installed or not in PATH",
			Hint:     "Install GNOME Shell: sudo apt install gnome-shell",
		})
	}

	// Check Session Type (Wayland / X11)
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType == "" {
		out, err := exec.Command("sh", "-c", "echo $XDG_SESSION_TYPE").Output()
		if err == nil {
			sessionType = strings.TrimSpace(string(out))
		}
	}

	if sessionType == "wayland" {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "Display Server",
			Status:   StatusOK,
			Message:  "Wayland session active (High performance & modern gestures)",
		})
	} else if sessionType == "x11" {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "Display Server",
			Status:   StatusWarn,
			Message:  "X11 session active (Wayland recommended for best fluidity)",
			Hint:     "Switch to 'Ubuntu on Wayland' at login screen for optimal experience",
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "Display Server",
			Status:   StatusWarn,
			Message:  fmt.Sprintf("Unknown session type: %s", sessionType),
		})
	}

	// Check dconf tool
	if _, err := exec.LookPath("dconf"); err == nil {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "dconf CLI",
			Status:   StatusOK,
			Message:  "dconf binary available for declarative configuration",
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Desktop Environment",
			Name:     "dconf CLI",
			Status:   StatusFail,
			Message:  "dconf not found",
			Hint:     "Install dconf: sudo apt install dconf-cli",
		})
	}
}

func checkAudioAndHardware(r *Report) {
	// Audio check
	if _, err := exec.LookPath("pactl"); err == nil || checkPipewireRunning() {
		r.Results = append(r.Results, CheckResult{
			Category: "Hardware & Audio",
			Name:     "Audio Subsystem",
			Status:   StatusOK,
			Message:  "PipeWire / PulseAudio detected (Ready for voice STT/TTS)",
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Hardware & Audio",
			Name:     "Audio Subsystem",
			Status:   StatusWarn,
			Message:  "Neither PipeWire nor PulseAudio command line tools found",
			Hint:     "Install pipewire-audio: sudo apt install pipewire-audio-client-libraries",
		})
	}

	// GPU / Acceleration
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		out, _ := exec.Command("nvidia-smi", "--query-gpu=name,driver_version", "--format=csv,noheader").Output()
		r.Results = append(r.Results, CheckResult{
			Category: "Hardware & Audio",
			Name:     "GPU Acceleration",
			Status:   StatusOK,
			Message:  fmt.Sprintf("NVIDIA GPU: %s", strings.TrimSpace(string(out))),
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Hardware & Audio",
			Name:     "GPU Acceleration",
			Status:   StatusOK,
			Message:  "Integrated / Vulkan standard GPU (CPU/Vulkan inference mode)",
		})
	}
}

func checkDevTools(r *Report) {
	tools := []struct {
		name string
		pkg  string
		hint string
	}{
		{"git", "Git VCS", "sudo apt install git"},
		{"curl", "cURL HTTP client", "sudo apt install curl"},
		{"make", "GNU Make", "sudo apt install build-essential"},
		{"go", "Go Toolchain", "Run 'vula install --module devtools' or install Go"},
	}

	for _, t := range tools {
		if _, err := exec.LookPath(t.name); err == nil {
			r.Results = append(r.Results, CheckResult{
				Category: "Developer Tooling",
				Name:     t.pkg,
				Status:   StatusOK,
				Message:  fmt.Sprintf("%s found", t.name),
			})
		} else {
			r.Results = append(r.Results, CheckResult{
				Category: "Developer Tooling",
				Name:     t.pkg,
				Status:   StatusWarn,
				Message:  fmt.Sprintf("%s missing", t.name),
				Hint:     t.hint,
			})
		}
	}
}

func checkAI(r *Report, cfg *config.Config) {
	// Check Ollama service
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	host := cfg.AI.OllamaHost
	if host == "" {
		host = "http://localhost:11434"
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", host+"/api/tags", nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == http.StatusOK {
		_ = resp.Body.Close()
		r.Results = append(r.Results, CheckResult{
			Category: "AI Subsystem",
			Name:     "Ollama Local Daemon",
			Status:   StatusOK,
			Message:  fmt.Sprintf("Connected to Ollama at %s", host),
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "AI Subsystem",
			Name:     "Ollama Local Daemon",
			Status:   StatusWarn,
			Message:  fmt.Sprintf("Cannot reach Ollama at %s", host),
			Hint:     "Install & run Ollama with 'vula ai setup' or visit https://ollama.ai",
		})
	}
}

func checkVoice(r *Report, cfg *config.Config) {
	if !cfg.Voice.Enabled {
		r.Results = append(r.Results, CheckResult{
			Category: "Voice Subsystem",
			Name:     "Voice Engine",
			Status:   StatusOK,
			Message:  "Voice engine disabled in configuration",
		})
		return
	}

	// Check if arecord or ffmpeg or pipewire record exists
	if _, err := exec.LookPath("arecord"); err == nil {
		r.Results = append(r.Results, CheckResult{
			Category: "Voice Subsystem",
			Name:     "Microphone Capture",
			Status:   StatusOK,
			Message:  "ALSA audio capture utility available",
		})
	} else if _, err := exec.LookPath("ffmpeg"); err == nil {
		r.Results = append(r.Results, CheckResult{
			Category: "Voice Subsystem",
			Name:     "Microphone Capture",
			Status:   StatusOK,
			Message:  "FFmpeg audio stream utility available",
		})
	} else {
		r.Results = append(r.Results, CheckResult{
			Category: "Voice Subsystem",
			Name:     "Microphone Capture",
			Status:   StatusWarn,
			Message:  "No standard audio capture utility (arecord/ffmpeg) found",
			Hint:     "Install alsa-utils: sudo apt install alsa-utils",
		})
	}
}

func checkPipewireRunning() bool {
	out, err := exec.Command("pgrep", "-x", "pipewire").Output()
	return err == nil && len(out) > 0
}

func (r *Report) Render() string {
	b := strings.Builder{}
	b.WriteString(ui.RenderHeader("Doctor Diagnostics", "Checking system readiness, desktop environment, and AI stack"))

	currentCategory := ""
	for _, res := range r.Results {
		if res.Category != currentCategory {
			currentCategory = res.Category
			b.WriteString(lipgloss.NewStyle().
				Bold(true).
				Foreground(ui.SecondaryColor).
				MarginTop(1).
				MarginBottom(0).
				Render(fmt.Sprintf("◆ %s", currentCategory)))
			b.WriteString("\n")
		}

		var badge string
		switch res.Status {
		case StatusOK:
			badge = ui.SuccessBadge.Render("  OK  ")
		case StatusWarn:
			badge = ui.WarnBadge.Render(" WARN ")
		case StatusFail:
			badge = ui.ErrorBadge.Render(" FAIL ")
		}

		line := fmt.Sprintf("  %s %-24s %s", badge, lipgloss.NewStyle().Bold(true).Render(res.Name), res.Message)
		b.WriteString(line)
		b.WriteString("\n")

		if res.Hint != "" && res.Status != StatusOK {
			hintStyle := lipgloss.NewStyle().Foreground(ui.MutedColor).PaddingLeft(12)
			b.WriteString(hintStyle.Render(fmt.Sprintf("↳ Hint: %s\n", res.Hint)))
		}
	}

	b.WriteString("\n")
	summaryBox := ui.CardStyle.Render(fmt.Sprintf(
		"Diagnostic Summary: %s Passed | %s Warnings | %s Failures",
		ui.SuccessStyle.Render(fmt.Sprintf("%d", r.Passed)),
		ui.WarnStyle.Render(fmt.Sprintf("%d", r.Warnings)),
		ui.ErrorStyle.Render(fmt.Sprintf("%d", r.Failures)),
	))
	b.WriteString(summaryBox)
	b.WriteString("\n")

	return b.String()
}
