package siidecrypt

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SerializeBSII converts decoded BSII data to text SII format (mirrors C# BSII_Serializer.Serialize)
func SerializeBSII(data *BSIIData) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("SiiNunit\n")
	sb.WriteString("{\n")

	indent := ""
	for _, block := range data.DecodedBlocks {
		if block.Name == "" || block.ID == nil || block.ID.Value == "" {
			continue
		}

		sb.WriteString(block.Name + " : " + block.ID.Value + " {\n")
		indent = " "

		for _, segment := range block.Segments {
			if segment.Type == 0 {
				continue
			}

			var output string
			switch segment.Type {
			case 0x36: // ArrayOfByteBool
				output = serializeByteBoolArray(&segment, indent)
			case 0x04: // ArrayOfEncodedString
				output = serializeEncodedStringArray(&segment, indent)
			case 0x3A, 0x3C, 0x3E: // ArrayOfIdA, ArrayOfIdC, ArrayOfIdE
				output = serializeIDArray(&segment, indent)
			case 0x26: // ArrayOfInt32
				output = serializeInt32Array(&segment, indent)
			case 0x06: // ArrayOfSingle
				output = serializeSingleArray(&segment, indent)
			case 0x2C: // ArrayOfUInt16
				output = serializeUInt16Array(&segment, indent)
			case 0x28: // ArrayOfUInt32
				output = serializeUInt32Array(&segment, indent)
			case 0x34: // ArrayOfUInt64
				output = serializeUInt64Array(&segment, indent)
			case 0x02: // ArrayOfUTF8String
				output = serializeUTF8StringArray(&segment, indent)
			case 0x12: // ArrayOfVectorOf3Int32
				output = serializeInt32Vector3Array(&segment, indent)
			case 0x0A: // ArrayOfVectorOf3Single
				output = serializeSingleVector3Array(&segment, indent)
			case 0x18: // ArrayOfVectorOf4Single
				output = serializeSingleVector4Array(&segment, indent)
			case 0x1A: // ArrayOfVectorOf8Single
				if data.Header.Version == BSIIVersion1 {
					output = serializeSingleVector7Array(&segment, indent)
				} else {
					output = serializeSingleVector8Array(&segment, indent)
				}
			case 0x35: // ByteBool
				output = serializeBool(&segment, indent)
			case 0x03: // EncodedString
				output = serializeEncodedString(&segment, indent)
			case 0x3D, 0x3B, 0x39: // IdType3, IdType2, Id
				output = serializeId(&segment, indent)
			case 0x25: // Int32
				output = serializeInt32(&segment, indent)
			case 0x31: // Int64
				output = serializeInt64(&segment, indent)
			case 0x2F, 0x27: // UInt32Type2, UInt32
				output = serializeUInt32(&segment, indent)
			case 0x33: // UInt64
				output = serializeUInt64(&segment, indent)
			case 0x2B: // UInt16
				output = serializeUInt16(&segment, indent)
			case 0x37: // OrdinalString
				output = serializeOrdinalString(&segment, indent)
			case 0x05: // Single
				output = serializeSingleValue(&segment, indent)
			case 0x01: // UTF8String
				output = serializeUTF8String(&segment, indent)
			case 0x07: // VectorOf2Single
				output = serializeSingleVector2(&segment, indent)
			case 0x11: // VectorOf3Int32
				output = serializeInt32Vector3(&segment, indent)
			case 0x09: // VectorOf3Single
				output = serializeSingleVector3(&segment, indent)
			case 0x17: // VectorOf4Single
				output = serializeSingleVector4(&segment, indent)
			case 0x19: // VectorOf8Single
				if data.Header.Version == BSIIVersion1 {
					output = serializeSingleVector7(&segment, indent)
				} else {
					output = serializeSingleVector8(&segment, indent)
				}
			case 0x32: // ArrayOfInt64
				output = serializeInt64Array(&segment, indent)
			case 0x08: // ArrayOfVectorOf2Single
				output = serializeSingleVector2Array(&segment, indent)
			case 0x29: // Int16
				output = serializeInt16(&segment, indent)
			case 0x2A: // ArrayOfInt16
				output = serializeInt16Array(&segment, indent)
			}
			sb.WriteString(output)
		}
		sb.WriteString("}\n")
		sb.WriteString("\n")
	}
	sb.WriteString("}")
	return []byte(sb.String()), nil
}

