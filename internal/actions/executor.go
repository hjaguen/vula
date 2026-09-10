package actions

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/vula-os/vula/internal/config"
	"github.com/vula-os/vula/internal/theme"
	"github.com/vula-os/vula/internal/ui"
)

// ExecutePlan runs an ActionPlan with safety checks, path sanitization, confirmation UI, and audit logging
func ExecutePlan(ctx context.Context, plan *ActionPlan, cfg *config.Config) error {
	if plan == nil || len(plan.Actions) == 0 {
		return fmt.Errorf("no actions to execute")
	}

	fmt.Println(ui.RenderHeader("Vula OS Control", plan.Summary))

	for i, action := range plan.Actions {
		risk := AssessRisk(action)

		// Ask confirmation for Sensitive or High risk actions
		if risk == RiskLevelSensitive || risk == RiskLevelHigh {
			confirmed, err := promptConfirmation(i+1, len(plan.Actions), action, risk)
			if err != nil || !confirmed {
				LogAudit(AuditEntry{
					Timestamp:      time.Now(),
					UserPrompt:     plan.OriginalPrompt,
					Action:         action,
					Risk:           risk,
					UserConfirmed:  false,
					ExecutedStatus: "CANCELLED",
				})
				fmt.Println(ui.WarnStyle.Render("Action cancelled by user."))
				return nil
			}
		}

		// Execute action
		fmt.Printf("  • Executing: %s (%s)...\n", ui.InfoStyle.Render(action.Description), risk)
		err := executeSingleAction(ctx, action, cfg)
		statusStr := "SUCCESS"
		errMsg := ""
		if err != nil {
			statusStr = "FAILED"
			errMsg = err.Error()
			fmt.Printf("    %s %v\n", ui.ErrorStyle.Render("Failed:"), err)
		} else {
			fmt.Printf("    %s Complete\n", ui.SuccessStyle.Render("✓"))
			notifyOSAction(action.Description)
		}

		LogAudit(AuditEntry{
			Timestamp:      time.Now(),
			UserPrompt:     plan.OriginalPrompt,
			Action:         action,
			Risk:           risk,
			UserConfirmed:  true,
			ExecutedStatus: statusStr,
			Error:          errMsg,
		})

		if err != nil {
			return err
		}
	}

	fmt.Println("\n" + ui.SuccessStyle.Render("✓ OS Control task batch executed successfully!"))
	return nil
}

func promptConfirmation(index, total int, action ActionPayload, risk RiskLevel) (bool, error) {
	fmt.Println()
	title := fmt.Sprintf("Action [%d/%d] Confirmation Required (%s Risk)", index, total, risk)
	details := fmt.Sprintf("Description: %s\nType:        %s", action.Description, action.Type)
	if action.Source != "" {
		details += fmt.Sprintf("\nSource:      %s", action.Source)
	}
	if action.Destination != "" {
		details += fmt.Sprintf("\nDestination: %s", action.Destination)
	}
	if action.App != "" {
		details += fmt.Sprintf("\nApp/Target:  %s", action.App)
	}

	fmt.Println(ui.RenderHeader(title, details))

	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want Vula to execute this operation?").
				Affirmative("Yes, Execute").
				Negative("No, Cancel").
				Value(&confirm),
		),
	)

	err := form.Run()
	return confirm, err
}

func executeSingleAction(ctx context.Context, action ActionPayload, cfg *config.Config) error {
	switch action.Type {
	case ActionTypeLaunchApp:
		return launchApplication(action.App)

	case ActionTypeControlSystem:
		return controlSystemSetting(ctx, action.Setting, action.Value, cfg)

	case ActionTypeManageFiles:
		return manageFiles(action.Operation, action.Source, action.Destination)

	case ActionTypeProcessControl:
		return controlProcess(action.Operation, action.App)

	case ActionTypeSearch:
		return searchSystem(action.Query)
	}

	return fmt.Errorf("unknown action type: %s", action.Type)
}

func launchApplication(appName string) error {
	trimmed := strings.ToLower(strings.TrimSpace(appName))
	if trimmed == "" {
		return fmt.Errorf("app name missing")
	}

	// 1. Common developer app aliases
	var binCandidate string
	switch trimmed {
	case "code", "vscode", "vs code":
		binCandidate = "code"
	case "dbeaver":
		binCandidate = "dbeaver-ce"
	case "obsidian":
		binCandidate = "obsidian"
	case "browser", "chrome", "brave":
		binCandidate = "brave-browser"
	case "terminal":
		binCandidate = "ghostty"
	default:
		binCandidate = trimmed
	}

	// Try running direct binary
	if path, err := exec.LookPath(binCandidate); err == nil {
		cmd := exec.Command(path)
		return cmd.Start()
	}

	// Try gtk-launch or gio launch
	if _, err := exec.LookPath("gtk-launch"); err == nil {
		cmd := exec.Command("gtk-launch", binCandidate)
		if err := cmd.Start(); err == nil {
			return nil
		}
	}

	return exec.Command("sh", "-c", appName+" &").Start()
}

