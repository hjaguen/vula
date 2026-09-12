package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vula-os/vula/internal/ai"
)

const systemPlannerPrompt = `You are the Vula OS Action Planner integrated into Ubuntu 24.04 LTS.
Vula HAS FULL NATIVE CAPABILITIES to execute desktop actions.

Your ONLY task is to return a raw JSON ActionPlan object for the user's request.

Example Output 1 (Adjust Brightness):
{
  "summary": "Ajustar brillo al 40%",
  "actions": [
    {
      "type": "control_system",
      "description": "Ajustar brillo al 40%",
      "setting": "brightness",
      "value": "40%"
    }
  ]
}

Example Output 2 (Launch App):
{
  "summary": "Abrir Visual Studio Code",
  "actions": [
    {
      "type": "launch_app",
      "description": "Abrir VS Code",
      "app": "code"
    }
  ]
}

Example Output 3 (Move File):
{
  "summary": "Mover informe a Documentos",
  "actions": [
    {
      "type": "manage_files",
      "description": "Mover archivo",
      "operation": "move",
      "source": "~/Downloads/informe.pdf",
      "destination": "~/Documents/informe.pdf"
    }
  ]
}

Allowed Enum Values:
- "type": "launch_app", "manage_files", "control_system", "process_control", "search"
- "setting": "brightness", "volume", "theme", "night_light", "screenshot", "lock", "tiling", "gaps"
- "operation": "move", "copy", "mkdir", "terminate", "find"

CRITICAL: Return ONLY valid unadorned JSON. No markdown backticks, no explanatory chat, no disclaimers.`

// ParseIntent translates a natural language instruction into a validated ActionPlan
func ParseIntent(ctx context.Context, client *ai.Client, userPrompt string) (*ActionPlan, error) {
	trimmedPrompt := strings.TrimSpace(userPrompt)
	if trimmedPrompt == "" {
		return nil, fmt.Errorf("empty instruction")
	}

	prompt := fmt.Sprintf("User Request: \"%s\"\n\nGenerate the JSON ActionPlan now:", trimmedPrompt)
	resp, err := client.AskTaskWithSystemPrompt(ctx, prompt, systemPlannerPrompt, "suggest_command", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse intent with AI: %w", err)
	}

	cleaned := extractJSON(resp)
	repaired := repairJSON(cleaned)

	var plan ActionPlan
	if err := json.Unmarshal([]byte(repaired), &plan); err != nil {
		return nil, fmt.Errorf("AI returned invalid action JSON payload: %w (raw response: %s)", err, resp)
	}

	sanitizeActionPlan(&plan, trimmedPrompt)

	if len(plan.Actions) == 0 {
		return nil, fmt.Errorf("no executable actions found for prompt: \"%s\"", userPrompt)
	}

	return &plan, nil
}

func sanitizeActionPlan(plan *ActionPlan, originalPrompt string) {
	lowerPrompt := strings.ToLower(originalPrompt)
	if plan.OriginalPrompt == "" {
		plan.OriginalPrompt = originalPrompt
	}
	if plan.Summary == "" || strings.Contains(plan.Summary, "...") {
		plan.Summary = fmt.Sprintf("Ejecutar acción: %s", originalPrompt)
	}

	for i := range plan.Actions {
		act := &plan.Actions[i]

		// Sanitize setting field if model echoed enum strings like "volume|brightness|..."
		if strings.Contains(act.Setting, "|") || act.Setting == "" {
			if strings.Contains(lowerPrompt, "brillo") || strings.Contains(lowerPrompt, "pantalla") {
				act.Setting = "brightness"
			} else if strings.Contains(lowerPrompt, "volumen") || strings.Contains(lowerPrompt, "audio") || strings.Contains(lowerPrompt, "sonido") {
				act.Setting = "volume"
			} else if strings.Contains(lowerPrompt, "tema") || strings.Contains(lowerPrompt, "modo") {
				act.Setting = "theme"
			} else if strings.Contains(lowerPrompt, "luz") || strings.Contains(lowerPrompt, "noche") {
				act.Setting = "night_light"
			} else if strings.Contains(lowerPrompt, "captura") || strings.Contains(lowerPrompt, "screenshot") {
				act.Setting = "screenshot"
			} else if strings.Contains(lowerPrompt, "bloquea") || strings.Contains(lowerPrompt, "lock") {
				act.Setting = "lock"
			} else if strings.Contains(lowerPrompt, "tiling") || strings.Contains(lowerPrompt, "ventana") || strings.Contains(lowerPrompt, "pantalla") {
				act.Setting = "tiling"
			} else if strings.Contains(lowerPrompt, "espacio") || strings.Contains(lowerPrompt, "gap") {
				act.Setting = "gaps"
			}
		}

		// Sanitize description if model left "..."
		if act.Description == "" || strings.Contains(act.Description, "...") {
			switch act.Type {
			case ActionTypeLaunchApp:
				act.Description = fmt.Sprintf("Abrir aplicación %s", act.App)
			case ActionTypeControlSystem:
				act.Description = fmt.Sprintf("Ajustar %s a %s", act.Setting, act.Value)
			case ActionTypeManageFiles:
				act.Description = fmt.Sprintf("%s %s -> %s", act.Operation, act.Source, act.Destination)
			case ActionTypeProcessControl:
				act.Description = fmt.Sprintf("%s proceso %s", act.Operation, act.App)
			case ActionTypeSearch:
				act.Description = fmt.Sprintf("Buscar %s", act.Query)
			default:
				act.Description = originalPrompt
			}
		}
	}
}

func extractJSON(input string) string {
	cleaned := strings.TrimSpace(input)

	// Strip markdown block markers if present
	if strings.Contains(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		var codeLines []string
		inBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				inBlock = !inBlock
				continue
			}
			if inBlock {
				codeLines = append(codeLines, line)
			}
		}
		if len(codeLines) > 0 {
			cleaned = strings.Join(codeLines, "\n")
		}
	}

	// Find first '{' and last '}'
	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start != -1 && end != -1 && end > start {
		return cleaned[start : end+1]
	}

	return cleaned
}

func repairJSON(input string) string {
	var sb strings.Builder
	inString := false
	escaped := false
	for _, r := range input {
		if r == '"' && !escaped {
			inString = !inString
			sb.WriteRune(r)
			escaped = false
			continue
		}
		if escaped {
			escaped = false
			sb.WriteRune(r)
			continue
		}
		if r == '\\' {
			escaped = true
			sb.WriteRune(r)
			continue
		}
		if inString && r == '\n' {
			sb.WriteString("\\n")
			continue
		}
		if inString && r == '\r' {
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}
