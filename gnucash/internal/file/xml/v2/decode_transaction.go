package v2

import (
	"encoding/xml"
	"fmt"
)

func (d *decoder) decodeTransaction() (Transaction, error) {
	var zero, trn Transaction

	err := d.skipSpace()
	if err != nil {
		return zero, fmt.Errorf("decodeTransaction: %w", err)
	}
	start, ok := d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsTrn && start.Name.Local == "id" {
		trn.Id, err = d.decodeGuid()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeTransaction: Expected id")
	}

	for {
		err := d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsTrn &&
			(start.Name.Local == "currency" ||
				start.Name.Local == "num") {
			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeTransaction: %w", err)
			}
		} else {
			break
		}
	}

	start, ok = d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsTrn && start.Name.Local == "date-posted" {
		trn.DatePosted, err = d.decodeTimeSpec()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
	} else {
		return zero, fmt.Errorf("decodeTransaction: Expected date-posted")
	}

	for {
		err := d.skipSpace()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
		start, ok := d.Token.(xml.StartElement)
		if ok && start.Name.Space == nsTrn &&
			(start.Name.Local == "date-entered" ||
				start.Name.Local == "description" ||
				start.Name.Local == "slots") {

			err = d.skipElement()
			if err != nil {
				return zero, fmt.Errorf("decodeTransaction: %w", err)
			}
		} else {
			break
		}
	}

	start, ok = d.Token.(xml.StartElement)
	if ok && start.Name.Space == nsTrn && start.Name.Local == "splits" {
		for {
			err := d.skipSpace()
			if err != nil {
				return zero, fmt.Errorf("decodeTransaction: %w", err)
			}
			start, ok = d.Token.(xml.StartElement)
			if ok && start.Name.Space == nsTrn && start.Name.Local == "split" {
				split, err := d.decodeSplit()
				if err != nil {
					return zero, fmt.Errorf("decodeTransaction: %w", err)
				}
				trn.Splits = append(trn.Splits, split)
			} else if len(trn.Splits) == 0 {
				return zero, fmt.Errorf("decodeTransaction: Expected split")
			} else {
				break
			}
		}
		if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsTrn && end.Name.Local == "splits") {
			return zero, fmt.Errorf("decodeTransaction: Expected splits")
		}
	}

	if end, ok := d.Token.(xml.EndElement); !(ok && end.Name.Space == nsGnc && end.Name.Local == "transaction") {
		err = d.skipElement()
		if err != nil {
			return zero, fmt.Errorf("decodeTransaction: %w", err)
		}
	}
	return trn, nil
}
