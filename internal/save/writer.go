package save

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// WriteGameSII writes a game.sii document to the specified profile and save slot.
// It creates a backup (game_backup.sii) before writing. The save directory must already exist.
// If encrypt is true, the file is written in encrypted format; otherwise, it's written as plaintext.
func WriteGameSII(profileDir, slot string, doc *sii.Document, encrypt bool) error {
	saveDir := filepath.Join(profileDir, "save", slot)
	gamePath := filepath.Join(saveDir, "game.sii")

	// Ensure save directory exists
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return fmt.Errorf("create save directory: %w", err)
	}

	// Backup existing file if it exists
	if err := BackupGameSII(profileDir, slot); err != nil {
		return fmt.Errorf("backup game.sii: %w", err)
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
// It creates a backup (info_backup.sii) before writing if the file exists.
// If the save directory doesn't exist, it will be created (useful for convoy tools).
// If encrypt is true, the file is written in encrypted format; otherwise, it's written as plaintext.
func WriteInfoSII(saveDir string, doc *sii.Document, encrypt bool) error {
	infoPath := filepath.Join(saveDir, "info.sii")

	// Ensure save directory exists
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return fmt.Errorf("create save directory: %w", err)
	}

	// Backup existing file if it exists
	if err := BackupInfoSII(saveDir); err != nil {
		return fmt.Errorf("backup info.sii: %w", err)
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
// It creates a backup (profile_backup.sii) before writing. The profile directory must already exist.
// profile.sii is always written in encrypted format, regardless of the encrypt parameter.
func WriteProfileSII(profileDir string, doc *sii.Document, encrypt bool) error {
	profilePath := filepath.Join(profileDir, "profile.sii")

	// Backup existing file if it exists
	if err := BackupProfileSII(profileDir); err != nil {
		return fmt.Errorf("backup profile.sii: %w", err)
	}

	// Serialize document to SII text
	plaintext, err := sii.WriteDocument(doc)
	if err != nil {
		return fmt.Errorf("serialize SII document: %w", err)
	}

	// profile.sii is always encrypted
	if err := siidecrypt.EncryptFile(profilePath, plaintext); err != nil {
		return fmt.Errorf("write encrypted profile.sii: %w", err)
	}

	return nil
}

// BackupGameSII creates a backup of game.sii as game_backup.sii.
// If backup already exists, it appends a timestamp to avoid conflicts.
func BackupGameSII(profileDir, slot string) error {
	saveDir := filepath.Join(profileDir, "save", slot)
	gamePath := filepath.Join(saveDir, "game.sii")
	backupPath := filepath.Join(saveDir, "game_backup.sii")

	// Check if game.sii exists
	if _, err := os.Stat(gamePath); err != nil {
		return fmt.Errorf("game.sii does not exist: %w", err)
	}

	// If backup already exists, append timestamp
	if _, err := os.Stat(backupPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupPath = filepath.Join(saveDir, fmt.Sprintf("game_backup_%s.sii", timestamp))
	}

	// Copy file
	return copyFile(gamePath, backupPath)
}

// BackupInfoSII creates a backup of info.sii as info_backup.sii.
// If backup already exists, it appends a timestamp to avoid conflicts.
func BackupInfoSII(saveDir string) error {
	infoPath := filepath.Join(saveDir, "info.sii")
	backupPath := filepath.Join(saveDir, "info_backup.sii")

	// Check if info.sii exists (it's optional, so return nil if it doesn't exist)
	if _, err := os.Stat(infoPath); err != nil {
		return nil // info.sii is optional, no error if missing
	}

	// If backup already exists, append timestamp
	if _, err := os.Stat(backupPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupPath = filepath.Join(saveDir, fmt.Sprintf("info_backup_%s.sii", timestamp))
	}

	// Copy file
	return copyFile(infoPath, backupPath)
}

// BackupProfileSII creates a backup of profile.sii as profile_backup.sii.
// If backup already exists, it appends a timestamp to avoid conflicts.
func BackupProfileSII(profileDir string) error {
	profilePath := filepath.Join(profileDir, "profile.sii")
	backupPath := filepath.Join(profileDir, "profile_backup.sii")

	// Check if profile.sii exists
	if _, err := os.Stat(profilePath); err != nil {
		return fmt.Errorf("profile.sii does not exist: %w", err)
	}

	// If backup already exists, append timestamp
	if _, err := os.Stat(backupPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupPath = filepath.Join(profileDir, fmt.Sprintf("profile_backup_%s.sii", timestamp))
	}

	// Copy file
	return copyFile(profilePath, backupPath)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copy file contents: %w", err)
	}

	return nil
}
