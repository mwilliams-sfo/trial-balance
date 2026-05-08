package xml

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xml/model"
	"trial_balance/gnucash/file/xml/names"
)

func parseAccount(d *xml.Decoder, start xml.StartElement) (model.Account, error) {
	var zero, a model.Account
	if start.Name != names.Account() {
		return zero, fmt.Errorf("parseAccount: invalid start tag")
	}

	tok, err := findChildElement(d, start, names.AccountName())
	if err != nil {
		return zero, err
	}
	if nameStart, ok := tok.(xml.StartElement); ok {
		text, err := parseTextElement(d, nameStart)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
		a.Name = text
	} else {
		return zero, fmt.Errorf("parseAccount: name not found")
	}

	tok, err = findChildElement(d, start, names.AccountId())
	if err != nil {
		return zero, err
	}
	if idStart, ok := tok.(xml.StartElement); ok {
		guid, err := parseGuid(d, idStart)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
		a.Id = guid
	} else {
		return zero, fmt.Errorf("parseAccount: id not found")
	}

	tok, err = findChildElement(d, start, names.AccountParent())
	if err != nil {
		return zero, err
	}
	if parentStart, ok := tok.(xml.StartElement); ok {
		guid, err := parseGuid(d, parentStart)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
		a.Parent = guid
	} else {
		return zero, fmt.Errorf("parseAccount: parent not found")
	}

	err = skipContent(d, start)
	if err != nil {
		return zero, err
	}
	return a, nil
}
