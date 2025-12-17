package items

import (
	"fmt"
	"strconv"
	"strings"
)

// EconomyEventQueue mirrors the C# Economy_event_Queue class from CustomClasses/Save/Items/Economy_event_Queue.cs.
type EconomyEventQueue struct {
	Data []string
}

// FromProperties populates the EconomyEventQueue from a map of SII properties.
func (e *EconomyEventQueue) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "data":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "data["):
			e.Data = append(e.Data, val)
		}
	}
	return nil
}

// ToProperties converts the EconomyEventQueue struct to a map of properties.
func (e *EconomyEventQueue) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["data"] = []string{strconv.Itoa(len(e.Data))}
	for i, v := range e.Data {
		props[fmt.Sprintf("data[%d]", i)] = []string{v}
	}

	return props
}

