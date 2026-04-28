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
	Id(name xml.Name, guid model.Guid) Xml2Builder
	GDate(d model.GDate) Xml2Builder
	TimeStamp(t time.Time) Xml2Builder
	TimeSpec(t time.Time) Xml2Builder
	Build() []xml.Token
}

type builder []xml.Token

func NewXml2Builder() Xml2Builder {
	return builder(nil)
}

func (b builder) Tokens(ts ...xml.Token) Xml2Builder {
	return append(b, ts...)
}

func (b builder) Id(name xml.Name, guid model.Guid) Xml2Builder {
	start := xml.StartElement{
		name,
		[]xml.Attr{
			{xml.Name{"", "guid"}, string(guid)},
		},
	}
	return b.Tokens(start, start.End())
}

func (b builder) GDate(d model.GDate) Xml2Builder {
	return b.Tokens(
		xml.StartElement{xml.Name{"", "gdate"}, nil},
		xml.CharData(fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)),
		xml.EndElement{xml.Name{"", "gdate"}},
	)
}

func (b builder) TimeStamp(t time.Time) Xml2Builder {
	return b.Tokens(
		xml.StartElement{xml.Name{"ts", "date"}, nil},
		xml.CharData(t.Format("2006-01-02 15:04:05 -0700")),
		xml.EndElement{xml.Name{"ts", "date"}},
	)
}

func (b builder) TimeSpec(t time.Time) Xml2Builder {
	return b.TimeStamp(t).Tokens(
		xml.StartElement{xml.Name{"ts", "ns"}, nil},
		xml.CharData(strconv.AppendInt(nil, int64(t.Nanosecond()), 10)),
		xml.EndElement{xml.Name{"ts", "ns"}},
	)
}

func (b builder) Build() []xml.Token {
	return b
}
