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
	for i := 0; ; i++ {
		if i < len(start.Attr) {
			attr := start.Attr[i]
			if attr.Name.Space == "" && attr.Name.Local == "type" {
				idType = attr.Value
				break
			}
		} else {
			return "", fmt.Errorf("decodeGuid: Expected type attribute")
		}
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
	var countType string
	for i := 0; ; i++ {
		if i < len(start.Attr) {
			attr := start.Attr[i]
			if attr.Name.Space == nsCd && attr.Name.Local == "type" {
				countType = attr.Value
				break
			}
		} else {
			return "", 0, fmt.Errorf("decodeCount: Expected type attribute")
		}
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

func (d *Decoder) decodeBook() (*Book, error) {
	book := &Book{}

	err := d.skipSpace()
	if err != nil {
		return nil, fmt.Errorf("decodeBook: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsBook && start.Name.Local == "id") {
		return nil, fmt.Errorf("decodeBook: Expected id")
	}
	book.Id, err = d.decodeGuid()
	if err != nil {
		return nil, fmt.Errorf("decodeBook: %w", err)
	}
	err = d.skipSpace()
	if err != nil {
		return nil, fmt.Errorf("decodeBook: %w", err)
	}
	return book, nil
}

func (d *Decoder) DecodeFile() (*Book, error) {
	err := d.skipSpace()
	if err != nil {
		return nil, fmt.Errorf("DecodeFile: %w", err)
	}
	if _, ok := d.Token.(xml.EndElement); ok {
		return nil, fmt.Errorf("DecodeFile: Expected count-data")
	}
	countType, count, err := d.decodeCount()
	if err != nil {
		return nil, fmt.Errorf("DecodeFile: %w", err)
	}
	if countType != "book" {
		return nil, fmt.Errorf("DecodeFile: Expected book count")
	}
	if count != 1 {
		return nil, fmt.Errorf("DecodeFile: Invalid book count")
	}

	err = d.skipSpace()
	if err != nil {
		return nil, fmt.Errorf("DecodeFile: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if !(ok && start.Name.Space == nsGnc && start.Name.Local == "book") {
		return nil, fmt.Errorf("DecodeFile: Expected book element")
	}
	return d.decodeBook()
}
