package items

import (
	"fmt"
	"strconv"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// PlayerJob mirrors the C# Player_Job class from CustomClasses/Save/Items/Player_Job.cs.
type PlayerJob struct {
	CompanyTruck            string
	CompanyTrailer          string
	TargetPlacement         dataformat.Vector4f // Vector_3f_4f stored as Vector4f
	TargetPlacementMedium   dataformat.Vector4f
	TargetPlacementHard     dataformat.Vector4f
	TargetPlacementRigid    dataformat.Vector4f
	SourcePlacement         dataformat.Vector4f
	SelectedTarget          *int // nullable int
	TimeLowerLimit          int
	TimeUpperLimit          *int // nullable int
	JobDistance             int
	FuelConsumed            dataformat.Float
	LastReportedFuel        dataformat.Float
	TotalFines              int
	IsTrailerLoaded         bool
	OnlineJobID             *int // nullable int
	OnlineJobTrailerModel   string
	AutoloadUsed            bool
	Cargo                   string
	SourceCompany           string
	TargetCompany           string
	CargoModelIndex         int
	IsArticulated           bool
	IsCargoMarketJob        bool
	StartTime               int
	PlannedDistanceKm       int
	FerryTime               int
	FerryPrice              int
	Urgency                 int
	Special                 string
	UnitsCount              int
	FillRatio               int
}

// FromProperties populates the PlayerJob from a map of SII properties.
func (p *PlayerJob) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "company_truck":
			p.CompanyTruck = val
		case key == "company_trailer":
			p.CompanyTrailer = val
		case key == "target_placement":
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse target_placement: %w", err)
			}
			p.TargetPlacement = vec
		case key == "target_placement_medium":
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse target_placement_medium: %w", err)
			}
			p.TargetPlacementMedium = vec
		case key == "target_placement_hard":
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse target_placement_hard: %w", err)
			}
			p.TargetPlacementHard = vec
		case key == "target_placement_rigid":
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse target_placement_rigid: %w", err)
			}
			p.TargetPlacementRigid = vec
		case key == "source_placement":
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse source_placement: %w", err)
			}
			p.SourcePlacement = vec
		case key == "selected_target":
			if val == "nil" {
				p.SelectedTarget = nil
			} else {
				target := parseInt(val)
				p.SelectedTarget = &target
			}
		case key == "time_lower_limit":
			p.TimeLowerLimit = parseInt(val)
		case key == "time_upper_limit":
			if val == "nil" {
				p.TimeUpperLimit = nil
			} else {
				limit := parseInt(val)
				p.TimeUpperLimit = &limit
			}
		case key == "job_distance":
			p.JobDistance = parseInt(val)
		case key == "fuel_consumed":
			p.FuelConsumed = parseFloat(val)
		case key == "last_reported_fuel":
			p.LastReportedFuel = parseFloat(val)
		case key == "total_fines":
			p.TotalFines = parseInt(val)
		case key == "is_trailer_loaded":
			p.IsTrailerLoaded = parseBool(val)
		case key == "online_job_id":
			if val == "nil" {
				p.OnlineJobID = nil
			} else {
				id := parseInt(val)
				p.OnlineJobID = &id
			}
		case key == "online_job_trailer_model":
			p.OnlineJobTrailerModel = val
		case key == "autoload_used":
			p.AutoloadUsed = parseBool(val)
		case key == "cargo":
			p.Cargo = val
		case key == "source_company":
			p.SourceCompany = val
		case key == "target_company":
			p.TargetCompany = val
		case key == "cargo_model_index":
			p.CargoModelIndex = parseInt(val)
		case key == "is_articulated":
			p.IsArticulated = parseBool(val)
		case key == "is_cargo_market_job":
			p.IsCargoMarketJob = parseBool(val)
		case key == "start_time":
			p.StartTime = parseInt(val)
		case key == "planned_distance_km":
			p.PlannedDistanceKm = parseInt(val)
		case key == "ferry_time":
			p.FerryTime = parseInt(val)
		case key == "ferry_price":
			p.FerryPrice = parseInt(val)
		case key == "urgency":
			p.Urgency = parseInt(val)
		case key == "special":
			p.Special = val
		case key == "units_count":
			p.UnitsCount = parseInt(val)
		case key == "fill_ratio":
			p.FillRatio = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the PlayerJob struct to a map of properties.
func (p *PlayerJob) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["company_truck"] = []string{p.CompanyTruck}
	props["company_trailer"] = []string{p.CompanyTrailer}
	props["target_placement"] = []string{formatVector3f4f(p.TargetPlacement)}
	props["target_placement_medium"] = []string{formatVector3f4f(p.TargetPlacementMedium)}
	props["target_placement_hard"] = []string{formatVector3f4f(p.TargetPlacementHard)}
	props["target_placement_rigid"] = []string{formatVector3f4f(p.TargetPlacementRigid)}
	props["source_placement"] = []string{formatVector3f4f(p.SourcePlacement)}
	if p.SelectedTarget == nil {
		props["selected_target"] = []string{"nil"}
	} else {
		props["selected_target"] = []string{strconv.Itoa(*p.SelectedTarget)}
	}
	props["time_lower_limit"] = []string{strconv.Itoa(p.TimeLowerLimit)}
	if p.TimeUpperLimit == nil {
		props["time_upper_limit"] = []string{"nil"}
	} else {
		props["time_upper_limit"] = []string{strconv.Itoa(*p.TimeUpperLimit)}
	}
	props["job_distance"] = []string{strconv.Itoa(p.JobDistance)}
	props["fuel_consumed"] = []string{formatFloat(p.FuelConsumed)}
	props["last_reported_fuel"] = []string{formatFloat(p.LastReportedFuel)}
	props["total_fines"] = []string{strconv.Itoa(p.TotalFines)}
	props["is_trailer_loaded"] = []string{formatBool(p.IsTrailerLoaded)}
	if p.OnlineJobID == nil {
		props["online_job_id"] = []string{"nil"}
	} else {
		props["online_job_id"] = []string{strconv.Itoa(*p.OnlineJobID)}
	}
	props["online_job_trailer_model"] = []string{p.OnlineJobTrailerModel}
	props["autoload_used"] = []string{formatBool(p.AutoloadUsed)}
	props["cargo"] = []string{p.Cargo}
	props["source_company"] = []string{p.SourceCompany}
	props["target_company"] = []string{p.TargetCompany}
	props["cargo_model_index"] = []string{strconv.Itoa(p.CargoModelIndex)}
	props["is_articulated"] = []string{formatBool(p.IsArticulated)}
	props["is_cargo_market_job"] = []string{formatBool(p.IsCargoMarketJob)}
	props["start_time"] = []string{strconv.Itoa(p.StartTime)}
	props["planned_distance_km"] = []string{strconv.Itoa(p.PlannedDistanceKm)}
	props["ferry_time"] = []string{strconv.Itoa(p.FerryTime)}
	props["ferry_price"] = []string{strconv.Itoa(p.FerryPrice)}
	props["urgency"] = []string{strconv.Itoa(p.Urgency)}
	props["special"] = []string{p.Special}
	props["units_count"] = []string{strconv.Itoa(p.UnitsCount)}
	props["fill_ratio"] = []string{strconv.Itoa(p.FillRatio)}

	return props
}

