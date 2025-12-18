package gameref

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robebs/ts-se-tool-go/internal/externaldata"
	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// BuildGameRefCache scans a gameref directory and builds a cache of
// external cargo/company data, roughly similar in spirit to the .sdf
// cache that TS SE Tool uses. The layout is expected to be:
//
//	root/<GameType>/<DLC>/def/...
func BuildGameRefCache(root, gameType string) (*externaldata.GameRefCache, error) {
	cache := &externaldata.GameRefCache{
		CargoesByID:   make(map[string]*externaldata.ExtCargo),
		CompaniesByID: make(map[string]*externaldata.ExtCompany),
	}

	gameRoot := filepath.Join(root, gameType)
	entries, err := os.ReadDir(gameRoot)
	if err != nil {
		return nil, err
	}

	for _, dlc := range entries {
		if !dlc.IsDir() {
			continue
		}
		dlcRoot := filepath.Join(gameRoot, dlc.Name(), "def")
		if err := scanCargoDefs(cache, dlcRoot); err != nil {
			return nil, err
		}
		if err := scanCompanyDefs(cache, dlcRoot); err != nil {
			return nil, err
		}
	}

	return cache, nil
}

func scanCargoDefs(cache *externaldata.GameRefCache, defRoot string) error {
	cargoRoot := filepath.Join(defRoot, "cargo")
	files, err := filepath.Glob(filepath.Join(cargoRoot, "*.sii"))
	if err != nil {
		return err
	}
	for _, f := range files {
		plain, err := siidecrypt.DecryptFile(f, true)
		if err != nil {
			return fmt.Errorf("decrypt cargo def %s: %w", f, err)
		}
		doc, err := sii.ReadDocument(plain)
		if err != nil {
			return fmt.Errorf("parse cargo def %s: %w", f, err)
		}
		for _, b := range doc.Blocks {
			if b.Type != "cargo_data" {
				continue
			}
			id := b.Name
			if _, exists := cache.CargoesByID[id]; exists {
				continue
			}
			cache.CargoesByID[id] = &externaldata.ExtCargo{ID: id}
		}
	}
	return nil
}

func scanCompanyDefs(cache *externaldata.GameRefCache, defRoot string) error {
	companyRoot := filepath.Join(defRoot, "company")
	entries, err := os.ReadDir(companyRoot)
	if err != nil {
		return nil
	}
	for _, dir := range entries {
		if !dir.IsDir() {
			continue
		}
		companyID := dir.Name()
		c := cache.CompaniesByID[companyID]
		if c == nil {
			c = &externaldata.ExtCompany{ID: companyID}
			cache.CompaniesByID[companyID] = c
		}
		// We do not parse individual in/out files in detail here; it is enough
		// for now to record the presence of the company in the cache.
	}
	return nil
}
