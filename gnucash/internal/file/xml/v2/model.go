package v2

import (
	"time"
)

type Account struct {
	Name   string
	Id     string
	Parent string
}

type Transaction struct {
	Id         string
	DatePosted time.Time
}

type Book struct {
	Id           string
	Accounts     []Account
	Transactions []Transaction
}
