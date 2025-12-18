package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robebs/ts-se-tool-go/internal/app"
)

func main() {
	profile := flag.String("profile", "./tmp", "profile directory (contains save/)")
	slot := flag.String("slot", "1", "save slot (e.g. 1, autosave)")
	city2country := flag.String("city2country", "", "optional CityToCountry.csv path")
	gamerefRoot := flag.String("gameref", "", "optional gameref root path")
	gameType := flag.String("gametype", "ETS2", "game type (ETS2 or ATS)")
	flag.Parse()

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
