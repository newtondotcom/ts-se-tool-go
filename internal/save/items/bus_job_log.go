package items

import (
	"fmt"
	"strconv"
	"strings"
)

// BusJobLog mirrors the C# Bus_job_Log class from CustomClasses/Save/Items/Bus_job_Log.cs.
type BusJobLog struct {
	Version int
	Entries []string
}

// FromProperties populates the BusJobLog from a map of SII properties.
func (b *BusJobLog) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "version":
			b.Version = parseInt(val)
		case key == "entries":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "entries["):
			b.Entries = append(b.Entries, val)
		}
	}
	return nil
}

// ToProperties converts the BusJobLog struct to a map of properties.
func (b *BusJobLog) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["version"] = []string{strconv.Itoa(b.Version)}

	props["entries"] = []string{strconv.Itoa(len(b.Entries))}
	for i, v := range b.Entries {
		props[fmt.Sprintf("entries[%d]", i)] = []string{v}
	}

	return props
}

