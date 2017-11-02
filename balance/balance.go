package balance

import (
	"errors"
	"time"
)

// ErrEmptyBalancesMessage is the error message used when a Balances object contains no balance items.
const ErrEmptyBalancesMessage = "empty Balances object"

// New creates a new balance
func New(date time.Time, options ...Option) (b *Balance, err error) {
	bb := balance{date: date}
	for _, o := range options {
		err = o(&bb)
		if err != nil {
			return
		}
	}
	return
}

// Balance holds the logic for a balance item.
type Balance interface {
	Date() time.Time
	Amount() int
}

// balance holds the logic for a balance item.
type balance struct {
	date   time.Time
	amount int
}

func (b balance) Date() time.Time {
	return b.date
}

func (b balance) Amount() int {
	return b.amount
}

// Equal returns true if two balance objects are logically equal
func (b balance) Equal(ob Balance) bool {
	return b.amount == ob.Amount() && b.date.Equal(ob.Date())
}

//type jsonHelper struct {
//	date   time.Time
//	amount int64
//}
//
//// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
//func (b balance) MarshalJSON() ([]byte, error) {
//	return json.Marshal(&jsonHelper{
//		date:   b.date,
//		amount: b.amount,
//	})
//}
//
//// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
//func (b *balance) UnmarshalJSON(data []byte) error {
//	aux := new(jsonHelper)
//	if err := json.Unmarshal(data, aux); err != nil {
//		return err
//	}
//	b.date = aux.date
//	b.amount = aux.amount
//	return nil
//}

//Balances holds multiple balance items.
type Balances []Balance

// Sum returns the value of all of the balances summed together.
func (bs Balances) Sum() (s int) {
	for _, b := range bs {
		s += b.Amount()
	}
	return
}

// Earliest returns the balance with the earliest date contained in a Balances set.
// If multiple balance object have the same date, the balance encountered first will be returned.
func (bs Balances) Earliest() (e *Balance, err error) {
	if len(bs) == 0 {
		return e, errors.New(ErrEmptyBalancesMessage)
	}
	e = new(Balance)
	*e = balance{date: time.Date(2000000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.Date().Before((*e).Date()) {
			*e = b
		}
	}
	return
}

// Latest returns the balance with the latest date contained in a Balances set.
// If multiple balance object have the same date, the balance encountered last will be returned.
func (bs Balances) Latest() (l *Balance, err error) {
	if len(bs) == 0 {
		return l, errors.New(ErrEmptyBalancesMessage)
	}
	l = new(Balance)
	*l = balance{date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if !(*l).Date().After(b.Date()) {
			*l = b
		}
	}
	return
}