func controlSystemSetting(ctx context.Context, setting, value string, cfg *config.Config) error {
	switch strings.ToLower(setting) {
	case "theme":
		tm := theme.NewManager(cfg)
		return tm.ApplyTheme(value)

	case "volume":
		val := strings.TrimSuffix(value, "%")
		if _, err := exec.LookPath("wpctl"); err == nil {
			return exec.Command("wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", val+"%").Run()
		}
		return exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", val+"%").Run()

	case "lock":
		if _, err := exec.LookPath("loginctl"); err == nil {
			return exec.Command("loginctl", "lock-session").Run()
		}
		return exec.Command("gnome-screensaver-command", "-l").Run()

	case "screenshot":
		if _, err := exec.LookPath("gnome-screenshot"); err == nil {
			return exec.Command("gnome-screenshot").Run()
		}
		return exec.Command("spectacle").Run()

	case "brightness":
		val := strings.TrimSuffix(value, "%")
		if _, err := exec.LookPath("brightnessctl"); err == nil {
			return exec.Command("brightnessctl", "set", val+"%").Run()
		}
		if _, err := exec.LookPath("gdbus"); err == nil {
			gdbusCmd := fmt.Sprintf("gdbus call --session --dest org.gnome.SettingsDaemon.Power --object-path /org/gnome/SettingsDaemon/Power --method org.freedesktop.DBus.Properties.Set org.gnome.SettingsDaemon.Power.Screen Brightness \"<int32 %s>\"", val)
			return exec.Command("sh", "-c", gdbusCmd).Run()
		}
		return fmt.Errorf("brightness control tool (brightnessctl or gdbus) not found")

	case "night_light":
		val := "true"
		if strings.ToLower(value) == "off" || value == "false" {
			val = "false"
		}
		return exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-enabled", val).Run()
	}

	return fmt.Errorf("unsupported system setting: %s", setting)
}

func manageFiles(op, src, dest string) error {
	cleanSrc, err := SanitizePath(src)
	if err != nil && op != "mkdir" {
		return err
	}

	cleanDest := ""
	if dest != "" {
		cleanDest, err = SanitizePath(dest)
		if err != nil {
			return err
		}
	}

	switch strings.ToLower(op) {
	case "mkdir":
		targetDir := cleanSrc
		if targetDir == "" {
			targetDir = cleanDest
		}
		return os.MkdirAll(targetDir, 0755)

	case "move":
		if cleanDest == "" {
			return fmt.Errorf("destination directory missing for move operation")
		}
		_ = os.MkdirAll(filepath.Dir(cleanDest), 0755)

		// Try gio move first for GNOME desktop trash safety
		if _, err := exec.LookPath("gio"); err == nil {
			if err := exec.Command("gio", "move", cleanSrc, cleanDest).Run(); err == nil {
				return nil
			}
		}
		return os.Rename(cleanSrc, cleanDest)

	case "copy":
		if cleanDest == "" {
			return fmt.Errorf("destination directory missing for copy operation")
		}
		_ = os.MkdirAll(filepath.Dir(cleanDest), 0755)
		return copyFileContents(cleanSrc, cleanDest)
	}

	return fmt.Errorf("unsupported file operation: %s", op)
}

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func controlProcess(op, target string) error {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" {
		return fmt.Errorf("target process name or PID missing")
	}

	switch strings.ToLower(op) {
	case "find", "list":
		out, err := exec.Command("pgrep", "-a", trimmed).Output()
		if err != nil {
			return fmt.Errorf("no process found matching '%s'", trimmed)
		}
		fmt.Printf("    Process List:\n%s\n", string(out))
		return nil

	case "terminate", "kill", "stop":
		return exec.Command("pkill", "-f", trimmed).Run()
	}

	return fmt.Errorf("unsupported process operation: %s", op)
}

func searchSystem(query string) error {
	if _, err := exec.LookPath("fd"); err == nil {
		out, err := exec.Command("fd", "-max-results", "10", query, os.Getenv("HOME")).Output()
		if err == nil {
			fmt.Printf("    Search Results for '%s':\n%s\n", query, string(out))
			return nil
		}
	}

	out, err := exec.Command("find", os.Getenv("HOME"), "-name", "*"+query+"*", "-maxdepth", "4").Output()
	if err != nil {
		return err
	}
	fmt.Printf("    Search Results for '%s':\n%s\n", query, string(out))
	return nil
}

func notifyOSAction(description string) {
	_ = exec.Command("notify-send", "-a", "Vula OS Control", "-i", "preferences-system", "⚡ Vula OS Control", description).Start()
}
