package items

import (
	"fmt"
	"strconv"
	"strings"
)

// FerryLog mirrors the C# Ferry_log class from CustomClasses/Save/Items/Ferry_log.cs.
type FerryLog struct {
	Entries []string
}

// FromProperties populates the FerryLog from a map of SII properties.
func (f *FerryLog) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "entries":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "entries["):
			f.Entries = append(f.Entries, val)
		}
	}
	return nil
}

// ToProperties converts the FerryLog struct to a map of properties.
func (f *FerryLog) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["entries"] = []string{strconv.Itoa(len(f.Entries))}
	for i, v := range f.Entries {
		props[fmt.Sprintf("entries[%d]", i)] = []string{v}
	}

	return props
}

