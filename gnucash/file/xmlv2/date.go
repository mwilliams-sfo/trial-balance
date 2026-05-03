package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
	"trial_balance/gnucash/file/xmlv2/names"
	"trial_balance/gnucash/model"
)

func parseGDate(d *xml.Decoder, start xml.StartElement) (model.GDate, error) {
	zero := model.GDate{}
	if start.Name != names.GDate() {
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
	if end.Name != names.GDate() {
		return zero, fmt.Errorf("parseGDate: invalid token, expected </gdate>")
	}
	return model.GDate{date.Year(), date.Month(), date.Day()}, nil

}

func parseTimeStamp(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	zero := time.Time{}
	if start.Name != names.TimeDate() {
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
	if end.Name != names.TimeDate() {
		return zero, fmt.Errorf("parseTimeStamp: invalid token, expected </ts:date>")
	}
	return ts, nil
}

func parseNs(d *xml.Decoder, start xml.StartElement) (int, error) {
	if start.Name != names.TimeNs() {
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
	if end.Name != names.TimeNs() {
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
	if nsStart.Name != names.TimeNs() {
		return zero, fmt.Errorf("parseTimeSpec: invalid token, expected <ts:ns>")
	}
	ns, err := parseNs(d, nsStart)
	if err != nil {
		return zero, err
	}
	ts = ts.Add(time.Duration(ns))
	return ts, nil
}
