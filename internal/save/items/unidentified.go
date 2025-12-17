package items

import (
	"strings"
)

// Unidentified mirrors the C# Unidentified class from CustomClasses/Save/Items/Unidentified.cs.
type Unidentified struct {
	Lines []string
}

// FromProperties populates the Unidentified from a map of SII properties.
// Note: Unidentified stores raw lines, so we reconstruct them from properties.
func (u *Unidentified) FromProperties(props map[string][]string) error {
	// For Unidentified, we store the raw lines
	// This is a simplified version - in practice, you might want to preserve the original format
	var lines []string
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		for _, val := range vals {
			lines = append(lines, key+": "+val)
		}
	}
	u.Lines = lines
	return nil
}

// ToProperties converts the Unidentified struct to a map of properties.
func (u *Unidentified) ToProperties() map[string][]string {
	// For Unidentified, we return the raw lines as a single property
	// This is a simplified version
	props := make(map[string][]string)
	if len(u.Lines) > 0 {
		props["_raw"] = []string{strings.Join(u.Lines, "\n")}
	}
	return props
}

