package items

import (
	"fmt"
	"strconv"
	"strings"
)

// OversizeRouteOffers mirrors the C# Oversize_Route_offers class from CustomClasses/Save/Items/Oversize_Route_offers.cs.
type OversizeRouteOffers struct {
	Offers []string
	Route  string
}

// FromProperties populates the OversizeRouteOffers from a map of SII properties.
func (o *OversizeRouteOffers) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "offers":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "offers["):
			o.Offers = append(o.Offers, val)
		case key == "route":
			o.Route = val
		}
	}
	return nil
}

// ToProperties converts the OversizeRouteOffers struct to a map of properties.
func (o *OversizeRouteOffers) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["offers"] = []string{strconv.Itoa(len(o.Offers))}
	for i, v := range o.Offers {
		props[fmt.Sprintf("offers[%d]", i)] = []string{v}
	}

	props["route"] = []string{o.Route}

	return props
}

