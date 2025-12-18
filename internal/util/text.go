package util

import (
	"encoding/hex"
	"strings"
)

// StringToHex converts a UTF-8 string to uppercase hexadecimal representation.
// This mirrors C# Utilities.TextUtilities.FromStringToHex.
// Example: "My Profile" → "4D792050726F66696C65"
func StringToHex(s string) string {
	bytes := []byte(s)
	hexStr := hex.EncodeToString(bytes)
	return strings.ToUpper(hexStr)
}
