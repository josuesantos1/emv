package tlv

import (
	"encoding/hex"
	"fmt"
)

// Parser parses EMV TLV encoded data according to EMV 4.3 Book 3
type Parser struct{}

// TLV represents a Tag-Length-Value structure.
type TLV struct {
	Tag   []byte
	Len   int
	Value []byte
}

const (
	maskMultiByteTag = 0x1F
	maskContinuation = 0x80
)

func (p *Parser) Parse(data []byte) ([]TLV, error) {

	result := make([]TLV, 0, 3)
	pos := 0

	for pos < len(data) {

		tag, usedTag, err := p.ParseTag(data[pos:])
		if err != nil {
			return nil, err
		}
		pos += usedTag

		L, usedLen, err := p.ParseLength(data[pos:])
		if err != nil {
			return nil, err
		}
		pos += usedLen

		if pos+L > len(data) {
			return nil, fmt.Errorf("value length %d exceeds buffer at position %d (total: %d)", L, pos, len(data))
		}

		value := data[pos : pos+L]
		pos += L

		result = append(result, TLV{
			Tag:   tag,
			Len:   L,
			Value: value,
		})
	}

	return result, nil
}

func (p *Parser) ParseTag(data []byte) ([]byte, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("buffer is empty")
	}

	b1 := data[0]
	tag := []byte{b1}
	pos := 1

	if (b1 & maskMultiByteTag) != maskMultiByteTag {
		return tag, pos, nil
	}

	for {
		if pos >= len(data) {
			return nil, 0, fmt.Errorf("incomplete tag: unexpected end of data at position %d", pos)
		}

		b := data[pos]
		tag = append(tag, b)
		pos++

		// Bit 8 = 0 indicates last byte of multi-byte tag
		if (b & maskContinuation) == 0 {
			break
		}
	}

	return tag, pos, nil
}

func (p *Parser) ParseLength(data []byte) (int, int, error) {
	if len(data) < 1 {
		return 0, 0, fmt.Errorf("cannot parse length from empty buffer")
	}

	firstByte := data[0]

	if firstByte <= 0x7F {
		return int(firstByte), 1, nil
	}

	if firstByte == 0x81 {
		if len(data) < 2 {
			return 0, 0, fmt.Errorf("incomplete length encoding: expected 2 bytes, got %d", len(data))
		}
		return int(data[1]), 2, nil
	}

	if firstByte == 0x82 {
		if len(data) < 3 {
			return 0, 0, fmt.Errorf("incomplete length encoding: expected 3 bytes, got %d", len(data))
		}
		length := int(data[1])<<8 | int(data[2])
		return length, 3, nil
	}

	return 0, 0, fmt.Errorf("unsupported length encoding: 0x%02X", firstByte)
}

func (t TLV) TagHex() string {
	return hex.EncodeToString(t.Tag)
}

func (t TLV) ValueHex() string {
	return hex.EncodeToString(t.Value)
}
