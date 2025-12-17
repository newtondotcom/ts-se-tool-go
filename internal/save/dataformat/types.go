package dataformat

import (
	"fmt"
	"math"
	"strconv"
)

// This package contains low-level data types used across save-game items,
// mirroring CustomClasses/Save/DataFormat in the C# project.

// Color approximates SCS_Color from the C# code. It stores channels as bytes.
type Color struct {
	A uint8
	R uint8
	G uint8
	B uint8
}

// NewColorFromString parses a Color from a string (mirrors SCS_Color(string _input))
func NewColorFromString(input string) Color {
	if input == "0" {
		return Color{}
	}
	if input == "nil" {
		input = "4294967295" // Corresponds to Color.FromArgb(255, 255, 255, 255)
	}

	// Convert decimal string to uint32, then to hex string
	val, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return Color{} // Default to transparent black on error
	}
	hexInput := fmt.Sprintf("%08X", val)

	// Split hex into parts and convert to int
	// C# code: hexColorParts[0]=Alpha, hexColorParts[1]=Blue, hexColorParts[2]=Green, hexColorParts[3]=Red
	// Go: A, R, G, B
	a, _ := strconv.ParseUint(hexInput[0:2], 16, 8)
	b, _ := strconv.ParseUint(hexInput[2:4], 16, 8)
	g, _ := strconv.ParseUint(hexInput[4:6], 16, 8)
	r, _ := strconv.ParseUint(hexInput[6:8], 16, 8)

	return Color{A: byte(a), R: byte(r), G: byte(g), B: byte(b)}
}

// ToString returns the string representation of Color (mirrors SCS_Color.ToString())
func (c Color) ToString() string {
	if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
		return "0"
	}
	if c.A == 255 && c.R == 255 && c.G == 255 && c.B == 255 {
		return "nil"
	}

	// C# BitConverter.ToUInt32(new byte[] { R, G, B, 255 }, 0)
	// Go: R, G, B, A (where A is always 255 for this conversion)
	temp := uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(255)<<24
	return strconv.FormatUint(uint64(temp), 10)
}

// Float is a thin wrapper around float32 mirroring SCS_Float. It allows us to
// add helper methods later if needed.
type Float float32

// String wraps a plain string; added for symmetry with SCS_string.
type String string

// Vector2f mirrors Vector_2f.
type Vector2f struct {
	X float32
	Y float32
}

// Vector3f mirrors Vector_3f.
type Vector3f struct {
	X float32
	Y float32
	Z float32
}

// Vector4f mirrors Vector_4f and the 3f_4f variants.
type Vector4f struct {
	X float32
	Y float32
	Z float32
	W float32
}

// Vector3i mirrors Vector_3i.
type Vector3i struct {
	X int32
	Y int32
	Z int32
}

// Uint32ToFloat32 converts a uint32 to float32 by reinterpreting the bits.
func Uint32ToFloat32(v uint32) float32 {
	return math.Float32frombits(v)
}

// Float32ToUint32 converts a float32 to uint32 by reinterpreting the bits.
func Float32ToUint32(f float32) uint32 {
	return math.Float32bits(f)
}

