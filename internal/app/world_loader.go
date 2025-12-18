package app

import (
	"fmt"

	"github.com/robebs/ts-se-tool-go/internal/externaldata"
	"github.com/robebs/ts-se-tool-go/internal/gameref"
	"github.com/robebs/ts-se-tool-go/internal/save/loader"
	"github.com/robebs/ts-se-tool-go/internal/save/world"
)

// LoadOptions describes how to load a world from a profile.
type LoadOptions struct {
	GameType          string // e.g. "ETS2" or "ATS"
	ProfilePath       string // folder containing save/
	SaveSlot          string // e.g. "1", "autosave"
	CityToCountryPath string // optional CityToCountry.csv path
	GameRefRoot       string // optional gameref root path
}

// LoadWorld assembles a full World structure from game.sii and optional
// external data (City->Country mapping, gameref cache).
func LoadWorld(opts LoadOptions) (*world.World, error) {
	gamePath := loader.ResolveGameSIIPath(opts.ProfilePath, opts.SaveSlot)

	w, err := loader.LoadWorldFromGameSII(gamePath)
	if err != nil {
		return nil, err
	}

	// Country dictionary (City -> Country)
	if opts.CityToCountryPath != "" {
		cd, err := externaldata.LoadCountryDictionary(opts.CityToCountryPath)
		if err != nil {
			return nil, fmt.Errorf("load CountryDictionary: %w", err)
		}
		w.CountryDict = cd
	}

	// GameRef cache (external cargo/company data)
	if opts.GameRefRoot != "" && opts.GameType != "" {
		cache, err := gameref.BuildGameRefCache(opts.GameRefRoot, opts.GameType)
		if err != nil {
			return nil, fmt.Errorf("build GameRef cache: %w", err)
		}
		w.GameRef = cache
	}

	return w, nil
}
