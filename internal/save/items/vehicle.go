package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Vehicle mirrors the C# Vehicle class from CustomClasses/Save/Items/Vehicle.cs.
type Vehicle struct {
	EngineWear                dataformat.Float
	TransmissionWear          dataformat.Float
	CabinWear                 dataformat.Float
	ChassisWear               dataformat.Float
	WheelsWear                []dataformat.Float
	Accessories               []string
	LicensePlate              dataformat.String
	FuelRelative              dataformat.Float
	Odometer                  uint32
	OdometerFloatPart         dataformat.Float
	RheostatFactor            dataformat.Float
	UserMirrorRot             []dataformat.Vector4f
	UserHeadOffset            dataformat.Vector3f
	UserFov                   dataformat.Float
	UserWheelUpDown           dataformat.Float
	UserWheelFrontBack        dataformat.Float
	UserMouseLeftRightDefault dataformat.Float
	UserMouseUpDownDefault    dataformat.Float
	TripFuelL                 uint32
	TripFuel                  dataformat.Float
	TripDistanceKm            uint32
	TripDistance              dataformat.Float
	TripTimeMin               uint32
	TripTime                  dataformat.Float
}

// FromProperties populates the Vehicle from a map of SII properties.
func (v *Vehicle) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "engine_wear":
			v.EngineWear = parseFloat(val)
		case key == "transmission_wear":
			v.TransmissionWear = parseFloat(val)
		case key == "cabin_wear":
			v.CabinWear = parseFloat(val)
		case key == "chassis_wear":
			v.ChassisWear = parseFloat(val)
		case key == "wheels_wear":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "wheels_wear["):
			v.WheelsWear = append(v.WheelsWear, parseFloat(val))
		case key == "fuel_relative":
			v.FuelRelative = parseFloat(val)
		case key == "license_plate":
			v.LicensePlate = dataformat.String(val)
		case key == "accessories":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "accessories["):
			v.Accessories = append(v.Accessories, val)
		case key == "user_mirror_rot":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "user_mirror_rot["):
			vec, err := parseVector4f(val)
			if err != nil {
				return fmt.Errorf("parse user_mirror_rot: %w", err)
			}
			v.UserMirrorRot = append(v.UserMirrorRot, vec)
		case key == "user_head_offset":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse user_head_offset: %w", err)
			}
			v.UserHeadOffset = vec
		case key == "user_fov":
			v.UserFov = parseFloat(val)
		case key == "user_wheel_up_down":
			v.UserWheelUpDown = parseFloat(val)
		case key == "user_wheel_front_back":
			v.UserWheelFrontBack = parseFloat(val)
		case key == "user_mouse_left_right_default":
			v.UserMouseLeftRightDefault = parseFloat(val)
		case key == "user_mouse_up_down_default":
			v.UserMouseUpDownDefault = parseFloat(val)
		case key == "rheostat_factor":
			v.RheostatFactor = parseFloat(val)
		case key == "odometer":
			v.Odometer = parseUint32(val)
		case key == "odometer_float_part":
			v.OdometerFloatPart = parseFloat(val)
		case key == "trip_fuel_l":
			v.TripFuelL = parseUint32(val)
		case key == "trip_fuel":
			v.TripFuel = parseFloat(val)
		case key == "trip_distance_km":
			v.TripDistanceKm = parseUint32(val)
		case key == "trip_distance":
			v.TripDistance = parseFloat(val)
		case key == "trip_time_min":
			v.TripTimeMin = parseUint32(val)
		case key == "trip_time":
			v.TripTime = parseFloat(val)
		}
	}
	return nil
}

// ToProperties converts the Vehicle struct to a map of properties.
func (v *Vehicle) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["engine_wear"] = []string{formatFloat(v.EngineWear)}
	props["transmission_wear"] = []string{formatFloat(v.TransmissionWear)}
	props["cabin_wear"] = []string{formatFloat(v.CabinWear)}
	props["fuel_relative"] = []string{formatFloat(v.FuelRelative)}
	props["rheostat_factor"] = []string{formatFloat(v.RheostatFactor)}

	props["user_mirror_rot"] = []string{strconv.Itoa(len(v.UserMirrorRot))}
	for i, vec := range v.UserMirrorRot {
		props[fmt.Sprintf("user_mirror_rot[%d]", i)] = []string{formatVector4f(vec)}
	}

	props["user_head_offset"] = []string{formatVector3f(v.UserHeadOffset)}
	props["user_fov"] = []string{formatFloat(v.UserFov)}
	props["user_wheel_up_down"] = []string{formatFloat(v.UserWheelUpDown)}
	props["user_wheel_front_back"] = []string{formatFloat(v.UserWheelFrontBack)}
	props["user_mouse_left_right_default"] = []string{formatFloat(v.UserMouseLeftRightDefault)}
	props["user_mouse_up_down_default"] = []string{formatFloat(v.UserMouseUpDownDefault)}

	props["accessories"] = []string{strconv.Itoa(len(v.Accessories))}
	for i, acc := range v.Accessories {
		props[fmt.Sprintf("accessories[%d]", i)] = []string{acc}
	}

	props["odometer"] = []string{strconv.FormatUint(uint64(v.Odometer), 10)}
	props["odometer_float_part"] = []string{formatFloat(v.OdometerFloatPart)}
	props["trip_fuel_l"] = []string{strconv.FormatUint(uint64(v.TripFuelL), 10)}
	props["trip_fuel"] = []string{formatFloat(v.TripFuel)}
	props["trip_distance_km"] = []string{strconv.FormatUint(uint64(v.TripDistanceKm), 10)}
	props["trip_distance"] = []string{formatFloat(v.TripDistance)}
	props["trip_time_min"] = []string{strconv.FormatUint(uint64(v.TripTimeMin), 10)}
	props["trip_time"] = []string{formatFloat(v.TripTime)}
	props["license_plate"] = []string{string(v.LicensePlate)}
	props["chassis_wear"] = []string{formatFloat(v.ChassisWear)}

	props["wheels_wear"] = []string{strconv.Itoa(len(v.WheelsWear))}
	for i, wear := range v.WheelsWear {
		props[fmt.Sprintf("wheels_wear[%d]", i)] = []string{formatFloat(wear)}
	}

	return props
}

