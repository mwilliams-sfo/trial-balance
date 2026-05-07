package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strconv"
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

func findAttr(start xml.StartElement, name xml.Name) (int, error) {
	for i, attr := range start.Attr {
		if attr.Name == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Attribute %s:%s not found", name.Space, name.Local)
}

func parseTextElement(d *xml.Decoder, start xml.StartElement) (string, error) {
	var text string
	tok, err := d.RawToken()
	if err != nil {
		return "", fmt.Errorf("parseTextElement: %w", err)
	}
	if cd, ok := tok.(xml.CharData); ok {
		text = string(cd)
	} else {
		return "", fmt.Errorf("Invalid token, expected character data")
	}
	err = skipContent(d, start)
	if err != nil {
		return "", fmt.Errorf("parseTextElement: %w", err)
	}
	return text, nil
}

func parseCountData(d *xml.Decoder, start xml.StartElement) (string, int, error) {
	i, err := findAttr(start, names.CdTypeAttr())
	if err != nil {
		return "", 0, fmt.Errorf("parseCountData: %w", err)
	}
	typeValue := start.Attr[i].Value
	text, err := parseTextElement(d, start)
	if err != nil {
		return "", 0, fmt.Errorf("parseCountData: %w", err)
	}
	count, err := strconv.ParseInt(text, 10, 0)
	if err != nil {
		return "", 0, fmt.Errorf("parseCountData: %w", err)
	}
	return typeValue, int(count), nil
}

func isGuidChar(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f'
}

func isValidGuid(s string) bool {
	return len(s) == 32 && strings.IndexFunc(s, func(r rune) bool { return !isGuidChar(r) }) == -1
}

func parseGuid(d *xml.Decoder, start xml.StartElement) (model.Guid, error) {
	var zero model.Guid
	i, err := findAttr(start, names.IdTypeAttr())
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	if start.Attr[i].Value != "guid" {
		return zero, fmt.Errorf("parseGuid: invaid type")
	}
	text, err := parseTextElement(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseGuid: %w", err)
	}
	if !isValidGuid(text) {
		return zero, fmt.Errorf("parseGuid: invalid guid")
	}
	return model.Guid(text), nil
}
