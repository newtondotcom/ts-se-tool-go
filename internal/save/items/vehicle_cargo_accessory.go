package items

import (
	"strconv"
)

// VehicleCargoAccessory mirrors the C# Vehicle_Cargo_Accessory class from CustomClasses/Save/Items/Vehicle_Cargo_Accessory.cs.
type VehicleCargoAccessory struct {
	CargoData string
	ModelSeed uint32
	DataPath  string
	Refund    uint32
}

// FromProperties populates the VehicleCargoAccessory from a map of SII properties.
func (v *VehicleCargoAccessory) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "cargo_data":
			v.CargoData = val
		case key == "model_seed":
			v.ModelSeed = parseUint32(val)
		case key == "data_path":
			v.DataPath = val
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehicleCargoAccessory struct to a map of properties.
func (v *VehicleCargoAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["cargo_data"] = []string{v.CargoData}
	props["model_seed"] = []string{strconv.FormatUint(uint64(v.ModelSeed), 10)}
	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

