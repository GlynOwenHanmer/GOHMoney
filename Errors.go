package GOHMoney

import (
	"fmt"
	"time"
)

// BalanceDateOutOfAccountTimeRange is a type returned when the date of a Balance is not contained within the TimeRange of the Account that holds it.
// BalanceDate and AccountTimeRange fields are present and provide the exact detail of the timings that have discrepancies.
type BalanceDateOutOfAccountTimeRange struct {
	BalanceDate      time.Time
	AccountTimeRange TimeRange
}

// Error ensures that BalanceDateOutOfAccountTimeRange adheres to the error interface.
func (e BalanceDateOutOfAccountTimeRange) Error() string {
	return fmt.Sprintf("Balance Date is outside of Account Time Range.")
}
