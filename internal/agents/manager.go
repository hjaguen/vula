package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/vula-os/vula/internal/config"
)

type AgentSession struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"status"` // running, idle, completed, error
	CreatedAt time.Time `json:"created_at"`
	WindowID  string    `json:"window_id,omitempty"`
}

type HerdrAgentListResponse struct {
	Result struct {
		Agents []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"agents"`
	} `json:"result"`
}

type FleetStatus struct {
	Installed    bool           `json:"installed"`
	BinaryPath   string         `json:"binary_path"`
	DaemonActive bool           `json:"daemon_active"`
	ActiveCount  int            `json:"active_count"`
	Sessions     []AgentSession `json:"sessions"`
}

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// IsInstalled returns true if herdr CLI binary is available on PATH
func (m *Manager) IsInstalled() bool {
	_, err := exec.LookPath("herdr")
	return err == nil
}

// GetBinaryPath returns full path to herdr binary
func (m *Manager) GetBinaryPath() string {
	path, err := exec.LookPath("herdr")
	if err != nil {
		return ""
	}
	return path
}

// EnsureDaemonRunning verifies or starts the herdr daemon service/socket
func (m *Manager) EnsureDaemonRunning() error {
	if !m.IsInstalled() {
		return fmt.Errorf("herdr is not installed. Run 'vula apps install herdr' or 'vula profile apply full-stack-dev'")
	}

	cmd := exec.Command("herdr", "status")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Daemon server might not be running; launch herdr server in background
	serverCmd := exec.Command("herdr", "server")
	serverCmd.Stdout = nil
	serverCmd.Stderr = nil
	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed starting Herdr server: %w", err)
	}
	time.Sleep(200 * time.Millisecond)
	return nil
}

// GetFleetStatus returns diagnostic status of the Herdr agent fleet
func (m *Manager) GetFleetStatus() FleetStatus {
	status := FleetStatus{
		Installed:  m.IsInstalled(),
		BinaryPath: m.GetBinaryPath(),
	}

	if !status.Installed {
		return status
	}

	sessions, err := m.ListSessions()
	if err == nil {
		status.DaemonActive = true
		status.Sessions = sessions
		status.ActiveCount = len(sessions)
	}

	return status
}

// ListSessions queries active Herdr agent panes/sessions
func (m *Manager) ListSessions() ([]AgentSession, error) {
	if !m.IsInstalled() {
		return nil, fmt.Errorf("herdr is not installed")
	}

	// Try herdr agent list first
	cmdAgent := exec.Command("herdr", "agent", "list")
	out, err := cmdAgent.Output()
	if err == nil && len(out) > 0 {
		var resp HerdrAgentListResponse
		if err := json.Unmarshal(out, &resp); err == nil && len(resp.Result.Agents) > 0 {
			var sessions []AgentSession
			for _, a := range resp.Result.Agents {
				sessions = append(sessions, AgentSession{
					ID:        a.ID,
					Name:      a.Name,
					Status:    a.Status,
					CreatedAt: time.Now(),
				})
			}
			return sessions, nil
		}
	}

	// Fallback to herdr session list
	cmdSession := exec.Command("herdr", "session", "list")
	textOut, err := cmdSession.Output()
	if err != nil {
		return []AgentSession{}, nil
	}

	var sessions []AgentSession
	lines := strings.Split(string(textOut), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "name") || strings.HasPrefix(line, "-") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			sessions = append(sessions, AgentSession{
				ID:        parts[0],
				Name:      parts[0],
				Status:    parts[1],
				CreatedAt: time.Now(),
			})
		} else if len(line) > 0 {
			sessions = append(sessions, AgentSession{
				ID:        fmt.Sprintf("session-%d", i),
				Name:      line,
				Status:    "running",
				CreatedAt: time.Now(),
			})
		}
	}

	return sessions, nil
}

// SpawnAgent launches a new agent session in Herdr with a prompt or task
func (m *Manager) SpawnAgent(name string, prompt string) (*AgentSession, error) {
	if err := m.EnsureDaemonRunning(); err != nil {
		return nil, err
	}

	if name == "" {
		name = fmt.Sprintf("agent-%d", time.Now().Unix())
	}

	// Execute herdr --session <name> or herdr agent start
	var cmd *exec.Cmd
	if prompt != "" {
		cmd = exec.Command("herdr", "--session", name, "--", prompt)
	} else {
		cmd = exec.Command("herdr", "--session", name)
	}
	cmd.Env = os.Environ()

	// Launch in background
	if err := cmd.Start(); err != nil {
		// Fallback: try herdr session create or agent start
		altCmd := exec.Command("herdr", "agent", "start", name)
		if altErr := altCmd.Run(); altErr != nil {
			return nil, fmt.Errorf("failed spawning agent session: %w", err)
		}
	}

	session := &AgentSession{
		ID:        name,
		Name:      name,
		Prompt:    prompt,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	return session, nil
}

// AttachSession attaches the current interactive terminal to a Herdr pane or session
func (m *Manager) AttachSession(sessionID string) error {
	if !m.IsInstalled() {
		return fmt.Errorf("herdr is not installed")
	}

	cmd := exec.Command("herdr", "session", "attach", sessionID)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		// Fallback to herdr agent attach
		altCmd := exec.Command("herdr", "agent", "attach", sessionID)
		altCmd.Stdin = os.Stdin
		altCmd.Stdout = os.Stdout
		altCmd.Stderr = os.Stderr
		return altCmd.Run()
	}
	return nil
}

// KillSession terminates an active Herdr agent session
func (m *Manager) KillSession(sessionID string) error {
	if !m.IsInstalled() {
		return fmt.Errorf("herdr is not installed")
	}

	cmd := exec.Command("herdr", "session", "stop", sessionID)
	if err := cmd.Run(); err != nil {
		altCmd := exec.Command("herdr", "session", "delete", sessionID)
		return altCmd.Run()
	}
	return nil
}
