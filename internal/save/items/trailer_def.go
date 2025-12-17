package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// TrailerDef mirrors the C# Trailer_Def class from CustomClasses/Save/Items/Trailer_Def.cs.
type TrailerDef struct {
	Trailer                dataformat.String
	GrossTrailerWeightLimit int
	ChassisMass            dataformat.Float
	BodyMass               dataformat.Float
	Axles                  int
	Volume                 dataformat.Float
	BodyType               string
	ChainType              string
	CountryValidity        []string
	MassRatio              []dataformat.Float
	Length                 dataformat.Float
	SourceName             dataformat.String
}

// FromProperties populates the TrailerDef from a map of SII properties.
func (t *TrailerDef) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "trailer":
			t.Trailer = dataformat.String(val)
		case key == "gross_trailer_weight_limit":
			t.GrossTrailerWeightLimit = parseInt(val)
		case key == "chassis_mass":
			t.ChassisMass = parseFloat(val)
		case key == "body_mass":
			t.BodyMass = parseFloat(val)
		case key == "axles":
			t.Axles = parseInt(val)
		case key == "volume":
			t.Volume = parseFloat(val)
		case key == "body_type":
			t.BodyType = val
		case key == "chain_type":
			t.ChainType = val
		case key == "country_validity":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "country_validity["):
			t.CountryValidity = append(t.CountryValidity, val)
		case key == "mass_ratio":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "mass_ratio["):
			t.MassRatio = append(t.MassRatio, parseFloat(val))
		case key == "length":
			t.Length = parseFloat(val)
		case key == "source_name":
			t.SourceName = dataformat.String(val)
		}
	}
	return nil
}

// ToProperties converts the TrailerDef struct to a map of properties.
func (t *TrailerDef) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["trailer"] = []string{string(t.Trailer)}
	props["gross_trailer_weight_limit"] = []string{strconv.Itoa(t.GrossTrailerWeightLimit)}
	props["chassis_mass"] = []string{formatFloat(t.ChassisMass)}
	props["body_mass"] = []string{formatFloat(t.BodyMass)}
	props["axles"] = []string{strconv.Itoa(t.Axles)}
	props["volume"] = []string{formatFloat(t.Volume)}
	props["body_type"] = []string{t.BodyType}
	props["chain_type"] = []string{t.ChainType}

	props["country_validity"] = []string{strconv.Itoa(len(t.CountryValidity))}
	for i, v := range t.CountryValidity {
		props[fmt.Sprintf("country_validity[%d]", i)] = []string{v}
	}

	props["mass_ratio"] = []string{strconv.Itoa(len(t.MassRatio))}
	for i, v := range t.MassRatio {
		props[fmt.Sprintf("mass_ratio[%d]", i)] = []string{formatFloat(v)}
	}

	props["length"] = []string{formatFloat(t.Length)}
	props["source_name"] = []string{string(t.SourceName)}

	return props
}

