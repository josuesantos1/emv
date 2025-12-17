package tlv

import (
	pkgtlv "github.com/josuesantos1/emv/pkg/tlv"
	"testing"
	"time"
)

func TestTlv_Populate(t *testing.T) {
	tests := []struct {
		name    string
		tlvs    []pkgtlv.TLV
		want    Tlv
		wantErr bool
	}{
		{
			name: "populate all fields successfully",
			tlvs: []pkgtlv.TLV{
				{Tag: "5A", Value: "1234567890123456"},
				{Tag: "5F24", Value: "251231"},
				{Tag: "9F34", Value: "1F0000"},
			},
			want: Tlv{
				Pan:          "1234567890123456",
				DataValidade: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				CVM:          "1F0000",
			},
			wantErr: false,
		},
		{
			name: "populate with partial data",
			tlvs: []pkgtlv.TLV{
				{Tag: "5A", Value: "9876543210987654"},
			},
			want: Tlv{
				Pan: "9876543210987654",
			},
			wantErr: false,
		},
		{
			name: "populate with invalid date format",
			tlvs: []pkgtlv.TLV{
				{Tag: "5F24", Value: "9999"},
			},
			want:    Tlv{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlv := &Tlv{}
			err := tlv.Populate(tt.tlvs)

			if (err != nil) != tt.wantErr {
				t.Errorf("Tlv.Populate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if tlv.Pan != tt.want.Pan {
					t.Errorf("Tlv.Populate() Pan = %v, want %v", tlv.Pan, tt.want.Pan)
				}
				if !tt.want.DataValidade.IsZero() && !tlv.DataValidade.Equal(tt.want.DataValidade) {
					t.Errorf("Tlv.Populate() DataValidade = %v, want %v", tlv.DataValidade, tt.want.DataValidade)
				}
				if tlv.CVM != tt.want.CVM {
					t.Errorf("Tlv.Populate() CVM = %v, want %v", tlv.CVM, tt.want.CVM)
				}
			}
		})
	}
}

func TestTlv_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tlv     Tlv
		wantErr bool
	}{
		{
			name: "valid Tlv with all fields",
			tlv: Tlv{
				Pan:          "4539578763621486",
				DataValidade: time.Now().AddDate(1, 0, 0),
				CVM:          "1F0000",
			},
			wantErr: false,
		},
		{
			name: "invalid Tlv - Pan too short",
			tlv: Tlv{
				Pan:          "123456789012",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: true,
		},
		{
			name: "invalid Tlv - Pan too long",
			tlv: Tlv{
				Pan:          "12345678901234567890",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: true,
		},
		{
			name: "invalid Tlv - missing Pan",
			tlv: Tlv{
				Pan:          "",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: true,
		},
		{
			name: "invalid Tlv - missing CVM",
			tlv: Tlv{
				Pan:          "1234567890123456",
				DataValidade: time.Now(),
				CVM:          "",
			},
			wantErr: true,
		},
		{
			name: "valid Tlv - Pan with 14 digits",
			tlv: Tlv{
				Pan:          "4539578763621486",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: false,
		},
		{
			name: "valid Tlv - Pan with 18 digits",
			tlv: Tlv{
				Pan:          "4539578763621486",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: false,
		},
		{
			name: "invalid Tlv - Pan is not valid",
			tlv: Tlv{
				Pan:          "4539578763621487",
				DataValidade: time.Now(),
				CVM:          "1F0000",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tlv.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tlv.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
