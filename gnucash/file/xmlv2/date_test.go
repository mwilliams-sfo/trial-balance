package xmlv2

import (
	"encoding/xml"
	"testing"
	"time"
	"trial_balance/gnucash/file/xmlv2/model"
	"trial_balance/gnucash/file/xmlv2/xmltest"
)

func TestParseGDate(t *testing.T) {
	decoder := xmltest.NewBuilder().
		StartElement("", "gdate").
		CharData("2026-04-25").
		EndElement("", "gdate").
		BuildDecoder()

	tok, _ := decoder.RawToken()
	date, err := parseGDate(decoder, tok.(xml.StartElement))
	if err != nil {

		t.Fatalf("Unexpected error: %v", err)
	}
	expected := model.GDate{2026, 4, 25}
	if date != expected {
		t.Errorf("parseGDate() = %v; want %v", date, expected)
	}
}

func TestParseTimeSpec(t *testing.T) {
	t.Run("without nanoseconds", func(t *testing.T) {
		expected := time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*3600))
		decoder := xmltest.NewBuilder().
			StartElement("act", "date-posted").
			TextElement("ts", "date", nil, "2006-01-02 15:04:05 -0700").
			EndElement("act", "date-posted").
			BuildDecoder()

		tok, _ := decoder.RawToken()
		ts, err := parseTimeSpec(decoder, tok.(xml.StartElement))

		if err != nil {
			t.Fatal(err)
		}
		if ts.Location() != expected.Location() {
			t.Fatalf("Location = %v, want %v", ts.Location(), expected.Location())
		}
		if !ts.Equal(expected) {
			t.Fatalf("TimeStamp = %v, want %v", ts, expected)
		}
	})
	t.Run("with nanoseconds", func(t *testing.T) {
		expected := time.Date(2006, 1, 2, 15, 4, 5, 123, time.FixedZone("", -7*3600))
		decoder := xmltest.NewBuilder().
			StartElement("act", "date-posted").
			TextElement("ts", "date", nil, "2006-01-02 15:04:05 -0700").
			TextElement("ts", "ns", nil, "123").
			EndElement("act", "date-posted").
			BuildDecoder()

		tok, _ := decoder.RawToken()
		ts, err := parseTimeSpec(decoder, tok.(xml.StartElement))

		if err != nil {
			t.Fatal(err)
		}
		if ts.Location() != expected.Location() {
			t.Fatalf("Location = %v, want %v", ts.Location(), expected.Location())
		}
		if !ts.Equal(expected) {
			t.Fatalf("TimeStamp = %v, want %v", ts, expected)
		}
	})
}
