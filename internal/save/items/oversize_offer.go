package items

import (
	"strconv"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// OversizeOffer mirrors the C# Oversize_Offer class from CustomClasses/Save/Items/Oversize_Offer.cs.
type OversizeOffer struct {
	OfferData     string
	Truck         dataformat.String
	Expiration    uint32
	IntroCutscene dataformat.String
	OutroCutscene dataformat.String
}

// FromProperties populates the OversizeOffer from a map of SII properties.
func (o *OversizeOffer) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "offer_data":
			o.OfferData = val
		case key == "truck":
			o.Truck = dataformat.String(val)
		case key == "expiration":
			o.Expiration = parseUint32(val)
		case key == "intro_cutscene":
			o.IntroCutscene = dataformat.String(val)
		case key == "outro_cutscene":
			o.OutroCutscene = dataformat.String(val)
		}
	}
	return nil
}

// ToProperties converts the OversizeOffer struct to a map of properties.
func (o *OversizeOffer) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["offer_data"] = []string{o.OfferData}
	props["truck"] = []string{string(o.Truck)}
	props["expiration"] = []string{strconv.FormatUint(uint64(o.Expiration), 10)}
	props["intro_cutscene"] = []string{string(o.IntroCutscene)}
	props["outro_cutscene"] = []string{string(o.OutroCutscene)}

	return props
}

