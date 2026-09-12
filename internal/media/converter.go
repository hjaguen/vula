package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// ConvertWebmToMp4 converts a .webm screen recording to a web-friendly .mp4 via ffmpeg
func ConvertWebmToMp4(inputFile, outputFile string) (string, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg is required for media conversion. Please install it with 'sudo apt install ffmpeg'")
	}

	targetInput := strings.TrimSpace(inputFile)
	if targetInput == "" {
		latest, err := findLatestWebmFile()
		if err != nil {
			return "", err
		}
		targetInput = latest
	}

	if _, err := os.Stat(targetInput); os.IsNotExist(err) {
		return "", fmt.Errorf("input video file '%s' does not exist", targetInput)
	}

	targetOutput := strings.TrimSpace(outputFile)
	if targetOutput == "" {
		ext := filepath.Ext(targetInput)
		targetOutput = strings.TrimSuffix(targetInput, ext) + ".mp4"
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", targetInput, "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", targetOutput)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w (output: %s)", err, string(out))
	}

	return targetOutput, nil
}

func findLatestWebmFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	searchDirs := []string{
		filepath.Join(home, "Videos", "Screencasts"),
		filepath.Join(home, "Videos"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Desktop"),
	}

	type fileModInfo struct {
		path    string
		modTime os.FileInfo
	}

	var candidates []fileModInfo
	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".webm") {
				fullPath := filepath.Join(dir, entry.Name())
				if info, err := entry.Info(); err == nil {
					candidates = append(candidates, fileModInfo{path: fullPath, modTime: info})
				}
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no .webm screen recordings found in ~/Videos or ~/Downloads")
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].modTime.ModTime().After(candidates[j].modTime.ModTime())
	})

	return candidates[0].path, nil
}
