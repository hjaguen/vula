package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/config"
)

// AssessRisk evaluates the risk level of an OS action payload
func AssessRisk(action ActionPayload) RiskLevel {
	switch action.Type {
	case ActionTypeLaunchApp, ActionTypeSearch:
		return RiskLevelSafe

	case ActionTypeControlSystem:
		switch strings.ToLower(action.Setting) {
		case "volume", "brightness", "theme", "night_light", "screenshot", "lock":
			return RiskLevelSafe
		case "shutdown", "reboot", "poweroff":
			return RiskLevelHigh
		default:
			return RiskLevelSensitive
		}

	case ActionTypeManageFiles:
		switch strings.ToLower(action.Operation) {
		case "move", "copy", "mkdir":
			// Elevate risk if targeting critical system paths
			if isSystemCriticalPath(action.Source) || isSystemCriticalPath(action.Destination) {
				return RiskLevelHigh
			}
			return RiskLevelSensitive
		case "delete", "rm", "remove", "purge":
			return RiskLevelHigh
		default:
			return RiskLevelSensitive
		}

	case ActionTypeProcessControl:
		switch strings.ToLower(action.Operation) {
		case "find", "list":
			return RiskLevelSafe
		case "terminate", "kill", "stop":
			return RiskLevelSensitive
		default:
			return RiskLevelSensitive
		}
	}

	return RiskLevelSensitive
}

// SanitizePath expands home directory, resolves absolute path, and validates system boundaries
func SanitizePath(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", fmt.Errorf("empty path")
	}

	// Expand ~ to user home
	if strings.HasPrefix(trimmed, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		trimmed = filepath.Join(home, strings.TrimPrefix(trimmed, "~"))
	}

	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s: %w", input, err)
	}

	// Block dangerous root system paths
	cleaned := filepath.Clean(absPath)
	if isSystemCriticalPath(cleaned) {
		return "", fmt.Errorf("access denied: %s is a restricted system path", cleaned)
	}

	return cleaned, nil
}

func isSystemCriticalPath(path string) bool {
	if path == "" {
		return false
	}
	clean := filepath.Clean(path)
	restricted := []string{
		"/",
		"/boot",
		"/etc",
		"/sys",
		"/proc",
		"/dev",
		"/usr",
		"/lib",
		"/lib64",
		"/sbin",
		"/bin",
	}
	for _, r := range restricted {
		if clean == r || (r != "/" && strings.HasPrefix(clean, r+"/")) {
			return true
		}
	}
	return false
}

// LogAudit appends an audit entry to ~/.config/vula/audit.log
func LogAudit(entry AuditEntry) {
	cfgDir, err := config.ConfigDir()
	if err != nil {
		return
	}

	auditFile := filepath.Join(cfgDir, "audit.log")
	f, err := os.OpenFile(auditFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	data, err := json.Marshal(entry)
	if err == nil {
		_, _ = f.Write(append(data, '\n'))
	}
}

// GetAuditLogs reads all audit entries from ~/.config/vula/audit.log
func GetAuditLogs() ([]AuditEntry, error) {
	cfgDir, err := config.ConfigDir()
	if err != nil {
		return nil, err
	}

	auditFile := filepath.Join(cfgDir, "audit.log")
	data, err := os.ReadFile(auditFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var entries []AuditEntry
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		var entry AuditEntry
		if err := json.Unmarshal([]byte(trimmed), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}
