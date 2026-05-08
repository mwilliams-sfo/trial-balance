package names

import (
	"encoding/xml"
)

func makeName(space, local string) func() xml.Name {
	name := xml.Name{space, local}
	return func() xml.Name { return name }
}

var (
	CountData = makeName("gnc", "count-data")
	Book      = makeName("gnc", "book")
	Account   = makeName("gnc", "account")

	CdTypeAttr = makeName("cd", "type")

	BookId = makeName("book", "id")

	AccountName   = makeName("act", "name")
	AccountId     = makeName("act", "id")
	AccountParent = makeName("act", "parent")

	TimeDate = makeName("ts", "date")
	TimeNs   = makeName("ts", "ns")

	GDate      = makeName("", "gdate")
	IdTypeAttr = makeName("", "type")
)
