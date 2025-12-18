package save

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// WriteGameSII writes a game.sii document to the specified profile and save slot.
// It overwrites the existing file if it exists. The save directory must already exist.
// If encrypt is true, the file is written in encrypted format; otherwise, it's written as plaintext.
func WriteGameSII(profileDir, slot string, doc *sii.Document, encrypt bool) error {
	saveDir := filepath.Join(profileDir, "save", slot)
	gamePath := filepath.Join(saveDir, "game.sii")

	// Ensure save directory exists
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return fmt.Errorf("create save directory: %w", err)
	}

	// Serialize document to SII text
	plaintext, err := sii.WriteDocument(doc)
	if err != nil {
		return fmt.Errorf("serialize SII document: %w", err)
	}

	// Write file (encrypted or plaintext)
	if encrypt {
		if err := siidecrypt.EncryptFile(gamePath, plaintext); err != nil {
			return fmt.Errorf("write encrypted game.sii: %w", err)
		}
	} else {
		if err := os.WriteFile(gamePath, plaintext, 0o644); err != nil {
			return fmt.Errorf("write game.sii: %w", err)
		}
	}

	return nil
}

// WriteInfoSII writes an info.sii document to the specified save directory.
// If the save directory doesn't exist, it will be created (useful for convoy tools).
// If encrypt is true, the file is written in encrypted format; otherwise, it's written as plaintext.
func WriteInfoSII(saveDir string, doc *sii.Document, encrypt bool) error {
	infoPath := filepath.Join(saveDir, "info.sii")

	// Ensure save directory exists
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return fmt.Errorf("create save directory: %w", err)
	}

	// Serialize document to SII text
	plaintext, err := sii.WriteDocument(doc)
	if err != nil {
		return fmt.Errorf("serialize SII document: %w", err)
	}

	// Write file (encrypted or plaintext)
	if encrypt {
		if err := siidecrypt.EncryptFile(infoPath, plaintext); err != nil {
			return fmt.Errorf("write encrypted info.sii: %w", err)
		}
	} else {
		if err := os.WriteFile(infoPath, plaintext, 0o644); err != nil {
			return fmt.Errorf("write info.sii: %w", err)
		}
	}

	return nil
}

// WriteProfileSII writes a profile.sii document to the specified profile directory.
// It overwrites the existing file if it exists. The profile directory must already exist.
// If encrypt is true, the file is written in encrypted format; otherwise, it's written as plaintext.
func WriteProfileSII(profileDir string, doc *sii.Document, encrypt bool) error {
	profilePath := filepath.Join(profileDir, "profile.sii")

	// Serialize document to SII text
	plaintext, err := sii.WriteDocument(doc)
	if err != nil {
		return fmt.Errorf("serialize SII document: %w", err)
	}

	// Write file (encrypted or plaintext)
	if encrypt {
		if err := siidecrypt.EncryptFile(profilePath, plaintext); err != nil {
			return fmt.Errorf("write encrypted profile.sii: %w", err)
		}
	} else {
		if err := os.WriteFile(profilePath, plaintext, 0o644); err != nil {
			return fmt.Errorf("write profile.sii: %w", err)
		}
	}

	return nil
}

// BackupSaveSlot creates a backup of a save slot by renaming it to {slot}.bckp.
// If a backup already exists, it appends a timestamp to avoid conflicts.
// This mirrors the "mv slot slot.bckp" behavior mentioned in TODO.md.
func BackupSaveSlot(profileDir, slot string) error {
	saveDir := filepath.Join(profileDir, "save", slot)
	backupDir := filepath.Join(profileDir, "save", slot+".bckp")

	// Check if save slot exists
	if _, err := os.Stat(saveDir); err != nil {
		return fmt.Errorf("save slot does not exist: %w", err)
	}

	// If backup already exists, append timestamp
	if _, err := os.Stat(backupDir); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupDir = filepath.Join(profileDir, "save", fmt.Sprintf("%s.bckp_%s", slot, timestamp))
	}

	// Rename directory
	if err := os.Rename(saveDir, backupDir); err != nil {
		return fmt.Errorf("rename save slot to backup: %w", err)
	}

	return nil
}
