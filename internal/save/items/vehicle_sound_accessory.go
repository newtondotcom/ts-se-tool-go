package items

import (
	"strconv"
)

// VehicleSoundAccessory mirrors the C# Vehicle_Sound_Accessory class from CustomClasses/Save/Items/Vehicle_Sound_Accessory.cs.
type VehicleSoundAccessory struct {
	DataPath string
	Refund   uint32
}

// FromProperties populates the VehicleSoundAccessory from a map of SII properties.
func (v *VehicleSoundAccessory) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "data_path":
			v.DataPath = val
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehicleSoundAccessory struct to a map of properties.
func (v *VehicleSoundAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

