package gnucash

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	xmlutil "trial_balance/gnucash/internal/file/xml/util"
)

func decode(d *xml.Decoder) (Book, error) {
	var zero Book

	tok, err := d.Token()
	if err != nil {
		return zero, err
	}
	if pi, ok := tok.(xml.ProcInst); !ok || pi.Target != "xml" {
		return zero, fmt.Errorf("Expected XML prolog")
	}

	tok, err = d.Token()
	if err != nil {
		return zero, err
	}
	if cd, ok := tok.(xml.CharData); ok {
		if xmlutil.IsXmlSpace(string(cd)) {
			tok, err = d.Token()
			if err != nil {
				return zero, err
			}
		}
	}
	start, ok := tok.(xml.StartElement)
	if !ok {
		return zero, fmt.Errorf("Expected root element")
	}
	switch {
	case start.Name.Space == "" && start.Name.Local == "gnc-v2":
		return decodeXmlV2(d)
	default:
		return zero, fmt.Errorf("Unrecognized file format")
	}
}

func ReadFile(name string) (Book, error) {
	var zero Book

	f, err := os.Open(name)
	if err != nil {
		return zero, fmt.Errorf("ReadFile: %w", err)
	}
	defer f.Close()

	var r io.Reader
	r, err = gzip.NewReader(f)
	switch err {
	case gzip.ErrHeader:
		_, err = f.Seek(0, os.SEEK_SET)
		if err != nil {
			return zero, fmt.Errorf("ReadFile: %w", err)
		}
		r = f
	case nil:
	default:
		return zero, fmt.Errorf("ReadFile: %w", err)
	}

	return decode(xml.NewDecoder(r))
}
