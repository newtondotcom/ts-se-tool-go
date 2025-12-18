package save

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/robebs/ts-se-tool-go/internal/save/items"
	"github.com/robebs/ts-se-tool-go/internal/save/world"
	"github.com/robebs/ts-se-tool-go/internal/sii"
)

// Common truck names for ETS2/ATS
var commonTrucks = []string{
	"truck.volvo.fh16",
	"truck.scania.streamline",
	"truck.mercedes.actros",
	"truck.daf.xf",
	"truck.iveco.hiway",
	"truck.man.tgx",
	"truck.renault.magnum",
}

// SetMoney sets the money amount in the bank block.
func SetMoney(doc *sii.Document, amount int64) error {
	bankBlock := findBlockByType(doc, "bank")
	if bankBlock == nil {
		return fmt.Errorf("bank block not found")
	}

	var bank items.Bank
	if err := bank.FromProperties(bankBlock.Properties); err != nil {
		return fmt.Errorf("load bank: %w", err)
	}

	bank.MoneyAccount = amount
	bankBlock.Properties = bank.ToProperties()

	return nil
}

// SetXP sets the experience points in the economy block.
func SetXP(doc *sii.Document, xp uint32) error {
	econBlock := findBlockByType(doc, "economy")
	if econBlock == nil {
		return fmt.Errorf("economy block not found")
	}

	var econ items.Economy
	if err := econ.FromProperties(econBlock.Properties); err != nil {
		return fmt.Errorf("load economy: %w", err)
	}

	econ.ExperiencePoints = xp
	econBlock.Properties = econ.ToProperties()

	return nil
}

// SetSkillsMax sets all skills to maximum (255) in the economy block.
func SetSkillsMax(doc *sii.Document) error {
	econBlock := findBlockByType(doc, "economy")
	if econBlock == nil {
		return fmt.Errorf("economy block not found")
	}

	var econ items.Economy
	if err := econ.FromProperties(econBlock.Properties); err != nil {
		return fmt.Errorf("load economy: %w", err)
	}

	// Set all skills to max (255)
	econ.Adr = 255
	econ.LongDist = 255
	econ.Heavy = 255
	econ.Fragile = 255
	econ.Urgent = 255
	econ.Mechanical = 255

	econBlock.Properties = econ.ToProperties()

	return nil
}

// BuyAllGarages adds all available garages to the economy block.
func BuyAllGarages(doc *sii.Document, w *world.World) error {
	econBlock := findBlockByType(doc, "economy")
	if econBlock == nil {
		return fmt.Errorf("economy block not found")
	}

	var econ items.Economy
	if err := econ.FromProperties(econBlock.Properties); err != nil {
		return fmt.Errorf("load economy: %w", err)
	}

	// Get all garage names from world or find all garage blocks
	garageMap := make(map[string]bool)
	for _, g := range econ.Garages {
		garageMap[g] = true
	}

	// Find all garage blocks in document
	for _, block := range doc.Blocks {
		if block.Type == "garage" {
			if !garageMap[block.Name] {
				econ.Garages = append(econ.Garages, block.Name)
				garageMap[block.Name] = true
			}
		}
	}

	// Also add garages from world if available
	if w != nil {
		for _, g := range w.Garages {
			garageName := "garage." + g.Name
			if !garageMap[garageName] {
				econ.Garages = append(econ.Garages, garageName)
				garageMap[garageName] = true
			}
		}
	}

	econBlock.Properties = econ.ToProperties()

	return nil
}

// UpgradeAllGarages sets all garage status to maximum (2 = fully upgraded).
func UpgradeAllGarages(doc *sii.Document) error {
	upgraded := 0
	for i := range doc.Blocks {
		if doc.Blocks[i].Type == "garage" {
			var garage items.Garage
			if err := garage.FromProperties(doc.Blocks[i].Properties); err != nil {
				continue
			}

			// Status 2 = fully upgraded garage
			garage.Status = 2
			doc.Blocks[i].Properties = garage.ToProperties()
			upgraded++
		}
	}

	if upgraded == 0 {
		return fmt.Errorf("no garage blocks found")
	}

	return nil
}

// PopulateGaragesWithTrucks adds random trucks to all garages.
func PopulateGaragesWithTrucks(doc *sii.Document, w *world.World) error {
	// Get available trucks
	trucks := getAvailableTrucks(w)

	if len(trucks) == 0 {
		// Fallback to common trucks
		trucks = commonTrucks
	}

	rand.Seed(time.Now().UnixNano())
	populated := 0

	for i := range doc.Blocks {
		if doc.Blocks[i].Type == "garage" {
			var garage items.Garage
			if err := garage.FromProperties(doc.Blocks[i].Properties); err != nil {
				continue
			}

			// Garage capacity is typically 5 slots
			capacity := 5
			currentCount := len(garage.Vehicles)

			// Add trucks up to capacity
			for currentCount < capacity {
				truck := trucks[rand.Intn(len(trucks))]
				garage.Vehicles = append(garage.Vehicles, truck)
				currentCount++
			}

			doc.Blocks[i].Properties = garage.ToProperties()
			populated++
		}
	}

	if populated == 0 {
		return fmt.Errorf("no garage blocks found")
	}

	return nil
}

