//go:build windows

package discovery

import (
	"os"
	"path/filepath"
)

// discoverDocumentsProfilesWindows returns profile locations under the user's
// Documents folder for the given game (ETS2/ATS).
func discoverDocumentsProfilesWindows(game GameType) ([]ProfileLocation, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var gameFolder string
	switch game {
	case GameETS2:
		gameFolder = "Euro Truck Simulator 2"
	case GameATS:
		gameFolder = "American Truck Simulator"
	default:
		return nil, nil
	}

	docsRoot := filepath.Join(home, "Documents")
	profilesDir := filepath.Join(docsRoot, gameFolder, "profiles")

	info, err := os.Stat(profilesDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	loc := ProfileLocation{
		Game:        game,
		Source:      SourceDocuments,
		Root:        docsRoot,
		ProfilesDir: profilesDir,
	}
	return []ProfileLocation{loc}, nil
}

// discoverSteamProfilesWindows looks for Steam userdata-based profiles
// following the pattern userdata/<SteamID>/<GameID>/remote/profiles.
func discoverSteamProfilesWindows(game GameType) ([]ProfileLocation, error) {
	// Simple heuristic: look under common Steam install roots. This mirrors
	// TS SE Tool's behaviour conceptually, without reading the registry.
	possible := []string{
		os.Getenv("STEAM_PATH"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramFiles"),
	}

	var roots []string
	for _, base := range possible {
		if base == "" {
			continue
		}
		// If base already points to Steam, use it directly; else append Steam.
		if filepath.Base(base) == "Steam" {
			roots = append(roots, base)
		} else {
			roots = append(roots, filepath.Join(base, "Steam"))
		}
	}

	var gameID string
	switch game {
	case GameETS2:
		gameID = "227300"
	case GameATS:
		gameID = "270880"
	default:
		return nil, nil
	}

	var out []ProfileLocation
	for _, root := range roots {
		userdataDir := filepath.Join(root, "userdata")
		entries, err := os.ReadDir(userdataDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			steamIDDir := filepath.Join(userdataDir, e.Name(), gameID, "remote", "profiles")
			info, err := os.Stat(steamIDDir)
			if err != nil || !info.IsDir() {
				continue
			}
			out = append(out, ProfileLocation{
				Game:        game,
				Source:      SourceSteamCloud,
				Root:        root,
				ProfilesDir: steamIDDir,
			})
		}
	}

	return out, nil
}
