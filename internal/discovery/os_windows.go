//go:build windows

package discovery

// discoverProfilesForOS aggregates all profile locations for the given
// game on Windows (Documents + Steam userdata).
func discoverProfilesForOS(game GameType) ([]ProfileLocation, error) {
	var all []ProfileLocation
	if docs, err := discoverDocumentsProfilesWindows(game); err == nil {
		all = append(all, docs...)
	}
	if steam, err := discoverSteamProfilesWindows(game); err == nil {
		all = append(all, steam...)
	}
	return all, nil
}
