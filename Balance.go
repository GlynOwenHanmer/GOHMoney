package GOHMoney

import (
	"time"
	"errors"
)

const EmptyBalancesMessage  = "Empty Balances Object"

// Balance holds the logic for a balance item.
type Balance struct {
	Date   time.Time
	Amount float32
}

// Validate checks the fields of a Balance and returns any logic errors that are present within it.
func (balance Balance) Validate() error {
	if balance.Date.IsZero() {
		return BalanceZeroDate
	}
	return nil
}

// BalanceFieldError represents an error with the logic of a Balance item.
type BalanceFieldError string

// A collection of possible BalanceFieldErrors
const (
	BalanceZeroDate = BalanceFieldError("Date of balance is zero.")
)

// Error ensures that BalanceFieldError adheres to the error interface.
func (e BalanceFieldError) Error() string {
	return string(e)
}


//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances amount summed together.
func (bs Balances) Sum() float32 {
	var sum float32
	for _, b := range bs {
		sum += b.Amount
	}
	return sum
}

// Earliest returns the Balance with the earliest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered first will be returned.
func (bs Balances) Earliest() (Balance, error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(EmptyBalancesMessage)
	}
	earliest := Balance{Date: time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.Date.Before(earliest.Date) {
			earliest = Balance(b)
		}
	}
	return earliest, nil
}

// Latest returns the Balance with the latest date contained in a Balances set.
// If multiple Balance object have the same date, the Balance encountered last will be returned.
func (bs Balances) Latest() (Balance, error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(EmptyBalancesMessage)
	}
	latest := Balance{Date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.Date.After(latest.Date) || b.Date.Equal(latest.Date) {
			latest = Balance(b)
		}
	}
	return latest, nil
}