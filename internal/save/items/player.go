package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Player mirrors the C# Player class from CustomClasses/Save/Items/Player.cs.
// It intentionally keeps field names close to the original for easier mapping.
type Player struct {
	HQCity                   string
	Trailers                 []string
	TrailerUtilizationLogs   []string
	TrailerDefs              []string
	AssignedTruck            string
	MyTruck                  string
	MyTruckPlacement         string // original is Vector_3f_4f; stored as raw string for now
	MyTruckPlacementValid    bool
	MyTrailerPlacement       string // Vector_3f_4f; raw string
	MySlaveTrailerPlacements dataformat.Float
	MyTrailerAttached        bool
	MyTrailerUsed            bool
	AssignedTrailer          string
	MyTrailer                string
	AssignedTrailerConnected bool
	TruckPlacement           string // Vector_3f_4f; raw string
	TrailerPlacement         string // Vector_3f_4f; raw string
	SlaveTrailerPlacements   dataformat.Float
	ScheduleTransferToHQ     bool
	Flags                    int
	GasPumpMoneyDebt         int
	CurrentJob               string
	CurrentBusJob            string
	SelectedJob              string
	DrivingTime              int
	SleepingCount            int
	FreeRoamDistance         int
	DiscovaryDistance        dataformat.Float
	DismissedDrivers         []string
	Trucks                   []string
	TruckProfitLogs          []string
	Drivers                  []string
	DriverReadinessTimer     []int
	DriverQuitWarned         []bool
}

// FromProperties populates the Player from a map of SII properties as produced
// by the sii parser.
func (p *Player) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "hq_city":
			p.HQCity = val
		case key == "trailers":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "trailers["):
			p.Trailers = append(p.Trailers, val)
		case key == "trailer_utilization_logs":
		case strings.HasPrefix(key, "trailer_utilization_logs["):
			p.TrailerUtilizationLogs = append(p.TrailerUtilizationLogs, val)
		case key == "trailer_defs":
		case strings.HasPrefix(key, "trailer_defs["):
			p.TrailerDefs = append(p.TrailerDefs, val)
		case key == "assigned_truck":
			p.AssignedTruck = val
		case key == "my_truck":
			p.MyTruck = val
		case key == "my_truck_placement":
			p.MyTruckPlacement = val
		case key == "my_truck_placement_valid":
			p.MyTruckPlacementValid = parseBool(val)
		case key == "my_trailer_placement":
			p.MyTrailerPlacement = val
		case key == "my_slave_trailer_placements":
			p.MySlaveTrailerPlacements = parseFloat(val)
		case key == "my_trailer_attached":
			p.MyTrailerAttached = parseBool(val)
		case key == "my_trailer_used":
			p.MyTrailerUsed = parseBool(val)
		case key == "assigned_trailer":
			p.AssignedTrailer = val
		case key == "my_trailer":
			p.MyTrailer = val
		case key == "assigned_trailer_connected":
			p.AssignedTrailerConnected = parseBool(val)
		case key == "truck_placement":
			p.TruckPlacement = val
		case key == "trailer_placement":
			p.TrailerPlacement = val
		case key == "slave_trailer_placements":
			p.SlaveTrailerPlacements = parseFloat(val)
		case key == "schedule_transfer_to_hq":
			p.ScheduleTransferToHQ = parseBool(val)
		case key == "flags":
			p.Flags = parseInt(val)
		case key == "gas_pump_money_debt":
			p.GasPumpMoneyDebt = parseInt(val)
		case key == "current_job":
			p.CurrentJob = val
		case key == "current_bus_job":
			p.CurrentBusJob = val
		case key == "selected_job":
			p.SelectedJob = val
		case key == "driving_time":
			p.DrivingTime = parseInt(val)
		case key == "sleeping_count":
			p.SleepingCount = parseInt(val)
		case key == "free_roam_distance":
			p.FreeRoamDistance = parseInt(val)
		case key == "discovary_distance":
			p.DiscovaryDistance = parseFloat(val)
		case key == "dismissed_drivers":
		case strings.HasPrefix(key, "dismissed_drivers["):
			p.DismissedDrivers = append(p.DismissedDrivers, val)
		case key == "trucks":
		case strings.HasPrefix(key, "trucks["):
			p.Trucks = append(p.Trucks, val)
		case key == "truck_profit_logs":
		case strings.HasPrefix(key, "truck_profit_logs["):
			p.TruckProfitLogs = append(p.TruckProfitLogs, val)
		case key == "drivers":
		case strings.HasPrefix(key, "drivers["):
			p.Drivers = append(p.Drivers, val)
		case key == "driver_readiness_timer":
		case strings.HasPrefix(key, "driver_readiness_timer["):
			p.DriverReadinessTimer = append(p.DriverReadinessTimer, parseInt(val))
		case key == "driver_quit_warned":
		case strings.HasPrefix(key, "driver_quit_warned["):
			p.DriverQuitWarned = append(p.DriverQuitWarned, parseBool(val))
		}
	}
	return nil
}

