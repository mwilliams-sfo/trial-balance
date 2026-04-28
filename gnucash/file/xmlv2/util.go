package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strings"
	"trial_balance/gnucash/model"
)

func requireStartElement(d *xml.Decoder) (xml.StartElement, error) {
	zero := xml.StartElement{}
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	start, ok := tok.(xml.StartElement)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected start tag")
	}
	return start, nil
}

func requireCharData(d *xml.Decoder) (xml.CharData, error) {
	zero := xml.CharData("")
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	cd, ok := tok.(xml.CharData)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected character data")
	}
	return cd, nil
}

func requireEndElement(d *xml.Decoder) (xml.EndElement, error) {
	zero := xml.EndElement{}
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	end, ok := tok.(xml.EndElement)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected end tag")
	}
	return end, nil
}

func isGuidChar(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f'
}

func isValidGuid(s string) bool {
	return len(s) == 32 && strings.IndexFunc(s, func(r rune) bool { return !isGuidChar(r) }) == -1
}

func requireGuid(d *xml.Decoder, start xml.StartElement) (model.Guid, error) {
	var zero model.Guid
	for _, attr := range start.Attr {
		if attr.Name.Local == "guid" && isValidGuid(attr.Value) {
			end, err := requireEndElement(d)
			if err != nil {
				return zero, err
			}
			if end != start.End() {
				return zero, fmt.Errorf("requireGuid: invalid token, expected </%s:%s>", start.Name.Space, start.Name.Local)
			}
			return model.Guid(attr.Value), nil
		}
	}
	return zero, fmt.Errorf("guid attribute not found")
}
