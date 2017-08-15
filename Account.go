package GOHMoney

import (
	"encoding/json"
	"strings"
	"github.com/lib/pq"
	"time"
)

// An Account holds the logic for an account.
type Account struct {
	Name       string
	timeRange TimeRange
}

// Start returns the start time of the Account's TimeRange
func (account Account) Start() time.Time {
	return account.timeRange.Start.Time
}

// End returns the end time of the Account's TimeRange. The End time is returned as a NullTime object as it may not be a valid time when the Account is not yet closed.
func (account Account) End() pq.NullTime {
	return account.timeRange.End
}

// IsOpen return true if the Account is open.
func (account Account) IsOpen() bool {
	return !account.timeRange.End.Valid
}

// String() ensures that Account conforms to the Stringer interface.
func (account Account) String() string {
	jsonBytes, err := json.Marshal(account)
	var ret string
	if err != nil {ret = "Unable to form Account string."} else {ret = string(jsonBytes)}
	return ret
}

// Validate checks the state of an account to see if it is has any logical errors. Validate returns a set of errors representing errors with different fields of the account.
func (account Account) Validate() AccountFieldError {
	var fieldErrorDescriptions []string
	if len(strings.TrimSpace(account.Name)) == 0 {
		fieldErrorDescriptions = append(fieldErrorDescriptions, EmptyNameError)
	}
	if err := account.timeRange.Validate(); err != nil {
		fieldErrorDescriptions = append(fieldErrorDescriptions, err.Error())
	}
	if !account.timeRange.Start.Valid || account.timeRange.Start.Time.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroDateOpenedError)
	}
	if account.timeRange.End.Valid && account.timeRange.End.Time.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroValidDateClosedError)
	}
	if len(fieldErrorDescriptions) > 0 {
		return AccountFieldError(fieldErrorDescriptions)
	}
	return nil
}

// ValidateBalance validates a Balance against an Account.
// ValidateBalance returns any logical errors between the Account and the Balance.
// ValidateBalance first attempts to validate the Account as an entity by itself. If there are any errors with the Account, these errors are returned and the Balance is not attempted to be validated against the account.
// If the Date of the Balance is outside of the TimeRange of the Account, a BalanceDateOutOfAccountTimeRange will be returned.
func (account Account) ValidateBalance(balance Balance) error {
	if err := account.Validate(); err != nil {
		return err
	}
	if err := balance.Validate(); err != nil {
		return err
	}
	if !account.timeRange.Contains(balance.Date) {
		return BalanceDateOutOfAccountTimeRange{
			BalanceDate:balance.Date,
			AccountTimeRange:account.timeRange,
		}
	}
	return nil
}

// NewAccount creates a new Account object with a Valid Start time and returns it, also returning any logical errors with the newly created account.
func NewAccount(name string, opened time.Time, closed pq.NullTime) (Account, error) {
	newAccount := Account{
		Name: name,
		timeRange: TimeRange{
			Start: pq.NullTime{
				Valid: true,
				Time: opened,
			},
			End: closed,
		},
	}
	var err error
	if accountErr := newAccount.Validate(); len(accountErr) > 0 {
		err = accountErr
	}
	return newAccount, err
}

// Accounts holds multiple Account items.
type Accounts []Account