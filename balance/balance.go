package balance

import (
	"errors"
	"time"

	"github.com/rhymond/go-money"
)

// EmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const EmptyBalancesMessage = "empty Balances object"

// New creates a new Balance object
func New(date time.Time, amount money.Money) (b Balance, err error) {
	b = Balance{amount: amount, date: date}
	return b, b.Validate()
}

// Balance holds the logic for a balance item.
type Balance struct {
	date   time.Time
	amount money.Money
}

// Date returns the Date of the Balance
func (b Balance) Date() time.Time {
	return b.date
}

// Amount returns the Amount of the Balance
func (b Balance) Amount() money.Money {
	return b.amount
}

// Validate checks the fields of a Balance and returns any logic errors that are present within it.
func (b Balance) Validate() error {
	if b.date.IsZero() {
		return ZeroDate
	}
	return nil
}

// FieldError represents an error with the logic of a Balance item.
type FieldError string

// A collection of possible BalanceFieldErrors
const (
	ZeroDate = FieldError("date of balance is zero.")
)

// Error ensures that FieldError adheres to the error interface.
func (e FieldError) Error() string {
	return string(e)
}

//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances amount summed together.
func (bs Balances) Sum() (money.Money, error) {
	sum := new(money.Money)
	*sum = NewMoney(0)
	var err error
	for _, b := range bs {
		newAmount := b.Amount()
		sum, err = (*sum).Add(&newAmount)
		if err != nil {
			break
		}
	}
	return *sum, err
}

// Earliest returns the Balance with the earliest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered first will be returned.
func (bs Balances) Earliest() (e Balance, err error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(EmptyBalancesMessage)
	}
	e = Balance{date: time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.date.Before(e.date) {
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
	l = Balance{date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if !l.date.After(b.date) {
			l = Balance(b)
		}
	}
	return
}

// NewMoney creates a new money.Money object with currency of GBP
func NewMoney(amount int64) money.Money {
	return *money.New(amount, "GBP")
}
