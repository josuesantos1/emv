package tlv

import (
	"bytes"
	"encoding/hex"
	"testing"
	"strings"
)

func TestParser_Parse(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    []TLV
		wantErr bool
	}{
		{
			name:  "parse valid TLV data with PAN, expiry date and CVM",
			input: "5A0812345678901234565F2404251200009F340442000000",
			want: []TLV{
				{Tag: hexToBytes("5A"), Len: 8, Value: hexToBytes("1234567890123456")},
				{Tag: hexToBytes("5F24"), Len: 4, Value: hexToBytes("25120000")},
				{Tag: hexToBytes("9F34"), Len: 4, Value: hexToBytes("42000000")},
			},
			wantErr: false,
		},
		{
			name:  "parse single TLV",
			input: "5A084539578763621486",
			want: []TLV{
				{Tag: hexToBytes("5A"), Len: 8, Value: hexToBytes("4539578763621486")},
			},
			wantErr: false,
		},
		{
			name:    "empty data",
			input:   "",
			want:    []TLV{},
			wantErr: false,
		},
		{
			name:    "incomplete TLV - missing value",
			input:   "5A08",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "incomplete TLV - missing length and value",
			input:   "5A",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{}
			data, _ := hex.DecodeString(tt.input)

			got, err := parser.Parse(data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("Parser.Parse() returned %d TLVs, want %d", len(got), len(tt.want))
					return
				}

				for i := range got {
					if !bytes.Equal(got[i].Tag, tt.want[i].Tag) {
						t.Errorf("Parser.Parse() TLV[%d].Tag = %X, want %X", i, got[i].Tag, tt.want[i].Tag)
					}
					if got[i].Len != tt.want[i].Len {
						t.Errorf("Parser.Parse() TLV[%d].Len = %v, want %v", i, got[i].Len, tt.want[i].Len)
					}
					if !bytes.Equal(got[i].Value, tt.want[i].Value) {
						t.Errorf("Parser.Parse() TLV[%d].Value = %X, want %X", i, got[i].Value, tt.want[i].Value)
					}
				}
			}
		})
	}
}

func TestParser_ParseTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTag  string
		wantUsed int
		wantErr  bool
	}{
		{
			name:     "single byte tag",
			input:    "5A08",
			wantTag:  "5A",
			wantUsed: 1,
			wantErr:  false,
		},
		{
			name:     "two byte tag",
			input:    "9F3404",
			wantTag:  "9F34",
			wantUsed: 2,
			wantErr:  false,
		},
		{
			name:     "three byte tag",
			input:    "5F2404",
			wantTag:  "5F24",
			wantUsed: 2,
			wantErr:  false,
		},
		{
			name:    "empty buffer",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{}
			data, _ := hex.DecodeString(tt.input)

			tag, used, err := parser.ParseTag(data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotTag := strings.ToUpper(hex.EncodeToString(tag))
				if gotTag != tt.wantTag {
					t.Errorf("Parser.ParseTag() tag = %v, want %v", gotTag, tt.wantTag)
				}
				if used != tt.wantUsed {
					t.Errorf("Parser.ParseTag() used = %v, want %v", used, tt.wantUsed)
				}
			}
		})
	}
}

func TestParser_ParseLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantUsed int
		wantErr  bool
	}{
		{
			name:     "single byte length",
			input:    "08",
			wantLen:  8,
			wantUsed: 1,
			wantErr:  false,
		},
		{
			name:     "length 4",
			input:    "04",
			wantLen:  4,
			wantUsed: 1,
			wantErr:  false,
		},
		{
			name:    "empty buffer",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{}
			data, _ := hex.DecodeString(tt.input)

			length, used, err := parser.ParseLength(data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if length != tt.wantLen {
					t.Errorf("Parser.ParseLength() length = %v, want %v", length, tt.wantLen)
				}
				if used != tt.wantUsed {
					t.Errorf("Parser.ParseLength() used = %v, want %v", used, tt.wantUsed)
				}
			}
		})
	}
}

func hexToBytes(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
