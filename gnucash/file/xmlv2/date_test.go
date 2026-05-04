package xmlv2

import (
	"encoding/xml"
	"testing"
	"time"
	"trial_balance/gnucash/file/internal/xmlstub"
	"trial_balance/gnucash/file/xmlv2/model"
)

func TestParseGDate(t *testing.T) {
	tokens := xmlstub.NewXml2Builder().
		StartElement("", "gdate").
		CharData("2026-04-25").
		EndElement("", "gdate").
		Build()
	decoder := xml.NewTokenDecoder(xmlstub.NewReader(tokens[1:]))

	date, err := parseGDate(decoder, tokens[0].(xml.StartElement))
	if err != nil {

		t.Fatalf("Unexpected error: %v", err)
	}
	expected := model.GDate{2026, 4, 25}
	if date != expected {
		t.Errorf("parseGDate() = %v; want %v", date, expected)
	}
}

func TestParseTimeStamp(t *testing.T) {
	tokens := xmlstub.NewXml2Builder().
		StartElement("ts", "date").
		CharData("2006-01-02 15:04:05 -0700").
		EndElement("ts", "date").
		Build()
	decoder := xml.NewTokenDecoder(xmlstub.NewReader(tokens[1:]))

	ts, err := parseTimeStamp(decoder, tokens[0].(xml.StartElement))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*3600))
	if ts.Location() != expected.Location() {
		t.Fatalf("Location = %v, want %v", ts.Location(), expected.Location())
	}
	if !ts.Equal(expected) {
		t.Fatalf("TimeStamp = %v, want %v", ts, expected)
	}
}

func TestParseTimeSpec(t *testing.T) {
	expected := time.Date(2006, 1, 2, 15, 4, 5, 123, time.FixedZone("", -7*3600))
	tokens := xmlstub.NewXml2Builder().
		TimeSpec(expected).
		StartElement("ts", "ns").
		CharData("123").
		EndElement("ts", "ns").
		Build()
	decoder := xml.NewTokenDecoder(xmlstub.NewReader(tokens[1:]))

	ts, err := parseTimeSpec(decoder, tokens[0].(xml.StartElement))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ts.Location() != expected.Location() {
		t.Fatalf("Location = %v, want %v", ts.Location(), expected.Location())
	}
	if !ts.Equal(expected) {
		t.Fatalf("TimeStamp = %v, want %v", ts, expected)
	}
}
