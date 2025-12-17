package items

import (
	"fmt"
	"strconv"
	"strings"
)

// DeliveryLog mirrors the C# Delivery_log class from CustomClasses/Save/Items/Delivery_log.cs.
type DeliveryLog struct {
	Version        int
	Entries        []string
	CachedJobsCount int
}

// FromProperties populates the DeliveryLog from a map of SII properties.
func (d *DeliveryLog) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "version":
			d.Version = parseInt(val)
		case key == "entries":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "entries["):
			d.Entries = append(d.Entries, val)
		case key == "cached_jobs_count":
			d.CachedJobsCount = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the DeliveryLog struct to a map of properties.
func (d *DeliveryLog) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["version"] = []string{strconv.Itoa(d.Version)}

	props["entries"] = []string{strconv.Itoa(len(d.Entries))}
	for i, v := range d.Entries {
		props[fmt.Sprintf("entries[%d]", i)] = []string{v}
	}

	props["cached_jobs_count"] = []string{strconv.Itoa(d.CachedJobsCount)}

	return props
}

