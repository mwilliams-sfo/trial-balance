package main

import (
	"fmt"
	"os"
	gc "trial_balance/gnucash"
)

func run() error {
	book, err := gc.ReadFile(
		"c:\\Users\\Myles\\Documents\\Finance\\Accounts\\Accounts.gnucash")
	if err != nil {
		return err
	}
	fmt.Println(book)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
