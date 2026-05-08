package xml

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xml/model"
	"trial_balance/gnucash/file/xml/names"
)

func formattedName(name xml.Name) string {
	s := ""
	if name.Space != "" {
		s += name.Space + ":"
	}
	s += name.Local
	return s
}

func skipSpace(d *xml.Decoder) (xml.Token, error) {
	var zero xml.Token
	for {
		tok, err := d.RawToken()
		if err != nil {
			return zero, err
		}
		cd, ok := tok.(xml.CharData)
		if !ok {
			return tok, nil
		}
		for _, r := range string(cd) {
			if !(r == '\t' || r == '\n' || r == '\r' || r == ' ') {
				return tok, nil
			}
		}
	}
}

func skipContent(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := skipSpace(d)
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

func findAttr(start xml.StartElement, name xml.Name) (int, error) {
	for i, attr := range start.Attr {
		if attr.Name == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Attribute %s not found", formattedName(name))
}

func findChildElement(d *xml.Decoder, parent xml.StartElement, name xml.Name) (xml.Token, error) {
	var zero xml.Token
	for {
		tok, err := skipSpace(d)
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
	tok, err := d.RawToken()
	if err != nil {
		return "", err
	}
	cd, ok := tok.(xml.CharData)
	if !ok {
		return "", fmt.Errorf("Invalid token, expected character data")
	}
	text := string(cd)

	tok, err = d.RawToken()
	if err != nil {
		return "", err
	}
	end, ok := tok.(xml.EndElement)
	if !ok || end.Name != start.Name {
		return "", fmt.Errorf("Invalid token, expected </%s>", formattedName(start.Name))
	}
	return text, nil
}

func validateGuid(s string) error {
	if len(s) != 32 {
		return fmt.Errorf("validateGuid: invalid length")
	}
	for _, r := range s {
		if !('0' <= r && r <= '9' || 'a' <= r && r <= 'f') {
			return fmt.Errorf("validateGuid: invalid character")
		}
	}
	return nil
}

func parseGuid(d *xml.Decoder, start xml.StartElement) (model.Guid, error) {
	var zero model.Guid
	i, err := findAttr(start, names.IdTypeAttr())
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	idType := start.Attr[i].Value
	if idType != "guid" {
		return zero, fmt.Errorf("parseGuid: unrecognized id type %s", idType)
	}
	text, err := parseTextElement(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	err = validateGuid(text)
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	return model.Guid(text), nil
}
