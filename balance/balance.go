package balance

import (
	"errors"
	"time"

	"encoding/json"
)

// EmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const EmptyBalancesMessage = "empty Balances object"

// New creates a new Balance
func New(date time.Time, options ...Option) (b Balance, err error) {
	b.date = date
	for _, o := range options {
		err = o(&b)
		if err != nil {
			return
		}
	}
	return
}

// Balance holds the logic for a Balance item.
type Balance struct {
	date   time.Time
	amount int64
}

// Date returns the Date of the Balance
func (b Balance) Date() time.Time {
	return b.date
}

// Amount returns the amount of the Balance
func (b Balance) Amount() int64 {
	return b.amount
}

// Equal returns true if two Balance objects are logically equal
func (b Balance) Equal(ob Balance) bool {
	return b.amount == ob.Amount() && b.Date().Equal(ob.Date())
}

type jsonHelper struct {
	Date   time.Time
	Amount int64
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (b Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonHelper{
		Date:   b.date,
		Amount: b.amount,
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
func (b *Balance) UnmarshalJSON(data []byte) error {
	aux := new(jsonHelper)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	b.date = aux.Date
	b.amount = aux.Amount
	var returnErr error
	if err := b.Validate(); err != nil {
		returnErr = err
	}
	return returnErr
}

// FieldError represents an error with the logic of a Balance item.
type FieldError string

// A collection of possible BalanceFieldErrors
const (
	ZeroDate = FieldError("date of Balance is zero.")
)

// Error ensures that FieldError adheres to the error interface.
func (e FieldError) Error() string {
	return string(e)
}

//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances summed together.
func (bs Balances) Sum() (s int64) {
	for _, b := range bs {
		s += b.amount
	}
	return
}

// Earliest returns the Balance with the earliest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered first will be returned.
func (bs Balances) Earliest() (e Balance, err error) {
	if len(bs) == 0 {
		return e, errors.New(EmptyBalancesMessage)
	}
	e = Balance{date: time.Date(2000000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.date.Before(e.date) {
			e = b
		}
	}
	return
}

// Latest returns the Balance with the latest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered last will be returned.
func (bs Balances) Latest() (l Balance, err error) {
	if len(bs) == 0 {
		return l, errors.New(EmptyBalancesMessage)
	}
	l = Balance{date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if !l.date.After(b.date) {
			l = b
		}
	}
	return
}
