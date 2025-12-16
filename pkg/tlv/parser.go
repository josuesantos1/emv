package tlv

import ("fmt"
 "strconv")

type TLV struct {
	Pan string 
	DataValidade string
	CVM string
}

func (t *TLV) Parse(message string) error {
	for _, tag := range IsoMessage {
		if tag.Tag == "5A" {
			tagLen := tag.Length+2
			length := message[tag.Length:tagLen]
			lenInteger, err := strconv.Atoi(length)
			if err != nil {
				return err
			}

			t.Pan = message[tagLen:lenInteger+2+tag.Length]
			message = message[tag.Length-1:]
		}

		fmt.Println(t.Pan, 123)
	}
	return nil
}

