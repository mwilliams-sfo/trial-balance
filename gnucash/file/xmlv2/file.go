package xmlv2

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	fileModel "trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/names"
	"trial_balance/gnucash/model"
)

func xmlDecoder(r io.Reader) (*xml.Decoder, error) {
	buf := make([]byte, 5)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	compressed := bytes.Equal(buf[:2], []byte{0x1f, 0x8b})
	r = io.MultiReader(bytes.NewBuffer(buf), r)
	buf = nil
	if compressed {
		r, err = gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
	}
	return xml.NewDecoder(r), nil
}

func buildBook(fileBook fileModel.Book) (model.Book, error) {
	var zero, book model.Book
	book.Accounts = make([]model.Account, len(fileBook.Accounts))
	accountMap := make(map[string]*model.Account, len(fileBook.Accounts))
	for i, act := range fileBook.Accounts {
		accountMap[string(act.Id)] = &book.Accounts[i]
	}
	for i, act := range fileBook.Accounts {
		parent, ok := accountMap[string(act.Parent)]
		if !ok {
			return zero, fmt.Errorf("Unknown parent account id")
		}
		book.Accounts[i] = model.Account{
			Book:   &book,
			Name:   act.Name,
			Parent: parent,
		}
	}
	return book, nil
}

func Read(r io.Reader) (model.Book, error) {
	var zero model.Book
	d, err := xmlDecoder(r)
	if err != nil {
		return zero, err
	}
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	if pi, ok := tok.(xml.ProcInst); !ok || pi.Target != "xml" {
		return zero, fmt.Errorf("Invalid token, expected XML declaration")
	}
	tok, err = d.RawToken()
	if err != nil {
		return zero, err
	}
	start, ok := tok.(xml.StartElement)
	if !ok || start.Name != names.GncV2() {
		return zero, fmt.Errorf("Invalid token, expected <gnc-v2>")
	}
	gc, err := parseGnuCash(d, start)
	if err != nil {
		return zero, err
	}
	books := []fileModel.Book(gc)
	return buildBook(books[0])
}
