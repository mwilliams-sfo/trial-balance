package model

type Numeric interface {
	Num() int64
	Den() int64
	Add(b Numeric) Numeric
	Float64() float64
}

type numeric struct{ num, den int64 }

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a < b {
		t := b
		b = a
		a = t
	}
	for {
		r := a % b
		if r == 0 {
			return b
		}
		a = b
		b = r
	}
}

func (n numeric) scale(m int64) numeric {
	return numeric{
		num: n.num * m,
		den: n.den * m,
	}
}

func (n numeric) normalize() numeric {
	d := gcd(n.num, n.den)
	if n.den < 0 {
		d = -d
	}
	return numeric{
		num: n.num / d,
		den: n.den / d,
	}
}

func NewNumeric(n, d int64) Numeric {
	return numeric{num: n, den: d}.normalize()
}

func (n numeric) Num() int64 { return n.num }

func (n numeric) Den() int64 { return n.den }

func (n numeric) Add(b Numeric) Numeric {
	n = n.normalize()
	m := n.den / gcd(n.den, b.Den()) * b.Den()
	n = n.scale(m / n.den)
	n.num += b.Num() * (m / b.Den())
	return n.normalize()
}

func (n numeric) Float64() float64 {
	return float64(n.num) / float64(n.den)
}
