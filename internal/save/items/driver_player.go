package items

// DriverPlayer mirrors the C# Driver_Player class from CustomClasses/Save/Items/Driver_Player.cs.
type DriverPlayer struct {
	ProfitLog string
}

// FromProperties populates the DriverPlayer from a map of SII properties.
func (d *DriverPlayer) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "profit_log":
			d.ProfitLog = val
		}
	}
	return nil
}

// ToProperties converts the DriverPlayer struct to a map of properties.
func (d *DriverPlayer) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["profit_log"] = []string{d.ProfitLog}

	return props
}

