package tlv

import ("time"
"github.com/josuesantos1/emv/pkg/tlv"
"fmt")

type Tlv struct {
	Pan string
	DataValidade time.Time
	CVM string
}

func (t *Tlv) Populate(tlvs []tlv.TLV) error {
	for _, tlvItem := range tlvs {
		if tlvItem.Tag == "5A" {
			t.Pan = tlvItem.Value
		}

		if tlvItem.Tag == "5F24" {
			yy := tlvItem.Value[0:2]
			mm := tlvItem.Value[2:4]
			dd := "01"
			year := "20" + yy

			dateStr := fmt.Sprintf("%s-%s-%s", year, mm, dd)

			parsedTime, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return err
			}
			
			t.DataValidade = parsedTime
		}

		if tlvItem.Tag == "9F34" {
			t.CVM = tlvItem.Value
		}
	}

	return nil
}

func (t *Tlv) Validate() error {
	fields := map[string]any{
		"Pan": t.Pan,
		"Data de validade": t.DataValidade,
		"CVM": t.CVM,
	}

	for _, field := range fields {
		if field == "" || field == nil {
			return fmt.Errorf("Field %s is required", field)
		}
	}

	if len(t.Pan) <= 13 || len(t.Pan) >= 19 {
		return fmt.Errorf("Field Pan is len...")
	}

	return nil
}
