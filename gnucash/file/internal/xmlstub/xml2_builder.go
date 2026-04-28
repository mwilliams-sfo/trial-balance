package xmlstub

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
	"trial_balance/gnucash/model"
)

type Xml2Builder interface {
	Tokens(ts ...xml.Token) Xml2Builder
	StartElement(space, name string, attr ...xml.Attr) Xml2Builder
	CharData(s string) Xml2Builder
	EndElement(space, name string) Xml2Builder
	Id(space, name string, guid model.Guid) Xml2Builder
	GDate(d model.GDate) Xml2Builder
	TimeStamp(t time.Time) Xml2Builder
	TimeSpec(t time.Time) Xml2Builder
	Account(a model.Account) Xml2Builder
	Build() []xml.Token
}

type builder []xml.Token

func NewXml2Builder() Xml2Builder {
	return builder(nil)
}

func (b builder) Tokens(ts ...xml.Token) Xml2Builder {
	return append(b, ts...)
}

func (b builder) StartElement(space, name string, attrs ...xml.Attr) Xml2Builder {
	return b.Tokens(xml.StartElement{xml.Name{space, name}, attrs})
}

func (b builder) CharData(s string) Xml2Builder {
	return b.Tokens(xml.CharData(s))
}

func (b builder) EndElement(space, name string) Xml2Builder {
	return b.Tokens(xml.EndElement{xml.Name{space, name}})
}

func (b builder) Id(space, name string, guid model.Guid) Xml2Builder {
	return b.
		StartElement(space, name, xml.Attr{xml.Name{Local: "guid"}, string(guid)}).
		EndElement(space, name)
}

func (b builder) GDate(d model.GDate) Xml2Builder {
	return b.
		StartElement("", "gdate").
		CharData(fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)).
		EndElement("", "gdate")
}

func (b builder) TimeStamp(t time.Time) Xml2Builder {
	return b.
		StartElement("ts", "date").
		CharData(t.Format("2006-01-02 15:04:05 -0700")).
		EndElement("ts", "date")
}

func (b builder) TimeSpec(t time.Time) Xml2Builder {
	return b.
		TimeStamp(t).
		StartElement("ts", "ns").
		CharData(string(strconv.AppendInt(nil, int64(t.Nanosecond()), 10))).
		EndElement("ts", "ns")
}

func (b builder) Account(a model.Account) Xml2Builder {
	return b.
		StartElement("gnc", "account").
		StartElement("act", "name").
		CharData(a.Name).
		EndElement("act", "name").
		Id("act", "id", a.Id).
		Id("act", "parent", a.Parent).
		EndElement("gnc", "account")
}

func (b builder) Build() []xml.Token {
	return b
}
