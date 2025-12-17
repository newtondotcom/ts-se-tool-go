package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Trailer mirrors the C# Trailer class from CustomClasses/Save/Items/Trailer.cs.
type Trailer struct {
	CargoMass                dataformat.Float
	CargoDamage              dataformat.Float
	TrailerBodyWear          dataformat.Float
	ChassisWear              dataformat.Float
	WheelsWear               []dataformat.Float
	Accessories              []string
	LicensePlate             dataformat.String
	TrailerDefinition        string
	SlaveTrailer             string
	IsPrivate                bool
	Oversize                 bool
	VirtualRearWheelsOffset  dataformat.Float
	Odometer                 uint32
	OdometerFloatPart        dataformat.Float
	TripFuelL                uint32
	TripFuel                 dataformat.Float
	TripDistanceKm           uint32
	TripDistance             dataformat.Float
	TripTimeMin              uint32
	TripTime                 dataformat.Float
}

// FromProperties populates the Trailer from a map of SII properties.
func (t *Trailer) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "trailer_definition":
			t.TrailerDefinition = val
		case key == "oversize":
			t.Oversize = parseBool(val)
		case key == "cargo_mass":
			t.CargoMass = parseFloat(val)
		case key == "cargo_damage":
			t.CargoDamage = parseFloat(val)
		case key == "virtual_rear_wheels_offset":
			t.VirtualRearWheelsOffset = parseFloat(val)
		case key == "slave_trailer":
			t.SlaveTrailer = val
		case key == "is_private":
			t.IsPrivate = parseBool(val)
		case key == "trailer_body_wear":
			t.TrailerBodyWear = parseFloat(val)
		case key == "accessories":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "accessories["):
			t.Accessories = append(t.Accessories, val)
		case key == "odometer":
			t.Odometer = parseUint32(val)
		case key == "odometer_float_part":
			t.OdometerFloatPart = parseFloat(val)
		case key == "trip_fuel_l":
			t.TripFuelL = parseUint32(val)
		case key == "trip_fuel":
			t.TripFuel = parseFloat(val)
		case key == "trip_distance_km":
			t.TripDistanceKm = parseUint32(val)
		case key == "trip_distance":
			t.TripDistance = parseFloat(val)
		case key == "trip_time_min":
			t.TripTimeMin = parseUint32(val)
		case key == "trip_time":
			t.TripTime = parseFloat(val)
		case key == "license_plate":
			t.LicensePlate = dataformat.String(val)
		case key == "chassis_wear":
			t.ChassisWear = parseFloat(val)
		case key == "wheels_wear":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "wheels_wear["):
			t.WheelsWear = append(t.WheelsWear, parseFloat(val))
		}
	}
	return nil
}

// ToProperties converts the Trailer struct to a map of properties.
func (t *Trailer) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["trailer_definition"] = []string{t.TrailerDefinition}
	props["oversize"] = []string{formatBool(t.Oversize)}
	props["cargo_mass"] = []string{formatFloat(t.CargoMass)}
	props["cargo_damage"] = []string{formatFloat(t.CargoDamage)}
	props["virtual_rear_wheels_offset"] = []string{formatFloat(t.VirtualRearWheelsOffset)}
	props["slave_trailer"] = []string{t.SlaveTrailer}
	props["is_private"] = []string{formatBool(t.IsPrivate)}
	props["trailer_body_wear"] = []string{formatFloat(t.TrailerBodyWear)}

	props["accessories"] = []string{strconv.Itoa(len(t.Accessories))}
	for i, v := range t.Accessories {
		props[fmt.Sprintf("accessories[%d]", i)] = []string{v}
	}

	props["odometer"] = []string{strconv.FormatUint(uint64(t.Odometer), 10)}
	props["odometer_float_part"] = []string{formatFloat(t.OdometerFloatPart)}
	props["trip_fuel_l"] = []string{strconv.FormatUint(uint64(t.TripFuelL), 10)}
	props["trip_fuel"] = []string{formatFloat(t.TripFuel)}
	props["trip_distance_km"] = []string{strconv.FormatUint(uint64(t.TripDistanceKm), 10)}
	props["trip_distance"] = []string{formatFloat(t.TripDistance)}
	props["trip_time_min"] = []string{strconv.FormatUint(uint64(t.TripTimeMin), 10)}
	props["trip_time"] = []string{formatFloat(t.TripTime)}
	props["license_plate"] = []string{string(t.LicensePlate)}
	props["chassis_wear"] = []string{formatFloat(t.ChassisWear)}

	props["wheels_wear"] = []string{strconv.Itoa(len(t.WheelsWear))}
	for i, v := range t.WheelsWear {
		props[fmt.Sprintf("wheels_wear[%d]", i)] = []string{formatFloat(v)}
	}

	return props
}

