package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// TransportData mirrors the C# Transport_Data class from CustomClasses/Save/Items/Transport_Data.cs.
type TransportData struct {
	Distance      int
	Time          uint32
	Money         int
	CountPerAdr   []int
	Docks         []dataformat.String
	CountPerDock  []int
}

// FromProperties populates the TransportData from a map of SII properties.
func (t *TransportData) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "distance":
			t.Distance = parseInt(val)
		case key == "time":
			t.Time = parseUint32(val)
		case key == "money":
			t.Money = parseInt(val)
		case key == "count_per_adr":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "count_per_adr["):
			t.CountPerAdr = append(t.CountPerAdr, parseInt(val))
		case key == "docks":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "docks["):
			t.Docks = append(t.Docks, dataformat.String(val))
		case key == "count_per_dock":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "count_per_dock["):
			t.CountPerDock = append(t.CountPerDock, parseInt(val))
		}
	}
	return nil
}

// ToProperties converts the TransportData struct to a map of properties.
func (t *TransportData) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["distance"] = []string{strconv.Itoa(t.Distance)}
	props["time"] = []string{strconv.FormatUint(uint64(t.Time), 10)}
	props["money"] = []string{strconv.Itoa(t.Money)}

	props["count_per_adr"] = []string{strconv.Itoa(len(t.CountPerAdr))}
	for i, v := range t.CountPerAdr {
		props[fmt.Sprintf("count_per_adr[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["docks"] = []string{strconv.Itoa(len(t.Docks))}
	for i, v := range t.Docks {
		props[fmt.Sprintf("docks[%d]", i)] = []string{string(v)}
	}

	props["count_per_dock"] = []string{strconv.Itoa(len(t.CountPerDock))}
	for i, v := range t.CountPerDock {
		props[fmt.Sprintf("count_per_dock[%d]", i)] = []string{strconv.Itoa(v)}
	}

	return props
}

