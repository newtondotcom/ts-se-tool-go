package items

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Company mirrors Company in C# CustomClasses/Save/Items/Company.cs
type Company struct {
	PermanentData       string
	DeliveredTrailer    string
	DeliveredPos        []dataformat.Vector3f
	JobOffer            []string
	CargoOfferSeeds     []uint32
	Discovered          bool
	ReservedTrailerSlot *int // nullable int
}

// FromProperties fills the Company struct from a map of properties.
func (c *Company) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "permanent_data":
			c.PermanentData = val
		case key == "delivered_trailer":
			c.DeliveredTrailer = val
		case key == "delivered_pos":
			// capacity hint ignored
		case strings.HasPrefix(key, "delivered_pos["):
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse delivered_pos: %w", err)
			}
			c.DeliveredPos = append(c.DeliveredPos, vec)
		case key == "job_offer":
			// capacity hint ignored
		case strings.HasPrefix(key, "job_offer["):
			c.JobOffer = append(c.JobOffer, val)
		case key == "cargo_offer_seeds":
			// capacity hint ignored
		case strings.HasPrefix(key, "cargo_offer_seeds["):
			seed, err := parseCompanyUint32(val)
			if err != nil {
				return fmt.Errorf("parse cargo_offer_seeds: %w", err)
			}
			c.CargoOfferSeeds = append(c.CargoOfferSeeds, seed)
		case key == "discovered":
			c.Discovered = parseCompanyBool(val)
		case key == "reserved_trailer_slot":
			if val == "nil" {
				c.ReservedTrailerSlot = nil
			} else {
				slot := parseCompanyInt(val)
				c.ReservedTrailerSlot = &slot
			}
		}
	}
	return nil
}

// ToProperties converts the Company struct to a map of properties.
func (c *Company) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["permanent_data"] = []string{c.PermanentData}
	props["delivered_trailer"] = []string{c.DeliveredTrailer}

	props["delivered_pos"] = []string{strconv.Itoa(len(c.DeliveredPos))}
	for i, v := range c.DeliveredPos {
		props[fmt.Sprintf("delivered_pos[%d]", i)] = []string{formatVector3f(v)}
	}

	props["job_offer"] = []string{strconv.Itoa(len(c.JobOffer))}
	for i, v := range c.JobOffer {
		props[fmt.Sprintf("job_offer[%d]", i)] = []string{v}
	}

	props["cargo_offer_seeds"] = []string{strconv.Itoa(len(c.CargoOfferSeeds))}
	for i, v := range c.CargoOfferSeeds {
		props[fmt.Sprintf("cargo_offer_seeds[%d]", i)] = []string{strconv.FormatUint(uint64(v), 10)}
	}

	props["discovered"] = []string{formatCompanyBool(c.Discovered)}

	if c.ReservedTrailerSlot == nil {
		props["reserved_trailer_slot"] = []string{"nil"}
	} else {
		props["reserved_trailer_slot"] = []string{strconv.Itoa(*c.ReservedTrailerSlot)}
	}

	return props
}

// parseVector3f parses a Vector3f from a string like "(&bd85bf17, &bd5ecfd4, &bb69a963)"
func parseVector3f(s string) (dataformat.Vector3f, error) {
	// Remove parentheses and split by comma
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return dataformat.Vector3f{}, fmt.Errorf("invalid vector3f format: expected 3 components, got %d", len(parts))
	}

	x, err := parseHexFloat(strings.TrimSpace(parts[0]))
	if err != nil {
		return dataformat.Vector3f{}, fmt.Errorf("parse X: %w", err)
	}
	y, err := parseHexFloat(strings.TrimSpace(parts[1]))
	if err != nil {
		return dataformat.Vector3f{}, fmt.Errorf("parse Y: %w", err)
	}
	z, err := parseHexFloat(strings.TrimSpace(parts[2]))
	if err != nil {
		return dataformat.Vector3f{}, fmt.Errorf("parse Z: %w", err)
	}

	return dataformat.Vector3f{X: x, Y: y, Z: z}, nil
}

