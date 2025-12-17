package siidecrypt

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// BSIICharTable is the character table for encoded strings (0x03, 0x04)
var BSIICharTable = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z', '_',
}

// BSIIHeader mirrors BSII_Header in C#
type BSIIHeader struct {
	Signature uint32
	Version   uint32
}

// BSIIStructureBlock mirrors BSII_StructureBlock in C#
type BSIIStructureBlock struct {
	Type        uint32
	StructureID uint32
	Validity    bool
	Name        string
	Segments    []BSIIDataSegment
	ID          *IDComplexType
}

// BSIIDataSegment mirrors BSII_DataSegment in C#
type BSIIDataSegment struct {
	Name  string
	Type  uint32
	Value interface{}
}

// IDComplexType mirrors IDComplexType in C#
type IDComplexType struct {
	PartCount byte
	Address   uint64
	Value     string
}

// BSIIData mirrors BSII_Data in C#
type BSIIData struct {
	Header        BSIIHeader
	Blocks        []BSIIStructureBlock
	DecodedBlocks []BSIIStructureBlock
}

// decodeUInt32 reads a uint32 from bytes at offset, advances offset
func decodeUInt32(bytes []byte, offset *int) (uint32, error) {
	if len(bytes)-*offset < 4 {
		return 0, fmt.Errorf("not enough bytes for uint32")
	}
	v := binary.LittleEndian.Uint32(bytes[*offset:])
	*offset += 4
	return v, nil
}

// decodeUTF8String reads a UTF-8 string (type 0x01)
func decodeUTF8String(bytes []byte, offset *int) (string, error) {
	length, err := decodeUInt32(bytes, offset)
	if err != nil {
		return "", err
	}
	if len(bytes)-*offset < int(length) {
		return "", fmt.Errorf("not enough bytes for UTF-8 string")
	}
	s := string(bytes[*offset : *offset+int(length)])
	*offset += int(length)
	return s, nil
}

