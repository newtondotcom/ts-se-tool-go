package items

import (
	"strconv"
)

// EconomyEvent mirrors the C# Economy_event class from CustomClasses/Save/Items/Economy_event.cs.
type EconomyEvent struct {
	Time      uint32
	UnitLink  string
	Param     int
}

// FromProperties populates the EconomyEvent from a map of SII properties.
func (e *EconomyEvent) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "time":
			e.Time = parseUint32(val)
		case key == "unit_link":
			e.UnitLink = val
		case key == "param":
			e.Param = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the EconomyEvent struct to a map of properties.
func (e *EconomyEvent) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["time"] = []string{strconv.FormatUint(uint64(e.Time), 10)}
	props["unit_link"] = []string{e.UnitLink}
	props["param"] = []string{strconv.Itoa(e.Param)}

	return props
}

