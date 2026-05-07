package model

import (
	"time"
)

type Split struct {
	Transaction *Transaction
	Value       Numeric
}

type Transaction struct {
	Account    *Account
	DatePosted time.Time
	Splits     []Split
}

type Account struct {
	Book   *Book
	Name   string
	Parent *Account
}

type Book struct {
	Accounts []Account
}
