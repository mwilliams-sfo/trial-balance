package xmlstub

import (
	"encoding/xml"
	"io"
)

type reader struct {
	tokens []xml.Token
	index  int
}

func NewReader(tokens []xml.Token) xml.TokenReader {
	return &reader{tokens, 0}
}
func (r *reader) Token() (xml.Token, error) {
	if r.index < len(r.tokens) {
		token := r.tokens[r.index]
		r.index++
		err, ok := token.(error)
		if ok {
			return nil, err
		}
		return token, nil
		r.index++
	}
	return io.EOF, nil
}
