package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Garage mirrors the C# Garage class from CustomClasses/Save/Items/Garage.cs.
type Garage struct {
	Vehicles    []string
	Drivers     []string
	Trailers    []string
	Status      int
	ProfitLog   string
	Productivity dataformat.Float
}

// FromProperties populates the Garage from a map of SII properties.
func (g *Garage) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "vehicles":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "vehicles["):
			if val == "null" {
				g.Vehicles = append(g.Vehicles, "")
			} else {
				g.Vehicles = append(g.Vehicles, val)
			}
		case key == "drivers":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "drivers["):
			if val == "null" {
				g.Drivers = append(g.Drivers, "")
			} else {
				g.Drivers = append(g.Drivers, val)
			}
		case key == "trailers":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "trailers["):
			g.Trailers = append(g.Trailers, val)
		case key == "status":
			g.Status = parseInt(val)
		case key == "profit_log":
			g.ProfitLog = val
		case key == "productivity":
			g.Productivity = parseFloat(val)
		}
	}
	return nil
}

// ToProperties converts the Garage struct to a map of properties.
func (g *Garage) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["vehicles"] = []string{strconv.Itoa(len(g.Vehicles))}
	for i, v := range g.Vehicles {
		if v == "" {
			props[fmt.Sprintf("vehicles[%d]", i)] = []string{"null"}
		} else {
			props[fmt.Sprintf("vehicles[%d]", i)] = []string{v}
		}
	}

	props["drivers"] = []string{strconv.Itoa(len(g.Drivers))}
	for i, v := range g.Drivers {
		if v == "" {
			props[fmt.Sprintf("drivers[%d]", i)] = []string{"null"}
		} else {
			props[fmt.Sprintf("drivers[%d]", i)] = []string{v}
		}
	}

	props["trailers"] = []string{strconv.Itoa(len(g.Trailers))}
	for i, v := range g.Trailers {
		props[fmt.Sprintf("trailers[%d]", i)] = []string{v}
	}

	props["status"] = []string{strconv.Itoa(g.Status)}
	props["profit_log"] = []string{g.ProfitLog}
	props["productivity"] = []string{formatFloat(g.Productivity)}

	return props
}

