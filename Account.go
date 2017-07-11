package GOHMoney

import (
	"encoding/json"
	"time"
	"github.com/lib/pq"
	"strings"
)

// An Account holds the logic for an account.
type Account struct {
	Name       string	`json:"name"`
	DateOpened time.Time	`json:"date_opened"`
	DateClosed pq.NullTime	`json:"date_closed"`
}

// IsOpen return true if the Account is open.
func (account Account) IsOpen() bool {
	return !account.DateClosed.Valid
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
	if account.DateOpened.IsZero() {
		fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroDateOpenedError)
	}
	if account.DateClosed.Valid {
		if account.DateClosed.Time.IsZero() {
			fieldErrorDescriptions = append(fieldErrorDescriptions, ZeroValidDateClosedError)
		} else if account.DateClosed.Time.Before(account.DateOpened) {
			fieldErrorDescriptions = append(fieldErrorDescriptions, DateClosedBeforeDateOpenedError)
		}
	}
	return AccountFieldError(fieldErrorDescriptions)
}

// Accounts holds multiple Account items.
type Accounts []Account