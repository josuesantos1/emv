package tlv

import (
	"fmt"
	//"strconv"
)

type Parser struct{}

type TLV struct {
	Tag   string
	Len   int
	Value string
}

func (p *Parser) Parse(data []byte) ([]TLV, error) {

	var result []TLV
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
			return nil, fmt.Errorf("value excede o buffer")
		}

		value := data[pos : pos+L]
		pos += L

		result = append(result, TLV{
			Tag:   fmt.Sprintf("%X", tag),
			Len:   L,
			Value: fmt.Sprintf("%X", value),
		})
	}

	return result, nil
}

func (*Parser) ParseTag(data []byte) ([]byte, int, error) {

	if len(data) == 0 {
		return nil, 0, fmt.Errorf("buffer is null")
	}

	b1 := data[0]
	tag := []byte{b1}
	pos := 1

	if (b1 & 0x1F) != 0x1F {
		return tag, pos, nil
	}

	for {
		if pos >= len(data) {
			return nil, 0, fmt.Errorf("tag its worng")
		}

		b := data[pos]
		tag = append(tag, b)
		pos++

		if (b & 0x80) == 0 {
			break
		}
	}

	return tag, pos, nil
}

func (*Parser) ParseLength(data []byte) (int, int, error) {

	if len(data) < 1 {
		return 0, 0, fmt.Errorf("lenght data is worng")
	}

	L := int(data[0])
	return L, 1, nil
}
