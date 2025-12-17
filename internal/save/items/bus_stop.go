package items

import (
	"fmt"
	"strconv"
	"strings"
)

// BusStop mirrors the C# Bus_stop class from CustomClasses/Save/Items/Bus_stop.cs.
type BusStop struct {
	Discovered bool
	LinesOffer []string
}

// FromProperties populates the BusStop from a map of SII properties.
func (b *BusStop) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "discovered":
			b.Discovered = parseBool(val)
		case key == "lines_offer":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "lines_offer["):
			b.LinesOffer = append(b.LinesOffer, val)
		}
	}
	return nil
}

// ToProperties converts the BusStop struct to a map of properties.
func (b *BusStop) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["discovered"] = []string{formatBool(b.Discovered)}

	props["lines_offer"] = []string{strconv.Itoa(len(b.LinesOffer))}
	for i, v := range b.LinesOffer {
		props[fmt.Sprintf("lines_offer[%d]", i)] = []string{v}
	}

	return props
}