// decodeUInt64String reads an encoded string (type 0x03)
func decodeUInt64String(bytes []byte, offset *int) (string, error) {
	value, err := decodeUInt64(bytes, offset)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	val := value
	for val != 0 {
		charIdx := int(val % 38)
		if charIdx < 0 {
			charIdx = -charIdx
		}
		charIdx--
		val = val / 38
		if charIdx >= 0 && charIdx < len(BSIICharTable) {
			result.WriteByte(BSIICharTable[charIdx])
		}
	}
	// Reverse the string
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

// decodeUInt64 reads a uint64 from bytes
func decodeUInt64(bytes []byte, offset *int) (uint64, error) {
	if len(bytes)-*offset < 8 {
		return 0, fmt.Errorf("not enough bytes for uint64")
	}
	v := binary.LittleEndian.Uint64(bytes[*offset:])
	*offset += 8
	return v, nil
}

// decodeInt32 reads an int32 from bytes
func decodeInt32(bytes []byte, offset *int) (int32, error) {
	if len(bytes)-*offset < 4 {
		return 0, fmt.Errorf("not enough bytes for int32")
	}
	v := int32(binary.LittleEndian.Uint32(bytes[*offset:]))
	*offset += 4
	return v, nil
}

// decodeSingle reads a float32 from bytes
func decodeSingle(bytes []byte, offset *int) (float32, error) {
	if len(bytes)-*offset < 4 {
		return 0, fmt.Errorf("not enough bytes for float32")
	}
	v := binary.LittleEndian.Uint32(bytes[*offset:])
	*offset += 4
	return float32FromBits(v), nil
}

// float32FromBits converts uint32 bits to float32
func float32FromBits(bits uint32) float32 {
	return math.Float32frombits(bits)
}

// decodeBool reads a bool from bytes
func decodeBool(bytes []byte, offset *int) (bool, error) {
	if len(bytes)-*offset < 1 {
		return false, fmt.Errorf("not enough bytes for bool")
	}
	v := bytes[*offset] != 0
	*offset++
	return v, nil
}

// decodeUInt16 reads a uint16 from bytes
func decodeUInt16(bytes []byte, offset *int) (uint16, error) {
	if len(bytes)-*offset < 2 {
		return 0, fmt.Errorf("not enough bytes for uint16")
	}
	v := binary.LittleEndian.Uint16(bytes[*offset:])
	*offset += 2
	return v, nil
}

// decodeInt16 reads an int16 from bytes
func decodeInt16(bytes []byte, offset *int) (int16, error) {
	if len(bytes)-*offset < 2 {
		return 0, fmt.Errorf("not enough bytes for int16")
	}
	v := int16(binary.LittleEndian.Uint16(bytes[*offset:]))
	*offset += 2
	return v, nil
}

// decodeInt64 reads an int64 from bytes
func decodeInt64(bytes []byte, offset *int) (int64, error) {
	if len(bytes)-*offset < 8 {
		return 0, fmt.Errorf("not enough bytes for int64")
	}
	v := int64(binary.LittleEndian.Uint64(bytes[*offset:]))
	*offset += 8
	return v, nil
}

// decodeID reads an ID complex type (0x39, 0x3B, 0x3D)
func decodeID(bytes []byte, offset *int) (*IDComplexType, error) {
	if len(bytes)-*offset < 1 {
		return nil, fmt.Errorf("not enough bytes for ID")
	}
	result := &IDComplexType{}
	result.PartCount = bytes[*offset]
	*offset++

	if result.PartCount == 0xFF {
		// Address-based ID
		addr, err := decodeUInt64(bytes, offset)
		if err != nil {
			return nil, err
		}
		result.Address = addr
		// Convert address to nameless format
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, addr)
		parts := make([]string, 4)
		currentPart := ""
		for i := 0; i < len(data); i++ {
			if i%2 == 0 && i > 0 {
				if i >= len(data)-2 {
					for strings.HasPrefix(currentPart, "0") {
						currentPart = currentPart[1:]
					}
				}
				if currentPart != "" {
					result.Value = currentPart + "." + result.Value
				}
				parts[3-(i/2)] = currentPart
				currentPart = ""
			}
			currentPart = fmt.Sprintf("%02x", data[i]) + currentPart
			if i == len(data)-1 {
				for strings.HasPrefix(currentPart, "0") {
					currentPart = currentPart[1:]
				}
				if currentPart != "" {
					result.Value = currentPart + "." + result.Value
				}
				parts[0] = currentPart
			}
		}
		if result.Value != "" {
			result.Value = "_nameless." + strings.TrimSuffix(result.Value, ".")
		}
	} else {
		// String-based ID
		for i := 0; i < int(result.PartCount); i++ {
			s, err := decodeUInt64String(bytes, offset)
			if err != nil {
				return nil, err
			}
			if i > 0 {
				result.Value += "."
			}
			result.Value += s
		}
		if result.PartCount == 0 {
			result.Value = "null"
		}
	}
	return result, nil
}

// decodeOrdinalStringList reads an ordinal string list (type 0x37)
func decodeOrdinalStringList(bytes []byte, offset *int) (map[uint32]string, error) {
	length, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	values := make(map[uint32]string)
	for i := uint32(0); i < length; i++ {
		ordinal, err := decodeUInt32(bytes, offset)
		if err != nil {
			return nil, err
		}
		str, err := decodeUTF8String(bytes, offset)
		if err != nil {
			return nil, err
		}
		values[ordinal] = str
	}
	return values, nil
}

// getOrdinalStringFromValues retrieves a string from ordinal values map
func getOrdinalStringFromValues(values map[uint32]string, bytes []byte, offset *int) (string, error) {
	index, err := decodeUInt32(bytes, offset)
	if err != nil {
		return "", err
	}
	if val, ok := values[index]; ok {
		return val, nil
	}
	return "", nil
}

// Array decoders
func decodeUTF8StringArray(bytes []byte, offset *int) ([]string, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]string, count)
	for i := uint32(0); i < count; i++ {
		s, err := decodeUTF8String(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}
	return result, nil
}

func decodeUInt64StringArray(bytes []byte, offset *int) ([]string, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]string, count)
	for i := uint32(0); i < count; i++ {
		s, err := decodeUInt64String(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}
	return result, nil
}

