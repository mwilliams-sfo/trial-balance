package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strings"
	"trial_balance/gnucash/model"
)

func requireStartElement(d *xml.Decoder) (xml.StartElement, error) {
	var start xml.StartElement
	tok, err := d.RawToken()
	if err != nil {
		return start, err
	}
	var ok bool
	if start, ok = tok.(xml.StartElement); !ok {
		return start, fmt.Errorf("Invalid token, expected start tag")
	}
	return start, nil
}

func requireCharData(d *xml.Decoder) (xml.CharData, error) {
	var cd xml.CharData
	tok, err := d.RawToken()
	if err != nil {
		return cd, err
	}
	var ok bool
	if cd, ok = tok.(xml.CharData); !ok {
		return cd, fmt.Errorf("Invalid token, expected character data")
	}
	return cd, nil
}

func requireEndElement(d *xml.Decoder) (xml.EndElement, error) {
	var end xml.EndElement
	tok, err := d.RawToken()
	if err != nil {
		return end, err
	}
	var ok bool
	if end, ok = tok.(xml.EndElement); !ok {
		return end, fmt.Errorf("Invalid token, expected end tag")
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

func skipContent(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.RawToken()
		if err != nil {
			return err
		}
		switch x := tok.(type) {
		case xml.StartElement:
			err = skipContent(d, x)
			if err != nil {
				return err
			}
		case xml.EndElement:
			if x.Name != start.Name {
				return fmt.Errorf("Mismatched end tag: </%s:%s>", x.Name.Space, x.Name.Local)
			}
			return nil
		}
	}
}

func findChildElement(d *xml.Decoder, parent xml.StartElement, name xml.Name) (xml.Token, error) {
	var zero xml.Token
	for {
		tok, err := d.RawToken()
		if err != nil {
			return zero, err
		}
		switch x := tok.(type) {
		case xml.StartElement:
			if x.Name == name {
				return x, nil
			}
			err = skipContent(d, x)
			if err != nil {
				return zero, err
			}
		case xml.EndElement:
			if x.Name != parent.Name {
				return zero, fmt.Errorf("Mismatched end tag: </%s:%s>", x.Name.Space, x.Name.Local)
			}
			return x, nil
		}
	}
}
