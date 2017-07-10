package GOHMoney

import (
	"encoding/json"
	"time"
	"github.com/lib/pq"
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

// Accounts holds multiple Account items.
type Accounts []Account