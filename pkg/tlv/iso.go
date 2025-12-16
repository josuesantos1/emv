package tlv

type Tag struct {
	Name   string
	Length int
	Tag    string
}

var IsoMessage = []Tag{
	{
		Name:   "Pan",
		Length: 2,
		Tag:    "5A",
	},
	{
		Name:   "Data de Validade",
		Length: 4,
		Tag:    "5F24",
	},
	{
		Name:   "CVM",
		Length: 4,
		Tag:    "9F34",
	},
}
