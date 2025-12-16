package tlv

import (
	"fmt"
	"strconv"
)

type TLV struct {
	Pan          string
	DataValidade string
	CVM          string
}

func (t *TLV) Parse(message string) error {
	for _, tag := range IsoMessage {
		tagLen := tag.Length + 2
		length := message[tag.Length:tagLen]
		lenInteger, err := strconv.Atoi(length)
		if err != nil {
			return err
		}

		dataLen := lenInteger + tagLen

		if tag.Tag == "5A" {
			t.Pan = message[tagLen:dataLen]
			message = message[dataLen:]
		}

		if tag.Tag == "5F24" {
			t.DataValidade = message[tagLen:dataLen]
			message = message[dataLen:]
		}

		if tag.Tag == "9F34" {
			t.CVM = message[tagLen:dataLen]
			message = message[tagLen:]
		}

		fmt.Println(message, 123)
	}
	return nil
}

func (t *TLV) String() string {
	return fmt.Sprintf("%s %s %s", t.Pan, t.DataValidade, t.CVM)
}
