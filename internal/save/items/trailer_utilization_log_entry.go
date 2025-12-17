package items

import (
	"strconv"
)

// TrailerUtilizationLogEntry mirrors the C# Trailer_Utilization_log_Entry class from CustomClasses/Save/Items/Trailer_Utilization_log_Entry.cs.
type TrailerUtilizationLogEntry struct {
	EconomyDay int
	UseTime    int
}

// FromProperties populates the TrailerUtilizationLogEntry from a map of SII properties.
func (t *TrailerUtilizationLogEntry) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "economy_day":
			t.EconomyDay = parseInt(val)
		case key == "use_time":
			t.UseTime = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the TrailerUtilizationLogEntry struct to a map of properties.
func (t *TrailerUtilizationLogEntry) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["economy_day"] = []string{strconv.Itoa(t.EconomyDay)}
	props["use_time"] = []string{strconv.Itoa(t.UseTime)}

	return props
}

