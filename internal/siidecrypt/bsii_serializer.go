package siidecrypt

import (
	"fmt"
	"strconv"
	"strings"
)

// SerializeBSII converts decoded BSII data to text SII format
func SerializeBSII(data *BSIIData) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString("SiiNunit\n{\n")

	// Serialize all decoded blocks
	for _, block := range data.DecodedBlocks {
		if block.ID == nil {
			continue
		}

		// Write block header: type : name {
		buf.WriteString(fmt.Sprintf("%s : %s {\n", block.Name, block.ID.Value))

		// Write all segments
		for _, seg := range block.Segments {
			if seg.Value == nil {
				continue
			}

			valueStr := serializeValue(seg.Type, seg.Value)
			if valueStr != "" {
				buf.WriteString(fmt.Sprintf(" %s: %s\n", seg.Name, valueStr))
			}
		}

		buf.WriteString("}\n")
	}

	buf.WriteString("}\n")
	return []byte(buf.String()), nil
}

// serializeValue converts a decoded value to its SII text representation
func serializeValue(dataType uint32, value interface{}) string {
	switch dataType {
	case 0x01: // UTF8String
		if s, ok := value.(string); ok {
			return s
		}
	case 0x02: // ArrayOfUTF8String
		if arr, ok := value.([]string); ok {
			return serializeStringArray(arr)
		}
	case 0x03: // EncodedString
		if s, ok := value.(string); ok {
			return s
		}
	case 0x04: // ArrayOfEncodedString
		if arr, ok := value.([]string); ok {
			return serializeStringArray(arr)
		}
	case 0x05: // Single
		if f, ok := value.(float32); ok {
			return formatFloat(f)
		}
	case 0x06: // ArrayOfSingle
		if arr, ok := value.([]float32); ok {
			return serializeFloatArray(arr)
		}
	case 0x07: // VectorOf2Single
		if v, ok := value.(SingleVector2); ok {
			return fmt.Sprintf("(%f, %f)", v.A, v.B)
		}
	case 0x08: // ArrayOfVectorOf2Single
		if arr, ok := value.([]SingleVector2); ok {
			return serializeVector2Array(arr)
		}
	case 0x09: // VectorOf3Single
		if v, ok := value.(SingleVector3); ok {
			return fmt.Sprintf("(%f, %f, %f)", v.A, v.B, v.C)
		}
	case 0x0A: // ArrayOfVectorOf3Single
		if arr, ok := value.([]SingleVector3); ok {
			return serializeVector3Array(arr)
		}
	case 0x11: // VectorOf3Int32
		if v, ok := value.(Int32Vector3); ok {
			return fmt.Sprintf("(%d, %d, %d)", v.A, v.B, v.C)
		}
	case 0x12: // ArrayOfVectorOf3Int32
		if arr, ok := value.([]Int32Vector3); ok {
			return serializeInt32Vector3Array(arr)
		}
	case 0x17: // VectorOf4Single
		if v, ok := value.(SingleVector4); ok {
			return fmt.Sprintf("(%f, %f, %f, %f)", v.A, v.B, v.C, v.D)
		}
	case 0x18: // ArrayOfVectorOf4Single
		if arr, ok := value.([]SingleVector4); ok {
			return serializeVector4Array(arr)
		}
	case 0x19: // VectorOf8Single
		if v7, ok := value.(SingleVector7); ok {
			return fmt.Sprintf("(%f, %f, %f, %f, %f, %f, %f)", v7.A, v7.B, v7.C, v7.D, v7.E, v7.F, v7.G)
		}
		if v8, ok := value.(SingleVector8); ok {
			return fmt.Sprintf("(%f, %f, %f, %f, %f, %f, %f, %f)", v8.A, v8.B, v8.C, v8.D, v8.E, v8.F, v8.G, v8.H)
		}
	case 0x1A: // ArrayOfVectorOf8Single
		if arr7, ok := value.([]SingleVector7); ok {
			return serializeVector7Array(arr7)
		}
		if arr8, ok := value.([]SingleVector8); ok {
			return serializeVector8Array(arr8)
		}
	case 0x25: // Int32
		if i, ok := value.(int32); ok {
			return strconv.FormatInt(int64(i), 10)
		}
	case 0x26: // ArrayOfInt32
		if arr, ok := value.([]int32); ok {
			return serializeInt32Array(arr)
		}
	case 0x27, 0x2F: // UInt32, UInt32Type2
		if u, ok := value.(uint32); ok {
			return strconv.FormatUint(uint64(u), 10)
		}
	case 0x28: // ArrayOfUInt32
		if arr, ok := value.([]uint32); ok {
			return serializeUInt32Array(arr)
		}
	case 0x29: // Int16
		if i, ok := value.(int16); ok {
			return strconv.FormatInt(int64(i), 10)
		}
	case 0x2A: // ArrayOfInt16
		if arr, ok := value.([]int16); ok {
			return serializeInt16Array(arr)
		}
	case 0x2B: // UInt16
		if u, ok := value.(uint16); ok {
			return strconv.FormatUint(uint64(u), 10)
		}
	case 0x2C: // ArrayOfUInt16
		if arr, ok := value.([]uint16); ok {
			return serializeUInt16Array(arr)
		}
	case 0x31: // Int64
		if i, ok := value.(int64); ok {
			return strconv.FormatInt(i, 10)
		}
	case 0x32: // ArrayOfInt64
		if arr, ok := value.([]int64); ok {
			return serializeInt64Array(arr)
		}
	case 0x33: // UInt64
		if u, ok := value.(uint64); ok {
			return strconv.FormatUint(u, 10)
		}
	case 0x34: // ArrayOfUInt64
		if arr, ok := value.([]uint64); ok {
			return serializeUInt64Array(arr)
		}
	case 0x35: // ByteBool
		if b, ok := value.(bool); ok {
			return formatBool(b)
		}
	case 0x36: // ArrayOfByteBool
		if arr, ok := value.([]bool); ok {
			return serializeBoolArray(arr)
		}
	case 0x37: // OrdinalString
		if s, ok := value.(string); ok {
			return s
		}
	case 0x39, 0x3B, 0x3D: // Id, IdType2, IdType3
		if id, ok := value.(*IDComplexType); ok {
			return id.Value
		}
	case 0x3A, 0x3C, 0x3E: // ArrayOfIdA, ArrayOfIdC, ArrayOfIdE
		if arr, ok := value.([]*IDComplexType); ok {
			return serializeIDArray(arr)
		}
	}
	return ""
}

