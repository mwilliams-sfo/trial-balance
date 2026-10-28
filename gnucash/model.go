package gnucash

import (
	"time"
	num "trial_balance/gnucash/numeric"
)

type Account struct {
	Name   string
	Id     string
	Parent *Account
}

type Split struct {
	Id    string
	Value num.Numeric
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
