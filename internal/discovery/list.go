package discovery

import (
	"os"
	"path/filepath"
)

// ListProfiles returns all profiles found under the given profile location.
func ListProfiles(loc ProfileLocation) ([]Profile, error) {
	entries, err := os.ReadDir(loc.ProfilesDir)
	if err != nil {
		return nil, err
	}
	var out []Profile
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		profPath := filepath.Join(loc.ProfilesDir, e.Name())
		// A valid profile folder must contain profile.sii.
		if info, err := os.Stat(filepath.Join(profPath, "profile.sii")); err != nil || info.IsDir() {
			continue
		}
		out = append(out, Profile{
			Game:     loc.Game,
			Location: loc,
			NameHex:  e.Name(),
			Path:     profPath,
		})
	}
	return out, nil
}

// ListSaveSlots returns all save slots for a given profile, i.e. each
// subdirectory under <profile>/save that contains a game.sii file.
func ListSaveSlots(p Profile) ([]SaveSlot, error) {
	saveDir := filepath.Join(p.Path, "save")
	entries, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, err
	}
	var out []SaveSlot
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		slotPath := filepath.Join(saveDir, e.Name())
		gamePath := filepath.Join(slotPath, "game.sii")
		infoPath := filepath.Join(slotPath, "info.sii")
		if info, err := os.Stat(gamePath); err != nil || info.IsDir() {
			continue
		}
		out = append(out, SaveSlot{
			Profile:  p,
			SlotName: e.Name(),
			Path:     slotPath,
			GameSII:  gamePath,
			InfoSII:  infoPath,
		})
	}
	return out, nil
}

// DiscoverProfiles aggregates all profile locations for the given game and
// returns all discovered profiles.
func DiscoverProfiles(game GameType, custom CustomConfig) ([]Profile, error) {
	var locs []ProfileLocation
	var err error

	// OS-specific discovery happens in platform-specific files that provide
	// discoverProfilesForOS. Here we just call the shared entry point.
	if osLocs, e := discoverProfilesForOS(game); e == nil {
		locs = append(locs, osLocs...)
	} else if err == nil {
		err = e
	}

	// Custom paths (both OSes).
	locs = append(locs, discoverCustomProfiles(custom, game)...)

	var profiles []Profile
	for _, loc := range locs {
		ps, e := ListProfiles(loc)
		if e != nil {
			if err == nil {
				err = e
			}
			continue
		}
		profiles = append(profiles, ps...)
	}
	return profiles, err
}

// DiscoverSaveSlots discovers all profiles for a game and returns all
// save slots found across them.
func DiscoverSaveSlots(game GameType, custom CustomConfig) ([]SaveSlot, error) {
	profiles, err := DiscoverProfiles(game, custom)
	if err != nil {
		return nil, err
	}
	var out []SaveSlot
	for _, p := range profiles {
		slots, e := ListSaveSlots(p)
		if e != nil {
			if err == nil {
				err = e
			}
			continue
		}
		out = append(out, slots...)
	}
	return out, err
}
