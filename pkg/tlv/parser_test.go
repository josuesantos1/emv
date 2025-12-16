package tlv

import (
	"testing"
	"fmt"
)

func TestTLVParse(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    TLV
		wantErr bool
	}{
		{
			name: "parsing válido",
			data: "5A08123456785F240425129F34044200",
			want: TLV{
				Pan:          "12345678",
				DataValidade: "2512",
				CVM:          "4200",
			},
			wantErr: false,
		},
		// {
		// 	name:    "dados não hexadecimais",
		// 	data:    "ZZZZZ",
		// 	wantErr: true,
		// },
		// {
		// 	name:    "dados incompletos",
		// 	data:    "5A",
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlv := &TLV{}
			fmt.Println(tt.data)
			err := tlv.Parse(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLV.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tlv.Pan != tt.want.Pan {
					t.Errorf("TLV.Parse() Pan = %v, want %v", tlv.Pan, tt.want.Pan)
				}
				if tlv.DataValidade != tt.want.DataValidade {
					t.Errorf("TLV.Parse() DataValidade = %v, want %v", tlv.DataValidade, tt.want.DataValidade)
				}
				if tlv.CVM != tt.want.CVM {
					t.Errorf("TLV.Parse() CVM = %v, want %v", tlv.CVM, tt.want.CVM)
				}
			}
		})
	}
}
