package items

import (
	"fmt"
	"strconv"
	"strings"
)

// OversizeOfferCtrl mirrors the C# Oversize_offer_Ctrl class from CustomClasses/Save/Items/Oversize_offer_Ctrl.cs.
type OversizeOfferCtrl struct {
	RouteOffers []string
}

// FromProperties populates the OversizeOfferCtrl from a map of SII properties.
func (o *OversizeOfferCtrl) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "route_offers":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "route_offers["):
			o.RouteOffers = append(o.RouteOffers, val)
		}
	}
	return nil
}

// ToProperties converts the OversizeOfferCtrl struct to a map of properties.
func (o *OversizeOfferCtrl) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["route_offers"] = []string{strconv.Itoa(len(o.RouteOffers))}
	for i, v := range o.RouteOffers {
		props[fmt.Sprintf("route_offers[%d]", i)] = []string{v}
	}

	return props
}

