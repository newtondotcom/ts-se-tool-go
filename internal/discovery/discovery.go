package discovery

// GameType identifies which game a profile/save belongs to.
type GameType string

const (
	GameETS2 GameType = "ETS2"
	GameATS  GameType = "ATS"
)

// SourceKind describes how a profile location was discovered.
type SourceKind string

const (
	SourceDocuments  SourceKind = "documents"
	SourceSteamCloud SourceKind = "steam_cloud"
	SourceCustom     SourceKind = "custom"
)

// ProfileLocation describes a root directory that contains one or more
// ETS2/ATS profiles (each as a hex-named subdirectory with profile.sii).
type ProfileLocation struct {
	Game        GameType
	Source      SourceKind
	Root        string // e.g. Documents root, Steam userdata root, or custom path
	ProfilesDir string // full path to the profiles directory
}

// Profile represents a single game profile (a hex-named folder that
// contains profile.sii and a save/ subdirectory).
type Profile struct {
	Game     GameType
	Location ProfileLocation
	NameHex  string // folder name, e.g. 4D792050726F66696C65
	Path     string // full path to the profile folder
}

// SaveSlot represents a single save slot inside a profile, such as
// save/1, save/autosave, etc.
type SaveSlot struct {
	Profile  Profile
	SlotName string // e.g. "1", "autosave"
	Path     string // full path to the save/<Slot> directory
	GameSII  string // full path to game.sii
	InfoSII  string // full path to info.sii
}
