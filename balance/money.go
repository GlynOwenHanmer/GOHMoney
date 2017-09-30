package balance

import (
	"github.com/rhymond/go-money"
)

// NewMoney returns a Money object with the default currency
func NewMoney(amount int64) Money {
	return Money(*money.New(amount, "GBP"))
}

// Money is a wrapper around a money.Money object as accessing pointer methods on the money.Money object was proving very difficult in some cases
type Money money.Money

// Add returns the value of two Money object added together.
func (m Money) Add(m2 Money) (Money, error) {
	a, b := moneysToCompare(m, m2)
	ret, err := (&a).Add(&b)
	if err != nil {
		return NewMoney(0), err
	}
	return Money(*ret), nil
}

// Equals returns true if 2 Money objects are identical
func (m Money) Equals(m2 Money) (bool, error) {
	a, b := moneysToCompare(m, m2)
	return (&a).Equals(&b)
}

func moneysToCompare(a, b Money) (c, d money.Money) {
	return money.Money(a), money.Money(b)
}
