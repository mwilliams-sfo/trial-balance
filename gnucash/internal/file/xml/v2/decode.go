package v2

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	xmlutil "trial_balance/gnucash/internal/file/xml/util"
)

var guidRe = regexp.MustCompile("\\A[0-9a-f]{32}\\z")

type Decoder struct {
	XmlDecoder *xml.Decoder
	Token      xml.Token
}

func (d *Decoder) token() error {
	var err error
	d.Token, err = d.XmlDecoder.Token()
	return err
}

func (d *Decoder) skipSpace() error {
	err := d.token()
	if err != nil {
		return fmt.Errorf("skipSpace: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if ok && xmlutil.IsXmlSpace(string(cd)) {
		err := d.token()
		if err != nil {
			return fmt.Errorf("skipSpace: %w", err)
		}
	}
	return nil
}

func (d *Decoder) skipElement() error {
	d.Token = nil
	return d.XmlDecoder.Skip()
}

func (d *Decoder) decodeGuid() (string, error) {
	start := d.Token.(xml.StartElement)
	var idType string
	idType, ok := xmlutil.GetAttr(start, xml.Name{Local: "type"})
	if !ok {
		return "", fmt.Errorf("decodeGuid: Expected type attribute")
	}
	if idType != "guid" {
		return "", fmt.Errorf("decodeGuid: Expected GUID")
	}
	err := d.token()
	if err != nil {
		return "", fmt.Errorf("decodeGuid: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if !ok {
		return "", fmt.Errorf("decodeGuid: Expected GUID text")
	}
	guid := string(cd)
	if !guidRe.MatchString(guid) {
		return "", fmt.Errorf("decodeGuid: Invalid GUID")
	}
	err = d.skipElement()
	if err != nil {
		return "", fmt.Errorf("decodeGuid: %w", err)
	}
	return guid, nil
}

func (d *Decoder) decodeCount() (string, int, error) {
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsGnc && start.Name.Local == "count-data") {
		return "", 0, fmt.Errorf("decodeCount: Expected count-data")
	}
	countType, ok := xmlutil.GetAttr(start, xml.Name{Space: nsCd, Local: "type"})
	if !ok {
		return "", 0, fmt.Errorf("decodeCount: Expected type attribute")
	}
	err := d.token()
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if !ok {
		return "", 0, fmt.Errorf("decodeCount: Expected count value")
	}
	count, err := strconv.Atoi(string(cd))
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: %w", err)
	}
	err = d.skipElement()
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: Invalid count value")
	}
	return countType, count, nil
}

func (d *Decoder) decodeAccount() (Account, error) {
	var zero, account Account

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeAccount: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsAct && start.Name.Local == "name") {
		return zero, fmt.Errorf("decodeAccount: Expected name")
	}
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

	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeAccount: %w", err)
	}
	start, ok = d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsAct && start.Name.Local == "id") {
		return zero, fmt.Errorf("decodeAccount: Expected id")
	}
	account.Id, err = d.decodeGuid()
	if err != nil {
		return zero, fmt.Errorf("decodeAccount: %w", err)
	}

	for {
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeAccount: %w", err)
		}
		start, ok = d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsAct && (start.Name.Local == "type" ||
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

func (d *Decoder) decodeBook() (Book, error) {
	var zero, book Book

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeBook: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsBook && start.Name.Local == "id") {
		return zero, fmt.Errorf("decodeBook: Expected id")
	}
	book.Id, err = d.decodeGuid()
	if err != nil {
		return zero, fmt.Errorf("decodeBook: %w", err)
	}
	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeBook: %w", err)
	}

	for {
		start, ok := d.Token.(xml.StartElement)
		if !(ok && start.Name.Space == nsBook && start.Name.Local == "slots") {
			break
		}

		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	}

	for {
		start, ok := d.Token.(xml.StartElement)
		if !(ok && start.Name.Space == nsGnc && start.Name.Local == "count-data") {
			break
		}
		countType, count, err := d.decodeCount()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
		switch countType {
		case "account":
			book.Accounts = make([]Account, count)
		default:
		}
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	}

	for {
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsGnc &&
			(start.Name.Local == "commodity" || start.Name.Local == "pricedb") {

			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
			err = d.skipSpace()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
		} else {
			break
		}
	}

	for i := range len(book.Accounts) {
		start, ok := d.Token.(xml.StartElement)
		if !(ok && start.Name.Space == nsGnc && start.Name.Local == "account") {
			return zero, fmt.Errorf("decodeBook: Expected account")
		}
		book.Accounts[i], err = d.decodeAccount()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	}
	if start, ok := d.Token.(xml.StartElement); ok && start.Name.Space == nsGnc && start.Name.Local == "account" {
		return zero, fmt.Errorf("decodeBook: Unexpected account")
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsGnc && end.Name.Local == "book") {
		err := d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	}

	return book, nil
}

func (d *Decoder) DecodeFile() (Book, error) {
	var zero Book

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("DecodeFile: %w", err)
	}
	if _, ok := d.Token.(xml.EndElement); ok {
		return zero, fmt.Errorf("DecodeFile: Expected count-data")
	}
	countType, count, err := d.decodeCount()
	if err != nil {
		return zero, fmt.Errorf("DecodeFile: %w", err)
	}
	if countType != "book" {
		return zero, fmt.Errorf("DecodeFile: Expected book count")
	}
	if count != 1 {
		return zero, fmt.Errorf("DecodeFile: Invalid book count")
	}

	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("DecodeFile: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsGnc && start.Name.Local == "book") {
		return zero, fmt.Errorf("DecodeFile: Expected book element")
	}
	book, err := d.decodeBook()
	if err != nil {
		return zero, fmt.Errorf("DecodeFile: %w", err)
	}

	err = d.skipElement()
	if err != nil {
		return zero, fmt.Errorf("DecodeFile: %w", err)
	}
	return book, nil
}