// RecruitEmployeesAndPopulateTrucks creates drivers and assigns them to trucks.
func RecruitEmployeesAndPopulateTrucks(doc *sii.Document, w *world.World) error {
	rand.Seed(time.Now().UnixNano())

	// Get player block to add drivers
	playerBlock := findBlockByType(doc, "player")
	if playerBlock == nil {
		return fmt.Errorf("player block not found")
	}

	var player items.Player
	if err := player.FromProperties(playerBlock.Properties); err != nil {
		return fmt.Errorf("load player: %w", err)
	}

	// Get available trucks
	trucks := getAvailableTrucks(w)
	if len(trucks) == 0 {
		trucks = commonTrucks
	}

	driverCounter := 1
	driversCreated := 0

	// Process all garages
	for i := range doc.Blocks {
		if doc.Blocks[i].Type == "garage" {
			var garage items.Garage
			if err := garage.FromProperties(doc.Blocks[i].Properties); err != nil {
				continue
			}

			// Assign drivers to trucks in this garage
			for j, truck := range garage.Vehicles {
				if truck == "" || truck == "null" {
					continue
				}

				// Check if driver already assigned
				if j < len(garage.Drivers) && garage.Drivers[j] != "" && garage.Drivers[j] != "null" {
					continue
				}

				// Create new driver
				driverName := fmt.Sprintf("driver.name_%04d", driverCounter)
				driverCounter++

				// Create driver_ai block
				driverBlock := createDriverAIBlock(driverName, truck)
				doc.Blocks = append(doc.Blocks, driverBlock)

				// Ensure drivers array is large enough
				for len(garage.Drivers) <= j {
					garage.Drivers = append(garage.Drivers, "")
				}
				garage.Drivers[j] = driverName

				// Add to player's drivers list
				player.Drivers = append(player.Drivers, driverName)
				player.DriverReadinessTimer = append(player.DriverReadinessTimer, 0)
				player.DriverQuitWarned = append(player.DriverQuitWarned, false)

				driversCreated++
			}

			doc.Blocks[i].Properties = garage.ToProperties()
		}
	}

	// Update player block
	playerBlock.Properties = player.ToProperties()

	if driversCreated == 0 {
		return fmt.Errorf("no trucks found to assign drivers")
	}

	return nil
}

// Helper functions

func findBlockByType(doc *sii.Document, blockType string) *sii.Block {
	for i := range doc.Blocks {
		if doc.Blocks[i].Type == blockType {
			return &doc.Blocks[i]
		}
	}
	return nil
}

func getAvailableTrucks(w *world.World) []string {
	var trucks []string

	if w != nil && w.Player != nil {
		// Get trucks from player
		trucks = append(trucks, w.Player.Trucks...)
	}

	// Also check world's player trucks map
	if w != nil {
		for truckName := range w.PlayerTrucks {
			// Avoid duplicates
			found := false
			for _, t := range trucks {
				if t == truckName {
					found = true
					break
				}
			}
			if !found {
				trucks = append(trucks, truckName)
			}
		}
	}

	return trucks
}

func createDriverAIBlock(name, truck string) sii.Block {
	driver := items.DriverAI{
		Adr:                     50, // Default skill values
		LongDist:                50,
		Heavy:                   50,
		Fragile:                 50,
		Urgent:                  50,
		Mechanical:              50,
		Hometown:                "berlin", // Default city
		CurrentCity:             "berlin",
		State:                   0,
		OnDutyTimer:             0,
		ExtraMaintenance:        0,
		DriverJob:               "",
		ExperiencePoints:        1000,
		TrainingPolicy:          0,
		AssignedTruck:           truck,
		AssignedTruckEfficiency: 1.0,
		AssignedTruckAxleCount:  2,
		AssignedTruckMass:       0.0,
		SlotTruckEfficiency:     1.0,
		SlotTruckAxleCount:      2,
		SlotTruckMass:           0.0,
		AdoptedTrailer:          "",
		AssignedTrailer:         "",
		OldHometown:             "",
		ProfitLog:               "",
	}

	return sii.Block{
		Type:       "driver_ai",
		Name:       name,
		Properties: driver.ToProperties(),
	}
}
