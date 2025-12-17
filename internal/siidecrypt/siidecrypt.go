package siidecrypt

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// SignatureType values copied from SIIDecryptSharp/FileType.cs
type SignatureType uint32

const (
	SignatureUnknown   SignatureType = 999
	SignaturePlainText SignatureType = 1315531091
	SignatureEncrypted SignatureType = 1131635539
	SignatureBinary    SignatureType = 1229542210
	Signature3nK       SignatureType = 21720627
)

// SIIHeader mirrors SII_Header in the C# code.
type SIIHeader struct {
	Signature uint32
	DataSize  uint32
}

// SIIData mirrors SII_Data in the C# code.
type SIIData struct {
	Header SIIHeader
	Data   []byte
}

// Key copied from Decryptor.SII_Key in C#.
var siiKey = []byte{
	0x2a, 0x5f, 0xcb, 0x17, 0x91, 0xd2, 0x2f, 0xb6,
	0x02, 0x45, 0xb3, 0xd8, 0x36, 0x9e, 0xd0, 0xb2,
	0xc2, 0x73, 0x71, 0x56, 0x3f, 0xbf, 0x1f, 0x3c,
	0x9e, 0xdf, 0x6b, 0x11, 0x82, 0x5a, 0x5d, 0x0a,
}

// DecryptFile is the Go equivalent of Decryptor.Decrypt(filePath, decode=true).
// It returns the decoded payload (typically starting with "SiiNunit" for
// plaintext SII files). Binary BSII and 3nK formats are not yet supported.
func DecryptFile(path string, decode bool) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	pos := 0
	fileType, ok := tryReadUint32(bytes, &pos)
	if !ok {
		return nil, errors.New("invalid file: cannot read signature")
	}

	if SignatureType(fileType) == SignatureEncrypted {
		// In the original C# implementation, the same buffer (starting at 0)
		// is passed to Decrypt, which in turn re-reads the signature and
		// header fields. We mirror that behaviour here and always start from 0.
		data, err := decrypt(bytes)
		if err != nil {
			return nil, err
		}

		// zlib uncompress, mirroring Zlib.uncompress in C#
		r, err := zlib.NewReader(bytesReader(data.Data))
		if err != nil {
			return nil, fmt.Errorf("zlib uncompress: %w", err)
		}
		decompressed, err := io.ReadAll(r)
		_ = r.Close()
		if err != nil {
			return nil, fmt.Errorf("zlib read: %w", err)
		}
		data.Data = decompressed
		bytes = decompressed
	}

	if !decode {
		return bytes, nil
	}

	// Decode by inner signature (plaintext / binary / 3nK)
	pos = 0
	dataType, ok := tryReadUint32(bytes, &pos)
	if !ok {
		return nil, errors.New("invalid data: cannot read inner signature")
	}

	switch SignatureType(dataType) {
	case SignaturePlainText:
		// Equivalent to DecodePlaintext: return as-is.
		return bytes, nil
	case SignatureBinary:
		// Decode BSII (binary SII format)
		return DecodeBSII(bytes[pos:])
	case Signature3nK:
		return nil, errors.New("_3nK decoding not implemented yet")
	default:
		return nil, fmt.Errorf("unknown inner signature: 0x%08X", dataType)
	}
}

// decrypt is the Go equivalent of Decryptor.Decrypt(ref byte[] encrypted, int offset)
// for the AES-encrypted SII payload (before zlib). Like the original C# code,
// it always starts reading from the beginning of the buffer.
func decrypt(encrypted []byte) (*SIIData, error) {
	header := SIIHeader{}
	var hmac []byte
	var iv []byte

	pos := 0

	if len(encrypted)-pos >= 4 {
		header.Signature = binary.LittleEndian.Uint32(encrypted[pos:])
		pos += 4
	}
	if len(encrypted)-pos >= 32 {
		hmac = encrypted[pos : pos+32]
		_ = hmac // HMAC is currently unused, same as original C# code.
		pos += 32
	}
	if len(encrypted)-pos >= 16 {
		iv = encrypted[pos : pos+16]
		pos += 16
	}
	if len(encrypted)-pos >= 4 {
		header.DataSize = binary.LittleEndian.Uint32(encrypted[pos:])
		pos += 4
	}

	finalEncrypted := encrypted[pos:]
	if len(finalEncrypted)%aes.BlockSize != 0 {
		return nil, errors.New("encrypted data is not a multiple of AES block size")
	}

	block, err := aes.NewCipher(siiKey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV size for AES")
	}

	dst := make([]byte, len(finalEncrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, finalEncrypted)

	// PKCS#7 unpadding
	if len(dst) == 0 {
		return nil, errors.New("decrypted data is empty")
	}
	pad := int(dst[len(dst)-1])
	if pad <= 0 || pad > len(dst) {
		return nil, errors.New("invalid PKCS#7 padding")
	}
	dst = dst[:len(dst)-pad]

	return &SIIData{
		Header: header,
		Data:   dst,
	}, nil
}

func tryReadUint32(b []byte, offset *int) (uint32, bool) {
	if len(b)-*offset < 4 {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(b[*offset:])
	*offset += 4
	return v, true
}

// bytesReader is a tiny helper to adapt a []byte to io.Reader without copying.
func bytesReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}
