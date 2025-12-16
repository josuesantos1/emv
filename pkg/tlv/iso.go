package tlv 

type Tag struct {
	Name string
	Length int64
}

const IsoMessage = []Tag{
	Tag{
		Name: "Pan",
		Length: 2,
	},
	Tag{
		Name: "Data de Validade",
		Length: 4,
	},
	Tag{
		Name: "CVM",
		Length: 4,
	},
}
