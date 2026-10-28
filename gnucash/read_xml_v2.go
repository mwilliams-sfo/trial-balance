package gnucash

import (
	"encoding/xml"
	"fmt"
	xmlv2 "trial_balance/gnucash/internal/file/xml/v2"
)

func decodeXmlV2(d *xml.Decoder) (Book, error) {
	var zero Book

	v2Book, err := xmlv2.Decode(d)
	if err != nil {
		return zero, err
	}

	book := Book{Id: v2Book.Id}
	book.Accounts = make([]Account, len(v2Book.Accounts))
	accountMap := make(map[string]*Account)
	for i := range len(v2Book.Accounts) {
		book.Accounts[i].Name = v2Book.Accounts[i].Name
		book.Accounts[i].Id = v2Book.Accounts[i].Id
		accountMap[book.Accounts[i].Id] = &book.Accounts[i]
	}
	for i, v2Account := range v2Book.Accounts {
		if v2Account.Parent != "" {
			var ok bool
			book.Accounts[i].Parent, ok = accountMap[v2Account.Parent]
			if !ok {
				return zero, fmt.Errorf("Invalid account parent ID")
			}
		}
	}

	book.Transactions = make([]Transaction, len(v2Book.Transactions))
	for i, v2Trn := range v2Book.Transactions {
		book.Transactions[i] = Transaction{
			Id:         v2Trn.Id,
			DatePosted: v2Trn.DatePosted,
			Splits:     make([]Split, len(v2Trn.Splits)),
		}
		for j, v2Split := range v2Trn.Splits {
			book.Transactions[i].Splits[j] = Split{
				Id:    v2Split.Id,
				Value: v2Split.Value,
			}
		}
	}

	return book, nil
}
