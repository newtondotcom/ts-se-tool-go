package items

import (
	"strconv"
)

// FerryLogEntry mirrors the C# Ferry_log_Entry class from CustomClasses/Save/Items/Ferry_log_Entry.cs.
type FerryLogEntry struct {
	Ferry      string
	Connection string
	LastVisit  uint32
	UseCount   uint32
}

// FromProperties populates the FerryLogEntry from a map of SII properties.
func (f *FerryLogEntry) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "ferry":
			f.Ferry = val
		case key == "connection":
			f.Connection = val
		case key == "last_visit":
			f.LastVisit = parseUint32(val)
		case key == "use_count":
			f.UseCount = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the FerryLogEntry struct to a map of properties.
func (f *FerryLogEntry) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["ferry"] = []string{f.Ferry}
	props["connection"] = []string{f.Connection}
	props["last_visit"] = []string{strconv.FormatUint(uint64(f.LastVisit), 10)}
	props["use_count"] = []string{strconv.FormatUint(uint64(f.UseCount), 10)}

	return props
}

