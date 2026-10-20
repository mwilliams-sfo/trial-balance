package main

import (
	"fmt"
	"os"
	gc "trial_balance/gnucash"
)

func run() error {
	_, err := gc.ReadFile(
		"c:\\Users\\Myles\\Documents\\Finance\\Accounts\\Accounts.gnucash")
	return err
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
