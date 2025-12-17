package siidecrypt

import (
	"fmt"
)

// BSIISupportedVersions are the supported BSII format versions
const (
	BSIIVersion0 = 0
	BSIIVersion1 = 1
	BSIIVersion2 = 2
	BSIIVersion3 = 3
)

// DecodeBSII decodes binary SII format to text SII format
func DecodeBSII(bytes []byte) ([]byte, error) {
	pos := 0

	fileData := &BSIIData{
		Blocks:        []BSIIStructureBlock{},
		DecodedBlocks: []BSIIStructureBlock{},
	}

	// Read header
	sig, err := decodeUInt32(bytes, &pos)
	if err != nil {
		return nil, fmt.Errorf("read signature: %w", err)
	}
	fileData.Header.Signature = sig

	ver, err := decodeUInt32(bytes, &pos)
	if err != nil {
		return nil, fmt.Errorf("read version: %w", err)
	}
	fileData.Header.Version = ver

	if ver != BSIIVersion0 && ver != BSIIVersion1 && ver != BSIIVersion2 && ver != BSIIVersion3 {
		return nil, fmt.Errorf("BSII version %d not supported", ver)
	}

	ordinalLists := make(map[uint32]map[uint32]string)

	// First pass: collect all structure definitions only
	posFirstPass := pos
	for posFirstPass < len(bytes) {
		blockType, err := decodeUInt32(bytes, &posFirstPass)
		if err != nil {
			break // End of data
		}

		if blockType == 0 {
			// Structure definition block
			currentBlock := BSIIStructureBlock{
				Type:     blockType,
				Segments: []BSIIDataSegment{},
			}

			valid, err := decodeBool(bytes, &posFirstPass)
			if err != nil {
				break // End of data
			}
			currentBlock.Validity = valid

			if !valid {
				fileData.Blocks = append(fileData.Blocks, currentBlock)
				continue
			}

			structureID, err := decodeUInt32(bytes, &posFirstPass)
			if err != nil {
				break // End of data
			}
			currentBlock.StructureID = structureID

			name, err := decodeUTF8String(bytes, &posFirstPass)
			if err != nil {
				break // End of data
			}
			currentBlock.Name = name

			// Read segments
			for {
				segment, err := readDataBlock(bytes, &posFirstPass)
				if err != nil {
					break // End of data or error
				}
				if segment.Type == 0 {
					break
				}

				if segment.Type == 0x37 { // OrdinalString
					if _, exists := ordinalLists[structureID]; !exists {
						if dict, ok := segment.Value.(map[uint32]string); ok {
							ordinalLists[structureID] = dict
						}
					}
				}

				currentBlock.Segments = append(currentBlock.Segments, segment)
			}

			// Check if we already have this structure ID
			found := false
			for i := range fileData.Blocks {
				if fileData.Blocks[i].StructureID == structureID {
					found = true
					break
				}
			}
			if !found {
				fileData.Blocks = append(fileData.Blocks, currentBlock)
			}
		} else {
			// Instance block - skip it for now (we'll process in second pass)
			// Find structure if it exists to skip properly
			var blockDataItem *BSIIStructureBlock
			for i := range fileData.Blocks {
				if fileData.Blocks[i].StructureID == blockType {
					blockDataItem = &fileData.Blocks[i]
					break
				}
			}
			if blockDataItem != nil {
				// We know the structure, skip the data
				_, err := decodeID(bytes, &posFirstPass)
				if err != nil {
					break // Can't continue
				}
				for _, seg := range blockDataItem.Segments {
					if err := skipSegmentData(bytes, &posFirstPass, seg.Type, fileData.Header.Version); err != nil {
						break // Can't continue
					}
				}
			} else {
				// Structure not defined yet - can't skip without structure
				// This means we hit an instance before its definition
				// We'll need to process it in second pass, but for now we can't skip it
				// So we break and will do a full second pass
				break
			}
		}
	}

	// Second pass: process all blocks (definitions already collected, now process instances)
	pos = 0
	// Skip header
	_, _ = decodeUInt32(bytes, &pos)
	_, _ = decodeUInt32(bytes, &pos)

	for pos < len(bytes) {
		blockType, err := decodeUInt32(bytes, &pos)
		if err != nil {
			break // End of data
		}

		if blockType != 0 {
			// Instance block - find the structure definition
			var blockDataItem *BSIIStructureBlock
			for i := range fileData.Blocks {
				if fileData.Blocks[i].StructureID == blockType {
					blockDataItem = &fileData.Blocks[i]
					break
				}
			}
			if blockDataItem == nil {
				// Structure not found - this shouldn't happen after first pass
				// Skip this instance by trying to decode ID and segments
				_, err := decodeID(bytes, &pos)
				if err != nil {
					break
				}
				// Can't skip segments without structure, so we break
				break
			}

			// Create a copy of the structure
			blockData := BSIIStructureBlock{
				StructureID: blockDataItem.StructureID,
				Name:        blockDataItem.Name,
				Type:        blockDataItem.Type,
				Validity:    blockDataItem.Validity,
				Segments:    make([]BSIIDataSegment, len(blockDataItem.Segments)),
			}
			for i, seg := range blockDataItem.Segments {
				blockData.Segments[i] = BSIIDataSegment{
					Name: seg.Name,
					Type: seg.Type,
				}
			}
			if blockDataItem.ID != nil {
				blockData.ID = &IDComplexType{
					PartCount: blockDataItem.ID.PartCount,
					Address:   blockDataItem.ID.Address,
					Value:     blockDataItem.ID.Value,
				}
			}

			// Load data for this instance
			list := ordinalLists[blockData.StructureID]
			err = loadDataBlockLocal(bytes, &pos, &blockData, fileData.Header.Version, list)
			if err != nil {
				return nil, fmt.Errorf("load data block for structure %d: %w", blockType, err)
			}

			fileData.DecodedBlocks = append(fileData.DecodedBlocks, blockData)
		} else {
			// Structure definition - skip it (already collected in first pass)
			_, _ = decodeBool(bytes, &pos)
			_, _ = decodeUInt32(bytes, &pos)     // structureID
			_, _ = decodeUTF8String(bytes, &pos) // name
			// Skip segments
			for {
				segType, err := decodeUInt32(bytes, &pos)
				if err != nil {
					break
				}
				if segType == 0 {
					break
				}
				_, _ = decodeUTF8String(bytes, &pos) // segment name
				if segType == 0x37 {
					_, _ = decodeOrdinalStringList(bytes, &pos)
				}
			}
		}
	}

	// Serialize to text SII
	return SerializeBSII(fileData)
}

