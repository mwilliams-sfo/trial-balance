package xml

import (
	"encoding/xml"
	"fmt"
	"trial_balance/gnucash/file/xml/model"
	"trial_balance/gnucash/file/xml/names"
)

func parseBook(d *xml.Decoder, start xml.StartElement) (model.Book, error) {
	var zero, book model.Book
	tok, err := skipSpace(d)
	if err != nil {
		return zero, fmt.Errorf("parseBook: %w", err)
	}
	childStart, ok := tok.(xml.StartElement)
	if !ok || childStart.Name != names.BookId() {
		return zero, fmt.Errorf("parseBook: invalid token, expected <%s>", formattedName(names.BookId()))
	}
	book.Id, err = parseGuid(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseBook: %w", err)
	}

	for {
		tok, err := skipSpace(d)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
		switch x := tok.(type) {
		case xml.StartElement:
			if x.Name == names.CountData() || x.Name == (xml.Name{"gnc", "pricedb"}) {
				goto countLoop
			}
			err = skipContent(d, x)
			if err != nil {
				return zero, fmt.Errorf("parseBook: %w", err)
			}
		case xml.EndElement:
			goto countLoop
		}
	}
countLoop:
	for {
		childStart, ok = tok.(xml.StartElement)
		if !ok || childStart.Name != names.CountData() {
			break
		}
		countType, count, err := parseCountData(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
		switch countType {
		case "account":
			book.Accounts = make([]model.Account, count)
		}
		tok, err = skipSpace(d)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
	}

	for {
		childStart, ok = tok.(xml.StartElement)
		if !ok || childStart.Name == names.Account() {
			break
		}
		err = skipContent(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
	}
	for i := 0; ; i++ {
		childStart, ok = tok.(xml.StartElement)
		if !ok || childStart.Name != names.Account() {
			if i < len(book.Accounts) {
				return zero, fmt.Errorf("parseBook: too few accounts")
			}
			break
		}
		if i == len(book.Accounts) {
			return zero, fmt.Errorf("parseBook: too many accounts")
		}
		var err error
		book.Accounts[i], err = parseAccount(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
	}

	if end, ok := tok.(xml.EndElement); ok {
		if end.Name != start.Name {
			return zero, fmt.Errorf("parseBook: invalid token, expected </%s>", formattedName(start.Name))
		}
	} else {
		err = skipContent(d, start)
		if err != nil {
			return zero, fmt.Errorf("parseBook: %w", err)
		}
	}
	return book, nil
}
