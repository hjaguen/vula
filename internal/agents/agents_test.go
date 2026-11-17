package agents

import (
	"os"
	"strings"
	"testing"
)

func TestAgentSkillGeneration(t *testing.T) {
	md := GenerateSkillMarkdown()
	if !strings.Contains(md, "name: vula") {
		t.Errorf("Skill markdown missing name frontmatter")
	}
	if !strings.Contains(md, "vula do") {
		t.Errorf("Skill markdown missing 'vula do'")
	}
	if !strings.Contains(md, "vula agents") {
		t.Errorf("Skill markdown missing 'vula agents'")
	}
}

func TestWriteSkillFile(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "vula-skill-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	path, err := WriteSkillFile()
	if err != nil {
		t.Fatalf("WriteSkillFile failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written skill file: %v", err)
	}
	if !strings.Contains(string(content), "Vula OS Agent Skill") {
		t.Errorf("Skill content unexpected: %s", string(content))
	}
}

func TestFindAgent(t *testing.T) {
	claude := FindAgent("claude")
	if claude == nil {
		t.Fatalf("expected to find 'claude'")
	}
	if claude.Command != "claude" {
		t.Errorf("expected command 'claude', got '%s'", claude.Command)
	}

	opencode := FindAgent("c")
	if opencode == nil {
		t.Fatalf("expected to find agent by alias 'c'")
	}
	if opencode.Name != "OpenCode" {
		t.Errorf("expected 'OpenCode', got '%s'", opencode.Name)
	}
}
