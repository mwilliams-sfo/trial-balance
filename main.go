package main

import (
	"fmt"
	"os"
	gcio "trial_balance/gnucash/io"
)

func main() {
	book, err := gcio.Read(
		"c:\\Users\\Myles\\Documents\\Finance\\Accounts\\Accounts.gnucash")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Book id: %v\n", book.Id)
}
