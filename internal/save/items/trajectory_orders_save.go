package items

import (
	"fmt"
	"strconv"
	"strings"
)

// TrajectoryOrdersSave mirrors the C# Trajectory_orders_Save class from CustomClasses/Save/Items/Trajectory_orders_Save.cs.
type TrajectoryOrdersSave struct {
	HandledArray []bool
	StageArray   []int
	TrajectoryUID uint64
}

// FromProperties populates the TrajectoryOrdersSave from a map of SII properties.
func (t *TrajectoryOrdersSave) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "handled_array":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "handled_array["):
			t.HandledArray = append(t.HandledArray, parseBool(val))
		case key == "stage_array":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "stage_array["):
			t.StageArray = append(t.StageArray, parseInt(val))
		case key == "trajectory_uid":
			t.TrajectoryUID = parseUint64(val)
		}
	}
	return nil
}

// ToProperties converts the TrajectoryOrdersSave struct to a map of properties.
func (t *TrajectoryOrdersSave) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["handled_array"] = []string{strconv.Itoa(len(t.HandledArray))}
	for i, v := range t.HandledArray {
		props[fmt.Sprintf("handled_array[%d]", i)] = []string{formatBool(v)}
	}

	props["stage_array"] = []string{strconv.Itoa(len(t.StageArray))}
	for i, v := range t.StageArray {
		props[fmt.Sprintf("stage_array[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["trajectory_uid"] = []string{strconv.FormatUint(t.TrajectoryUID, 10)}

	return props
}

