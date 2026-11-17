package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupCreateListAndRestore(t *testing.T) {
	// Create isolated temporary home directory
	tmpHome, err := os.MkdirTemp("", "vula-backup-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	// Create mock dotfiles
	mockBashrc := filepath.Join(tmpHome, ".bashrc")
	_ = os.WriteFile(mockBashrc, []byte("# original bashrc"), 0644)

	mgr := NewManager()

	// 1. Test Create Backup
	info, err := mgr.Create("test-snap", "Test snapshot description")
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	if info.Name != "test-snap" {
		t.Errorf("expected name 'test-snap', got '%s'", info.Name)
	}
	if len(info.Dotfiles) == 0 {
		t.Errorf("expected .bashrc to be included in dotfiles backup")
	}

	// 2. Test List Backups
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("failed to list backups: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 backup in list, got %d", len(list))
	}
	if list[0].ID != info.ID {
		t.Errorf("expected ID '%s', got '%s'", info.ID, list[0].ID)
	}

	// 3. Modify original .bashrc to simulate change
	_ = os.WriteFile(mockBashrc, []byte("# modified by vula"), 0644)

	// 4. Test Restore Backup
	restored, err := mgr.Restore("test-snap")
	if err != nil {
		t.Fatalf("failed to restore backup: %v", err)
	}
	if restored.ID != info.ID {
		t.Errorf("expected restored ID '%s', got '%s'", info.ID, restored.ID)
	}

	// Verify content was restored
	restoredContent, err := os.ReadFile(mockBashrc)
	if err != nil {
		t.Fatalf("failed to read restored bashrc: %v", err)
	}
	if string(restoredContent) != "# original bashrc" {
		t.Errorf("expected '# original bashrc', got '%s'", string(restoredContent))
	}

	// 5. Test Delete Backup
	deleted, err := mgr.Delete("test-snap")
	if err != nil {
		t.Fatalf("failed to delete backup: %v", err)
	}
	if deleted.ID != info.ID {
		t.Errorf("expected deleted ID '%s', got '%s'", info.ID, deleted.ID)
	}

	listAfterDelete, _ := mgr.List()
	if len(listAfterDelete) != 0 {
		t.Errorf("expected 0 backups after delete, got %d", len(listAfterDelete))
	}
}
