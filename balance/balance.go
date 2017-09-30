package balance

import (
	"errors"
	"time"
)

// EmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const EmptyBalancesMessage = "empty Balances object"

// New creates a new Balance object
func New(date time.Time, amount Money) (b Balance, err error) {
	return Balance{Amount: amount, Date: date}, b.Validate()
}

// Balance holds the logic for a balance item.
type Balance struct {
	Date   time.Time
	Amount Money
}

// Validate checks the fields of a Balance and returns any logic errors that are present within it.
func (b Balance) Validate() error {
	if b.Date.IsZero() {
		return ZeroDate
	}
	return nil
}

// FieldError represents an error with the logic of a Balance item.
type FieldError string

// A collection of possible BalanceFieldErrors
const (
	ZeroDate = FieldError("Date of balance is zero.")
)

// Error ensures that FieldError adheres to the error interface.
func (e FieldError) Error() string {
	return string(e)
}

//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances Amount summed together.
func (bs Balances) Sum() (Money, error) {
	sum := NewMoney(0)
	var err error
	for _, b := range bs {
		sum, err = sum.Add(b.Amount)
		if err != nil {
			break
		}
	}
	return sum, err
}

// Earliest returns the Balance with the earliest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered first will be returned.
func (bs Balances) Earliest() (e Balance, err error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(EmptyBalancesMessage)
	}
	e = Balance{Date: time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.Date.Before(e.Date) {
			e = Balance(b)
		}
	}
	return
}

// Latest returns the Balance with the latest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered last will be returned.
func (bs Balances) Latest() (l Balance, err error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(EmptyBalancesMessage)
	}
	l = Balance{Date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if !l.Date.After(b.Date) {
			l = Balance(b)
		}
	}
	return
}
