package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	gc "trial_balance/gnucash"
)

func accountName(a gc.Account) []string {
	var len int
	for it := &a; it.Parent != nil; it = it.Parent {
		len++
	}
	name := make([]string, len)
	for i, it := len-1, &a; it.Parent != nil; i, it = i-1, it.Parent {
		name[i] = it.Name
	}
	return name
}

func compareAccountNames(a, b []string) int {
	for i := 0; ; i++ {
		switch {
		case i == len(a) && i == len(b):
			return 0
		case i == len(a):
			return -1
		case i == len(b):
			return 1
		case a[i] != b[i]:
			return strings.Compare(a[i], b[i])
		}
	}
}

func dumpAccountNames(accounts []gc.Account) error {
	accountNames := make([][]string, len(accounts)-1)
	for i := 1; i < len(accounts); i++ {
		accountNames[i-1] = accountName(accounts[i])
	}
	slices.SortStableFunc(accountNames, compareAccountNames)
	for _, name := range accountNames {
		sep := ""
		for _, s := range name {
			fmt.Printf("%v%v", sep, s)
			sep = ":"
		}
		fmt.Println()
	}
	return nil
}

func run() error {
	book, err := gc.ReadFile(
		"c:\\Users\\Myles\\Documents\\Finance\\Accounts\\Accounts.gnucash")
	if err != nil {
		return err
	}
	err = dumpAccountNames(book.Accounts)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
