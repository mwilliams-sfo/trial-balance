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

func parseTimeSpec(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	var zero time.Time

	tok, err := skipSpace(d)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}
	childStart, ok := tok.(xml.StartElement)
	if !ok || childStart.Name != names.TimeDate() {
		return zero, fmt.Errorf("parseTimeSpec: invalid token, expected <%s>", names.TimeDate())
	}
	text, err := parseTextElement(d, childStart)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", text)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}

	tok, err = skipSpace(d)
	if err != nil {
		return zero, fmt.Errorf("parseTimeSpec: %w", err)
	}
	childStart, ok = tok.(xml.StartElement)
	if ok {
		if childStart.Name != names.TimeNs() {
			return zero, fmt.Errorf("parseTimeSpec: invalid token, expected <%s>", formattedName(names.TimeNs()))
		}
		text, err = parseTextElement(d, childStart)
		if err != nil {
			return zero, fmt.Errorf("parseTimeSpec: %w", err)
		}
		ns, err := strconv.ParseInt(text, 10, 0)
		if err != nil || ns < 0 || ns > 999_999_999 {
			return zero, fmt.Errorf("parseTimeSpec: invalid nanoseconds")
		}
		t = t.Add(time.Duration(ns))
		tok, err = skipSpace(d)
		if err != nil {
			return zero, fmt.Errorf("parseTimeSpec: %w", err)
		}
	}
	end, ok := tok.(xml.EndElement)
	if !ok || end.Name != start.Name {
		return zero, fmt.Errorf("parseTimeSpec: invalid token, expected </%s>", formattedName(start.Name))
	}
	return t, nil
}
