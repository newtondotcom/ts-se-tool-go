package save

// This package defines high-level save file structures that correspond to
// CustomClasses/Save/SaveFileInfoData.cs and SaveFileProfileData.cs.

import (
	"fmt"
	"path/filepath"

	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// FileInfoData is a placeholder for the information held in SaveFileInfoData.
// Once the C# class is fully ported we can replace this with a richer struct.
type FileInfoData struct{}

// FileProfileData is a placeholder for the information held in SaveFileProfileData.
type FileProfileData struct{}

// Documents groups the decoded SII documents that make up a save slot.
// It is the UI-free equivalent of what LoadSaveFile + NewDecodeFile do
// in the C# code.
type Documents struct {
	Profile *sii.Document
	Info    *sii.Document
	Game    *sii.Document
}

// decodeSiiDocument decrypts (if needed) and parses a single SII file into
// a generic sii.Document.
func decodeSiiDocument(path string) (*sii.Document, error) {
	plain, err := siidecrypt.DecryptFile(path, true)
	if err != nil {
		return nil, fmt.Errorf("decrypt %s: %w", path, err)
	}
	doc, err := sii.ReadDocument(plain)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return doc, nil
}

// LoadProfileDataFile is a Go analogue of the C# LoadProfileDataFile method,
// but it only deals with data: given a profile directory, it decodes
// profile.sii and returns a generic SII document.
func LoadProfileDataFile(profileDir string) (*sii.Document, error) {
	path := filepath.Join(profileDir, "profile.sii")
	return decodeSiiDocument(path)
}

// LoadSaveFile is a data-only variant of the C# LoadSaveFile method. It
// decodes profile.sii, info.sii and game.sii for a given profile directory
// and save slot, and returns them as parsed SII documents.
func LoadSaveFile(profileDir, slot string) (*Documents, error) {
	profilePath := filepath.Join(profileDir, "profile.sii")
	saveDir := filepath.Join(profileDir, "save", slot)
	infoPath := filepath.Join(saveDir, "info.sii")
	gamePath := filepath.Join(saveDir, "game.sii")

	profileDoc, err := decodeSiiDocument(profilePath)
	if err != nil {
		return nil, err
	}
	infoDoc, err := decodeSiiDocument(infoPath)
	if err != nil {
		return nil, err
	}
	gameDoc, err := decodeSiiDocument(gamePath)
	if err != nil {
		return nil, err
	}

	return &Documents{
		Profile: profileDoc,
		Info:    infoDoc,
		Game:    gameDoc,
	}, nil
}

// WriteSaveFile writes all three SII files (profile.sii, info.sii, game.sii) for a given
// profile directory and save slot. This is a convenience function that calls the individual
// write functions. If encrypt is true, files are written in encrypted format.
func WriteSaveFile(profileDir, slot string, docs *Documents, encrypt bool) error {
	// Write profile.sii
	if docs.Profile != nil {
		if err := WriteProfileSII(profileDir, docs.Profile, encrypt); err != nil {
			return fmt.Errorf("write profile.sii: %w", err)
		}
	}

	// Write game.sii and info.sii
	if docs.Game != nil {
		if err := WriteGameSII(profileDir, slot, docs.Game, encrypt); err != nil {
			return fmt.Errorf("write game.sii: %w", err)
		}
	}

	if docs.Info != nil {
		saveDir := filepath.Join(profileDir, "save", slot)
		if err := WriteInfoSII(saveDir, docs.Info, encrypt); err != nil {
			return fmt.Errorf("write info.sii: %w", err)
		}
	}

	return nil
}
