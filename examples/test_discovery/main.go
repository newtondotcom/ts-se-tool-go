package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/robebs/ts-se-tool-go/internal/discovery"
)

func main() {
	gameType := flag.String("gametype", "ETS2", "Game type: ETS2 or ATS")
	customPath := flag.String("custom", "", "Optional custom profile path to test")
	verbose := flag.Bool("verbose", false, "Show detailed information")
	flag.Parse()

	game := discovery.GameType(*gameType)
	if game != discovery.GameETS2 && game != discovery.GameATS {
		log.Fatalf("invalid game type: %s (must be ETS2 or ATS)", *gameType)
	}

	fmt.Printf("=== Testing Game File Discovery for %s ===\n\n", game)

	// Setup custom config if provided
	custom := discovery.CustomConfig{Paths: map[discovery.GameType][]string{}}
	if *customPath != "" {
		custom.Paths[game] = []string{*customPath}
		fmt.Printf("Custom path configured: %s\n\n", *customPath)
	}

	// Test profile discovery
	fmt.Println("--- Discovering Profiles ---")
	profiles, err := discovery.DiscoverProfiles(game, custom)
	if err != nil {
		log.Printf("Error discovering profiles: %v", err)
	}
	if len(profiles) == 0 {
		fmt.Println("❌ No profiles found!")
		fmt.Println("\nExpected locations:")
		printExpectedLocations(game)
		os.Exit(1)
	}

	fmt.Printf("✅ Found %d profile(s):\n\n", len(profiles))
	for i, p := range profiles {
		fmt.Printf("[%d] Profile: %s\n", i+1, p.NameHex)
		fmt.Printf("    Source: %s\n", p.Location.Source)
		fmt.Printf("    Path: %s\n", p.Path)
		if *verbose {
			fmt.Printf("    Root: %s\n", p.Location.Root)
			fmt.Printf("    Profiles Dir: %s\n", p.Location.ProfilesDir)
		}

		// Check if profile.sii exists
		profileSII := filepath.Join(p.Path, "profile.sii")
		if info, err := os.Stat(profileSII); err == nil && !info.IsDir() {
			fmt.Printf("    ✅ profile.sii found (%d bytes)\n", info.Size())
		} else {
			fmt.Printf("    ❌ profile.sii NOT found\n")
		}

		// Check save directory
		saveDir := filepath.Join(p.Path, "save")
		if info, err := os.Stat(saveDir); err == nil && info.IsDir() {
			fmt.Printf("    ✅ save/ directory found\n")
		} else {
			fmt.Printf("    ❌ save/ directory NOT found\n")
		}
		fmt.Println()
	}

	// Test save slot discovery
	fmt.Println("--- Discovering Save Slots ---")
	saveSlots, err := discovery.DiscoverSaveSlots(game, custom)
	if err != nil {
		log.Printf("Error discovering save slots: %v", err)
	}
	if len(saveSlots) == 0 {
		fmt.Println("❌ No save slots found!")
		os.Exit(1)
	}

	fmt.Printf("✅ Found %d save slot(s):\n\n", len(saveSlots))
	for i, slot := range saveSlots {
		fmt.Printf("[%d] Save Slot: %s\n", i+1, slot.SlotName)
		fmt.Printf("    Profile: %s\n", slot.Profile.NameHex)
		fmt.Printf("    Path: %s\n", slot.Path)

		// Check game.sii
		if info, err := os.Stat(slot.GameSII); err == nil && !info.IsDir() {
			fmt.Printf("    ✅ game.sii found (%d bytes)\n", info.Size())
		} else {
			fmt.Printf("    ❌ game.sii NOT found at: %s\n", slot.GameSII)
		}

		// Check info.sii
		if info, err := os.Stat(slot.InfoSII); err == nil && !info.IsDir() {
			fmt.Printf("    ✅ info.sii found (%d bytes)\n", info.Size())
		} else {
			fmt.Printf("    ⚠️  info.sii NOT found (optional)\n")
		}
		fmt.Println()
	}

	// Summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Profiles discovered: %d\n", len(profiles))
	fmt.Printf("Save slots discovered: %d\n", len(saveSlots))
	fmt.Println("\n✅ Discovery test completed successfully!")
}

func printExpectedLocations(game discovery.GameType) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	var gameFolder string
	var gameID string
	switch game {
	case discovery.GameETS2:
		gameFolder = "Euro Truck Simulator 2"
		gameID = "227300"
	case discovery.GameATS:
		gameFolder = "American Truck Simulator"
		gameID = "270880"
	}

	fmt.Printf("1. Documents folder:\n")
	fmt.Printf("   %s\n", filepath.Join(home, "Documents", gameFolder, "profiles"))
	fmt.Printf("\n2. Steam userdata folders:\n")

	possibleSteamRoots := []string{
		os.Getenv("STEAM_PATH"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramFiles"),
	}

	for _, base := range possibleSteamRoots {
		if base == "" {
			continue
		}
		var steamRoot string
		if filepath.Base(base) == "Steam" {
			steamRoot = base
		} else {
			steamRoot = filepath.Join(base, "Steam")
		}
		fmt.Printf("   %s\\userdata\\<SteamID>\\%s\\remote\\profiles\n", steamRoot, gameID)
	}
	fmt.Println()
}
