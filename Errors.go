package GOHMoney

import (
	"bytes"
	"fmt"
	"time"
)

// AccountFieldError holds zero or more descriptions of things that are wrong with potential new Account items.
type AccountFieldError []string

// Error ensures that AccountFieldError adheres to the error interface.
func (e AccountFieldError) Error() string {
	var errorString bytes.Buffer
	errorString.WriteString("AccountFieldError: ")
	for i, field := range e {
		errorString.WriteString(field)
		if i < len(e) - 1 {
			errorString.WriteByte(' ')
		}
	}
	return string(errorString.String())
}

// Various error strings describing possible errors with potential new Account items.
const (
	EmptyNameError                   = "Empty name."
	ZeroDateOpenedError              = "No opened date given."
	ZeroValidDateClosedError         = "Closed date marked as valid but not set."
)

// BalanceDateOutOfAccountTimeRange is a type returned when the date of a Balance is not contained within the TimeRange of the Account that holds it.
// BalanceDate and AccountTimeRange fields are present and provide the exact detail of the timings that have discrepancies.
type BalanceDateOutOfAccountTimeRange struct {
	BalanceDate time.Time
	AccountTimeRange TimeRange
}

// Error ensures that BalanceDateOutOfAccountTimeRange adheres to the error interface.
func (e BalanceDateOutOfAccountTimeRange) Error() string {
	return fmt.Sprintf("Balance Date is outside of Account Time Range.")
}
