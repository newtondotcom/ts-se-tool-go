package externaldata

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// This package mirrors classes from CustomClasses/ExternalData.
// We only implement the parts that are directly useful for the
// world-loading pipeline (CountryDictionary, GameRefCache, heavy
// cargo lists, localisation files, ...).

// -------------------------------------------------------------------
// Countries (CityToCountry + CountryProperties)
// -------------------------------------------------------------------

// CountryDictionary corresponds to CountryDictionary.cs. It provides a
// mapping city -> country loaded from a simple CSV/TSV style file.
type CountryDictionary struct {
	byCity map[string]string
}

// LoadCountryDictionary loads a City->Country mapping from a text file.
// The expected format is one mapping per line using either ';' or ',' as a
// separator, e.g.:
//
//	bruxelles;belgium
//	lyon;france
func LoadCountryDictionary(path string) (*CountryDictionary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cd := &CountryDictionary{byCity: make(map[string]string)}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		sep := ";"
		if strings.Contains(line, ",") {
			sep = ","
		}
		parts := strings.SplitN(line, sep, 2)
		if len(parts) != 2 {
			continue
		}
		city := strings.TrimSpace(parts[0])
		country := strings.TrimSpace(parts[1])
		if city == "" {
			continue
		}
		cd.byCity[city] = country
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cd, nil
}

// GetCountry returns the country identifier for the given city or
// an empty string if unknown.
func (cd *CountryDictionary) GetCountry(city string) string {
	if cd == nil {
		return ""
	}
	return cd.byCity[city]
}

// CountryProperties mirrors the data loaded from CountryProperties.csv in
// TS SE Tool: a small lookup table of extra attributes per country.
type CountryProperties struct {
	ID          string
	DisplayName string
	Value       string
}

// LoadCountryProperties loads CountryProperties.csv into a map keyed by ID.
// The expected format per line is:
//
//	id;displayName;numericOrTextValue
func LoadCountryProperties(path string) (map[string]CountryProperties, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := make(map[string]CountryProperties)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ";")
		if len(parts) < 3 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])
		val := strings.TrimSpace(parts[2])
		if id == "" {
			continue
		}
		out[id] = CountryProperties{ID: id, DisplayName: name, Value: val}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// LoadExtCountries is a convenience helper mirroring the C# LoadExtCountries
// method: it loads both the City->Country mapping and the CountryProperties
// table in one call.
func LoadExtCountries(cityToCountryPath, countryPropsPath string) (*CountryDictionary, map[string]CountryProperties, error) {
	cd, err := LoadCountryDictionary(cityToCountryPath)
	if err != nil {
		return nil, nil, err
	}
	props, err := LoadCountryProperties(countryPropsPath)
	if err != nil {
		return cd, nil, err
	}
	return cd, props, nil
}

// -------------------------------------------------------------------
// Heavy cargo list (heavy_cargoes.csv)
// -------------------------------------------------------------------

var defaultHeavyCargos = []string{
	"asph_miller", "cable_reel", "concr_beams", "dozer", "locomotive", "metal_center", "mobile_crane", "transformat",
	"case600", "cat627", "coil", "kalmar240", "kalmar240_s", "komatsu155", "terex3160", "transformer", "wirtgen250",
}

// LoadHeavyCargoes implements the behaviour of LoadExtCargoes in C#:
// - tries to read heavy_cargoes.csv (one cargo id per line)
// - if missing, creates the file with a default list and returns that list.
func LoadHeavyCargoes(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Create default file and return defaults
		if mkErr := os.MkdirAll(filepath.Dir(path), 0o755); mkErr != nil {
			return defaultHeavyCargos, mkErr
		}
		if wrErr := os.WriteFile(path, []byte(strings.Join(defaultHeavyCargos, "\n")), 0o644); wrErr != nil {
			return defaultHeavyCargos, wrErr
		}
		return append([]string(nil), defaultHeavyCargos...), nil
	}

	lines := strings.Split(string(data), "\n")
	var out []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		out = append(out, l)
	}
	return out, nil
}

// -------------------------------------------------------------------
// Localisation helpers (truck brands, driver names)
// -------------------------------------------------------------------

// LoadTruckBrands mirrors LoadTruckBrandsLng: it reads
// lang/Default/truck_brands.txt and returns a map brandID -> display name.
func LoadTruckBrands(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	out := make(map[string]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.SplitN(l, ";", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		out[key] = val
	}
	return out, nil
}

// LoadDriverNames mirrors LoadDriverNamesLng: it reads
// lang/Default/<GameType>/driver_names.csv and returns a map id -> name.
func LoadDriverNames(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	out := make(map[string]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		parts := strings.SplitN(l, ";", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		out[key] = val
	}
	return out, nil
}

// -------------------------------------------------------------------
// External cargo/company definitions (gameref cache)
// -------------------------------------------------------------------

// ExtCargo corresponds to ExtCargo.cs. Only a subset of fields is
// modelled for now; we can extend it as needed.
type ExtCargo struct {
	ID         string
	Fragility  int
	ADRClass   int
	Mass       int
	UnitReward int
}

// ExtCompany corresponds to ExtCompany.cs.
type ExtCompany struct {
	ID        string
	CargosIn  []string
	CargosOut []string
}

// GameRefCache groups parsed external game data (cargo and company
// definitions). This acts as a replacement for the original .sdf cache.
type GameRefCache struct {
	CargoesByID   map[string]*ExtCargo
	CompaniesByID map[string]*ExtCompany
}

// -------------------------------------------------------------------
// Other shells kept for future parity with TS SE Tool
// -------------------------------------------------------------------

// LevelNames corresponds to LevelNames.cs.
type LevelNames struct{}

// Routes corresponds to Routes.cs.
type Routes struct{}

// ScsFont corresponds to ScsFont.cs.
type ScsFont struct{}

// ErrNotFound is a generic sentinel error for optional external data.
var ErrNotFound = errors.New("external data not found")
