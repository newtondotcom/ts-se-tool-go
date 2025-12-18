package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/app"
	"github.com/robebs/ts-se-tool-go/internal/save"
	"github.com/robebs/ts-se-tool-go/internal/save/loader"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "ts-se-tool",
		Usage:  "Euro Truck Simulator 2 / American Truck Simulator Save Editor",
		Action: runInteractive,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(c *cli.Context) error {
	// Display welcome message
	displayWelcome()

	// Step 1: Game/Profile selection
	selected, err := selectGameAndProfile()
	if err != nil {
		return fmt.Errorf("select game and profile: %w", err)
	}

	// Step 2: Load save file
	fmt.Println("\nLoading save file...")
	docs, err := save.LoadSaveFile(selected.ProfileDir, selected.SaveSlot)
	if err != nil {
		return fmt.Errorf("load save file: %w", err)
	}
	fmt.Println("Save file loaded successfully")

	// Step 3: Backup save slot
	fmt.Println("\nCreating backup...")
	if err := save.BackupSaveSlot(selected.ProfileDir, selected.SaveSlot); err != nil {
		return fmt.Errorf("backup save slot: %w", err)
	}
	fmt.Println("Backup created successfully")

	// Step 4: Load world for garage/truck information
	fmt.Println("\nLoading world data...")
	opts := app.LoadOptions{
		GameType:    selected.GameType,
		ProfilePath: selected.ProfileDir,
		SaveSlot:    selected.SaveSlot,
	}
	w, err := app.LoadWorld(opts)
	if err != nil {
		// World loading is optional, continue without it
		fmt.Printf("Warning: Could not load world data: %v\n", err)
		w = nil
	} else {
		fmt.Println("World data loaded successfully")
	}

	// Step 5: Main menu loop
	modified := false
	for {
		displayMainMenu()
		choice := getUserChoice()

		var err error
		switch choice {
		case 1:
			amount, err := promptMoney()
			if err == nil {
				err = save.SetMoney(docs.Game, amount)
				if err == nil {
					fmt.Printf("Money set to %d\n", amount)
					modified = true
				}
			}
		case 2:
			xp, err := promptXP()
			if err == nil {
				err = save.SetXP(docs.Game, xp)
				if err == nil {
					fmt.Printf("XP set to %d\n", xp)
					modified = true
				}
			}
		case 3:
			err = save.SetSkillsMax(docs.Game)
			if err == nil {
				fmt.Println("All skills set to maximum")
				modified = true
			}
		case 4:
			err = save.BuyAllGarages(docs.Game, w)
			if err == nil {
				fmt.Println("All garages purchased")
				modified = true
			}
		case 5:
			err = save.UpgradeAllGarages(docs.Game)
			if err == nil {
				fmt.Println("All garages upgraded")
				modified = true
			}
		case 6:
			err = save.PopulateGaragesWithTrucks(docs.Game, w)
			if err == nil {
				fmt.Println("Garages populated with trucks")
				modified = true
			}
		case 7:
			err = save.RecruitEmployeesAndPopulateTrucks(docs.Game, w)
			if err == nil {
				fmt.Println("Employees recruited and assigned to trucks")
				modified = true
			}
		case 8:
			// Save and exit
			if modified {
				return saveChanges(selected, docs)
			}
			fmt.Println("No changes to save")
			return nil
		case 9:
			// Exit without saving
			if modified {
				fmt.Println("Warning: You have unsaved changes!")
				fmt.Print("Are you sure you want to exit? (y/n): ")
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					continue
				}
			}
			return nil
		default:
			fmt.Println("Invalid choice. Please select 1-9.")
			continue
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		// Ask if user wants to continue
		if !confirmContinue() {
			if modified {
				return saveChanges(selected, docs)
			}
			return nil
		}
	}
}

func saveChanges(selected *SelectedSave, docs *save.Documents) error {
	fmt.Println("\nSaving changes...")

	// Detect encryption format by checking original file
	encrypt := false
	gamePath := loader.ResolveGameSIIPath(selected.ProfileDir, selected.SaveSlot)
	if info, err := os.Stat(gamePath); err == nil && !info.IsDir() {
		firstBytes := make([]byte, 4)
		if f, err := os.Open(gamePath); err == nil {
			if n, _ := f.Read(firstBytes); n == 4 {
				sig := binary.LittleEndian.Uint32(firstBytes)
				if sig == 1131635539 { // SignatureEncrypted
					encrypt = true
				}
			}
			f.Close()
		}
	}

	// Write all files
	if err := save.WriteSaveFile(selected.ProfileDir, selected.SaveSlot, docs, encrypt); err != nil {
		return fmt.Errorf("write save file: %w", err)
	}

	fmt.Println("Changes saved successfully!")
	return nil
}

func displayWelcome() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                          ║")
	fmt.Println("║     TS SE Tool - Save Game Editor for ETS2/ATS         ║")
	fmt.Println("║                                                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
}
