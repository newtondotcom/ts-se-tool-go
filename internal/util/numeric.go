package util

import "fmt"

// This file will gradually host Go equivalents of NumericUtilities from
// CustomClasses/Utilities/NumericUtilities.cs. For now only minimal placeholders
// needed by early ports are defined.

// IntegerToHexString is a small helper corresponding to
// Utilities.NumericUtilities.IntegerToHexString in C#.
func IntegerToHexString(v uint32) string {
	return fmt.Sprintf("%08X", v)
}
