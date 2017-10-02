package money

import (
	"fmt"

	"github.com/rhymond/go-money"
)

// New creates a new money.Money object with currency of GBP
func New(amount int64) Money {
	return Money{defaultMoney(amount)}
}

type Money struct {
	inner *money.Money
}

func (m Money) Currency() money.Currency {
	if m.inner == nil {
		m.inner = defaultMoney(0)
	}
	return *m.inner.Currency()
}

func (m Money) Amount() int64 {
	if m.inner == nil {
		return 0
	}
	return m.inner.Amount()
}

func (m Money) Equal(om Money) (bool, error) {
	if m.Amount() != om.Amount() {
		return false, nil
	}
	if err := assertSameCurrency(m.Currency(), om.Currency()); err != nil {
		return false, err
	}
	return true, nil
}

func (m Money) Add(om Money) (Money, error) {
	if err := assertSameCurrency(m.Currency(), om.Currency()); err != nil {
		return Money{}, err
	}
	return Money{inner: money.New(m.Amount()+om.Amount(), m.Currency().Code)}, nil
}

func defaultMoney(amount int64) *money.Money {
	return money.New(amount, "GBP")
}

func assertSameCurrency(c1, c2 money.Currency) error {
	if c1 != c2 {
		return fmt.Errorf("currency mismatch: %s, %s", c1.Code, c2.Code)
	}
	return nil
}
