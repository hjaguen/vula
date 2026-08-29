package doctor

import (
	"testing"

	"github.com/vula-os/vula/internal/config"
)

func TestRunDiagnostics(t *testing.T) {
	cfg := config.DefaultConfig()
	report := RunDiagnostics(cfg)

	if report == nil {
		t.Fatal("RunDiagnostics returned nil")
	}

	if len(report.Results) == 0 {
		t.Fatal("expected at least one diagnostic result")
	}

	rendered := report.Render()
	if len(rendered) == 0 {
		t.Error("expected non-empty rendered diagnostic output")
	}
}
