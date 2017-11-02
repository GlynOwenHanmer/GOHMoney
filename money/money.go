package money

import (
	"github.com/glynternet/GOHMoney/money/currency"
)

func New(amount int64, currency currency.Code) *money {
	return &money{amount:amount, currency:currency}
}

type money struct {
	amount   int64
	currency currency.Code
}

// Money is an object representing a value and currency
type Money interface {
	Amount() int64
	Currency() currency.Code
}

// amount returns the value of the Money formed from the currency's lowest
// denominator.
// e.g. For £45.67, amount() would return 4567
func (m money) Amount() int64 {
	return m.amount
}

// currency returns the currency of the Money. If the Money has no currency, an error will also be returned.
func (m money) Currency() currency.Code {
	return m.Currency()
}

//// Moneys is a group of Moneys
//type Moneys []Money
//
//// currencies returns an array of the Currencies present within a Moneys.
//// currencies will only have one occurrence of each currency present.
//// currencies will return an error as soon as one occurs whilst retrieving
//// the currency of any Money.
//func (ms Moneys) currencies() ([]money.currency, error) {
//	var cs []money.currency
//	for _, m := range ms {
//		cur, err := m.currency()
//		if err != nil {
//			return cs, err
//		}
//		var found bool
//		for _, c := range cs {
//			if cur == c {
//				found = true
//			}
//		}
//		if found {
//			continue
//		}
//		cs = append(cs, cur)
//	}
//	return cs, nil
//}

// Validate returns an error if part of a Money is not valid.
//func (m money) Validate() error {
//	switch {
//	case m.inner == nil,
//		m.inner.currency() == nil, //todo check for
//		m.inner.currency().Code == "":
//		return ErrNoCurrency
//	}
//	return nil
//}



// SameCurrency returns true if the money and provided Money arguments all have the same currency.
// If a only one Money object is provided, m, SameCurrency will always return true with no error.
// If any Money has no currency assigned, SameCurrency will return false.
//func (m money) SameCurrency(oms ...Money) (bool, error) {
//	if len(oms) == 0 {
//		return true, nil
//	}
//	moneys := append([]Money{m}, oms...)
//	cs, err := Moneys(moneys).currencies()
//	 if error is not nil, currencies may not return all currencies
	//if err != nil {
	//	return false, err
	//}
	//return len(cs) == 1, err
//}



// Equal returns true if both Money objects are equal.
// Equal will return false and an ErrNoCurrency if either Money has no currency
// set.
//func (m Money) Equal(om Money) (bool, error) {
//	if same, err := m.SameCurrency(om); !same || err != nil {
//		return false, err
//	}
//	return m.Amount() == om.Amount(), nil
//}

// Add returns the sum of both Money objects
// If the Money objects are of different currencies, an error will be returned.
//func (m Money) Add(om Money) (Money, error) {
//	for _, mon := range []Money{m, om} {
//		if err := mon.Validate(); err != nil {
//			return Money{}, err
//		}
//	}
//	cs, err := Moneys{m, om}.currencies()
//	if err != nil {
//		return Money{}, err
//	}
//	if err := assertSameCurrency(cs...); err != nil {
//		return Money{}, err
//	}
//	return Money{inner: money.New(m.Amount()+om.Amount(), cs[0].Code)}, nil
//}

// CurrencyMismatchError is an error that is returned when an operation that
// requires multiple Money objects to be of the same currency is called.
//type CurrencyMismatchError struct {
//	A, B money.Currency
//}

// Error returns a string describing the mismatch between multiple Moneys
//func (e CurrencyMismatchError) Error() string {
//	return fmt.Sprintf("currency mismatch: %s, %s", e.A.Code, e.B.Code)
//}

// MarshalJSON marshals an Account into A json blob, returning the blob with
// any errors that occur during the marshalling.
//func (m Money) MarshalJSON() ([]byte, error) {
//	c, err := m.Currency()
//	if err != nil && err != ErrNoCurrency {
//		log.Printf("Error getting currency: %s", err)
//	}
//	type Alias Money
//	return json.Marshal(&struct {
//		Amount   int64
//		Currency string
//	}{
//		Amount:   m.Amount(),
//		Currency: c.Code,
//	})
//}
//
// UnmarshalJSON attempts to unmarshal A json blob into an Account object, returning any errors that occur during the unmarshalling.
//func (m *Money) UnmarshalJSON(data []byte) error {
//	type Alias Money
//	aux := &struct {
//		Amount   int64
//		Currency string
//	}{}
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//	m.inner = money.New(aux.Amount, aux.Currency)
//	return nil
//}
//
// String displays information about the underlying state of a Money
//func (m Money) String() string {
//	if m.inner == nil {
//		return "NIL"
//	}
//	return fmt.Sprintf("%v, %v", m.inner.Currency(), m.inner.Amount())
//}
//
//func initialiseIfRequired(m *Money) {
//	if m == nil || m.inner == nil {
//		aux, err := New(0, "")
//		if err != nil {
//			log.Printf("Error calling New: %s", err)
//		}
//		*m = *aux
//	}
//}

//func assertSameCurrency(cs ...money.Currency) error {
//	for i := 1; i < len(cs); i++ {
//		if cs[0] != cs[i] {
//			return CurrencyMismatchError{A: cs[0], B: cs[i]}
//		}
//	}
//	return nil
//}
