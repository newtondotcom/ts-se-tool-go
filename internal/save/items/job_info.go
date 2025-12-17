package items

import (
	"strconv"
)

// JobInfo mirrors the C# Job_Info class from CustomClasses/Save/Items/Job_Info.cs.
type JobInfo struct {
	Cargo            string
	SourceCompany    string
	TargetCompany    string
	CargoModelIndex  int
	IsArticulated    bool
	IsCargoMarketJob bool
	StartTime        int
	PlannedDistanceKm int
	FerryTime        int
	FerryPrice       int
	Urgency          *int // nullable int
	Special          string
	UnitsCount       int
	FillRatio        int
}

// FromProperties populates the JobInfo from a map of SII properties.
func (j *JobInfo) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "cargo":
			j.Cargo = val
		case key == "source_company":
			j.SourceCompany = val
		case key == "target_company":
			j.TargetCompany = val
		case key == "cargo_model_index":
			j.CargoModelIndex = parseInt(val)
		case key == "is_articulated":
			j.IsArticulated = parseBool(val)
		case key == "is_cargo_market_job":
			j.IsCargoMarketJob = parseBool(val)
		case key == "start_time":
			j.StartTime = parseInt(val)
		case key == "planned_distance_km":
			j.PlannedDistanceKm = parseInt(val)
		case key == "ferry_time":
			j.FerryTime = parseInt(val)
		case key == "ferry_price":
			j.FerryPrice = parseInt(val)
		case key == "urgency":
			if val == "nil" {
				j.Urgency = nil
			} else {
				urgency := parseInt(val)
				j.Urgency = &urgency
			}
		case key == "special":
			j.Special = val
		case key == "units_count":
			j.UnitsCount = parseInt(val)
		case key == "fill_ratio":
			j.FillRatio = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the JobInfo struct to a map of properties.
func (j *JobInfo) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["cargo"] = []string{j.Cargo}
	props["source_company"] = []string{j.SourceCompany}
	props["target_company"] = []string{j.TargetCompany}
	props["cargo_model_index"] = []string{strconv.Itoa(j.CargoModelIndex)}
	props["is_articulated"] = []string{formatBool(j.IsArticulated)}
	props["is_cargo_market_job"] = []string{formatBool(j.IsCargoMarketJob)}
	props["start_time"] = []string{strconv.Itoa(j.StartTime)}
	props["planned_distance_km"] = []string{strconv.Itoa(j.PlannedDistanceKm)}
	props["ferry_time"] = []string{strconv.Itoa(j.FerryTime)}
	props["ferry_price"] = []string{strconv.Itoa(j.FerryPrice)}
	if j.Urgency == nil {
		props["urgency"] = []string{"nil"}
	} else {
		props["urgency"] = []string{strconv.Itoa(*j.Urgency)}
	}
	props["special"] = []string{j.Special}
	props["units_count"] = []string{strconv.Itoa(j.UnitsCount)}
	props["fill_ratio"] = []string{strconv.Itoa(j.FillRatio)}

	return props
}

