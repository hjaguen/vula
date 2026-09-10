package actions

import "time"

type ActionType string

const (
	ActionTypeLaunchApp     ActionType = "launch_app"
	ActionTypeManageFiles   ActionType = "manage_files"
	ActionTypeControlSystem ActionType = "control_system"
	ActionTypeProcessControl ActionType = "process_control"
	ActionTypeSearch        ActionType = "search"
)

type RiskLevel string

const (
	RiskLevelSafe      RiskLevel = "SAFE"
	RiskLevelSensitive RiskLevel = "SENSITIVE"
	RiskLevelHigh      RiskLevel = "HIGH"
)

type ActionPayload struct {
	Type        ActionType             `json:"type"`
	Description string                 `json:"description"`
	App         string                 `json:"app,omitempty"`
	Operation   string                 `json:"operation,omitempty"` // move, copy, mkdir, find, terminate
	Source      string                 `json:"source,omitempty"`
	Destination string                 `json:"destination,omitempty"`
	Setting     string                 `json:"setting,omitempty"`   // volume, brightness, night_light, theme, lock, screenshot
	Value       string                 `json:"value,omitempty"`
	Query       string                 `json:"query,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

type ActionPlan struct {
	OriginalPrompt string          `json:"original_prompt"`
	Summary        string          `json:"summary"`
	Actions        []ActionPayload `json:"actions"`
}

type AuditEntry struct {
	Timestamp      time.Time     `json:"timestamp"`
	UserPrompt     string        `json:"user_prompt"`
	Action         ActionPayload `json:"action"`
	Risk           RiskLevel     `json:"risk"`
	UserConfirmed  bool          `json:"user_confirmed"`
	ExecutedStatus string        `json:"executed_status"`
	Error          string        `json:"error,omitempty"`
}
