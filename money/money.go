package money

import (
	"fmt"

	"encoding/json"

	"errors"
	"log"

	"github.com/rhymond/go-money"
	"strings"
)

// New creates a new Money and returns a pointer to it,
// along with any errors associated with the Money whilst creating it.
// New will always convert any currency code to upper case.
func New(amount int64, currency string) (*Money, error) {
	//todo ensure that currency contains no whitespace
	currency = strings.ToUpper(currency)
	if len(currency) != 3 && currency != "" {
		return nil, fmt.Errorf(`invalid currency code: "%s". Must be 3 or 0 in length`, currency)
	}
	m := newMoney(amount, currency)
	return &m, nil
}

func newMoney(amount int64, currency string) Money {
	return Money{inner: money.New(amount, currency)}
}

// GBP creates A new money.Money object with currency of gbp
func GBP(amount int64) Money {
	return Money{gbp(amount)}
}

// Money is an object representing a value and currency
type Money struct {
	inner *money.Money
}

// Moneys is a group of Moneys
type Moneys []Money

// currencies returns an array of the Currencies present within a Moneys.
// currencies will only have one occurrence of each Currency present.
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

// Validate returns an error if part of a Money is not valid.
func (m Money) Validate() error {
	switch {
	case m.inner == nil,
		m.inner.Currency() == nil,			//todo check for

		m.inner.Currency().Code == "":
		return ErrNoCurrency
	}
	return nil
}

// Display returns a string representing the value of the Money, including currency symbol.
func (m Money) Display() string {
	initialiseIfRequired(&m)
	return m.inner.Display()
}

// Currency returns the Currency of the Money. If the Money has no Currency, an error will also be returned.
func (m Money) Currency() (money.Currency, error) {
	initialiseIfRequired(&m)
	if err := m.Validate(); err != nil {
		return money.Currency{}, err
	}
	return *m.inner.Currency(), nil
}

// SameCurrency returns true if the money and provided Money arguments all have the same Currency.
func (m Money) SameCurrency(oms ...Money) (bool, error) {
	moneys := []Money{m}
	moneys = append(moneys, oms...)
	cs, err := Moneys(moneys).currencies()
	return len(cs) < 2, err
}

// Amount returns the value of the Money formed from the currency's lowest denominator.
// e.g. For £45.67, Amount() would return 4567
func (m Money) Amount() int64 {
	initialiseIfRequired(&m)
	return m.inner.Amount()
}

// Equal returns true if both Money objects are equal.
func (m Money) Equal(om Money) (bool, error) {
	if m.Amount() != om.Amount() {
		return false, nil
	}
	return m.SameCurrency(om)
}

// Add returns the sum of both Money objects
// If the Money objects are of different currencies, an error will be returned.
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

// CurrencyMismatchError is an error that is returned when an operation that requires
// multiple Money objects to be of the same currency is called.
type CurrencyMismatchError struct {
	A, B money.Currency
}

// Error returns a string describing the mismatch between multiple Moneys
func (e CurrencyMismatchError) Error() string {
	return fmt.Sprintf("currency mismatch: %s, %s", e.A.Code, e.B.Code)
}

// MarshalJSON marshals an Account into A json blob, returning the blob with any errors that occur during the marshalling.
func (m Money) MarshalJSON() ([]byte, error) {
	c, err := m.Currency()
	if err != nil && err != ErrNoCurrency {
		log.Printf("Error getting currency: %s", err)
	}
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
		aux, err := New(0, "")
		if err != nil {
			log.Printf("Error calling New: %s", err)
		}
		*m = *aux
	}
}

func gbp(amount int64) *money.Money {
	return money.New(amount, "GBP")
}

func assertSameCurrency(cs ...money.Currency) error {
	for i := 1; i < len(cs); i++ {
		if cs[0] != cs[i] {
			return CurrencyMismatchError{A: cs[0], B: cs[i]}
		}
	}
	return nil
}

var (
	// ErrNoCurrency is returned when a Money's currency is fetched but it has not been assigned a currency.
	ErrNoCurrency = errors.New("currency is not set")
)
