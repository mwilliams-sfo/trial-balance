package xmlv2

import (
	"encoding/xml"
	"testing"
	"trial_balance/gnucash/file/internal/xmlstub"
	"trial_balance/gnucash/file/xmlv2/model"
)

func TestParseAccount(t *testing.T) {
	tokens := xmlstub.NewXml2Builder().
		StartElement("gnc", "account").
		TextElement("act", "name", nil, "Assets").
		Id("act", "id", model.Guid("01234567890123456789012345678901")).
		Id("act", "parent", model.Guid("98765432109876543210987654321098")).
		EndElement("gnc", "account").
		Build()
	decoder := xml.NewTokenDecoder(xmlstub.NewReader(tokens[1:]))

	account, err := parseAccount(decoder, tokens[0].(xml.StartElement))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := model.Account{
		Name:   "Assets",
		Id:     "01234567890123456789012345678901",
		Parent: "98765432109876543210987654321098",
	}
	if account.Name != expected.Name {
		t.Fatalf("Name = %v, want %v", account.Name, expected.Name)
	}
	if account.Id != expected.Id {
		t.Fatalf("Id = %v, want %v", account.Id, expected.Id)
	}
	if account.Parent != expected.Parent {
		t.Fatalf("Parent = %v, want %v", account.Parent, expected.Parent)
	}
}
