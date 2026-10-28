package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	gc "trial_balance/gnucash"
	num "trial_balance/gnucash/numeric"
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

func trialBalances(w io.Writer, transactions []gc.Transaction) {
	sorted := make([]gc.Transaction, len(transactions))
	copy(sorted, transactions)
	slices.SortStableFunc(sorted, func(a, b gc.Transaction) int {
		aTime := a.DatePosted
		bTime := b.DatePosted
		switch {
		case aTime.Before(bTime):
			return -1
		case bTime.Before(aTime):
			return 1
		default:
			return 0
		}
	})

	debits := num.Int(0)
	credits := num.Int(0)
	for _, trn := range sorted {
		for _, split := range trn.Splits {
			sign := split.Value.Sign()
			switch {
			case sign > 0:
				debits = debits.Plus(split.Value)
			case sign < 0:
				credits = credits.Plus(split.Value.Negate())
			}
		}
		fmt.Fprintf(w, "%v %20.02f %20.02f\n",
			trn.DatePosted.Local().Format("2006-01-02 15:04:05.999999999"),
			debits.Float64(), credits.Float64())
	}
}

func run() error {
	book, err := gc.ReadFile(
		"c:\\Users\\Myles\\Documents\\Finance\\Accounts\\Accounts.gnucash")
	if err != nil {
		return err
	}
	trialBalances(os.Stdout, book.Transactions)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
