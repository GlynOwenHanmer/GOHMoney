package balance

import (
	"errors"
	"time"

	"encoding/json"

	"github.com/GlynOwenHanmer/GOHMoney/money"
)

// EmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const EmptyBalancesMessage = "empty Balances object"

// New creates a new Balance
func New(date time.Time, amount money.Money) (Balance, error) {
	if err := amount.Validate(); err != nil {
		return Balance{}, err
	}
	b := Balance{money: amount, date: date}
	return b, b.Validate()
}

// Balance holds the logic for a balance item.
type Balance struct {
	date  time.Time
	money money.Money
}

// Date returns the Date of the Balance
func (b Balance) Date() time.Time {
	return b.date
}

// Money returns the Money of the Balance
func (b Balance) Money() money.Money {
	return b.money
}

// Equal returns true if two Balance objects are logically equal
func (b Balance) Equal(ob Balance) bool {
	if amountEqual, err := b.Money().Equal(ob.Money()); !amountEqual || !b.Date().Equal(ob.Date()) || err != nil {
		return false
	}
	return true
}

// Validate checks the fields of a Balance and returns any logic errors that are present within it.
func (b Balance) Validate() error {
	if b.date.IsZero() {
		return ZeroDate
	}
	return nil
}

type jsonHelper struct {
	Date  time.Time
	Money money.Money
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (b Balance) MarshalJSON() ([]byte, error) {
	type Alias Balance
	return json.Marshal(&jsonHelper{
		Date:  b.Date(),
		Money: b.Money(),
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
func (b *Balance) UnmarshalJSON(data []byte) error {
	type Alias Balance
	aux := new(jsonHelper)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	b.date = aux.Date
	b.money = aux.Money
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
	ZeroDate = FieldError("date of balance is zero.")
)

// Error ensures that FieldError adheres to the error interface.
func (e FieldError) Error() string {
	return string(e)
}

//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances money summed together.
func (bs Balances) Sum() (money.Money, error) {
	var initialised bool
	var s money.Money
	var err error
	if len(bs) < 1 {
		return money.Money{}, nil
	}
	for _, b := range bs {
		if !initialised {
			s = b.Money()
			initialised = true
			continue
		}
		s, err = s.Add(b.Money())
		if err != nil {
			break
		}
	}
	return s, err
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
			e = b
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
			l = b
		}
	}
	return
}
