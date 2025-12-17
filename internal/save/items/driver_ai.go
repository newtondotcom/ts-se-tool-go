package items

import (
	"strconv"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// DriverAI mirrors the C# Driver_AI class from CustomClasses/Save/Items/Driver_AI.cs.
type DriverAI struct {
	Adr                      uint8
	LongDist                 uint8
	Heavy                    uint8
	Fragile                  uint8
	Urgent                   uint8
	Mechanical               uint8
	Hometown                 string
	CurrentCity              string
	State                    uint32
	OnDutyTimer              int
	ExtraMaintenance         int
	DriverJob                string
	ExperiencePoints         int
	TrainingPolicy           int
	AdoptedTruck             string
	AssignedTruck            string
	AssignedTruckEfficiency  dataformat.Float
	AssignedTruckAxleCount   int
	AssignedTruckMass        dataformat.Float
	SlotTruckEfficiency      dataformat.Float
	SlotTruckAxleCount       int
	SlotTruckMass            dataformat.Float
	AdoptedTrailer           string
	AssignedTrailer          string
	OldHometown              string
	ProfitLog                string
}

// FromProperties populates the DriverAI from a map of SII properties.
func (d *DriverAI) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "adr":
			d.Adr = parseUint8(val)
		case key == "long_dist":
			d.LongDist = parseUint8(val)
		case key == "heavy":
			d.Heavy = parseUint8(val)
		case key == "fragile":
			d.Fragile = parseUint8(val)
		case key == "urgent":
			d.Urgent = parseUint8(val)
		case key == "mechanical":
			d.Mechanical = parseUint8(val)
		case key == "hometown":
			d.Hometown = val
		case key == "current_city":
			d.CurrentCity = val
		case key == "state":
			d.State = parseUint32(val)
		case key == "on_duty_timer":
			d.OnDutyTimer = parseInt(val)
		case key == "extra_maintenance":
			d.ExtraMaintenance = parseInt(val)
		case key == "driver_job":
			d.DriverJob = val
		case key == "experience_points":
			d.ExperiencePoints = parseInt(val)
		case key == "training_policy":
			d.TrainingPolicy = parseInt(val)
		case key == "adopted_truck":
			d.AdoptedTruck = val
		case key == "assigned_truck":
			d.AssignedTruck = val
		case key == "assigned_truck_efficiency":
			d.AssignedTruckEfficiency = parseFloat(val)
		case key == "assigned_truck_axle_count":
			d.AssignedTruckAxleCount = parseInt(val)
		case key == "assigned_truck_mass":
			d.AssignedTruckMass = parseFloat(val)
		case key == "slot_truck_efficiency":
			d.SlotTruckEfficiency = parseFloat(val)
		case key == "slot_truck_axle_count":
			d.SlotTruckAxleCount = parseInt(val)
		case key == "slot_truck_mass":
			d.SlotTruckMass = parseFloat(val)
		case key == "adopted_trailer":
			d.AdoptedTrailer = val
		case key == "assigned_trailer":
			d.AssignedTrailer = val
		case key == "old_hometown":
			d.OldHometown = val
		case key == "profit_log":
			d.ProfitLog = val
		}
	}
	return nil
}

// ToProperties converts the DriverAI struct to a map of properties.
func (d *DriverAI) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["adr"] = []string{strconv.Itoa(int(d.Adr))}
	props["long_dist"] = []string{strconv.Itoa(int(d.LongDist))}
	props["heavy"] = []string{strconv.Itoa(int(d.Heavy))}
	props["fragile"] = []string{strconv.Itoa(int(d.Fragile))}
	props["urgent"] = []string{strconv.Itoa(int(d.Urgent))}
	props["mechanical"] = []string{strconv.Itoa(int(d.Mechanical))}
	props["hometown"] = []string{d.Hometown}
	props["current_city"] = []string{d.CurrentCity}
	props["state"] = []string{strconv.FormatUint(uint64(d.State), 10)}
	props["on_duty_timer"] = []string{strconv.Itoa(d.OnDutyTimer)}
	props["extra_maintenance"] = []string{strconv.Itoa(d.ExtraMaintenance)}
	props["driver_job"] = []string{d.DriverJob}
	props["experience_points"] = []string{strconv.Itoa(d.ExperiencePoints)}
	props["training_policy"] = []string{strconv.Itoa(d.TrainingPolicy)}
	props["adopted_truck"] = []string{d.AdoptedTruck}
	props["assigned_truck"] = []string{d.AssignedTruck}
	props["assigned_truck_efficiency"] = []string{formatFloat(d.AssignedTruckEfficiency)}
	props["assigned_truck_axle_count"] = []string{strconv.Itoa(d.AssignedTruckAxleCount)}
	props["assigned_truck_mass"] = []string{formatFloat(d.AssignedTruckMass)}
	props["slot_truck_efficiency"] = []string{formatFloat(d.SlotTruckEfficiency)}
	props["slot_truck_axle_count"] = []string{strconv.Itoa(d.SlotTruckAxleCount)}
	props["slot_truck_mass"] = []string{formatFloat(d.SlotTruckMass)}
	props["adopted_trailer"] = []string{d.AdoptedTrailer}
	props["assigned_trailer"] = []string{d.AssignedTrailer}
	props["old_hometown"] = []string{d.OldHometown}
	props["profit_log"] = []string{d.ProfitLog}

	return props
}

func parseUint8(s string) uint8 {
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}

func parseUint32(s string) uint32 {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

