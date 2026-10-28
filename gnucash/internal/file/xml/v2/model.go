package v2

import (
	"time"
)

type Numeric struct {
	Num, Den int64
}

type Account struct {
	Name   string
	Id     string
	Parent string
}

type Split struct {
	Id    string
	Value Numeric
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
