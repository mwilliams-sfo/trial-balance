package xml

import (
	"encoding/xml"
	"testing"
	"trial_balance/gnucash/file/xml/model"
	"trial_balance/gnucash/file/xml/xmltest"
)

func TestParseAccount(t *testing.T) {
	decoder := xmltest.NewBuilder().
		StartElement("gnc", "account").
		TextElement("act", "name", nil, "Assets").
		Guid("act", "id", model.Guid("01234567890123456789012345678901")).
		Guid("act", "parent", model.Guid("98765432109876543210987654321098")).
		EndElement("gnc", "account").
		BuildDecoder()

	tok, _ := decoder.RawToken()
	account, err := parseAccount(decoder, tok.(xml.StartElement))

	if err != nil {
		t.Fatal(err)
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
