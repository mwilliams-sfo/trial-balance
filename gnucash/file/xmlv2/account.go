package xmlv2

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xmlv2/names"
	"trial_balance/gnucash/model"
)

func parseAccountTextChild(d *xml.Decoder, start xml.StartElement) (string, error) {
	cd, err := requireCharData(d)
	if err != nil {
		return "", err
	}
	end, err := requireEndElement(d)
	if err != nil {
		return "", err
	}
	if end.Name != start.End().Name {
		return "", fmt.Errorf("parseAccountTextChild: invalid token, expected </%s:%s>", start.Name.Space, start.Name.Local)
	}
	return string(cd), nil
}

func parseAccountGuid(d *xml.Decoder, start xml.StartElement) (model.Guid, error) {
	var zero model.Guid
	for _, attr := range start.Attr {
		if attr.Name.Local == "guid" {
			id := model.Guid(attr.Value)
			err := skipContent(d, start)
			if err != nil {
				return zero, err
			}
			return id, nil
		}
	}
	return zero, fmt.Errorf("<%s:%s>: guid attribute not found", start.Name.Space, start.Name.Local)
}

func parseAccount(d *xml.Decoder, start xml.StartElement) (model.Account, error) {
	var zero, a model.Account
	if start.Name != names.Account() {
		return zero, fmt.Errorf("parseAccount: invalid start tag")
	}

	tok, err := findChildElement(d, start, names.AccountName())
	if err != nil {
		return zero, err
	}
	switch x := tok.(type) {
	case xml.StartElement:
		text, err := parseAccountTextChild(d, x)
		if err != nil {
			return zero, err
		}
		a.Name = text
	default:
		return zero, fmt.Errorf("parseAccount: name not found")
	}

	tok, err = findChildElement(d, start, names.AccountId())
	if err != nil {
		return zero, err
	}
	switch x := tok.(type) {
	case xml.StartElement:

		id, err := parseAccountGuid(d, x)
		if err != nil {
			return zero, err
		}
		a.Id = id
	default:
		return zero, fmt.Errorf("parseAccount: id not found")
	}

	tok, err = findChildElement(d, start, names.AccountParent())
	if err != nil {
		return zero, err
	}
	switch x := tok.(type) {
	case xml.StartElement:
		id, err := parseAccountGuid(d, x)
		if err != nil {
			return zero, err
		}
		a.Parent = id
	default:
		return zero, fmt.Errorf("parseAccount: parent not found")
	}

	err = skipContent(d, start)
	if err != nil {
		return zero, err
	}
	return a, nil
}