// formatSingle formats a float32 according to C# logic: "nil" if null, "&hex" if decimal or >= 1e7, else integer
func formatSingle(f *float32) string {
	if f == nil {
		return "nil"
	}
	val := *f
	if val-float32(math.Trunc(float64(val))) != 0.00 || val >= 1e7 {
		// Format as hex
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, math.Float32bits(val))
		hexStr := ""
		for i := len(bytes) - 1; i >= 0; i-- {
			hexStr += fmt.Sprintf("%02x", bytes[i])
		}
		return "&" + hexStr
	}
	return strconv.FormatInt(int64(val), 10)
}

func serializeByteBoolArray(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]bool)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent)
	sb.WriteString(seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strings.ToLower(strconv.FormatBool(v)) + "\n")
	}
	return sb.String()
}

func serializeEncodedStringArray(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]string)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + v + "\n")
	}
	return sb.String()
}

func serializeIDArray(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]*IDComplexType)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + v.Value + "\n")
	}
	return sb.String()
}

func serializeInt32Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]int32)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatInt(int64(v), 10) + "\n")
	}
	return sb.String()
}

func serializeSingleArray(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]float32)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + formatSingle(&v) + "\n")
	}
	return sb.String()
}

func serializeUInt16Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]uint16)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatUint(uint64(v), 10) + "\n")
	}
	return sb.String()
}

func serializeUInt32Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]uint32)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatUint(uint64(v), 10) + "\n")
	}
	return sb.String()
}

func serializeUInt64Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]uint64)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatUint(v, 10) + "\n")
	}
	return sb.String()
}

func serializeUTF8StringArray(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]string)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		if _, err := strconv.ParseInt(v, 10, 32); err == nil {
			sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + v + "\n")
		} else {
			if v == "" {
				sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: \"\"\n")
			} else {
				if isLimitedAlphabet(v) {
					sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + v + "\n")
				} else {
					sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: \"" + v + "\"\n")
				}
			}
		}
	}
	return sb.String()
}

func serializeInt32Vector3Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]Int32Vector3)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + strconv.FormatInt(int64(v.A), 10) + ", " + strconv.FormatInt(int64(v.B), 10) + ", " + strconv.FormatInt(int64(v.C), 10) + ")\n")
	}
	return sb.String()
}

func serializeSingleVector3Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]SingleVector3)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ")\n")
	}
	return sb.String()
}

func serializeSingleVector2Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]SingleVector2)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ")\n")
	}
	return sb.String()
}

func serializeSingleVector4Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]SingleVector4)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + formatSingle(&v.A) + "; " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ", " + formatSingle(&v.D) + ")\n")
	}
	return sb.String()
}

func serializeSingleVector8Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]SingleVector8)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ") (" +
			formatSingle(&v.E) + "; " + formatSingle(&v.F) + ", " + formatSingle(&v.G) + ", " + formatSingle(&v.H) + ")\n")
	}
	return sb.String()
}

func serializeSingleVector7Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]SingleVector7)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ") (" +
			formatSingle(&v.D) + "; " + formatSingle(&v.E) + ", " + formatSingle(&v.F) + ", " + formatSingle(&v.G) + ")\n")
	}
	return sb.String()
}

func serializeBool(seg *BSIIDataSegment, indent string) string {
	b, ok := seg.Value.(bool)
	if !ok {
		b = false
	}
	return indent + seg.Name + ": " + strings.ToLower(strconv.FormatBool(b)) + "\n"
}

func serializeEncodedString(seg *BSIIDataSegment, indent string) string {
	s, ok := seg.Value.(string)
	if !ok {
		s = ""
	}
	if s == "" {
		s = "\"\""
	}
	return indent + seg.Name + ": " + s + "\n"
}

func serializeId(seg *BSIIDataSegment, indent string) string {
	id, ok := seg.Value.(*IDComplexType)
	if !ok {
		return ""
	}
	return indent + seg.Name + ": " + id.Value + "\n"
}

