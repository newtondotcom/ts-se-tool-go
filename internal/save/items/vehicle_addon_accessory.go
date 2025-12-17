package items

import (
	"fmt"
	"strconv"
	"strings"
)

// VehicleAddonAccessory mirrors the C# Vehicle_Addon_Accessory class from CustomClasses/Save/Items/Vehicle_Addon_Accessory.cs.
type VehicleAddonAccessory struct {
	SlotName  []string
	SlotHookup []string
	DataPath  string
	Refund    uint32
}

// FromProperties populates the VehicleAddonAccessory from a map of SII properties.
func (v *VehicleAddonAccessory) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "slot_name":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "slot_name["):
			v.SlotName = append(v.SlotName, val)
		case key == "slot_hookup":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "slot_hookup["):
			v.SlotHookup = append(v.SlotHookup, val)
		case key == "data_path":
			v.DataPath = val
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehicleAddonAccessory struct to a map of properties.
func (v *VehicleAddonAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["slot_name"] = []string{strconv.Itoa(len(v.SlotName))}
	for i, s := range v.SlotName {
		props[fmt.Sprintf("slot_name[%d]", i)] = []string{s}
	}

	props["slot_hookup"] = []string{strconv.Itoa(len(v.SlotHookup))}
	for i, s := range v.SlotHookup {
		props[fmt.Sprintf("slot_hookup[%d]", i)] = []string{s}
	}

	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

