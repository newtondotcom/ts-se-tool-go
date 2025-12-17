package items

import (
	"fmt"
	"strconv"
	"strings"
)

// DeliveryLogEntry mirrors the C# Delivery_log_Entry class from CustomClasses/Save/Items/Delivery_log_Entry.cs.
type DeliveryLogEntry struct {
	Params []string
}

// FromProperties populates the DeliveryLogEntry from a map of SII properties.
func (d *DeliveryLogEntry) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "params":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "params["):
			d.Params = append(d.Params, val)
		}
	}
	return nil
}

// ToProperties converts the DeliveryLogEntry struct to a map of properties.
func (d *DeliveryLogEntry) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["params"] = []string{strconv.Itoa(len(d.Params))}
	for i, v := range d.Params {
		props[fmt.Sprintf("params[%d]", i)] = []string{v}
	}

	return props
}

