package apps

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/vula-os/vula/internal/ui"
)

// RunInteractiveAppStore launches a Charm Huh multi-select TUI store
func (m *Manager) RunInteractiveAppStore() error {
	fmt.Println(ui.RenderHeader("Developer App Store", "Select software recipes to install or update"))

	var selectedApps []string

	// Build options with pre-selected state if binary is already installed
	var options []huh.Option[string]
	for _, app := range Catalog {
		installed := isAppInstalled(app)
		label := fmt.Sprintf("%-16s [%-8s] %s", app.Name, app.Category, app.Description)
		if installed {
			label += " (Installed)"
		}
		options = append(options, huh.NewOption(label, app.ID).Selected(installed))
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
	if _, err := exec.LookPath(app.ID); err == nil {
		return true
	}
	// Fallback check for snap or apt
	if app.InstallType == "snap" {
		out, err := exec.Command("snap", "list", app.ID).Output()
		return err == nil && len(out) > 0
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
		return cmd.Run()
	case "curl-sh":
		if app.ID == "starship" {
			return m.InstallCLIStack()
		}
		cmd := exec.Command("sh", "-c", "curl -fsSL "+app.PackageName+" | sh -s -- -y")
		return cmd.Run()
	}
	return nil
}
