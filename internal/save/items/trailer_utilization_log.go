package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// TrailerUtilizationLog mirrors the C# Trailer_Utilization_log class from CustomClasses/Save/Items/Trailer_Utilization_log.cs.
type TrailerUtilizationLog struct {
	Entries                []string
	TotalDrivenDistanceKm  int
	TotalTransportedCargoes int
	TotalTransportedWeight  dataformat.Float
}

// FromProperties populates the TrailerUtilizationLog from a map of SII properties.
func (t *TrailerUtilizationLog) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "entries":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "entries["):
			t.Entries = append(t.Entries, val)
		case key == "total_driven_distance_km":
			t.TotalDrivenDistanceKm = parseInt(val)
		case key == "total_transported_cargoes":
			t.TotalTransportedCargoes = parseInt(val)
		case key == "total_transported_weight":
			t.TotalTransportedWeight = parseFloat(val)
		}
	}
	return nil
}

// ToProperties converts the TrailerUtilizationLog struct to a map of properties.
func (t *TrailerUtilizationLog) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["entries"] = []string{strconv.Itoa(len(t.Entries))}
	for i, v := range t.Entries {
		props[fmt.Sprintf("entries[%d]", i)] = []string{v}
	}

	props["total_driven_distance_km"] = []string{strconv.Itoa(t.TotalDrivenDistanceKm)}
	props["total_transported_cargoes"] = []string{strconv.Itoa(t.TotalTransportedCargoes)}
	props["total_transported_weight"] = []string{formatFloat(t.TotalTransportedWeight)}

	return props
}

