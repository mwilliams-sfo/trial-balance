package v2

import (
	"encoding/xml"
	"fmt"
)

func (d *decoder) decodeAccount() (Account, error) {
	var zero, account Account

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeAccount: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsAct && start.Name.Local == "name" {
		err = d.token()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
		cd, ok := d.Token.(xml.CharData)
		if !ok {
			return zero, fmt.Errorf("decodeAccount: Expected name text")
		}
		account.Name = string(cd)
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeAccount: Expected name")
	}

	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeAccount: %w", err)
	}
	start, ok = d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsAct && start.Name.Local == "id" {
		account.Id, err = d.decodeGuid()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeAccount: Expected id")
	}

	for {
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
		start, ok = d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsAct &&
			(start.Name.Local == "type" ||
				start.Name.Local == "commodity" ||
				start.Name.Local == "commodity-scu" ||
				start.Name.Local == "non-standard-scu" ||
				start.Name.Local == "code" ||
				start.Name.Local == "description" ||
				start.Name.Local == "slots") {
			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeAccount: %w", err)
			}
		} else {
			break
		}
	}

	start, ok = d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsAct && start.Name.Local == "parent" {
		account.Parent, err = d.decodeGuid()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsGnc && end.Name.Local == "account") {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
	}
	return account, nil
}
