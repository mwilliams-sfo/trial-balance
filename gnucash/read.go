package gnucash

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func decode(d *xml.Decoder) (*Book, error) {
	return nil, fmt.Errorf("Not implemented")
}

func ReadFile(name string) (*Book, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}
	defer f.Close()

	var r io.Reader
	r, err = gzip.NewReader(f)
	if err == gzip.ErrHeader {
		_, err = f.Seek(0, os.SEEK_SET)
		if err != nil {
			return nil, fmt.Errorf("ReadFile: %w", err)
		}
		r = f
	} else if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}

	return decode(xml.NewDecoder(r))
}
