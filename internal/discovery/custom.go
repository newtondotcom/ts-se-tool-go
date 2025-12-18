package discovery

import (
	"os"
	"path/filepath"
)

// CustomConfig describes user-provided additional profile roots, grouped
// by game type, similar in spirit to ProgSettingsV.CustomPaths.
type CustomConfig struct {
	Paths map[GameType][]string
}

// discoverCustomProfiles returns profile locations for custom paths
// supplied by the user in CustomConfig.
func discoverCustomProfiles(custom CustomConfig, game GameType) []ProfileLocation {
	var out []ProfileLocation
	roots := custom.Paths[game]
	for _, root := range roots {
		if root == "" {
			continue
		}
		// First try <root>/profiles.
		profilesDir := filepath.Join(root, "profiles")
		if info, err := os.Stat(profilesDir); err == nil && info.IsDir() {
			out = append(out, ProfileLocation{
				Game:        game,
				Source:      SourceCustom,
				Root:        root,
				ProfilesDir: profilesDir,
			})
			continue
		}
		// Fallback: treat root itself as a profiles dir if it contains
		// at least one subdirectory with a profile.sii file.
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			candidate := filepath.Join(root, e.Name(), "profile.sii")
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				out = append(out, ProfileLocation{
					Game:        game,
					Source:      SourceCustom,
					Root:        root,
					ProfilesDir: root,
				})
				break
			}
		}
	}
	return out
}
