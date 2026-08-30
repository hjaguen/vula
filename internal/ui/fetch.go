package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vula-os/vula/internal/config"
)

const FetchLogoASCII = `
  ⚡ V U L A
 ▄▄   ▄▄ ▄▄  ▄▄ ▄▄   ▄▄▄▄ 
 ██   ██ ██  ██ ██  ██  ██
 ▀█▄ ▄█▀ ██  ██ ██  ██████
  ▀███▀  ▀████▀ ██▄▄██  ██
`

// RenderFetchCard generates a clean, Neofetch-style system summary with Vula metrics
func RenderFetchCard(cfg *config.Config) string {
	// 1. User & Host Header
	user := os.Getenv("USER")
	if user == "" {
		user = "vula"
	}
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "ubuntu"
	}
	headerText := fmt.Sprintf("%s@%s", user, hostname)

	headerStyle := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	sepStyle := lipgloss.NewStyle().Foreground(MutedColor)
	separator := strings.Repeat("-", len(headerText)+6)

	// 2. System & Hardware Metrics
	hostModel := getHostModel()
	kernel := getKernelVersion()
	uptime := getUptime()
	packages := getPackageCount()
	shellName := getShellWithVersion()
	deWM := getDesktopWM()
	cpuModel := getCPUModel()
	gpuModel := getGPUModel()
	memUsage := getMemoryUsage()
	aiModel := cfg.AI.DefaultModel
	themeName := strings.Title(cfg.Theme.Palette)

	// 3. Format Key-Value Lines
	lblStyle := lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(TextLightColor)
	accentStyle := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)

	lines := []string{
		headerStyle.Render(headerText),
		sepStyle.Render(separator),
		fmt.Sprintf("%s %s", lblStyle.Render("OS:       "), valStyle.Render("Ubuntu 24.04 LTS (Noble Numbat)")),
		fmt.Sprintf("%s %s", lblStyle.Render("Host:     "), valStyle.Render(hostModel)),
		fmt.Sprintf("%s %s", lblStyle.Render("Kernel:   "), valStyle.Render(kernel)),
		fmt.Sprintf("%s %s", lblStyle.Render("Uptime:   "), valStyle.Render(uptime)),
		fmt.Sprintf("%s %s", lblStyle.Render("Packages: "), valStyle.Render(packages)),
		fmt.Sprintf("%s %s", lblStyle.Render("Shell:    "), valStyle.Render(shellName)),
		fmt.Sprintf("%s %s", lblStyle.Render("DE/WM:    "), valStyle.Render(deWM)),
		fmt.Sprintf("%s %s", lblStyle.Render("CPU:      "), valStyle.Render(cpuModel)),
		fmt.Sprintf("%s %s", lblStyle.Render("GPU:      "), valStyle.Render(gpuModel)),
		fmt.Sprintf("%s %s", lblStyle.Render("Memory:   "), valStyle.Render(memUsage)),
		sepStyle.Render(separator),
		fmt.Sprintf("%s %s", lblStyle.Render("Vula:     "), accentStyle.Render(fmt.Sprintf("Go %s | %s", runtime.Version(), Version))),
		fmt.Sprintf("%s %s", lblStyle.Render("Theme:    "), valStyle.Render(themeName)),
		fmt.Sprintf("%s %s", lblStyle.Render("Local AI: "), valStyle.Render(fmt.Sprintf("Ollama (%s)", aiModel))),
		fmt.Sprintf("%s %s", lblStyle.Render("Voice AI: "), valStyle.Render("Whisper STT + Piper TTS")),
	}

	// 4. Render Dual-Row Color Palette (Standard + Bright)
	colorPalette := renderDualColorPalette()
	infoBlock := strings.Join(lines, "\n") + "\n\n" + colorPalette

	// 5. Render Logo on Left
	logo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		MarginRight(4).
		Render(strings.TrimSpace(FetchLogoASCII))

	content := lipgloss.JoinHorizontal(lipgloss.Top, logo, infoBlock)

	// Clean, unclipped container with subtle outer margin
	return lipgloss.NewStyle().Margin(1, 1).Render(content) + "\n"
}