// Helper serialization functions
func serializeStringArray(arr []string) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, s := range arr {
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func serializeFloatArray(arr []float32) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, f := range arr {
		parts = append(parts, formatFloat(f))
	}
	return strings.Join(parts, ", ")
}

func serializeInt32Array(arr []int32) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, i := range arr {
		parts = append(parts, strconv.FormatInt(int64(i), 10))
	}
	return strings.Join(parts, ", ")
}

func serializeUInt32Array(arr []uint32) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, u := range arr {
		parts = append(parts, strconv.FormatUint(uint64(u), 10))
	}
	return strings.Join(parts, ", ")
}

func serializeUInt16Array(arr []uint16) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, u := range arr {
		parts = append(parts, strconv.FormatUint(uint64(u), 10))
	}
	return strings.Join(parts, ", ")
}

func serializeInt16Array(arr []int16) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, i := range arr {
		parts = append(parts, strconv.FormatInt(int64(i), 10))
	}
	return strings.Join(parts, ", ")
}

func serializeInt64Array(arr []int64) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, i := range arr {
		parts = append(parts, strconv.FormatInt(i, 10))
	}
	return strings.Join(parts, ", ")
}

func serializeUInt64Array(arr []uint64) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, u := range arr {
		parts = append(parts, strconv.FormatUint(u, 10))
	}
	return strings.Join(parts, ", ")
}

func serializeBoolArray(arr []bool) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, b := range arr {
		parts = append(parts, formatBool(b))
	}
	return strings.Join(parts, ", ")
}

func serializeVector2Array(arr []SingleVector2) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%f, %f)", v.A, v.B))
	}
	return strings.Join(parts, ", ")
}

func serializeVector3Array(arr []SingleVector3) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%f, %f, %f)", v.A, v.B, v.C))
	}
	return strings.Join(parts, ", ")
}

func serializeInt32Vector3Array(arr []Int32Vector3) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%d, %d, %d)", v.A, v.B, v.C))
	}
	return strings.Join(parts, ", ")
}

func serializeVector4Array(arr []SingleVector4) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%f, %f, %f, %f)", v.A, v.B, v.C, v.D))
	}
	return strings.Join(parts, ", ")
}

func serializeVector7Array(arr []SingleVector7) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%f, %f, %f, %f, %f, %f, %f)", v.A, v.B, v.C, v.D, v.E, v.F, v.G))
	}
	return strings.Join(parts, ", ")
}

func serializeVector8Array(arr []SingleVector8) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, v := range arr {
		parts = append(parts, fmt.Sprintf("(%f, %f, %f, %f, %f, %f, %f, %f)", v.A, v.B, v.C, v.D, v.E, v.F, v.G, v.H))
	}
	return strings.Join(parts, ", ")
}

func serializeIDArray(arr []*IDComplexType) string {
	if len(arr) == 0 {
		return "0"
	}
	var parts []string
	parts = append(parts, strconv.Itoa(len(arr)))
	for _, id := range arr {
		parts = append(parts, id.Value)
	}
	return strings.Join(parts, ", ")
}

func formatFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'g', -1, 32)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

