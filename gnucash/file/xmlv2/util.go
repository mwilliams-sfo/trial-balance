package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strings"
	"trial_balance/gnucash/file/xmlv2/names"
	"trial_balance/gnucash/model"
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
attrLoop:
	for i := 0; ; i++ {
		switch {
		case i == len(start.Attr):
			return zero, fmt.Errorf("parseGuid: guid attribute not found")
		case start.Attr[i].Name == names.GuidAttr():
			if !isValidGuid(start.Attr[i].Value) {
				return zero, fmt.Errorf("parseGuid: invalid guid")
			}
			guid = model.Guid(start.Attr[i].Value)
			break attrLoop
		}
	}
	err := skipContent(d, start)
	if err != nil {
		return zero, err
	}
	return guid, nil
}
