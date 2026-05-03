package xmlv2

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/names"
)

func parseBook(d *xml.Decoder, start xml.StartElement) (model.Book, error) {
	var zero model.Book
	if start.Name != names.Book() {
		return zero, fmt.Errorf("parseBook: invalid start tag")
	}
	var book model.Book

	tok, err := findChildElement(d, start, names.BookId())
	if err != nil {
		return zero, err
	}
	switch x := tok.(type) {
	case xml.StartElement:
		guid, err := parseGuid(d, x)
		if err != nil {
			return zero, err
		}
		book.Id = guid
	default:
		return zero, fmt.Errorf("parseBook: id not found")
	}

	for {
		tok, err = findChildElement(d, start, names.Account())
		if err != nil {
			return zero, err
		}
		switch x := tok.(type) {
		case xml.StartElement:
			account, err := parseAccount(d, x)
			if err != nil {
				return zero, err
			}
			book.Accounts = append(book.Accounts, account)
		case xml.EndElement:
			return book, nil
		default:
			return zero, fmt.Errorf("parseBook: unexpected error: %w", err)
		}
	}
}