func decodeSingleArray(bytes []byte, offset *int) ([]float32, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]float32, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingle(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeInt32Array(bytes []byte, offset *int) ([]int32, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]int32, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeInt32(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeUInt32Array(bytes []byte, offset *int) ([]uint32, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]uint32, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeUInt32(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeUInt16Array(bytes []byte, offset *int) ([]uint16, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]uint16, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeUInt16(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeInt16Array(bytes []byte, offset *int) ([]int16, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]int16, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeInt16(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeInt64Array(bytes []byte, offset *int) ([]int64, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]int64, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeInt64(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeUInt64Array(bytes []byte, offset *int) ([]uint64, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]uint64, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeUInt64(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeBoolArray(bytes []byte, offset *int) ([]bool, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]bool, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeBool(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeIDArray(bytes []byte, offset *int) ([]*IDComplexType, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]*IDComplexType, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeID(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

// Vector types
type SingleVector2 struct {
	A, B float32
}

type SingleVector3 struct {
	A, B, C float32
}

type SingleVector4 struct {
	A, B, C, D float32
}

type SingleVector7 struct {
	A, B, C, D, E, F, G float32
}

type SingleVector8 struct {
	A, B, C, D, E, F, G, H float32
}

type Int32Vector3 struct {
	A, B, C int32
}

func decodeSingleVector2(bytes []byte, offset *int) (SingleVector2, error) {
	a, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector2{}, err
	}
	b, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector2{}, err
	}
	return SingleVector2{A: a, B: b}, nil
}

func decodeSingleVector2Array(bytes []byte, offset *int) ([]SingleVector2, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]SingleVector2, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingleVector2(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeSingleVector3(bytes []byte, offset *int) (SingleVector3, error) {
	a, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector3{}, err
	}
	b, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector3{}, err
	}
	c, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector3{}, err
	}
	return SingleVector3{A: a, B: b, C: c}, nil
}

func decodeSingleVector3Array(bytes []byte, offset *int) ([]SingleVector3, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]SingleVector3, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingleVector3(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeSingleVector4(bytes []byte, offset *int) (SingleVector4, error) {
	a, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector4{}, err
	}
	b, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector4{}, err
	}
	c, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector4{}, err
	}
	d, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector4{}, err
	}
	return SingleVector4{A: a, B: b, C: c, D: d}, nil
}

func decodeSingleVector4Array(bytes []byte, offset *int) ([]SingleVector4, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]SingleVector4, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingleVector4(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeSingleVector7(bytes []byte, offset *int) (SingleVector7, error) {
	a, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	b, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	c, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	d, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	e, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	f, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	g, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector7{}, err
	}
	return SingleVector7{A: a, B: b, C: c, D: d, E: e, F: f, G: g}, nil
}

func decodeSingleVector7Array(bytes []byte, offset *int) ([]SingleVector7, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]SingleVector7, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingleVector7(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeSingleVector8(bytes []byte, offset *int) (SingleVector8, error) {
	a, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	b, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	c, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	d, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	e, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	f, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	g, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}
	h, err := decodeSingle(bytes, offset)
	if err != nil {
		return SingleVector8{}, err
	}

	// Apply bias decoding (from C# code)
	bias := int64(d)
	bits := bias & 0xFFF
	bits -= 2048
	bits = bits << 9
	a += float32(bits)

	bits2 := bias >> 12
	bits2 &= 0xFFF
	bits2 -= 2048
	bits2 = bits2 << 9
	c += float32(bits2)

	return SingleVector8{A: a, B: b, C: c, D: d, E: e, F: f, G: g, H: h}, nil
}

func decodeSingleVector8Array(bytes []byte, offset *int) ([]SingleVector8, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]SingleVector8, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeSingleVector8(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func decodeInt32Vector3(bytes []byte, offset *int) (Int32Vector3, error) {
	a, err := decodeInt32(bytes, offset)
	if err != nil {
		return Int32Vector3{}, err
	}
	b, err := decodeInt32(bytes, offset)
	if err != nil {
		return Int32Vector3{}, err
	}
	c, err := decodeInt32(bytes, offset)
	if err != nil {
		return Int32Vector3{}, err
	}
	return Int32Vector3{A: a, B: b, C: c}, nil
}

func decodeInt32Vector3Array(bytes []byte, offset *int) ([]Int32Vector3, error) {
	count, err := decodeUInt32(bytes, offset)
	if err != nil {
		return nil, err
	}
	result := make([]Int32Vector3, count)
	for i := uint32(0); i < count; i++ {
		v, err := decodeInt32Vector3(bytes, offset)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}
