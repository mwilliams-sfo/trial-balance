package v2

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"time"
	xmlutil "trial_balance/gnucash/internal/file/xml/util"
)

var guidRe = regexp.MustCompile("\\A[0-9a-f]{32}\\z")
var numericRe = regexp.MustCompile("\\A-?[0-9]+/-?[0-9]+\\z")

func (d *decoder) decodeGuid() (string, error) {
	start := d.Token.(xml.StartElement)
	var idType string
	idType, ok := xmlutil.GetAttr(start, xml.Name{Local: "type"})
	if !ok {
		return "", fmt.Errorf("decodeGuid: Expected type attribute")
	}
	if idType != "guid" {
		return "", fmt.Errorf("decodeGuid: Expected GUID")
	}
	err := d.token()
	if err != nil {
		return "", fmt.Errorf("decodeGuid: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if !ok {
		return "", fmt.Errorf("decodeGuid: Expected GUID text")
	}
	guid := string(cd)
	if !guidRe.MatchString(guid) {
		return "", fmt.Errorf("decodeGuid: Invalid GUID")
	}
	err = d.skipElement()
	if err != nil {
		return "", fmt.Errorf("decodeGuid: %w", err)
	}
	return guid, nil
}

func (d *decoder) decodeTimeStamp() (time.Time, error) {
	var zero, t time.Time
	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	if start, ok := d.Token.(xml.StartElement); ok && start.Name.Space == nsTs && start.Name.Local == "date" {
		err = d.token()
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
		if cd, ok := d.Token.(xml.CharData); ok {
			t, err = time.Parse("2006-01-02 15:04:05 -0700", string(cd))
		} else {
			return zero, fmt.Errorf("decodeTimeSpec: Expected date string", err)
		}
	} else {
		return zero, fmt.Errorf("decodeTimeSpec: Expected date", err)
	}
	err = d.skipElement()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	return t, nil
}

func (d *decoder) decodeTimeSpec() (time.Time, error) {
	var zero time.Time
	outerName := d.Token.(xml.StartElement).Name

	t, err := d.decodeTimeStamp()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	err = d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTimeSpec: %w", err)
	}
	if start, ok := d.Token.(xml.StartElement); ok && start.Name.Space == nsTs && start.Name.Local == "ns" {
		cd, ok := d.Token.(xml.CharData)
		if !ok {
			return zero, fmt.Errorf("decodeTimeSpec: Expected date string", err)
		}
		ns, err := strconv.Atoi(string(cd))
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
		if !(0 <= ns && ns <= 999999999) {
			return zero, fmt.Errorf("decodeTimeSpec: Invalid nanoseconds value", err)
		}
		t = t.Add(time.Duration(ns))
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name == outerName) {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTimeSpec: %w", err)
		}
	}
	return t, nil
}

func (d *decoder) decodeCount() (string, int, error) {
	start := d.Token.(xml.StartElement)
	countType, ok := xmlutil.GetAttr(start, xml.Name{Space: nsCd, Local: "type"})
	if !ok {
		return "", 0, fmt.Errorf("decodeCount: Expected type attribute")
	}
	err := d.token()
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: %w", err)
	}
	cd, ok := d.Token.(xml.CharData)
	if !ok {
		return "", 0, fmt.Errorf("decodeCount: Expected count value")
	}
	count, err := strconv.Atoi(string(cd))
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: %w", err)
	}
	err = d.skipElement()
	if err != nil {
		return "", 0, fmt.Errorf("decodeCount: Invalid count value")
	}
	return countType, count, nil
}
