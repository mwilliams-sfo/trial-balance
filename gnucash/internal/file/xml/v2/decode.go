package v2

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"time"
	xmlutil "trial_balance/gnucash/internal/file/xml/util"
)

var guidRe = regexp.MustCompile("\\A[0-9a-f]{32}\\z")

type decoder struct {
	XmlDecoder *xml.Decoder
	Token      xml.Token
}

func (d *decoder) token() error {
	var err error
	d.Token, err = d.XmlDecoder.Token()
	return err
}

func (d *decoder) skipSpace() error {
	d.Token = nil
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

func (d *decoder) skipElement() error {
	d.Token = nil
	return d.XmlDecoder.Skip()
}

func (d *decoder) decodeGuid() (string, error) {
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

func (d *decoder) decodeTimeSpec() (time.Time, error) {
	var zero time.Time

	outerName := d.Token.(xml.StartElement).Name
	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsTs && start.Name.Local == "date") {
		return zero, fmt.Errorf("decodeTimeSpec: Expected date", err)
	}
	err = d.token()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if !ok {
		return zero, fmt.Errorf("decodeTimeSpec: Expected date string", err)
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", string(cd))
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	err = d.skipElement()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}

	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	if ok && start.Name.Space == nsTs && start.Name.Local == "ns" {
		cd, ok := d.Token.(xml.CharData)
		if !ok {
			return zero, fmt.Errorf("decodeTimeSpec: Expected date string", err)
		}
		ns, err := strconv.Atoi(string(cd))
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
		if !(0 <= ns && ns <= 999999999) {
			return zero, fmt.Errorf("decodeTimeSpec: Invalid nanoseconds value", err)
		}
		t = t.Add(time.Duration(ns))
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name == outerName) {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
	}
	return t, nil
}

func (d *decoder) decodeCount() (string, int, error) {
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

func (d *decoder) decodeAccount() (Account, error) {
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

func (d *decoder) decodeTransaction() (Transaction, error) {
	var zero, trn Transaction

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTransaction: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsTrn && start.Name.Local == "id") {
		return zero, fmt.Errorf("decodeTransaction: Expected id")
	}
	trn.Id, err = d.decodeGuid()
	if err != nil {
		return zero, fmt.Errorf("decodeTransaction: %w", err)
	}

	for {
		err := d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
		start, ok := d.Token.(xml.StartElement)
		if !(ok && start.Name.Space == nsTrn && (start.Name.Local == "currency" || start.Name.Local == "num")) {
			break
		}
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
	}

	start, ok = d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsTrn && start.Name.Local == "date-posted") {
		return zero, fmt.Errorf("decodeTransaction: Expected date-posted")
	}
	trn.DatePosted, err = d.decodeTimeSpec()
	if err != nil {
		return zero, fmt.Errorf("decodeTransaction: %w", err)
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsGnc && end.Name.Local == "transaction") {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
	}
	return trn, nil
}

func (d *decoder) decodeBook() (Book, error) {
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
		case "transaction":
			book.Transactions = make([]Transaction, count)
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

	for i := range len(book.Transactions) {
		start, ok := d.Token.(xml.StartElement)
		if !(ok && start.Name.Space == nsGnc && start.Name.Local == "transaction") {
			return zero, fmt.Errorf("decodeBook: Expected transaction")
		}
		book.Transactions[i], err = d.decodeTransaction()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	}
	if start, ok := d.Token.(xml.StartElement); ok && start.Name.Space == nsGnc && start.Name.Local == "transaction" {
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

func (d *decoder) decode() (Book, error) {
	var zero Book

	err := d.skipSpace()
	if err != nil {
		return zero, err
	}
	if _, ok := d.Token.(xml.EndElement); ok {
		return zero, fmt.Errorf("Expected count-data")
	}
	countType, count, err := d.decodeCount()
	if err != nil {
		return zero, err
	}
	if countType != "book" {
		return zero, fmt.Errorf("Expected book count")
	}
	if count != 1 {
		return zero, fmt.Errorf("Invalid book count")
	}

	err = d.skipSpace()
	if err != nil {
		return zero, err
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsGnc && start.Name.Local == "book") {
		return zero, fmt.Errorf("Expected book element")
	}
	book, err := d.decodeBook()
	if err != nil {
		return zero, err
	}

	err = d.skipElement()
	if err != nil {
		return zero, err
	}
	return book, nil
}

func Decode(d *xml.Decoder) (Book, error) {
	book, err := (&decoder{XmlDecoder: d}).decode()
	if err != nil {
		return Book{}, fmt.Errorf("Decode: %w", err)
	}
	return book, nil
}
