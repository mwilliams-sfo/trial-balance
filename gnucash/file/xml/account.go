package xml

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xml/model"
	"trial_balance/gnucash/file/xml/names"
)

func parseAccount(d *xml.Decoder, start xml.StartElement) (model.Account, error) {
	var zero model.Account
	if start.Name != names.Account() {
		return zero, fmt.Errorf("parseAccount: invalid start tag")
	}
	var account model.Account

	tok, err := skipSpace(d)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: %w", err)
	}
	childStart, ok := tok.(xml.StartElement)
	if !ok || childStart.Name != names.AccountName() {
		return zero, fmt.Errorf("parseAccount: invalid token, expected <%s>", formattedName(names.AccountName()))
	}
	account.Name, err = parseTextElement(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: %w", err)
	}

	tok, err = skipSpace(d)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: %w", err)
	}
	childStart, ok = tok.(xml.StartElement)
	if !ok || childStart.Name != names.AccountId() {
		return zero, fmt.Errorf("parseAccount: invalid token, expected <%s>", formattedName(names.AccountId()))
	}
	account.Id, err = parseGuid(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseAccount: %w", err)
	}

	for {
		tok, err = skipSpace(d)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
		childStart, ok = tok.(xml.StartElement)
		if !ok {
			break
		}
		if childStart.Name == names.AccountParent() {
			account.Parent, err = parseGuid(d, childStart)
			if err != nil {
				return zero, fmt.Errorf("parseAccount: %w", err)
			}
			account.HasParent = true
			tok, err = skipSpace(d)
			if err != nil {
				return zero, fmt.Errorf("parseAccount: %w", err)
			}
			break
		}
		err = skipContent(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
	}

	if childStart, ok = tok.(xml.StartElement); ok {
		err = skipContent(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
	}
	if end, ok := tok.(xml.EndElement); ok {
		if end.Name != start.Name {
			return zero, fmt.Errorf("parseAccount: invalid token, expected </%s>", formattedName(start.Name))
		}
	} else {
		err = skipContent(d, start)
		if err != nil {
			return zero, fmt.Errorf("parseAccount: %w", err)
		}
	}
	return account, nil
}
