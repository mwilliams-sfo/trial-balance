package gnucash

import (
	"time"
)

type Account struct {
	Name   string
	Id     string
	Parent *Account
}

type Split struct {
	Id               string
	Amount, Quantity Numeric
	Account          string
}

type Transaction struct {
	Id         string
	DatePosted time.Time
	Splits     []Split
}

type Book struct {
	Id           string
	Accounts     []Account
	Transactions []Transaction
}
