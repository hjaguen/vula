package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vula-os/vula/internal/ai"
)

const systemPlannerPrompt = `You are the Vula OS Action Planner integrated into Ubuntu 24.04 LTS.
Your goal is to parse the user's natural language request into a valid JSON ActionPlan.

Valid Action Types and fields:
1. "launch_app": {"type": "launch_app", "description": "...", "app": "code|dbeaver|obsidian|spotify|browser|terminal|gimp|inkscape|blender|vlc|..."}
2. "manage_files": {"type": "manage_files", "description": "...", "operation": "move|copy|mkdir", "source": "...", "destination": "..."}
3. "control_system": {"type": "control_system", "description": "...", "setting": "volume|brightness|theme|night_light|screenshot|lock", "value": "..."}
4. "process_control": {"type": "process_control", "description": "...", "operation": "find|terminate", "app": "name_or_pid"}
5. "search": {"type": "search", "description": "...", "query": "..."}

OUTPUT ONLY VALID UNADORNED JSON MATCHING THIS EXACT SCHEME:
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
}
DO NOT INCLUDE MARKDOWN CODE BLOCKS OR EXPLANATORY TEXT OUTSIDE THE JSON.`

// ParseIntent translates a natural language instruction into a validated ActionPlan
func ParseIntent(ctx context.Context, client *ai.Client, userPrompt string) (*ActionPlan, error) {
	trimmedPrompt := strings.TrimSpace(userPrompt)
	if trimmedPrompt == "" {
		return nil, fmt.Errorf("empty instruction")
	}

	prompt := fmt.Sprintf("User Request: \"%s\"", trimmedPrompt)
	resp, err := client.AskTask(ctx, prompt, "suggest_command", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse intent with AI: %w", err)
	}

	cleaned := strings.TrimSpace(resp)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var plan ActionPlan
	if err := json.Unmarshal([]byte(cleaned), &plan); err != nil {
		return nil, fmt.Errorf("AI returned invalid action JSON payload: %w (raw response: %s)", err, cleaned)
	}

	plan.OriginalPrompt = trimmedPrompt
	if len(plan.Actions) == 0 {
		return nil, fmt.Errorf("no executable actions found for prompt: \"%s\"", userPrompt)
	}

	return &plan, nil
}
