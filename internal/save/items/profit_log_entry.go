package items

import (
	"strconv"
)

// ProfitLogEntry mirrors the C# Profit_log_Entry class from CustomClasses/Save/Items/Profit_log_Entry.cs.
type ProfitLogEntry struct {
	Revenue          int
	Wage             int
	Maintenance      int
	Fuel             int
	Distance         int
	DistanceOnJob    bool
	CargoCount       int
	Cargo            string
	SourceCity       string
	SourceCompany    string
	DestinationCity  string
	DestinationCompany string
	TimestampDay     int
}

// FromProperties populates the ProfitLogEntry from a map of SII properties.
func (p *ProfitLogEntry) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "revenue":
			p.Revenue = parseInt(val)
		case key == "wage":
			p.Wage = parseInt(val)
		case key == "maintenance":
			p.Maintenance = parseInt(val)
		case key == "fuel":
			p.Fuel = parseInt(val)
		case key == "distance":
			p.Distance = parseInt(val)
		case key == "distance_on_job":
			p.DistanceOnJob = parseBool(val)
		case key == "cargo_count":
			p.CargoCount = parseInt(val)
		case key == "cargo":
			p.Cargo = val
		case key == "source_city":
			p.SourceCity = val
		case key == "source_company":
			p.SourceCompany = val
		case key == "destination_city":
			p.DestinationCity = val
		case key == "destination_company":
			p.DestinationCompany = val
		case key == "timestamp_day":
			p.TimestampDay = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the ProfitLogEntry struct to a map of properties.
func (p *ProfitLogEntry) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["revenue"] = []string{strconv.Itoa(p.Revenue)}
	props["wage"] = []string{strconv.Itoa(p.Wage)}
	props["maintenance"] = []string{strconv.Itoa(p.Maintenance)}
	props["fuel"] = []string{strconv.Itoa(p.Fuel)}
	props["distance"] = []string{strconv.Itoa(p.Distance)}
	props["distance_on_job"] = []string{formatBool(p.DistanceOnJob)}
	props["cargo_count"] = []string{strconv.Itoa(p.CargoCount)}
	props["cargo"] = []string{p.Cargo}
	props["source_city"] = []string{p.SourceCity}
	props["source_company"] = []string{p.SourceCompany}
	props["destination_city"] = []string{p.DestinationCity}
	props["destination_company"] = []string{p.DestinationCompany}
	props["timestamp_day"] = []string{strconv.Itoa(p.TimestampDay)}

	return props
}

