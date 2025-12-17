package items

import (
	"strconv"
)

// OversizeBlockRuleSave mirrors the C# Oversize_Block_rule_Save class from CustomClasses/Save/Items/Oversize_Block_rule_Save.cs.
type OversizeBlockRuleSave struct {
	EscortCharType      int
	ParentTrajectoryUID uint64
	ParentTrajectoryIdx int
}

// FromProperties populates the OversizeBlockRuleSave from a map of SII properties.
func (o *OversizeBlockRuleSave) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "escort_char_type":
			o.EscortCharType = parseInt(val)
		case key == "parent_trajectory_uid":
			o.ParentTrajectoryUID = parseUint64(val)
		case key == "parent_trajectory_idx":
			o.ParentTrajectoryIdx = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the OversizeBlockRuleSave struct to a map of properties.
func (o *OversizeBlockRuleSave) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["escort_char_type"] = []string{strconv.Itoa(o.EscortCharType)}
	props["parent_trajectory_uid"] = []string{strconv.FormatUint(o.ParentTrajectoryUID, 10)}
	props["parent_trajectory_idx"] = []string{strconv.Itoa(o.ParentTrajectoryIdx)}

	return props
}

func parseUint64(s string) uint64 {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

