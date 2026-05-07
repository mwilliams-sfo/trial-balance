package xmlv2

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/names"
)

func parseGnuCash(d *xml.Decoder, start xml.StartElement) (model.GnuCash, error) {
	var err error
	tok, err := d.RawToken()
	if err != nil {
		return nil, fmt.Errorf("parseGnuCash: %w", err)
	}

	var bookCount int
	if childStart, ok := tok.(xml.StartElement); ok && childStart.Name == names.CountData() {
		countType, count, err := parseCountData(d, childStart)
		if err != nil {
			return nil, fmt.Errorf("parseGnuCash: %w", err)
		}
		if countType != "book" {
			return nil, fmt.Errorf("parseGnuCash: wrong count type, expected book")
		}
		bookCount = count
	} else {
		return nil, fmt.Errorf("parseGnuCash: invalid token, expected <gnc:count-data>")
	}

	books := make([]model.Book, bookCount)
	for i := range bookCount {
		tok, err := d.RawToken()
		if err != nil {
			return nil, fmt.Errorf("parseGnuCash: %w", err)
		}
		if childStart, ok := tok.(xml.StartElement); ok && childStart.Name == names.Book() {
			book, err := parseBook(d, start)
			if err != nil {
				return nil, fmt.Errorf("parseGnuCash: %w", err)
			}
			books[i] = book
		} else {
			return nil, fmt.Errorf("parseGnuCash: invalid token, expected <gnc:book>")
		}
	}

	err = skipContent(d, start)
	if err != nil {
		return nil, fmt.Errorf("parseGnuCash: %w", err)
	}
	return model.GnuCash(books), nil
}
