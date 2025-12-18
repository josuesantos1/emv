package domain

import (
	"encoding/hex"
	"fmt"
	"github.com/josuesantos1/emv/pkg/tlv"
	"time"
)

type Tlv struct {
	Pan          string
	DataValidade time.Time
	CVM          string
}

const (
	Pan     = "5A"
	ExpDate = "5F24"
	CVM     = "9F34"
)

func (t *Tlv) Populate(tlvs []tlv.TLV) error {
	for _, tlvItem := range tlvs {
		tagHex := tlvItem.TagHex()
		valueHex := tlvItem.ValueHex()

		if tagHex == Pan {
			t.Pan = valueHex
		}

		if tagHex == ExpDate {
			if len(valueHex) < 4 {
				return fmt.Errorf("invalid expiry date format: value too short")
			}
			yy := valueHex[0:2]
			mm := valueHex[2:4]
			dd := "01"
			year := "20" + yy

			dateStr := fmt.Sprintf("%s-%s-%s", year, mm, dd)

			parsedTime, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("failed to parse expiry date '%s': %w", dateStr, err)
			}

			t.DataValidade = parsedTime
		}

		if tagHex == CVM {
			t.CVM = valueHex
		}
	}

	return nil
}

func (t *Tlv) Validate() error {
	if t.Pan == "" {
		return fmt.Errorf("field Pan is required")
	}
	if t.DataValidade.IsZero() {
		return fmt.Errorf("field Data de validade is required")
	}
	if t.CVM == "" {
		return fmt.Errorf("field CVM is required")
	}

	if len(t.Pan) < 13 || len(t.Pan) > 19 {
		return fmt.Errorf("PAN must be between 13 and 19 digits")
	}

	if !t.ValidatePan() {
		return fmt.Errorf("PAN failed Luhn algorithm validation")
	}

	now := time.Now()

	year := t.DataValidade.Year()
	month := t.DataValidade.Month()

	if year < now.Year() || (year == now.Year() && month < now.Month()) {
		return fmt.Errorf("card expired: expiry date %s is before current date %s", t.DataValidade.Format("01/2006"), now.Format("01/2006"))
	}

	if err := t.ValidateCVM(); err != nil {
		return err
	}

	return nil
}

func (t *Tlv) ValidatePan() bool {
	sum := 0
	alt := false

	for i := len(t.Pan) - 1; i >= 0; i-- {
		if t.Pan[i] < '0' || t.Pan[i] > '9' {
			return false
		}

		digit := int(t.Pan[i] - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}

	return sum%10 == 0
}

var bit01Values = map[string]struct{}{
	"00": {},
	"01": {},
	"02": {},
	"03": {},
	"04": {},
	"05": {},
	"06": {},
	"07": {},
	"1D": {},
	"1E": {},
	"1F": {},
	"20": {},
	"FF": {},
}

var bit02Values = map[string]struct{}{
	"00": {},
	"01": {},
	"02": {},
	"03": {},
	"04": {},
	"05": {},
	"06": {},
	"07": {},
	"08": {},
	"09": {},
	"FF": {},
}

var bit03Values = map[string]struct{}{
	"00": {},
	"01": {},
	"02": {},
	"03": {},
	"04": {},
	"05": {},
	"FF": {},
}

func (t *Tlv) ValidateCVM() error {
	if len(t.CVM) < 6 {
		return fmt.Errorf("CVM must be at least 6 characters, got %d", len(t.CVM))
	}

	bit01 := t.CVM[:2]
	bit02 := t.CVM[2:4]
	bit03 := t.CVM[4:6]

	if _, exists := bit01Values[bit01]; !exists {
		return fmt.Errorf("invalid CVM bit 1 value '%s': not a supported method", bit01)
	}

	if _, exists := bit02Values[bit02]; !exists {
		return fmt.Errorf("invalid CVM bit 2 value '%s': not a supported condition", bit02)
	}

	if _, exists := bit03Values[bit03]; !exists {
		return fmt.Errorf("invalid CVM bit 3 value '%s': not a supported value", bit03)
	}

	return nil
}

func (t *Tlv) ValueHex(value []byte) string {
	return hex.EncodeToString(value)
}
