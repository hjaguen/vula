package keys

import (
	"testing"
)

func TestFormatKeycap(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<Super>space", "[Super] + [SPACE]"},
		{"<Super><Alt>v", "[Super] + [Alt] + [V]"},
		{"<Super><Shift>k", "[Super] + [Shift] + [K]"},
		{"", "[Unset]"},
	}

	for _, tt := range tests {
		result := FormatKeycap(tt.input)
		if result != tt.expected {
			t.Errorf("FormatKeycap(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetDefaultKeybindings(t *testing.T) {
	bindings := GetKeybindings()
	if len(bindings) < 5 {
		t.Errorf("Expected at least 5 default keybindings, got %d", len(bindings))
	}
}
