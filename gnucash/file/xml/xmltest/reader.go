package xmltest

import (
	"encoding/xml"
	"io"
)

type reader struct {
	tokens []any
	index  int
}

func (r *reader) Token() (xml.Token, error) {
	if r.index < len(r.tokens) {
		token := r.tokens[r.index]
		r.index++

		switch x := token.(type) {
		case xml.Token:
			return x, nil
		case error:
			return nil, x
		}
		panic("Invalid token type")
	}
	return nil, io.EOF
}
