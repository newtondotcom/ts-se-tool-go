package items

import (
	"fmt"
	"strconv"
	"strings"
)

// ProfitLog mirrors the C# Profit_log class from CustomClasses/Save/Items/Profit_log.cs.
type ProfitLog struct {
	StatsData        []string
	AccDistanceFree  int
	AccDistanceOnJob int
	HistoryAge       *int // nullable int
}

// FromProperties populates the ProfitLog from a map of SII properties.
func (p *ProfitLog) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "stats_data":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "stats_data["):
			p.StatsData = append(p.StatsData, val)
		case key == "acc_distance_free":
			p.AccDistanceFree = parseInt(val)
		case key == "acc_distance_on_job":
			p.AccDistanceOnJob = parseInt(val)
		case key == "history_age":
			if val == "nil" {
				p.HistoryAge = nil
			} else {
				age := parseInt(val)
				p.HistoryAge = &age
			}
		}
	}
	return nil
}

// ToProperties converts the ProfitLog struct to a map of properties.
func (p *ProfitLog) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["stats_data"] = []string{strconv.Itoa(len(p.StatsData))}
	for i, v := range p.StatsData {
		props[fmt.Sprintf("stats_data[%d]", i)] = []string{v}
	}

	props["acc_distance_free"] = []string{strconv.Itoa(p.AccDistanceFree)}
	props["acc_distance_on_job"] = []string{strconv.Itoa(p.AccDistanceOnJob)}
	if p.HistoryAge == nil {
		props["history_age"] = []string{"nil"}
	} else {
		props["history_age"] = []string{strconv.Itoa(*p.HistoryAge)}
	}

	return props
}

