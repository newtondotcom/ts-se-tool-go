package items

import (
	"fmt"
	"strconv"
	"strings"
)

// Registry mirrors the C# Registry class from CustomClasses/Save/Items/Registry.cs.
type Registry struct {
	Data  []int
	Valid []bool
	Keys  []int
	Index []int
}

// FromProperties populates the Registry from a map of SII properties.
func (r *Registry) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "data":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "data["):
			r.Data = append(r.Data, parseInt(val))
		case key == "valid":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "valid["):
			r.Valid = append(r.Valid, parseBool(val))
		case key == "keys":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "keys["):
			r.Keys = append(r.Keys, parseInt(val))
		case key == "index":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "index["):
			r.Index = append(r.Index, parseInt(val))
		}
	}
	return nil
}

// ToProperties converts the Registry struct to a map of properties.
func (r *Registry) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["data"] = []string{strconv.Itoa(len(r.Data))}
	for i, v := range r.Data {
		props[fmt.Sprintf("data[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["valid"] = []string{strconv.Itoa(len(r.Valid))}
	for i, v := range r.Valid {
		props[fmt.Sprintf("valid[%d]", i)] = []string{formatBool(v)}
	}

	props["keys"] = []string{strconv.Itoa(len(r.Keys))}
	for i, v := range r.Keys {
		props[fmt.Sprintf("keys[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["index"] = []string{strconv.Itoa(len(r.Index))}
	for i, v := range r.Index {
		props[fmt.Sprintf("index[%d]", i)] = []string{strconv.Itoa(v)}
	}

	return props
}