// readDataBlock reads a data segment definition
func readDataBlock(bytes []byte, offset *int) (BSIIDataSegment, error) {
	result := BSIIDataSegment{}

	segType, err := decodeUInt32(bytes, offset)
	if err != nil {
		return result, err
	}
	result.Type = segType

	if result.Type != 0 {
		name, err := decodeUTF8String(bytes, offset)
		if err != nil {
			return result, err
		}
		result.Name = name
	}

	if result.Type == 0x37 { // OrdinalString
		dict, err := decodeOrdinalStringList(bytes, offset)
		if err != nil {
			return result, err
		}
		result.Value = dict
	}

	return result, nil
}

// loadDataBlockLocal loads actual data values for a block instance
func loadDataBlockLocal(bytes []byte, offset *int, block *BSIIStructureBlock, formatVersion uint32, ordinalValues map[uint32]string) error {
	id, err := decodeID(bytes, offset)
	if err != nil {
		return err
	}
	block.ID = id

	for i := range block.Segments {
		dataType := block.Segments[i].Type
		var value interface{}
		var err error

		switch dataType {
		case 0x01: // UTF8String
			value, err = decodeUTF8String(bytes, offset)
		case 0x02: // ArrayOfUTF8String
			value, err = decodeUTF8StringArray(bytes, offset)
		case 0x03: // EncodedString
			value, err = decodeUInt64String(bytes, offset)
		case 0x04: // ArrayOfEncodedString
			value, err = decodeUInt64StringArray(bytes, offset)
		case 0x05: // Single
			value, err = decodeSingle(bytes, offset)
		case 0x06: // ArrayOfSingle
			value, err = decodeSingleArray(bytes, offset)
		case 0x07: // VectorOf2Single
			value, err = decodeSingleVector2(bytes, offset)
		case 0x08: // ArrayOfVectorOf2Single
			value, err = decodeSingleVector2Array(bytes, offset)
		case 0x09: // VectorOf3Single
			value, err = decodeSingleVector3(bytes, offset)
		case 0x0A: // ArrayOfVectorOf3Single
			value, err = decodeSingleVector3Array(bytes, offset)
		case 0x11: // VectorOf3Int32
			value, err = decodeInt32Vector3(bytes, offset)
		case 0x12: // ArrayOfVectorOf3Int32
			value, err = decodeInt32Vector3Array(bytes, offset)
		case 0x17: // VectorOf4Single
			value, err = decodeSingleVector4(bytes, offset)
		case 0x18: // ArrayOfVectorOf4Single
			value, err = decodeSingleVector4Array(bytes, offset)
		case 0x19: // VectorOf8Single (or VectorOf7Single for version 0 and 1)
			if formatVersion == BSIIVersion0 || formatVersion == BSIIVersion1 {
				value, err = decodeSingleVector7(bytes, offset)
			} else {
				value, err = decodeSingleVector8(bytes, offset)
			}
		case 0x1A: // ArrayOfVectorOf8Single
			if formatVersion == BSIIVersion0 || formatVersion == BSIIVersion1 {
				value, err = decodeSingleVector7Array(bytes, offset)
			} else {
				value, err = decodeSingleVector8Array(bytes, offset)
			}
		case 0x25: // Int32
			value, err = decodeInt32(bytes, offset)
		case 0x26: // ArrayOfInt32
			value, err = decodeInt32Array(bytes, offset)
		case 0x27, 0x2F: // UInt32, UInt32Type2
			value, err = decodeUInt32(bytes, offset)
		case 0x28: // ArrayOfUInt32
			value, err = decodeUInt32Array(bytes, offset)
		case 0x29: // Int16
			value, err = decodeInt16(bytes, offset)
		case 0x2A: // ArrayOfInt16
			value, err = decodeInt16Array(bytes, offset)
		case 0x2B: // UInt16
			value, err = decodeUInt16(bytes, offset)
		case 0x2C: // ArrayOfUInt16
			value, err = decodeUInt16Array(bytes, offset)
		case 0x31: // Int64
			value, err = decodeInt64(bytes, offset)
		case 0x32: // ArrayOfInt64
			value, err = decodeInt64Array(bytes, offset)
		case 0x33: // UInt64
			value, err = decodeUInt64(bytes, offset)
		case 0x34: // ArrayOfUInt64
			value, err = decodeUInt64Array(bytes, offset)
		case 0x35: // ByteBool
			value, err = decodeBool(bytes, offset)
		case 0x36: // ArrayOfByteBool
			value, err = decodeBoolArray(bytes, offset)
		case 0x37: // OrdinalString
			value, err = getOrdinalStringFromValues(ordinalValues, bytes, offset)
		case 0x39, 0x3B, 0x3D: // Id, IdType2, IdType3
			value, err = decodeID(bytes, offset)
		case 0x3A, 0x3C, 0x3E: // ArrayOfIdA, ArrayOfIdC, ArrayOfIdE
			value, err = decodeIDArray(bytes, offset)
		default:
			// Unknown type, skip
			continue
		}

		if err != nil {
			return fmt.Errorf("decode type 0x%02X for %s: %w", dataType, block.Segments[i].Name, err)
		}
		block.Segments[i].Value = value
	}

	return nil
}

