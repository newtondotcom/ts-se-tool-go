package items

import (
	"strconv"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// VehiclePaintJobAccessory mirrors the C# Vehicle_Paint_job_Accessory class from CustomClasses/Save/Items/Vehicle_Paint_job_Accessory.cs.
type VehiclePaintJobAccessory struct {
	MaskRColor dataformat.Vector3f
	MaskGColor dataformat.Vector3f
	MaskBColor dataformat.Vector3f
	FlakeColor dataformat.Vector3f
	FlipColor  dataformat.Vector3f
	BaseColor  dataformat.Vector3f
	DataPath   string
	Refund     uint32
}

// FromProperties populates the VehiclePaintJobAccessory from a map of SII properties.
func (v *VehiclePaintJobAccessory) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "mask_r_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.MaskRColor = vec
		case key == "mask_g_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.MaskGColor = vec
		case key == "mask_b_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.MaskBColor = vec
		case key == "flake_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.FlakeColor = vec
		case key == "flip_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.FlipColor = vec
		case key == "base_color":
			vec, err := parseVector3f(val)
			if err != nil {
				return err
			}
			v.BaseColor = vec
		case key == "data_path":
			v.DataPath = val
		case key == "refund":
			v.Refund = parseUint32(val)
		}
	}
	return nil
}

// ToProperties converts the VehiclePaintJobAccessory struct to a map of properties.
func (v *VehiclePaintJobAccessory) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["mask_r_color"] = []string{formatVector3f(v.MaskRColor)}
	props["mask_g_color"] = []string{formatVector3f(v.MaskGColor)}
	props["mask_b_color"] = []string{formatVector3f(v.MaskBColor)}
	props["flake_color"] = []string{formatVector3f(v.FlakeColor)}
	props["flip_color"] = []string{formatVector3f(v.FlipColor)}
	props["base_color"] = []string{formatVector3f(v.BaseColor)}
	props["data_path"] = []string{v.DataPath}
	props["refund"] = []string{strconv.FormatUint(uint64(v.Refund), 10)}

	return props
}

