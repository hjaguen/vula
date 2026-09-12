package profiles

import (
	"testing"
)

func TestGetProfileByID(t *testing.T) {
	p, found := GetProfileByID("full-stack-dev")
	if !found || p == nil {
		t.Fatalf("Expected profile 'full-stack-dev' to be found")
	}
	if p.Name != "Full-Stack Web & Backend Developer" {
		t.Errorf("Unexpected profile name: %s", p.Name)
	}

	_, foundInvalid := GetProfileByID("non-existent-profile")
	if foundInvalid {
		t.Errorf("Expected invalid profile to return false")
	}
}

func TestCatalogIntegrity(t *testing.T) {
	if len(Catalog) < 5 {
		t.Errorf("Expected at least 5 profiles in Catalog, got %d", len(Catalog))
	}

	for _, p := range Catalog {
		if p.ID == "" {
			t.Errorf("Profile ID cannot be empty")
		}
		if p.Name == "" {
			t.Errorf("Profile Name cannot be empty")
		}
		if p.Category == "" {
			t.Errorf("Profile Category cannot be empty")
		}
	}
}
