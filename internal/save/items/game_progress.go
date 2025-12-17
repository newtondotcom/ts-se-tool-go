package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// GameProgress mirrors the C# Game_Progress class from CustomClasses/Save/Items/Game_Progress.cs.
type GameProgress struct {
	GenericTransports  string
	UndamagedTransports string
	CleanTransports    string
	OwnedTrucks        []dataformat.String
}

// FromProperties populates the GameProgress from a map of SII properties.
func (g *GameProgress) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "generic_transports":
			g.GenericTransports = val
		case key == "undamaged_transports":
			g.UndamagedTransports = val
		case key == "clean_transports":
			g.CleanTransports = val
		case key == "owned_trucks":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "owned_trucks["):
			g.OwnedTrucks = append(g.OwnedTrucks, dataformat.String(val))
		}
	}
	return nil
}

// ToProperties converts the GameProgress struct to a map of properties.
func (g *GameProgress) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["generic_transports"] = []string{g.GenericTransports}
	props["undamaged_transports"] = []string{g.UndamagedTransports}
	props["clean_transports"] = []string{g.CleanTransports}

	props["owned_trucks"] = []string{strconv.Itoa(len(g.OwnedTrucks))}
	for i, v := range g.OwnedTrucks {
		props[fmt.Sprintf("owned_trucks[%d]", i)] = []string{string(v)}
	}

	return props
}

