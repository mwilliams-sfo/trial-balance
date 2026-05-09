package model

import (
	"time"
)

type Guid string

type Numeric struct{ num, den int64 }

type GDate struct {
	Year  int
	Month time.Month
	Day   int
}

type Account struct {
	Name           string
	Id             Guid
	HasDescription bool
	Description    string
	HasParent      bool
	Parent         Guid
}

type CurrencyRef struct{ Space, Id string }

type Split struct {
	Id       Guid
	Value    Numeric
	Quantity Numeric
	Account  string
}

type Transaction struct {
	Id         Guid
	Currency   CurrencyRef
	Num        string
	DatePosted time.Time
}

type Book struct {
	Id           Guid
	Accounts     []Account
	Transactions []Transaction
}

type GnuCashXml []Book
