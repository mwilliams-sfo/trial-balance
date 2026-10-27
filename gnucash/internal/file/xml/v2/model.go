package v2

type Account struct {
	Name   string
	Id     string
	Parent string
}

type Transaction struct {
	Id string
}

type Book struct {
	Id           string
	Accounts     []*Account
	Transactions []*Transaction
}
