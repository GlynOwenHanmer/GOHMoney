package GOHMoney

import (
	"encoding/json"
	"strings"
	"github.com/lib/pq"
	"time"
)

// An Account holds the logic for an account.
type Account struct {
	Name       string	`json:"name"`
	TimeRange
}

// IsOpen return true if the Account is open.
func (account Account) IsOpen() bool {
	return !account.TimeRange.End.Valid
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
	if err := account.TimeRange.Validate(); err != nil {
		fieldErrorDescriptions = append(fieldErrorDescriptions, err.Error())
	}
	if !account.TimeRange.Start.Valid || account.TimeRange.Start.Time.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroDateOpenedError)
	}
	if account.TimeRange.End.Valid && account.TimeRange.End.Time.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroValidDateClosedError)
	}
	return AccountFieldError(fieldErrorDescriptions)
}

// ValidateBalance validates a Balance against an Account.
// ValidateBalance returns any logical errors between the Account and the Balance.
// ValidateBalance first attempts to validate the Account as an entity by itself. If there are any errors with the Account, these errors are returned and the Balance is not attempted to be validated against the account.
// If the Date of the Balance is outside of the TimeRange of the Account, a BalanceDateOutOfAccountTimeRange will be returned.
func (account Account) ValidateBalance(balance Balance) error {
	if err := account.Validate(); err != nil {
		return err
	}
	if !account.TimeRange.Contains(balance.Date) {
		return BalanceDateOutOfAccountTimeRange{
			BalanceDate:balance.Date,
			AccountTimeRange:account.TimeRange,
		}
	}
	return nil
}

// NewAccount creates a new Account object with a Valid Start time and returns it, also returning any logical errors with the newly created account.
func NewAccount(name string, opened time.Time, closed pq.NullTime) (Account, error) {
	newAccount := Account{
		Name: name,
		TimeRange: TimeRange{
			Start: pq.NullTime{
				Valid: true,
				Time: opened,
			},
			End: closed,
		},
	}
	return newAccount, newAccount.Validate()
}

// Accounts holds multiple Account items.
type Accounts []Account