func getHostModel() string {
	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		model := strings.TrimSpace(string(data))
		if vendor, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
			model = strings.TrimSpace(string(vendor)) + " " + model
		}
		if model != "" {
			return model
		}
	}
	return "PC / Laptop"
}

func getKernelVersion() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "Linux"
	}
	return strings.TrimSpace(string(out))
}

func getUptime() string {
	out, err := exec.Command("uptime", "-p").Output()
	if err == nil {
		up := strings.TrimSpace(string(out))
		up = strings.TrimPrefix(up, "up ")
		return up
	}
	return "unknown"
}

func getPackageCount() string {
	dpkgOut, err := exec.Command("sh", "-c", "dpkg-query -f '.' -W 2>/dev/null | wc -c").Output()
	dpkgCount := "0"
	if err == nil {
		dpkgCount = strings.TrimSpace(string(dpkgOut))
	}
	return fmt.Sprintf("%s (dpkg)", dpkgCount)
}

func getShellWithVersion() string {
	shell := os.Getenv("SHELL")
	name := "bash"
	if shell != "" {
		name = filepath.Base(shell)
	}
	out, err := exec.Command(name, "--version").Output()
	if err == nil {
		firstLine := strings.Split(string(out), "\n")[0]
		fields := strings.Fields(firstLine)
		if len(fields) >= 4 {
			ver := fields[3]
			if idx := strings.Index(ver, "("); idx != -1 {
				ver = ver[:idx]
			}
			return fmt.Sprintf("%s %s", name, ver)
		}
	}
	return name
}

func getDesktopWM() string {
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType == "" {
		sessionType = "x11"
	}
	return fmt.Sprintf("GNOME 46 (Mutter / %s)", strings.ToUpper(sessionType))
}

func getCPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "x86_64 CPU"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				model := strings.TrimSpace(parts[1])
				fields := strings.Fields(model)
				return strings.Join(fields, " ")
			}
		}
	}
	return "x86_64 CPU"
}

func getGPUModel() string {
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		out, err := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return "Integrated / OpenGL / Vulkan"
}

func getMemoryUsage() string {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if fields[0] == "MemTotal:" {
				memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
			} else if fields[0] == "MemAvailable:" {
				memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		}
	}

	if memTotal == 0 {
		return "unknown"
	}

	memUsed := memTotal - memAvailable
	memUsedMB := memUsed / 1024
	memTotalMB := memTotal / 1024

	return fmt.Sprintf("%d MiB / %d MiB", memUsedMB, memTotalMB)
}

func renderDualColorPalette() string {
	row1 := []lipgloss.Color{
		lipgloss.Color("#181825"), // Black
		lipgloss.Color("#F38BA8"), // Red
		lipgloss.Color("#A6E3A1"), // Green
		lipgloss.Color("#F9E2AF"), // Yellow
		lipgloss.Color("#89B4FA"), // Blue
		lipgloss.Color("#CBA6F7"), // Magenta
		lipgloss.Color("#94E2D5"), // Cyan
		lipgloss.Color("#BAC2DE"), // White
	}

	row2 := []lipgloss.Color{
		lipgloss.Color("#585B70"), // Bright Black
		lipgloss.Color("#EB6F92"), // Bright Red
		lipgloss.Color("#31748F"), // Bright Green
		lipgloss.Color("#F6C177"), // Bright Yellow
		lipgloss.Color("#9CCFD8"), // Bright Blue
		lipgloss.Color("#C4A7E7"), // Bright Magenta
		lipgloss.Color("#EA9A97"), // Bright Cyan
		lipgloss.Color("#E0DEF4"), // Bright White
	}

	var sb strings.Builder
	for _, c := range row1 {
		sb.WriteString(lipgloss.NewStyle().Background(c).Render("   "))
	}
	sb.WriteString("\n")
	for _, c := range row2 {
		sb.WriteString(lipgloss.NewStyle().Background(c).Render("   "))
	}

	return sb.String()
}
