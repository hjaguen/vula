package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vula-os/vula/internal/ai"
)

const systemPlannerPrompt = `You are the Vula OS Action Planner integrated into Ubuntu 24.04 LTS.
You are connected directly to Vula's Linux Desktop Action Engine. Vula HAS FULL NATIVE CAPABILITIES to execute desktop actions, including adjusting screen brightness, controlling volume, changing desktop themes, launching apps, moving files, searching, and managing processes.

Your ONLY job is to parse the user's instruction into a JSON ActionPlan.

Valid Action Types and fields:
1. "launch_app": {"type": "launch_app", "description": "...", "app": "code|dbeaver|obsidian|spotify|browser|terminal|gimp|inkscape|blender|vlc|..."}
2. "manage_files": {"type": "manage_files", "description": "...", "operation": "move|copy|mkdir", "source": "...", "destination": "..."}
3. "control_system": {"type": "control_system", "description": "...", "setting": "volume|brightness|theme|night_light|screenshot|lock", "value": "..."}
4. "process_control": {"type": "process_control", "description": "...", "operation": "find|terminate", "app": "name_or_pid"}
5. "search": {"type": "search", "description": "...", "query": "..."}

CRITICAL RULES:
- You MUST respond ONLY with a raw JSON object matching the scheme below.
- Do NOT output any conversational chat explanations, apologies, or disclaimers.
- Do NOT say "Lo siento" or claim you cannot perform the action. Vula WILL execute the action.
- Output ONLY UNADORNED JSON matching this exact scheme:

{
  "original_prompt": "<original user text>",
  "summary": "<1 sentence Spanish summary of what will be done>",
  "actions": [
    {
      "type": "...",
      "description": "...",
      "app": "...",
      "operation": "...",
      "source": "...",
      "destination": "...",
      "setting": "...",
      "value": "...",
      "query": "..."
    }
  ]
}`

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

	var plan ActionPlan
	if err := json.Unmarshal([]byte(cleaned), &plan); err != nil {
		return nil, fmt.Errorf("AI returned invalid action JSON payload: %w (raw response: %s)", err, resp)
	}

	plan.OriginalPrompt = trimmedPrompt
	if len(plan.Actions) == 0 {
		return nil, fmt.Errorf("no executable actions found for prompt: \"%s\"", userPrompt)
	}

	return &plan, nil
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
