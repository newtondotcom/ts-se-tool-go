package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// JobOfferData mirrors the C# Job_offer_Data class from CustomClasses/Save/Items/Job_offer_Data.cs.
type JobOfferData struct {
	Target              string
	ExpirationTime      *uint32 // nullable uint
	Urgency             *int    // nullable int
	ShortestDistanceKm  int
	FerryTime           int
	FerryPrice          int
	Cargo               string
	CompanyTruck        string
	TrailerVariant      string
	TrailerDefinition   string
	UnitsCount          int
	FillRatio           int
	TrailerPlace        []dataformat.Vector4f // Vector_3f_4f is stored as Vector4f
}

// FromProperties populates the JobOfferData from a map of SII properties.
func (j *JobOfferData) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "target":
			j.Target = val
		case key == "expiration_time":
			if val == "nil" {
				j.ExpirationTime = nil
			} else {
				expTime := parseUint32(val)
				j.ExpirationTime = &expTime
			}
		case key == "urgency":
			if val == "nil" {
				j.Urgency = nil
			} else {
				urgency := parseInt(val)
				j.Urgency = &urgency
			}
		case key == "shortest_distance_km":
			j.ShortestDistanceKm = parseInt(val)
		case key == "ferry_time":
			j.FerryTime = parseInt(val)
		case key == "ferry_price":
			j.FerryPrice = parseInt(val)
		case key == "cargo":
			j.Cargo = val
		case key == "company_truck":
			j.CompanyTruck = val
		case key == "trailer_variant":
			j.TrailerVariant = val
		case key == "trailer_definition":
			j.TrailerDefinition = val
		case key == "units_count":
			j.UnitsCount = parseInt(val)
		case key == "fill_ratio":
			j.FillRatio = parseInt(val)
		case key == "trailer_place":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "trailer_place["):
			vec, err := parseVector3f4f(val)
			if err != nil {
				return fmt.Errorf("parse trailer_place: %w", err)
			}
			j.TrailerPlace = append(j.TrailerPlace, vec)
		}
	}
	return nil
}

// ToProperties converts the JobOfferData struct to a map of properties.
func (j *JobOfferData) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["target"] = []string{j.Target}
	if j.ExpirationTime == nil {
		props["expiration_time"] = []string{"nil"}
	} else {
		props["expiration_time"] = []string{strconv.FormatUint(uint64(*j.ExpirationTime), 10)}
	}
	if j.Urgency == nil {
		props["urgency"] = []string{"nil"}
	} else {
		props["urgency"] = []string{strconv.Itoa(*j.Urgency)}
	}
	props["shortest_distance_km"] = []string{strconv.Itoa(j.ShortestDistanceKm)}
	props["ferry_time"] = []string{strconv.Itoa(j.FerryTime)}
	props["ferry_price"] = []string{strconv.Itoa(j.FerryPrice)}
	props["cargo"] = []string{j.Cargo}
	props["company_truck"] = []string{j.CompanyTruck}
	props["trailer_variant"] = []string{j.TrailerVariant}
	props["trailer_definition"] = []string{j.TrailerDefinition}
	props["units_count"] = []string{strconv.Itoa(j.UnitsCount)}
	props["fill_ratio"] = []string{strconv.Itoa(j.FillRatio)}

	props["trailer_place"] = []string{strconv.Itoa(len(j.TrailerPlace))}
	for i, v := range j.TrailerPlace {
		props[fmt.Sprintf("trailer_place[%d]", i)] = []string{formatVector3f4f(v)}
	}

	return props
}

// parseVector3f4f parses a Vector_3f_4f from a string like "(&c75137e5, &40dfd91d, &453d1a51) (&3f59911d; &ba5459ea, &bf06e8dd, &ba03acdd)"
// Note: Vector_3f_4f contains both placement (Vector3f) and direction (Vector4f)
// We store it as Vector4f where X,Y,Z are from placement and W is from direction
func parseVector3f4f(s string) (dataformat.Vector4f, error) {
	// Split by parentheses to get placement and direction parts
	// Format: "(placement) (direction)"
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '(' || r == ')'
	})
	if len(parts) < 4 {
		return dataformat.Vector4f{}, fmt.Errorf("invalid vector3f4f format: expected 2 parts in parentheses")
	}

	// parts[0] is empty, parts[1] is placement, parts[2] is empty, parts[3] is direction
	placementStr := strings.TrimSpace(parts[1])
	directionStr := strings.TrimSpace(parts[3])

	// Parse placement (Vector3f)
	placement, err := parseVector3f("(" + placementStr + ")")
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse placement: %w", err)
	}

	// Parse direction (Vector4f) - format: "&3f59911d; &ba5459ea, &bf06e8dd, &ba03acdd"
	directionParts := strings.FieldsFunc(directionStr, func(r rune) bool {
		return r == ';' || r == ','
	})
	if len(directionParts) != 4 {
		return dataformat.Vector4f{}, fmt.Errorf("invalid direction format: expected 4 components, got %d", len(directionParts))
	}

	w, err := parseHexFloat(strings.TrimSpace(directionParts[0]))
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse W: %w", err)
	}
	// Note: direction X, Y, Z are not used in the stored Vector4f, only W is stored
	// The placement X, Y, Z are stored in the Vector4f X, Y, Z

	return dataformat.Vector4f{
		X: placement.X,
		Y: placement.Y,
		Z: placement.Z,
		W: w,
	}, nil
}

// formatVector3f4f formats a Vector4f to string like "(&c75137e5, &40dfd91d, &453d1a51) (&3f59911d; &ba5459ea, &bf06e8dd, &ba03acdd)"
// Note: X,Y,Z are from placement, W is from direction
func formatVector3f4f(v dataformat.Vector4f) string {
	placement := formatVector3f(dataformat.Vector3f{X: v.X, Y: v.Y, Z: v.Z})
	// For direction, we only have W, so we use W for all components (matching C# behavior)
	direction := fmt.Sprintf("(%s; %s, %s, %s)", formatHexFloat(v.W), formatHexFloat(v.W), formatHexFloat(v.W), formatHexFloat(v.W))
	return placement + " " + direction
}

