package gnucash

import (
	"fmt"
)

type Id interface {
	Equals(other Id) bool
}

type Guid string

type Transaction struct {
	Id Id
}

type Book struct {
	Id           Id
	Transactions []Transaction
}

func (guid Guid) Equals(other Id) bool {
	oguid, ok := other.(Guid)
	return ok && guid == oguid
}

func (guid Guid) String() string {
	return fmt.Sprintf("GUID{%s}", string(guid))
}
