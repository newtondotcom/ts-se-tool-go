package loader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/items"
	"github.com/robebs/ts-se-tool-go/internal/save/itemsextra"
	"github.com/robebs/ts-se-tool-go/internal/save/world"
	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// LoadWorldFromGameSII loads a World structure from a single game.sii file.
// It is inspired by the C# NewPrepareData() + Prepare*Initial methods.
func LoadWorldFromGameSII(path string) (*world.World, error) {
	plain, err := siidecrypt.DecryptFile(path, true)
	if err != nil {
		return nil, fmt.Errorf("decrypt game.sii: %w", err)
	}

	doc, err := sii.ReadDocument(plain)
	if err != nil {
		return nil, fmt.Errorf("parse game.sii: %w", err)
	}

	w := world.NewWorld(doc)

	// Core typed blocks
	econ, err := doc.ToEconomy()
	if err != nil {
		return nil, fmt.Errorf("load economy: %w", err)
	}
	w.Economy = econ

	player, err := doc.ToPlayer()
	if err != nil {
		return nil, fmt.Errorf("load player: %w", err)
	}
	w.Player = player

	w.BlocksByID = doc.IndexByName()

	if econ != nil {
		prepareCitiesAndCompanies(w, econ)
		prepareGarages(w, econ)
	}
	if econ != nil {
		prepareCargoTrailerAndCompanyTrucks(w, econ)
	}
	if player != nil {
		preparePlayerFleet(w, player)
	}

	return w, nil
}

// prepareCitiesAndCompanies mirrors the logic of PrepareCitiesInitial().
func prepareCitiesAndCompanies(w *world.World, econ *items.Economy) {
	cityByName := make(map[string]*itemsextra.City)
	companyByName := make(map[string]*itemsextra.Company)

	for _, compName := range econ.Companies {
		parts := strings.Split(compName, ".")
		if len(parts) < 4 {
			continue
		}
		// Example: company.ikea.city.bruxelles
		companyID := parts[1]
		cityID := parts[3]

		city := cityByName[cityID]
		if city == nil {
			city = &itemsextra.City{
				Name:                    cityID,
				Companies:               []string{},
				JobOffersCountByCompany: make(map[string]int),
			}
			cityByName[cityID] = city
			w.Cities = append(w.Cities, city)
		}
		// link company to city
		city.Companies = append(city.Companies, companyID)

		comp := companyByName[companyID]
		if comp == nil {
			comp = &itemsextra.Company{
				Name:     companyID,
				CityName: cityID,
			}
			companyByName[companyID] = comp
			w.Companies = append(w.Companies, comp)
		}

		// Optional: count job offers for this company
		if blk, ok := w.BlocksByID[compName]; ok && blk.Type == "company" {
			var c items.Company
			if err := c.FromProperties(blk.Properties); err == nil {
				count := len(c.JobOffer)
				city.JobOffersCountByCompany[companyID] += count
			}
		}
	}
}

// prepareGarages mirrors PrepareGaragesInitial().
func prepareGarages(w *world.World, econ *items.Economy) {
	for _, gname := range econ.Garages {
		blk, ok := w.BlocksByID[gname]
		if !ok || blk.Type != "garage" {
			continue
		}
		var g items.Garage
		if err := g.FromProperties(blk.Properties); err != nil {
			continue
		}

		nameParts := strings.Split(gname, ".")
		shortName := gname
		if len(nameParts) > 1 {
			shortName = nameParts[1]
		}

		wg := &itemsextra.Garages{
			Name:     shortName,
			Status:   g.Status,
			Vehicles: append([]string{}, g.Vehicles...),
			Drivers:  append([]string{}, g.Drivers...),
			Trailers: append([]string{}, g.Trailers...),
		}
		w.Garages = append(w.Garages, wg)
	}
}

