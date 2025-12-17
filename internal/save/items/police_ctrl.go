package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// PoliceCtrl mirrors the C# Police_Ctrl class from CustomClasses/Save/Items/Police_Ctrl.cs.
type PoliceCtrl struct {
	OffenceTimer   []dataformat.Float
	OffenceCounter []int
	OffenceValid   []bool
}

// FromProperties populates the PoliceCtrl from a map of SII properties.
func (p *PoliceCtrl) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "offence_timer":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "offence_timer["):
			p.OffenceTimer = append(p.OffenceTimer, parseFloat(val))
		case key == "offence_counter":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "offence_counter["):
			p.OffenceCounter = append(p.OffenceCounter, parseInt(val))
		case key == "offence_valid":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "offence_valid["):
			p.OffenceValid = append(p.OffenceValid, parseBool(val))
		}
	}
	return nil
}

// ToProperties converts the PoliceCtrl struct to a map of properties.
func (p *PoliceCtrl) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["offence_timer"] = []string{strconv.Itoa(len(p.OffenceTimer))}
	for i, v := range p.OffenceTimer {
		props[fmt.Sprintf("offence_timer[%d]", i)] = []string{formatFloat(v)}
	}

	props["offence_counter"] = []string{strconv.Itoa(len(p.OffenceCounter))}
	for i, v := range p.OffenceCounter {
		props[fmt.Sprintf("offence_counter[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["offence_valid"] = []string{strconv.Itoa(len(p.OffenceValid))}
	for i, v := range p.OffenceValid {
		props[fmt.Sprintf("offence_valid[%d]", i)] = []string{formatBool(v)}
	}

	return props
}

