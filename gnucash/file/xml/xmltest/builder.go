package xmltest

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
	"trial_balance/gnucash/file/xml/model"
)

type Builder interface {
	Token(tok xml.Token) Builder
	Error(err error) Builder
	StartElement(space, name string, attr ...xml.Attr) Builder
	CharData(s string) Builder
	EndElement(space, name string) Builder
	TextElement(space, name string, attrs []xml.Attr, text string) Builder
	Guid(space, name string, guid model.Guid) Builder
	GDate(d model.GDate) Builder
	TimeStamp(t time.Time) Builder
	TimeSpec(t time.Time) Builder
	Account(a model.Account) Builder
	BuildDecoder() *xml.Decoder
}

type builder []any

func NewBuilder() Builder {
	return builder(nil)
}

func (b builder) Token(tok xml.Token) Builder {
	return append(b, tok)
}

func (b builder) Error(err error) Builder {
	return append(b, err)
}

func (b builder) StartElement(space, name string, attrs ...xml.Attr) Builder {
	return b.Token(xml.StartElement{xml.Name{space, name}, attrs})
}

func (b builder) CharData(s string) Builder {
	return b.Token(xml.CharData(s))
}

func (b builder) EndElement(space, name string) Builder {
	return b.Token(xml.EndElement{xml.Name{space, name}})
}

func (b builder) TextElement(space, name string, attrs []xml.Attr, text string) Builder {
	return b.
		StartElement(space, name, attrs...).
		CharData(text).
		EndElement(space, name)
}

func (b builder) Guid(space, name string, guid model.Guid) Builder {
	return b.
		StartElement(space, name,
			xml.Attr{
				Name:  xml.Name{Local: "type"},
				Value: "guid",
			}).
		CharData(string(guid)).
		EndElement(space, name)
}

func (b builder) GDate(d model.GDate) Builder {
	return b.TextElement("", "gdate", nil,
		fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day))
}

func (b builder) TimeStamp(t time.Time) Builder {
	return b.TextElement("ts", "date", nil,
		t.Format("2006-01-02 15:04:05 -0700"))
}

func (b builder) TimeSpec(t time.Time) Builder {
	return b.
		TimeStamp(t).
		TextElement("ts", "ns", nil,
			string(strconv.AppendInt(nil, int64(t.Nanosecond()), 10)))
}

func (b builder) Account(a model.Account) Builder {
	b2 := b.
		StartElement("gnc", "account").
		TextElement("act", "name", nil, a.Name).
		Guid("act", "id", a.Id)
	if a.HasDescription {
		b2 = b2.TextElement("act", "description", nil, a.Description)
	}
	if a.HasParent {
		b2 = b2.Guid("act", "parent", a.Parent)
	}
	return b2.EndElement("gnc", "account")
}

func (b builder) BuildDecoder() *xml.Decoder {
	return xml.NewTokenDecoder(&reader{b, 0})
}
