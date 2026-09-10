package ai

import "strings"

type TaskComplexity string

const (
	ComplexityLight TaskComplexity = "light"
	ComplexityHeavy TaskComplexity = "heavy"
)

// EvaluateComplexity determines if a prompt/task should be handled locally or offloaded to Cloud AI
func EvaluateComplexity(prompt string, taskType string, threshold int) TaskComplexity {
	if threshold <= 0 {
		threshold = 1200
	}

	// 1. Explicit task type heuristics
	switch taskType {
	case "git_commit", "suggest_command":
		// Git commits and 1-line command suggestions are inherently lightweight
		if len(prompt) < threshold*3 {
			return ComplexityLight
		}
	case "explain_selection", "diagnose_error":
		// Deep code explanations or long stack traces with > 1200 chars/tokens are heavy
		if len(prompt) > threshold {
			return ComplexityHeavy
		}
	}

	// 2. Prompt length heuristic
	// Rough estimation: 1 token ~ 4 characters
	estimatedTokens := len(prompt) / 4
	if estimatedTokens >= threshold || len(prompt) >= threshold*3 {
		return ComplexityHeavy
	}

	// 3. Keyword / complexity pattern analysis
	heavyKeywords := []string{
		"refactor", "architect", "optimize performance", "debug memory leak",
		"multi-file", "complex algorithm", "convert codebase", "security audit",
	}

	lower := strings.ToLower(prompt)
	for _, kw := range heavyKeywords {
		if strings.Contains(lower, kw) {
			return ComplexityHeavy
		}
	}

	return ComplexityLight
}