// skipSegmentData skips the data for a segment without decoding it
func skipSegmentData(bytes []byte, offset *int, dataType uint32, formatVersion uint32) error {
	switch dataType {
	case 0x01: // UTF8String
		length, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		if len(bytes)-*offset < int(length) {
			return fmt.Errorf("not enough bytes")
		}
		*offset += int(length)
	case 0x02: // ArrayOfUTF8String
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		for i := uint32(0); i < count; i++ {
			if err := skipSegmentData(bytes, offset, 0x01, formatVersion); err != nil {
				return err
			}
		}
	case 0x03: // EncodedString
		_, err := decodeUInt64(bytes, offset)
		return err
	case 0x04: // ArrayOfEncodedString
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		for i := uint32(0); i < count; i++ {
			if err := skipSegmentData(bytes, offset, 0x03, formatVersion); err != nil {
				return err
			}
		}
	case 0x05: // Single
		if len(bytes)-*offset < 4 {
			return fmt.Errorf("not enough bytes")
		}
		*offset += 4
	case 0x06: // ArrayOfSingle
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 4
	case 0x07: // VectorOf2Single
		*offset += 8
	case 0x08: // ArrayOfVectorOf2Single
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 8
	case 0x09: // VectorOf3Single
		*offset += 12
	case 0x0A: // ArrayOfVectorOf3Single
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 12
	case 0x11: // VectorOf3Int32
		*offset += 12
	case 0x12: // ArrayOfVectorOf3Int32
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 12
	case 0x17: // VectorOf4Single
		*offset += 16
	case 0x18: // ArrayOfVectorOf4Single
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 16
	case 0x19: // VectorOf8Single
		if formatVersion == BSIIVersion0 || formatVersion == BSIIVersion1 {
			*offset += 28 // Vector7
		} else {
			*offset += 32 // Vector8
		}
	case 0x1A: // ArrayOfVectorOf8Single
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		if formatVersion == BSIIVersion0 || formatVersion == BSIIVersion1 {
			*offset += int(count) * 28
		} else {
			*offset += int(count) * 32
		}
	case 0x25: // Int32
		*offset += 4
	case 0x26: // ArrayOfInt32
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 4
	case 0x27, 0x2F: // UInt32, UInt32Type2
		*offset += 4
	case 0x28: // ArrayOfUInt32
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 4
	case 0x29: // Int16
		*offset += 2
	case 0x2A: // ArrayOfInt16
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 2
	case 0x2B: // UInt16
		*offset += 2
	case 0x2C: // ArrayOfUInt16
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 2
	case 0x31: // Int64
		*offset += 8
	case 0x32: // ArrayOfInt64
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 8
	case 0x33: // UInt64
		*offset += 8
	case 0x34: // ArrayOfUInt64
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 8
	case 0x35: // ByteBool
		*offset += 1
	case 0x36: // ArrayOfByteBool
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		*offset += int(count) * 1
	case 0x37: // OrdinalString
		*offset += 4 // Just the index
	case 0x39, 0x3B, 0x3D: // Id, IdType2, IdType3
		// Skip ID - complex, need to decode partially
		if len(bytes)-*offset < 1 {
			return fmt.Errorf("not enough bytes")
		}
		partCount := bytes[*offset]
		*offset++
		if partCount == 0xFF {
			*offset += 8 // Address
		} else {
			for i := 0; i < int(partCount); i++ {
				_, err := decodeUInt64(bytes, offset)
				if err != nil {
					return err
				}
			}
		}
	case 0x3A, 0x3C, 0x3E: // ArrayOfIdA, ArrayOfIdC, ArrayOfIdE
		count, err := decodeUInt32(bytes, offset)
		if err != nil {
			return err
		}
		for i := uint32(0); i < count; i++ {
			if err := skipSegmentData(bytes, offset, 0x39, formatVersion); err != nil {
				return err
			}
		}
	}
	return nil
}
