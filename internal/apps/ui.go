package apps

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/vula-os/vula/internal/packages"
	"github.com/vula-os/vula/internal/ui"
)

// RunInteractiveAppStore launches a Charm Huh multi-select TUI store
func (m *Manager) RunInteractiveAppStore() error {
	fmt.Println(ui.RenderHeader("Developer App Store", "Select software recipes to install or update"))

	var selectedApps []string

	// Build options with [Installed] label, but DO NOT pre-select them so only user-chosen apps install
	var options []huh.Option[string]
	for _, app := range Catalog {
		installed := isAppInstalled(app)
		label := fmt.Sprintf("%-16s [%-8s] %s", app.Name, app.Category, app.Description)
		if installed {
			label += " (Installed)"
		}
		// Unselected by default so only user-checked options will execute
		options = append(options, huh.NewOption(label, app.ID))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select developer software to configure:").
				Options(options...).
				Value(&selectedApps),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedApps) == 0 {
		fmt.Println(ui.WarnStyle.Render("No apps selected."))
		return nil
	}

	fmt.Printf("\n%s Installing selected software recipes...\n\n", ui.InfoStyle.Render("⚡"))

	for _, id := range selectedApps {
		app := findAppByID(id)
		if app == nil {
			continue
		}

		fmt.Printf("  • Installing %s (%s)...\n", ui.InfoStyle.Render(app.Name), app.InstallType)
		if err := m.installApp(*app); err != nil {
			fmt.Printf("    %s %v\n", ui.ErrorStyle.Render("Failed:"), err)
		} else {
			fmt.Printf("    %s Installed successfully\n", ui.SuccessStyle.Render("✓"))
		}
	}

	fmt.Println("\n" + ui.SuccessStyle.Render("✓ App Store installation batch complete!"))
	return nil
}

func isAppInstalled(app AppRecipe) bool {
	// 1. Binary lookups across known aliases
	binaries := getBinaryCandidates(app.ID)
	for _, bin := range binaries {
		if _, err := exec.LookPath(bin); err == nil {
			return true
		}
	}

	// 2. APT package checks
	aptPkgs := getAptCandidates(app.ID, app.PackageName)
	for _, pkg := range aptPkgs {
		if packages.IsAptInstalled(pkg) {
			return true
		}
	}

	// 3. Snap package checks
	snapNames := getSnapCandidates(app.ID, app.PackageName)
	for _, snapName := range snapNames {
		if checkSnapInstalled(snapName) {
			return true
		}
	}

	// 4. Flatpak checks
	flatpakIDs := getFlatpakCandidates(app.ID)
	for _, fID := range flatpakIDs {
		if checkFlatpakInstalled(fID) {
			return true
		}
	}

	// 5. Desktop file checks
	if checkDesktopFileInstalled(app.ID, binaries) {
		return true
	}

	return false
}

func getBinaryCandidates(id string) []string {
	switch id {
	case "vscode":
		return []string{"code", "code-insiders", "codium", "vscodium"}
	case "dbeaver":
		return []string{"dbeaver", "dbeaver-ce"}
	case "obsidian":
		return []string{"obsidian"}
	case "ripgrep":
		return []string{"rg", "ripgrep"}
	case "neovim":
		return []string{"nvim", "neovim"}
	case "docker":
		return []string{"docker", "dockerd"}
	case "brave":
		return []string{"brave", "brave-browser"}
	default:
		return []string{id}
	}
}

func getAptCandidates(id, packageName string) []string {
	cleanPkg := cleanPackageName(packageName)
	switch id {
	case "docker":
		return []string{"docker-ce", "docker.io", "docker"}
	case "vscode":
		return []string{"code", "code-insiders", "codium"}
	case "dbeaver":
		return []string{"dbeaver-ce", "dbeaver"}
	default:
		if cleanPkg != "" {
			return []string{id, cleanPkg}
		}
		return []string{id}
	}
}

func getSnapCandidates(id, packageName string) []string {
	cleanPkg := cleanPackageName(packageName)
	switch id {
	case "vscode":
		return []string{"code", "code-insiders"}
	case "dbeaver":
		return []string{"dbeaver-ce"}
	default:
		if cleanPkg != "" {
			return []string{id, cleanPkg}
		}
		return []string{id}
	}
}

func getFlatpakCandidates(id string) []string {
	switch id {
	case "vscode":
		return []string{"com.visualstudio.code", "com.visualstudio.code.insiders", "com.vscodium.codium"}
	case "dbeaver":
		return []string{"io.dbeaver.DBeaverCommunity"}
	case "obsidian":
		return []string{"md.obsidian.Obsidian"}
	case "blender":
		return []string{"org.blender.Blender"}
	case "gimp":
		return []string{"org.gimp.GIMP"}
	case "inkscape":
		return []string{"org.inkscape.Inkscape"}
	case "krita":
		return []string{"org.kde.krita"}
	case "brave":
		return []string{"com.brave.Browser"}
	case "vlc":
		return []string{"org.videolan.VLC"}
	default:
		return nil
	}
}

func cleanPackageName(pkg string) string {
	for i, char := range pkg {
		if char == ' ' {
			return pkg[:i]
		}
	}
	return pkg
}

func checkSnapInstalled(snapName string) bool {
	if _, err := exec.LookPath("snap"); err != nil {
		return false
	}
	out, err := exec.Command("snap", "list", snapName).Output()
	return err == nil && len(out) > 0
}

func checkFlatpakInstalled(flatpakID string) bool {
	if _, err := exec.LookPath("flatpak"); err != nil {
		return false
	}
	out, err := exec.Command("flatpak", "info", flatpakID).Output()
	return err == nil && len(out) > 0
}

func checkDesktopFileInstalled(id string, binaries []string) bool {
	searchPaths := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		os.Getenv("HOME") + "/.local/share/applications",
		"/var/lib/flatpak/exports/share/applications",
		"/var/lib/snapd/desktop/applications",
	}

	targets := append([]string{id}, binaries...)
	for _, dir := range searchPaths {
		for _, target := range targets {
			desktopPath := fmt.Sprintf("%s/%s.desktop", dir, target)
			if _, err := os.Stat(desktopPath); err == nil {
				return true
			}
		}
	}
	return false
}

func findAppByID(id string) *AppRecipe {
	for _, a := range Catalog {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

func (m *Manager) installApp(app AppRecipe) error {
	switch app.InstallType {
	case "apt":
		return m.pkg.InstallAptPackages([]string{app.PackageName})
	case "snap":
		cmd := exec.Command("sh", "-c", "sudo snap install "+app.PackageName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case "curl-sh":
		if app.ID == "starship" {
			return m.InstallCLIStack()
		}
		cmd := exec.Command("sh", "-c", "curl -fsSL "+app.PackageName+" | sh -s -- -y")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}