// prepareCargoTrailerAndCompanyTrucks mirrors PrepareCargoTrailerDefsVariantsLists().
func prepareCargoTrailerAndCompanyTrucks(w *world.World, econ *items.Economy) {
	cargoIndex := make(map[string]*itemsextra.Cargo)
	trailerDefIndex := make(map[string]*itemsextra.TrailerDefinition)
	companyTrucks := make(map[string]*itemsextra.CompanyTruck)

	for _, compName := range econ.Companies {
		blk, ok := w.BlocksByID[compName]
		if !ok || blk.Type != "company" {
			continue
		}
		var c items.Company
		if err := c.FromProperties(blk.Properties); err != nil {
			continue
		}

		for _, jobName := range c.JobOffer {
			jobBlk, ok := w.BlocksByID[jobName]
			if !ok || jobBlk.Type != "job_offer_data" {
				continue
			}
			var j items.JobOfferData
			if err := j.FromProperties(jobBlk.Properties); err != nil {
				continue
			}

			// cargo: "cargo.furniture" -> "furniture"
			cargoID := j.Cargo
			if parts := strings.Split(j.Cargo, "."); len(parts) > 1 {
				cargoID = parts[1]
			}

			cargotype := 0
			if strings.Contains(j.CompanyTruck, "\"heavy") {
				cargotype = 1
			} else if strings.Contains(j.CompanyTruck, "\"double") {
				cargotype = 2
			}

			cg := cargoIndex[cargoID]
			if cg == nil {
				cg = &itemsextra.Cargo{
					ID:             cargoID,
					CargoType:      cargotype,
					TrailerDefName: j.TrailerDefinition,
					UnitsCount:     j.UnitsCount,
				}
				cargoIndex[cargoID] = cg
				w.Cargoes = append(w.Cargoes, cg)
			}

			// trailer definitions and variants
			if j.TrailerDefinition != "" {
				td := trailerDefIndex[j.TrailerDefinition]
				if td == nil {
					td = &itemsextra.TrailerDefinition{Name: j.TrailerDefinition}
					trailerDefIndex[j.TrailerDefinition] = td
					w.TrailerDefs = append(w.TrailerDefs, td)
				}
				if j.TrailerVariant != "" {
					// avoid duplicates
					found := false
					for _, v := range td.Variants {
						if v == j.TrailerVariant {
							found = true
							break
						}
					}
					if !found {
						td.Variants = append(td.Variants, j.TrailerVariant)
					}
				}
			}

			// company trucks
			if j.CompanyTruck != "" {
				ct := companyTrucks[j.CompanyTruck]
				if ct == nil {
					ct = &itemsextra.CompanyTruck{
						TruckID:   j.CompanyTruck,
						CargoType: cargotype,
					}
					companyTrucks[j.CompanyTruck] = ct
					w.CompanyTrucks = append(w.CompanyTrucks, ct)
				}
			}
		}
	}
}

// preparePlayerFleet mirrors PreparePlayerDictionariesInitial().
func preparePlayerFleet(w *world.World, p *items.Player) {
	// Player trucks
	for _, tname := range p.Trucks {
		w.PlayerTrucks[tname] = &itemsextra.UserCompanyTruckData{TruckName: tname}
	}

	// Player trailers
	for _, trName := range p.Trailers {
		w.PlayerTrailers[trName] = &itemsextra.UserCompanyTrailerData{TrailerName: trName}
	}

	// Trailer defs
	for _, defName := range p.TrailerDefs {
		// We keep the mapping on the trailer side if both are present.
		if td, ok := w.PlayerTrailers[defName]; ok {
			// defName is itself a trailer name in some cases; do nothing special.
			_ = td
		}
	}
}

// ResolveGameSIIPath is a small helper which builds a game.sii path
// from a profile directory and slot name (e.g. "1", "autosave").
func ResolveGameSIIPath(profileDir, slot string) string {
	return filepath.Join(profileDir, "save", slot, "game.sii")
}
