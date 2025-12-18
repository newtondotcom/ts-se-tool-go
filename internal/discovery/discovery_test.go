package discovery

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestListProfilesAndSaves_FakeTree(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Choose testdata tree depending on OS to exercise ListProfiles/ListSaveSlots.
	var profilesDir string
	if runtime.GOOS == "windows" {
		profilesDir = filepath.Join(cwd, "testdata", "windows", "Documents", "Euro Truck Simulator 2", "profiles")
	} else {
		profilesDir = filepath.Join(cwd, "testdata", "linux", ".local", "share", "Euro Truck Simulator 2", "profiles")
	}

	// Ensure the directory exists; if not, skip.
	if _, err := os.Stat(profilesDir); err != nil {
		t.Skip("testdata profiles dir not present, skipping")
	}

	loc := ProfileLocation{
		Game:        GameETS2,
		Source:      SourceDocuments,
		Root:        filepath.Dir(filepath.Dir(profilesDir)),
		ProfilesDir: profilesDir,
	}

	profiles, err := ListProfiles(loc)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(profiles) == 0 {
		t.Fatalf("expected at least one profile, got 0")
	}

	slots, err := ListSaveSlots(profiles[0])
	if err != nil {
		t.Fatalf("ListSaveSlots: %v", err)
	}
	if len(slots) == 0 {
		t.Fatalf("expected at least one save slot, got 0")
	}

	if filepath.Base(slots[0].GameSII) != "game.sii" {
		t.Errorf("expected game.sii, got %s", slots[0].GameSII)
	}
	if filepath.Base(slots[0].InfoSII) != "info.sii" {
		// We don't create info.sii in the fake tree, but the path should still end with info.sii.
		if filepath.Base(slots[0].InfoSII) != "info.sii" {
			t.Errorf("expected info.sii path, got %s", slots[0].InfoSII)
		}
	}
}
