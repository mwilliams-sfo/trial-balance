package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strings"
	"trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/names"
)

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

func parseTextElement(d *xml.Decoder, start xml.StartElement) (string, error) {
	var text string
	tok, err := d.RawToken()
	if err != nil {
		return "", err
	}
	if cd, ok := tok.(xml.CharData); ok {
		text = string(cd)
	} else {
		return "", fmt.Errorf("Invalid token, expected character data")
	}
	err = skipContent(d, start)
	if err != nil {
		return "", err
	}
	return text, nil
}

func isGuidChar(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f'
}

func isValidGuid(s string) bool {
	return len(s) == 32 && strings.IndexFunc(s, func(r rune) bool { return !isGuidChar(r) }) == -1
}

func parseGuid(d *xml.Decoder, start xml.StartElement) (model.Guid, error) {
	var zero, guid model.Guid
	var idType string
typeLoop:
	for i := 0; ; i++ {
		if i == len(start.Attr) {
			return zero, fmt.Errorf("parseGuid: type attribute not found")
		}
		if start.Attr[i].Name == names.IdTypeAttr() {
			idType = start.Attr[i].Value
			break typeLoop
		}
	}
	if idType != "guid" {
		return zero, fmt.Errorf("parseGuid: invaid type")
	}

	tok, err := d.RawToken()
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	if cd, ok := tok.(xml.CharData); ok {
		if !isValidGuid(string(cd)) {
			return zero, fmt.Errorf("parseGuid: invalid guid")
		}
		guid = model.Guid(string(cd))
	} else {
		return zero, fmt.Errorf("parseGuid: invalid token, expected char data")
	}

	tok, err = d.RawToken()
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	if end, ok := tok.(xml.EndElement); !ok || end.Name != start.Name {
		return zero, fmt.Errorf("parseGuid: invalid token, expected </%s:%s>", start.Name.Space, start.Name.Local)
	}

	return guid, nil
}
