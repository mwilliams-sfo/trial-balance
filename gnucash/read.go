package gnucash

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func isXmlSpace(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("\t\n\r ", r) {
			return false
		}
	}
	return true
}

func decode(d *xml.Decoder) (*Book, error) {
	tok, err := d.Token()
	if err != nil {
		return nil, err
	}
	if pi, ok := tok.(xml.ProcInst); !ok || pi.Target != "xml" {
		return nil, fmt.Errorf("Expected XML prolog")
	}

	tok, err = d.Token()
	if err != nil {
		return nil, err
	}
	if cd, ok := tok.(xml.CharData); ok {
		if isXmlSpace(string(cd)) {
			tok, err = d.Token()
			if err != nil {
				return nil, err
			}
		}
	}
	start, ok := tok.(xml.StartElement)
	if !ok {
		return nil, fmt.Errorf("Expected root element")
	}
	switch {
	case start.Name.Space == "" && start.Name.Local == "gnc-v2":
		return &Book{}, nil
	default:
		return nil, fmt.Errorf("Unrecognized file format")
	}
}

func ReadFile(name string) (*Book, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}
	defer f.Close()

	var r io.Reader
	r, err = gzip.NewReader(f)
	switch err {
	case gzip.ErrHeader:
		_, err = f.Seek(0, os.SEEK_SET)
		if err != nil {
			return nil, fmt.Errorf("ReadFile: %w", err)
		}
		r = f
	case nil:
	default:
		return nil, fmt.Errorf("ReadFile: %w", err)
	}

	return decode(xml.NewDecoder(r))
}
