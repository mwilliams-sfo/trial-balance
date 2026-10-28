package v2

import (
	"encoding/xml"
	"fmt"
	xmlutil "trial_balance/gnucash/internal/file/xml/util"
)

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

func (d *decoder) decode() (Book, error) {
	var zero, book Book

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
	if ok && start.Name.Space == nsGnc && start.Name.Local == "book" {
		book, err = d.decodeBook()
		if err != nil {
			return zero, err
		}
	} else {
		return zero, fmt.Errorf("Expected book element")
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == "" && end.Name.Local == "gnc-v2") {
		err = d.skipElement()
		if err != nil {
			return zero, err
		}
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
