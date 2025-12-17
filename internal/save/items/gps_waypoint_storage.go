package items

import (
	"fmt"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// GPSWaypointStorage mirrors the C# GPS_waypoint_Storage class from CustomClasses/Save/Items/GPS_waypoint_Storage.cs.
type GPSWaypointStorage struct {
	NavNodePosition dataformat.Vector3i
	Direction       string
}

// FromProperties populates the GPSWaypointStorage from a map of SII properties.
func (g *GPSWaypointStorage) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "nav_node_position":
			vec, err := parseVector3i(val)
			if err != nil {
				return fmt.Errorf("parse nav_node_position: %w", err)
			}
			g.NavNodePosition = vec
		case key == "direction":
			g.Direction = val
		}
	}
	return nil
}

// ToProperties converts the GPSWaypointStorage struct to a map of properties.
func (g *GPSWaypointStorage) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["nav_node_position"] = []string{formatVector3i(g.NavNodePosition)}
	props["direction"] = []string{g.Direction}

	return props
}

