package v2

import (
	"encoding/xml"
	"fmt"
)

func (d *decoder) decodeBook() (Book, error) {
	var zero, book Book

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeBook: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsBook && start.Name.Local == "id" {
		book.Id, err = d.decodeGuid()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeBook: Expected id")
	}

	for {
		err = d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeBook: %w", err)
		}
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsBook && start.Name.Local == "slots" {
			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
		} else {
			break
		}
	}

	for {
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsGnc && start.Name.Local == "count-data" {
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
		} else {
			break
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
		if ok && start.Name.Space == nsGnc && start.Name.Local == "account" {
			book.Accounts[i], err = d.decodeAccount()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
			err = d.skipSpace()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
		} else {
			return zero, fmt.Errorf("decodeBook: Expected account")
		}
	}
	if start, ok := d.Token.(xml.StartElement); ok && start.Name.Space == nsGnc && start.Name.Local == "account" {
		return zero, fmt.Errorf("decodeBook: Unexpected account")
	}

	for i := range len(book.Transactions) {
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsGnc && start.Name.Local == "transaction" {
			book.Transactions[i], err = d.decodeTransaction()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
			err = d.skipSpace()
			if err != nil {
				return zero, fmt.Errorf("decodeBook: %w", err)
			}
		} else {
			return zero, fmt.Errorf("decodeBook: Expected transaction")
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
