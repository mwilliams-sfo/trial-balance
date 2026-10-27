package gnucash

import (
	"encoding/xml"
	xmlv2 "trial_balance/gnucash/internal/file/xml/v2"
)

func decodeXmlV2(d *xml.Decoder, start xml.StartElement) (*Book, error) {
	v2Decoder := &xmlv2.Decoder{XmlDecoder: d}
	v2Book, err := v2Decoder.DecodeFile()
	if err != nil {
		return nil, err
	}
	return &Book{Id: v2Book.Id}, nil
}
