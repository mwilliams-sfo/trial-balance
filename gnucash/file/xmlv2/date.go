package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
	"trial_balance/gnucash/model"
)

func requireStartElement(d *xml.Decoder) (xml.StartElement, error) {
	zero := xml.StartElement{}
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	start, ok := tok.(xml.StartElement)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected start tag")
	}
	return start, nil
}

func requireCharData(d *xml.Decoder) (xml.CharData, error) {
	zero := xml.CharData("")
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	cd, ok := tok.(xml.CharData)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected character data")
	}
	return cd, nil
}

func requireEndElement(d *xml.Decoder) (xml.EndElement, error) {
	zero := xml.EndElement{}
	tok, err := d.RawToken()
	if err != nil {
		return zero, err
	}
	end, ok := tok.(xml.EndElement)
	if !ok {
		return zero, fmt.Errorf("Invalid token, expected end tag")
	}
	return end, nil
}

var gDateName = xml.Name{"", "gdate"}

func parseGDate(d *xml.Decoder, start xml.StartElement) (model.GDate, error) {
	zero := model.GDate{}
	if start.Name != gDateName {
		return zero, fmt.Errorf("parseGDate: invalid start tag, expected <gdate>")
	}

	text, err := requireCharData(d)
	if err != nil {
		return zero, err
	}
	date, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return zero, fmt.Errorf("parseGDate: invalid text: %w", err)
	}

	end, err := requireEndElement(d)
	if err != nil {
		return zero, err
	}
	if end.Name != gDateName {
		return zero, fmt.Errorf("parseGDate: invalid token, expected </gdate>")
	}
	return model.GDate{date.Year(), date.Month(), date.Day()}, nil

}

var timeStampName = xml.Name{"ts", "date"}

func parseTimeStamp(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	zero := time.Time{}
	if start.Name != timeStampName {
		return zero, fmt.Errorf("parseTimeStamp: invalid start tag")
	}

	cd, err := requireCharData(d)
	if err != nil {
		return zero, err
	}
	ts, err := time.Parse("2006-01-02 15:04:05 -0700", string(cd))
	if err != nil {
		return zero, fmt.Errorf("parseTimeStamp: invalid text: %w", err)
	}

	end, err := requireEndElement(d)
	if err != nil {
		return zero, err
	}
	if end.Name != timeStampName {
		return zero, fmt.Errorf("parseTimeStamp: invalid token, expected </ts:date>")
	}
	return ts, nil
}

var nsName = xml.Name{"ts", "ns"}

func parseNs(d *xml.Decoder, start xml.StartElement) (int, error) {
	if start.Name != nsName {
		return 0, fmt.Errorf("parseNs: invalid start tag")
	}
	cd, err := requireCharData(d)
	if err != nil {
		return 0, err
	}
	ns, err := strconv.ParseInt(string(cd), 10, 0)
	if err != nil || ns < 0 || 999_999_999 < ns {
		return 0, fmt.Errorf("parseNs: invalid nanoseconds")
	}
	end, err := requireEndElement(d)
	if err != nil {
		return 0, err
	}
	if end.Name != nsName {
		return 0, fmt.Errorf("parseNs: invalid token, expected </ts:ns>")
	}
	return int(ns), nil
}

func parseTimeSpec(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	zero := time.Time{}
	ts, err := parseTimeStamp(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: invalid time stamp: %w", err)
	}
	nsStart, err := requireStartElement(d)
	if err != nil {
		return zero, err
	}
	if nsStart.Name != nsName {
		return zero, fmt.Errorf("parseTimeSpec: invalid token, expected <ts:ns>")
	}
	ns, err := parseNs(d, nsStart)
	if err != nil {
		return zero, err
	}
	ts = ts.Add(time.Duration(ns))
	return ts, nil
}
