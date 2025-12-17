package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// OversizeJobSave mirrors the C# Oversize_Job_save class from CustomClasses/Save/Items/Oversize_Job_save.cs.
type OversizeJobSave struct {
	FrontEscortWSPosition          dataformat.Vector3f
	BackEscortWSPosition           dataformat.Vector3f
	FrontTrajectoryUID             *uint64 // nullable uint64
	BackTrajectoryUID              *uint64 // nullable uint64
	FrontTrajectoryPosition        dataformat.Float
	BackTrajectoryPosition         dataformat.Float
	FrontEscortRotation            dataformat.Vector4f
	BackEscortRotation             dataformat.Vector4f
	FrontEscortSpeed               dataformat.Float
	BackEscortSpeed                dataformat.Float
	SpawnEscortActive              bool
	TrajectoryOrders               []string
	FrontCharType                  string
	BackCharType                   string
	OversizeManagerState           int
	OversizeManagerCurrentKdopIdx  *int // nullable int
	OversizeManagerLastValidPos    dataformat.Vector3f
	ActiveBlocksRules              []string
	FrontTypeState                 int
	BackTypeState                  int
	FrontVehicleSeed               uint32
	BackVehicleSeed                uint32
	MapRouteHash                   []int
	Offer                          string
}

// FromProperties populates the OversizeJobSave from a map of SII properties.
func (o *OversizeJobSave) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "front_escort_ws_position":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse front_escort_ws_position: %w", err)
			}
			o.FrontEscortWSPosition = vec
		case key == "back_escort_ws_position":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse back_escort_ws_position: %w", err)
			}
			o.BackEscortWSPosition = vec
		case key == "front_trajectory_uid":
			if val == "nil" {
				o.FrontTrajectoryUID = nil
			} else {
				uid := parseUint64(val)
				o.FrontTrajectoryUID = &uid
			}
		case key == "back_trajectory_uid":
			if val == "nil" {
				o.BackTrajectoryUID = nil
			} else {
				uid := parseUint64(val)
				o.BackTrajectoryUID = &uid
			}
		case key == "front_trajectory_position":
			o.FrontTrajectoryPosition = parseFloat(val)
		case key == "back_trajectory_position":
			o.BackTrajectoryPosition = parseFloat(val)
		case key == "front_escort_rotation":
			vec, err := parseVector4f(val)
			if err != nil {
				return fmt.Errorf("parse front_escort_rotation: %w", err)
			}
			o.FrontEscortRotation = vec
		case key == "back_escort_rotation":
			vec, err := parseVector4f(val)
			if err != nil {
				return fmt.Errorf("parse back_escort_rotation: %w", err)
			}
			o.BackEscortRotation = vec
		case key == "front_escort_speed":
			o.FrontEscortSpeed = parseFloat(val)
		case key == "back_escort_speed":
			o.BackEscortSpeed = parseFloat(val)
		case key == "spawn_escort_active":
			o.SpawnEscortActive = parseBool(val)
		case key == "trajectory_orders":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "trajectory_orders["):
			o.TrajectoryOrders = append(o.TrajectoryOrders, val)
		case key == "front_char_type":
			o.FrontCharType = val
		case key == "back_char_type":
			o.BackCharType = val
		case key == "oversize_manager_state":
			o.OversizeManagerState = parseInt(val)
		case key == "oversize_manager_current_kdop_idx":
			if val == "nil" {
				o.OversizeManagerCurrentKdopIdx = nil
			} else {
				idx := parseInt(val)
				o.OversizeManagerCurrentKdopIdx = &idx
			}
		case key == "oversize_manager_last_valid_pos":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse oversize_manager_last_valid_pos: %w", err)
			}
			o.OversizeManagerLastValidPos = vec
		case key == "active_blocks_rules":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "active_blocks_rules["):
			o.ActiveBlocksRules = append(o.ActiveBlocksRules, val)
		case key == "front_type_state":
			o.FrontTypeState = parseInt(val)
		case key == "back_type_state":
			o.BackTypeState = parseInt(val)
		case key == "front_vehicle_seed":
			o.FrontVehicleSeed = parseUint32(val)
		case key == "back_vehicle_seed":
			o.BackVehicleSeed = parseUint32(val)
		case key == "map_route_hash":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "map_route_hash["):
			o.MapRouteHash = append(o.MapRouteHash, parseInt(val))
		case key == "offer":
			o.Offer = val
		}
	}
	return nil
}

// ToProperties converts the OversizeJobSave struct to a map of properties.
func (o *OversizeJobSave) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["front_escort_ws_position"] = []string{formatVector3f(o.FrontEscortWSPosition)}
	props["back_escort_ws_position"] = []string{formatVector3f(o.BackEscortWSPosition)}

	if o.FrontTrajectoryUID == nil {
		props["front_trajectory_uid"] = []string{"nil"}
	} else {
		props["front_trajectory_uid"] = []string{strconv.FormatUint(*o.FrontTrajectoryUID, 10)}
	}
	if o.BackTrajectoryUID == nil {
		props["back_trajectory_uid"] = []string{"nil"}
	} else {
		props["back_trajectory_uid"] = []string{strconv.FormatUint(*o.BackTrajectoryUID, 10)}
	}

	props["front_trajectory_position"] = []string{formatFloat(o.FrontTrajectoryPosition)}
	props["back_trajectory_position"] = []string{formatFloat(o.BackTrajectoryPosition)}
	props["front_escort_rotation"] = []string{formatVector4f(o.FrontEscortRotation)}
	props["back_escort_rotation"] = []string{formatVector4f(o.BackEscortRotation)}
	props["front_escort_speed"] = []string{formatFloat(o.FrontEscortSpeed)}
	props["back_escort_speed"] = []string{formatFloat(o.BackEscortSpeed)}
	props["spawn_escort_active"] = []string{formatBool(o.SpawnEscortActive)}

	props["trajectory_orders"] = []string{strconv.Itoa(len(o.TrajectoryOrders))}
	for i, v := range o.TrajectoryOrders {
		props[fmt.Sprintf("trajectory_orders[%d]", i)] = []string{v}
	}

	props["front_char_type"] = []string{o.FrontCharType}
	props["back_char_type"] = []string{o.BackCharType}
	props["oversize_manager_state"] = []string{strconv.Itoa(o.OversizeManagerState)}
	if o.OversizeManagerCurrentKdopIdx == nil {
		props["oversize_manager_current_kdop_idx"] = []string{"nil"}
	} else {
		props["oversize_manager_current_kdop_idx"] = []string{strconv.Itoa(*o.OversizeManagerCurrentKdopIdx)}
	}
	props["oversize_manager_last_valid_pos"] = []string{formatVector3f(o.OversizeManagerLastValidPos)}

	props["active_blocks_rules"] = []string{strconv.Itoa(len(o.ActiveBlocksRules))}
	for i, v := range o.ActiveBlocksRules {
		props[fmt.Sprintf("active_blocks_rules[%d]", i)] = []string{v}
	}

	props["front_type_state"] = []string{strconv.Itoa(o.FrontTypeState)}
	props["back_type_state"] = []string{strconv.Itoa(o.BackTypeState)}
	props["front_vehicle_seed"] = []string{strconv.FormatUint(uint64(o.FrontVehicleSeed), 10)}
	props["back_vehicle_seed"] = []string{strconv.FormatUint(uint64(o.BackVehicleSeed), 10)}

	props["map_route_hash"] = []string{strconv.Itoa(len(o.MapRouteHash))}
	for i, v := range o.MapRouteHash {
		props[fmt.Sprintf("map_route_hash[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["offer"] = []string{o.Offer}

	return props
}