// ToProperties produces a map of SII properties equivalent to C# PrintOut.
// The caller can use this to construct an sii.Block{Type: "player", ...}.
func (p *Player) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["hq_city"] = []string{p.HQCity}

	props["trailers"] = []string{strconv.Itoa(len(p.Trailers))}
	for i, v := range p.Trailers {
		props[fmt.Sprintf("trailers[%d]", i)] = []string{v}
	}

	props["trailer_utilization_logs"] = []string{strconv.Itoa(len(p.TrailerUtilizationLogs))}
	for i, v := range p.TrailerUtilizationLogs {
		props[fmt.Sprintf("trailer_utilization_logs[%d]", i)] = []string{v}
	}

	props["trailer_defs"] = []string{strconv.Itoa(len(p.TrailerDefs))}
	for i, v := range p.TrailerDefs {
		props[fmt.Sprintf("trailer_defs[%d]", i)] = []string{v}
	}

	props["assigned_truck"] = []string{p.AssignedTruck}
	props["my_truck"] = []string{p.MyTruck}
	props["my_truck_placement"] = []string{p.MyTruckPlacement}
	props["my_truck_placement_valid"] = []string{formatBool(p.MyTruckPlacementValid)}
	props["my_trailer_placement"] = []string{p.MyTrailerPlacement}
	props["my_slave_trailer_placements"] = []string{formatFloat(p.MySlaveTrailerPlacements)}
	props["my_trailer_attached"] = []string{formatBool(p.MyTrailerAttached)}
	props["my_trailer_used"] = []string{formatBool(p.MyTrailerUsed)}
	props["assigned_trailer"] = []string{p.AssignedTrailer}
	props["my_trailer"] = []string{p.MyTrailer}
	props["assigned_trailer_connected"] = []string{formatBool(p.AssignedTrailerConnected)}
	props["truck_placement"] = []string{p.TruckPlacement}
	props["trailer_placement"] = []string{p.TrailerPlacement}
	props["slave_trailer_placements"] = []string{formatFloat(p.SlaveTrailerPlacements)}
	props["schedule_transfer_to_hq"] = []string{formatBool(p.ScheduleTransferToHQ)}
	props["flags"] = []string{strconv.Itoa(p.Flags)}
	props["gas_pump_money_debt"] = []string{strconv.Itoa(p.GasPumpMoneyDebt)}
	props["current_job"] = []string{p.CurrentJob}
	props["current_bus_job"] = []string{p.CurrentBusJob}
	props["selected_job"] = []string{p.SelectedJob}
	props["driving_time"] = []string{strconv.Itoa(p.DrivingTime)}
	props["sleeping_count"] = []string{strconv.Itoa(p.SleepingCount)}
	props["free_roam_distance"] = []string{strconv.Itoa(p.FreeRoamDistance)}
	props["discovary_distance"] = []string{formatFloat(p.DiscovaryDistance)}

	props["dismissed_drivers"] = []string{strconv.Itoa(len(p.DismissedDrivers))}
	for i, v := range p.DismissedDrivers {
		props[fmt.Sprintf("dismissed_drivers[%d]", i)] = []string{v}
	}

	props["trucks"] = []string{strconv.Itoa(len(p.Trucks))}
	for i, v := range p.Trucks {
		props[fmt.Sprintf("trucks[%d]", i)] = []string{v}
	}

	props["truck_profit_logs"] = []string{strconv.Itoa(len(p.TruckProfitLogs))}
	for i, v := range p.TruckProfitLogs {
		props[fmt.Sprintf("truck_profit_logs[%d]", i)] = []string{v}
	}

	props["drivers"] = []string{strconv.Itoa(len(p.Drivers))}
	for i, v := range p.Drivers {
		props[fmt.Sprintf("drivers[%d]", i)] = []string{v}
	}

	props["driver_readiness_timer"] = []string{strconv.Itoa(len(p.DriverReadinessTimer))}
	for i, v := range p.DriverReadinessTimer {
		props[fmt.Sprintf("driver_readiness_timer[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["driver_quit_warned"] = []string{strconv.Itoa(len(p.DriverQuitWarned))}
	for i, v := range p.DriverQuitWarned {
		props[fmt.Sprintf("driver_quit_warned[%d]", i)] = []string{formatBool(v)}
	}

	return props
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseBool(s string) bool {
	// C# ToString().ToLower() produit "true"/"false".
	return strings.EqualFold(s, "true")
}

func parseFloat(s string) dataformat.Float {
	f, _ := strconv.ParseFloat(s, 32)
	return dataformat.Float(f)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func formatFloat(f dataformat.Float) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}
