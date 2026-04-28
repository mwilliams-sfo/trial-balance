package xmlv2

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/model"
)

var (
	accountName       = xml.Name{"gnc", "account"}
	accountNameName   = xml.Name{"act", "name"}
	accountIdName     = xml.Name{"act", "id"}
	accountParentName = xml.Name{"act", "parent"}
)

func parseAccountName(d *xml.Decoder) (string, error) {
	start, err := requireStartElement(d)
	if err != nil {
		return "", err
	}
	if start.Name != accountNameName {
		return "", fmt.Errorf("parseAccountName: invalid start tag")
	}
	cd, err := requireCharData(d)
	if err != nil {
		return "", err
	}
	end, err := requireEndElement(d)
	if err != nil {
		return "", err
	}
	if end.Name != start.End().Name {
		return "", fmt.Errorf("parseAccountName: invalid token, expected </act:name>")
	}
	return string(cd), nil
}

func parseAccount(d *xml.Decoder, start xml.StartElement) (model.Account, error) {
	var zero model.Account
	if start.Name != accountName {
		return zero, fmt.Errorf("parseAccount: invalid start tag")
	}

	name, err := parseAccountName(d)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: invalid name: %w", err)
	}
	a := model.Account{}
	a.Name = name

	childStart, err := requireStartElement(d)
	if err != nil {
		return zero, err
	}
	if childStart.Name != accountIdName {
		return zero, fmt.Errorf("parseAccount: invalid token, expected <act:id>")
	}
	a.Id, err = requireGuid(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: invalid id: %w", err)
	}

	depth := 1
parentLoop:
	for {
		tok, err := d.RawToken()
		if err != nil {
			return zero, err
		}
		switch x := tok.(type) {
		case xml.StartElement:
			depth++
			if x.Name == accountParentName {
				childStart = x
				break parentLoop
			}
		case xml.EndElement:
			depth--
			if depth == 0 {
				return zero, fmt.Errorf("parseAccount: parent id not found")
			}
		}
	}
	a.Parent, err = requireGuid(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: invalid parent id: %w", err)
	}

	var end xml.EndElement
	depth = 1
endLoop:
	for {
		tok, err := d.RawToken()
		if err != nil {
			return zero, err
		}
		switch x := tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
			if depth == 0 {
				end = x
				break endLoop
			}
		}
	}
	if end != start.End() {
		return zero, fmt.Errorf("parseAccount: invalid token, expected </gnc:account>")
	}
	return a, nil
}
