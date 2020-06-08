package markup

import "encoding/xml"

type node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []node     `xml:",any"`
	Body    string     `xml:",chardata"`
}

func (n *node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type nn node

	return d.DecodeElement((*nn)(n), &start)
}