func serializeInt32(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(int32)
	text := "nil"
	if ok {
		text = strconv.FormatInt(int64(val), 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

func serializeInt64(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(int64)
	text := "nil"
	if ok {
		text = strconv.FormatInt(val, 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

func serializeUInt32(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(uint32)
	text := "nil"
	if ok && val != 4294967295 {
		text = strconv.FormatUint(uint64(val), 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

func serializeUInt64(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(uint64)
	text := "nil"
	if ok {
		text = strconv.FormatUint(val, 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

func serializeUInt16(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(uint16)
	text := "nil"
	if ok && val != 65535 {
		text = strconv.FormatUint(uint64(val), 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

func serializeOrdinalString(seg *BSIIDataSegment, indent string) string {
	s, ok := seg.Value.(string)
	if !ok {
		s = ""
	}
	return indent + seg.Name + ": " + s + "\n"
}

func serializeSingleValue(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(float32)
	var f *float32
	if ok {
		f = &val
	}
	return indent + seg.Name + ": " + formatSingle(f) + "\n"
}

func serializeUTF8String(seg *BSIIDataSegment, indent string) string {
	s, ok := seg.Value.(string)
	if !ok {
		s = ""
	}
	sb := strings.Builder{}
	sb.WriteString(indent + seg.Name + ": ")
	if _, err := strconv.ParseInt(s, 10, 32); err == nil {
		sb.WriteString(s + "\n")
	} else {
		if s == "" {
			sb.WriteString("\"\"\n")
		} else {
			if isLimitedAlphabet(s) {
				sb.WriteString(s + "\n")
			} else {
				sb.WriteString("\"" + s + "\"\n")
			}
		}
	}
	return sb.String()
}

func serializeSingleVector2(seg *BSIIDataSegment, indent string) string {
	var v SingleVector2
	if vPtr, ok := seg.Value.(*SingleVector2); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(SingleVector2); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ")\n"
}

func serializeSingleVector3(seg *BSIIDataSegment, indent string) string {
	var v SingleVector3
	if vPtr, ok := seg.Value.(*SingleVector3); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(SingleVector3); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ")\n"
}

func serializeSingleVector4(seg *BSIIDataSegment, indent string) string {
	var v SingleVector4
	if vPtr, ok := seg.Value.(*SingleVector4); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(SingleVector4); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + formatSingle(&v.A) + "; " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ", " + formatSingle(&v.D) + ")\n"
}

func serializeInt32Vector3(seg *BSIIDataSegment, indent string) string {
	var v Int32Vector3
	if vPtr, ok := seg.Value.(*Int32Vector3); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(Int32Vector3); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + strconv.FormatInt(int64(v.A), 10) + ", " + strconv.FormatInt(int64(v.B), 10) + ", " + strconv.FormatInt(int64(v.C), 10) + ")\n"
}

func serializeSingleVector8(seg *BSIIDataSegment, indent string) string {
	var v SingleVector8
	if vPtr, ok := seg.Value.(*SingleVector8); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(SingleVector8); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ") (" +
		formatSingle(&v.E) + "; " + formatSingle(&v.F) + ", " + formatSingle(&v.G) + ", " + formatSingle(&v.H) + ")\n"
}

func serializeSingleVector7(seg *BSIIDataSegment, indent string) string {
	var v SingleVector7
	if vPtr, ok := seg.Value.(*SingleVector7); ok {
		v = *vPtr
	} else if vVal, ok := seg.Value.(SingleVector7); ok {
		v = vVal
	} else {
		return ""
	}
	return indent + seg.Name + ": (" + formatSingle(&v.A) + ", " + formatSingle(&v.B) + ", " + formatSingle(&v.C) + ") (" +
		formatSingle(&v.D) + "; " + formatSingle(&v.E) + ", " + formatSingle(&v.F) + ", " + formatSingle(&v.G) + ")\n"
}

func serializeInt64Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]int64)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatInt(v, 10) + "\n")
	}
	return sb.String()
}

func serializeInt16Array(seg *BSIIDataSegment, indent string) string {
	arr, ok := seg.Value.([]int16)
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(indent + seg.Name + ": " + strconv.Itoa(len(arr)) + "\n")
	for i, v := range arr {
		sb.WriteString(indent + seg.Name + "[" + strconv.Itoa(i) + "]: " + strconv.FormatInt(int64(v), 10) + "\n")
	}
	return sb.String()
}

func serializeInt16(seg *BSIIDataSegment, indent string) string {
	val, ok := seg.Value.(int16)
	text := "nil"
	if ok && val != 32767 {
		text = strconv.FormatInt(int64(val), 10)
	}
	return indent + seg.Name + ": " + text + "\n"
}

// isLimitedAlphabet checks if string contains only limited alphabet characters
func isLimitedAlphabet(s string) bool {
	limitedAlphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	for _, r := range s {
		if !strings.ContainsRune(limitedAlphabet, r) {
			return false
		}
	}
	return true
}
