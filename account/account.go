package account

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoney/balance"
)

// New creates a new Account object with a Valid Start time and returns it, also returning any logical errors with the newly created account.
func New(name string, opened time.Time, closed GOHMoney.NullTime) (a Account, err error) {
	a = Account{
		Name: name,
		timeRange: GOHMoney.TimeRange{
			Start: GOHMoney.NullTime{
				Valid: true,
				Time:  opened,
			},
			End: closed,
		},
	}
	if aErr := a.Validate(); len(aErr) > 0 {
		err = aErr
	}
	return a, err
}

// An Account holds the logic for an account.
type Account struct {
	Name      string
	timeRange GOHMoney.TimeRange
}

// Start returns the start time that the Account opened.
func (a Account) Start() time.Time {
	return a.timeRange.Start.Time
}

// End returns the a NullTime object that is Valid if the account has been closed.
func (a Account) End() GOHMoney.NullTime {
	return a.timeRange.End
}

// IsOpen return true if the Account is open.
func (a Account) IsOpen() bool {
	return !a.timeRange.End.Valid
}

// String() ensures that Account conforms to the Stringer interface.
func (a Account) String() string {
	jsonBytes, err := json.Marshal(a)
	var ret string
	if err != nil {
		ret = "Unable to form Account string."
	} else {
		ret = string(jsonBytes)
	}
	return ret
}

// Validate checks the state of an account to see if it is has any logical errors. Validate returns a set of errors representing errors with different fields of the account.
func (a Account) Validate() FieldError {
	var fieldErrorDescriptions []string
	if len(strings.TrimSpace(a.Name)) == 0 {
		fieldErrorDescriptions = append(fieldErrorDescriptions, EmptyNameError)
	}
	if err := a.timeRange.Validate(); err != nil {
		fieldErrorDescriptions = append(fieldErrorDescriptions, err.Error())
	}
	if a.Start().IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroDateOpenedError)
	}
	if a.End().Valid && a.End().Time.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroValidDateClosedError)
	}
	if len(fieldErrorDescriptions) > 0 {
		return FieldError(fieldErrorDescriptions)
	}
	return nil
}

// ValidateBalance validates a Balance against an Account.
// ValidateBalance returns any logical errors between the Account and the Balance.
// ValidateBalance first attempts to validate the Account as an entity by itself. If there are any errors with the Account, these errors are returned and the Balance is not attempted to be validated against the account.
// If the Date of the Balance is outside of the TimeRange of the Account, a DateOutOfAccountTimeRange will be returned.
func (a Account) ValidateBalance(b balance.Balance) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if err := b.Validate(); err != nil {
		return err
	}
	if !a.timeRange.Contains(b.Date) && (!a.End().Valid || !a.End().Time.Equal(b.Date)) {
		return balance.DateOutOfAccountTimeRange{
			BalanceDate:      b.Date,
			AccountTimeRange: a.timeRange,
		}
	}
	return nil
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	return json.Marshal(&struct {
		*Alias
		Start time.Time
		End   GOHMoney.NullTime
	}{
		Alias: (*Alias)(&a),
		Start: a.Start(),
		End:   a.End(),
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
func (a *Account) UnmarshalJSON(data []byte) error {
	type Alias Account
	aux := &struct {
		Start time.Time
		End   GOHMoney.NullTime
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.timeRange = GOHMoney.TimeRange{
		Start: GOHMoney.NullTime{Valid: true, Time: aux.Start},
		End:   aux.End,
	}
	var returnErr error
	if err := a.Validate(); err != nil {
		returnErr = err
	}
	return returnErr
}

// Equal returns true if both accounts a and b are logically the same.
func (a Account) Equal(b Account) bool {
	switch {
	case a.Name != b.Name:
		return false
	case !a.timeRange.Equal(b.timeRange):
		return false
	}
	return true
}

// Accounts holds multiple Account items.
type Accounts []Account
