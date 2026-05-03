package xmlv2

import (
	"encoding/xml"
	"testing"
	"trial_balance/gnucash/file/internal/xmlstub"
	"trial_balance/gnucash/file/xmlv2/model"
)

func TestParseBook(t *testing.T) {
	expected := model.Book{
		Id: "01234567890123456789012345678901",
		Accounts: []model.Account{
			{
				Name:   "Assets",
				Id:     "12345678901234567890123456789012",
				Parent: "23456789012345678901234567890123",
			},
		},
	}
	tokens := xmlstub.NewXml2Builder().
		StartElement("gnc", "book").
		Guid("book", "id", expected.Id).
		Account(expected.Accounts[0]).
		EndElement("gnc", "book").
		Build()
	decoder := xml.NewTokenDecoder(xmlstub.NewReader(tokens[1:]))

	book, err := parseBook(decoder, tokens[0].(xml.StartElement))
	if err != nil {
		t.Fatal(err)
	}
	if book.Id != expected.Id {
		t.Fatalf("Id = %v, want %v", book.Id, expected.Id)
	}
	if len(book.Accounts) != 1 || book.Accounts[0] != expected.Accounts[0] {
		t.Fatalf("Accounts = %v, want %v", book.Accounts, expected.Accounts)
	}
}
