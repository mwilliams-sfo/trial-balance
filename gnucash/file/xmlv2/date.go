package xmlv2

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
	"trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/names"
)

func parseGDate(d *xml.Decoder, start xml.StartElement) (model.GDate, error) {
	var zero model.GDate
	if start.Name != names.GDate() {
		return zero, fmt.Errorf("parseGDate: invalid start tag")
	}
	text, err := parseTextElement(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseGDate: %w", err)
	}
	date, err := time.Parse(time.DateOnly, text)
	return model.GDate{date.Year(), date.Month(), date.Day()}, nil

}

func parseTimeStamp(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	var zero time.Time
	if start.Name != names.TimeDate() {
		return zero, fmt.Errorf("parseTimeStamp: invalid start tag")
	}
	text, err := parseTextElement(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseTimeStamp: %w", err)
	}
	ts, err := time.Parse("2006-01-02 15:04:05 -0700", text)
	return ts, nil
}

func parseNs(d *xml.Decoder, start xml.StartElement) (int, error) {
	if start.Name != names.TimeNs() {
		return 0, fmt.Errorf("parseNs: invalid start tag")
	}
	text, err := parseTextElement(d, start)
	if err != nil {
		return 0, fmt.Errorf("parseNs: %w", err)
	}
	ns, err := strconv.ParseInt(text, 10, 0)
	if err != nil || ns < 0 || 999_999_999 < ns {
		return 0, fmt.Errorf("parseNs: invalid nanoseconds")
	}
	return int(ns), nil
}

func parseTimeSpec(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	var zero time.Time
	ts, err := parseTimeStamp(d, start)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}
	tok, err := d.RawToken()
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}
	if nsStart, ok := tok.(xml.StartElement); ok && nsStart.Name == names.TimeNs() {
		ns, err := parseNs(d, nsStart)
		if err != nil {
			return zero, fmt.Errorf("parseTimeStamp: %w", err)
		}
		ts = ts.Add(time.Duration(ns))
	} else {
		return zero, fmt.Errorf("parseTimeStamp: invalid token, expected <ts:ns>")
	}
	return ts, nil
}
