package ai

import (
	"testing"
)

func TestEvaluateComplexity(t *testing.T) {
	tests := []struct {
		name       string
		prompt     string
		taskType   string
		threshold  int
		wantResult TaskComplexity
	}{
		{
			name:       "Short git commit prompt",
			prompt:     "Fix typo in README.md",
			taskType:   "git_commit",
			threshold:  1200,
			wantResult: ComplexityLight,
		},
		{
			name:       "Short command prompt",
			prompt:     "find all pdf files",
			taskType:   "suggest_command",
			threshold:  1200,
			wantResult: ComplexityLight,
		},
		{
			name:       "Long explanation prompt exceeding threshold",
			prompt:     makeString(5000),
			taskType:   "explain_selection",
			threshold:  1200,
			wantResult: ComplexityHeavy,
		},
		{
			name:       "Prompt with heavy keyword refactor",
			prompt:     "Please refactor this complex module",
			taskType:   "general",
			threshold:  1200,
			wantResult: ComplexityHeavy,
		},
		{
			name:       "Normal short general question",
			prompt:     "What is the capital of France?",
			taskType:   "general",
			threshold:  1200,
			wantResult: ComplexityLight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateComplexity(tt.prompt, tt.taskType, tt.threshold)
			if got != tt.wantResult {
				t.Errorf("EvaluateComplexity() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func makeString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
