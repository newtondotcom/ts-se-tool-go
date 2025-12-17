package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// MailDef mirrors the C# Mail_Def class from CustomClasses/Save/Items/Mail_Def.cs.
type MailDef struct {
	ID          int
	MailTextRef dataformat.String
	ParamKeys   []dataformat.String
	ParamValues []dataformat.String
	Read        bool
	Accepted    bool
	Expired     bool
	CustomData  int
}

// FromProperties populates the MailDef from a map of SII properties.
func (m *MailDef) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "id":
			m.ID = parseInt(val)
		case key == "mail_text_ref":
			m.MailTextRef = dataformat.String(val)
		case key == "param_keys":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "param_keys["):
			m.ParamKeys = append(m.ParamKeys, dataformat.String(val))
		case key == "param_values":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "param_values["):
			m.ParamValues = append(m.ParamValues, dataformat.String(val))
		case key == "read":
			m.Read = parseBool(val)
		case key == "accepted":
			m.Accepted = parseBool(val)
		case key == "expired":
			m.Expired = parseBool(val)
		case key == "custom_data":
			m.CustomData = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the MailDef struct to a map of properties.
func (m *MailDef) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["id"] = []string{strconv.Itoa(m.ID)}
	props["mail_text_ref"] = []string{string(m.MailTextRef)}

	props["param_keys"] = []string{strconv.Itoa(len(m.ParamKeys))}
	for i, v := range m.ParamKeys {
		props[fmt.Sprintf("param_keys[%d]", i)] = []string{string(v)}
	}

	props["param_values"] = []string{strconv.Itoa(len(m.ParamValues))}
	for i, v := range m.ParamValues {
		props[fmt.Sprintf("param_values[%d]", i)] = []string{string(v)}
	}

	props["read"] = []string{formatBool(m.Read)}
	props["accepted"] = []string{formatBool(m.Accepted)}
	props["expired"] = []string{formatBool(m.Expired)}
	props["custom_data"] = []string{strconv.Itoa(m.CustomData)}

	return props
}

