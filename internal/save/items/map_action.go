package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// MapAction mirrors the C# Map_action class from CustomClasses/Save/Items/Map_action.cs.
type MapAction struct {
	IDParams    []uint64
	Name        dataformat.String
	Command     string
	NumParams   []string
	StrParams   []string
	TargetTags  int
	TargetRange int
	Type        string
}

// FromProperties populates the MapAction from a map of SII properties.
func (m *MapAction) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "id_params":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "id_params["):
			v, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("parse id_params: %w", err)
			}
			m.IDParams = append(m.IDParams, v)
		case key == "name":
			m.Name = dataformat.String(val)
		case key == "command":
			m.Command = val
		case key == "num_params":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "num_params["):
			m.NumParams = append(m.NumParams, val)
		case key == "str_params":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "str_params["):
			m.StrParams = append(m.StrParams, val)
		case key == "target_tags":
			m.TargetTags = parseInt(val)
		case key == "target_range":
			m.TargetRange = parseInt(val)
		case key == "type":
			m.Type = val
		}
	}
	return nil
}

// ToProperties converts the MapAction struct to a map of properties.
func (m *MapAction) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["id_params"] = []string{strconv.Itoa(len(m.IDParams))}
	for i, v := range m.IDParams {
		props[fmt.Sprintf("id_params[%d]", i)] = []string{strconv.FormatUint(v, 10)}
	}

	props["name"] = []string{string(m.Name)}
	props["command"] = []string{m.Command}

	props["num_params"] = []string{strconv.Itoa(len(m.NumParams))}
	for i, v := range m.NumParams {
		props[fmt.Sprintf("num_params[%d]", i)] = []string{v}
	}

	props["str_params"] = []string{strconv.Itoa(len(m.StrParams))}
	for i, v := range m.StrParams {
		props[fmt.Sprintf("str_params[%d]", i)] = []string{v}
	}

	props["target_tags"] = []string{strconv.Itoa(m.TargetTags)}
	props["target_range"] = []string{strconv.Itoa(m.TargetRange)}
	props["type"] = []string{m.Type}

	return props
}

