package items

import (
	"strconv"
	"strings"
)

// VehicleAccessory mirrors the C# Vehicle_Accessory class from CustomClasses/Save/Items/Vehicle_Accessory.cs.
type VehicleAccessory struct {
	DataPath string
	Refund   uint32
	AccType  string // derived from data_path, defaults to "generalpart"
}

// FromProperties populates the VehicleAccessory from a map of SII properties.
func (v *VehicleAccessory) FromProperties(props map[string][]string) error {
	v.AccType = "generalpart" // default
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "data_path":
			v.DataPath = val
			// Determine accType from data_path (mirrors C# logic)
			if strings.Contains(val, "/data.sii") {
				v.AccType = "basepart"
			} else if strings.Contains(val, "chassis") {
				v.AccType = "chassis"
			} else if strings.Contains(val, "body") {
				v.AccType = "body"
			} else if strings.Contains(val, "cabin") {
				v.AccType = "cabin"
			} else if strings.Contains(val, "engine") {
				v.AccType = "engine"
			} else if strings.Contains(val, "transmission") {
				v.AccType = "transmission"
			}
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehicleAccessory struct to a map of properties.
func (v *VehicleAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

