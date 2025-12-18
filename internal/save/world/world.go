package world

import (
	"github.com/robebs/ts-se-tool-go/internal/externaldata"
	"github.com/robebs/ts-se-tool-go/internal/save/items"
	"github.com/robebs/ts-se-tool-go/internal/save/itemsextra"
	"github.com/robebs/ts-se-tool-go/internal/sii"
)

// World is a high-level view of a save-game world, built from a
// decoded game.sii and optional external data (country dictionary,
// gameref cache, ...). It is intentionally smaller than the full C#
// object graph but exposes the core entities needed by the CLI.
type World struct {
	// Raw document and typed core blocks
	Document   *sii.Document
	Economy    *items.Economy
	Player     *items.Player
	BlocksByID map[string]sii.Block

	// Aggregated entities from ItemsExtra
	Cities        []*itemsextra.City
	Countries     []*itemsextra.Country
	Companies     []*itemsextra.Company
	Cargoes       []*itemsextra.Cargo
	TrailerDefs   []*itemsextra.TrailerDefinition
	CompanyTrucks []*itemsextra.CompanyTruck
	Garages       []*itemsextra.Garages

	// Player-owned fleet
	PlayerTrucks   map[string]*itemsextra.UserCompanyTruckData
	PlayerTrailers map[string]*itemsextra.UserCompanyTrailerData

	// Optional external data
	CountryDict *externaldata.CountryDictionary
	GameRef     *externaldata.GameRefCache
}

// NewWorld constructs an empty World wrapper for a given document.
func NewWorld(doc *sii.Document) *World {
	return &World{
		Document:       doc,
		BlocksByID:     make(map[string]sii.Block),
		PlayerTrucks:   make(map[string]*itemsextra.UserCompanyTruckData),
		PlayerTrailers: make(map[string]*itemsextra.UserCompanyTrailerData),
	}
}
