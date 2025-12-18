package save

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
	"github.com/robebs/ts-se-tool-go/internal/util"
)

// RenameProfile renames a profile by creating a new directory with the hex-encoded
// new name, copying all contents, writing the updated profile.sii, and deleting
// the old directory. This mirrors the C# FormProfileEditorRenameClone behavior.
func RenameProfile(oldProfileDir, newProfileName, profilesRoot string) (string, error) {
	// Step 1: Convert new name to hex
	hexName := util.StringToHex(newProfileName)
	newProfileDir := filepath.Join(profilesRoot, hexName)

	// Step 2: Check if new profile already exists
	if _, err := os.Stat(newProfileDir); err == nil {
		return "", fmt.Errorf("profile with name %s already exists at %s", newProfileName, newProfileDir)
	}

	// Step 3: Copy entire directory tree
	if err := util.CopyDirectory(oldProfileDir, newProfileDir); err != nil {
		return "", fmt.Errorf("copy profile directory: %w", err)
	}

	// Step 4: Load and update profile.sii
	profileDoc, err := LoadProfileDataFile(newProfileDir)
	if err != nil {
		// Clean up on error
		_ = os.RemoveAll(newProfileDir)
		return "", fmt.Errorf("load profile.sii: %w", err)
	}

	// Update profile name in the document
	// The profile name is typically stored in a "profile" block with a "name" property
	// We'll update all blocks that might contain the profile name
	updated := updateProfileNameInDocument(profileDoc, newProfileName)

	// Step 5: Write updated profile.sii
	// Determine if original was encrypted by checking if it exists and trying to read it
	encrypt := false
	originalProfilePath := filepath.Join(oldProfileDir, "profile.sii")
	if _, err := os.Stat(originalProfilePath); err == nil {
		// Try to detect encryption by reading first 4 bytes
		firstBytes := make([]byte, 4)
		if f, err := os.Open(originalProfilePath); err == nil {
			if n, _ := f.Read(firstBytes); n == 4 {
				sig := binary.LittleEndian.Uint32(firstBytes)
				if sig == uint32(siidecrypt.SignatureEncrypted) {
					encrypt = true
				}
			}
			f.Close()
		}
	}

	if updated {
		if err := WriteProfileSII(newProfileDir, profileDoc, encrypt); err != nil {
			// Clean up on error
			_ = os.RemoveAll(newProfileDir)
			return "", fmt.Errorf("write updated profile.sii: %w", err)
		}
	}

	// Step 6: Delete old directory only after successful copy and write
	if err := os.RemoveAll(oldProfileDir); err != nil {
		// Log error but don't fail - the rename was successful
		return newProfileDir, fmt.Errorf("remove old profile directory (rename succeeded): %w", err)
	}

	return newProfileDir, nil
}

// CloneProfile creates a copy of a profile with a new name, similar to RenameProfile
// but without deleting the source profile. Useful for creating backup profiles.
func CloneProfile(sourceProfileDir, newProfileName, profilesRoot string) (string, error) {
	// Step 1: Convert new name to hex
	hexName := util.StringToHex(newProfileName)
	newProfileDir := filepath.Join(profilesRoot, hexName)

	// Step 2: Check if new profile already exists
	if _, err := os.Stat(newProfileDir); err == nil {
		return "", fmt.Errorf("profile with name %s already exists at %s", newProfileName, newProfileDir)
	}

	// Step 3: Copy entire directory tree
	if err := util.CopyDirectory(sourceProfileDir, newProfileDir); err != nil {
		return "", fmt.Errorf("copy profile directory: %w", err)
	}

	// Step 4: Load and update profile.sii
	profileDoc, err := LoadProfileDataFile(newProfileDir)
	if err != nil {
		// Clean up on error
		_ = os.RemoveAll(newProfileDir)
		return "", fmt.Errorf("load profile.sii: %w", err)
	}

	// Update profile name in the document
	updated := updateProfileNameInDocument(profileDoc, newProfileName)

	// Step 5: Write updated profile.sii
	// Determine if original was encrypted
	encrypt := false
	originalProfilePath := filepath.Join(sourceProfileDir, "profile.sii")
	if _, err := os.Stat(originalProfilePath); err == nil {
		firstBytes := make([]byte, 4)
		if f, err := os.Open(originalProfilePath); err == nil {
			if n, _ := f.Read(firstBytes); n == 4 {
				sig := binary.LittleEndian.Uint32(firstBytes)
				if sig == uint32(siidecrypt.SignatureEncrypted) {
					encrypt = true
				}
			}
			f.Close()
		}
	}

	if updated {
		if err := WriteProfileSII(newProfileDir, profileDoc, encrypt); err != nil {
			// Clean up on error
			_ = os.RemoveAll(newProfileDir)
			return "", fmt.Errorf("write updated profile.sii: %w", err)
		}
	}

	return newProfileDir, nil
}

// updateProfileNameInDocument updates the profile name in the SII document.
// It searches for blocks that might contain the profile name and updates them.
// Returns true if any updates were made.
func updateProfileNameInDocument(doc *sii.Document, newName string) bool {
	updated := false
	for i := range doc.Blocks {
		block := &doc.Blocks[i]
		// Common property names that might contain the profile name
		if vals, ok := block.Properties["name"]; ok && len(vals) > 0 {
			block.Properties["name"] = []string{newName}
			updated = true
		}
		if vals, ok := block.Properties["profile_name"]; ok && len(vals) > 0 {
			block.Properties["profile_name"] = []string{newName}
			updated = true
		}
		// For "profile" type blocks, also check common variations
		if block.Type == "profile" {
			if vals, ok := block.Properties["name"]; ok && len(vals) > 0 {
				block.Properties["name"] = []string{newName}
				updated = true
			}
		}
	}
	return updated
}
