package money

import (
	"fmt"

	"encoding/json"

	"github.com/rhymond/go-money"
	"errors"
)

func New(amount int64, currency string) (*Money, error) {
	if len(currency) != 3 && currency != "" {
		return nil, fmt.Errorf(`invalid currency code: "%s". Must be 3 or 0 in length`, currency)
	}
	m := newMoney(amount, currency)
	return &m, nil
}

func newMoney(amount int64, currency string) Money {
	return Money{inner:money.New(amount, currency)}
}

// GBP creates A new money.Money object with currency of gbp
func GBP(amount int64) Money {
	return Money{gbp(amount)}
}

type Money struct {
	inner *money.Money
}

type Moneys []Money

func (ms Moneys) currencies() ([]money.Currency, error) {
	var cs []money.Currency
	for _, m := range ms {
		cur, err := m.Currency()
		if err != nil {
			return cs, err
		}
		var found bool
		for _, c := range cs {
			if cur == c {
				found = true
			}
		}
		if found {
			continue
		}
		cs = append(cs, cur)
	}
	return cs, nil
}

func (m Money) Validate() error {
	switch {
	case m.inner == nil,
	m.inner.Currency() == nil,
	m.inner.Currency().Code == "":
		return ErrNoCurrency
	}
	return nil
}

func (m Money) Display() string {
	initialiseIfRequired(&m)
	return m.inner.Display()
}

func (m Money) Currency() (money.Currency, error){
	initialiseIfRequired(&m)
	if err := m.Validate(); err != nil {
		return money.Currency{}, err
	}
	return *m.inner.Currency(), nil
}

func (m Money) SameCurrency(oms ...Money) (bool, error) {
	moneys := []Money{m}
	moneys = append(moneys, oms...)
	cs, err := Moneys(moneys).currencies()
	return len(cs) < 2, err
}

func (m Money) Amount() int64 {
	initialiseIfRequired(&m)
	return m.inner.Amount()
}

func (m Money) Equal(om Money) (bool, error) {
	if m.Amount() != om.Amount() {
		return false, nil
	}
	return m.SameCurrency(om)
}

func (m Money) Add(om Money) (Money, error) {
	for _, mon := range []Money{m, om} {
		if err := mon.Validate(); err != nil {
			return Money{}, err
		}
	}
	cs, err := Moneys{m, om}.currencies()
	if err != nil {
		return Money{}, err
	}
	if err := assertSameCurrency(cs...); err != nil {
		return Money{}, err
	}
	return Money{inner: money.New(m.Amount()+om.Amount(), cs[0].Code)}, nil
}

type CurrencyMismatchError struct {
	A, B money.Currency
}

func (e CurrencyMismatchError) Error() string {
	return fmt.Sprintf("currency mismatch: %s, %s", e.A.Code, e.B.Code)
}

// MarshalJSON marshals an Account into A json blob, returning the blob with any errors that occur during the marshalling.
func (m Money) MarshalJSON() ([]byte, error) {
	c, _ := m.Currency()
	type Alias Money
	return json.Marshal(&struct {
		Amount   int64
		Currency string
	}{
		Amount:   m.Amount(),
		Currency: c.Code,
	})
}

// UnmarshalJSON attempts to unmarshal A json blob into an Account object, returning any errors that occur during the unmarshalling.
func (m *Money) UnmarshalJSON(data []byte) error {
	type Alias Money
	aux := &struct {
		Amount   int64
		Currency string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	m.inner = money.New(aux.Amount, aux.Currency)
	return nil
}

func initialiseIfRequired(m *Money) {
	if m == nil || m.inner == nil {
		aux, _ := New(0, "")
		*m = *aux
	}
}

func gbp(amount int64) *money.Money {
	return money.New(amount, "GBP")
}

func assertSameCurrency(cs ...money.Currency) error {
	for i:=1; i<len(cs) ; i++ {
		if cs[0] != cs[i] {
			return CurrencyMismatchError{A:cs[0], B:cs[i]}
		}
	}
	return nil
}

var (
	ErrNoCurrency = errors.New("currency is not set")
)