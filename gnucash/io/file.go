package io

import (
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	gc "trial_balance/gnucash"
)

var (
	xmlNsGnc  = "http://www.gnucash.org/XML/gnc"
	xmlNsBook = "http://www.gnucash.org/XML/book"
	xmlNsTrn  = "http://www.gnucash.org/XML/trn"
	xmlNsCd   = "http://www.gnucash.org/XML/cd"
)

var errUnrecognizedFormat = errors.New("Unrecognized file format")
var errUnexpectedText = errors.New("Unexpected character data")

func errUnexpectedElement(name xml.Name) error {
	return fmt.Errorf("Unexpected element: %s:%s", name.Space, name.Local)
}

func errUnclosedElement(name xml.Name) error {
	return fmt.Errorf("Unclosed element: %s:%s", name.Space, name.Local)
}

func isXmlSpace(charData xml.CharData) bool {
	for _, r := range string(charData) {
		if !(r == '\t' || r == '\n' || r == '\r' || r == ' ') {
			return false
		}
	}
	return true
}

func trimCharData(s string) string {
	return strings.Trim(s, "\t\n\r ")
}

func readXmlV2CountData(decoder *xml.Decoder, start xml.StartElement) (string, int, error) {
	typeIndex := -1
	for i, attr := range start.Attr {
		if attr.Name.Space == xmlNsCd && attr.Name.Local == "type" {
			typeIndex = i
		}
	}
	if typeIndex == -1 {
		return "", 0, fmt.Errorf("Invalid count-data (no type)")
	}
	typeValue := start.Attr[typeIndex].Value

	token, err := decoder.Token()
	if err != nil {
		return "", 0, err
	}
	cd, ok := token.(xml.CharData)
	if !ok {
		return "", 0, fmt.Errorf("Invalid count-data")
	}
	count, err := strconv.ParseInt(trimCharData(string(cd)), 10, 32)
	if err != nil {
		return "", 0, err
	}
	if count < 0 || math.MaxInt < count {
		return "", 0, fmt.Errorf("count-data value out of range")
	}

	token, err = decoder.Token()
	if err != nil {
		return "", 0, err
	}
	if _, ok := token.(xml.EndElement); !ok {
		return "", 0, errUnclosedElement(start.Name)
	}
	return typeValue, int(count), nil
}

func readXmlV2Guid(decoder *xml.Decoder, start xml.StartElement) (gc.Guid, error) {
	noGuid := gc.Guid("")
	token, err := decoder.Token()
	if err != nil {
		return noGuid, err
	}
	cd, ok := token.(xml.CharData)
	if !ok || len(string(cd)) == 0 {
		return noGuid, fmt.Errorf("Invalid GUID")
	}
	guid := gc.Guid(string(cd))
	token, err = decoder.Token()
	if err != nil {
		return noGuid, err
	}
	if _, ok := token.(xml.EndElement); !ok {
		return noGuid, errUnclosedElement(start.Name)
	}
	return guid, nil
}

func readXmlV2Id(decoder *xml.Decoder, start xml.StartElement) (gc.Id, error) {
	typeIndex := -1
	for i, attr := range start.Attr {
		if attr.Name.Local == "type" {
			typeIndex = i
			break
		}
	}
	if typeIndex == -1 {
		return nil, fmt.Errorf("Invalid id (no type)")
	}
	typeValue := start.Attr[typeIndex].Value
	if typeValue == "guid" {
		return readXmlV2Guid(decoder, start.Copy())
	}
	return nil, fmt.Errorf("Invalid id type: %s", typeValue)
}

func readXmlV2Transaction(decoder *xml.Decoder, start xml.StartElement) (*gc.Transaction, error) {
	fmt.Println("Reading a transaction")
	trn := gc.Transaction{}
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Space == xmlNsTrn && se.Name.Local == "id" {
				trn.Id, err = readXmlV2Id(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
			} else if err = decoder.Skip(); err != nil {
				return nil, err
			}
		} else if cd, ok := token.(xml.CharData); ok {
			if !isXmlSpace(cd) {
				return nil, errUnexpectedText
			}
		} else if _, ok = token.(xml.EndElement); ok {
			break
		}
	}
	return &trn, nil
}

func readXmlV2Book(decoder *xml.Decoder, start xml.StartElement) (*gc.Book, error) {
	fmt.Println("Reading a book")
	book := gc.Book{}
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Space == xmlNsBook && se.Name.Local == "id" {
				book.Id, err = readXmlV2Id(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
			} else if se.Name.Space == xmlNsCd && se.Name.Local == "count-data" {
				countType, count, err := readXmlV2CountData(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
				if countType == "transaction" {
					trns := make([]gc.Transaction, count)
					book.Transactions = trns[:0]
				}
			} else if se.Name.Space == xmlNsGnc && se.Name.Local == "transaction" {
				trn, err := readXmlV2Transaction(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
				book.Transactions = append(book.Transactions, *trn)
			} else if err = decoder.Skip(); err != nil {
				return nil, err
			}
		} else if cd, ok := token.(xml.CharData); ok {
			if !isXmlSpace(cd) {
				return nil, errUnexpectedText
			}
		} else if _, ok = token.(xml.EndElement); ok {
			break
		}
	}
	return &book, nil
}

func readXmlV2(decoder *xml.Decoder, start xml.StartElement) (*gc.Book, error) {
	fmt.Println("Reading XML v2")
	var book *gc.Book
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Space == xmlNsGnc && se.Name.Local == "count-data" {
				countType, count, err := readXmlV2CountData(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
				if countType == "book" && count != 1 {
					return nil, fmt.Errorf("Book count is %d, expected 1", count)
				}
			} else if se.Name.Space == xmlNsGnc && se.Name.Local == "book" {
				if book != nil {
					return nil, fmt.Errorf("Too many books")
				}
				book, err = readXmlV2Book(decoder, se.Copy())
				if err != nil {
					return nil, err
				}
			} else if err = decoder.Skip(); err != nil {
				return nil, err

			}
		} else if cd, ok := token.(xml.CharData); ok {
			if !isXmlSpace(cd) {
				return nil, errUnrecognizedFormat
			}
		} else if _, ok = token.(xml.EndElement); ok {
			break
		}
	}
	return book, nil
}

func readXml(r io.Reader) (*gc.Book, error) {
	decoder := xml.NewDecoder(r)
	token, err := decoder.Token()
	if err != nil {
		return nil, errUnrecognizedFormat
	}
	if pi, ok := token.(xml.ProcInst); !ok || pi.Target != "xml" {
		return nil, errUnrecognizedFormat
	}

	token, err = decoder.Token()
	if err != nil {
		return nil, err
	}
	if cd, ok := token.(xml.CharData); ok {
		if !isXmlSpace(cd) {
			return nil, errUnexpectedText
		}
		token, err = decoder.Token()
		if err != nil {
			return nil, err
		}
	}

	if se, ok := token.(xml.StartElement); ok {
		if se.Name.Space == "" && se.Name.Local == "gnc-v2" {
			return readXmlV2(decoder, se.Copy())
		}
		return nil, fmt.Errorf("Unrecognized document element: %s", se.Name.Local)
	}
	return nil, fmt.Errorf("No document element")
}

func Read(name string) (*gc.Book, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	book, err := readXml(r)
	if err == errUnrecognizedFormat {
		_, err = r.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		var gzr io.Reader
		gzr, err = gzip.NewReader(r)
		if err != nil {
			return nil, errUnrecognizedFormat
		}
		book, err = readXml(gzr)
	}
	return book, err
}
