package actions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAssessRisk(t *testing.T) {
	tests := []struct {
		name     string
		action   ActionPayload
		expected RiskLevel
	}{
		{
			name: "Launch App",
			action: ActionPayload{
				Type: ActionTypeLaunchApp,
				App:  "code",
			},
			expected: RiskLevelSafe,
		},
		{
			name: "Volume Safe Control",
			action: ActionPayload{
				Type:    ActionTypeControlSystem,
				Setting: "volume",
				Value:   "80%",
			},
			expected: RiskLevelSafe,
		},
		{
			name: "Shutdown High Risk Control",
			action: ActionPayload{
				Type:    ActionTypeControlSystem,
				Setting: "shutdown",
			},
			expected: RiskLevelHigh,
		},
		{
			name: "File Move Sensitive",
			action: ActionPayload{
				Type:        ActionTypeManageFiles,
				Operation:   "move",
				Source:      "~/file.txt",
				Destination: "~/Documents/file.txt",
			},
			expected: RiskLevelSensitive,
		},
		{
			name: "File Delete High Risk",
			action: ActionPayload{
				Type:      ActionTypeManageFiles,
				Operation: "delete",
				Source:    "~/file.txt",
			},
			expected: RiskLevelHigh,
		},
		{
			name: "System Path File Move High Risk",
			action: ActionPayload{
				Type:        ActionTypeManageFiles,
				Operation:   "move",
				Source:      "/etc/passwd",
				Destination: "~/passwd",
			},
			expected: RiskLevelHigh,
		},
		{
			name: "Terminate Process Sensitive",
			action: ActionPayload{
				Type:      ActionTypeProcessControl,
				Operation: "terminate",
				App:       "firefox",
			},
			expected: RiskLevelSensitive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssessRisk(tt.action)
			if got != tt.expected {
				t.Errorf("AssessRisk(%s) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	t.Run("Expand Home Tilde", func(t *testing.T) {
		got, err := SanitizePath("~/testfile.txt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := filepath.Join(home, "testfile.txt")
		if got != expected {
			t.Errorf("SanitizePath(~/testfile.txt) = %s, want %s", got, expected)
		}
	})

	t.Run("Restricted System Path Root", func(t *testing.T) {
		_, err := SanitizePath("/etc")
		if err == nil {
			t.Error("expected error for /etc, got nil")
		}
	})

	t.Run("Empty Path", func(t *testing.T) {
		_, err := SanitizePath("   ")
		if err == nil {
			t.Error("expected error for empty path, got nil")
		}
	})
}

func TestSanitizeActionPlan(t *testing.T) {
	plan := &ActionPlan{
		Actions: []ActionPayload{
			{
				Type:        ActionTypeControlSystem,
				Description: "...",
				Setting:     "volume|brightness|theme|night_light|screenshot|lock",
				Value:       "40%",
			},
		},
	}

	sanitizeActionPlan(plan, "ajusta el brillo al 40%")

	if plan.Actions[0].Setting != "brightness" {
		t.Errorf("expected setting 'brightness', got '%s'", plan.Actions[0].Setting)
	}

	if plan.Actions[0].Description == "..." || plan.Actions[0].Description == "" {
		t.Errorf("expected valid description, got '%s'", plan.Actions[0].Description)
	}
}

func TestRepairJSON(t *testing.T) {
	input := `{"summary": "line1\nline2"}`
	repaired := repairJSON(input)
	if !testing.Verbose() && repaired == "" {
		t.Error("failed repairing JSON")
	}
}
