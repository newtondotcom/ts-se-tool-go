//go:build !windows

package discovery

import (
	"os"
	"path/filepath"
)

// discoverProfilesLinux returns profile locations for ETS2/ATS on Linux.
// It looks under ~/.local/share/<Game>/profiles and optionally under an
// XDG documents directory, plus Steam userdata equivalents.
func discoverProfilesLinux(game GameType) ([]ProfileLocation, error) {
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

	var out []ProfileLocation

	// ~/.local/share/<Game>/profiles
	localShare := filepath.Join(home, ".local", "share", gameFolder, "profiles")
	if info, err := os.Stat(localShare); err == nil && info.IsDir() {
		out = append(out, ProfileLocation{
			Game:        game,
			Source:      SourceDocuments,
			Root:        filepath.Dir(filepath.Dir(localShare)), // ~/.local/share
			ProfilesDir: localShare,
		})
	}

	// XDG_DOCUMENTS_DIR/<Game>/profiles if XDG_DOCUMENTS_DIR is set.
	if docs := os.Getenv("XDG_DOCUMENTS_DIR"); docs != "" {
		profilesDir := filepath.Join(docs, gameFolder, "profiles")
		if info, err := os.Stat(profilesDir); err == nil && info.IsDir() {
			out = append(out, ProfileLocation{
				Game:        game,
				Source:      SourceDocuments,
				Root:        docs,
				ProfilesDir: profilesDir,
			})
		}
	}

	// Steam userdata style: ~/.steam/steam/userdata or ~/.local/share/Steam/userdata
	steamRoots := []string{
		filepath.Join(home, ".steam", "steam"),
		filepath.Join(home, ".local", "share", "Steam"),
	}

	var gameID string
	switch game {
	case GameETS2:
		gameID = "227300"
	case GameATS:
		gameID = "270880"
	}

	for _, root := range steamRoots {
		userdataDir := filepath.Join(root, "userdata")
		entries, err := os.ReadDir(userdataDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			profilesDir := filepath.Join(userdataDir, e.Name(), gameID, "remote", "profiles")
			if info, err := os.Stat(profilesDir); err == nil && info.IsDir() {
				out = append(out, ProfileLocation{
					Game:        game,
					Source:      SourceSteamCloud,
					Root:        root,
					ProfilesDir: profilesDir,
				})
			}
		}
	}

	return out, nil
}
