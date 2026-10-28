package v2

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

func (d *decoder) decodeSplit() (Split, error) {
	var zero, split Split

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeSplit: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsSplit && start.Name.Local == "id" {
		split.Id, err = d.decodeGuid()
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeSplit: Expected id")
	}

	for {
		err := d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsSplit &&
			(start.Name.Local == "memo" ||
				start.Name.Local == "action" ||
				start.Name.Local == "reconciled-state" ||
				start.Name.Local == "reconcile-date") {
			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeSplit: %w", err)
			}
		} else {
			break
		}
	}

	start, ok = d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsSplit && start.Name.Local == "value" {
		err := d.token()
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
		cd, ok := d.Token.(xml.CharData)
		if !(ok && numericRe.MatchString(string(cd))) {
			return zero, fmt.Errorf("decodeSplit: Expected numeric string")
		}
		fields := strings.SplitN(string(cd), "/", 2)
		split.Value.Num, err = strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
		split.Value.Den, err = strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeSplit: Expected value")
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsTrn && end.Name.Local == "split") {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeSplit: %w", err)
		}
	}
	return split, nil
}
