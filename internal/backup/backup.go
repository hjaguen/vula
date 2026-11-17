package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BackupInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	HasDconf    bool      `json:"has_dconf"`
	Dotfiles    []string  `json:"dotfiles"`
	Path        string    `json:"path"`
}

type Manager struct {
	baseDir string
}

func NewManager() *Manager {
	home := os.Getenv("HOME")
	return &Manager{
		baseDir: filepath.Join(home, ".config", "vula", "backups"),
	}
}

// GetBaseDir returns the backups directory
func (m *Manager) GetBaseDir() string {
	return m.baseDir
}

// Create creates a new system state snapshot
func (m *Manager) Create(name, description string) (*BackupInfo, error) {
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backups directory: %w", err)
	}

	timestamp := time.Now()
	cleanName := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(name)), " ", "-")
	if cleanName == "" {
		cleanName = "snapshot"
	}
	id := fmt.Sprintf("%s_%s", timestamp.Format("20060102-150405"), cleanName)
	backupPath := filepath.Join(m.baseDir, id)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup folder: %w", err)
	}

	info := &BackupInfo{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   timestamp,
		Path:        backupPath,
		Dotfiles:    make([]string, 0),
	}

	// 1. Dump GNOME Dconf settings if dconf command is available
	if _, err := exec.LookPath("dconf"); err == nil {
		dconfPath := filepath.Join(backupPath, "gnome_dconf.dump")
		cmd := exec.Command("dconf", "dump", "/org/gnome/")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			if err := os.WriteFile(dconfPath, out, 0644); err == nil {
				info.HasDconf = true
			}
		}

		// Also dump custom keybindings
		kbPath := filepath.Join(backupPath, "custom_keybindings.dump")
		kbCmd := exec.Command("dconf", "dump", "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/")
		kbOut, err := kbCmd.Output()
		if err == nil && len(kbOut) > 0 {
			_ = os.WriteFile(kbPath, kbOut, 0644)
		}
	}

	// 2. Backup critical user dotfiles
	home := os.Getenv("HOME")
	dotfilesToSave := []string{
		".bashrc",
		".zshrc",
		".tmux.conf",
		filepath.Join(".config", "fish", "config.fish"),
		filepath.Join(".config", "starship.toml"),
		filepath.Join(".config", "nvim", "init.lua"),
		filepath.Join(".config", "vula", "config.yaml"),
		filepath.Join(".config", "vula", "shell.sh"),
	}

	dotfilesDir := filepath.Join(backupPath, "dotfiles")
	for _, relPath := range dotfilesToSave {
		src := filepath.Join(home, relPath)
		if fi, err := os.Stat(src); err == nil && !fi.IsDir() {
			dst := filepath.Join(dotfilesDir, relPath)
			_ = os.MkdirAll(filepath.Dir(dst), 0755)
			if err := copyFile(src, dst); err == nil {
				info.Dotfiles = append(info.Dotfiles, relPath)
			}
		}
	}

	// 3. Save metadata
	metaPath := filepath.Join(backupPath, "metadata.json")
	metaBytes, err := json.MarshalIndent(info, "", "  ")
	if err == nil {
		_ = os.WriteFile(metaPath, metaBytes, 0644)
	}

	return info, nil
}

// List returns all recorded backups sorted by newest first
func (m *Manager) List() ([]BackupInfo, error) {
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	var results []BackupInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		metaPath := filepath.Join(m.baseDir, entry.Name(), "metadata.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}

		var info BackupInfo
		if err := json.Unmarshal(data, &info); err == nil {
			info.Path = filepath.Join(m.baseDir, entry.Name())
			results = append(results, info)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	return results, nil
}

// FindBackup looks up a backup by ID or partial name
func (m *Manager) FindBackup(query string) (*BackupInfo, error) {
	list, err := m.List()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	// 1. Exact ID match
	for _, b := range list {
		if strings.EqualFold(b.ID, q) {
			return &b, nil
		}
	}

	// 2. Exact Name match
	for _, b := range list {
		if strings.EqualFold(b.Name, q) {
			return &b, nil
		}
	}

	// 3. Prefix or contains match
	for _, b := range list {
		if strings.Contains(strings.ToLower(b.ID), q) || strings.Contains(strings.ToLower(b.Name), q) {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("no backup found matching '%s'", query)
}

// Restore restores dotfiles and GNOME dconf settings from a snapshot
func (m *Manager) Restore(query string) (*BackupInfo, error) {
	target, err := m.FindBackup(query)
	if err != nil {
		return nil, err
	}

	home := os.Getenv("HOME")

	// 1. Restore Dotfiles
	dotfilesDir := filepath.Join(target.Path, "dotfiles")
	for _, relPath := range target.Dotfiles {
		src := filepath.Join(dotfilesDir, relPath)
		dst := filepath.Join(home, relPath)
		if _, err := os.Stat(src); err == nil {
			_ = os.MkdirAll(filepath.Dir(dst), 0755)
			_ = copyFile(src, dst)
		}
	}

	// 2. Restore GNOME Dconf settings
	dconfDump := filepath.Join(target.Path, "gnome_dconf.dump")
	if _, err := os.Stat(dconfDump); err == nil {
		dumpData, err := os.ReadFile(dconfDump)
		if err == nil && len(dumpData) > 0 {
			cmd := exec.Command("dconf", "load", "/org/gnome/")
			cmd.Stdin = strings.NewReader(string(dumpData))
			_ = cmd.Run()
		}
	}

	// 3. Restore custom keybindings
	kbDump := filepath.Join(target.Path, "custom_keybindings.dump")
	if _, err := os.Stat(kbDump); err == nil {
		kbData, err := os.ReadFile(kbDump)
		if err == nil && len(kbData) > 0 {
			cmd := exec.Command("dconf", "load", "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/")
			cmd.Stdin = strings.NewReader(string(kbData))
			_ = cmd.Run()
		}
	}

	return target, nil
}

// Delete removes a backup snapshot directory
func (m *Manager) Delete(query string) (*BackupInfo, error) {
	target, err := m.FindBackup(query)
	if err != nil {
		return nil, err
	}

	if err := os.RemoveAll(target.Path); err != nil {
		return nil, fmt.Errorf("failed to delete backup at %s: %w", target.Path, err)
	}

	return target, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
