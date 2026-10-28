package gnucash

import (
	"fmt"
)

type Numeric struct {
	Num, Den int64
}

func Int(i int64) Numeric {
	return Numeric{Num: i, Den: 1}
}

func (n Numeric) Float64() float64 {
	n = n.Reduce()
	return float64(n.Num) / float64(n.Den)
}

func (n Numeric) Sign() int {
	if n.Den < 0 {
		n.Num = -n.Num
	}
	switch {
	case n.Num == 0:
		return 0
	case n.Num < 0:
		return -1
	default:
		return 1
	}
}

func Compare(a, b Numeric) int {
	if a.Den != b.Den {
		a = a.Reduce()
		b = b.Reduce()
		gcd := gcd(a.Den, b.Den)
		a.Num *= b.Den / gcd
		b.Num *= a.Den / gcd
	}
	switch {
	case a.Num < b.Num:
		return -1
	case a.Num > b.Num:
		return 1
	default:
		return 0
	}
}

func gcd(a, b int64) int64 {
	switch {
	case a == 0 && b == 0:
		return 1
	case a == 0:
		return b
	case b == 0:
		return a
	}
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a < b {
		a, b = b, a
	}
	for {
		a, b = b, a%b
		if b == 0 {
			return a
		}
	}
}

func (n Numeric) Negate() Numeric {
	n.Num = -n.Num
	return n
}

func (n Numeric) Reduce() Numeric {
	if n.Den < 0 {
		n.Num, n.Den = -n.Num, -n.Den
	}
	gcd := gcd(n.Num, n.Den)
	if gcd == 0 {
		fmt.Printf("About to panic for %v\n", n)
	}

	n.Num /= gcd
	n.Den /= gcd
	return n
}

func (a Numeric) Plus(b Numeric) Numeric {
	switch {
	case b.Num == 0:
	case a.Num == 0:
		a = b
	case a.Den == b.Den:
		a.Num += b.Num
	default:
		a = a.Reduce()
		b = b.Reduce()
		gcd := gcd(a.Den, b.Den)
		a.Num = a.Num*(b.Den/gcd) + b.Num*(a.Den/gcd)
		a.Den = a.Den / gcd * b.Den
	}
	return a
}
