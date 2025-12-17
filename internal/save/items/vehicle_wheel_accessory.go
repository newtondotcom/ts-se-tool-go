package items

import (
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// VehicleWheelAccessory mirrors the C# Vehicle_Wheel_Accessory class from CustomClasses/Save/Items/Vehicle_Wheel_Accessory.cs.
type VehicleWheelAccessory struct {
	Offset    int
	PaintColor dataformat.Vector3f
	DataPath  string
	Refund    uint32
	AccType   string // derived from data_path, defaults to "generalpart"
}

// FromProperties populates the VehicleWheelAccessory from a map of SII properties.
func (v *VehicleWheelAccessory) FromProperties(props map[string][]string) error {
	v.AccType = "generalpart" // default
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "offset":
			v.Offset = parseInt(val)
		case key == "paint_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.PaintColor = vec
		case key == "data_path":
			v.DataPath = val
			// Determine accType from data_path (mirrors C# logic)
			if strings.Contains(val, "/f_tire/") || strings.Contains(val, "/r_tire/") ||
				strings.Contains(val, "/f_wheel/") || strings.Contains(val, "/r_wheel/") ||
				strings.Contains(val, "/t_wheel/") {
				v.AccType = "tire"
			}
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehicleWheelAccessory struct to a map of properties.
func (v *VehicleWheelAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["offset"] = []string{strconv.Itoa(v.Offset)}
	props["paint_color"] = []string{formatVector3f(v.PaintColor)}
	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

