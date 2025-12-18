package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/discovery"
	"github.com/robebs/ts-se-tool-go/internal/save"
	"github.com/robebs/ts-se-tool-go/internal/util"
)

type SelectedSave struct {
	GameType   string
	ProfileDir string
	SaveSlot   string
}

func selectGameAndProfile() (*SelectedSave, error) {
	// Discover both ETS2 and ATS saves
	custom := discovery.CustomConfig{Paths: map[discovery.GameType][]string{}}

	ets2Saves, _ := discovery.DiscoverSaveSlots(discovery.GameETS2, custom)
	atsSaves, _ := discovery.DiscoverSaveSlots(discovery.GameATS, custom)

	var allSaves []discovery.SaveSlot
	var gameType string

	// Determine which game(s) are available
	if len(ets2Saves) > 0 && len(atsSaves) > 0 {
		// Both games found - ask user to choose
		fmt.Println("\nMultiple games detected:")
		fmt.Println("1. Euro Truck Simulator 2")
		fmt.Println("2. American Truck Simulator")
		fmt.Print("\nSelect game (1 or 2): ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "2" {
			allSaves = atsSaves
			gameType = "ATS"
		} else {
			allSaves = ets2Saves
			gameType = "ETS2"
		}
	} else if len(ets2Saves) > 0 {
		allSaves = ets2Saves
		gameType = "ETS2"
		fmt.Println("\nEuro Truck Simulator 2 detected")
	} else if len(atsSaves) > 0 {
		allSaves = atsSaves
		gameType = "ATS"
		fmt.Println("\nAmerican Truck Simulator detected")
	} else {
		// No games found - ask for profile path
		fmt.Println("\nNo games detected. Please provide profile path manually.")
		fmt.Print("Enter profile directory path: ")
		reader := bufio.NewReader(os.Stdin)
		profilePath, _ := reader.ReadString('\n')
		profilePath = strings.TrimSpace(profilePath)

		fmt.Print("Enter save slot (e.g., 1, autosave): ")
		slot, _ := reader.ReadString('\n')
		slot = strings.TrimSpace(slot)
		if slot == "" {
			slot = "1"
		}

		// Try to determine game type from profile
		gameType = "ETS2" // default

		return &SelectedSave{
			GameType:   gameType,
			ProfileDir: profilePath,
			SaveSlot:   slot,
		}, nil
	}

	// Group saves by profile
	profilesMap := make(map[string][]discovery.SaveSlot)
	for _, save := range allSaves {
		profileKey := save.Profile.Path
		profilesMap[profileKey] = append(profilesMap[profileKey], save)
	}

	// Get unique profiles
	var profiles []discovery.Profile
	seenProfiles := make(map[string]bool)
	for _, save := range allSaves {
		if !seenProfiles[save.Profile.Path] {
			profiles = append(profiles, save.Profile)
			seenProfiles[save.Profile.Path] = true
		}
	}

	// Display profiles with readable names
	fmt.Println("\nAvailable profiles:")
	for i, profile := range profiles {
		profileName := profile.NameHex
		// Try to decode hex name
		if decoded, err := util.HexToString(profile.NameHex); err == nil {
			profileName = fmt.Sprintf("%s (%s)", decoded, profile.NameHex)
		}
		fmt.Printf("[%d] %s\n", i+1, profileName)
	}

	// Let user select profile
	fmt.Print("\nSelect profile (number): ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(profiles) {
		return nil, fmt.Errorf("invalid profile selection")
	}

	selectedProfile := profiles[choice-1]

	// Get save slots for selected profile
	saveSlots := profilesMap[selectedProfile.Path]
	fmt.Println("\nAvailable save slots:")
	for i, slot := range saveSlots {
		fmt.Printf("[%d] %s\n", i+1, slot.SlotName)
	}

	// Let user select save slot
	fmt.Print("\nSelect save slot (number): ")
	slotChoiceStr, _ := reader.ReadString('\n')
	slotChoiceStr = strings.TrimSpace(slotChoiceStr)
	slotChoice, err := strconv.Atoi(slotChoiceStr)
	if err != nil || slotChoice < 1 || slotChoice > len(saveSlots) {
		return nil, fmt.Errorf("invalid save slot selection")
	}

	selectedSlot := saveSlots[slotChoice-1]

	fmt.Printf("\nSelected: %s - Profile: %s - Slot: %s\n",
		gameType, selectedProfile.NameHex, selectedSlot.SlotName)

	return &SelectedSave{
		GameType:   gameType,
		ProfileDir: selectedProfile.Path,
		SaveSlot:   selectedSlot.SlotName,
	}, nil
}

// getProfileName attempts to get a readable profile name from profile.sii
func getProfileName(profileDir string) string {
	doc, err := save.LoadProfileDataFile(profileDir)
	if err != nil {
		return ""
	}

	// Look for profile name in profile blocks
	for _, block := range doc.Blocks {
		if block.Type == "profile" {
			if vals, ok := block.Properties["name"]; ok && len(vals) > 0 {
				return vals[0]
			}
		}
	}

	return ""
}
