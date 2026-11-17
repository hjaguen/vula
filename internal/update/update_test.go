package update

import "testing"

func TestIsNewer(t *testing.T) {
	tests := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"v1.0.0", "v0.6.0", true},
		{"v0.6.1", "v0.6.0", true},
		{"v0.7.0", "v0.6.5", true},
		{"v0.6.0", "v0.6.0", false},
		{"v0.5.9", "v0.6.0", false},
		{"1.0.0", "0.9.0", true},
		{"v1.0.0-rc1", "v1.0.0", false},
	}

	for _, tt := range tests {
		got := IsNewer(tt.latest, tt.current)
		if got != tt.expected {
			t.Errorf("IsNewer(%q, %q) = %v; want %v", tt.latest, tt.current, got, tt.expected)
		}
	}
}