// formatVector3f formats a Vector3f to string like "(&bd85bf17, &bd5ecfd4, &bb69a963)"
func formatVector3f(v dataformat.Vector3f) string {
	return fmt.Sprintf("(%s, %s, %s)", formatHexFloat(v.X), formatHexFloat(v.Y), formatHexFloat(v.Z))
}

// parseHexFloat parses a hex float string like "&bd85bf17" or "123.45"
// Mirrors NumericUtilities.HexFloatToSingleFloat
func parseHexFloat(s string) (float32, error) {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "&") {
		// Hex float format: "&bd85bf17"
		hexStr := strings.TrimPrefix(s, "&")

		// Split into 2-byte chunks
		if len(hexStr) != 8 {
			return 0, fmt.Errorf("invalid hex float: expected 8 hex digits, got %d", len(hexStr))
		}

		chunks := make([]string, 4)
		for i := 0; i < 4; i++ {
			chunks[i] = hexStr[i*2 : (i+1)*2]
		}

		// Reverse order
		for i, j := 0, len(chunks)-1; i < j; i, j = i+1, j-1 {
			chunks[i], chunks[j] = chunks[j], chunks[i]
		}

		// Parse as uint32
		hexConcat := strings.Join(chunks, "")
		val, err := strconv.ParseUint(hexConcat, 16, 32)
		if err != nil {
			return 0, err
		}

		// Convert to bytes and reverse for little-endian
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(val))
		for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
			bytes[i], bytes[j] = bytes[j], bytes[i]
		}

		// Convert to float32
		return math.Float32frombits(binary.LittleEndian.Uint32(bytes)), nil
	}
	// Regular float
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

// formatHexFloat formats a float32 to hex string like "&bd85bf17" or regular format
// Mirrors NumericUtilities.SingleFloatToString
func formatHexFloat(f float32) string {
	// Use the same logic as Float.ToString() for formatting
	if math.IsNaN(float64(f)) {
		// Format as hex
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, math.Float32bits(f))
		hexStr := ""
		for i := len(bytes) - 1; i >= 0; i-- {
			hexStr += fmt.Sprintf("%02x", bytes[i])
		}
		return "&" + hexStr
	}

	intFloat := int32(f)
	if float32(intFloat)-f != 0 || f >= 1e7 {
		// Format as hex
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, math.Float32bits(f))
		hexStr := ""
		for i := len(bytes) - 1; i >= 0; i-- {
			hexStr += fmt.Sprintf("%02x", bytes[i])
		}
		return "&" + hexStr
	}

	// Format as regular float
	return strconv.FormatFloat(float64(f), 'g', 6, 32)
}

func parseCompanyInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseCompanyUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

func parseCompanyBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func formatCompanyBool(b bool) string {
	return strconv.FormatBool(b)
}

// parseVector4f parses a Vector4f from a string like "(&3f59911d; &ba5459ea, &bf06e8dd, &ba03acdd)"
func parseVector4f(s string) (dataformat.Vector4f, error) {
	// Remove parentheses and split by semicolon/comma
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ';' || r == ','
	})
	if len(parts) != 4 {
		return dataformat.Vector4f{}, fmt.Errorf("invalid vector4f format: expected 4 components, got %d", len(parts))
	}

	w, err := parseHexFloat(strings.TrimSpace(parts[0]))
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse W: %w", err)
	}
	x, err := parseHexFloat(strings.TrimSpace(parts[1]))
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse X: %w", err)
	}
	y, err := parseHexFloat(strings.TrimSpace(parts[2]))
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse Y: %w", err)
	}
	z, err := parseHexFloat(strings.TrimSpace(parts[3]))
	if err != nil {
		return dataformat.Vector4f{}, fmt.Errorf("parse Z: %w", err)
	}

	return dataformat.Vector4f{W: w, X: x, Y: y, Z: z}, nil
}

// formatVector4f formats a Vector4f to string like "(&3f59911d; &ba5459ea, &bf06e8dd, &ba03acdd)"
func formatVector4f(v dataformat.Vector4f) string {
	return fmt.Sprintf("(%s; %s, %s, %s)", formatHexFloat(v.W), formatHexFloat(v.X), formatHexFloat(v.Y), formatHexFloat(v.Z))
}
