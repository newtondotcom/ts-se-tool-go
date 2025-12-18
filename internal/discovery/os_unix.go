//go:build !windows

package discovery

// discoverProfilesForOS aggregates all profile locations for the given
// game on Unix-like systems (Linux, etc.).
func discoverProfilesForOS(game GameType) ([]ProfileLocation, error) {
	return discoverProfilesLinux(game)
}
