package xmlv2

import (
	"encoding/xml"
	"fmt"
	"time"
	"trial_balance/gnucash/model"
)

func requireCharData(d *xml.Decoder) (string, error) {
	tok, err := d.RawToken()
	if err != nil {
		return "", err
	}
	cd, ok := tok.(xml.CharData)
	if !ok {
		return "", fmt.Errorf("Character data expected")
	}
	return string(cd), nil
}

func requireEndElement(d *xml.Decoder, name xml.Name) error {
	tok, err := d.RawToken()
	if err != nil {
		return err
	}
	end, ok := tok.(xml.EndElement)
	if !ok || end.Name != name {
		return fmt.Errorf("</%s:%s> expected", name.Space, name.Local)
	}
	return nil
}

var gDateName = xml.Name{"", "gdate"}

func parseGDate(d *xml.Decoder, start xml.StartElement) (model.GDate, error) {
	if start.Name != gDateName {
		return model.GDate{}, fmt.Errorf("parseGDate: invalid start tag: %s", start.Name.Local)
	}

	text, err := requireCharData(d)
	if err != nil {
		return model.GDate{}, err
	}
	date, err := time.Parse(time.DateOnly, text)
	if err != nil {
		return model.GDate{}, fmt.Errorf("parseGDate: invalid text: %w", err)
	}

	err = requireEndElement(d, gDateName)
	if err != nil {
		return model.GDate{}, err
	}
	return model.GDate{date.Year(), date.Month(), date.Day()}, nil

}

var timeStampName = xml.Name{"ts", "date"}

func parseTimeStamp(d *xml.Decoder, start xml.StartElement) (time.Time, error) {
	zero:=time.Time{}
	if start.Name != timeStampName {
		return zero, fmt.Errorf("Parsemodel.TimeStamp: invalid start tag")
	}

	text, err := requireCharData(d)
	if err != nil {
		return zero, err
	}
	if err != nil {
		return zero, err
	}
	ts, err := time.Parse("2006-01-02 15:04:05 -0700", text)
	if err != nil {
		return zero, fmt.Errorf("Parsemodel.TimeStamp: invalid text: %w", err)
	}

	err = requireEndElement(d, timeStampName)
	if err != nil {
		return zero, err
	}
	return ts, nil
}
