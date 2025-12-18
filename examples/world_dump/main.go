package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robebs/ts-se-tool-go/internal/app"
	"github.com/robebs/ts-se-tool-go/internal/discovery"
)

func main() {
	profile := flag.String("profile", "", "profile directory (contains save/)")
	slot := flag.String("slot", "", "save slot (e.g. 1, autosave); if empty, auto-select")
	city2country := flag.String("city2country", "", "optional CityToCountry.csv path")
	gamerefRoot := flag.String("gameref", "", "optional gameref root path")
	gameType := flag.String("gametype", "ETS2", "game type (ETS2 or ATS)")
	listOnly := flag.Bool("list", false, "only list detected profiles and saves, do not load world")
	flag.Parse()

	game := discovery.GameType(*gameType)

	// If no profile is provided, try to discover profiles/saves automatically.
	if *profile == "" {
		custom := discovery.CustomConfig{Paths: map[discovery.GameType][]string{}}
		saves, err := discovery.DiscoverSaveSlots(game, custom)
		if err != nil {
			log.Fatalf("discover saves: %v", err)
		}
		if len(saves) == 0 {
			log.Fatalf("no saves found for game %s", *gameType)
		}

		fmt.Println("Discovered save slots:")
		for i, s := range saves {
			fmt.Printf("[%d] game=%s source=%s profile=%s slot=%s path=%s\n",
				i,
				s.Profile.Game,
				s.Profile.Location.Source,
				s.Profile.NameHex,
				s.SlotName,
				s.Path,
			)
		}

		if *listOnly {
			return
		}

		// Auto-select the first save slot.
		chosen := saves[0]
		*profile = chosen.Profile.Path
		*slot = chosen.SlotName
		fmt.Printf("\nAuto-selected profile=%s slot=%s\n", chosen.Profile.NameHex, chosen.SlotName)
	} else if *slot == "" {
		// profile provided but no slot: default to "1" as a common case.
		*slot = "1"
	}

	opts := app.LoadOptions{
		GameType:          *gameType,
		ProfilePath:       *profile,
		SaveSlot:          *slot,
		CityToCountryPath: *city2country,
		GameRefRoot:       *gamerefRoot,
	}

	w, err := app.LoadWorld(opts)
	if err != nil {
		log.Fatalf("load world: %v", err)
	}

	fmt.Printf("Cities: %d\n", len(w.Cities))
	for i, c := range w.Cities {
		if i >= 5 {
			break
		}
		fmt.Printf("  - %s (country=%s, companies=%d)\n", c.Name, c.Country, len(c.Companies))
	}

	fmt.Printf("Companies: %d\n", len(w.Companies))
	for i, c := range w.Companies {
		if i >= 5 {
			break
		}
		fmt.Printf("  - %s (city=%s, jobOffers=%d)\n", c.Name, c.CityName, len(c.JobOffers))
	}

	fmt.Printf("Cargoes: %d\n", len(w.Cargoes))
	fmt.Printf("Trailer definitions: %d\n", len(w.TrailerDefs))
	fmt.Printf("Company trucks: %d\n", len(w.CompanyTrucks))
	fmt.Printf("Garages: %d\n", len(w.Garages))
	fmt.Printf("Player trucks: %d\n", len(w.PlayerTrucks))
	fmt.Printf("Player trailers: %d\n", len(w.PlayerTrailers))
}
